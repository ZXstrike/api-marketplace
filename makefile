# Makefile (must be named exactly “Makefile” with no extension)

.PHONY: help gen-keys dev-up run-gateway run-market migrate test

help:
	@echo "Available targets:"
	@echo "  make gen-keys       # Generate ECDSA key pair"
	@echo "  make dev-up         # podman-compose up -d"
	@echo "  make run-gateway    # go run api-gateway/cmd/main.go"
	@echo "  make run-market     # go run marketplace-app/cmd/main.go"
	@echo "  make migrate        # Run DB migrations"
	@echo "  make test           # Run all Go tests"

gen-keys:
	@./shared/secrets/generate_ecdsa_keys.sh

dev-up:
	@docker compose up -d

run-gateway:
	@cd api-gateway && go run cmd/main.go

run-market:
	@cd marketplace-app && go run cmd/main.go

migrate:
	@cd marketplace-app && go run cmd/migrate/main.go

test:
	@go test ./...
