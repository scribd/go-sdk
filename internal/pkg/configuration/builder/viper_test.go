package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViperBuilder(t *testing.T) {
	t.Run("Build", func(t *testing.T) {
		cases := []struct {
			name       string
			configName string
			wantError  bool
		}{
			{
				name:       "ValidFile",
				configName: "valid",
				wantError:  false,
			},
			{
				name:       "ValidWithoutEnvs",
				configName: "valid-no-envs",
				wantError:  true,
			},
			{
				name:       "InvalidFile",
				configName: "invalid",
				wantError:  true,
			},
			{
				name:       "NonExistentFile",
				configName: "nonexistent",
				wantError:  true,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				b := New(c.configName).ConfigPath("testdata")

				_, err := b.Build()
				gotError := err != nil

				assert.Equal(t, gotError, c.wantError)
			})
		}
	})
}
