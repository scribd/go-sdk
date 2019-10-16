package configuration

import (
	app "git.lo/microservices/sdk/go-sdk/pkg/app"
)

// Config is an app-wide configuration
type Config struct {
	App *app.Config
}

// NewConfig returns a new Config instance
func NewConfig() (*Config, error) {
	config := &Config{}

	appConfig, err := app.NewConfig()
	if err != nil {
		return config, err
	}

	config.App = appConfig

	return config, nil
}
