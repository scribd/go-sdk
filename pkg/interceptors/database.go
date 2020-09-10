package interceptors

import (
	"context"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	gorm "github.com/jinzhu/gorm"
	grpc "google.golang.org/grpc"

	sdkcontext "git.lo/microservices/sdk/go-sdk/pkg/context/database"
)

// DatabaseUnaryServerInterceptor returns a unary server interceptor that adds gorm.DB to the context.
func DatabaseUnaryServerInterceptor(db *gorm.DB) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		newCtx := sdkcontext.ToContext(ctx, db)
		return handler(newCtx, req)
	}
}

// DatabaseStreamServerInterceptor returns a streaming server interceptor that adds gorm.DB to the context.
func DatabaseStreamServerInterceptor(db *gorm.DB) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		newCtx := sdkcontext.ToContext(stream.Context(), db)
		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}
