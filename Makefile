# ========================================
# 💼 Scout Makefile (Simplified)
# ========================================

APP_NAME := scout
CORE_BINARY := bin/scout-core
WRAPPER_BINARY := bin/$(APP_NAME)

PWD := $(shell pwd)
LIB_PATH := $(PWD)/.scout/llama
MODEL_PATH := $(PWD)/.scout/model

GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
RESET := \033[0m

.DEFAULT_GOAL := help

help:
	@echo "$(YELLOW)Available commands:$(RESET)"
	@echo "  make setup           📦 First time setup"
	@echo "  make check           🧐 Verify files exist"
	@echo "  make build           🏗️  Build binary & wrapper"
	@echo "  make run             🚀 Run Scout"
	@echo "  make clean           🧼 Remove artifacts"

setup:
	@mkdir -p .scout/llama
	@mkdir -p .scout/model
	@echo "$(GREEN)✅ Created .scout directory structure.$(RESET)"
	@echo "$(YELLOW)⚠️  ACTION REQUIRED:$(RESET)"
	@echo "1. Download Model to: $(MODEL_PATH)/llama-3.2-3b-instruct-q4_k_m.gguf"
	@echo "2. Copy 'libllama.dylib' (or .so) to: $(LIB_PATH)/"

check:
	@echo "$(CYAN)Checking configuration...$(RESET)"
	@ls -lh $(LIB_PATH)/libllama.* 2>/dev/null && echo "$(GREEN)✅ Library found.$(RESET)" || echo "$(RED)❌ Library missing in .scout/llama/$(RESET)"
	@ls -lh $(MODEL_PATH)/*.gguf 2>/dev/null && echo "$(GREEN)✅ Model found.$(RESET)" || echo "$(RED)❌ Model missing in .scout/model/$(RESET)"

build:
	@echo "$(GREEN)Building $(APP_NAME)...$(RESET)"
	@mkdir -p bin
	@go build -o $(CORE_BINARY) ./cmd/scout
	@cp scout.sh $(WRAPPER_BINARY)
	@chmod +x $(WRAPPER_BINARY)
	@echo "$(GREEN)✅ Build complete!$(RESET)"
	@echo "   Run it with: ./bin/scout"

run: build
	@echo "$(GREEN)🚀 Launching $(APP_NAME)...$(RESET)"
	@$(WRAPPER_BINARY)

clean:
	@rm -rf bin tmp