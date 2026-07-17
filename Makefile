# Load variables from .env when it exists so `make` targets can reuse them.
ifneq (,$(wildcard ./.env))
include .env
export
endif

# Postgres URL built from the same variables the app uses.
DATABASE_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
MIGRATIONS_DIR := db/migrations

TEST_DATABASE_URL ?= postgres://sauser:sapass@localhost:5444/sadb_test?sslmode=disable

.PHONY: help run build format lint test test-integration tidy migrate-up migrate-down migrate-create docker-up docker-down test-db-up test-db-down

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'

run: ## Run the server locally
	go run ./cmd/short-access

build: ## Build the server binary into ./bin
	go build -o bin/short-access ./cmd/short-access

format: ## Format all Go code
	go fmt ./...

lint: ## Run the linter (needs golangci-lint installed)
	golangci-lint run

test: ## Run all tests (database tests skip unless TEST_DATABASE_URL is set)
	go test ./...

test-integration: ## Run tests against the local test database (run test-db-up first)
	TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test ./...

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

test-db-up: ## Start the throwaway test database
	docker compose --profile test up -d test-db

test-db-down: ## Stop and remove the throwaway test database
	docker compose --profile test rm -sf test-db
