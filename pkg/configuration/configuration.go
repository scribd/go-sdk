package configuration

import (
	"fmt"
	"log"
	"os"

	"github.com/scribd/go-sdk/pkg/configuration/apps"
	"github.com/scribd/go-sdk/pkg/configuration/builder"
)

const (
	defaultEnv       = "development"
	defaultConfigDir = "config"

	appRootEnvKey = "APP_ROOT"
	appNameEnvKey = "APP_NAME"
	envEnvKey     = "APP_ENV"
)

// Configuration is an app-wide configuration
type Configuration struct {
	AppName string
	AppEnv  string
	AppRoot string

	Database        apps.Database
	Server          apps.Server
	Instrumentation apps.Instrumentation
	Logger          apps.Logger
	Tracking        apps.Tracking
	PubSub          apps.PubSub
}

// NewConfig returns a new Config instance in default location.
func NewConfig(configs ...Configurable) (*Configuration, error) {
	return NewConfigWithPath(defaultConfigDir, configs...)
}

// NewConfigWithPath returns a new Config instance within a directory.
// configDir is relative path of config directory.
func NewConfigWithPath(confDir string, configurables ...Configurable) (*Configuration, error) {
	appName := getAppNameMust()
	appRoot := getAppRootMust()
	appEnv := getEnvMust()

	config := &Configuration{
		AppName: appName,
		AppEnv:  appEnv,
		AppRoot: appRoot,
	}

	for i := range configurables {
		// using viper as main builder.
		builder, err := builder.NewViper(confDir, appName, appEnv, appRoot)
		if err != nil {
			return nil, fmt.Errorf("new viper builder. err: %w", err)
		}

		if err := configurables[i](builder, config); err != nil {
			return nil, fmt.Errorf("applying configurable. err: %w", err)
		}
	}

	return config, nil
}

func getAppNameMust() string {
	appRoot := os.Getenv(appNameEnvKey)
	if appRoot == "" {
		log.Fatalf("env key %s missing", appNameEnvKey)
	}

	return appRoot
}

func getAppRootMust() string {
	appRoot := os.Getenv(appRootEnvKey)
	if appRoot == "" {
		log.Fatalf("env key %s missing", appRootEnvKey)
	}

	return appRoot
}

func getEnvMust() string {
	appRoot := os.Getenv(envEnvKey)
	if appRoot == "" {
		return defaultEnv
	}

	return appRoot
}
