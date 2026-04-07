# Variables
APP_NAME := reference-service
BINARY_NAME := reference-service-go
BUILD_DIR := ./bin
CHART_PATH := ./chart
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build run generate test coverage fmt lint clean docker-build helm-lint help

.DEFAULT_GOAL := help

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application binary
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -trimpath -ldflags "-X reference-service-go/internal/build.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/reference-service-go

run: ## Run the application locally
	go run -ldflags "-X reference-service-go/internal/build.Version=$(VERSION)" ./cmd/reference-service-go -config config/config.yaml

generate: ## Run go generate to generate code from OpenAPI specs
	go generate ./...

test: ## Run all tests with race detection (requires Docker)
	TESTCONTAINERS_RYUK_DISABLED=true go test -v -race -tags integration ./...

coverage: ## Run tests with coverage report (requires Docker)
	@mkdir -p $(BUILD_DIR)
	TESTCONTAINERS_RYUK_DISABLED=true go test -v -race -tags integration -coverprofile=$(BUILD_DIR)/coverage.out -covermode=atomic ./...
	go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html

fmt: ## Format Go code
	golangci-lint fmt

lint: ## Run golangci-lint
	golangci-lint run --timeout=5m

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)
	go clean

docker-build: ## Build the Docker image
	docker build --build-arg VERSION=$(VERSION) -t $(APP_NAME):latest .

helm-lint: ## Lint the Helm chart
	helm lint $(CHART_PATH)
