package cache

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
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
				Store: "redis",
				Redis: Redis{
					Addrs:    []string{"localhost:6379"},
					Username: "test",
					Password: "test",
				},
			},
		},
		{
			name: "NewWithConfigFileWorks, URL set",
			env:  "test",
			cfg: &Config{
				Store: "redis",
				Redis: Redis{
					Addrs:    []string{},
					URL:      "redis://user:password@localhost:6379/0?protocol=3",
					Username: "test",
					Password: "test",
				},
			},
			envOverrides: [][]string{
				{"APP_CACHE_REDIS_ADDRS", " "},
				{"APP_CACHE_REDIS_URL", "redis://user:password@localhost:6379/0?protocol=3"},
			},
		},
		{
			name: "NewWithConfigFileWorks, incorrect store",
			env:  "test",
			cfg: &Config{
				Store: "memcached",
				Redis: Redis{
					Addrs:    []string{"localhost:6379"},
					Username: "test",
					Password: "test",
				},
			},
			wantErr:      true,
			envOverrides: [][]string{{"APP_CACHE_STORE", "memcached"}},
		},
		{
			name: "NewWithConfigFileWorks, neither URL nor Addrs set",
			env:  "test",
			cfg: &Config{
				Store: "redis",
				Redis: Redis{
					Addrs:    []string{},
					Username: "test",
					Password: "test",
				},
			},
			wantErr:      true,
			envOverrides: [][]string{{"APP_CACHE_REDIS_ADDRS", " "}, {"APP_CACHE_REDIS_URL", ""}},
		},
	}

	currentAppRoot := os.Getenv("APP_ROOT")
	defer os.Setenv("APP_ROOT", currentAppRoot)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var envVariables [][]string

			if len(tc.envOverrides) > 0 {
				for _, o := range tc.envOverrides {
					currentVal := os.Getenv(o[0])
					envVariables = append(envVariables, []string{o[0], currentVal})

					os.Setenv(o[0], o[1])
				}
			}

			_, filename, _, _ := runtime.Caller(0)
			tmpRootParent := filepath.Dir(filename)
			os.Setenv("APP_ROOT", filepath.Join(tmpRootParent, "testdata"))

			c, err := NewConfig()
			if tc.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}

			assert.Equal(t, tc.cfg, c)

			// teardown
			if len(envVariables) > 0 {
				for _, o := range envVariables {
					os.Clearenv()
					os.Setenv(o[0], o[1])
				}
			}
		})
	}
}
