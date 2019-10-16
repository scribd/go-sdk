package configuration

import (
	app "git.lo/microservices/sdk/go-sdk/pkg/app"
	server "git.lo/microservices/sdk/go-sdk/pkg/server"
)

// Config is an app-wide configuration
type Config struct {
	App    *app.Config
	Server *server.Config
}

// NewConfig returns a new Config instance
func NewConfig() (*Config, error) {
	config := &Config{}

	appConfig, err := app.NewDefaultConfig()
	if err != nil {
		return config, err
	}

	serverConfig, err := server.NewConfig()
	if err != nil {
		return config, err
	}

	config.App = appConfig
	config.Server = serverConfig

	return config, nil
}
