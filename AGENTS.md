# AGENTS.md — TikTokUp Repository Guide

## Project Overview
TikTokUp is a 100% Dockerized full-stack web platform:
- **Backend (`backend/`)**: Go `1.26` (`module tiktokup/backend`), stdlib `net/http` + `log/slog` + `pgx/v5` (`v5.13.0`), layered architecture (`cmd/{server,migrate}`, `internal/{config,database,handler,middleware,router,server}`, `migrations/`, `pkg/response`).
- **Database (`postgres`)**: PostgreSQL `18` (`postgres:18-alpine`, `18.6`), persistent volume `postgres_data:/var/lib/postgresql`, embedded transactional SQL migrations tracked in `schema_migrations`.
- **Frontend (`frontend/`)**: Next.js `16.4.6` (App Router, RSC, `output: "standalone"`), React `19.3.0`, TypeScript `5`, Tailwind CSS `4.3.3`, Shadcn UI `4.23.1` (`style: "base-nova"`, `@base-ui/react`).
- **Reverse Proxy (`nginx/`)**: `nginx:alpine` routing `/api/*` and `/health` to `backend:8080`, and `/*` to `frontend:3000`.

## Toolchain Rule (Docker-First)
Do **not** assume Go, Node.js, or npm are installed on the host. Run all language commands inside Docker or via `Makefile`:
- Backend checks: `make test-backend` (`golang:1.26-alpine` running `go vet ./... && go test -v -count=1 ./...`)
- Frontend checks: `make typecheck-frontend` & `make lint-frontend` (`node:22-alpine`)
- Database commands: `make db-status`, `make db-migrate`, `make db-shell`
- Add Shadcn UI components: `make shadcn-add COMPONENT=<name>`
- Start production-like stack: `docker compose up -d --build --wait`
- Start dev hot-reload stack: `docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build` (exposes Postgres on `POSTGRES_PORT`, default `5433`)

## Coding & Architecture Conventions
1. **Backend (Go 1.26 + PostgreSQL 18)**:
   - Use `pkg/response` (`OK`, `Created`, `Error`) for all HTTP responses so the JSON envelope (`success`, `data`, `error`, `meta`) stays consistent.
   - Always log with `log/slog` and include `request_id` from `response.RequestIDFromContext(r.Context())`.
   - Keep `cmd/server/main.go` thin; place new domain code under `internal/{handler,service,repository}` and versioned SQL files in `backend/migrations/` (lexicographical `000002_*.sql`).
   - Keep handlers testable by injecting interfaces (e.g., `handler.DatabaseInspector`) instead of concrete `*pgxpool.Pool` where appropriate.
2. **Frontend (Next.js 16 + Shadcn UI v4)**:
   - Shadcn UI uses `@base-ui/react` (`style: "base-nova"`). Use the `render` prop instead of Radix `asChild` (and pass `nativeButton={false}` when rendering `<Link>` inside `<Button>`).
   - Keep API contracts synchronized between `backend/pkg/response/json.go` + `backend/internal/handler/*` and `frontend/src/types/api.ts`.
   - Use `src/lib/api.ts` for API calls so SSR uses `INTERNAL_API_URL` (`http://backend:8080`) and browser calls use `NEXT_PUBLIC_API_URL` (`/api`).
   - Respect React 19 ESLint rules (avoid calling `setState` synchronously inside `useEffect`).
3. **Nginx & Networking**:
   - Keep `backend` and `frontend` internal (`expose` only) in `docker-compose.yml`; route external traffic through `nginx` (`NGINX_PORT`, default `8000`).
   - Always preserve/propagate `X-Request-ID`.

## Documentation
Keep `docs/` in sync whenever architecture, endpoints, env vars, or dependencies change. See [docs/codebase-summary.md](docs/codebase-summary.md) and [docs/system-architecture.md](docs/system-architecture.md).
