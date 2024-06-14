package interceptors

import (
	"context"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	sdkcontext "github.com/scribd/go-sdk/pkg/context/database"
)

// DatabaseUnaryServerInterceptor returns a unary server interceptor that adds gorm.DB to the context.
func DatabaseUnaryServerInterceptor(db *gorm.DB) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		instrumentedDB := db.WithContext(ctx)
		newCtx := sdkcontext.ToContext(ctx, instrumentedDB)

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
		instrumentedDB := db.WithContext(stream.Context())
		newCtx := sdkcontext.ToContext(stream.Context(), instrumentedDB)

		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}
