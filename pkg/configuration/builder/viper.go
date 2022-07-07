package builder

import (
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const defaultConfType = "yaml"

// viperBuilder is a Builder to streamline Viper configuration and building.
// implements builder.Builder.
type viperBuilder struct {
	vConf *viper.Viper
	name  string
}

// NewViper initializes and returns a new ViperBuilder.
// configDir is the reletive path for configuration directory.
func NewViper(configDir, appName, appEnv, appRoot string) (Builder, error) {
	vConf := viper.New()
	vConf.SetDefault("ENV", appEnv)
	vConf.SetDefault("NAME", appName)

	if err := vConf.BindEnv("ENV", "APP_ENV"); err != nil {
		log.Fatalf("Could not bind ENV for APP_ENV")
	}

	vConf.AddConfigPath(path.Join(appRoot, configDir))
	vConf.SetConfigType(defaultConfType)

	return &viperBuilder{vConf: vConf}, nil
}

// SetConfigName sets the path argument as the Viper config path.
func (vb *viperBuilder) SetConfigName(name string) {
	vb.vConf.SetConfigName(name)
	vb.name = name
}

// SetDefault sets a default value for a configuration key.
// Any default value set will be available in the `viper.Viper` configuration
// instance that is returned after calling the `Build()` function.
func (vb *viperBuilder) SetDefault(key string, value interface{}) {
	vb.vConf.SetDefault(key, value)
}

// Build creates a Conf based on viper builder.
// It first extracts a Viper instance for the specific environment it's running
// in, then explicitly calls BindEnv for each of the attributes of the
// configuration. This is done to force Viper to be aware of the ENV variables
// for each of those configuration attributes. The Viper instance returned by
// this function can be unmarshalled by the caller in a configuration-specific
// type while respecting the precedence order.
func (vb *viperBuilder) Build() (Conf, error) {
	if vb.name == "" {
		return nil, fmt.Errorf("viperBuilder.SetConfigName() must be called.")
	}

	if err := vb.vConf.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("readInConfig from viper. err: %w", err)
	}

	env := vb.vConf.GetString("ENV")
	vb.vConf = vb.vConf.Sub(env)
	if vb.vConf == nil {
		return nil, fmt.Errorf("No %s configuration for ENV %s", vb.name, env)
	}

	vb.vConf.Set("ENV", env)
	vb.vConf.SetEnvPrefix(fmt.Sprintf("APP_%s", strings.ToUpper(vb.name)))
	vb.vConf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vb.vConf.AutomaticEnv()

	return newViperConfig(vb.vConf), nil
}

// viperConfig implements Builder.Conf.
type viperConfig struct {
	viperConfig *viper.Viper
}

func newViperConfig(v *viper.Viper) Conf {
	return &viperConfig{viperConfig: v}
}

// Unmarshal returns vipers unmarshaller.
func (vu *viperConfig) Unmarshal(i interface{}) error {
	return vu.viperConfig.Unmarshal(i)
}

// Bool returns a key's value as bool.
func (c *viperConfig) Bool(key string) bool {
	return c.viperConfig.GetBool(key)
}

// Float64 returns a key's value as float64.
func (c *viperConfig) Float64(key string) float64 {
	return c.viperConfig.GetFloat64(key)
}

// Int returns a key's value as int.
func (c *viperConfig) Int(key string) int {
	return c.viperConfig.GetInt(key)
}

// String returns a key's value as string.
func (c *viperConfig) String(key string) string {
	return c.viperConfig.GetString(key)
}

// StringMap returns a key's value as map[string]interface{}.
func (c *viperConfig) StringMap(key string) map[string]interface{} {
	return c.viperConfig.GetStringMap(key)
}

// StringMapString returns a key's value as map[string]string.
func (c *viperConfig) StringMapString(key string) map[string]string {
	return c.viperConfig.GetStringMapString(key)
}

// StringSlice returns a key's value as []string.
func (c *viperConfig) StringSlice(key string) []string {
	return c.viperConfig.GetStringSlice(key)
}

// Time returns a key's value as time.Time.
func (c *viperConfig) Time(key string) time.Time {
	return c.viperConfig.GetTime(key)
}

// Duration returns a key's value as time.Duration.
func (c *viperConfig) Duration(key string) time.Duration {
	return c.viperConfig.GetDuration(key)
}

// Set sets a value to a key.
func (c *viperConfig) Set(key string, value interface{}) {
	c.viperConfig.Set(key, value)
}

// IsSet checks if the key has assigned value.
func (c *viperConfig) IsSet(key string) bool {
	return c.viperConfig.IsSet(key)
}

// AllSettings returns all settings as map.
func (c *viperConfig) AllSettings() map[string]interface{} {
	return c.viperConfig.AllSettings()
}
