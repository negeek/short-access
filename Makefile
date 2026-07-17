# Load variables from .env when it exists so `make` targets can reuse them.
ifneq (,$(wildcard ./.env))
include .env
export
endif

# Postgres URL built from the same variables the app uses.
DATABASE_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
MIGRATIONS_DIR := db/migrations

.PHONY: help run build fmt lint test tidy migrate-up migrate-down migrate-create docker-up docker-down

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'

run: ## Run the server locally
	go run ./cmd/short-access

build: ## Build the server binary into ./bin
	go build -o bin/short-access ./cmd/short-access

fmt: ## Format all Go code
	go fmt ./...

lint: ## Run the linter (needs golangci-lint installed)
	golangci-lint run

test: ## Run all tests
	go test ./...

tidy: ## Tidy up go.mod and go.sum
	go mod tidy

migrate-up: ## Apply all pending migrations
	go run ./cmd/short-access migrate up

migrate-down: ## Roll back the last migration
	go run ./cmd/short-access migrate down

migrate-create: ## Create a new migration, e.g. make migrate-create name=add_widgets (needs golang-migrate)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

docker-up: ## Start the stack with Docker Compose
	docker compose up --build

docker-down: ## Stop the stack and remove containers
	docker compose down
