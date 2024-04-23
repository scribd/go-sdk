package redis

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/scribd/go-sdk/pkg/cache"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     cache.Redis
		wantErr bool
	}{
		{
			name: "Config without URL set",
			cfg: cache.Redis{
				Addrs: []string{"localhost:6379"},
			},
		},
		{
			name: "Config with URL set",
			cfg: cache.Redis{
				URL: "redis://localhost:6379",
			},
		},
		{
			name: "Config with URL set to cluster URL",
			cfg: cache.Redis{
				URL: "redis://user:password@localhost:6789?dial_timeout=3&read_timeout=6s&addr=localhost:6790&addr=localhost:6791",
			},
		},
		{
			name: "Config with invalid URL",
			cfg: cache.Redis{
				URL: "localhost:6379",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(&tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCfgToRedisClientOptions(t *testing.T) {
	tests := []struct {
		name    string
		cfg     cache.Redis
		check   func(t *testing.T, opts *redis.UniversalOptions)
		wantErr bool
	}{
		{
			name: "Config without URL set",
			cfg: cache.Redis{
				Addrs: []string{"localhost:6379"},
			},
			check: func(t *testing.T, opts *redis.UniversalOptions) {
				assert.Equal(t, []string{"localhost:6379"}, opts.Addrs)
			},
		},
		{
			name: "Config with URL set",
			cfg: cache.Redis{
				URL: "redis://localhost:6379",
			},
			check: func(t *testing.T, opts *redis.UniversalOptions) {
				assert.Equal(t, []string{"localhost:6379"}, opts.Addrs)
			},
		},
		{
			name: "Config with TLS enabled",
			cfg: cache.Redis{
				URL: "rediss://localhost:6379",
				TLS: cache.TLS{
					Enabled: true,
				},
			},
			check: func(t *testing.T, opts *redis.UniversalOptions) {
				assert.NotNil(t, opts.TLSConfig)
				assert.False(t, opts.TLSConfig.InsecureSkipVerify)
			},
		},
		{
			name: "Config with URL set to cluster URL",
			cfg: cache.Redis{
				URL: "redis://user:password@localhost:6789?dial_timeout=3&read_timeout=6s&addr=localhost:6790&addr=localhost:6791",
			},
			check: func(t *testing.T, opts *redis.UniversalOptions) {
				assert.Equal(t, []string{"localhost:6789", "localhost:6790", "localhost:6791"}, opts.Addrs)
				assert.Equal(t, 3*time.Second, opts.DialTimeout)
				assert.Equal(t, 6*time.Second, opts.ReadTimeout)
				assert.Equal(t, "user", opts.Username)
				assert.Equal(t, "password", opts.Password)
			},
		},
		{
			name: "Config with invalid URL",
			cfg: cache.Redis{
				URL: "localhost:6379",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := cfgToRedisClientOptions(&tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				tt.check(t, opts)
			}
		})
	}
}
