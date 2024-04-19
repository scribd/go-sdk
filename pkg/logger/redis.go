package logger

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/scribd/go-sdk/pkg/instrumentation"
)

type (
	RedisLogger struct {
		logger Logger
	}
)

func NewRedisLogger(l Logger) *RedisLogger {
	return &RedisLogger{l}
}

func (r *RedisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	logContext := instrumentation.TraceLogs(ctx)

	r.logger.WithFields(Fields{
		"dd": Fields{
			"trace_id": logContext.TraceID,
			"span_id":  logContext.SpanID,
		},
	}).Errorf(format, v...)
}

func SetRedisLogger(logger Logger) {
	redis.SetLogger(NewRedisLogger(logger))
}
