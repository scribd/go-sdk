# =============================================================================
# build stage
# =============================================================================

FROM golang:1.13.3 as builder

ENV GO111MODULE=on

RUN go get -u github.com/rakyll/gotest

WORKDIR /sdk

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

# =============================================================================
# development stage
# =============================================================================

FROM builder AS development

RUN go get -u -v github.com/go-delve/delve/cmd/dlv \
	github.com/golangci/golangci-lint/cmd/golangci-lint@v1.20.1
