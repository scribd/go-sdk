package database

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("RunningInTestEnvironment", func(t *testing.T) {
		expected := "test"
		actual := os.Getenv("APP_ENV")
		assert.Equal(t, expected, actual)
	})

	testCases := []struct {
		name      string
		wantError bool
		host      string
		port      int
		username  string
		password  string
		database  string
		pool      int
		timeout   string
	}{
		{
			name:      "NewWithoutConfigFileFails",
			wantError: true,
			host:      "",
			port:      0,
			username:  "",
			password:  "",
			database:  "",
			pool:      0,
			timeout:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewConfig()

			gotError := err != nil
			assert.Equal(t, gotError, tc.wantError)

			assert.Equal(t, c.Host, tc.host)
			assert.Equal(t, c.Port, tc.port)
			assert.Equal(t, c.Username, tc.username)
			assert.Equal(t, c.Password, tc.password)
			assert.Equal(t, c.Database, tc.database)
			assert.Equal(t, c.Pool, tc.pool)
			assert.Equal(t, c.Timeout, tc.timeout)
		})
	}
}
