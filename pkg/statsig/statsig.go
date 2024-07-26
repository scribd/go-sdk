package statsig

import (
	statsig "github.com/statsig-io/go-sdk"
)

func Initialize(c *Config) {
	opts := &statsig.Options{
		Environment: statsig.Environment{Tier: c.environment},
	}

	if c.LocalMode {
		opts.LocalMode = true
	}

	if c.ConfigSyncInterval > 0 {
		opts.ConfigSyncInterval = c.ConfigSyncInterval
	}

	if c.IDListSyncInterval > 0 {
		opts.IDListSyncInterval = c.IDListSyncInterval
	}

	statsig.InitializeWithOptions(c.SecretKey, opts)
}

func GetExperiment(user statsig.User, experiment string) statsig.DynamicConfig {
	return statsig.GetExperiment(user, experiment)
}

func GetFeatureFlag(user statsig.User, flag string) statsig.FeatureGate {
	return statsig.GetGate(user, flag)
}

func Shutdown() {
	statsig.Shutdown()
}
