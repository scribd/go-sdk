# Go SDK

SDK, the Go version.

## Table Of Contents

- [Table Of Contents](#table-of-contents)
- [Prerequisites](#prerequisites)
- [SDK functionality](#sdk-functionality)
    - [Application Configuration](#application-configuration)
        - [Predefined application-agnostic configurations](#predefined-application-agnostic-configurations)
        - [Custom application-specific configurations](#custom-application-specific-configurations)
    - [Environment-awareness](#environment-awareness)
- [Developing the SDK](#developing-the-sdk)
    - [Building the docker environment](#building-the-docker-environment)
    - [Running the docker environment](#running-the-docker-environment)
    - [Running tests within the docker environment](#running-tests-within-the-docker-environment)
    - [Entering the docker environment](#entering-the-docker-environment)
        - [Makefile](#makefile)
        - [Debugging](#debugging)
    - [Multi-container environment support](#multi-container-environment-support)
- [GitLab Pipeline](#gitlab-pipeline)
- [Maintainers](#maintainers)

## Prerequisites

* [Go](https://golang.org) (version `1.13.3`).
* [Docker](https://www.docker.com/) (version `19.03.2`).

## SDK functionality

The SDK provides the following out-of-the-box functionality:

* Application Configuration & Environment-awareness
* <insert package here>

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
  `config/database.yml` configuration file. (TBD)
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
$ docker-compose run --rm sdk make test
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

Made with ❤️  by the Core Services team.
