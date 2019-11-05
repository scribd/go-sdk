package tracking

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

// TestNewConfig can only validate that, without a config file the configuration
// cannot be created and so returns error.
func TestNewConfig(t *testing.T) {
	t.Run("RunningInTestEnvironment", func(t *testing.T) {
		expected := "test"
		actual := os.Getenv("APP_ENV")
		assert.Equal(t, expected, actual)
	})

	testCases := []struct {
		name      string
		sentryDSN string
		wantError bool
	}{
		{
			name:      "NewWithoutConfigFileFails",
			sentryDSN: "",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewConfig()

			gotError := err != nil
			assert.Equal(t, gotError, tc.wantError)

			assert.Equal(t, c.SentryDSN, tc.sentryDSN)
		})
	}
}