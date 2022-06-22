package kafka

import (
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/scribd/go-sdk/pkg/logger"
)

type (
	KafkaLogger struct {
		logger  logger.Logger
		levelFn func() kgo.LogLevel
	}
)

func NewKafkaLogger(l logger.Logger, opts ...Opt) *KafkaLogger {
	kafkaLogger := &KafkaLogger{
		logger: l,
		levelFn: func() kgo.LogLevel {
			return kgo.LogLevelError
		},
	}

	for _, opt := range opts {
		opt.apply(kafkaLogger)
	}

	return kafkaLogger
}

// Opt applies options to the kafka logger.
type Opt interface {
	apply(*KafkaLogger)
}

type opt struct{ fn func(*KafkaLogger) }

func (o opt) apply(l *KafkaLogger) { o.fn(l) }

// WithLevelFn sets a function that can dynamically change the log level.
func WithLevelFn(fn func() kgo.LogLevel) Opt {
	return opt{func(l *KafkaLogger) { l.levelFn = fn }}
}

// WithLevel sets a static level for the kgo.Logger Level function.
func WithLevel(level logger.Level) Opt {
	kgoLevel := kgo.LogLevelNone

	switch level {
	case logger.Panic, logger.Error, logger.Fatal:
		kgoLevel = kgo.LogLevelError
	case logger.Warn:
		kgoLevel = kgo.LogLevelWarn
	case logger.Info:
		kgoLevel = kgo.LogLevelInfo
	case logger.Trace, logger.Debug:
		kgoLevel = kgo.LogLevelDebug
	}

	return WithLevelFn(func() kgo.LogLevel { return kgoLevel })
}

// Level is for the kgo.Logger interface.
func (l *KafkaLogger) Level() kgo.LogLevel {
	return l.levelFn()
}

func (l *KafkaLogger) Log(level kgo.LogLevel, msg string, keyvals ...interface{}) {
	fields := logger.Fields{}
	for i := 0; i < len(keyvals); i += 2 {
		k, v := keyvals[i], keyvals[i+1]

		kStr, ok := k.(string)
		if !ok {
			continue
		}

		fields[kStr] = v
	}

	logEntry := l.logger.WithFields(fields)

	switch level {
	case kgo.LogLevelError:
		logEntry.Errorf(msg)
	case kgo.LogLevelWarn:
		logEntry.Warnf(msg)
	case kgo.LogLevelInfo:
		logEntry.Infof(msg)
	case kgo.LogLevelDebug:
		logEntry.Debugf(msg)
	default:
		// do nothing
	}
}
