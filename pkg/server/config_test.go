package server

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

			assert.Equal(t, c.Cors.Enabled, tc.enabled)
		})
	}
}

func TestNewConfigWithAppRoot(t *testing.T) {
	testCases := []struct {
		name     string
		enabled  bool
		settings []CorsSetting
	}{
		{
			name:    "NewWithConfigFileWorks",
			enabled: true,
			settings: []CorsSetting{{
				Path:             "*",
				AllowCredentials: true,
				AllowedHeaders:   []string{"Allowed-Header"},
				AllowedMethods:   []string{"GET"},
				AllowedOrigins:   []string{"*"},
				ExposedHeaders:   []string{"Exposed-Header"},
				MaxAge:           600,
			}},
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

			assert.Equal(t, tc.settings, c.GetCorsSettings())
			assert.True(t, c.Cors.Settings[0].Matches("/test"))

			assert.Equal(t, c.Cors.Enabled, tc.enabled)
			assert.Equal(t, c.Cors.Settings, tc.settings)
		})
	}
}
