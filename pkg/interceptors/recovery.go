package interceptors

import (
	"context"
	"fmt"
	stdlog "log"
	"runtime/debug"

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
		l, err := sdkloggercontext.Extract(ctx)
		if err != nil {
			debug.PrintStack()
			stdlog.Printf("logger not found in context: %v\n", err)
			stdlog.Fatalf("grpc: panic error: %v", rec)
		}

		l.WithError(fmt.Errorf("%v", rec)).Fatalf("panic error: %v", rec)
		return status.Errorf(codes.Internal, "")
	}),
}
