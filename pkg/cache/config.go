package cache

import (
	"fmt"
	"os"
	"strings"
	"time"

	cbuilder "github.com/scribd/go-sdk/internal/pkg/configuration/builder"
)

type (
	StoreType int

	// Redis provides configuration for redis cache.
	Redis struct {
		// URL into Redis ClusterOptions that can be used to connect to Redis
		URL string `mapstructure:"url"`

		// Either a single address or a seed list of host:port addresses
		// of cluster/sentinel nodes.
		Addrs []string `mapstructure:"addrs"`

		// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
		ClientName string `mapstructure:"client_name"`

		// Database to be selected after connecting to the server.
		// Only single-node and failover clients.
		DB int `mapstructure:"db"`

		// Protocol 2 or 3. Use the version to negotiate RESP version with redis-server.
		Protocol int `mapstructure:"protocol"`
		// Use the specified Username to authenticate the current connection
		// with one of the connections defined in the ACL list when connecting
		// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
		Username string `mapstructure:"username"`
		// Optional password. Must match the password specified in the
		// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
		// or the User Password when connecting to a Redis 6.0 instance, or greater,
		// that is using the Redis ACL system.
		Password string `mapstructure:"password"`

		// If specified with SentinelPassword, enables ACL-based authentication (via
		// AUTH <user> <pass>).
		SentinelUsername string `mapstructure:"sentinel_username"`
		// Sentinel password from "requirepass <password>" (if enabled) in Sentinel
		// configuration, or, if SentinelUsername is also supplied, used for ACL-based
		// authentication.
		SentinelPassword string `mapstructure:"sentinel_password"`

		// Maximum number of retries before giving up.
		MaxRetries int `mapstructure:"max_retries"`
		// Minimum backoff between each retry.
		MinRetryBackoff time.Duration `mapstructure:"min_retry_backoff"`
		// Maximum backoff between each retry.
		MaxRetryBackoff time.Duration `mapstructure:"max_retry_backoff"`

		// Dial timeout for establishing new connections.
		DialTimeout time.Duration `mapstructure:"dial_timeout"`
		// Timeout for socket reads. If reached, commands will fail
		// with a timeout instead of blocking. Supported values:
		//   - `0` - default timeout (3 seconds).
		//   - `-1` - no timeout (block indefinitely).
		//   - `-2` - disables SetReadDeadline calls completely.
		ReadTimeout time.Duration `mapstructure:"read_timeout"`
		// Timeout for socket writes. If reached, commands will fail
		// with a timeout instead of blocking.  Supported values:
		//   - `0` - default timeout (3 seconds).
		//   - `-1` - no timeout (block indefinitely).
		//   - `-2` - disables SetWriteDeadline calls completely.
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		// ContextTimeoutEnabled controls whether the client respects context timeouts and deadlines.
		// See https://redis.uptrace.dev/guide/go-redis-debugging.html#timeouts
		ContextTimeoutEnabled bool `mapstructure:"context_timeout_enabled"`

		// Base number of socket connections.
		// If there is not enough connections in the pool, new connections will be allocated in excess of PoolSize,
		// you can limit it through MaxActiveConns
		PoolSize int `mapstructure:"pool_size"`
		// Amount of time client waits for connection if all connections
		// are busy before returning an error.
		PoolTimeout time.Duration `mapstructure:"pool_timeout"`
		// Maximum number of idle connections.
		MaxIdleConns int `mapstructure:"max_idle_conns"`
		// Minimum number of idle connections which is useful when establishing
		// new connection is slow.
		MinIdleConns int `mapstructure:"min_idle_conns"`
		// Maximum number of connections allocated by the pool at a given time.
		// When zero, there is no limit on the number of connections in the pool.
		MaxActiveConns int `mapstructure:"max_active_conns"`
		// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
		// Should be less than server's timeout.
		//
		// Expired connections may be closed lazily before reuse.
		// If d <= 0, connections are not closed due to a connection's idle time.
		ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
		// ConnMaxLifetime is the maximum amount of time a connection may be reused.
		//
		// Expired connections may be closed lazily before reuse.
		// If <= 0, connections are not closed due to a connection's age.
		ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`

		// Only cluster clients.

		// The maximum number of retries before giving up. Command is retried
		// on network errors and MOVED/ASK redirects.
		MaxRedirects int `mapstructure:"max_redirects"`
		// Enables read-only commands on slave nodes.
		ReadOnly bool `mapstructure:"read_only"`
		// Allows routing read-only commands to the closest master or slave node.
		// It automatically enables ReadOnly.
		RouteByLatency bool `mapstructure:"route_by_latency"`
		// Allows routing read-only commands to the random master or slave node.
		// It automatically enables ReadOnly.
		RouteRandomly bool `mapstructure:"route_randomly"`

		// The sentinel master name.
		// Only failover clients.

		// The master name.
		MasterName string `mapstructure:"master_name"`

		// Disable set-lib on connect.
		DisableIndentity bool `mapstructure:"disable_indentity"`

		// Add suffix to client name.
		IdentitySuffix string `mapstructure:"identity_suffix"`

		// TLS configuration
		TLS TLS `mapstructure:"tls"`
	}

	TLS struct {
		// Enabled whether the TLS connection is enabled or not
		Enabled bool `mapstructure:"enabled"`

		// Ca Root CA certificate
		Ca string `mapstructure:"ca"`
		// Cert is a PEM certificate string
		Cert string `mapstructure:"cert_pem"`
		// CertKey is a PEM key certificate string
		CertKey string `mapstructure:"cert_pem_key"`
		// Passphrase is used in case the private key needs to be decrypted
		Passphrase string `mapstructure:"passphrase"`
		// InsecureSkipVerify whether to skip TLS verification or not
		InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
	}

	// Config provides configuration for cache.
	Config struct {
		Store string `mapstructure:"store"`
		Redis Redis  `mapstructure:"redis"`
	}
)

const (
	storeTypeRedisName = "redis"
)

func NewConfig() (*Config, error) {
	config := &Config{}
	viperBuilder := cbuilder.New("cache")

	appName := strings.ReplaceAll(os.Getenv("APP_SETTINGS_NAME"), "-", "_")
	viperBuilder.SetDefault("cache", fmt.Sprintf("%s_%s", appName, os.Getenv("APP_ENV")))

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %s", err.Error())
	}

	config.Redis.Addrs = vConf.GetStringSlice("redis.addrs")

	if err := config.validate(); err != nil {
		return config, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.Store == "" {
		return fmt.Errorf("store is required")
	}

	switch c.Store {
	case storeTypeRedisName:
		if c.Redis.URL == "" && len(c.Redis.Addrs) == 0 {
			return fmt.Errorf("url or addrs is required for redis")
		}
	default:
		return fmt.Errorf("store %s is not supported", c.Store)
	}

	return nil
}
