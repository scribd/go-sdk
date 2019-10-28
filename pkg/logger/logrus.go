package logger

import (
	"bytes"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	// Key name for logging time of event.
	fieldKeyTime = "timestamp"

	// Key name for logging event message.
	fieldKeyMsg = "message"
)

const (
	// fileMaxSize is the maximum size in megabytes of the log file
	// before it gets rotated. It defaults to 100 megabytes.
	fileMaxSize = 100
	// fileWillCompress determines if the rotated log files should
	// be compressed using gzip.
	fileWillCompress = true
	// fileMaxAge is the maximum number of days to retain old log
	// files based on the timestamp encoded in their filename. Note
	// that a day is defined as 24 hours and may not exactly
	// correspond to calendar days due to daylight savings, leap
	// seconds, etc. The default is not to remove old log files
	// based on age.
	fileMaxAge = 28
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

func newLogrusLogger(config *Config) (*logrus.Logger, error) {
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
		Filename: path.Join(config.FileLocation, config.FileName),
		MaxSize:  fileMaxSize,
		Compress: fileWillCompress,
		MaxAge:   fileMaxAge,
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

	return lLogger, nil
}

func newTestLogrusLogger(config *Config, out *bytes.Buffer) (*logrus.Logger, error) {
	lLogger, err := newLogrusLogger(config)
	lLogger.Out = out

	return lLogger, err
}

// An entry is the final or intermediate Logrus logging entry. It
// contains all the fields passed with WithField{,s}. It's finally
// logged when Trace, Debug, Info, Warn, Error, Fatal or Panic is called
// on it.
//
// logrusEntry implements the `Logger` interface.
type logrusLogEntry struct {
	entry *logrus.Entry
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
