.PHONY: help up down logs ps api-test api-lint dash-install dash-lint dash-build migrate seed

COMPOSE = docker compose -p marketlens -f infra/docker/docker-compose.yml
API_DIR = services/api
DASH_DIR = apps/admin-dashboard

help:
	@echo "MarketLens commands:"
	@echo "  make up            Start local stack"
	@echo "  make down          Stop local stack"
	@echo "  make logs          Tail logs"
	@echo "  make ps            Show running containers"
	@echo "  make api-test       Run Go tests"
	@echo "  make api-lint       Run basic Go formatting checks"
	@echo "  make dash-install   Install dashboard deps"
	@echo "  make dash-lint      Lint dashboard"
	@echo "  make dash-build     Build dashboard"
# 	@echo "  make migrate        Run DB migrations (placeholder)"
	@echo "  make seed           Seed DB (placeholder)"

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

api-test:
	cd $(API_DIR) && go test ./... -race

api-lint:
	cd $(API_DIR) && test -z "$$(gofmt -l .)"

dash-install:
	cd $(DASH_DIR) && npm ci

dash-lint:
	cd $(DASH_DIR) && npm run lint

dash-build:
	cd $(DASH_DIR) && npm run build

migrate-up:
	cd services/api && go run ./cmd/migrate -action=up -dir=migrations

migrate-down:
	cd services/api && go run ./cmd/migrate -action=down -dir=migrations

migrate-status:
	cd services/api && go run ./cmd/migrate -action=status -dir=migrations

migrate-redo:
	cd services/api && go run ./cmd/migrate -action=redo -dir=migrations

seed:
	@echo "TODO: implement scripts/seed"; exit 1
