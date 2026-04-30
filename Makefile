default: help

.PHONY: help
help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST)  | fgrep -v fgrep | sed -e 's/:.*##/:##/' | awk -F':##' '{printf "%-12s %s\n",$$1, $$2}'

.PHONY: test
test: ## Run tests.
	go test -v ./...

.PHONY: test-race
test-race: ## Run tests with the race detector.
	go test -race ./...

.PHONY: test-cover
test-cover: ## Run tests with package coverage.
	go test -cover ./...

.PHONY: test-fuzz-document
test-fuzz-document: ## Run document fuzz tests briefly.
	go test -fuzz=FuzzPositionOffsetRoundTrip -fuzztime=30s ./document

.PHONY: lint_install
lint_install: ## Install golangci-lint
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

.PHONY: lint
lint: lint_install ## Run linting operations
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format the code
	@gofmt -w .
