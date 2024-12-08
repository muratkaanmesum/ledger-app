# Project variables
APP_NAME := ptm
MAIN_FILE := cmd/main.go
BUILD_DIR := ./tmp
BUILD_FILE := $(BUILD_DIR)/$(APP_NAME)
DOCKER_COMPOSE := docker-compose

# Environment variables
ENV_FILE := .env
DB_CONTAINER := ptm-postgres
REDIS_CONTAINER := ptm-redis

# Default target
.DEFAULT_GOAL := help

# Build the application
.PHONY: build
build: ## Build the Go application
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_FILE) $(MAIN_FILE)
	@echo "Built $(BUILD_FILE)"

# Run the application locally
.PHONY: run
run: build ## Run the application locally
	@echo "Running $(APP_NAME)..."
	$(BUILD_FILE)

# Clean build artifacts
.PHONY: clean
clean: ## Remove build artifacts
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned build artifacts"

# Run the application with Docker Compose
.PHONY: docker-up
docker-up: ## Start the application with Docker Compose
	$(DOCKER_COMPOSE) up --build

# Stop the application with Docker Compose
.PHONY: docker-down
docker-down: ## Stop the application and remove containers
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

# Rebuild and restart Docker containers
.PHONY: docker-rebuild
docker-rebuild: ## Rebuild and restart the Docker containers
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	$(DOCKER_COMPOSE) up --build

# Check the status of Docker services
.PHONY: docker-ps
docker-ps: ## Check the status of running Docker containers
	$(DOCKER_COMPOSE) ps

# Access the PostgreSQL database container
.PHONY: db-shell
db-shell: ## Access the PostgreSQL container
	$(DOCKER_COMPOSE) exec $(DB_CONTAINER) psql -U postgres -d ptmdb

# Access the Redis container
.PHONY: redis-shell
redis-shell: ## Access the Redis container
	$(DOCKER_COMPOSE) exec $(REDIS_CONTAINER) redis-cli

# Format the code
.PHONY: fmt
fmt: ## Format the code with gofmt
	@gofmt -s -w .

# Run lint checks
.PHONY: lint
lint: ## Run lint checks
	@golangci-lint run

# Show help message
.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'