all: test

test:
	go test ./... -v

fmt:
	go fmt ./...

.PHONY: all test fmt
