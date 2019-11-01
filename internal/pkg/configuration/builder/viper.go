package builder

import (
	"os"
	"path"

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

	vConf.SetDefault("ENV", "development")
	vConf.SetConfigName(name)
	vConf.AddConfigPath(path.Join(os.Getenv("APP_ROOT"), "config"))
	vConf.SetConfigType("yaml")
	vConf.SetEnvPrefix("APP")
	vConf.AutomaticEnv()

	return &ViperBuilder{
		vConf: vConf,
		name:  name,
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
// Any default value should be set using the `SetDefault()` function of
// a `ViperBuilder` instance, before calling `Build()`.
func (vb *ViperBuilder) Build() (*viper.Viper, error) {
	if err := vb.vConf.ReadInConfig(); err != nil {
		return nil, err
	}

	vb.vConf = vb.vConf.Sub(vb.vConf.GetString("ENV"))
	vb.vConf.SetConfigName(vb.name)

	for key, val := range vb.defaults {
		vb.vConf.SetDefault(key, val)
	}

	return vb.vConf, nil
}
