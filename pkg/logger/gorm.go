package logger

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

const (
	gormLoggerTraceFieldKey = "sql"
	gormLoggerMsg           = "gorm DB Logger"
)

func NewGormLogger(l Logger) gormLogger {
	return gormLogger{l}
}

type gormLogger struct {
	logger Logger
}

func (g gormLogger) LogMode(logger.LogLevel) logger.Interface {
	// we ignore changes to the log level.
	// making changes to logger level causes a change to the root logger.
	return g
}

func (g gormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	g.logger.Infof(msg, args...)
}

func (g gormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	g.logger.Warnf(msg, args...)
}

func (g gormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	g.logger.WithError(fmt.Errorf(msg, args...)).Errorf(gormLoggerMsg)
}

func (g gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()

	l := g.logger.WithFields(Fields{
		gormLoggerTraceFieldKey: Fields{
			"duration":      time.Since(begin),
			"affected_rows": rows,
			"sql":           sql,
		},
	})

	if err != nil {
		l.WithError(err).Tracef(gormLoggerMsg)

		return
	}

	l.Tracef(gormLoggerMsg)
}

func (g gormLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	return sql, nil
}
