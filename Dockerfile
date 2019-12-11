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

# Build Mage: a make-like build tool written in Go
ENV GOBIN=/go/bin
RUN go install cmd/mage/mage.go

# =============================================================================
# development stage
# =============================================================================

FROM builder AS development

RUN go get -u -v github.com/go-delve/delve/cmd/dlv \
	github.com/golangci/golangci-lint/cmd/golangci-lint@v1.20.1

ENV GOBIN=/go/bin
COPY --from=builder $GOBIN/mage /usr/bin/mage
