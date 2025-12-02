# ========================================
# 💼 Scout Makefile
# ========================================

APP_NAME := scout
BINARY := bin/$(APP_NAME)
PKG := ./...
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*")

# Default target
.DEFAULT_GOAL := help

# Colors for logs
YELLOW := \033[1;33m
GREEN := \033[0;32m
RED := \033[0;31m
RESET := \033[0m

help:
	@echo "$(YELLOW)Available commands:$(RESET)"
	@echo "  make lint            🧹 Run Revive linter"
	@echo "  make test            🧪 Run Go tests"
	@echo "  make tidy            🧩 Run go mod tidy"
	@echo "  make build           🏗️  Build binary"
	@echo "  make run             🚀 Run binary in release mode"
	@echo "  make clean           🧼 Remove build artifacts"


# Run revive linter (requires revive.toml config)
lint:
	@echo "$(GREEN)Linting code with Revive...$(RESET)"
	@revive -config revive.toml ./... || echo "$(RED)Linting failed!$(RESET)"

# Run tests
test:
	@echo "$(GREEN)Running tests...$(RESET)"
	@go test -v $(PKG)

# Tidy up dependencies
tidy:
	@echo "$(GREEN)Tidying up modules...$(RESET)"
	@go mod tidy

# Build production binary
build:
	@echo "$(GREEN)Building $(APP_NAME)...$(RESET)"
	@mkdir -p bin
	@go build -o $(BINARY) ./cmd/scout

# Run in production mode
run:build
	@echo "$(GREEN)Running $(APP_NAME)...$(RESET)"
	@./$(BINARY)

# Clean build artifacts
clean:
	@echo "$(GREEN)Cleaning build artifacts...$(RESET)"
	@rm -rf $(BINARY) tmp
