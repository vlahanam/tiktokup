.PHONY: help up down restart build dev dev-down logs ps test test-backend test-python lint-frontend typecheck-frontend shadcn-add db-shell db-status db-migrate python-shell python-run clean

COMPOSE := docker compose
COMPOSE_DEV := docker compose -f docker-compose.yml -f docker-compose.dev.yml

help: ## Show available Docker-only commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

up: ## Build and start production-like stack via Docker Compose
	$(COMPOSE) up -d --build --wait

down: ## Stop and remove containers
	$(COMPOSE) down

restart: ## Restart running services
	$(COMPOSE) restart

build: ## Build all Docker images
	$(COMPOSE) build

dev: ## Start development stack with hot reload (Air + Next.js + FastAPI + Postgres)
	$(COMPOSE_DEV) up --build

dev-down: ## Stop development stack
	$(COMPOSE_DEV) down

logs: ## Tail logs from all services
	$(COMPOSE) logs -f

ps: ## Show container status and health
	$(COMPOSE) ps

test-backend: ## Run Go 1.26 unit tests and go vet inside Docker
	docker run --rm -v "$(CURDIR)/backend:/app" -w /app golang:1.26-alpine sh -c "go vet ./... && go test -v -count=1 ./..."

test-python: ## Run Ruff lint and Pytest for Python service inside Docker
	docker build -q --target dev -t tiktokup-python:test ./python >/dev/null && docker run --rm tiktokup-python:test sh -c "ruff check . && pytest -v"

lint-frontend: ## Run ESLint on frontend inside Docker
	docker run --rm -v "$(CURDIR)/frontend:/app" -w /app node:22-alpine npm run lint

typecheck-frontend: ## Run TypeScript typecheck on frontend inside Docker
	docker run --rm -v "$(CURDIR)/frontend:/app" -w /app node:22-alpine npx tsc --noEmit

test: test-backend test-python typecheck-frontend lint-frontend ## Run all backend, python & frontend checks via Docker

shadcn-add: ## Add a Shadcn UI component via Docker (usage: make shadcn-add COMPONENT=table)
	@test -n "$(COMPONENT)" || (echo "Usage: make shadcn-add COMPONENT=<name>" && exit 1)
	docker run --rm -u "$$(id -u):$$(id -g)" -e HOME=/tmp -e npm_config_cache=/tmp/.npm -v "$(CURDIR)/frontend:/app" -w /app node:22-alpine npx --yes shadcn@latest add $(COMPONENT) --yes

db-shell: ## Open an interactive psql shell in the PostgreSQL container
	$(COMPOSE) exec postgres sh -c 'psql -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"'

db-status: ## Check PostgreSQL readiness and applied schema migrations
	$(COMPOSE) exec -T postgres sh -c 'pg_isready -U "$$POSTGRES_USER" -d "$$POSTGRES_DB" && psql -U "$$POSTGRES_USER" -d "$$POSTGRES_DB" -c "SELECT version, applied_at FROM schema_migrations ORDER BY version;"'

db-migrate: ## Run SQL migrations using the backend container binary
	$(COMPOSE) exec -T backend /app/migrate

python-shell: ## Open an interactive shell inside the Python container
	$(COMPOSE) exec python bash

python-run: ## Run a Python script inside the container (usage: make python-run SCRIPT=scripts/example.py)
	@test -n "$(SCRIPT)" || (echo "Usage: make python-run SCRIPT=scripts/example.py" && exit 1)
	$(COMPOSE) exec -T python python "$(SCRIPT)"

clean: ## Stop containers and remove volumes/build artifacts (destroys DB data!)
	$(COMPOSE_DEV) down -v --remove-orphans
