package builder

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

// ViperBuilder is a builder to streamline Viper configuration and building.
type ViperBuilder struct {
	vConf    *viper.Viper
	defaults map[string]string
	name     string
}

// New initializes and returns a new ViperBuilder.
func New(name string) *ViperBuilder {
	vConf := viper.New()

	vConf.SetDefault("APP_ENV", "development")
	if err := vConf.BindEnv("APP_ENV"); err != nil {
		log.Fatalf("Could not bind ENV for APP_ENV")
	}
	vConf.SetConfigName(name)
	vConf.AddConfigPath(path.Join(os.Getenv("APP_ROOT"), "config"))
	vConf.SetConfigType("yaml")

	return &ViperBuilder{
		vConf:    vConf,
		name:     name,
		defaults: make(map[string]string),
	}
}

// ConfigPath sets the path argument as the Viper config path.
func (vb *ViperBuilder) ConfigPath(path string) *ViperBuilder {
	vb.vConf.AddConfigPath(path)
	return vb
}

// SetDefault sets a default value for a configuration key.
// Any default value set will be available in the `viper.Viper` configuration
// instance that is returned after calling the `Build()` function.
func (vb *ViperBuilder) SetDefault(key string, value string) *ViperBuilder {
	vb.defaults[key] = value
	return vb
}

// Build builds the Viper config and returns it.
// It first extracts a Viper instance for the specific environment it's running
// in, then explicitly calls BindEnv for each of the attributes of the
// configuration. This is done to force Viper to be aware of the ENV variables
// for each of those configuration attributes. The Viper instance returned by
// this function can be unmarshalled by the caller in a configuration-specific
// type while respecting the precedence order.
func (vb *ViperBuilder) Build() (*viper.Viper, error) {
	if err := vb.vConf.ReadInConfig(); err != nil {
		return nil, err
	}

	env := vb.vConf.GetString("APP_ENV")
	vb.vConf = vb.vConf.Sub(env)
	if vb.vConf == nil {
		return nil, fmt.Errorf("No %s configuration for ENV %s", vb.name, env)
	}

	if err := vb.vConf.BindEnv("APP_ENV"); err != nil {
		return nil, fmt.Errorf("Could not bind ENV for APP_ENV")
	}

	vb.vConf.SetEnvPrefix(fmt.Sprintf("APP_%s", strings.ToUpper(vb.name)))
	vb.vConf.AutomaticEnv()

	// Bind the ENV values manually.
	//
	// Workaround for `Unmarshal` which doesn't respect the environment
	// variables when loading the values.
	//
	// See: https://github.com/spf13/viper/issues/761
	for _, k := range vb.vConf.AllKeys() {
		if err := vb.vConf.BindEnv(strings.ToUpper(k)); err != nil {
			return nil, fmt.Errorf("Could not configure %s for ENV %s", k, env)
		}

	}

	for key, val := range vb.defaults {
		vb.vConf.SetDefault(key, val)
	}

	return vb.vConf, nil
}
