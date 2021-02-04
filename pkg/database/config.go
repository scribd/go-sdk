package database

import (
	"fmt"
	"os"
	"strings"

	cbuilder "github.com/scribd/go-sdk/internal/pkg/configuration/builder"
)

// Config is the database connection configuration.
type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Pool     int    `mapstructure:"pool"`
	Timeout  string `mapstructure:"timeout"`
}

// NewConfig returns a new Config instance.
func NewConfig() (*Config, error) {
	config := &Config{}
	viperBuilder := cbuilder.New("database")

	appName := strings.ReplaceAll(os.Getenv("APP_SETTINGS_NAME"), "-", "_")
	viperBuilder.SetDefault("database", fmt.Sprintf("%s_%s", appName, os.Getenv("APP_ENV")))

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("Unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
