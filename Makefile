.PHONY: build test run clean validate lint help

# Build configuration
APP_NAME := torus-neighbors
BUILD_DIR := ./bin
MAIN_PATH := ./cmd/main.go

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

run: build ## Build and run the application
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

validate: build ## Run local validation only
	@echo "Running local validation..."
	@$(BUILD_DIR)/$(APP_NAME) -validate

solve: build ## Solve challenge (interactive mode)
	@echo "Enter your user identifier:"
	@read user && DEBUG_HTTP=1 $(BUILD_DIR)/$(APP_NAME) -user "$$user"

lint: ## Run code linters
	@echo "Running linters..."
	@go fmt ./...
	@go vet ./...

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Development targets
dev-setup: deps ## Set up development environment
	@echo "Development environment setup complete"

check: lint test ## Run all checks (lint + test)
	@echo "All checks passed!"

install: build ## Install the binary to $GOPATH/bin
	@echo "Installing $(APP_NAME)..."
	@go install $(MAIN_PATH)
	@echo "Installed to $(shell go env GOPATH)/bin/$(APP_NAME)"