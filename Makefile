GOCACHE ?= /tmp/go-build
GOLANGCI_LINT ?= golangci-lint

.PHONY: fmt fmt-check lint vet build ci

fmt: ## Format Go files
	gofmt -w .

fmt-check: ## Check Go formatting
	test -z "$$(gofmt -l .)"

lint: ## Run golangci-lint
	GOCACHE=$(GOCACHE) $(GOLANGCI_LINT) run ./...

vet: ## Run go vet
	GOCACHE=$(GOCACHE) go vet ./...

build: ## Build all packages
	GOCACHE=$(GOCACHE) go build ./...

ci: fmt-check lint vet build ## Run all CI checks locally
