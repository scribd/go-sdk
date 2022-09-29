package interceptors

import (
	"fmt"

	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

const (
	datadogServiceServerSuffix = "grpc"
	datadogServiceClientSuffix = "grpc-client"
)

// TracingUnaryServerInterceptor returns a unary server interceptor that will
// trace requests to the given gRPC server.
func TracingUnaryServerInterceptor(applicationName string) grpc.UnaryServerInterceptor {
	serviceName := fmt.Sprintf("%s-%s", applicationName, datadogServiceServerSuffix)
	return grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName(serviceName))
}

// TracingStreamServerInterceptor returns a stream server interceptor that will
// trace streaming requests to the given gRPC server.
func TracingStreamServerInterceptor(applicationName string) grpc.StreamServerInterceptor {
	serviceName := fmt.Sprintf("%s-%s", applicationName, datadogServiceServerSuffix)
	return grpctrace.StreamServerInterceptor(grpctrace.WithServiceName(serviceName))
}

// TracingUnaryClientInterceptor returns a unary client interceptor that will
// trace requests performed by gRPC client.
func TracingUnaryClientInterceptor(applicationName string) grpc.UnaryClientInterceptor {
	serviceName := fmt.Sprintf("%s-%s", applicationName, datadogServiceClientSuffix)
	return grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName(serviceName))
}

// TracingStreamClientInterceptor returns a stream server interceptor that will
// trace streaming requests performed by gRPC client.
func TracingStreamClientInterceptor(applicationName string) grpc.StreamClientInterceptor {
	serviceName := fmt.Sprintf("%s-%s", applicationName, datadogServiceClientSuffix)
	return grpctrace.StreamClientInterceptor(grpctrace.WithServiceName(serviceName))
}
