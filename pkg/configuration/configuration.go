package configuration

import (
	"fmt"

	app "github.com/scribd/go-sdk/pkg/app"
	"github.com/scribd/go-sdk/pkg/aws"
	"github.com/scribd/go-sdk/pkg/cache"
	database "github.com/scribd/go-sdk/pkg/database"
	instrumentation "github.com/scribd/go-sdk/pkg/instrumentation"
	logger "github.com/scribd/go-sdk/pkg/logger"
	"github.com/scribd/go-sdk/pkg/pubsub"
	server "github.com/scribd/go-sdk/pkg/server"
	"github.com/scribd/go-sdk/pkg/statsig"
	tracking "github.com/scribd/go-sdk/pkg/tracking"
)

// Config is an app-wide configuration
type Config struct {
	App             *app.Config
	Database        *database.Config
	Instrumentation *instrumentation.Config
	Logger          *logger.Config
	Server          *server.Config
	Tracking        *tracking.Config
	PubSub          *pubsub.Config
	Cache           *cache.Config
	AWS             *aws.Config
	Statsig         *statsig.Config
}

// NewConfig returns a new Config instance
func NewConfig() (*Config, error) {
	var errGroup error
	config := &Config{}

	appConfig, err := app.NewDefaultConfig()
	if err != nil {
		errGroup = wrapErrors(errGroup, fmt.Errorf("default config err: %w", err))
	}

	dbConfig, err := database.NewConfig()
	if err != nil {
		errGroup = wrapErrors(errGroup, fmt.Errorf("database config err: %w", err))
	}

	instrumentationConfig, err := instrumentation.NewConfig()
	if err != nil {
		errGroup = wrapErrors(errGroup, fmt.Errorf("instrumentation config err: %w", err))
	}

	loggerConfig, err := logger.NewConfig()
	if err != nil {
		errGroup = wrapErrors(errGroup, fmt.Errorf("logger config err: %w", err))
	}

	serverConfig, err := server.NewConfig()
	if err != nil {
		errGroup = wrapErrors(errGroup, fmt.Errorf("server config err: %w", err))
	}

	trackingConfig, err := tracking.NewConfig()
	if err != nil {
		errGroup = wrapErrors(errGroup, fmt.Errorf("tracking config err: %w", err))
	}

	pubsubConfig, err := pubsub.NewConfig()
	if err != nil {
		errGroup = wrapErrors(errGroup, fmt.Errorf("pubsub config err: %w", err))
	}

	cacheConfig, err := cache.NewConfig()
	if err != nil {
		errGroup = wrapErrors(errGroup, fmt.Errorf("cache config err: %w", err))
	}

	awsConfig, err := aws.NewConfig()
	if err != nil {
		errGroup = wrapErrors(errGroup, fmt.Errorf("aws config err: %w", err))
	}

	statsigConfig, err := statsig.NewConfig()
	if err != nil {
		errGroup = wrapErrors(errGroup, fmt.Errorf("statsig config err: %w", err))
	}

	config.App = appConfig
	config.Database = dbConfig
	config.Instrumentation = instrumentationConfig
	config.Logger = loggerConfig
	config.Server = serverConfig
	config.Tracking = trackingConfig
	config.PubSub = pubsubConfig
	config.Cache = cacheConfig
	config.AWS = awsConfig
	config.Statsig = statsigConfig

	return config, errGroup
}

func wrapErrors(baseError, err error) error {
	if baseError == nil {
		return err
	}

	return fmt.Errorf("%s. %w", baseError.Error(), err)
}
