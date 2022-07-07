package apps

import (
	"os"
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestDatabaseFetchConfig(t *testing.T) {
	t.Run("RunningInTestEnvironment", func(t *testing.T) {
		expected := expectedTestEnv
		actual := os.Getenv("APP_ENV")
		assert.Equal(t, expected, actual)
	})

	testCases := []struct {
		name      string
		filename  string
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
			name:      "NewDatabaseConfig",
			filename:  "database",
			wantError: false,
			host:      "mysql",
			port:      3306,
			username:  "root",
			password:  "zeus",
			database:  "testdb",
			pool:      5,
			timeout:   "1s",
		},
		{
			name:      "NewDatabaseConfigWithSubChange",
			filename:  "database2",
			wantError: false,
			host:      "mysql2",
			port:      3306,
			username:  "notrootforsure",
			password:  "zeus",
			database:  "testdb",
			pool:      5,
			timeout:   "1s",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := newTestBuilder(tc.filename, "database", t)
			db := &Database{}

			if err := db.FetchConfig(b); err != nil {
				if tc.wantError {
					return
				}

				t.Fatalf("unmarshalling config, %s", err.Error())
			}

			assert.Equal(t, db.Host, tc.host)
			assert.Equal(t, db.Port, tc.port)
			assert.Equal(t, db.Username, tc.username)
			assert.Equal(t, db.Password, tc.password)
			assert.Equal(t, db.Database, tc.database)
			assert.Equal(t, db.Pool, tc.pool)
			assert.Equal(t, db.Timeout, tc.timeout)
		})
	}
}
