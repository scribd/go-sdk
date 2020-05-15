package tracking

import (
	"os"
	"path/filepath"
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

	currentAppRoot := os.Getenv("APP_ROOT")
	defer os.Setenv("APP_ROOT", currentAppRoot)

	currentAppVersion := os.Getenv("APP_VERSION")
	defer os.Setenv("APP_VERSION", currentAppVersion)

	currentAppServerName := os.Getenv("APP_SERVER_NAME")
	defer os.Setenv("APP_SERVER_NAME", currentAppServerName)

	testCases := []struct {
		name       string
		release    string
		sentryDSN  string
		serverName string
		withConfig bool
		wantError  bool
	}{
		{
			name:       "NewWithoutConfigFileFails",
			release:    "",
			sentryDSN:  "",
			serverName: "",
			withConfig: false,
			wantError:  true,
		},
		{
			name:    "NewWithConfigFile",
			release: "releaseTagExample",
			// The expected value is the DSN defined in ./testfiles/config/sentry.yml
			sentryDSN:  "https://aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@0000000.ingest.sentry.io/0000000",
			serverName: "serverHostnameExample",
			withConfig: true,
			wantError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("APP_VERSION", tc.release)
			os.Setenv("APP_SERVER_NAME", tc.serverName)

			if tc.withConfig {
				os.Setenv("APP_ROOT", filepath.Join("/", "sdk", "pkg", "tracking", "testfiles"))
			}

			c, err := NewConfig()

			assert.Equal(t, os.Getenv("APP_ENV"), c.environment)
			assert.Equal(t, tc.release, c.release)
			assert.Equal(t, tc.sentryDSN, c.SentryDSN)
			assert.Equal(t, tc.serverName, c.serverName)

			gotError := err != nil
			assert.Equal(t, tc.wantError, gotError)

		})
	}
}
