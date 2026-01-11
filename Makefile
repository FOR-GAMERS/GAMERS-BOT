.PHONY: help build run test clean docker-build docker-up docker-stop docker-down fmt lint vet deps tidy setup env

# Variables
APP_NAME := gamers-discord-bot
BIN_DIR := bin
DOCKER_IMAGE := $(APP_NAME):latest
GO_FILES := $(shell find . -name '*.go' -type f)

# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "$(COLOR_BOLD)GAMERS Discord Bot - Makefile Commands$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_GREEN)Build & Run:$(COLOR_RESET)"
	@echo "  $(COLOR_YELLOW)make build$(COLOR_RESET)           - Build the application binary"
	@echo "  $(COLOR_YELLOW)make run$(COLOR_RESET)             - Run the application locally"
	@echo "  $(COLOR_YELLOW)make dev$(COLOR_RESET)             - Run in development mode with auto-reload"
	@echo "  $(COLOR_YELLOW)make clean$(COLOR_RESET)           - Remove build artifacts"
	@echo ""
	@echo "$(COLOR_GREEN)Testing:$(COLOR_RESET)"
	@echo "  $(COLOR_YELLOW)make test$(COLOR_RESET)            - Run all tests"
	@echo "  $(COLOR_YELLOW)make test-verbose$(COLOR_RESET)    - Run tests with verbose output"
	@echo "  $(COLOR_YELLOW)make test-coverage$(COLOR_RESET)   - Run tests with coverage report"
	@echo ""
	@echo "$(COLOR_GREEN)Code Quality:$(COLOR_RESET)"
	@echo "  $(COLOR_YELLOW)make fmt$(COLOR_RESET)             - Format Go code"
	@echo "  $(COLOR_YELLOW)make lint$(COLOR_RESET)            - Run linter (requires golangci-lint)"
	@echo "  $(COLOR_YELLOW)make vet$(COLOR_RESET)             - Run go vet"
	@echo "  $(COLOR_YELLOW)make check$(COLOR_RESET)           - Run fmt, vet, and lint"
	@echo ""
	@echo "$(COLOR_GREEN)Dependencies:$(COLOR_RESET)"
	@echo "  $(COLOR_YELLOW)make deps$(COLOR_RESET)            - Download dependencies"
	@echo "  $(COLOR_YELLOW)make deps-update$(COLOR_RESET)     - Update dependencies"
	@echo "  $(COLOR_YELLOW)make tidy$(COLOR_RESET)            - Tidy go.mod and go.sum"
	@echo ""
	@echo "$(COLOR_GREEN)Docker:$(COLOR_RESET)"
	@echo "  $(COLOR_YELLOW)make docker-build$(COLOR_RESET)    - Build Docker image"
	@echo "  $(COLOR_YELLOW)make docker-up$(COLOR_RESET)       - Start Docker container with docker-compose"
	@echo "  $(COLOR_YELLOW)make docker-stop$(COLOR_RESET)     - Stop Docker container"
	@echo "  $(COLOR_YELLOW)make docker-down$(COLOR_RESET)     - Stop and remove all Docker resources (containers, networks, volumes)"
	@echo "  $(COLOR_YELLOW)make docker-clean$(COLOR_RESET)    - Remove Docker images and containers"
	@echo "  $(COLOR_YELLOW)make docker-logs$(COLOR_RESET)     - Show Docker container logs"
	@echo ""
	@echo "$(COLOR_GREEN)Setup & Utilities:$(COLOR_RESET)"
	@echo "  $(COLOR_YELLOW)make setup$(COLOR_RESET)           - Initial project setup"
	@echo "  $(COLOR_YELLOW)make env$(COLOR_RESET)             - Copy .env.example to .env"
	@echo "  $(COLOR_YELLOW)make install-tools$(COLOR_RESET)   - Install development tools"
	@echo ""

## build: Build the application
build:
	@echo "$(COLOR_BLUE)Building $(APP_NAME)...$(COLOR_RESET)"
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) ./cmd
	@echo "$(COLOR_GREEN)✓ Build complete: $(BIN_DIR)/$(APP_NAME)$(COLOR_RESET)"

## build-linux: Build for Linux (useful for Docker)
build-linux:
	@echo "$(COLOR_BLUE)Building $(APP_NAME) for Linux...$(COLOR_RESET)"
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-linux ./cmd
	@echo "$(COLOR_GREEN)✓ Linux build complete: $(BIN_DIR)/$(APP_NAME)-linux$(COLOR_RESET)"

## run: Run the application
run:
	@echo "$(COLOR_BLUE)Running $(APP_NAME)...$(COLOR_RESET)"
	@go run ./cmd/main.go

## dev: Run in development mode
dev:
	@echo "$(COLOR_BLUE)Running in development mode...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Tip: Install 'air' for auto-reload: go install github.com/cosmtrek/air@latest$(COLOR_RESET)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(COLOR_YELLOW)Air not found, running without auto-reload$(COLOR_RESET)"; \
		go run ./cmd/main.go; \
	fi

