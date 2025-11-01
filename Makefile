.PHONY: help install swagger build run test docker-up docker-down migrate-up migrate-info logs

APP_NAME=takehome-go
DOCKER_COMPOSE=docker-compose
GO=go
SWAG=swag

help:
	@echo "Comandos essenciais:"
	@echo "  make install       - Instala deps Go e ferramentas (swag)"
	@echo "  make swagger       - Gera docs Swagger em ./docs"
	@echo "  make build         - Compila a aplicação"
	@echo "  make run           - Roda a aplicação localmente"
	@echo "  make test          - Executa testes"
	@echo "  make docker-up     - Sobe Postgres, aplica migrations (Flyway) e API"
	@echo "  make docker-down   - Para os containers"
	@echo "  make migrate-up    - Executa migrations (Flyway)"
	@echo "  make migrate-info  - Mostra status das migrations"
	@echo "  make logs          - Segue logs dos containers"

install:
	$(GO) mod download
	$(GO) mod tidy
	$(GO) install github.com/swaggo/swag/cmd/swag@latest

swagger:
	$(SWAG) init -g cmd/api/main.go -o docs

build: swagger
	$(GO) build -o bin/$(APP_NAME) cmd/api/main.go

run: swagger
	$(GO) run cmd/api/main.go

test:
	$(GO) test ./...

docker-up:
	$(DOCKER_COMPOSE) up -d --build

docker-down:
	$(DOCKER_COMPOSE) down -v

migrate-up:
	$(DOCKER_COMPOSE) run --rm flyway migrate

migrate-info:
	$(DOCKER_COMPOSE) run --rm flyway info

logs:
	$(DOCKER_COMPOSE) logs -f