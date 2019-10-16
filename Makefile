all: test

test: ## Run the tests.
	go test ./... -v

fmt: ## Run gofmt.
	go fmt ./...

help:
	@awk -F '##' '!/help/ && / ## / {sub(/:.*/, "", $$1); printf "\033[33m%-15s\033[0m %s\n", $$1, $$2}' Makefile

.PHONY: all test fmt
