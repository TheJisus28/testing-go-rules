.PHONY: help deps mocks swagger generate build run test test-cover db-up db-down

SWAG ?= swag

# Default: local binary. Alternative: make mocks MOCKERY="go run github.com/vektra/mockery/v2@v2.53.5"
MOCKERY ?= mockery

export GOTOOLCHAIN ?= auto

help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

deps: ## Download Go dependencies
	go mod tidy

mocks: ## Generate mocks with mockery (reads .mockery.yaml)
	$(MOCKERY)

swagger: ## Generate OpenAPI specs under docs/ with swag
	$(SWAG) init -g cmd/main.go -o docs --parseDependency --parseInternal

generate: mocks swagger ## Generate mocks and Swagger documentation

build: ## Build the binary
	go build -o bin/socialnet ./cmd

run: build ## Run the server
	./bin/socialnet

test: generate ## Run all tests (generates mocks first)
	go test ./...

test-cover: generate ## Per-domain coverage for service/
	@echo "=== users ==="
	@go test -cover ./internal/users/service/...
	@echo "=== auth ==="
	@go test -cover ./internal/auth/service/...
	@echo "=== posts ==="
	@go test -cover ./internal/posts/service/...
	@echo "=== friendships ==="
	@go test -cover ./internal/friendships/service/...

db-up: ## Start PostgreSQL with Docker
	docker compose up -d postgres

db-down: ## Stop PostgreSQL
	docker compose down
