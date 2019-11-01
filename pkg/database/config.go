package database

import (
	"fmt"

	cbuilder "git.lo/microservices/sdk/go-sdk/internal/pkg/configuration/builder"
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

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("Unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
