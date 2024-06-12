package aws

import (
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
				AWSConfig: AWSConfig{
					Region: "us-east-2",
				},
				S3: map[string]S3Config{
					"default": {
						Region: "us-east-1",
					},
					"test": {
						Region: "us-west-2",
					},
				},
			},
		},
		{
			name: "NewWithConfigFileWorks, overrides",
			env:  "test",
			cfg: &Config{
				AWSConfig: AWSConfig{
					Region: "us-west-2",
				},
				S3: map[string]S3Config{
					"default": {
						Region: "us-east-1",
						Credentials: CredentialsConfig{
							AssumeRole: AssumeRoleConfig{
								ARN: "test",
							},
						},
					},
					"test": {
						Region: "us-west-1",
					},
				},
			},
			envOverrides: [][]string{
				{"APP_AWS_CONFIG_REGION", "us-west-2"},
				{"APP_AWS_S3_TEST_REGION", "us-west-1"},
				{"APP_AWS_S3_DEFAULT_CREDENTIALS_ASSUME_ROLE_ARN", "test"}},
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
