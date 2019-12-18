package tracking

import (
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

// NewSentryHook creates a hook to be added to an instance of logger
// and initializes the raven client.
// This method sets the timeout to 100 milliseconds.
func NewSentryHook(config *Config) (*logrus_sentry.SentryHook, error) {
	hook, err := logrus_sentry.NewSentryHook(
		config.SentryDSN,
		[]logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})

	if err != nil {
		return hook, err
	}

	hook.Timeout = time.Duration(config.SentryTimeout) * time.Millisecond
	hook.StacktraceConfiguration.Enable = true

	return hook, err
}
