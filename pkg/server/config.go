package server

import (
	"fmt"
	"net/http"

	cbuilder "github.com/scribd/go-sdk/internal/pkg/configuration/builder"
)

const (
	matchAllPaths = "*"
)

type (
	// Config represents a web server configuration
	Config struct {
		Host     string `mapstructure:"host"`
		GRPCPort string `mapstructure:"grpc_port"`
		HTTPPort string `mapstructure:"http_port"`

		Cors Cors `mapstructure:"cors"`
	}

	// Cors struct represents a flag indicating if CORS feature is enabled or not
	// Also, Cors struct contains a list of CORS Settings
	Cors struct {
		Enabled bool `mapstructure:"enabled"`

		Settings []CorsSetting `mapstructure:"settings"`
	}

	// CorsSetting struct contains CORS settings.
	CorsSetting struct {
		// Path represents a server route string, for example "/example/{id}" for which the following CORS settings
		// will be applied
		Path string `mapstructure:"path"`
		// AllowedOrigins is a list of origins a cross-domain request can be executed from.
		// If the special "*" value is present in the list, all origins will be allowed.
		// An origin may contain a wildcard (*) to replace 0 or more characters
		// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
		// Only one wildcard can be used per origin.
		AllowedOrigins []string `mapstructure:"allowed_origins"`
		// AllowOriginFunc is a custom function to validate the origin. It take the origin
		// as argument and returns true if allowed or false otherwise. If this option is
		// set, the content of AllowedOrigins is ignored.
		AllowOriginFunc func(origin string) bool
		// AllowOriginRequestFunc is a custom function to validate the origin.
		// It takes the HTTP Request object and the origin as
		// argument and returns true if allowed or false otherwise.
		// If this option is set, the content of `AllowedOrigins` and `AllowOriginFunc` is ignored.
		AllowOriginRequestFunc func(r *http.Request, origin string) bool
		// AllowedHeaders is list of non simple headers the client is allowed to use with
		// cross-domain requests.
		// If the special "*" value is present in the list, all headers will be allowed.
		// "Origin" is always appended to the list.
		AllowedMethods []string `mapstructure:"allowed_methods"`
		// AllowedHeaders is list of non simple headers the client is allowed to use with
		// cross-domain requests.
		// If the special "*" value is present in the list, all headers will be allowed.
		// "Origin" is always appended to the list.
		AllowedHeaders []string `mapstructure:"allowed_headers"`
		// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
		// API specification
		ExposedHeaders []string `mapstructure:"exposed_headers"`
		// AllowCredentials indicates whether the request can include user credentials like
		// cookies, HTTP authentication or client side SSL certificates.
		AllowCredentials bool `mapstructure:"allow_credentials"`
		// MaxAge indicates how long (in seconds) the results of a preflight request
		// can be cached
		MaxAge int `mapstructure:"max_age"`
		// AllowCredentials indicates whether the request can include user credentials like
		// cookies, HTTP authentication or client side SSL certificates.
		OptionsPassthrough bool `mapstructure:"options_passthrough"`
	}
)

// NewConfig returns a new ServerConfig instance
func NewConfig() (*Config, error) {
	config := &Config{}
	viperBuilder := cbuilder.New("server")

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("Unable to decode into struct: %s", err.Error())
	}

	return config, nil
}

// GetCorsSettings returns list of CORS settings
func (c *Config) GetCorsSettings() []CorsSetting {
	return c.Cors.Settings
}

// Matches returns true if the provided path string equals to Path setting or equals
// to "*". Returns false otherwise
func (s CorsSetting) Matches(path string) bool {
	if path == s.Path || s.Path == matchAllPaths {
		return true
	}

	return false
}
