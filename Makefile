# Makefile for reference-service-go

# Variables
APP_NAME := reference-service
BINARY_NAME := service
CMD_PATH := cmd/main.go
BUILD_DIR := ./build
CHART_PATH := ./chart
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOGENERATE := $(GOCMD) generate

# Build flags
CGO_ENABLED := 0

# Phony targets
.PHONY: all build run generate test coverage fmt lint clean docker-build docker-run helm-lint helm-template mod-tidy help

# Default target
all: help

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  build         - Build the application binary"
	@echo "  run           - Run the application locally"
	@echo "  generate      - Run go generate to generate code from OpenAPI specs"
	@echo "  test          - Run tests with race detection"
	@echo "  coverage      - Run tests with coverage report"
	@echo "  fmt           - Format Go code using golangci-lint"
	@echo "  lint          - Run golangci-lint"
	@echo "  clean         - Remove build artifacts"
	@echo "  docker-build  - Build the Docker image"
	@echo "  docker-run    - Run the Docker container"
	@echo "  helm-lint     - Lint the Helm chart"
	@echo "  helm-template - Render Helm chart templates"
	@echo "  mod-tidy      - Tidy Go module dependencies"

## build: Build the application binary
build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

## run: Run the application locally
run:
	VERSION=$(VERSION) $(GORUN) $(CMD_PATH) -config config/config.yaml

## generate: Run go generate to generate code from OpenAPI specs
generate:
	$(GOGENERATE) ./...

## test: Run tests with race detection
test:
	$(GOTEST) -v -race ./...

## coverage: Run tests with coverage report
coverage:
	@mkdir -p $(BUILD_DIR)
	$(GOTEST) -v -race -coverprofile=$(BUILD_DIR)/coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html

## fmt: Format Go code
fmt:
	golangci-lint fmt

## lint: Run golangci-lint
lint:
	golangci-lint run --timeout=5m

## clean: Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)
	$(GOCLEAN)

## docker-build: Build the Docker image
docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(APP_NAME):latest .

## docker-run: Run the Docker container
docker-run:
	docker run --rm -p 8080:8080 $(APP_NAME):latest

## helm-lint: Lint the Helm chart
helm-lint:
	helm lint $(CHART_PATH)

## helm-template: Render Helm chart templates
helm-template:
	helm template $(APP_NAME) $(CHART_PATH)

## mod-tidy: Tidy Go module dependencies
mod-tidy:
	$(GOMOD) tidy
