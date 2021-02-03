package interceptors

import (
	"context"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc "google.golang.org/grpc"

	sdkcontext "github.com/scribd/go-sdk/pkg/context/metrics"
	sdkmetrics "github.com/scribd/go-sdk/pkg/metrics"
)

// MetricsUnaryServerInterceptor returns a unary server interceptors that adds sdkmetrics.Metrics to the context.
func MetricsUnaryServerInterceptor(metrics sdkmetrics.Metrics) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		newCtx := sdkcontext.ToContext(ctx, metrics)
		return handler(newCtx, req)
	}
}

// MetricsStreamServerInterceptor returns a streaming server interceptor that adds sdkmetrics.Metrics to the context.
func MetricsStreamServerInterceptor(metrics sdkmetrics.Metrics) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		newCtx := sdkcontext.ToContext(stream.Context(), metrics)
		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}
