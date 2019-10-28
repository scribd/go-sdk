// logger includes a configurable logger that enforces Structured Logging,
// formatters, File or Console output and de-facto standard logging levels.

package logger

import (
	"bytes"
)

type Level string

const (
	// Panic level, highest level of severity. Logs and then calls
	// panic with the message passed to Debug, Info, ...
	Panic Level = "panic"
	// Fatal level. Logs and then calls `logger.Exit(1)`. It will
	// exit even if the logging level is set to Panic.
	Fatal Level = "fatal"
	// Error level. Logs. Used for errors that should definitely be
	// noted. Commonly used for hooks to send errors to an error
	// tracking service.
	Error Level = "error"
	// Warn level. Non-critical entries that deserve eyes.
	Warn Level = "warn"
	// Info level. General operational entries about what's going on
	// inside the application.
	Info Level = "info"
	// Debug level. Usually only enabled when debugging. Very
	// verbose logging.
	Debug Level = "debug"
	// Trace level. Designates finer-grained informational events
	// than the Debug.
	Trace Level = "trace"
)

// Logger is the interface that defines the API/contract exposed by the
// SDK Logger.
type Logger interface {
	// Panicf logs a message at level Panic.
	Panicf(format string, args ...interface{})
	// Fatalf logs a message at level Fatal.
	Fatalf(format string, args ...interface{})
	// Errorf logs a message at level Error.
	Errorf(format string, args ...interface{})
	// Warnf logs a message at level Warning.
	Warnf(format string, args ...interface{})
	// Infof logs a message at level Info.
	Infof(format string, args ...interface{})
	// Debugf logs a message at level Debug.
	Debugf(format string, args ...interface{})
	// Trace logs a message at level Trace.
	Tracef(format string, args ...interface{})
	// WithFields creates an entry from the logger and adds multiple
	// fields to it. This is simply a helper for `WithField`,
	// invoking it once for each field.
	//
	// Note that it doesn't log until you call Debug, Print, Info,
	// Warn, Fatal or Panic on the Entry it returns.
	WithFields(keyValues Fields) Logger
}

// NewLogger returns a Logger instance with the given configuration.
func NewLogger(config *Config) (Logger, error) {
	return newLogrus(config)
}

// NewTestLogger returns a Logger instance that will write into the bytes buffer
// passed as parameter.
// NewTestLogger is recommended only for testing.
func NewTestLogger(config *Config, out *bytes.Buffer) (Logger, error) {
	return newTestLogrus(config, out)
}
