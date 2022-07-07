package configuration

import (
	"github.com/scribd/go-sdk/pkg/configuration/builder"
)

// default filenames for configurables.
const (
	databaseDefaultFileName        = "database"
	pubsubDefaultFileName          = "pubsub"
	serverDefaultFileName          = "server"
	loggerDefaultFileName          = "logger"
	trackingDefaultFileName        = "sentry"
	instrumentationDefaultFileName = "datadog"
)

// Configurable type indicates apps to be fetched.
type Configurable func(builder.Builder, *Configuration) error

var (
	// DatabaseConf is a Configurable for database pointing at default directory.
	DatabaseConf Configurable = DatabaseConfWithFilename(databaseDefaultFileName)
	// ServerConf is a Configurable that fetches server pointing at default directory.
	ServerConf Configurable = ServerConfWithFilename(serverDefaultFileName)
	// LoggerConf is a Configurable that fetches logger pointing at default directory.
	LoggerConf Configurable = LoggerConfWithFilename(loggerDefaultFileName)
	// TrackingConf is a Configurable that fetches tracking pointing at default directory.
	TrackingConf Configurable = TrackingConfWithFilename(trackingDefaultFileName)
	// PubSubConf is a Configurable that fetches pubsub pointing at default directory.
	PubSubConf Configurable = PubsubConfWithFilename(pubsubDefaultFileName)
	// InstrumentationConf is a Configurable that fetches instrumentation pointing at default directory.
	InstrumentationConf Configurable = InstrumentationConfWithFilename(instrumentationDefaultFileName)
)

// DatabaseConfWithFilename returns database configurable based on filename.
func DatabaseConfWithFilename(filename string) Configurable {
	return func(b builder.Builder, conf *Configuration) error {
		b.SetConfigName(filename)

		return conf.Database.FetchConfig(b)
	}
}

// ServerConfWithFilename returns server configurable based on filename.
func ServerConfWithFilename(filename string) Configurable {
	return func(b builder.Builder, conf *Configuration) error {
		b.SetConfigName(filename)

		return conf.Server.FetchConfig(b)
	}
}

// TrackingConfWithFilename returns tracking based on filename.
func TrackingConfWithFilename(filename string) Configurable {
	return func(b builder.Builder, conf *Configuration) error {
		b.SetConfigName(filename)

		return conf.Tracking.FetchConfig(b)
	}
}

// InstrumentationConfWithFilename returns Instrumentation based on filename.
func InstrumentationConfWithFilename(filename string) Configurable {
	return func(b builder.Builder, conf *Configuration) error {
		b.SetConfigName(filename)

		return conf.Instrumentation.FetchConfig(b)
	}
}

// PubsubConfWithFilename returns pubsub based on filename.
func PubsubConfWithFilename(filename string) Configurable {
	return func(c builder.Builder, conf *Configuration) error {
		c.SetConfigName(filename)

		return conf.PubSub.FetchConfig(c)
	}
}

// LoggerConfWithFilename returns logger based on filename.
func LoggerConfWithFilename(filename string) Configurable {
	return func(c builder.Builder, conf *Configuration) error {
		c.SetConfigName(filename)

		return conf.Logger.FetchConfig(c)
	}
}
