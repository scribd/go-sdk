package interceptors

import (
	grpc "google.golang.org/grpc"

	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

// TracingUnaryServerInterceptor returns a unary server interceptors that will
// trace requests to the given gRPC server.
func TracingUnaryServerInterceptor(applicationName string) grpc.UnaryServerInterceptor {
	return grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName(applicationName))
}

// TracingStreamServerInterceptor returns a stream server interceptors that will
// trace streaming requests to the given gRPC server.
func TracingStreamServerInterceptor(applicationName string) grpc.StreamServerInterceptor {
	return grpctrace.StreamServerInterceptor(grpctrace.WithServiceName(applicationName))
}
