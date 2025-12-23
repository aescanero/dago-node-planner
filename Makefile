# Makefile for dago-node-planner

# Variables
APP_NAME := node-planner
BINARY := $(APP_NAME)
VERSION ?= dev
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt

# Docker parameters
DOCKER_IMAGE := aescanero/dago-node-planner
DOCKER_TAG := $(VERSION)

# Paths
CMD_PATH := ./cmd/$(APP_NAME)
BUILD_DIR := ./bin

.PHONY: all build clean test coverage lint fmt deps tidy run dev docker-build docker-push help

## help: Display this help message
help:
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## all: Build the binary
all: build

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOGET) -v ./...
	$(GOMOD) download

## tidy: Tidy go.mod
tidy:
	@echo "Tidying go.mod..."
	$(GOMOD) tidy

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout 5m; \
	else \
		echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

## coverage: Run tests with coverage report
coverage: test
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## build: Build the binary
build: deps
	@echo "Building $(BINARY)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) $(CMD_PATH)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY)"

## build-linux: Build for Linux
build-linux: deps
	@echo "Building $(BINARY) for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY)-linux-amd64 $(CMD_PATH)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY)-linux-amd64"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

## run: Build and run locally
run: build
	@echo "Running $(BINARY)..."
	$(BUILD_DIR)/$(BINARY)

## run-local: Run locally without building (requires 'go run')
run-local:
	@echo "Running $(APP_NAME) locally..."
	$(GOCMD) run $(CMD_PATH)/main.go

## dev: Run with live reload (requires 'air')
dev:
	@echo "Starting development server with live reload..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f deployments/docker/Dockerfile .
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

## docker-push: Push Docker image
docker-push: docker-build
	@echo "Pushing Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest
	@echo "Docker image pushed"

## release: Create a release (tag version)
release:
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "VERSION is required. Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release $(VERSION)..."
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	@echo "Release $(VERSION) created"

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/cosmtrek/air@latest
	@echo "Tools installed"
