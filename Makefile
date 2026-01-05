.PHONY: help run build test clean dev

# Variables
APP_NAME=matchaciee-api
BINARY_NAME=bin/api
MAIN_PATH=cmd/api/main.go

help:
	@echo "Available commands:"
	@echo "  make run          - Run the application"
	@echo "  make dev          - Run with hot reload"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Remove build artifacts"

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)

# Run with hot reload
dev:
	@echo "Running with hot reload..."
	@air

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean
	@echo "Clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download

# Lint code
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# CI checks (what CI runs)
ci: lint vet test
	@echo "CI checks passed!"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"

# Docker commands
docker-dev-up:
	@echo "Starting development database..."
	@docker-compose -f docker-compose.dev.yml up -d

docker-dev-down:
	@echo "Stopping development database..."
	@docker-compose -f docker-compose.dev.yml down

docker-dev-logs:
	@echo "Showing database logs..."
	@docker-compose -f docker-compose.dev.yml logs -f

docker-prod-up:
	@echo "Starting production services..."
	@docker-compose --env-file .env.production up -d --build

docker-prod-down:
	@echo "Stopping production services..."
	@docker-compose --env-file .env.production down