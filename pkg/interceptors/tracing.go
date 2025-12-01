package interceptors

import (
	"fmt"

	grpctrace "github.com/DataDog/dd-trace-go/contrib/google.golang.org/grpc/v2"
	"google.golang.org/grpc"
)

const (
	datadogServiceServerSuffix = "grpc"
	datadogServiceClientSuffix = "grpc-client"
)

// TracingUnaryServerInterceptor returns a unary server interceptor that will
// trace requests to the given gRPC server.
func TracingUnaryServerInterceptor(applicationName string) grpc.UnaryServerInterceptor {
	serviceName := fmt.Sprintf("%s-%s", applicationName, datadogServiceServerSuffix)
	return grpctrace.UnaryServerInterceptor(grpctrace.WithService(serviceName))
}

// TracingStreamServerInterceptor returns a stream server interceptor that will
// trace streaming requests to the given gRPC server.
func TracingStreamServerInterceptor(applicationName string) grpc.StreamServerInterceptor {
	serviceName := fmt.Sprintf("%s-%s", applicationName, datadogServiceServerSuffix)
	return grpctrace.StreamServerInterceptor(grpctrace.WithService(serviceName))
}

// TracingUnaryClientInterceptor returns a unary client interceptor that will
// trace requests performed by gRPC client.
func TracingUnaryClientInterceptor(applicationName string) grpc.UnaryClientInterceptor {
	serviceName := fmt.Sprintf("%s-%s", applicationName, datadogServiceClientSuffix)
	return grpctrace.UnaryClientInterceptor(grpctrace.WithService(serviceName))
}

// TracingStreamClientInterceptor returns a stream server interceptor that will
// trace streaming requests performed by gRPC client.
func TracingStreamClientInterceptor(applicationName string) grpc.StreamClientInterceptor {
	serviceName := fmt.Sprintf("%s-%s", applicationName, datadogServiceClientSuffix)
	return grpctrace.StreamClientInterceptor(grpctrace.WithService(serviceName))
}
