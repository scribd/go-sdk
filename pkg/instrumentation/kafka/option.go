package kafka

import (
	"math"
	"os"
	"strconv"
)

type config struct {
	consumerServiceName string
	producerServiceName string
	analyticsRate       float64
}

// An Option customizes the config.
type Option func(cfg *config)

func newConfig(opts ...Option) *config {
	cfg := &config{
		consumerServiceName: "kafka",
		producerServiceName: "kafka",
		analyticsRate:       math.NaN(),
	}

	var datadogAnalyticsEnabled bool
	datadogAnalyticsEnabledString := os.Getenv("DD_TRACE_KAFKA_ANALYTICS_ENABLED")
	v, err := strconv.ParseBool(datadogAnalyticsEnabledString)
	if err == nil {
		datadogAnalyticsEnabled = v
	}

	if datadogAnalyticsEnabled {
		cfg.analyticsRate = 1.0
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// WithServiceName sets the config service name to serviceName.
func WithServiceName(serviceName string) Option {
	return func(cfg *config) {
		cfg.consumerServiceName = serviceName
		cfg.producerServiceName = serviceName
	}
}

// WithAnalytics enables Trace Analytics for all started spans.
func WithAnalytics(on bool) Option {
	return func(cfg *config) {
		if on {
			cfg.analyticsRate = 1.0
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}

// WithAnalyticsRate sets the sampling rate for Trace Analytics events
// correlated to started spans.
func WithAnalyticsRate(rate float64) Option {
	return func(cfg *config) {
		if rate >= 0.0 && rate <= 1.0 {
			cfg.analyticsRate = rate
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}
