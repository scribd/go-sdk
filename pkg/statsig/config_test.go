package statsig

import (
	"os"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("RunningInTestEnvironment", func(t *testing.T) {
		expected := "test"
		actual := os.Getenv("APP_ENV")
		assert.Equal(t, expected, actual)
	})

	testCases := []struct {
		name               string
		secretKey          string
		localMode          bool
		configSyncInterval time.Duration
		idListSyncInterval time.Duration
		wantError          bool
	}{
		{
			name:               "default",
			secretKey:          "",
			localMode:          false,
			configSyncInterval: 0,
			idListSyncInterval: 0,
			wantError:          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewConfig()
			if tc.wantError {
				assert.Error(t, err)
			}

			assert.Equal(t, c.SecretKey, tc.secretKey)
			assert.Equal(t, c.LocalMode, tc.localMode)
			assert.Equal(t, c.ConfigSyncInterval, tc.configSyncInterval)
			assert.Equal(t, c.IDListSyncInterval, tc.idListSyncInterval)
		})
	}
}
