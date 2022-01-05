package interceptors

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/scribd/go-sdk/pkg/testing/testproto"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	golog "log"
	"net"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	sdktesting "github.com/scribd/go-sdk/pkg/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
)

type TestRecord struct {
	ID   int
	Name string
}

func TestDatabaseLoggingUnaryServerInterceptor(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	dbFile := "/tmp/test_db"
	defer os.Remove(dbFile)

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("Failed to open DB: %s", err)
	}
	defer db.Close()

	var (
		testRecordOne = TestRecord{ID: 1, Name: "test_name"}
		testRecordTwo = TestRecord{ID: 2, Name: "new_test_name"}
	)

	errors := db.Begin().
		CreateTable(TestRecord{}).
		Create(testRecordOne).
		Create(testRecordTwo).
		Commit().GetErrors()

	for _, err := range errors {
		t.Fatalf("Errors: %v", err)
	}

	var buffer bytes.Buffer

	l, err := getLogger("debug", &buffer)
	require.Nil(t, err)

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer([]grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
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

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx,
		"bufnet",
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

	dbFile := "/tmp/test_db"
	defer os.Remove(dbFile)

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatalf("Failed to open DB: %s", err)
	}
	defer db.Close()

	var (
		testRecordOne = TestRecord{ID: 1, Name: "test_name"}
		testRecordTwo = TestRecord{ID: 2, Name: "new_test_name"}
	)

	errors := db.Begin().
		CreateTable(TestRecord{}).
		Create(testRecordOne).
		Create(testRecordTwo).
		Commit().GetErrors()

	for _, err := range errors {
		t.Fatalf("Errors: %v", err)
	}

	var buffer bytes.Buffer

	l, err := getLogger("debug", &buffer)
	require.Nil(t, err)

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer([]grpc.ServerOption{
		grpc.ChainStreamInterceptor(
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

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx,
		"bufnet",
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

	var sql = (fields["sql"]).(map[string]interface{})

	assert.NotEmpty(t, sql["duration"])
	assert.NotEmpty(t, sql["affected_rows"])
	assert.NotEmpty(t, sql["file_location"])
	assert.NotNil(t, sql["values"])
}
