package logger

import (
	"fmt"
	"os"
	"path"

	cbuilder "git.lo/microservices/sdk/go-sdk/internal/pkg/configuration/builder"

	"github.com/spf13/viper"
)

// Config stores the configuration for the logger.
// For some loggers there can only be one level across writers, for such
// the level of Console is picked by default.
type Config struct {
	ConsoleEnabled    bool   `mapstructure:"console_enabled"`
	ConsoleJSONFormat bool   `mapstructure:"console_json_format"`
	ConsoleLevel      string `mapstructure:"console_level"`
	FileEnabled       bool   `mapstructure:"file_enabled"`
	FileJSONFormat    bool   `mapstructure:"file_json_format"`
	FileLevel         string `mapstructure:"file_level"`
	FileLocation      string `mapstructure:"file_location"`
	FileName          string `mapstructure:"file_name"`
}

var vConf *viper.Viper

// NewConfig returns a new ServerConfig instance
func NewConfig() (*Config, error) {
	config := &Config{}
	viperBuilder := cbuilder.
		New().
		ConfigName("logger").
		SetDefault("file_location", path.Join(os.Getenv("APP_ROOT"), "log")).
		SetDefault("file_name", os.Getenv("APP_ENV"))

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("Unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
