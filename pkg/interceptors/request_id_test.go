package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/scribd/go-sdk/pkg/context/requestid"
)

var (
	unaryInfo = &grpc.UnaryServerInfo{
		FullMethod: "TestService.UnaryMethod",
	}
)

func TestRequestIDUnaryServerInterceptor(t *testing.T) {
	testRequestID := newRequestID()

	tests := []struct {
		name    string
		set     func() context.Context
		handler func(ctx context.Context, req interface{}) (interface{}, error)
	}{
		{
			name: "without request ID",
			set: func() context.Context {
				return context.Background()
			},
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				requestID, err := requestid.Extract(ctx)
				assert.NoError(t, err)

				assert.NotEmpty(t, requestID)

				return "test", nil
			},
		},
		{
			name: "with request ID",
			set: func() context.Context {
				ctx := context.Background()
				md := metadata.Pairs(RequestIDKey, testRequestID)

				return metadata.NewIncomingContext(ctx, md)
			},
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				requestID, err := requestid.Extract(ctx)
				assert.NoError(t, err)

				assert.Equal(t, testRequestID, requestID)

				return "test", nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.set()

			_, err := RequestIDUnaryServerInterceptor()(ctx, "test", unaryInfo, tt.handler)
			assert.NoError(t, err)
		})
	}
}

type testServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *testServerStream) Context() context.Context {
	return ss.ctx
}

func (ss *testServerStream) SendMsg(m interface{}) error {
	return nil
}

func (ss *testServerStream) RecvMsg(m interface{}) error {
	return nil
}

var (
	streamInfo = &grpc.StreamServerInfo{
		FullMethod:     "TestService.StreamMethod",
		IsServerStream: true,
	}
)

func TestRequestIDStreamServerInterceptor(t *testing.T) {
	testRequestID := newRequestID()

	tests := []struct {
		name    string
		set     func() context.Context
		handler func(srv interface{}, stream grpc.ServerStream) error
	}{
		{
			name: "without request id",
			set: func() context.Context {
				return context.Background()
			},
			handler: func(srv interface{}, stream grpc.ServerStream) error {
				requestID, err := requestid.Extract(stream.Context())
				assert.NoError(t, err)

				assert.NotEmpty(t, requestID)

				return nil
			},
		},
		{
			name: "with request id",
			set: func() context.Context {
				ctx := context.Background()
				md := metadata.Pairs(RequestIDKey, testRequestID)

				return metadata.NewIncomingContext(ctx, md)
			},
			handler: func(srv interface{}, stream grpc.ServerStream) error {
				requestID, err := requestid.Extract(stream.Context())
				assert.NoError(t, err)

				assert.Equal(t, testRequestID, requestID)

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.set()

			testStream := &testServerStream{ctx: ctx}
			err := RequestIDStreamServerInterceptor()(struct{}{}, testStream, streamInfo, tt.handler)
			assert.NoError(t, err)
		})
	}
}
