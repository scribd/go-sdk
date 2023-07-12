# Go SDK

SDK, the Go version.

## Table Of Contents

- [Prerequisites](#prerequisites)
- [SDK functionality](#sdk-functionality)
    - [Application Configuration](#application-configuration)
        - [Predefined application-agnostic configurations](#predefined-application-agnostic-configurations)
        - [Custom application-specific configurations](#custom-application-specific-configurations)
        - [Complex values representation](#complex-values-representation)
        - [Environment-awareness](#environment-awareness)
        - [Using application configuration in tests](#using-application-configuration-in-tests)
    - [Logger](#logger)
        - [Initialization and default configuration](#initialization-and-default-configuration)
        - [Environment-awareness](#environment-awareness-1)
        - [Log levels](#log-levels)
        - [Structured logging](#structured-logging)
    - [Logging & tracing middleware](#logging---tracing-middleware)
        - [HTTP server middleware](#http-server-middleware)
        - [gRPC server interceptors](#grpc-server-interceptors)
        - [Formatting and handlers](#formatting-and-handlers)
        - [Sentry error reporting](#sentry-error-reporting)
    - [Database Connection](#database-connection)
    - [Server](#server)
        - [CORS settings](#cors-settings)
        - [CORS middleware](#cors-middleware)
    - [ORM Integration](#orm-integration)
        - [Usage of ORM](#usage-of-orm)
    - [PubSub](#pubsub)
        - [Kafka specific configuration](#kafka-specific-configuration)
- [APM & Instrumentation](#apm---instrumentation)
    - [Request ID middleware](#request-id-middleware)
        - [HTTP server Request ID middleware](#http-server-request-id-middleware)
        - [gRPC server Request ID interceptors](#grpc-server-request-id-interceptors)
    - [HTTP Router Instrumentation](#http-router-instrumentation)
    - [gRPC server and client interceptors](#grpc-server-and-client-interceptors)
    - [Database instrumentation & ORM logging](#database-instrumentation---orm-logging)
        - [HTTP server middleware example](#http-server-middleware-example)
        - [gRPC server interceptors example](#grpc-server-interceptors-example)
    - [AWS Session instrumentation](#aws-session-instrumentation)
    - [PubSub instrumentation and logging](#pubsub-instrumentation-and-logging)
      - [Kafka](#kafka)
    - [Profiling](#profiling)
    - [Custom Metrics](#custom-metrics)
- [Using the `go-sdk` in isolation](#using-the--go-sdk--in-isolation)
- [Developing the SDK](#developing-the-sdk)
    - [Building the docker environment](#building-the-docker-environment)
    - [Running tests within the docker environment](#running-tests-within-the-docker-environment)
    - [Entering the docker environment](#entering-the-docker-environment)
    - [Using a development build of `go-sdk`](#using-a-development-build-of--go-sdk-)
    - [Commit messages](#commit-messages)
- [Release](#release)
- [Maintainers](#maintainers)

## Prerequisites

* [Go](https://golang.org) (version `1.20.6`).
* [Docker](https://www.docker.com/) (version `19.03.2`).

## SDK functionality

The SDK provides the following out-of-the-box functionality:

* Application Configuration & Environment-awareness
* Logger

### Application Configuration

The SDK provides the framework for managing your application's configuration.
As it uses the excellent [Viper](https://github.com/spf13/viper) under the
hood, `go-sdk` knows how to read configuration files (**YAML only**) and `ENV`
variables. `go-sdk` facilitates usage of predefined static configurations and
custom dynamic configurations. Lastly, the `go-sdk` is environment-aware. This
means that it supports environment scoped configuration in the configuration
files out-of-the-box.

:warning: Environment variables will not be picked up by Viper if the keys are
not present in the respective YAML file. This is due to Viper loading the
configuration by reading a YAML file **first**, and binding the respective ENV
variables on the fly.

#### Predefined application-agnostic configurations

The predefined configurations have an associated type and expect their
respective configuration to be present in a file corresponding to their name.

The list of predefined top-level configurations:
* `Server`, containing the application server configurations, expects a
  `config/server.yml` configuration file.
* `Logger`, containing the application logger configuration, expects a
  `config/logger.yml` configuration file.
* `Database`, containing the application database configuration, expects a
  `config/database.yml` configuration file.
* `Redis`, containing the application Redis configuration, expects a
  `config/redis.yml` configuration file. (TBD)
* `PubSub`, containing the application pubsub configuration, expects a
  `config/pubsub.yml` configuration file.
* <insert new configuration here>

For example, to get the host and the port at which the HTTP server listens on
in an application:

```go
host     := sdk.Config.Server.Host
httpPort := sdk.Config.Server.HTTPPort
```

#### Custom application-specific configurations

Additionally, the `go-sdk` allows developers to add custom configuration data
in a `go-chassis`-powered application, through the `config/settings.yml` file.
As these configurations are custom and their type cannot (read: shouldn't) be
inferred by the SDK, it exposes a familiar Viper-like interface for interaction
with them.

All of the application-specific configurations can be accessed through the
`sdk.Config.App` object. For example, given a `config/settings.yml` with the
following contents:

```yaml
# config/settings.yml
common: &common
  name: "my-awesome-app"

development:
  <<: *common

staging:
  <<: *common

production:
  <<: *common
```

To get the application name from the configuration one would use the following
statement:

```go
applicationName := sdk.Config.App.GetString("name")
```

The configuration variables can be overridden by a corresponding environment
variable; these variables must have the following format:

```
APP_SETTINGS_NAME=my-really-awesome-app
^^^ ^^^^^^^^ ^^^^ ^^^^^^^^^^^^^^^^^^^^^
|   |        |    + ----------- variable_value
|   |        + ---------------- variable_name
|   + ------------------------- config_file
+ ----------------------------- prefix
```

The environment variable has the precedence over the configuration file.

#### Complex values representation

Since `yaml` data type support is much richer than environment variables we have to take extra care if we want to
override `yaml` complex data types such as `list` and `dictionary`.

To represent the `list` we can use space-separated values:

```bash
APP_SETTINGS_NAME="value1 value2"
```

And then, on the application side we have to convert it to a string slice manually:

```go
// if the value is the part of `settings.yml`
stringSlice := sdk.Config.App.StringSlice("name")

// or calling the Viper's API directly
stringSlice := viper.GetStringSlice("app.settings.name")
```

To represent the `dictionary` we can use a JSON:

```bash
APP_SETTINGS_NAME="{\"key\":\"value\"}"
```

And then, on the application side we have to convert it to a string map manually:

```go
// if the value is the part of `settings.yml`
stringMap := sdk.Config.App.StringMapString("name")

// or calling the Viper's API directly
stringMap := viper.GetStringMapString("app.settings.name")
```

#### Environment-awareness

Application configurations vary between environments, therefore the `go-sdk` is
environment-aware. This means that any application that uses the `go-sdk`
should have an `APP_ENV` environment variable set. Applications built using the
`go-chassis` have this out of the box. In absence of the `APP_ENV` variable,
it's defaulted to `development`.

This means that any configuration files, placed in the `config` directory,
should contain environment namespacing. For example:

```yaml
# config/logger.yml
common: &common
  console_enabled: true
  console_json_format: true
  console_level: "info"
  file_enabled: false
  file_json_format: true
  file_level: "trace"

development:
  <<: *common
  console_level: "debug"
  file_enabled: true
  file_level: "debug"

staging:
  <<: *common

production:
  <<: *common
```

When run in a particular environment, by setting `APP_ENV`, it will load the
respective section from the YAML file. For example, when the application is
loaded in `development` environment, `go-sdk` will automatically load the
values from the `development` section. This convention is applied to all
configurations supported by `go-sdk`.

#### Using application configuration in tests

When an application is using the SDK to load configurations, that includes
loading application configurations in test environment as well. This means that
due to its [environment-aware](#environment-awareness) nature, by default the
SDK will load the configuration from the YAML files in the `test` namespace.

This provides two ways to configure your application for testing:
1. By adding a `test` namespace to the respective YAML file and adding the test
   configuration there, or
2. Using the provided structs and their setters/getters from the test files
   themselves.

Adding a `test` namespace to any configuration file, looks like:

```yaml
# config/logger.yml
common: &common
  console_enabled: true
  console_json_format: true
  console_level: "info"
  file_enabled: false
  file_json_format: true
  file_level: "trace"

development:
  <<: *common

test:
  <<: *common
  console_level: "debug"
  file_enabled: true
  file_level: "debug"

staging:
  <<: *common

production:
  <<: *common
```

Given the configuration above, the SDK will load out-of-the-box the `test`
configuration and apply it to the logger.

In cases where we want to modify the configuration of an application in a test
file, we can simply use the constructor that `go-sdk` provides:

```go
package main

import (
    "testing"

    sdkconfig "github.com/scribd/go-sdk/pkg/configuration"
}

var (
    // Config is SDK-powered application configuration.
    Config *sdkconfig.Config
)

func SomeTest(t *testing.T) {
    if Config, err = sdkconfig.NewConfig(); err != nil {
        log.Fatalf("Failed to load SDK config: %s", err.Error())
    }

    // Change application settings
    Config.App.Set("key", "value")

    // Continue with testing...
    setting := Config.App.GetString("key") // returns "value"
}
```

### Logger

`go-sdk` ships with a logger, configured with sane defaults out-of-the-box.
Under the hood, it uses the popular
[logrus](https://github.com/sirupsen/logrus) as a logger, in combination with
[lumberjack](https://github.com/natefinch/lumberjack) for log rolling. It
supports log levels, formatting (plain or JSON), structured logging and
storing the output to a file and/or `STDOUT`.

#### Initialization and default configuration

`go-sdk` comes with a configuration built-in for the logger. That means that
initializing a new logger is as easy as:

```go
package main

import (
	"log"

	sdklogger "github.com/scribd/go-sdk/pkg/logger"
)

var (
	// Logger is SDK-powered application logger.
	Logger sdklogger.Logger
	err    error
)

func main() {
	if loggerConfig, err := sdklogger.NewConfig(); err != nil {
		log.Fatalf("Failed to initialize SDK logger configuration: %s", err.Error())
	}

	if Logger, err = sdklogger.NewBuilder(loggerConfig).Build(); err != nil {
		log.Fatalf("Failed to load SDK logger: %s", err.Error())
	}
}
```

The logger initialized with the default configuration will use the `log/`
directory in the root of the project to save the log files, with the name of
the current application environment and a `.log` extension.

#### Environment-awareness

Much like with the application configuration, the logger follows the convention
of loading a YAML file placed in the `config/logger.yml` file. This means that
any project which imports the `logger` package from `go-sdk` can place the
logger configuration in their `config/logger.yml` file and the `go-sdk` will
load that configuration when initializing the logger.

Since the logger is also environment-aware, it will assume the presence of the
`APP_ENV` environment variable and use it to set the name of the log file to
the environment name. For example, an application running with
`APP_ENV=development` will have its log entries in `log/development.log` by
default.

Also, it expects the respective `logger.yml` file to be environment namespaced.
For example:

```yaml
# config/logger.yml
common: &common
  console_enabled: true
  console_json_format: false
  console_level: "debug"
  file_enabled: true
  file_json_format: false
  file_level: "debug"

development:
  <<: *common
  file_enabled: false

test:
  <<: *common
  console_enabled: false

staging:
  <<: *common
  console_json_format: true
  console_level: "debug"
  file_enabled: false

production:
  <<: *common
  console_json_format: true
  console_level: "info"
  file_enabled: false
```

Given the configuration above, the logger package will load out-of-the-box the
configuration for the respective environment and apply it to the logger
instance.

#### Log levels

The SDK's logger follows best practices when it comes to logging levels. It
exposes multiple log levels, in order from lowest to highest:

* Trace, invoked with `logger.Tracef`
* Debug, invoked with `logger.Debugf`
* Info, invoked with `logger.Infof`
* Warn, invoked with `logger.Warnf`
* Error, invoked with `logger.Errorf`
* Fatal, invoked with `logger.Fatalf`
* Panic, invoked with `logger.Panicf`

Each of these log levels will produce log entries, while two of them have
additional functionality:

* `logger.Fatalf` will add a log entry and exit the program with error code 1
  (i.e. `exit(1)`)
* `logger.Panicf` will add a log entry and invoke `panic` with all of the
  arguments passed to the `Panicf` call

#### Structured logging

Loggers created using the `go-sdk` logger package, define a list of hardcoded
fields that every log entry will be consisted of. This is done by design, with
the goal of a uniform log line structure across all Go services that use the
`go-sdk`.

Adding more field is made possible by the `logger.WithFields` function
and by the `Builder` API:

```go
fields := map[string]string{ "role": "server" }

if Logger, err = sdklogger.NewBuilder(loggerConfig).SetFields(fields).Build(); err != nil {
	log.Fatalf("Failed to load SDK logger: %s", err.Error())
}
```

While adding more fields is easy to do, removing the three default
fields from the log lines is, by design, very hard to do and highly
discouraged.

The list of fields are:

* `level`, indicating the log level of the log line
* `message`, representing the actual log message
* `timestamp`, the date & time of the log entry in ISO 8601 UTC format

### Logging & tracing middleware

`go-sdk` ships with a `Logger` middleware. When used, it tries to retrieve the `RequestID`, `TraceID` and `SpanID`
from the incoming context. Then, middleware assigns those values to the log entries for further correlation.

#### HTTP server middleware

```go
func main() {
	loggingMiddleware := sdkmiddleware.NewLoggingMiddleware(sdk.Logger)

	httpServer := server.
		NewHTTPServer(host, httpPort, applicationEnv, applicationName).
		MountMiddleware(loggingMiddleware.Handler).
		MountRoutes(routes)
}
```

#### gRPC server interceptors

```go
func main() {
    rpcServer, err := server.NewGrpcServer(
        host,
        grpcPort,
        []grpc.ServerOption{
            grpc.ChainUnaryInterceptor(
                sdkinterceptors.TracingUnaryServerInterceptor(applicationName),
                sdkinterceptors.RequestIDUnaryServerInterceptor(),
                sdkinterceptors.LoggerUnaryServerInterceptor(logger),
            ),
            grpc.ChainStreamInterceptor(
                sdkinterceptors.TracingStreamServerInterceptor(applicationName),
                sdkinterceptors.RequestIDStreamServerInterceptor(),
                sdkinterceptors.LoggerStreamServerInterceptor(logger),
            ),
        }...)
}
```

#### Formatting and handlers

The logger ships with two different formats: a plaintext and JSON format. This
means that the log entries can have a simple plaintext format, or a JSON
format. These options are configurable using the `ConsoleJSONFormat` and
`FileJSONFormat` attributes of the logger `Config`.

An example of the plain text format:

```
timestamp="2019-10-23T15:28:54Z" level=info message="GET  HTTP/1.1 200"
```

An example of the JSON format:

```json
{"level":"info","message":"GET  HTTP/1.1 200","timestamp":"2019-10-23T15:29:26Z"}
```

The logger handles the log entries and can store them in a file or send them to
`STDOUT`. These options are configurable using the `ConsoleEnabled` and
`FileEnabled` attributes of the logger `Config`.

#### Sentry error reporting

The logger can be further instrumented to report error messages to
[Sentry](https://docs.sentry.io).

The following instructions assume that a project has been
[setup in Sentry](https://docs.sentry.io/error-reporting/quickstart/?platform=go)
and that the corresponding DSN, or Data Source Name, is available.

The respective configuration file is `sentry.yml` and it should include
the following content:

```yaml
# config/sentry.yml
common: &common
  dsn: ""

development:
  <<: *common

staging:
  <<: *common

production:
  <<: *common
  dsn: "https://<key>@sentry.io/<project>"
```

The tracking can be enabled from the `Logger Builder` with the `SetTracking`
function. To trigger tracking, use `logger` with `WithError(err)` as follows:

```go
package main

import (
	"log"
	"errors"

	sdklogger   "github.com/scribd/go-sdk/pkg/logger"
	sdktracking "github.com/scribd/go-sdk/pkg/tracking"
)

var (
	// Logger is SDK-powered application logger.
	Logger sdklogger.Logger
	err    error
)

func main() {
	if loggerConfig, err := sdklogger.NewConfig(); err != nil {
		log.Fatalf("Failed to initialize SDK logger configuration: %s", err.Error())
	}

	if trackingConfig, err := sdktracking.NewConfig(); err != nil {
		log.Fatalf("Failed to initialize SDK tracking configuration: %s", err.Error())
	}

	if Logger, err = sdklogger.NewBuilder(loggerConfig).SetTracking(trackingConfig).Build(); err != nil {
		log.Fatalf("Failed to load SDK logger: %s", err.Error())
	}

	err = errors.New("new error")
	logger.WithError(err).Errorf("trying sentry error functionality")
```

A logger build with a valid tracking configuration will automatically
report to Sentry any errors emitted from the following log levels:

- `Error`;
- `Fatal`;
- `Panic`;

The following environment variables are automatically used in the configuration
of the Sentry client to enrich the error data:

```
environment: APP_ENV
release: APP_VERSION
serverName: APP_SERVER_NAME
```

More about the
["environment" configuration](https://docs.sentry.io/enriching-error-data/environments/?platform=go)
the
["server name" configuration](https://docs.sentry.io/error-reporting/configuration/?platform=go#server-name)
and the
["release" configuration](https://docs.sentry.io/workflow/releases/?platform=go)
can be found in the Sentry documentation.

### Database Connection

`go-sdk` ships with a default setup for a database connection, built on top of
the built-in [database
configuration](#predefined-application-agnostic-configurations). The
configuration that is read from the `config/database.yml` file is then used to
create connection details, which are then used to compose a [data source
name](https://en.wikipedia.org/wiki/Data_source_name) (DSN), for example:

```
username:password@tcp(192.168.1.1:8080)/dbname?timeout=10s&charset=utf8&parseTime=True&loc=Local
```

At the moment, the database connection established using the `go-sdk` can be
only to a MySQL database. This is subject to change as the `go-sdk` evolves and
becomes more feature-complete.

The database connection can also be configured through a YAML file. `go-sdk`
expects this file to be placed at the `config/database.yml` path, within the
root of the project.

Each of these connection details can be overriden by an `ENV` variable.

| Setting  | Description                      | YAML variable | Environment variable (ENV) | Default     |
| -------- | -------------------------------- | ------------- | -------------------------- | ----------- |
| Host     | The database host                | `host`        | `APP_DATABASE_HOST`        | `localhost` |
| Port     | The database port                | `port`        | `APP_DATABASE_PORT`        | `3306`      |
| Database | The database name                | `database`    | `APP_DATABASE_DATABASE`    |             |
| Username | App user name                    | `username`    | `APP_DATABASE_USERNAME`    |             |
| Password | App user password                | `password`    | `APP_DATABASE_PASSWORD`    |             |
| Pool     | Connection pool size             | `pool`        | `APP_DATABASE_POOL`        | `5`         |
| Timeout  | Connection timeout (in seconds)  | `timeout`     | `APP_DATABASE_TIMEOUT`     | `1s`        |

An example `database.yml`:

```yaml
common: &common
  host: db
  port: 3306
  username: username
  password: password
  pool: 5
  timeout: 1s

development:
  <<: *common
  database: application_development

test:
  <<: *common
  database: application_test

production:
  <<: *common
```

### Server

`go-sdk` provides a convenient way to create a basic Server configuration.

| Setting     | Description                      | YAML variable    | Environment variable (ENV) |
| --------    | -------------------------------- | -------------    | -------------------------- |
| Host        | Server host                      | `host`           | `APP_SERVER_HOST`          |
| HTTPPort    | HTTP port                        | `http_port`      | `APP_SERVER_HTTP_PORT`     |
| GRPCPort    | gRPC port                        | `grpc_port`      | `APP_SERVER_GRPC_PORT`     |
| HTTPTimeout | HTTP related timeouts            | `http_timeout`   |                            |
| CORS        | CORS settings                    | `cors`           |                            |

An example `server.yml`:

```yaml
common: &common
  http_port: 8080
  http_timeout:
    write: 2s
    read: 1s
    idle: 90s
  cors:
    enabled: true
    settings:
      - path: "*"
        allowed_origins: ["*"]
        allowed_methods: ["GET"]
        allowed_headers: ["Allowed-Header"]
        exposed_headers: ["Exposed-Header"]
        allow_credentials: true
        max_age: 600
```

#### CORS settings

CORS stands for [Cross Origin Resource Sharing](http://www.w3.org/TR/cors/). `go-sdk` provides a basic
optional configuration for the CORS settings that are passed to the HTTP middleware.

| Setting  | Description                      | YAML variable | Environment variable (ENV) | Default |
| -------- | -------------------------------- | ------------- | -------------------------- | ------- |
| Enabled  | Whether CORS enabled or not      | `enabled`     | `APP_SERVER_CORS_ENABLED`  | false   |
| Settings | List of CORS Settings            | `settings`    |                            |         |

**PLEASE NOTE:** there is no way to specify `Settings` via environment variables as it is presented as
a nested structure. To configure the CORS, use the `server.yaml` file

For the full list of the CORS settings please refer to the inline documentation of the [server package](https://github.com/scribd/go-sdk/blob/main/pkg/server/config.go)
Also, consider looking into the documentation of the [cors library](https://github.com/rs/cors#parameters) which
currently lays under the hood of the [CORS middleware](https://github.com/scribd/go-sdk/tree/main/pkg/middleware/cors.go).

#### CORS middleware

CORS [middleware](https://github.com/scribd/go-sdk/tree/main/pkg/middleware/cors.go) is a tiny wrapper around the [cors library](https://github.com/rs/cors#parameters).
It aims to provide an extensive way to configure CORS and at the same time not bind services to a particular
implementation.

Below is an example of the CORS middleware initialization:

```go
package main

import (
	"github.com/scribd/go-sdk/pkg/server"
	"log"
)

func main() {
	config, err := server.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	corsMiddleware := middleware.NewCorsMiddleware(config.Cors.Settings[0])

	// possible Server implementation
	httpServer := server.
		NewHTTPServer(host, httpPort, applicationEnv, applicationName).
		MountMiddleware(corsMiddleware.Handler).
		MountRoutes(routes)
}
```

### ORM Integration

`go-sdk` comes with an integration with the popular
[gorm](https://gorm.io/gorm) as an object-relational mapper (ORM).
Using the configuration details, namely the [data source
name](https://en.wikipedia.org/wiki/Data_source_name) (DSN) as their product,
gorm is able to open a connection and give the `go-sdk` users a preconfigured
ready-to-use database connection with an ORM attached. This can be done as
follows:

```go
package main

import (
	sdkdb "github.com/scribd/go-sdk/pkg/database"
)

func main() {
	// Loads the database configuration.
	dbConfig, err := sdkdb.NewConfig()

	// Establishes a gorm database connection using the connection details.
	dbConn, err := sdkdb.NewConnection(dbConfig)
}
```

The connection details are handled internally by the gorm integration, in other
words the `NewConnection` function, so they remain opaque for the user.

#### Usage of ORM

Invoking the constructor for a database connection, `go-sdk` returns a
[Gorm-powered](https://gorm.io/gorm) database connection. It can be
used right away to query the database:

```go
package main

import (
	"fmt"

	sdkdb "github.com/scribd/go-sdk/pkg/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name string `gorm:"type:varchar(255)"`
	Age  uint   `gorm:"type:int"`
}

func main() {
	dbConfig, err := sdkdb.NewConfig()
	dbConn, err := sdkdb.NewConnection(dbConfig)

	user := User{Name: name, Age: age}

	errs := dbConn.Create(&user).GetErrors()
	if errs != nil {
		fmt.Println(errs)
	}
	fmt.Println(user)
}
```

To learn more about Gorm, you can start with its [official
documentation](https://gorm.io/docs/).

### PubSub

`go-sdk` provides a convenient way to create a basic [PubSub](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern) configuration.

Top-level configuration is split into sections specific for the vendor. Let's look into the sample `pubsub.yml` file:

```yaml
common: &common
  kafka:
    # Set by APP_PUBSUB_KAFKA_BROKER_URLS env variable
    broker_urls:
      - "localhost:9092"
    # Set by APP_PUBSUB_KAFKA_CLIENT_ID env variable
    client_id: "test-app"
    # Set by APP_PUBSUB_KAFKA_CERT_PEM env variable
    cert_pem: "pem string"
    # Set by APP_PUBSUB_KAFKA_CERT_PEM_KEY env variable
    cert_pem_key: "pem key"
    security_protocol: "ssl"
    ssl_verification_enabled: true
    publisher:
      # Set by APP_PUBSUB_KAFKA_PUBLISHER_MAX_ATTEMPTS env variable
      max_attempts: 3
      write_timeout: "10s"
      topic: "test-topic"
    subscriber:
      topic: "test-topic"
      group_id: ""
```

#### Kafka specific configuration

Kafka top-level configuration contains generic options that fit both publisher and subscriber:

| Setting                  | Description                                                       | YAML variable              | Environment variable (ENV)                  | Type         | Possible Values                                       |
|--------------------------|-------------------------------------------------------------------|----------------------------|---------------------------------------------|--------------|-------------------------------------------------------|
| Broker URLs              | Broker URLs to connect to                                         | `broker_urls`              | `APP_PUBSUB_KAFKA_BROKER_URLS`              | list(string) | localhost:9093                                        |
| Client ID                | Client identifier                                                 | `client_id`                | `APP_PUBSUB_KAFKA_CLIENT_ID`                | string       | go-chassis-app                                        |
| Cert PEM                 | Client's public key string (PEM format) used for authentication   | `cert_pem`                 | `APP_PUBSUB_KAFKA_CERT_PEM`                 | string       | long PEM string (deprecated)                          |
| Cert PEM Key             | Client's private key string (PEM format) used for authentication  | `cert_pem_key`             | `APP_PUBSUB_KAFKA_CERT_PEM_KEY`             | string       | long PEM key string (deprecated)                      |
| Security Protocol        | Protocol used to communicate with brokers                         | `security_protocol`        | `APP_PUBSUB_KAFKA_SECURITY_PROTOCOL`        | string       | plaintext, ssl, sasl_plaintext, sasl_ssl (deprecated) |
| SSL verification enabled | Enable OpenSSL's builtin broker (server) certificate verification | `ssl_verification_enabled` | `APP_PUBSUB_KAFKA_SSL_VERIFICATION_ENABLED` | bool         | true, false (deprecated)                              |
| Enable Metrics           | Enable publishing metrics for Kafka client broker connection      | `metrics_enabled`          | `APP_PUBSUB_KAFKA_METRICS_ENABLED`          | bool         | true, false                                           |

To break it down further `Publisher` and `Subscriber` have their own set of configuration options:

**Publisher**:

| Setting        | Description                                                                                                                     | YAML variable     | Environment variable (ENV)                   | Type   | Possible Values |
|----------------|---------------------------------------------------------------------------------------------------------------------------------|-------------------|----------------------------------------------|--------|-----------------|
| Max attempts   | Maximum amount of times to retry sending a failing Message                                                                      | `max_attempts`    | `APP_PUBSUB_KAFKA_PUBLISHER_MAX_ATTEMPTS`    | int    | 3               |
| Write timeout  | Local message timeout. This value is only enforced locally and limits the time a produced message waits for successful delivery | `write_timeout`   | `APP_PUBSUB_KAFKA_PUBLISHER_WRITE_TIMEOUT`   | string | 10s             |
| Topic          | Topic name                                                                                                                      | `topic`           | `APP_PUBSUB_KAFKA_PUBLISHER_TOPIC`           | string | topic           |
| Enable Metrics | Enable publishing metrics for Kafka producer                                                                                    | `metrics_enabled` | `APP_PUBSUB_KAFKA_PUBLISHER_METRICS_ENABLED` | bool   | true, false     |

**Subscriber**:

| Setting           | Description                                                                                                                                                                                               | YAML variable         | Environment variable (ENV)                        | Type    | Possible Values |
|-------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------|---------------------------------------------------|---------|-----------------|
| Topic             | Topic name                                                                                                                                                                                                | `topic`               | `APP_PUBSUB_KAFKA_SUBSCRIBER_TOPIC`               | string  | topic           |
| Group ID          | Client group id string. All clients sharing the same group id belong to the same group                                                                                                                    | `group_id`            | `APP_PUBSUB_KAFKA_SUBSCRIBER_GROUP_ID`            | string  | service-name    |
| Enable Metrics    | Enable publishing metrics for Kafka consumer                                                                                                                                                              | `metrics_enabled`     | `APP_PUBSUB_KAFKA_SUBSCRIBER_METRICS_ENABLED`     | bool    | true, false     |
| Enable AutoCommit | Enable AutoCommit option for Kafka consumer (default `true`)                                                                                                                                              | `auto_commit.enabled` | `APP_PUBSUB_KAFKA_SUBSCRIBER_AUTO_COMMIT_ENABLED` | bool    | true, false     |
| Number of workers | Number of workers processing incoming messages                                                                                                                                                            | `workers`             | `APP_PUBSUB_KAFKA_SUBSCRIBER_WORKERS`             | number  | 1,2...          |
| Block rebalance   | Indicates if the subscriber will block rebalance event when polling. For more information refer to the [franz-go documentation](https://pkg.go.dev/github.com/twmb/franz-go/pkg/kgo#BlockRebalanceOnPoll) | `block_rebalance`     | `APP_PUBSUB_KAFKA_SUBSCRIBER_BLOCK_REBALANCE`     | bool    | true, false     |
| Max records       | Number of messages to retrieve from Kafka topic per polling iteration                                                                                                                                     | `max_records`         | `APP_PUBSUB_KAFKA_SUBSCRIBER_MAX_RECORDS`         | number  | 1,2...          |

**IMPORTANT NOTES**:
1. `APP_PUBSUB_KAFKA_SUBSCRIBER_WORKERS` is a setting for subscriber to increase the throughput. If the number is more than 1, messages will be processed concurrently in a separate goroutine. If the order of messages is important, this setting should be set to 1.
2. `APP_PUBSUB_KAFKA_SUBSCRIBER_BLOCK_REBALANCE` is a setting for subscriber to block rebalance event when polling. If the setting is set to `true`, the subscriber will block the rebalance event when polling. This option helps to reduce the amount of duplicate messages. But it is important to be sure that the message processing is fast enough to avoid blocking the rebalance event for a long time. For more information refer to the [franz-go documentation](https://pkg.go.dev/github.com/twmb/franz-go/pkg/kgo#BlockRebalanceOnPoll).

To authenticate the requests to Kafka, Go SDK provides a configuration set for TLS and [SASL](https://en.wikipedia.org/wiki/Simple_Authentication_and_Security_Layer)

**TLS**:

| Setting               | Description                                                      | YAML variable          | Environment variable (ENV)                  | Type   | Possible Values                  |
|-----------------------|------------------------------------------------------------------|------------------------|---------------------------------------------|--------|----------------------------------|
| Enabled               | Whether TLS authentication is enabled or not                     | `enabled`              | `APP_PUBSUB_KAFKA_TLS_ENABLED`              | bool   | true, false                      |
| Root certificate      | Ca Root CA certificate                                           | `ca`                   | `APP_PUBSUB_KAFKA_TLS_CA`                   | string | Root CA                          |
| Cert PEM              | Client's public key string (PEM format) used for authentication  | `cert_pem`             | `APP_PUBSUB_KAFKA_TLS_CERT_PEM`             | string | long PEM string (deprecated)     |
| Cert PEM Key          | Client's private key string (PEM format) used for authentication | `cert_pem_key`         | `APP_PUBSUB_KAFKA_TLS_CERT_PEM_KEY`         | string | long PEM key string (deprecated) |
| Passphrase            | Passphrase is used in case the private key needs to be decrypted | `passphrase`           | `APP_PUBSUB_KAFKA_TLS_PASSPHRASE`           | string | pass phrase                      |
| Skip TLS verification | Turn on / off TLS verification                                   | `insecure_skip_verify` | `APP_PUBSUB_KAFKA_TLS_INSECURE_SKIP_VERIFY` | bool   | true, false                      |

**SASL**:

| Setting                | Description                                                                                                                                                                                             | YAML variable   | Environment variable (ENV)                        | Type   | Possible Values    |
|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------|---------------------------------------------------|--------|--------------------|
| Enabled                | Whether SASL authentication is enabled or not                                                                                                                                                           | `enabled`       | `APP_PUBSUB_KAFKA_SASL_ENABLED`                   | bool   | true, false        |
| Mechanism              | Mechanism is a string representation of the SASL mechanism. Only "plain" and "aws_msk_iam" are supported                                                                                                | `mechanism`     | `APP_PUBSUB_KAFKA_SASL_MECHANISM`                 | string | plain, aws_msk_iam |
| Username               | The username to authenticate Kafka requests                                                                                                                                                             | `username`      | `APP_PUBSUB_KAFKA_SASL_USERNAME`                  | string | username           |
| Password               | The password to authenticate Kafka requests                                                                                                                                                             | `paswword`      | `APP_PUBSUB_KAFKA_SASL_PASSWORD`                  | string | password           |
| AWS MSK IAM access key | AWS MSK IAM access key to authenticate AWS MSK requests                                                                                                                                                 | `access_key`    | `APP_PUBSUB_KAFKA_SASL_AWS_MSK_IAM_ACCESS_KEY`    | string | access key         |
| AWS MSK IAM secret key | AWS MSK IAM secret key to authenticate AWS MSK requests                                                                                                                                                 | `secret_key`    | `APP_PUBSUB_KAFKA_SASL_AWS_MSK_IAM_SECRET_KEY`    | string | secret key         |
| AWS STS Session Token  | SessionToken is used to authenticate AWS MSK requests via AWS STS service.<br/>For more information please check the [documentation](https://docs.aws.amazon.com/STS/latest/APIReference/welcome.html). | `session_token` | `APP_PUBSUB_KAFKA_SASL_AWS_MSK_IAM_SESSION_TOKEN` | string | token              |
| User agent             | The client's user agent string                                                                                                                                                                          | `user_agent`    | `APP_PUBSUB_KAFKA_SASL_AWS_MSK_IAM_USER_AGENT`    | string | user agent         |
| Assumable role         | This role will be used to establish connection to AWS MSK ignoring the static credentials                                                                                                               | `role`          | `APP_PUBSUB_KAFKA_SASL_AWS_MSK_IAM_ROLE`          | string | AWS ARN string     |
| Session name           | Will be passed to AWS STS when assuming the role                                                                                                                                                        | `session_name`  | `APP_PUBSUB_KAFKA_SASL_AWS_MSK_IAM_SESSION_NAME`  | string | application        |

## APM & Instrumentation

The `go-sdk` provides an easy way to add application performance monitoring
(APM) & instrumentation to a service. It provides DataDog APM using the
[`dd-trace-go`](https://github.com/DataDog/dd-trace-go) library. `dd-trace-go`
provides gRPC server & client interceptors, HTTP router instrumentation, database connection & ORM instrumentation
and AWS session instrumentation. All of the traces and data are opaquely sent
to DataDog.

### Request ID middleware

For easier identification of requests and their tracing within components of a
single service, and across-services, `go-sdk` has a `RequestID` middleware.

#### HTTP server Request ID middleware

As an HTTP middleware, it checks every incoming request for a `X-Request-Id` HTTP header
and sets the value of the header as a field in the request `Context`. In case
there's no `RequestID` present, it will generate a UUID and assign it in the
request `Context`.

Example usage of the middleware:

```go
func main() {
	requestIDMiddleware := sdkmiddleware.NewRequestIDMiddleware()

	httpServer := server.
		NewHTTPServer(host, httpPort, applicationEnv, applicationName).
		MountMiddleware(requestIDMiddleware.Handler).
		MountRoutes(routes)
}
```

#### gRPC server Request ID interceptors

For the gRPC server, sdk provides unary and stream interceptors. It checks every incoming request
for a `x-request-id` [grpc metadata](https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md) header and sets the value of the header as a field in the request `Context`.
In case there's no `RequestID` present, it will generate a UUID and assign it in the request `Context`.

Example usage of interceptors:

```go
func main() {
    grpcServer, err := server.NewGrpcServer(
        host,
        grpcPort,
        []grpc.ServerOption{
            grpc.ChainUnaryInterceptor(
                sdkinterceptors.TracingUnaryServerInterceptor(applicationName),
                sdkinterceptors.RequestIDUnaryServerInterceptor(),
                sdkinterceptors.LoggerUnaryServerInterceptor(logger),
            ),
            grpc.ChainStreamInterceptor(
                sdkinterceptors.TracingStreamServerInterceptor(applicationName),
                sdkinterceptors.RequestIDStreamServerInterceptor(),
                sdkinterceptors.LoggerStreamServerInterceptor(logger),
            ),
        }...)
}
```

### HTTP Router Instrumentation

`go-sdk` ships with HTTP router instrumentation, based on DataDog's
`dd-trace-go` library. It spawns a new instrumented router where each of the
requests will create traces that will be sent opaquely to the DataDog agent.

Example usage of the instrumentation:

```go
// Example taken from go-chassis
func NewHTTPServer(host, port, applicationEnv, applicationName string) *HTTPServer {
	router := sdk.Tracer.Router(applicationName)
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: router,
	}
	return &HTTPServer{
		srv:            srv,
		applicationEnv: applicationEnv,
	}
}
```

### gRPC server and client interceptors

`go-sdk` ships with gRPC instrumentation, based on DataDog's
`dd-trace-go` library. It creates new gRPC interceptors for the server & client. Thus,
requests will create traces that will be sent opaquely to the DataDog agent.

Example usage of the instrumentation:

```go
// gRPC client example

si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName("client-application"))
ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName("client-application"))

conn, err := grpc.Dial("<gRPC host>", grpc.WithStreamInterceptor(si), grpc.WithUnaryInterceptor(ui))
if err != nil {
    log.Fatalf("Failed to init gRPC connection: %s", err)
}
defer conn.Close()
```

```go
// gRPC server example

grpcServer, err := server.NewGrpcServer(
    host,
    grpcPort,
    []grpc.ServerOption{
        grpc.ChainUnaryInterceptor(
            sdkinterceptors.TracingUnaryServerInterceptor(applicationName),
            sdkinterceptors.LoggerUnaryServerInterceptor(logger),
            sdkinterceptors.MetricsUnaryServerInterceptor(metrics),
            sdkinterceptors.DatabaseUnaryServerInterceptor(database),
            kitgrpc.Interceptor,
        ),
        grpc.ChainStreamInterceptor(
            sdkinterceptors.TracingStreamServerInterceptor(applicationName),
            sdkinterceptors.LoggerStreamServerInterceptor(logger),
            sdkinterceptors.MetricsStreamServerInterceptor(metrics),
            sdkinterceptors.DatabaseStreamServerInterceptor(database),
        ),
    }...)
if err != nil {
    logger.Fatalf("Failed to create gRPC Server: %s", err)
}
```

### Database instrumentation & ORM logging

`go-sdk` ships with two database-related middlewares: `Database` &
`DatabaseLogging` for both HTTP and gRPC servers.

The `Database` middleware which instruments the
[Gorm-powered](https://gorm.io/gorm) database connection. It utilizes
Gorm-specific callbacks that report spans and traces to Datadog. The
instrumented Gorm database connection is injected in the request `Context` and
it is always scoped within the request.

The `DatabaseLogging` middleware checks for a logger injected in the request
`context`. If found, the logger is passed to the Gorm database connection,
which in turn uses the logger to produce database query logs. A nice
side-effect of this approach is that, if the logger is tagged with a
`request_id`, there's a logs correlation between the HTTP requests and the
database queries. Also, if the logger is tagged with `trace_id` we can easily
correlate logs with traces and see corresponding database queries. Keep in mind 
that ORM logging happens at the same level as the base logger. The additional
database fields ('duration', 'affected_rows' and 'sql') are available
only when using the `trace` levels.

#### HTTP server middleware example

```go
func main() {
	databaseMiddleware := sdkmiddleware.NewDatabaseMiddleware(sdk.Database)
	databaseLoggingMiddleware := sdkmiddleware.NewDatabaseLoggingMiddleware()

	httpServer := server.
		NewHTTPServer(host, httpPort, applicationEnv, applicationName).
		MountMiddleware(databaseMiddleware.Handler).
		MountMiddleware(databaseLoggingMiddleware.Handler).
		MountRoutes(routes)
}
```

#### gRPC server interceptors example

```go
func main() {
    grpcServer, err := server.NewGrpcServer(
        host,
        grpcPort,
        []grpc.ServerOption{
            grpc.ChainUnaryInterceptor(
                sdkinterceptors.TracingUnaryServerInterceptor(applicationName),
                sdkinterceptors.LoggerUnaryServerInterceptor(logger),
                sdkinterceptors.MetricsUnaryServerInterceptor(metrics),
                sdkinterceptors.DatabaseUnaryServerInterceptor(database),
                sdkinterceptors.DatabaseLoggingUnaryServerInterceptor(),
                kitgrpc.Interceptor,
            ),
            grpc.ChainStreamInterceptor(
                sdkinterceptors.TracingStreamServerInterceptor(applicationName),
                sdkinterceptors.LoggerStreamServerInterceptor(logger),
                sdkinterceptors.MetricsStreamServerInterceptor(metrics),
                sdkinterceptors.DatabaseStreamServerInterceptor(database),
                sdkinterceptors.DatabaseLoggingStreamServerInterceptor(),
            ),
        }...)
}
```

### AWS Session instrumentation

`go-sdk` instruments the AWS session by wrapping it with a DataDog trace and
tagging it with the service name. In addition, this registers AWS as a separate
service in DataDog.

Example usage of the instrumentation:

```go
func main() {
    s := session.NewSession(&aws.Config{
		Endpoint: aws.String(config.GetString("s3_endpoint")),
		Region:   aws.String(config.GetString("default_region")),
		Credentials: credentials.NewStaticCredentials(
			config.GetString("access_key_id"),
			config.GetString("secret_access_key"),
			"",
		),
	})

    session = instrumentation.InstrumentAWSSession(s, instrumentation.Settings{AppName: "MyServiceName"})

    // Use the session...
}
```

To correlate the traces that are registered with DataDog with the corresponding
requests in other DataDog services, the
[`aws-sdk-go`](https://aws.amazon.com/sdk-for-go) provides functions with the
`WithContext` suffix. These functions expect the request `Context` as the first
arugment of the function, which allows the tracing chain to be continued inside
the SDK call stack. To learn more about these functions you can start by
reading about them on [the AWS developer
blog](https://aws.amazon.com/blogs/developer/v2-aws-sdk-for-go-adds-context-to-api-operations).

### PubSub instrumentation and logging

#### Kafka

`go-sdk` provides a flexible way to instrument the Kafka client with logging, tracing and metrics capabilities.

The [franz-go](https://github.com/twmb/franz-go/) Kafka client library is wrapped by the `go-sdk` and traces calls to Kafka.

For logging, `go-sdk` provides a [franz-go compatible](https://pkg.go.dev/github.com/twmb/franz-go/pkg/kgo#WithLogger) logger which is basically a plugin
around the SDK's `logger.Logger` interface. This logger will print `franz-go` log stubs depending on a configured log level.

```go
import (
    "github.com/twmb/franz-go/pkg/kgo"

    sdkkafkalogger "github.com/scribd/go-sdk/pkg/logger/kafka"
)

func main() {
    kafkaClientOpts := []kgo.Opt{
        kgo.WithLogger(
            sdkkafkalogger.NewKafkaLogger(
                logger,
                sdkkafkalogger.WithLevel(sdklogger.Level(config.Logger.ConsoleLevel)),
            ),
        ),
    }
}
```

For metrics, `go-sdk` provides a `franz-go` plugin to publish metrics to DataDog via DogStatsd. It is implemented as a
[hook](https://pkg.go.dev/github.com/twmb/franz-go/pkg/kgo#WithHooks). Metrics are split into three categories: `broker`, `producer` and `consumer`.

**Broker metrics**

Broker metrics represent a set of metrics around the broker connection. If enabled, the following metrics will be published:

```
count metrics
#{service_name}.kafka_client.broker.read_errors_total{node="#{node}"}
#{service_name}.kafka_client.broker.read_bytes_total{node="#{node}"}
#{service_name}.kafka_client.broker.write_errors_total{node="#{node}"}
#{service_name}.kafka_client.broker.write_bytes_total{node="#{node}"}
#{service_name}.kafka_client.broker.disconnects_total{node="#{node}"}
#{service_name}.kafka_client.broker.connect_errors_total{node="#{node}"}
#{service_name}.kafka_client.broker.connects_total{node="#{node}"}

time metrics
#{service_name}.kafka_client.broker.write_latency{node="#{node}"}
#{service_name}.kafka_client.broker.write_wait_latency{node="#{node}"}
#{service_name}.kafka_client.broker.read_latency{node="#{node}"}
#{service_name}.kafka_client.broker.read_wait_latency{node="#{node}"}

histogram metrics
#{service_name}.kafka_client.broker.throttle_latency{node="#{node}"}
```

**Producer metrics**

Producer metrics represent a set of metrics around the producer. If enabled, the following metrics will be published:

```
count metrics
#{service_name}.kafka_client.producer.records_error
#{service_name}.kafka_client.producer.records_unbuffered
#{service_name}.kafka_client.producer.records_buffered
#{service_name}.kafka_client.producer.produce_bytes_uncompressed_total{node="#{node}",topic="#{topic}",partition="#{partition}"}
#{service_name}.kafka_client.producer.produce_bytes_compressed_total{node="#{node}",topic="#{topic}",partition="#{partition}"}
```

**Consumer metrics**

Consumer metrics represent a set of metrics around the consumer. If enabled, the following metrics will be published:

```
count metrics
#{service_name}.kafka_client.consumer.records_unbuffered
#{service_name}.kafka_client.consumer.records_buffered
#{service_name}.kafka_client.consumer.fetch_bytes_uncompressed_total{node="#{node}",topic="#{topic}",partition="#{partition}"}
#{service_name}.kafka_client.consumer.fetch_bytes_compressed_total{node="#{node}",topic="#{topic}",partition="#{partition}"}
```

By default, all metrics will be published without sampling (rate `1.0`). The client can specify the sampling rate per
metric type and/or per metric name:

```go
import (
    "github.com/twmb/franz-go/pkg/kgo"

    sdkkafkamertics "github.com/scribd/go-sdk/pkg/metrics/kafka"
)

func main() {
    kafkaClientOpts := []kgo.Opt{
        // no sampling rate specified
        kgo.WithHooks(sdkkafkamertics.NewBrokerMetrics(metrics)),
        // sampling rate of 0.5 will be applied for all producer metrics
        kgo.WithHooks(sdkkafkamertics.NewProducerMetrics(metrics, sdkkafkamertics.WithSampleRate(0.5))),
        // sampling rate of 0.3 will be applied to consumer metric name provided only,
        // other metrics will be published with default sample rate
        kgo.WithHooks(sdkkafkamertics.NewConsumerMetrics(
            metrics,
            sdkkafkamertics.WithSampleRatesPerMetric(map[string]float64{
                "kafka_client.consumer.records_buffered": 0.3,
            }),
        )),
    }
}
```

### Profiling

You can send `pprof` samples to DataDog by enabling the profiler.
Under the hood the DataDog profiler will continuously take heap, CPU and mutex profiles, [every 1 minute by default](https://godoc.org/gopkg.in/DataDog/dd-trace-go.v1/profiler#pkg-constants).
The [default CPU profile duration is 15 seconds](https://godoc.org/gopkg.in/DataDog/dd-trace-go.v1/profiler#pkg-constants). Keep in mind that the profiler introduces overhead when it is being executed.
The default DataDog configuration, which go-sdk uses by default, is following [good practices](https://groups.google.com/g/golang-nuts/c/e6lB8ENbIw8?pli=1).

```go
func main() {
    // Build profiler.
    sdkProfiler := instrumentation.NewProfiler(
	config.Instrumentation,
	profiler.WithService(appName),
	profiler.WithVersion(version),
    )
    if err := sdkProfiler.Start(); err != nil {
        log.Fatalf("Failed to start profiler: %s", err)
    }
    defer sdkProfiler.Stop()
}
```

### Custom Metrics

`go-sdk` provides a way to send custom metrics to Datadog.

The `metrics.Configuration` configuration can be used to define the set of
tags to attach to every metric emitted by the client.

This client is configured by default with:

- `service:$APP_NAME`
- `env:$ENVIRONMENT`

Datadog tags documentation is available [here][ddtags].

See the [metric submission documentation][submit-metric] on how to
submit custom metrics.

Metric names must only contain ASCII alphanumerics, underscores, and
periods. The client will not replace nor check for invalid characters.

Some options are suppported when submitting metrics, like applying a
[sample rate][rate] to your metrics or tagging your metrics with your
[custom tags][custom-tags]. Find all the available functions to report
metrics in the Datadog Go [client GoDoc documentation][client-go].

Example usage of the custom metrics:

```go
package main

import (
	"log"
	"time"

	"github.com/scribd/go-sdk/pkg/metrics"
)

func main() {
	applicationEnv := "development"
	applicationName := "go-sdk-example"

	metricsConfig := &metrics.Config{
		Environment: applicationEnv,
		App:         applicationName,
	}
	client, err := metrics.NewBuilder(metricsConfig).Build()
	if err != nil {
		log.Fatalf("Could not initialize Metrics client: %s", err)
	}

	_ = client.Incr("example.increment", []string{""}, 1)
	_ = client.Decr("example.decrement", []string{""}, 1)
	_ = client.Count("example.count", 2, []string{""}, 1)
}
```

[ddtags]: <https://docs.datadoghq.com/getting_started/tagging/>
[client-go]: <https://godoc.org/github.com/DataDog/datadog-go/statsd#Client>
[custom-tags]: <https://docs.datadoghq.com/developers/metrics/dogstatsd_metrics_submission/?tab=go#metric-tagging>
[rate]: <https://docs.datadoghq.com/developers/metrics/dogstatsd_metrics_submission/?tab=go#metric-submission-options>
[submit-metric]: <https://docs.datadoghq.com/developers/metrics/dogstatsd_metrics_submission/?tab=go>

## Using the `go-sdk` in isolation

The `go-sdk` is a standalone Go module. This means that it can be imported and
used in virtually any Go project. Still, there are four conventions that the
`go-sdk` enforces which **must be present** in the host application:
1. The presence of the `APP_ENV` environment variable, used for
   [environment-awareness](#environment-awareness),
2. Support of a single file format (YAML) for storing configurations,
3. Support of only a single path to store the configuration, `config/` in the
   application root, and
4. The presence of the `APP_ROOT` environment variable, set to the absolute
   path of the application that it's used in. This is used to locate the
   `config` directory on disk and load the enclosed YAML files.

A good way to approach initialization of the `go-sdk` in your application can be
seen in the `go-chassis` itself:

```go
// internal/pkg/sdk/sdk.go
package sdk

import (
	"log"

	sdkconfig "github.com/scribd/go-sdk/pkg/configuration"
	sdklogger "github.com/scribd/go-sdk/pkg/logger"
)

var (
	// Config is SDK-powered application configuration.
	Config *sdkconfig.Config
	// Logger is SDK-powered application logger.
	Logger sdklogger.Logger
	err    error
)

func init() {
	if Config, err = sdkconfig.NewConfig(); err != nil {
		log.Fatalf("Failed to load SDK config: %s", err.Error())
	}

	if Logger, err = sdklogger.NewBuilder(loggerConfig).Build(); err != nil {
		log.Fatalf("Failed to load SDK logger: %s", err.Error())
	}
}
```

**Please note** that while using the `go-sdk` in isolation is possible, it is
**highly recommended** to use it in combination with the `go-chassis` for the
best development, debugging and maintenance experience.

## Developing the SDK

* The SDK provides a [Docker
  Compose](https://docs.docker.com/compose/) development environment to
  run, develop and test a service (version `1.24.1`).
* See [`docker-compose.yml`](./docker-compose.yml) for the service definitions.

### Building the docker environment

```sh
$ docker-compose build [--no-cache|--pull|--parallel]
```

Refer to the
[Compose CLI reference](https://docs.docker.com/compose/reference/build/)
for the full list of option and the synopsis of the command.

### Running tests within the docker environment

Compose provides a way to create and destroy isolated testing environments:

```sh
$ docker-compose run --rm sdk mage test:run
```

### Entering the docker environment

You can enter the docker environment to build, run and debug your service:

```
$ docker-compose run --rm sdk /bin/bash
root@1f31fa8e5c49:/sdk# go version
go version go1.20.6 linux/amd64
```

Refer to the
[Compose CLI reference](https://docs.docker.com/compose/reference/run/)
for the full list of option and the synopsis of the command.

### Using a development build of `go-sdk`

When developing a project that uses `go-sdk` as a dependency, you might need
to introduce changes to the `go-sdk` codebase.

In this case, perform the following steps:

1. Create a branch in the `go-sdk` repository with the changes you want,
   and push the branch to the repository remote:

   ```sh
   git push -u origin <username/branch-name>
   ```

2. Add or change the following line in the `go.mod` file of the project
   that uses `go-sdk` as a dependency:

   ```go
   replace github.com/scribd/go-sdk => github.com/scribd/go-sdk.git <username/branch-name>
   ```

3. From the project root, fetch the new branch by running:

   ```sh
   go mod tidy
   ```

4. Note that running `go mod tidy` will tie `go-sdk` to the specific git commit.
   So after running it, the `replace` statement will look like this:

   ```go
   replace github.com/scribd/go-sdk => github.com/scribd/go-sdk.git <pseudo-version>
   ```

   Therefore, you will need to repeat steps 1, 2 and 3 each time you add new
   changes to your branch in `go-sdk`.

5. After you are done with the required changes, create a merge request to the
   `go-sdk` repository. After the merge request is merged and a release is done,
   you need to, once again, alter the `replace` statement in your `go.mod` file:

   ```go
   replace github.com/scribd/go-sdk => github.com/scribd/go-sdk.git <tag-name>
   ```

### Commit messages

Please follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) and [Angular commit message guidelines](https://github.com/angular/angular/blob/22b96b9/CONTRIBUTING.md#-commit-message-guidelines) for writing commit messages.

## Release

This project is using [semantic-release](https://semantic-release.gitbook.io/semantic-release/)
and [conventional-commits](https://www.conventionalcommits.org/en/v1.0.0/),
with the [`angular` preset](https://github.com/conventional-changelog/conventional-changelog/tree/master/packages/conventional-changelog-angular).

Releases are done from the `origin/main` branch using a manual step at the end of the CI/CD pipeline.

In order to create a new release:

1. Merge / push changes to `origin/main`
2. Open the `origin/main` [GitHub Action Release pipeline](https://github.com/scribd/go-sdk/actions/workflows/release.yml)
3. Press ▶️ on the "Run workflow" step

A version bump will happen automatically and the type of version bump
(patch, minor, major) depends on the commits introduced since the last release.

The `semantic-release` configuration is in [`.releaserc.yml`](https://github.com/scribd/go-sdk/blob/main/.releaserc.yml).

## Maintainers

Made with ❤️ by the Core Services team.
