package metrics

import (
	"fmt"
	datadogstatsd "github.com/DataDog/datadog-go/statsd"
)

// Builder is a Metrics builder.
type Builder struct {
	config *Config
}

// NewBuilder initializes a Metrics builder with the given configuration.
func NewBuilder(config *Config) *Builder {
	return &Builder{
		config: config,
	}
}

// Build applies the given configuration and returns a Metrics instance.
func (b *Builder) Build() (Metrics, error) {
	// New returns a pointer to a new Client given an addr in the
	// format "hostname:port" or "unix:///path/to/socket".
	//
	// If the addr parameter is empty, the client uses the
	// DD_AGENT_HOST and (optionally) the DD_DOGSTATSD_PORT
	// environment variables to build a target address.
	dogstatsd, err := datadogstatsd.New("")
	if err != nil {
		return nil, err
	}

	// Namespace to prepend to all statsd calls.
	dogstatsd.Namespace = b.config.App + "."

	// Tags are global tags to be added to every statsd call.
	dogstatsd.Tags = []string{
		fmt.Sprintf("service:%s", b.config.App),
		fmt.Sprintf("env:%s", b.config.Environment),
	}

	return dogstatsd, nil
}
