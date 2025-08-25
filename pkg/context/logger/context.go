package logger

import (
	"context"
	"fmt"

	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	sdklogger "github.com/scribd/go-sdk/pkg/logger"
)

type ctxLoggerMarker struct{}

type ctxLogger struct {
	logger sdklogger.Logger
	fields sdklogger.Fields
}

var (
	ctxLoggerKey = &ctxLoggerMarker{}
)

// AddFields adds sdklogger fields to the logger.
func AddFields(ctx context.Context, fields sdklogger.Fields) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	for k, v := range fields {
		l.fields[k] = v
	}
}

// Extract takes the call-scoped sdklogger.Logger from the context.
// If the ctxLogger  wasn't used, an error is returned.
func Extract(ctx context.Context) (sdklogger.Logger, error) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return nil, fmt.Errorf("unable to get the logger")
	}

	fields := sdklogger.Fields{}

	// Add grpcctxtags tags metadata until now.
	tags := grpcctxtags.Extract(ctx)
	for k, v := range tags.Values() {
		fields[k] = v
	}

	// Add sdklogger fields added until now.
	for k, v := range l.fields {
		fields[k] = v
	}

	return l.logger.WithFields(fields), nil
}

// ToContext adds the sdklogger.Logger to the context for extraction later.
// Returning the new context that has been created.
func ToContext(ctx context.Context, logger sdklogger.Logger) context.Context {
	l := &ctxLogger{
		logger: logger,
		fields: sdklogger.Fields{},
	}
	return context.WithValue(ctx, ctxLoggerKey, l)
}
