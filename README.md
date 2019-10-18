# Go SDK

SDK, the Go version.

## Table Of Contents

- [Table Of Contents](#table-of-contents)
- [Prerequisites](#prerequisites)
- [SDK functionality](#sdk-functionality)
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

* <insert package here>

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