## clean: Remove build artifacts
clean:
	@echo "$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

## test: Run tests
test:
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	@go test ./... -v

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "$(COLOR_BLUE)Running tests (verbose)...$(COLOR_RESET)"
	@go test ./... -v -count=1

## test-coverage: Run tests with coverage
test-coverage:
	@echo "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)"
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✓ Coverage report generated: coverage.html$(COLOR_RESET)"

## fmt: Format Go code
fmt:
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	@go fmt ./...
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

## lint: Run linter
lint:
	@echo "$(COLOR_BLUE)Running linter...$(COLOR_RESET)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
		echo "$(COLOR_GREEN)✓ Linting complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)golangci-lint not installed. Install: https://golangci-lint.run/usage/install/$(COLOR_RESET)"; \
	fi

## vet: Run go vet
vet:
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	@go vet ./...
	@echo "$(COLOR_GREEN)✓ Vet complete$(COLOR_RESET)"

## check: Run all code quality checks
check: fmt vet lint
	@echo "$(COLOR_GREEN)✓ All checks passed$(COLOR_RESET)"

## deps: Download dependencies
deps:
	@echo "$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)"
	@go mod download
	@echo "$(COLOR_GREEN)✓ Dependencies downloaded$(COLOR_RESET)"

## deps-update: Update dependencies
deps-update:
	@echo "$(COLOR_BLUE)Updating dependencies...$(COLOR_RESET)"
	@go get -u ./...
	@go mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies updated$(COLOR_RESET)"

## tidy: Tidy go.mod and go.sum
tidy:
	@echo "$(COLOR_BLUE)Tidying go.mod...$(COLOR_RESET)"
	@go mod tidy
	@echo "$(COLOR_GREEN)✓ go.mod tidied$(COLOR_RESET)"

## docker-build: Build Docker image
docker-build:
	@echo "$(COLOR_BLUE)Building Docker image...$(COLOR_RESET)"
	@docker build -f docker/Dockerfile -t $(DOCKER_IMAGE) .
	@echo "$(COLOR_GREEN)✓ Docker image built: $(DOCKER_IMAGE)$(COLOR_RESET)"

## docker-up: Start Docker container with docker-compose
docker-up:
	@echo "$(COLOR_BLUE)Starting Docker container with docker-compose...$(COLOR_RESET)"
	@if [ ! -f env/.env ]; then \
		echo "$(COLOR_YELLOW)Warning: env/.env file not found. Run 'make env' first$(COLOR_RESET)"; \
		exit 1; \
	fi
	@docker network inspect gamers-network >/dev/null 2>&1 || docker network create gamers-network
	@docker compose -f docker/docker-compose.yml up -d
	@echo "$(COLOR_GREEN)✓ Container started: $(APP_NAME)$(COLOR_RESET)"

## docker-stop: Stop Docker container
docker-stop:
	@echo "$(COLOR_BLUE)Stopping Docker container...$(COLOR_RESET)"
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true
	@echo "$(COLOR_GREEN)✓ Container stopped$(COLOR_RESET)"

## docker-down: Stop and remove all Docker resources
docker-down:
	@echo "$(COLOR_BLUE)Stopping and removing all Docker resources...$(COLOR_RESET)"
	@docker compose -f docker/docker-compose.yml down -v --remove-orphans
	@echo "$(COLOR_GREEN)✓ All Docker resources stopped and removed$(COLOR_RESET)"

## docker-clean: Remove Docker images and containers
docker-clean: docker-stop
	@echo "$(COLOR_BLUE)Cleaning Docker resources...$(COLOR_RESET)"
	@docker rmi $(DOCKER_IMAGE) || true
	@echo "$(COLOR_GREEN)✓ Docker resources cleaned$(COLOR_RESET)"

## docker-logs: Show Docker container logs
docker-logs:
	@docker logs -f $(APP_NAME)

## setup: Initial project setup
setup: env deps
	@echo "$(COLOR_GREEN)✓ Project setup complete$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Next steps:$(COLOR_RESET)"
	@echo "  1. Edit env/.env file with your configuration"
	@echo "  2. Run 'make docker-up' to start the bot"

## env: Copy .env.example to .env
env:
	@if [ -f env/.env ]; then \
		echo "$(COLOR_YELLOW)env/.env file already exists$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_BLUE)Creating env/.env file...$(COLOR_RESET)"; \
		cp env/.env.example env/.env; \
		echo "$(COLOR_GREEN)✓ env/.env file created$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)Please edit env/.env with your configuration$(COLOR_RESET)"; \
	fi

## install-tools: Install development tools
install-tools:
	@echo "$(COLOR_BLUE)Installing development tools...$(COLOR_RESET)"
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing air (hot reload)..."
	@go install github.com/cosmtrek/air@latest
	@echo "$(COLOR_GREEN)✓ Development tools installed$(COLOR_RESET)"

