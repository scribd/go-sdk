package interceptors

import (
	"context"

	"github.com/google/uuid"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc "google.golang.org/grpc"

	sdkcontext "git.lo/microservices/sdk/go-sdk/pkg/context/requestid"
)

// RequestIDUnaryServerInterceptor returns a unary server interceptors that adds
// sdkrequestid.RequestID to the context.
func RequestIDUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if uuid, err := uuid.NewRandom(); err != nil {
			requestID := uuid.String()
			newCtx := sdkcontext.ToContext(ctx, requestID)
			return handler(newCtx, req)
		}
		return handler(ctx, req)
	}
}

// RequestIDStreamServerInterceptor returns a streaming server interceptor that adds
// sdkrequestid.RequestID to the context.
func RequestIDStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if uuid, err := uuid.NewRandom(); err != nil {
			requestID := uuid.String()
			newCtx := sdkcontext.ToContext(stream.Context(), requestID)
			wrapped := grpcmiddleware.WrapServerStream(stream)
			wrapped.WrappedContext = newCtx

			return handler(srv, wrapped)
		}

		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = stream.Context()

		return handler(srv, wrapped)
	}
}
