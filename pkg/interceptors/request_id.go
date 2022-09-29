package interceptors

import (
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc/metadata"

	"github.com/scribd/go-sdk/pkg/context/requestid"

	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// RequestIDKey is metadata key name for request ID
var RequestIDKey = "x-request-id"

func RequestIDUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		requestID := handleRequestID(ctx)

		newCtx := requestid.ToContext(ctx, requestID)

		return handler(newCtx, req)
	}
}

func RequestIDStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := stream.Context()
		requestID := handleRequestID(ctx)

		newCtx := requestid.ToContext(ctx, requestID)

		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}

func handleRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return newRequestID()
	}

	header, ok := md[RequestIDKey]
	if !ok || len(header) == 0 {
		return newRequestID()
	}

	requestID := header[0]
	if requestID == "" {
		return newRequestID()
	}

	return requestID
}

func newRequestID() string {
	var uuidString string
	if s, err := uuid.NewRandom(); err == nil {
		uuidString = s.String()
	}

	return uuidString
}
