package statsig

import (
	"fmt"
	"os"
	"time"

	cbuilder "github.com/scribd/go-sdk/internal/pkg/configuration/builder"
)

// Config stores the configuration for the statsig.
type Config struct {
	SecretKey          string        `mapstructure:"secret_key"`
	LocalMode          bool          `mapstructure:"local_mode"`
	ConfigSyncInterval time.Duration `mapstructure:"config_sync_interval"`
	IDListSyncInterval time.Duration `mapstructure:"id_list_sync_interval"`

	environment string
}

// NewConfig returns a new StatsigConfig instance.
func NewConfig() (*Config, error) {
	config := &Config{}
	config.environment = os.Getenv("APP_ENV")

	viperBuilder := cbuilder.New("statsig")

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
