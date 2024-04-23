package redis

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/redis/go-redis/v9"

	"github.com/scribd/go-sdk/pkg/cache"
)

func New(cfg *cache.Redis) (redis.UniversalClient, error) {
	opts, err := cfgToRedisClientOptions(cfg)
	if err != nil {
		return nil, err
	}

	return redis.NewUniversalClient(opts), nil
}

func cfgToRedisClientOptions(cfg *cache.Redis) (*redis.UniversalOptions, error) {
	var err error
	var clusterOptions *redis.ClusterOptions
	if cfg.URL != "" {
		clusterOptions, err = redis.ParseClusterURL(cfg.URL)
		if err != nil {
			return nil, err
		}
	}

	opts := &redis.UniversalOptions{
		Addrs:      cfg.Addrs,
		DB:         cfg.DB,
		ClientName: cfg.ClientName,

		Protocol: cfg.Protocol,
		Username: cfg.Username,
		Password: cfg.Password,

		SentinelUsername: cfg.SentinelUsername,
		SentinelPassword: cfg.SentinelPassword,

		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: cfg.MinRetryBackoff,
		MaxRetryBackoff: cfg.MaxRetryBackoff,

		DialTimeout:           cfg.DialTimeout,
		ReadTimeout:           cfg.ReadTimeout,
		WriteTimeout:          cfg.WriteTimeout,
		ContextTimeoutEnabled: cfg.ContextTimeoutEnabled,

		PoolSize:        cfg.PoolSize,
		PoolTimeout:     cfg.PoolTimeout,
		MaxIdleConns:    cfg.MaxIdleConns,
		MinIdleConns:    cfg.MinIdleConns,
		MaxActiveConns:  cfg.MaxActiveConns,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		ConnMaxLifetime: cfg.ConnMaxLifetime,

		MaxRedirects:   cfg.MaxRedirects,
		ReadOnly:       cfg.ReadOnly,
		RouteByLatency: cfg.RouteByLatency,
		RouteRandomly:  cfg.RouteRandomly,

		MasterName:       cfg.MasterName,
		DisableIndentity: cfg.DisableIndentity,
		IdentitySuffix:   cfg.IdentitySuffix,
	}
	if clusterOptions != nil {
		opts.Addrs = clusterOptions.Addrs
		opts.ClientName = clusterOptions.ClientName

		opts.Protocol = clusterOptions.Protocol
		opts.Username = clusterOptions.Username
		opts.Password = clusterOptions.Password

		if clusterOptions.MaxRetries != 0 {
			opts.MaxRetries = clusterOptions.MaxRetries
		}
		if clusterOptions.MinRetryBackoff != 0 {
			opts.MinRetryBackoff = clusterOptions.MinRetryBackoff
		}
		if clusterOptions.MaxRetryBackoff != 0 {
			opts.MaxRetryBackoff = clusterOptions.MaxRetryBackoff
		}

		if clusterOptions.DialTimeout != 0 {
			opts.DialTimeout = clusterOptions.DialTimeout
		}
		if clusterOptions.ReadTimeout != 0 {
			opts.ReadTimeout = clusterOptions.ReadTimeout
		}
		if clusterOptions.WriteTimeout != 0 {
			opts.WriteTimeout = clusterOptions.WriteTimeout
		}
		if clusterOptions.ContextTimeoutEnabled {
			opts.ContextTimeoutEnabled = clusterOptions.ContextTimeoutEnabled
		}

		if clusterOptions.PoolSize != 0 {
			opts.PoolSize = clusterOptions.PoolSize
		}
		if clusterOptions.PoolTimeout != 0 {
			opts.PoolTimeout = clusterOptions.PoolTimeout
		}
		if clusterOptions.MaxIdleConns != 0 {
			opts.MaxIdleConns = clusterOptions.MaxIdleConns
		}
		if clusterOptions.MinIdleConns != 0 {
			opts.MinIdleConns = clusterOptions.MinIdleConns
		}
		if clusterOptions.MaxActiveConns != 0 {
			opts.MaxActiveConns = clusterOptions.MaxActiveConns
		}
		if clusterOptions.ConnMaxIdleTime != 0 {
			opts.ConnMaxIdleTime = clusterOptions.ConnMaxIdleTime
		}
		if clusterOptions.ConnMaxLifetime != 0 {
			opts.ConnMaxLifetime = clusterOptions.ConnMaxLifetime
		}

		if clusterOptions.MaxRedirects != 0 {
			opts.MaxRedirects = clusterOptions.MaxRedirects
		}
		if clusterOptions.ReadOnly {
			opts.ReadOnly = clusterOptions.ReadOnly
		}
		if clusterOptions.RouteByLatency {
			opts.RouteByLatency = clusterOptions.RouteByLatency
		}
		if clusterOptions.RouteRandomly {
			opts.RouteRandomly = clusterOptions.RouteRandomly
		}
	}

	if cfg.TLS.Enabled {
		var caCertPool *x509.CertPool

		if cfg.TLS.Ca != "" {
			caCertPool = x509.NewCertPool()
			caCertPool.AppendCertsFromPEM([]byte(cfg.TLS.Ca))
		}

		var certificates []tls.Certificate
		if cfg.TLS.Cert != "" && cfg.TLS.CertKey != "" {
			cert, err := tls.X509KeyPair([]byte(cfg.TLS.Cert), []byte(cfg.TLS.CertKey))
			if err != nil {
				return nil, err
			}
			certificates = []tls.Certificate{cert}
		}

		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: cfg.TLS.InsecureSkipVerify,
			Certificates:       certificates,
			RootCAs:            caCertPool,
		}
	}

	return opts, nil
}
