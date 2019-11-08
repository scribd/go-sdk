package instrumentation

import (
	"fmt"
	"os"
	"path/filepath"
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
			os.Setenv("APP_ROOT", filepath.Join("/", "sdk", "pkg", "instrumentation", "testfiles"))

			c, err := NewConfig()
			require.Nil(t, err)

			assert.Equal(t, c.Enabled, tc.enabled)
		})
	}
}

func overrideAppRootAndTest(testedVariable string, testFunc func(string)) {
	currentAppRoot := os.Getenv("APP_ROOT")
	defer os.Setenv("APP_ROOT", currentAppRoot)
	os.Setenv("APP_ROOT", filepath.Join("/", "sdk", "pkg", "instrumentation", "testfiles"))

	overwrittenValue := "false"
	currentEnvValue := os.Getenv(testedVariable)
	os.Setenv(testedVariable, overwrittenValue)
	defer os.Setenv(testedVariable, currentEnvValue)

	testFunc(overwrittenValue)
}

func TestNewConfigWithAppRootAndOverwriteFromEnvTheEnableFlag(t *testing.T) {
	testCases := []struct {
		name    string
		keyName string
	}{
		{
			name:    "NewWithConfigFileWorks",
			keyName: "APP_DATADOG_ENABLED",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			overrideAppRootAndTest(tc.keyName, func(overwrittenValue string) {
				c, err := NewConfig()
				require.Nil(t, err)

				assert.Equal(t, fmt.Sprintf("%t", c.Enabled), overwrittenValue)
			})
		})
	}
}
