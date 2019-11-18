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

// Bool returns a key's value as bool.
func (c *Config) Bool(key string) bool {
	return c.vConf.GetBool(key)
}

// Float64 returns a key's value as float64.
func (c *Config) Float64(key string) float64 {
	return c.vConf.GetFloat64(key)
}

// Int returns a key's value as int.
func (c *Config) Int(key string) int {
	return c.vConf.GetInt(key)
}

// String returns a key's value as string.
func (c *Config) String(key string) string {
	return c.vConf.GetString(key)
}

// StringMap returns a key's value as map[string]interface{}.
func (c *Config) StringMap(key string) map[string]interface{} {
	return c.vConf.GetStringMap(key)
}

// StringMapString returns a key's value as map[string]string.
func (c *Config) StringMapString(key string) map[string]string {
	return c.vConf.GetStringMapString(key)
}

// StringSlice returns a key's value as []string.
func (c *Config) StringSlice(key string) []string {
	return c.vConf.GetStringSlice(key)
}

// Time returns a key's value as time.Time.
func (c *Config) Time(key string) time.Time {
	return c.vConf.GetTime(key)
}

// Duration returns a key's value as time.Duration.
func (c *Config) Duration(key string) time.Duration {
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
