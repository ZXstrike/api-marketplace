# Makefile for API Marketplace

# Variables
GATEWAY_DIR := ./api-gateway
MARKET_DIR := ./marketplace-app
SHARED_DIR := ./shared

# Go commands
GO := go
GO_RUN := $(GO) run
GO_TEST := $(GO) test

.PHONY: help all build run-gateway run-market run-web test gen-keys dev-up dev-down prod-up prod-down build-prod migrate clean

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all                Build all services"
	@echo "  run-gateway        Run the API Gateway service"
	@echo "  run-market         Run the Marketplace App service"
	@echo "  run-web            Run the frontend development server"
	@echo "  test               Run tests for all Go services"
	@echo "  gen-keys           Generate ECDSA key pair for JWT"
	@echo "  dev-up             Start all services with Docker Compose"
	@echo "  dev-down           Stop all services started with Docker Compose"
	@echo "  prod-up            Start all services with Docker Compose (production)"
	@echo "  prod-down          Stop all production services"
	@echo "  build-prod         Build all production Docker images"
	@echo "  migrate            Run database migrations"
	@echo "  clean              Clean build artifacts"

all: build

build:
	@echo "Building services..."
	@$(GO) build -o $(GATEWAY_DIR)/bin/gateway $(GATEWAY_DIR)/cmd/main.go
	@$(GO) build -o $(MARKET_DIR)/bin/market $(MARKET_DIR)/cmd/main.go

run-gateway:
	@echo "Starting API Gateway..."
	@$(GO_RUN) $(GATEWAY_DIR)/cmd/main.go

run-market:
	@echo "Starting Marketplace App..."
	@$(GO_RUN) $(MARKET_DIR)/cmd/main.go

run-web:
	@echo "Starting frontend development server..."
	@cd marketplace-web && npm run dev

test:
	@echo "Running tests..."
	@$(GO_TEST) ./...

gen-keys:
	@echo "Generating ECDSA keys..."
	@$(GO_RUN) $(SHARED_DIR)/secrets/generate_keys.go

dev-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

dev-down:
	@echo "Stopping services..."
	@docker-compose down

prod-up:
	@echo "Starting production services with Docker Compose..."
	@docker-compose -f docker-compose.prod.yml up -d --no-deps --build
	@echo "Done..."

prod-down:
	@echo "Stopping production services..."
	@docker-compose -f docker-compose.prod.yml down

build-prod:
	@echo "Building production Docker images..."
	@docker-compose -f docker-compose.prod.yml build

migrate:
	@echo "Running database migrations..."
	@$(GO_RUN) $(MARKET_DIR)/cmd/main.go --migrate

clean:
	@echo "Cleaning up..."
	@rm -f $(GATEWAY_DIR)/bin/gateway
	@rm -f $(MARKET_DIR)/bin/market
