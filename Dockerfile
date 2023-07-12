# =============================================================================
# build stage
# =============================================================================

FROM golang:1.20.6-alpine AS builder

WORKDIR /sdk

RUN apk add --no-cache \
	build-base \
	# curl version should be higher than 7.74 to mitigate SNYK-DEBIAN11-CURL-3320493
	curl \
	git \
	openssh-client \
	tzdata

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

# Build Mage: a make-like build tool written in Go
ENV GOBIN=/go/bin
RUN go install cmd/mage/mage.go

RUN CGO_ENABLED=0 GOOS=linux go build ./...

# =============================================================================
# linter stage
# =============================================================================

FROM builder AS linter

# binary will be $(go env GOPATH)/bin/golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
	| sh -s -- -b $(go env GOPATH)/bin v1.51.1

# install goimports
RUN go install golang.org/x/tools/cmd/goimports@v0.1.12

# =============================================================================
# development stage
# =============================================================================

FROM linter AS development

RUN go get -v \
	github.com/go-delve/delve/cmd/dlv@v1.4.0
