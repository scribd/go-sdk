package configuration

import (
	app "git.lo/microservices/sdk/go-sdk/pkg/app"
	database "git.lo/microservices/sdk/go-sdk/pkg/database"
	logger "git.lo/microservices/sdk/go-sdk/pkg/logger"
	server "git.lo/microservices/sdk/go-sdk/pkg/server"
	tracking "git.lo/microservices/sdk/go-sdk/pkg/tracking"
)

// Config is an app-wide configuration
type Config struct {
	App      *app.Config
	Database *database.Config
	Logger   *logger.Config
	Server   *server.Config
	Tracking *tracking.Config
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

	config.App = appConfig
	config.Database = dbConfig
	config.Logger = loggerConfig
	config.Server = serverConfig
	config.Tracking = trackingConfig

	return config, nil
}
