package interceptors

import (
	"context"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	sdkdatabasecontext "github.com/scribd/go-sdk/pkg/context/database"
	sdkloggercontext "github.com/scribd/go-sdk/pkg/context/logger"
	sdklogger "github.com/scribd/go-sdk/pkg/logger"
)

// DatabaseLoggingUnaryServerInterceptor returns a unary server interceptor.
// This interceptor extracts the logger from the request context, attaches it to
// gorm's database connection pool and logs the database queries with their
// meta-information using the logger.
func DatabaseLoggingUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		db, err := sdkdatabasecontext.Extract(ctx)
		if err != nil {
			return nil, err
		}

		l, err := sdkloggercontext.Extract(ctx)
		if err != nil {
			return nil, err
		}

		newDB := db.Session(&gorm.Session{
			Logger: sdklogger.NewGormLogger(l),
			NewDB:  true,
		})

		newCtx := sdkdatabasecontext.ToContext(ctx, newDB)

		return handler(newCtx, req)
	}
}

// DatabaseLoggingStreamServerInterceptor returns a streaming server interceptor.
// This interceptor extracts the logger from the request context, attaches it to
// gorm's database connection pool and logs the database queries with their
// meta-information using the logger.
func DatabaseLoggingStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		_ *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		db, err := sdkdatabasecontext.Extract(stream.Context())
		if err != nil {
			return err
		}

		l, err := sdkloggercontext.Extract(stream.Context())
		if err != nil {
			return err
		}

		newDB := db.Session(&gorm.Session{
			Logger: sdklogger.NewGormLogger(l),
			NewDB:  true,
		})

		newCtx := sdkdatabasecontext.ToContext(stream.Context(), newDB)
		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}
