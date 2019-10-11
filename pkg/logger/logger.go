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

