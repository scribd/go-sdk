package database

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
		name                          string
		wantError                     bool
		host                          string
		port                          int
		username                      string
		password                      string
		database                      string
		timeout                       string
		pool                          int
		maxOpenConnections            int
		connectionMaxIdleTime         time.Duration
		connectionMaxLifetime         time.Duration
		disableDefaultGormTransaction bool
		cachePreparedStatements       bool
		mysqlInterpolateParams        bool
	}{
		{
			name:                          "NewWithoutConfigFileFails",
			wantError:                     true,
			host:                          "",
			port:                          0,
			username:                      "",
			password:                      "",
			database:                      "",
			timeout:                       "",
			pool:                          0,
			maxOpenConnections:            0,
			connectionMaxIdleTime:         0,
			connectionMaxLifetime:         0,
			disableDefaultGormTransaction: false,
			cachePreparedStatements:       false,
			mysqlInterpolateParams:        false,
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
			assert.Equal(t, c.Timeout, tc.timeout)
			assert.Equal(t, c.Pool, tc.pool)
			assert.Equal(t, c.MaxOpenConnections, tc.maxOpenConnections)
			assert.Equal(t, c.ConnectionMaxIdleTime, tc.connectionMaxIdleTime)
			assert.Equal(t, c.ConnectionMaxLifetime, tc.connectionMaxLifetime)
			assert.Equal(t, c.DisableDefaultGormTransaction, tc.disableDefaultGormTransaction)
			assert.Equal(t, c.CachePreparedStatements, tc.cachePreparedStatements)
			assert.Equal(t, c.MysqlInterpolateParams, tc.mysqlInterpolateParams)
		})
	}
}
