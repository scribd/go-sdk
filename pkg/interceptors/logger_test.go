package interceptors

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	golog "log"
	"net"
	"testing"

	grpc_testing "github.com/grpc-ecosystem/go-grpc-middleware/testing"
	mwitkow_testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"

	sdklogger "github.com/scribd/go-sdk/pkg/logger"
)

var (
	goodPing = &mwitkow_testproto.PingRequest{Value: "something", SleepTimeMs: 9999}
)

const bufSize = 1024 * 1024

func TestLoggerUnaryServerInterceptors(t *testing.T) {
	var buffer bytes.Buffer

	l, err := getLogger("info", &buffer)
	require.Nil(t, err)

	mt := mocktracer.Start()
	defer mt.Stop()

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer([]grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			TracingUnaryServerInterceptor("test"),
			RequestIDUnaryServerInterceptor(),
			LoggerUnaryServerInterceptor(l),
		),
	}...)
	mwitkow_testproto.RegisterTestServiceServer(s, &grpc_testing.TestPingService{T: t})
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

	client := mwitkow_testproto.NewTestServiceClient(conn)
	_, err = client.Ping(context.Background(), &mwitkow_testproto.PingRequest{Value: "test"})
	assert.Nil(t, err)

	var fieldsUnary map[string]interface{}
	err = json.Unmarshal(buffer.Bytes(), &fieldsUnary)
	require.Nil(t, err)

	checkLoggerFields(t, fieldsUnary)
}

func TestLoggerStreamServerInterceptors(t *testing.T) {
	var buffer bytes.Buffer

	l, err := getLogger("info", &buffer)
	require.Nil(t, err)

	mt := mocktracer.Start()
	defer mt.Stop()

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer([]grpc.ServerOption{
		grpc.ChainStreamInterceptor(
			TracingStreamServerInterceptor("test"),
			RequestIDStreamServerInterceptor(),
			LoggerStreamServerInterceptor(l),
		),
	}...)

	mwitkow_testproto.RegisterTestServiceServer(s, &grpc_testing.TestPingService{T: t})

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

	client := mwitkow_testproto.NewTestServiceClient(conn)
	stream, err := client.PingStream(context.Background())
	assert.Nil(t, err)

	require.NoError(t, stream.Send(goodPing), "sending must succeed")
	require.NoError(t, stream.CloseSend(), "no error on close of stream")

	for {
		pong := &mwitkow_testproto.PingResponse{}
		recvErr := stream.RecvMsg(pong)
		if recvErr == io.EOF {
			break
		}
		require.NoError(t, recvErr, "no error on receive")
	}
	assert.Nil(t, err)

	var fieldsStream map[string]interface{}
	err = json.Unmarshal(buffer.Bytes(), &fieldsStream)
	require.Nil(t, err)

	checkLoggerFields(t, fieldsStream)
}

func getLogger(logLevel string, buf *bytes.Buffer) (sdklogger.Logger, error) {
	config := &sdklogger.Config{
		ConsoleEnabled:    true,
		ConsoleJSONFormat: true,
		ConsoleLevel:      logLevel,
		FileEnabled:       false,
	}

	return sdklogger.NewBuilder(config).BuildTestLogger(buf)
}

func checkLoggerFields(t *testing.T, fields map[string]interface{}) {
	assert.NotEmpty(t, fields["message"])
	assert.Equal(t, "info", fields["level"])
	assert.NotEmpty(t, fields["timestamp"])

	assert.NotEmpty(t, fields["span.kind"])
	assert.NotEmpty(t, fields["grpc.service"])
	assert.NotEmpty(t, fields["grpc.method"])
	assert.NotEmpty(t, fields["grpc.start_time"])
	assert.NotEmpty(t, fields["grpc.code"])
	assert.NotEmpty(t, fields["grpc.time_ms"])
	assert.NotEmpty(t, fields["grpc.request_id"])

	var dd = (fields["dd"]).(map[string]interface{})

	assert.NotEmpty(t, dd["trace_id"])
	assert.NotEmpty(t, dd["span_id"])
}
