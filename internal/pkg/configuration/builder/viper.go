package builder

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

// ViperBuilder is a builder to streamline Viper configuration and building.
type ViperBuilder struct {
	vConf *viper.Viper
}

// New initializes and returns a new ViperBuilder.
func New() *ViperBuilder {
	vConf := viper.New()

	vConf.SetDefault("APP_ENV", "development")
	vConf.AddConfigPath(path.Join(os.Getenv("APP_ROOT"), "config"))
	vConf.SetConfigType("yaml")
	vConf.SetEnvPrefix("APP")
	vConf.AutomaticEnv()

	return &ViperBuilder{
		vConf: vConf,
	}
}

// ConfigPath sets the path argument as the Viper config path.
func (vb *ViperBuilder) ConfigPath(path string) *ViperBuilder {
	vb.vConf.AddConfigPath(path)
	return vb
}

// ConfigName sets the name argument as the Viper config name.
func (vb *ViperBuilder) ConfigName(name string) *ViperBuilder {
	vb.vConf.SetConfigName(name)
	return vb
}

// SetDefault sets a default value for a configuration key.
func (vb *ViperBuilder) SetDefault(key string, value string) *ViperBuilder {
	vb.vConf.SetDefault(key, value)
	return vb
}

// Build builds the Viper config and returns it.
func (vb *ViperBuilder) Build() (*viper.Viper, error) {
	if err := vb.vConf.ReadInConfig(); err != nil {
		return nil, err
	}

	return vb.vConf.Sub(vb.vConf.GetString("APP_ENV")), nil
}
