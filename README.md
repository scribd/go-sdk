# Go SDK

SDK, the Go version.

## Table Of Contents

- [Prerequisites](#prerequisites)
- [SDK functionality](#sdk-functionality)
   - [Application Configuration](#application-configuration)
      - [Predefined application-agnostic configurations](#predefined-application-agnostic-configurations)
      - [Custom application-specific configurations](#custom-application-specific-configurations)
      - [Environment-awareness](#environment-awareness)
      - [Using application configuration in tests](#using-application-configuration-in-tests)
   - [Logger](#logger)
      - [Initialization and default configuration](#initialization-and-default-configuration)
      - [Environment-awareness](#environment-awareness-1)
      - [Log levels](#log-levels)
      - [Structured logging](#structured-logging)
   - [Logging &amp; tracing middleware](#logging--tracing-middleware)
      - [Formatting and handlers](#formatting-and-handlers)
      - [Sentry error reporting](#sentry-error-reporting)
   - [Database Connection](#database-connection)
   - [ORM Integration](#orm-integration)
      - [Usage of ORM](#usage-of-orm)
- [APM &amp; Instrumentation](#apm--instrumentation)
   - [Request ID middleware](#request-id-middleware)
   - [HTTP Router Instrumentation](#http-router-instrumentation)
   - [Database instrumentation &amp; ORM logging](#database-instrumentation--orm-logging)
   - [AWS Session instrumentation](#aws-session-instrumentation)
- [Using the go-sdk in isolation](#using-the-go-sdk-in-isolation)
- [Developing the SDK](#developing-the-sdk)
   - [Building the docker environment](#building-the-docker-environment)
   - [Running tests within the docker environment](#running-tests-within-the-docker-environment)
   - [Entering the docker environment](#entering-the-docker-environment)
- [Maintainers](#maintainers)

## Prerequisites

* [Go](https://golang.org) (version `1.13.3`).
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
  app_name: "my-awesome-app"

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
applicationName := sdk.Config.App.GetString("app_name")
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

    sdkconfig "git.lo/microservices/sdk/go-sdk/pkg/configuration"
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

	sdklogger "git.lo/microservices/sdk/go-sdk/pkg/logger"
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

`go-sdk` ships with a `Logger` middleware. When used, it assigns a random
UUID as a `RequestID`, if none exists, and it opens a new span using Datadog's
`ddtrace` library. The span has a corresponding `SpanID` that is used to easily
create sub-spans for application performance monitoring (APM). The `RequestID`
is used for request-scoped logs correlation.

Example usage of the middleware:

```go
func main() {
	loggingMiddleware := sdkmiddleware.NewLoggingMiddleware(sdk.Logger)

	httpServer := server.
		NewHTTPServer(host, httpPort, applicationEnv, applicationName).
		MountMiddleware(loggingMiddleware.Handler).
		MountRoutes(routes)
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
  timeout: 200 # milliseconds

development:
  <<: *common

staging:
  <<: *common

production:
  <<: *common
  dsn: "https://<key>@sentry.io/<project>"
```

The tracking can be enabled from the `Builder` with the `SetTracking`
function:

```go
package main

import (
	"log"

	sdklogger   "git.lo/microservices/sdk/go-sdk/pkg/logger"
	sdktracking "git.lo/microservices/sdk/go-sdk/pkg/tracking"
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
```

A logger build with a valid tracking configuration will automatically
report to Sentry any errors emitted from the following log levels:

- `Error`;
- `Fatal`;
- `Panic`;

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

### ORM Integration

`go-sdk` comes with an integration with the popular
[gorm](https://github.com/jinzhu/gorm) as an object-relational mapper (ORM).
Using the configuration details, namely the [data source
name](https://en.wikipedia.org/wiki/Data_source_name) (DSN) as their product,
gorm is able to open a connection and give the `go-sdk` users a preconfigured
ready-to-use database connection with an ORM attached. This can be done as
follows:

```go
package main

import (
	sdkdb "git.lo/microservices/chassis/go-sdk/pkg/database"
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
[Gorm-powered](https://github.com/jinzhu/gorm) database connection. It can be
used right away to query the database:

```go
package main

import (
	"fmt"

	sdkdb "git.lo/microservices/chassis/go-sdk/pkg/database"
	"github.com/jinzhu/gorm"
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

## APM & Instrumentation

The `go-sdk` provides an easy way to add application performance monitoring
(APM) & instrumentation to a service. It provides DataDog APM using the
[`dd-trace-go`](https://github.com/DataDog/dd-trace-go) library. `dd-trace-go`
provides HTTP router instrumentation, database connection & ORM instrumentation
and AWS session instrumentation. All of the traces and data are opaquely sent
to DataDog.

### Request ID middleware

For easier identification of requests and their tracing within components of a
single service, and across-services, `go-sdk` ships has a `RequestID`
middleware. It checks every incoming request for a `X-Request-Id` HTTP header
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

### Database instrumentation & ORM logging

`go-sdk` ships with two database-related middlewares: `Database` &
`DatabaseLogging`.

The `Database` middleware which instruments the
[Gorm-powered](https://github.com/jinzhu/gorm) database connection. It utilizes
Gorm-specific callbacks that report spans and traces to Datadog. The
instrumented Gorm database connection is injected in the request `Context` and
it is always scoped within the request.

The `DatabaseLogging` middleware checks for a logger injected in the request
`context`. If found, the logger is passed to the Gorm database connection,
which in turn uses the logger to produce database query logs. A nice
side-effect of this approach is that, if the logger is tagged with a
`request_id`, there's a logs correlation between the HTTP requests and the
database queries.

Example usage of the middlewares:

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

    session = instrumentation.InstrumentSession(s)

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

	sdkconfig "git.lo/microservices/sdk/go-sdk/pkg/configuration"
	sdklogger "git.lo/microservices/sdk/go-sdk/pkg/logger"
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
go version go1.13.3 linux/amd64
```

Refer to the
[Compose CLI reference](https://docs.docker.com/compose/reference/run/)
for the full list of option and the synopsis of the command.

## Maintainers

Made with ❤️ by the Core Services team.
