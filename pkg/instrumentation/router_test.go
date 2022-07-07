package instrumentation

import (
	"testing"

	"github.com/scribd/go-sdk/pkg/configuration/apps"
	assert "github.com/stretchr/testify/assert"
)

func TestNewTracer(t *testing.T) {
	testCases := []struct {
		name      string
		enabled   bool
		wantError bool
	}{
		{
			name:      "WithAConfigurationIsAlwaysOkWhenEnabled",
			enabled:   true,
			wantError: false,
		},
		{
			name:      "WithAConfigurationIsAlwaysOkWhenDisabled",
			enabled:   false,
			wantError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := apps.Instrumentation{
				Enabled: tc.enabled,
			}
			trcr := NewTracer(config)
			assert.Equal(t, trcr.Enabled, tc.enabled)
		})
	}
}
