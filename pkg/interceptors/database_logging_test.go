package interceptors

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	golog "log"
	"net"
	"path"
	"testing"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	sdktesting "github.com/scribd/go-sdk/pkg/testing"
	"github.com/scribd/go-sdk/pkg/testing/testproto"
)

type TestRecord struct {
	ID   int
	Name string
}

func TestDatabaseLoggingUnaryServerInterceptor(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	tempDBPath := path.Join(t.TempDir(), "test_db")

	db, err := gorm.Open(sqlite.Open(tempDBPath))
	if err != nil {
		t.Fatalf("Failed to open DB: %s", err)
	}

	testRecordOne := TestRecord{ID: 1, Name: "test_name"}
	testRecordTwo := TestRecord{ID: 2, Name: "new_test_name"}

	if err = db.AutoMigrate(TestRecord{}); err != nil {
		t.Fatalf("Failed to migrate DB: %s", err)
	}

	err = db.Begin().
		Create(&testRecordOne).
		Create(&testRecordTwo).
		Commit().Error
	if err != nil {
		t.Fatalf("Failed to create record: %s", err)
	}

	var buffer bytes.Buffer

	l, err := getLogger("trace", &buffer)
	require.Nil(t, err)

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer([]grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			RequestIDUnaryServerInterceptor(),
			TracingUnaryServerInterceptor("test"),
			LoggerUnaryServerInterceptor(l),
			DatabaseUnaryServerInterceptor(db),
			DatabaseLoggingUnaryServerInterceptor(),
		),
	}...)

	testproto.RegisterTestServiceServer(s, sdktesting.NewTestService(db))

	go func() {
		if serveErr := s.Serve(lis); serveErr != nil {
			golog.Fatalf("Server exited with error: %v", serveErr)
		}
	}()

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := testproto.NewTestServiceClient(conn)
	_, err = client.Get(context.Background(), &testproto.GetRequest{})
	assert.Nil(t, err)

	// read first log entry
	var fieldsUnary map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(buffer.Bytes()))
	err = dec.Decode(&fieldsUnary)
	require.Nil(t, err)

	checkGormLoggerFields(t, fieldsUnary)
}

func TestDatabaseLoggingStreamServerInterceptors(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	tempDBPath := path.Join(t.TempDir(), "test_db")
	db, err := gorm.Open(sqlite.Open(tempDBPath))
	if err != nil {
		t.Fatalf("Failed to open DB: %s", err)
	}

	testRecordOne := TestRecord{ID: 1, Name: "test_name"}
	testRecordTwo := TestRecord{ID: 2, Name: "new_test_name"}

	if err = db.AutoMigrate(TestRecord{}); err != nil {
		t.Fatalf("Failed to migrate DB: %s", err)
	}

	err = db.Begin().
		Create(&testRecordOne).
		Create(&testRecordTwo).
		Commit().Error
	if err != nil {
		t.Fatalf("Failed to create record: %s", err)
	}

	var buffer bytes.Buffer

	l, err := getLogger("trace", &buffer)
	require.Nil(t, err)

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer([]grpc.ServerOption{
		grpc.ChainStreamInterceptor(
			RequestIDStreamServerInterceptor(),
			TracingStreamServerInterceptor("test"),
			LoggerStreamServerInterceptor(l),
			DatabaseStreamServerInterceptor(db),
			DatabaseLoggingStreamServerInterceptor(),
		),
	}...)

	testproto.RegisterTestServiceServer(s, sdktesting.NewTestService(db))

	go func() {
		if serveErr := s.Serve(lis); serveErr != nil {
			golog.Fatalf("Server exited with error: %v", serveErr)
		}
	}()

	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := testproto.NewTestServiceClient(conn)
	stream, err := client.GetList(context.Background(), &testproto.GetListRequest{})
	assert.Nil(t, err)

	require.NoError(t, stream.CloseSend(), "no error on close of stream")

	for {
		res := &testproto.GetResponse{}
		recvErr := stream.RecvMsg(res)
		if recvErr == io.EOF {
			break
		}
		require.NoError(t, recvErr, "no error on receive")
	}
	assert.Nil(t, err)

	// read first log entry
	var fieldsStream map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(buffer.Bytes()))
	err = dec.Decode(&fieldsStream)
	require.Nil(t, err)

	checkGormLoggerFields(t, fieldsStream)
}

func checkGormLoggerFields(t *testing.T, fields map[string]interface{}) {
	assert.NotEmpty(t, fields["sql"])

	dbFields, ok := (fields["sql"]).(map[string]interface{})
	assert.True(t, ok, "%s not found in log fields", "trace")
	assert.NotEmpty(t, dbFields)

	assert.Contains(t, dbFields, "duration")
	assert.Contains(t, dbFields, "affected_rows")
	assert.Contains(t, dbFields, "sql")
}
