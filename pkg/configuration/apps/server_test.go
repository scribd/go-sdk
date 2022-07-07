package apps

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerFetchConfig(t *testing.T) {
	testCases := []struct {
		name        string
		filename    string
		enabled     bool
		httpTimeout HTTPTimeout
		wantError   bool
		settings    []CorsSetting
	}{
		{
			name:     "NewWithConfigFileWorks",
			filename: "server",
			enabled:  true,
			httpTimeout: HTTPTimeout{
				Write: time.Second * 2,
				Read:  time.Second * 1,
				Idle:  time.Second * 90,
			},
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
		{
			name:      "NewWithoutConfigFileFails",
			filename:  "imaginaryfile",
			enabled:   false,
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := newTestBuilder(tc.filename, "server", t)
			srv := &Server{}

			if err := srv.FetchConfig(b); err != nil {
				if tc.wantError {
					return
				}

				t.Fatalf("unmarshalling config, %s", err.Error())
			}

			// asserting http timeouts.
			assert.NotEmpty(t, srv.HTTPTimeout)

			assert.Equal(t, srv.HTTPTimeout.Write, tc.httpTimeout.Write)
			assert.Equal(t, srv.HTTPTimeout.Read, tc.httpTimeout.Read)
			assert.Equal(t, srv.HTTPTimeout.Idle, tc.httpTimeout.Idle)

			// asserting cors
			assert.Equal(t, tc.settings, srv.GetCorsSettings())
			assert.True(t, srv.Cors.Settings[0].Matches("/test"))

			assert.Equal(t, srv.Cors.Enabled, tc.enabled)
			assert.Equal(t, srv.Cors.Settings, tc.settings)
		})
	}
}
