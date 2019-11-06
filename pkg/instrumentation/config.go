package instrumentation

import (
	"fmt"

	cbuilder "git.lo/microservices/sdk/go-sdk/internal/pkg/configuration/builder"
)

// DataDog Agent, APP and logs available configurations.
//
// References:
// - https://docs.datadoghq.com/agent/docker/
// - https://docs.datadoghq.com/agent/docker/apm/
// - https://github.com/DataDog/datadog-agent/tree/master/Dockerfiles/agent#environment-variables
// - https://github.com/DataDog/datadog-agent/blob/master/pkg/config/config_template.yaml
// - https://github.com/DataDog/datadog-agent/blob/master/pkg/trace/config/env.go
//
// DD_AGENT_HOST            // The host of Dogstatsd and traces.
// DD_API_KEY               // <YOUR_DATADOG_API_KEY>
// DD_APM_ENABLED           // Set to true to enable the APM Agent.
// DD_APM_ENV	            // Sets the default environment for your traces.
// DD_APM_MAX_EPS	    // Sets the maximum Analyzed Spans per second. Default value: 200
// DD_APM_MAX_TPS	    // Sets the maximum traces per second. Default value: 10
// DD_APM_NON_LOCAL_TRAFFIC // Allow non-local traffic when tracing from other containers.
// DD_APM_RECEIVER_PORT     // Port that the Datadog Agent’s trace receiver listens on. Default value: 8126.
// DD_BIND_HOST             // Set the StatsD & receiver hostname.
// DD_DOGSTATSD_PORT        // The Agent DogStatsD port. The default: 8125.
// DD_DOGSTATSD_SOCKET      // Path to the unix socket to listen to.
// DD_DOGSTATSD_TAGS        // Additional tags to append to all metrics, events received by this DogStatsD server.
// DD_HISTOGRAM_AGGREGATES  // Default value: "max median avg count".
// DD_HISTOGRAM_PERCENTILES // Default value: "0.95".
// DD_LOGS_ENABLED          // Enables log collection when set to true.
// DD_LOG_LEVEL             // Set the logging level (trace/debug/info/warn/error/critical/off)
// DD_TRACE_AGENT_PORT      // The default is 8126

type Config struct {
	Enabled bool `mapstructure:"enabled"`
}

// NewConfig returns a new ServerConfig instance.
func NewConfig() (*Config, error) {
	config := &Config{}
	viperBuilder := cbuilder.New("datadog")

	vConf, err := viperBuilder.Build()
	if err != nil {
		return config, err
	}

	if err = vConf.Unmarshal(config); err != nil {
		return config, fmt.Errorf("Unable to decode into struct: %s", err.Error())
	}

	return config, nil
}
