package apps

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

// TestLoggerFetchConfig can only validate that, without a config file the configuration
// cannot be created and so returns error.
func TestLoggerFetchConfig(t *testing.T) {
	t.Run("RunningInTestEnvironment", func(t *testing.T) {
		expected := "test"
		actual := os.Getenv("APP_ENV")
		assert.Equal(t, expected, actual)
	})

	testCases := []struct {
		testName          string
		configFilename    string
		consoleEnabled    bool
		consoleJSONFormat bool
		consoleLevel      string
		fileEnabled       bool
		fileJSONFormat    bool
		fileLevel         string
		fileLocation      string
		filename          string
		wantError         bool
	}{
		{
			testName:          "NewWithoutConfigFile",
			configFilename:    "logger",
			consoleEnabled:    false,
			consoleJSONFormat: false,
			consoleLevel:      "debug",
			fileEnabled:       true,
			fileJSONFormat:    false,
			fileLevel:         "debug",
			fileLocation:      "/opt",
			filename:          "logger.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			b := newTestBuilder(tc.configFilename, "logger", t)

			lg := &Logger{}
			if err := lg.FetchConfig(b); err != nil {
				if tc.wantError {
					return
				}

				t.Fatalf("unmarshalling config, %s", err.Error())
			}

			assert.Equal(t, lg.ConsoleEnabled, tc.consoleEnabled)
			assert.Equal(t, lg.ConsoleJSONFormat, tc.consoleJSONFormat)
			assert.Equal(t, lg.ConsoleLevel, tc.consoleLevel)
			assert.Equal(t, lg.FileEnabled, tc.fileEnabled)
			assert.Equal(t, lg.FileJSONFormat, tc.fileJSONFormat)
			assert.Equal(t, lg.FileLevel, tc.fileLevel)
			assert.Equal(t, lg.FileLocation, tc.fileLocation)
			assert.Equal(t, lg.FileName, tc.filename)
		})
	}
}
