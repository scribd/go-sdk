package configuration

import (
	app "github.com/scribd/go-sdk/pkg/app"
	database "github.com/scribd/go-sdk/pkg/database"
	instrumentation "github.com/scribd/go-sdk/pkg/instrumentation"
	logger "github.com/scribd/go-sdk/pkg/logger"
	"github.com/scribd/go-sdk/pkg/pubsub"
	server "github.com/scribd/go-sdk/pkg/server"
	tracking "github.com/scribd/go-sdk/pkg/tracking"
)

// Config is an app-wide configuration
type Config struct {
	App             *app.Config
	Database        *database.Config
	Instrumentation *instrumentation.Config
	Logger          *logger.Config
	Server          *server.Config
	Tracking        *tracking.Config
	PubSub          *pubsub.Config
}

// NewConfig returns a new Config instance
func NewConfig() (*Config, error) {
	config := &Config{}

	appConfig, err := app.NewDefaultConfig()
	if err != nil {
		return config, err
	}

	dbConfig, err := database.NewConfig()
	if err != nil {
		return config, err
	}

	instrumentationConfig, err := instrumentation.NewConfig()
	if err != nil {
		return config, err
	}

	loggerConfig, err := logger.NewConfig()
	if err != nil {
		return config, err
	}

	serverConfig, err := server.NewConfig()
	if err != nil {
		return config, err
	}

	trackingConfig, err := tracking.NewConfig()
	if err != nil {
		return config, err
	}

	pubsubConfig, err := pubsub.NewConfig()
	if err != nil {
		return config, err
	}

	config.App = appConfig
	config.Database = dbConfig
	config.Instrumentation = instrumentationConfig
	config.Logger = loggerConfig
	config.Server = serverConfig
	config.Tracking = trackingConfig
	config.PubSub = pubsubConfig

	return config, nil
}
