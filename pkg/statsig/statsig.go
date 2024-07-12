package statsig

import (
	statsigsdk "github.com/statsig-io/go-sdk"
)

type (
	GetFeatureFlagFunc func(user statsigsdk.User) bool
	GetExperimentFunc  func(user statsigsdk.User) statsigsdk.DynamicConfig
)

func Initialize(c *Config) {
	opts := &statsigsdk.Options{
		Environment: statsigsdk.Environment{Tier: c.environment},
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

	statsigsdk.InitializeWithOptions(c.StatsigDSN, opts)
}

func GetExperiment(gate string) GetExperimentFunc {
	return func(user statsigsdk.User) statsigsdk.DynamicConfig {
		return statsigsdk.GetExperiment(user, gate)
	}
}

func GetFeatureFlag(flag string) GetFeatureFlagFunc {
	return func(user statsigsdk.User) bool {
		return statsigsdk.GetGate(user, flag).Value
	}
}

func Shutdown() {
	statsigsdk.Shutdown()
}
