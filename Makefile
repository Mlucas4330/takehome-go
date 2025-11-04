.PHONY: help up down logs rebuild migrate swagger clean deps

help:
	@echo Comandos disponiveis:
	@echo   deps      - Instala dependências Go e ferramentas (swag)
	@echo   up        - Inicia todos os serviços
	@echo   down      - Para todos os serviços
	@echo   logs      - Mostra os logs
	@echo   rebuild   - Reconstrói e reinicia os serviços
	@echo   migrate   - Executa as migrações (Flyway)
	@echo   swagger   - Gera/atualiza a documentação Swagger em ./docs (host)
	@echo   clean     - Remove containers e volumes

deps:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest

up:
	docker compose up -d

down:
	docker compose down --remove-orphans

logs:
	docker compose logs -f

rebuild:
	docker compose down
	docker compose up -d --build

migrate:
	docker compose run --rm flyway migrate

swagger:
	swag init -g cmd/api/main.go -o ./docs

clean:
	docker compose down -v

.DEFAULT_GOAL := help