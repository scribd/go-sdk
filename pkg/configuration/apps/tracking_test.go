package apps

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

// TestNewConfig can only validate that, without a config file the configuration
// cannot be created and so returns error.
func TestTrackingFetchConfig(t *testing.T) {
	t.Run("RunningInTestEnvironment", func(t *testing.T) {
		expected := "test"
		actual := os.Getenv("APP_ENV")
		assert.Equal(t, expected, actual)
	})

	testCases := []struct {
		name       string
		filename   string
		release    string
		sentryDSN  string
		serverName string
		wantError  bool
	}{
		{
			name:       "NewWithoutConfigFileFails",
			filename:   "imaginarytracking",
			release:    "",
			sentryDSN:  "",
			serverName: "",
			wantError:  true,
		},
		{
			name:     "NewWithConfigFile",
			filename: "sentry",
			release:  "releaseTagExample",
			// The expected value is the DSN defined in ./testdata/config/sentry.yml
			sentryDSN:  "https://aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@0000000.ingest.sentry.io/0000000",
			serverName: "serverHostnameExample",
			wantError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("APP_VERSION", tc.release)
			t.Setenv("APP_SERVER_NAME", tc.serverName)

			trc := &Tracking{}
			b := newTestBuilder(tc.filename, "tracking", t)

			if err := trc.FetchConfig(b); err != nil {
				if tc.wantError {
					return
				}

				t.Fatalf("getching tracking config, %s", err.Error())
			}

			assert.Equal(t, "test", trc.Environment)
			assert.Equal(t, tc.release, trc.Release)
			assert.Equal(t, tc.sentryDSN, trc.SentryDSN)
			assert.Equal(t, tc.serverName, trc.ServerName)
		})
	}
}
