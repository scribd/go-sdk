package database

import (
	"fmt"
	"os"
	"strings"
	"time"

	cbuilder "github.com/scribd/go-sdk/internal/pkg/configuration/builder"
)

// Config is the database connection configuration.
type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Timeout  string `mapstructure:"timeout"`
	// Connection settings
	// TODO Pool field name must be modified in the next major change.
	Pool                  int           `mapstructure:"pool"`
	MaxOpenConnections    int           `mapstructure:"max_open_connections"`
	ConnectionMaxIdleTime time.Duration `mapstructure:"connection_max_idle_time"`
	ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`

	// Performance settings
	DisableDefaultGormTransaction bool `mapstructure:"disable_default_gorm_transaction"`
	CachePreparedStatements       bool `mapstructure:"cache_prepared_statements"`
	MysqlInterpolateParams        bool `mapstructure:"mysql_interpolate_params"`
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
		return config, fmt.Errorf("unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
