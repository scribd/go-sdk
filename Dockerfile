# =============================================================================
# build stage
# =============================================================================

FROM golang:1.15.6 AS builder

ENV GO111MODULE=on

WORKDIR /sdk

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

# Build Mage: a make-like build tool written in Go
ENV GOBIN=/go/bin
RUN go install cmd/mage/mage.go

# =============================================================================
# linter stage
# =============================================================================

FROM builder AS linter

# binary will be $(go env GOPATH)/bin/golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
	| sh -s -- -b $(go env GOPATH)/bin v1.39.0

# =============================================================================
# development stage
# =============================================================================

FROM linter AS development

RUN go get -v \
	github.com/go-delve/delve/cmd/dlv@v1.4.0
