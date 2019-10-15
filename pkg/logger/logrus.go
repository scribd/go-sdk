package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type logrusLogEntry struct {
	entry *logrus.Entry
}

type logrusLogger struct {
	logger *logrus.Logger
}

const (
	// Key name for logging time of event
	fieldKeyTime = "timestamp"

	// Key name for logging event message
	fieldKeyMsg = "message"
)

func getFormatter(isJSON bool) logrus.Formatter {
	fieldMap := logrus.FieldMap{
		logrus.FieldKeyTime: fieldKeyTime,
		logrus.FieldKeyMsg:  fieldKeyMsg,
	}

	if isJSON {
		return &logrus.JSONFormatter{
			FieldMap: fieldMap,
		}
	}

	return &logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		FieldMap:               fieldMap,
	}
}

func newLogrus(config Config) (Logger, error) {
	logLevel := config.ConsoleLevel
	if logLevel == "" {
		logLevel = config.FileLevel
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	stdOutHandler := os.Stdout
	fileHandler := &lumberjack.Logger{
		Filename: fmt.Sprintf(config.FileLocation, config.Environment),
		MaxSize:  100,
		Compress: true,
		MaxAge:   28,
	}

	lLogger := &logrus.Logger{
		Hooks: make(logrus.LevelHooks),
		Level: level,
	}

	if config.ConsoleEnabled && config.FileEnabled {
		lLogger.SetOutput(io.MultiWriter(stdOutHandler, fileHandler))

		// Logrus can handle MultiWriter but not (easily) MultiFormat.
		// In this case the JSON format wins to ease the log processing.
		lLogger.SetFormatter(getFormatter(true))
	}

	if config.ConsoleEnabled && !config.FileEnabled {
		lLogger.SetOutput(stdOutHandler)
		lLogger.SetFormatter(getFormatter(config.ConsoleJSONFormat))
	}

	if !config.ConsoleEnabled && config.FileEnabled {
		lLogger.SetOutput(fileHandler)
		lLogger.SetFormatter(getFormatter(config.FileJSONFormat))
	}

	return &logrusLogger{
		logger: lLogger,
	}, nil
}

func (l *logrusLogger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Paniclf(format, args...)
}

func (l *logrusLogger) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.logger.WithFields(convertToLogrusFields(fields)),
	}
}

func (l *logrusLogEntry) Tracef(format string, args ...interface{}) {
	l.entry.Tracef(format, args...)
}

func (l *logrusLogEntry) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logrusLogEntry) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logrusLogEntry) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logrusLogEntry) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logrusLogEntry) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogEntry) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

func (l *logrusLogEntry) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.entry.WithFields(convertToLogrusFields(fields)),
	}
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, val := range fields {
		logrusFields[index] = val
	}
	return logrusFields
}
