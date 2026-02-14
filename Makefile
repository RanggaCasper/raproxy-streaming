.PHONY: help build run test clean docker-build docker-run

help: ## Show this help message
	@echo Usage: make [target]
	@echo Available targets:
	@echo   build           Build the application
	@echo   run             Run the application
	@echo   test            Run tests
	@echo   clean           Clean build artifacts
	@echo   deps            Download dependencies
	@echo   lint            Run linter
	@echo   docker-build    Build Docker image
	@echo   docker-run      Run Docker container
	@echo   dev             Run in development mode (air)
	@echo   install-tools   Install development tools

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/server cmd/server/main.go
	@echo "Build complete: bin/server"

run: ## Run the application
	@echo "Starting server..."
	@go run cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean
	@echo "Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t raproxy-streaming:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p 3000:3000 raproxy-streaming:latest

dev: ## Run in development mode with live reload (requires air)
	@echo "Starting development server..."
	@air

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
