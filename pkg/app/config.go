package app

import (
	"os"
	"path"
	"time"

	cbuilder "git.lo/microservices/sdk/go-sdk/internal/pkg/configuration/builder"

	"github.com/spf13/viper"
)

// Config is custom application configuration.
type Config struct {
	vConf *viper.Viper
}

// NewDefaultConfig returns a new Config with default values
func NewDefaultConfig() (*Config, error) {
	return NewConfig(
		path.Join(os.Getenv("APP_ROOT"), "config"),
		"settings",
	)
}

// NewConfig sets up the app configuration, setting default values and configurations.
func NewConfig(configPath string, configName string) (*Config, error) {
	conf := &Config{}
	viperBuilder := cbuilder.New(configName).ConfigPath(configPath)

	vConf, err := viperBuilder.Build()
	if err != nil {
		return nil, err
	}

	conf.vConf = vConf
	return conf, nil
}

// GetBool returns a key's value as bool.
func (c *Config) GetBool(key string) bool {
	return c.vConf.GetBool(key)
}

// GetFloat64 returns a key's value as float64.
func (c *Config) GetFloat64(key string) float64 {
	return c.vConf.GetFloat64(key)
}

// GetInt returns a key's value as int.
func (c *Config) GetInt(key string) int {
	return c.vConf.GetInt(key)
}

// GetString returns a key's value as string.
func (c *Config) GetString(key string) string {
	return c.vConf.GetString(key)
}

// GetStringMap returns a key's value as map[string]interface{}.
func (c *Config) GetStringMap(key string) map[string]interface{} {
	return c.vConf.GetStringMap(key)
}

// GetStringMapString returns a key's value as map[string]string.
func (c *Config) GetStringMapString(key string) map[string]string {
	return c.vConf.GetStringMapString(key)
}

// GetStringSlice returns a key's value as []string.
func (c *Config) GetStringSlice(key string) []string {
	return c.vConf.GetStringSlice(key)
}

// GetTime returns a key's value as time.Time.
func (c *Config) GetTime(key string) time.Time {
	return c.vConf.GetTime(key)
}

// GetDuration returns a key's value as time.Duration.
func (c *Config) GetDuration(key string) time.Duration {
	return c.vConf.GetDuration(key)
}

// Set sets a value to a key.
func (c *Config) Set(key string, value interface{}) {
	c.vConf.Set(key, value)
}

// IsSet checks if the key has assigned value.
func (c *Config) IsSet(key string) bool {
	return c.vConf.IsSet(key)
}

// AllSettings returns all settings as map.
func (c *Config) AllSettings() map[string]interface{} {
	return c.vConf.AllSettings()
}
