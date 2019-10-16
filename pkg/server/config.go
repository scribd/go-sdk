package server

import (
	"fmt"

	cbuilder "git.lo/microservices/sdk/go-sdk/internal/pkg/configuration/builder"

	"github.com/spf13/viper"
)

// Config represents a web server configuration
type Config struct {
	Host     string `mapstructure:"host"`
	GRPCPort string `mapstructure:"grpc_port"`
	HTTPPort string `mapstructure:"http_port"`
}

var vConf *viper.Viper

// NewConfig returns a new ServerConfig instance
func NewConfig() (*Config, error) {
	config := &Config{}
	viperBuilder := cbuilder.New().ConfigName("server")

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("Unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
