SHELL=/bin/bash
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOLINT=golangci-lint

GO111MODULE=on

.DEFAULT_GOAL := all

.PHONY: all
all: format tidy lint build test ## Run all checks and tests

.PHONY: help
help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: format
format: ## Format the code
	gofmt -s -w .

.PHONY: tidy
tidy: ## Tidy and verify dependencies
	$(GOCMD) mod tidy
	$(GOCMD) mod verify

.PHONY: download
download: ## Download module dependencies
	$(GOCMD) mod download

.PHONY: install-lint
install-lint: ## Install golangci-lint
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.7.2

.PHONY: lint
lint: install-lint ## Run the linter
	$(GOLINT) run --timeout=5m ./...

.PHONY: build
build: ## Build the library
	$(GOBUILD) -v  ./...

.PHONY: test
test: ## Run all the tests with race detection
	$(GOTEST) -v -race -covermode=atomic -coverprofile=coverage.out -timeout=30s ./...

.PHONY: test-short
test-short: ## Run short tests
	$(GOTEST) -v -short ./...

.PHONY: bench
bench: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./...

.PHONY: fuzz
fuzz: ## Run fuzz tests for 30s each
	$(GOTEST) -fuzz=. -fuzztime=30s ./...

.PHONY: coverage
coverage: test ## Generate coverage report
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	$(GOCMD) tool cover -func=coverage.out

.PHONY: ci
ci: lint test ## Run all the tests and code checks (CI pipeline)

.PHONY: clean
clean: ## Clean build artifacts and test cache
	$(GOCLEAN) -testcache
	rm -f coverage.out coverage.html
