package apps

import (
	"fmt"
	"os"

	"github.com/scribd/go-sdk/pkg/configuration/builder"
)

// Config stores the configuration for the tracking.
type Tracking struct {
	SentryDSN string `mapstructure:"dsn"`

	Environment string
	Release     string
	ServerName  string
}

// FetchConfig fetches data from builder into it's caller.
func (t *Tracking) FetchConfig(b builder.Builder) error {
	config, err := b.Build()
	if err != nil {
		return err
	}

	if err = config.Unmarshal(t); err != nil {
		return fmt.Errorf("Unable to decode into struct: %w", err)
	}

	t.Environment = config.String("ENV")
	// TODO os must be fetched from config not os.
	t.Release = os.Getenv("APP_VERSION")
	t.ServerName = os.Getenv("APP_SERVER_NAME")

	return nil
}
