# Canonical hawk-eco Makefile for Go LIBRARY repos.
# Source of truth: https://github.com/GrayCodeAI/hawk/blob/main/.shared-templates/Makefile.library.tmpl
# hawk-mcpkit is a foundation library consumed by engines (sight, inspect);
# no standalone binary, no release beyond tagging.

# ---------------------------------------------------------------------------
# Project metadata
# ---------------------------------------------------------------------------
NAME      := hawk-mcpkit

# ---------------------------------------------------------------------------
# Versioning — sourced from VERSION file; falls back to git describe.
# ---------------------------------------------------------------------------
VERSION ?= $(shell v=$$(cat VERSION 2>/dev/null | head -n1 | tr -d '[:space:]'); if [ -n "$$v" ]; then echo "$$v"; else git describe --tags --always --dirty 2>/dev/null || echo "dev"; fi)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# ---------------------------------------------------------------------------
# Tooling — pinned, install if missing.
# ---------------------------------------------------------------------------
GOBIN_DIR    := $(shell go env GOPATH)/bin
GOLANGCI     := $(GOBIN_DIR)/golangci-lint
GOFUMPT      := $(GOBIN_DIR)/gofumpt
GOIMPORTS    := $(GOBIN_DIR)/goimports
GOVULNCHECK  := $(GOBIN_DIR)/govulncheck

# ---------------------------------------------------------------------------
# Phony declarations (alphabetical).
# ---------------------------------------------------------------------------
.PHONY: all boundaries build ci clean cover fmt help hooks lint lint-fix \
        security test test-race tidy version vet

boundaries: ## Enforce foundation-repo import boundaries (zero GrayCodeAI/* deps).
	bash ./scripts/check-ecosystem-boundaries.sh

# ---------------------------------------------------------------------------
# Default target.
# ---------------------------------------------------------------------------
all: lint test build ## Default — lint, test, build.

# ---------------------------------------------------------------------------
# Build.
# ---------------------------------------------------------------------------
build: ## Build the library package.
	go build ./...

# ---------------------------------------------------------------------------
# Tests.
# ---------------------------------------------------------------------------
test: ## Run unit tests.
	go test ./... -count=1 -timeout=60s

test-race: ## Run unit tests with the race detector.
	go test ./... -race -count=1 -timeout=120s

cover: ## Generate a coverage report (coverage.out + coverage.html).
	go test ./... -race -coverprofile=coverage.out -covermode=atomic -timeout=120s
	@go tool cover -func=coverage.out | grep "^total:"
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# ---------------------------------------------------------------------------
# Quality gates.
# ---------------------------------------------------------------------------
fmt: ## Format source files (gofumpt + goimports).
	@command -v $(GOFUMPT)   >/dev/null 2>&1 || (echo "install: go install mvdan.cc/gofumpt@latest"   && exit 1)
	@command -v $(GOIMPORTS) >/dev/null 2>&1 || (echo "install: go install golang.org/x/tools/cmd/goimports@latest" && exit 1)
	$(GOFUMPT) -w .
	$(GOIMPORTS) -w .

vet: ## Run go vet.
	go vet ./...

lint: ## Run golangci-lint.
	@command -v $(GOLANGCI) >/dev/null 2>&1 || (echo "install: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest" && exit 1)
	$(GOLANGCI) run ./... --timeout=5m

lint-fix: ## Run golangci-lint with --fix.
	@command -v $(GOLANGCI) >/dev/null 2>&1 || (echo "install: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest" && exit 1)
	$(GOLANGCI) run ./... --fix --timeout=5m

security: ## Run govulncheck.
	@command -v $(GOVULNCHECK) >/dev/null 2>&1 || (echo "install: go install golang.org/x/vuln/cmd/govulncheck@latest" && exit 1)
	$(GOVULNCHECK) ./...

tidy: ## Tidy go.mod / go.sum.
	go mod tidy
	go mod verify

# ---------------------------------------------------------------------------
# Composite gate used by CI and pre-push.
# ---------------------------------------------------------------------------
ci: tidy fmt vet lint boundaries test-race security ## Run everything CI runs.
	@echo "All CI checks passed."

# ---------------------------------------------------------------------------
# Misc.
# ---------------------------------------------------------------------------
version: ## Print the version that will be embedded.
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"

clean: ## Remove build artefacts.
	rm -rf coverage.out coverage.html
	go clean -testcache

help: ## Show this help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

hooks: ## Install git hooks via lefthook (boundary guard, tests, co-author strip).
	@command -v lefthook >/dev/null 2>&1 || (echo "install: go install github.com/evilmartians/lefthook@latest" && exit 1)
	git config --unset core.hooksPath 2>/dev/null || true
	lefthook install
