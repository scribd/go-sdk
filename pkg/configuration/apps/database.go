package apps

import (
	"fmt"
	"os"
	"strings"

	"github.com/scribd/go-sdk/pkg/configuration/builder"
)

// Database is the database connection configuration.
type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Pool     int    `mapstructure:"pool"`
	Timeout  string `mapstructure:"timeout"`
}

// FetchConfig fetches data from builder into it's caller.
func (d *Database) FetchConfig(b builder.Builder) error {
	appName := strings.ReplaceAll(os.Getenv("APP_SETTINGS_NAME"), "-", "_")
	b.SetDefault("database", fmt.Sprintf("%s_%s", appName, os.Getenv("APP_ENV")))

	config, err := b.Build()
	if err != nil {
		return fmt.Errorf("building for database. err: %w", err)
	}

	if err = config.Unmarshal(d); err != nil {
		return fmt.Errorf("Unable to decode database into struct. err: %w", err)
	}

	return nil
}
