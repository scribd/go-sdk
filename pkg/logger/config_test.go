package logger

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
		consoleEnabled    bool
		consoleJSONFormat bool
		consoleLevel      string
		fileEnabled       bool
		fileJSONFormat    bool
		fileLevel         string
		fileLocation      string
		fileName          string
		name              string
		wantError         bool
	}{
		{
			name:              "NewWithoutConfigFileFails",
			wantError:         true,
			consoleEnabled:    false,
			consoleJSONFormat: false,
			consoleLevel:      "",
			fileEnabled:       false,
			fileJSONFormat:    false,
			fileLevel:         "",
			fileLocation:      "",
			fileName:          "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewConfig()

			gotError := err != nil
			assert.Equal(t, gotError, tc.wantError)

			assert.Equal(t, c.ConsoleEnabled, tc.consoleEnabled)
			assert.Equal(t, c.ConsoleJSONFormat, tc.consoleJSONFormat)
			assert.Equal(t, c.ConsoleLevel, tc.consoleLevel)
			assert.Equal(t, c.FileEnabled, tc.fileEnabled)
			assert.Equal(t, c.FileJSONFormat, tc.fileJSONFormat)
			assert.Equal(t, c.FileLevel, tc.fileLevel)
			assert.Equal(t, c.FileLocation, tc.fileLocation)
			assert.Equal(t, c.FileName, tc.fileName)
		})
	}
}
