package kafka

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestAnalyticsSettings(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := newConfig()
		assert.True(t, math.IsNaN(cfg.analyticsRate))
	})

	t.Run("enabled", func(t *testing.T) {
		cfg := newConfig(WithAnalytics(true))
		assert.Equal(t, 1.0, cfg.analyticsRate)
	})

	t.Run("provide analytics rate", func(t *testing.T) {
		tests := []struct {
			name        string
			val         float64
			expectedNan bool
		}{
			{
				name:        "less than 0.0",
				val:         -0.1,
				expectedNan: true,
			},
			{
				name:        "more than 1.0",
				val:         1.1,
				expectedNan: true,
			},
			{
				name: "in allowed range",
				val:  0.1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				cfg := newConfig(WithAnalyticsRate(tt.val))
				if tt.expectedNan {
					assert.True(t, math.IsNaN(cfg.analyticsRate))
				} else {
					assert.Equal(t, cfg.analyticsRate, tt.val)
				}
			})
		}

	})
}

func TestWithServiceName(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := newConfig()
		assert.Equal(t, cfg.consumerServiceName, "kafka")
		assert.Equal(t, cfg.producerServiceName, "kafka")
	})

	t.Run("service name provided", func(t *testing.T) {
		cfg := newConfig(WithServiceName("test"))
		assert.Equal(t, cfg.consumerServiceName, "test")
		assert.Equal(t, cfg.producerServiceName, "test")
	})
}
