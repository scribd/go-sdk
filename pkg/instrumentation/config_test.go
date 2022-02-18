package instrumentation

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

// TestNewConfig can only validate that, without a config file the configuration
// cannot be created and so returns error.
func TestNewConfig(t *testing.T) {
	testCases := []struct {
		name      string
		enabled   bool
		wantError bool
	}{
		{
			name:      "NewWithoutConfigFileFails",
			enabled:   false,
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewConfig()

			gotError := err != nil
			assert.Equal(t, gotError, tc.wantError)

			assert.Equal(t, c.Enabled, tc.enabled)
		})
	}
}

func TestNewConfigWithAppRoot(t *testing.T) {
	testCases := []struct {
		name    string
		enabled bool
	}{
		{
			name:    "NewWithConfigFileWorks",
			enabled: true,
		},
	}

	currentAppRoot := os.Getenv("APP_ROOT")
	defer os.Setenv("APP_ROOT", currentAppRoot)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, filename, _, _ := runtime.Caller(0)
			tmpRootParent := filepath.Dir(filename)
			os.Setenv("APP_ROOT", filepath.Join(tmpRootParent, "testdata"))

			c, err := NewConfig()
			require.Nil(t, err)

			assert.Equal(t, c.Enabled, tc.enabled)
			assert.Equal(t, c.CodeHotspotsEnabled, tc.enabled)
		})
	}
}

func TestNewConfigWithAppRootAndOverwriteFromEnvTheEnableFlag(t *testing.T) {
	type keyValue struct {
		key, value string

		check func(c *Config) bool
	}

	testCases := []struct {
		name string
		keys []keyValue
	}{
		{
			name: "NewWithConfigWithEnvVariablesOverwritten",
			keys: []keyValue{
				{
					key:   "APP_DATADOG_ENABLED",
					value: "false",
					check: func(c *Config) bool {
						return !c.Enabled
					},
				},
				{
					key:   "APP_DATADOG_CODE_HOTSPOTS_ENABLED",
					value: "false",
					check: func(c *Config) bool {
						return !c.CodeHotspotsEnabled
					},
				},
			},
		},
	}

	currentAppRoot := os.Getenv("APP_ROOT")
	defer os.Setenv("APP_ROOT", currentAppRoot)

	_, filename, _, _ := runtime.Caller(0)
	tmpRootParent := filepath.Dir(filename)
	os.Setenv("APP_ROOT", filepath.Join(tmpRootParent, "testdata"))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, kv := range tc.keys {
				func(kv keyValue) {
					currentEnvValue := os.Getenv(kv.key)
					os.Setenv(kv.key, kv.value)
					defer os.Setenv(kv.key, currentEnvValue)

					c, err := NewConfig()
					require.Nil(t, err)

					assert.True(t, kv.check(c))
				}(kv)
			}
		})
	}
}
