all: test

test: ## Run the tests.
	APP_ENV=test gotest ./... -v

fmt: ## Run gofmt.
	go fmt ./...

fmt-check: ## Check the code formatting.
	@if [ -n "$$(gofmt -l .)" ]; then echo "Go code is not formatted:\n\n$$(gofmt -d .)"; exit 1; fi

lint: ## Run the linter.
	golangci-lint run

help:
	@awk -F '##' '!/help/ && / ## / {sub(/:.*/, "", $$1); printf "\033[33m%-15s\033[0m %s\n", $$1, $$2}' Makefile

.PHONY: all test fmt fmt-check lint help
