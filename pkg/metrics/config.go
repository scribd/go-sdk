package metrics

type Config struct {
	Environment string
	App         string
}

// NewConfig returns a new Config instance.
func NewConfig(environment string, app string) (*Config, error) {
	config := &Config{
		Environment: environment,
		App:         app,
	}

	return config, nil
}
