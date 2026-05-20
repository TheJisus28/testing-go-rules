.PHONY: help deps mocks swagger generate build run test test-cover db-up db-down

SWAG ?= swag

# Binario local por defecto; alternativa: make mocks MOCKERY="go run github.com/vektra/mockery/v2@v2.53.5"
MOCKERY ?= mockery

export GOTOOLCHAIN ?= auto

help: ## Muestra esta ayuda
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

deps: ## Descarga dependencias Go
	go mod tidy

mocks: ## Genera mocks con mockery (lee .mockery.yaml)
	$(MOCKERY)

swagger: ## Genera OpenAPI en docs/ con swag
	$(SWAG) init -g cmd/main.go -o docs --parseDependency --parseInternal

generate: mocks swagger ## Genera mocks y documentación Swagger

build: ## Compila el binario
	go build -o bin/socialnet ./cmd

run: build ## Ejecuta el servidor
	./bin/socialnet

test: generate ## Ejecuta todos los tests (genera mocks antes)
	go test ./...

test-cover: generate ## Cobertura por dominio en service/
	@echo "=== users ==="
	@go test -cover ./internal/users/service/...
	@echo "=== auth ==="
	@go test -cover ./internal/auth/service/...
	@echo "=== posts ==="
	@go test -cover ./internal/posts/service/...
	@echo "=== friendships ==="
	@go test -cover ./internal/friendships/service/...

db-up: ## Levanta PostgreSQL con Docker
	docker compose up -d postgres

db-down: ## Detiene PostgreSQL
	docker compose down
