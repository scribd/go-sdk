// logger includes a configurable logger that enforces Structured Logging,
// Environment aware formatters, File or Console output and de-facto
// standard logging levels.

package logger

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

	WithFields(keyValues Fields) Logger
}

// Config stores the configuration for the logger.
// For some loggers there can only be one level across writers, for such
// the level of Console is picked by default.
type Config struct {
	ConsoleEnabled    bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	Environment       string
	FileEnabled       bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
}

func NewLogger(config Config) (Logger, error) {
	return newLogrus(config)
}
