package interceptors

import (
	"context"
	"fmt"
	"path"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc "google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	sdkcontext "git.lo/microservices/sdk/go-sdk/pkg/context/logger"
	sdklogger "git.lo/microservices/sdk/go-sdk/pkg/logger"
)

// LoggerUnaryServerInterceptor returns a unary server interceptors that adds the sdklogger.Logger to the context.
func LoggerUnaryServerInterceptor(logger sdklogger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()
		newCtx := newLoggerForCall(ctx, logger, info.FullMethod, startTime)

		resp, err := handler(newCtx, req)

		log(newCtx, err, startTime)

		return resp, err
	}
}

// LoggerStreamServerInterceptor returns a streaming server interceptor that adds the sdklogger.Logger to the context.
func LoggerStreamServerInterceptor(logger sdklogger.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		startTime := time.Now()
		newCtx := newLoggerForCall(stream.Context(), logger, info.FullMethod, startTime)

		wrapped := grpcmiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		err := handler(srv, wrapped)

		log(newCtx, err, startTime)

		return err
	}
}

func newLoggerForCall(
	ctx context.Context,
	logger sdklogger.Logger,
	method string,
	startTime time.Time,
) context.Context {
	callLog := logger.WithFields(
		sdklogger.Fields{
			"system":          "grpc",
			"span.kind":       "server",
			"grpc.service":    path.Dir(method)[1:],
			"grpc.method":     path.Base(method),
			"grpc.start_time": startTime.Format(time.RFC3339),
		})

	if d, ok := ctx.Deadline(); ok {
		callLog = callLog.WithFields(
			sdklogger.Fields{
				"grpc.request.deadline": d.Format(time.RFC3339),
			})
	}

	return sdkcontext.ToContext(ctx, callLog)
}

// errorToCode function determines the error code of an error.
func errorToCode(err error) grpccodes.Code {
	return grpcstatus.Code(err)
}

func log(ctx context.Context, err error, startTime time.Time) {
	code := errorToCode(err)
	level := grpcCodeToLevel(code)
	fields := sdklogger.Fields{
		"grpc.code":    code.String(),
		"grpc.time_ms": float32(time.Since(startTime).Nanoseconds()/1000) / 1000,
	}

	msg := fmt.Sprintf("finished gRPC call with code %s", code.String())
	l, extractErr := sdkcontext.Extract(ctx)
	if extractErr == nil {
		l = l.WithFields(fields)

		switch level {
		case sdklogger.Debug:
			l.Debugf(msg)
		case sdklogger.Info:
			l.Infof(msg)
		case sdklogger.Warn:
			l.Warnf(msg)
		case sdklogger.Error:
			l.Errorf(msg)
		case sdklogger.Fatal:
			l.Fatalf(msg)
		case sdklogger.Panic:
			l.Panicf(msg)
		}
	}
}

// grpcCodeToLevel is the default implementation of gRPC return grpccodes to log levels.
func grpcCodeToLevel(code grpccodes.Code) sdklogger.Level {
	switch code {
	case grpccodes.OK:
		return sdklogger.Info
	case grpccodes.Canceled:
		return sdklogger.Info
	case grpccodes.Unknown:
		return sdklogger.Error
	case grpccodes.InvalidArgument:
		return sdklogger.Info
	case grpccodes.DeadlineExceeded:
		return sdklogger.Warn
	case grpccodes.NotFound:
		return sdklogger.Info
	case grpccodes.AlreadyExists:
		return sdklogger.Info
	case grpccodes.PermissionDenied:
		return sdklogger.Warn
	case grpccodes.Unauthenticated:
		return sdklogger.Info
	case grpccodes.ResourceExhausted:
		return sdklogger.Warn
	case grpccodes.FailedPrecondition:
		return sdklogger.Warn
	case grpccodes.Aborted:
		return sdklogger.Warn
	case grpccodes.OutOfRange:
		return sdklogger.Warn
	case grpccodes.Unimplemented:
		return sdklogger.Error
	case grpccodes.Internal:
		return sdklogger.Error
	case grpccodes.Unavailable:
		return sdklogger.Warn
	case grpccodes.DataLoss:
		return sdklogger.Error
	default:
		return sdklogger.Error
	}
}
