GOCACHE ?= /tmp/go-build
GOLANGCI_LINT ?= golangci-lint
COURSE_DIR ?= ../go-courses

.PHONY: fmt fmt-check lint grade vet build run ci

fmt: ## Format Go files
	gofmt -w .

fmt-check: ## Check Go formatting
	test -z "$$(gofmt -l .)"

lint: ## Run golangci-lint
	GOCACHE=$(GOCACHE) $(GOLANGCI_LINT) run ./...

grade: ## Run course acceptance tests for the homework specified by TASK
	@test -n "$(TASK)" || { echo "Usage: make grade TASK=<number>"; exit 1; }
	GOCACHE=$(GOCACHE) python3 "$(COURSE_DIR)/autotests/scripts/ci.py" --repo "$(CURDIR)" --task "$(TASK)"

vet: ## Run go vet
	GOCACHE=$(GOCACHE) go vet ./...

build: ## Build all packages
	GOCACHE=$(GOCACHE) go build ./...

run: ## Run the application
	GOCACHE=$(GOCACHE) go run ./cmd

ci: fmt-check lint vet build ## Run all CI checks locally
