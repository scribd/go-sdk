package server

import (
	"fmt"
	"os"
	"path"

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

	vConf = viper.New()

	vConf.SetConfigType("yaml")
	vConf.SetConfigName("server")
	vConf.AddConfigPath(path.Join(os.Getenv("APP_ROOT"), "config"))
	vConf.SetEnvPrefix("APP")
	vConf.AutomaticEnv()

	if err := vConf.ReadInConfig(); err != nil {
		return config, err
	}

	err := vConf.Unmarshal(config)
	if err != nil {
		return config, fmt.Errorf("Unable to decode into struct: %v", err)
	}

	return config, nil
}
