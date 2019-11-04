# =============================================================================
# build stage
# =============================================================================

FROM golang:1.13.3

ENV GO111MODULE=on

RUN go get -u github.com/rakyll/gotest

WORKDIR /sdk

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
