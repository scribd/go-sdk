package apps

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestInstrumentationFetchConfig(t *testing.T) {
	t.Run("RunningInTestEnvironment", func(t *testing.T) {
		expected := expectedTestEnv
		actual := os.Getenv("APP_ENV")
		assert.Equal(t, expected, actual)
	})

	testCases := []struct {
		name      string
		filename  string
		wantError bool

		// conf related fields.
		environment         string
		enabled             bool
		codeHotspotsEnabled bool
	}{
		{
			name:                "NewWithoutConfigFileFails",
			filename:            "imaginaryfile",
			wantError:           true,
			environment:         "test",
			enabled:             true,
			codeHotspotsEnabled: true,
		},
		{
			name:                "NewDatabaseConfig",
			filename:            "instrumentation",
			wantError:           false,
			environment:         "test",
			enabled:             true,
			codeHotspotsEnabled: true,
		},
		{
			name:      "NewWithoutConfigFileFails",
			filename:  "imaginaryinstrumentation",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := newTestBuilder(tc.filename, "instrumentation", t)

			ins := &Instrumentation{}

			if err := ins.FetchConfig(b); err != nil {
				if tc.wantError {
					return
				}

				t.Fatalf("unmarshalling config, %s", err.Error())
			}

			assert.Equal(t, tc.environment, ins.Environment)
			assert.Equal(t, tc.enabled, ins.Enabled)
			assert.Equal(t, tc.codeHotspotsEnabled, ins.CodeHotspotsEnabled)
		})
	}
}
