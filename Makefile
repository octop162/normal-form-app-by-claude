# Makefile for normal-form-app

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
BINARY_NAME=normal-form-app
BINARY_UNIX=$(BINARY_NAME)_unix

# Linting tools
GOLANGCI_LINT=golangci-lint

# Main directories
CMD_DIR=./cmd/server
BUILD_DIR=./build

.PHONY: help build clean test coverage lint fmt vet deps tidy run dev install-tools check-tools

# Default target
all: clean deps test lint build

# Help target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the application
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(CMD_DIR)

# Build for Linux
build-linux: ## Build the application for Linux
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v $(CMD_DIR)

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f logs/*.log
	@rm -f logs/*.pid

# Run tests
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint the code
lint: ## Run golangci-lint
	@echo "Running linters..."
	$(GOLANGCI_LINT) run ./...

# Format the code
fmt: ## Format the code
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Run go vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Download dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Tidy dependencies
tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Run the application
run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run $(CMD_DIR)

# Development mode (with auto-reload)
dev: ## Run in development mode
	@echo "Starting development environment..."
	@./scripts/dev-start.sh

# Stop development environment
dev-stop: ## Stop development environment
	@echo "Stopping development environment..."
	@./scripts/dev-stop.sh

# Health check
health: ## Run health check
	@echo "Running health check..."
	@./scripts/health-check.sh

# Install development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	}
	@command -v goimports >/dev/null 2>&1 || { \
		echo "Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	}
	@command -v staticcheck >/dev/null 2>&1 || { \
		echo "Installing staticcheck..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
	}
	@command -v gosec >/dev/null 2>&1 || { \
		echo "Installing gosec..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	}
	@echo "Development tools installed successfully!"

# Check if tools are installed
check-tools: ## Check if development tools are installed
	@echo "Checking development tools..."
	@command -v go >/dev/null 2>&1 || { echo "Go is not installed"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "Docker is not installed"; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose is not installed"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "Node.js is not installed"; exit 1; }
	@command -v npm >/dev/null 2>&1 || { echo "npm is not installed"; exit 1; }
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is not installed. Run 'make install-tools'"; exit 1; }
	@echo "All required tools are installed!"

# Full check (format, vet, lint, test)
check: fmt vet lint test ## Run all checks (format, vet, lint, test)

# CI/CD pipeline simulation
ci: clean deps check build ## Simulate CI/CD pipeline

# Docker operations
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME):latest

# Database operations
db-up: ## Start database
	@echo "Starting database..."
	docker-compose up -d postgres

db-down: ## Stop database
	@echo "Stopping database..."
	docker-compose stop postgres

db-reset: ## Reset database
	@echo "Resetting database..."
	docker-compose down -v postgres
	docker-compose up -d postgres

# Environment setup
setup: install-tools deps ## Setup development environment
	@echo "Setting up development environment..."
	@cp -n .env.example .env 2>/dev/null || true
	@echo "Development environment setup complete!"
	@echo "Run 'make dev' to start the development server"

# Show environment info
info: ## Show environment information
	@echo "Environment Information:"
	@echo "======================="
	@echo "Go version: $$(go version)"
	@echo "Docker version: $$(docker --version)"
	@echo "Node.js version: $$(node --version)"
	@echo "npm version: $$(npm --version)"
	@echo "Project root: $$(pwd)"
	@echo "Git branch: $$(git branch --show-current 2>/dev/null || echo 'Not a git repository')"
	@echo "Git status: $$(git status --porcelain 2>/dev/null | wc -l || echo 'N/A') files changed"