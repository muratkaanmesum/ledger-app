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

.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	$(BUILD_FILE)

.PHONY: clean
clean: ## Remove build artifacts
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned build artifacts"

.PHONY: docker-up
docker-up:
	$(DOCKER_COMPOSE) up --build

.PHONY: docker-down
docker-down:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

.PHONY: docker-rebuild
docker-rebuild:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	$(DOCKER_COMPOSE) up --build

.PHONY: docker-ps
docker-ps:
	$(DOCKER_COMPOSE) ps

.PHONY: db-shell
db-shell:
	$(DOCKER_COMPOSE) exec $(DB_CONTAINER) psql -U postgres -d ptmdb

.PHONY: redis-shell
redis-shell:
	$(DOCKER_COMPOSE) exec $(REDIS_CONTAINER) redis-cli

.PHONY: fmt
fmt:
	@gofmt -s -w .

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: debug
debug:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	DEBUG=true $(DOCKER_COMPOSE) up --build
	@echo "Application started in debug mode. Attach your debugger to port 40000."

# Show help message
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'