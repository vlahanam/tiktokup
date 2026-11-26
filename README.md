# TikTokUp

Nền tảng web full-stack được đóng gói 100% bằng Docker, gồm **Golang 1.26 Backend**, **Python 3.13 Service/Runtime**, **PostgreSQL 18 Database**, **Next.js 16 + TypeScript + Shadcn UI Frontend** và **Nginx Reverse Proxy**.

## Tech Stack

| Thành phần | Công nghệ | Chi tiết |
| :--- | :--- | :--- |
| **Backend** | Golang `1.26` | `net/http` (stdlib), `log/slog` JSON, `pgx/v5`, Embedded SQL Migrations, multi-stage `alpine:3.21` |
| **Python Service** | Python `3.13.18` (`python:3.13-slim`) | `/opt/venv`, FastAPI + Uvicorn, `psycopg[binary]`, `scripts/` CLI runner, `pytest` + `ruff` |
| **Database** | PostgreSQL `18` (`postgres:18-alpine`) | Persistent volume `postgres_data:/var/lib/postgresql`, `pg_isready` healthcheck, nội bộ mạng Docker |
| **Frontend** | Next.js `16.4.6` + React `19.3` + TypeScript `5` | App Router, Tailwind CSS `4.3`, Shadcn UI `4.23` (`base-nova`), Standalone output |
| **Reverse Proxy** | Nginx (`nginx:alpine`) | Single entrypoint, Gzip, rate limiting, security headers, WebSocket HMR |
| **Orchestration** | Docker Compose | Production-ready (`docker-compose.yml`) + Dev hot-reload (`docker-compose.dev.yml`) |

## Quick Start (Chỉ cần Docker)

1. Kiểm tra file `.env` (hoặc tạo từ mẫu):
   ```bash
   cp -n .env.example .env
   ```
2. Build và khởi chạy toàn bộ stack (chờ tất cả container `healthy`):
   ```bash
   docker compose up -d --build --wait
   ```
3. Mở trình duyệt tại: **http://localhost:8000** (cổng mặc định `NGINX_PORT=8000`).

### Chạy chế độ Development (Hot Reload)
```bash
make dev
# hoặc:
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

## Các Endpoints mặc định (qua Nginx `:8000`)

| Method | Path | Mô tả |
| :--- | :--- | :--- |
| `GET` | `/` | Giao diện Next.js (TikTokUp System Telemetry Dashboard) |
| `GET` | `/nginx-health` | Healthcheck của Nginx Reverse Proxy |
| `GET` | `/health` | Direct liveness probe của Backend |
| `GET` | `/api/v1/health` | Trạng thái sức khỏe và uptime của Backend |
| `GET` | `/api/v1/ready` | Readiness probe của Backend (HTTP server, PostgreSQL & Python service) |
| `GET` | `/api/v1/info` | Metadata runtime Go 1.26, PostgreSQL 18 và Python 3.13 |

## Các lệnh Makefile hữu ích (100% chạy trong Docker)

```bash
make up                  # Build và chạy stack production-like (5 containers)
make down                # Dừng stack
make dev                 # Chạy stack dev có hot-reload (Air + FastAPI + Next.js + Postgres)
make test                # Chạy Go vet/tests + Python ruff/pytest + TypeScript + ESLint
make test-python         # Chạy Ruff lint + Pytest cho service Python trong Docker
make python-run SCRIPT=scripts/example.py  # Chạy file script Python trong container
make python-shell        # Mở bash shell tương tác trong container Python
make db-status           # Kiểm tra trạng thái PostgreSQL và danh sách migrations đã chạy
make db-migrate          # Chạy SQL migrations thông qua binary trong container backend
make db-shell            # Mở psql shell tương tác vào PostgreSQL container
make logs                # Xem log của toàn bộ container
make ps                  # Xem trạng thái và health của container
make shadcn-add COMPONENT=table  # Cài thêm Shadcn UI component qua Docker
```

## Tài liệu dự án (`docs/`)

- [Project Overview & PDR](docs/project-overview-pdr.md)
- [Business Requirements (BRD)](docs/business-requirements.md)
- [System Architecture](docs/system-architecture.md)
- [Codebase Summary](docs/codebase-summary.md)
- [Design Guidelines](docs/design-guidelines.md)
- [Deployment Guide](docs/deployment-guide.md)
- [Project Roadmap](docs/project-roadmap.md)
