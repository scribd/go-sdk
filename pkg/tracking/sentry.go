package tracking

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

var (
	levelsMap = map[logrus.Level]sentry.Level{
		logrus.PanicLevel: sentry.LevelFatal,
		logrus.FatalLevel: sentry.LevelFatal,
		logrus.ErrorLevel: sentry.LevelError,
		logrus.WarnLevel:  sentry.LevelWarning,
		logrus.InfoLevel:  sentry.LevelInfo,
		logrus.DebugLevel: sentry.LevelDebug,
		logrus.TraceLevel: sentry.LevelDebug,
	}
)

// Hook is a service hook for the Logrus logger.
//
// It's used for sending errors and messages to Sentry on specific
// log levels. It wraps a default Sentry client.
type Hook struct {
	levels      []logrus.Level
	tags        map[string]string
	release     string
	environment string
}

// Levels returns the list of Logrus levels for which this hook is configured
// to report errors.
func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}

// Fire uses the configured Sentry client to report the given Logrus Entry as
// Sentry Event.
func (hook *Hook) Fire(entry *logrus.Entry) error {
	entryError, ok := entry.Data[logrus.ErrorKey].(error)
	if !ok || entryError == nil {
		return nil
	}

	stacktrace := sentry.ExtractStacktrace(entryError)
	if stacktrace == nil {
		stacktrace = sentry.NewStacktrace()
	}

	exceptions := []sentry.Exception{{
		Type:       entry.Message,
		Value:      entryError.Error(),
		Stacktrace: stacktrace,
	}}

	event := &sentry.Event{
		Level:       levelsMap[entry.Level],
		Message:     entry.Message,
		Extra:       map[string]interface{}(entry.Data),
		Tags:        hook.tags,
		Environment: hook.environment,
		Release:     hook.release,
		Exception:   exceptions,
	}

	sentry.CaptureEvent(event)

	if entry.Level <= logrus.FatalLevel {
		sentry.Flush(time.Second * 5)
	}

	return nil
}

// SetTags sets the given map of tags to every Sentry Event handled by this hook.
func (hook *Hook) SetTags(tags map[string]string) {
	hook.tags = tags
}

// AddTag add a pair (key, value) in the map of tags attached to every
// Sentry Event handled by this hook.
func (hook *Hook) AddTag(key, value string) {
	hook.tags[key] = value
}

// SetRelease sets the release that every Sentry Event handled by this
// hook refers to.
func (hook *Hook) SetRelease(release string) {
	hook.release = release
}

// SetEnvironment sets the environment that every Sentry Event handled by this
// hook refers to.
func (hook *Hook) SetEnvironment(environment string) {
	hook.environment = environment
}

// NewSentryHook creates a hook to be added to an instance of logger
// and initializes Sentry in package level. SentryHook will be triggered
// on Panic, Fatal and Error levels.
func NewSentryHook(config *Config) (*Hook, error) {
	if err := sentry.Init(sentry.ClientOptions{
		// The DSN to use. If the DSN is not set, the client is effectively disabled.
		Dsn: config.SentryDSN,
		// In debug mode, the debug information is printed to stdout to help you understand what
		// sentry is doing.
		Debug: false,
		// Configures whether SDK should generate and attach stacktraces to pure capture message calls.
		AttachStacktrace: true,
		// The sample rate for event submission (0.0 - 1.0, defaults to 1.0)
		SampleRate: 1.0,
		// The server name to be reported.
		ServerName: config.serverName,
		// The release to be sent with events.
		Release: config.release,
		// The environment to be sent with events.
		Environment: config.environment,
	}); err != nil {
		return nil, fmt.Errorf("initializing sentry. err: %w", err)
	}

	levels := []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
	}

	return &Hook{
		levels: levels,
		tags:   map[string]string{},
	}, nil
}
