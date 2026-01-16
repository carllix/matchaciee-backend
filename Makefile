.PHONY: help run build test clean dev migrate-up migrate-down migrate-create migrate-version migrate-drop migrate-force swagger swagger-fmt

# Variables
APP_NAME=matchaciee-api
BINARY_NAME=bin/api
MAIN_PATH=cmd/api/main.go

help:
	@echo "Available commands:"
	@echo "  make run             - Run the application"
	@echo "  make dev             - Run with hot reload"
	@echo "  make build           - Build the application"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Remove build artifacts"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt             - Format code"
	@echo "  make fix-align       - Fix field alignment for better memory usage"
	@echo "  make vet             - Run go vet"
	@echo "  make lint            - Run linter"
	@echo "  make ci              - Run all CI checks (lint, vet, test)"
	@echo ""
	@echo "Documentation:"
	@echo "  make swagger         - Generate Swagger documentation"
	@echo "  make swagger-fmt     - Format Swagger annotations"
	@echo ""
	@echo "Database Migrations:"
	@echo "  make migrate-up      - Run all pending migrations"
	@echo "  make migrate-down    - Rollback the last migration"
	@echo "  make migrate-version - Show current migration version"
	@echo "  make migrate-create  - Create new migration (use name=<migration_name>)"
	@echo "  make migrate-drop    - Drop all tables (dangerous!)"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-dev-up   - Start development database"
	@echo "  make docker-dev-down - Stop development database"

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

# Fix field alignment
fix-align:
	@echo "Fixing field alignment..."
	@fieldalignment -fix ./...

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

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger documentation generated in docs/"

# Format Swagger annotations
swagger-fmt:
	@echo "Formatting Swagger annotations..."
	@swag fmt
	@echo "Swagger annotations formatted"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
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

# Database Migration commands
migrate-up:
	@echo "Running database migrations..."
	@./scripts/migrate.sh up

migrate-down:
	@echo "Rolling back last migration..."
	@./scripts/migrate.sh down

migrate-version:
	@./scripts/migrate.sh version

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Please specify migration name"; \
		echo "Usage: make migrate-create name=<migration_name>"; \
		exit 1; \
	fi
	@./scripts/migrate.sh create $(name)

migrate-drop:
	@./scripts/migrate.sh drop

migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Error: Please specify version number"; \
		echo "Usage: make migrate-force version=<version>"; \
		exit 1; \
	fi
	@./scripts/migrate.sh force $(version)