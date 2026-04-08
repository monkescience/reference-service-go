# Variables
APP_NAME := reference-service
BINARY_NAME := reference-service-go
BUILD_DIR := ./bin
CHART_PATH := ./chart
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COVERAGE_PROFILE := $(BUILD_DIR)/coverage.out
COVERAGE_HTML := $(BUILD_DIR)/coverage.html
SPECTRAL_RULESET := openapi/spectral.ruleset.yaml
SPECTRAL_SPECS := openapi/reference-api.yaml

.PHONY: build run generate test coverage fmt lint spectral clean docker-build helm-lint help

.DEFAULT_GOAL := help

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application binary
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -trimpath -ldflags "-X reference-service-go/internal/build.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/reference-service-go

run: ## Run the application locally
	go run -ldflags "-X reference-service-go/internal/build.Version=$(VERSION)" ./cmd/reference-service-go -config config/config.yaml

generate: ## Run code generation for OpenAPI specs and SQLC
	sqlc generate
	go generate -tags tools ./...

test: ## Run all tests with race detection (requires Docker)
	TESTCONTAINERS_RYUK_DISABLED=true go test -v -race -tags integration ./...

coverage: ## Run black-box coverage report (requires Docker)
	@mkdir -p $(BUILD_DIR)
	@rm -f $(COVERAGE_PROFILE) $(COVERAGE_HTML)
	TESTCONTAINERS_RYUK_DISABLED=true go test -v -tags integration ./tests
	go tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)

fmt: ## Format Go code
	golangci-lint fmt

lint: ## Run golangci-lint
	golangci-lint run --timeout=5m

spectral: ## Lint owned OpenAPI specs with Spectral
	spectral lint --ruleset $(SPECTRAL_RULESET) $(SPECTRAL_SPECS)

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)
	go clean

docker-build: ## Build the Docker image
	docker build --build-arg VERSION=$(VERSION) -t $(APP_NAME):latest .

helm-lint: ## Lint the Helm chart
	helm lint $(CHART_PATH)
