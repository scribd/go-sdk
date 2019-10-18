package builder

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
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
				wantError:  false,
			},
			{
				name:       "InvalidFile",
				configName: "invalid",
				wantError:  true,
			},
			{
				name:       "NonexistentFile",
				configName: "nonexistent",
				wantError:  true,
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				b := New().
					ConfigPath("testdata").
					ConfigName(c.configName)

				_, err := b.Build()
				gotError := err != nil

				assert.Equal(t, gotError, c.wantError)
			})
		}
	})
}
