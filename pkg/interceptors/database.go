package interceptors

import (
	"context"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"

	sdkcontext "github.com/scribd/go-sdk/pkg/context/database"
	sdkinstrumentation "github.com/scribd/go-sdk/pkg/instrumentation"
)

// DatabaseUnaryServerInterceptor returns a unary server interceptor that adds gorm.DB to the context.
func DatabaseUnaryServerInterceptor(db *gorm.DB) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		newDb := sdkinstrumentation.TraceDatabase(ctx, db)
		newCtx := sdkcontext.ToContext(ctx, newDb)

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
		newDb := sdkinstrumentation.TraceDatabase(stream.Context(), db)
		newCtx := sdkcontext.ToContext(stream.Context(), newDb)
		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}
