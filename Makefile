COMPOSE ?= docker compose
SERVICE_API ?= api
SERVICE_DB ?= db
SERVICE_FLYWAY ?= flyway

.PHONY: up
up:
	$(COMPOSE) up --build

.PHONY: up-d
up-d:
	$(COMPOSE) up -d --build

.PHONY: down
down:
	$(COMPOSE) down

.PHONY: down-v
down-v:
	$(COMPOSE) down -v

.PHONY: logs
logs:
	$(COMPOSE) logs -f

.PHONY: logs-api
logs-api:
	$(COMPOSE) logs -f $(SERVICE_API)

.PHONY: logs-db
logs-db:
	$(COMPOSE) logs -f $(SERVICE_DB)

.PHONY: logs-flyway
logs-flyway:
	$(COMPOSE) logs -f $(SERVICE_FLYWAY)

.PHONY: build-api
build-api:
	$(COMPOSE) build $(SERVICE_API)

.PHONY: restart-api
restart-api:
	$(COMPOSE) restart $(SERVICE_API)

.PHONY: sh-api
sh-api:
	$(COMPOSE) exec $(SERVICE_API) sh -c '$(cmd)'

.PHONY: psql
psql:
	$(COMPOSE) exec $(SERVICE_DB) psql -U $${DB_USER:-app} -d $${DB_NAME:-appdb} -c '$(cmd)'

.PHONY: flyway-info
flyway-info:
	$(COMPOSE) run --rm $(SERVICE_FLYWAY) info

.PHONY: flyway-migrate
flyway-migrate:
	$(COMPOSE) run --rm $(SERVICE_FLYWAY) migrate

.PHONY: flyway-clean
flyway-clean:
	$(COMPOSE) run --rm $(SERVICE_FLYWAY) clean

.PHONY: swagger
swagger:
	swag init -g cmd/api/main.go -o internal/docs

.PHONY: test
test:
	go test ./...

.PHONY: help
help:
	@echo "Comandos disponíveis:"
	@echo "  make up            - Sobe todos os serviços (build incluso)"
	@echo "  make up-d          - Sobe em background"
	@echo "  make down          - Para e remove containers"
	@echo "  make down-v        - Para e remove containers e volumes (reset DB)"
	@echo "  make logs          - Segue logs de todos os serviços"
	@echo "  make logs-api      - Logs da API"
	@echo "  make build-api     - Rebuild da API"
	@echo "  make restart-api   - Restart da API"
	@echo "  make psql cmd='\\l' - Executa comando psql no DB"
	@echo "  make flyway-info   - Status das migrations"
	@echo "  make flyway-migrate- Roda as migrations"
	@echo "  make flyway-clean  - Limpa o schema (CUIDADO)"
	@echo "  make swagger       - Gera Swagger localmente"
	@echo "  make test          - Roda testes Go"