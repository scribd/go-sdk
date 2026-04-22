package logger

import (
	"context"
	"fmt"
	"maps"

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
	maps.Copy(l.fields, fields)
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
	maps.Copy(fields, tags.Values())

	// Add sdklogger fields added until now.
	maps.Copy(fields, l.fields)

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
