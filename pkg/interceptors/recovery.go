package interceptors

import (
	"context"
	stdlog "log"
	"runtime/debug"
	"time"

	"github.com/getsentry/sentry-go"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdkloggercontext "github.com/scribd/go-sdk/pkg/context/logger"
)

// RecoveryUnaryServerInterceptor returns a unary server interceptor that recovers from panics,
// sends a sentry event, log in fatal level and halts the service.
// IMPORTANT: This interceptor should be the last one in the interceptor chain.
func RecoveryUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpcrecovery.UnaryServerInterceptor(recoveryOption...)
}

// RecoveryStreamServerInterceptor returns a streaming server interceptor that recovers from panics,
// sends a sentry event, log in fatal level and halts the service.
// IMPORTANT: This interceptor should be the last one in the interceptor chain.
func RecoveryStreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpcrecovery.StreamServerInterceptor(recoveryOption...)
}

var recoveryOption = []grpcrecovery.Option{
	grpcrecovery.WithRecoveryHandlerContext(func(ctx context.Context, rec interface{}) (err error) {
		sentry.CurrentHub().Recover(rec)
		sentry.Flush(time.Second * 5)

		l, err := sdkloggercontext.Extract(ctx)
		if err != nil {
			debug.PrintStack()
			stdlog.Printf("logger not found in context: %v\n", err)
			stdlog.Fatalf("grpc: panic error: %v", rec)
		}

		l.Fatalf("panic error: %v", rec)
		return status.Errorf(codes.Internal, "")
	}),
}
