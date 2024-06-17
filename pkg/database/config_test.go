package database

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	/*t.Run("RunningInTestEnvironment", func(t *testing.T) {
		expected := "test"
		actual := os.Getenv("APP_ENV")
		assert.Equal(t, expected, actual)
	})*/

	testCases := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "NewWithoutConfigFileFails",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewConfig()

			gotError := err != nil
			assert.Equal(t, gotError, tc.wantError)
		})
	}
}

func TestNewConfigWithAppRoot(t *testing.T) {
	testCases := []struct {
		name    string
		env     string
		cfg     *Config
		wantErr bool

		envOverrides [][]string
	}{
		{
			name: "NewWithConfigFileWorks",
			env:  "test",
			cfg: &Config{
				Host:     "mysql",
				Port:     3306,
				Username: "root",
				Password: "",
				Database: "test",
				Timeout:  "1s",
				Pool:     5,
				DBs: map[string]Config{
					"primary_replica": {
						Host:     "mysql-replica",
						Port:     3306,
						Username: "root",
						Password: "",
						Database: "test",
						Timeout:  "1s",
						Pool:     5,
						Replica:  true,
					},
				},
			},
		},
		{
			name: "NewWithConfigFileWorks, overrides",
			env:  "test",
			cfg: &Config{
				Host:     "mysql",
				Port:     3306,
				Username: "root",
				Password: "test",
				Database: "test",
				Timeout:  "1s",
				Pool:     5,
				DBs: map[string]Config{
					"primary_replica": {
						Host:     "mysql-replica",
						Port:     3306,
						Username: "root",
						Password: "test-replica",
						Database: "test",
						Timeout:  "1s",
						Pool:     5,
						Replica:  true,
					},
				},
			},
			envOverrides: [][]string{
				{"APP_DATABASE_PASSWORD", "test"},
				{"APP_DATABASE_DBS_PRIMARY_REPLICA_PASSWORD", "test-replica"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			if len(tc.envOverrides) > 0 {
				for _, o := range tc.envOverrides {
					t.Setenv(o[0], o[1])
				}
			}

			_, filename, _, _ := runtime.Caller(0)
			tmpRootParent := filepath.Dir(filename)
			t.Setenv("APP_ROOT", filepath.Join(tmpRootParent, "testdata"))

			c, err := NewConfig()
			if tc.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}

			assert.Equal(t, tc.cfg, c)
		})
	}
}
