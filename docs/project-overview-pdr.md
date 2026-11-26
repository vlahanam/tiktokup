# Project Overview & Product Development Requirements (PDR)

## 1. Tổng quan dự án

**TikTokUp** là một nền tảng ứng dụng web full-stack hiện đại được xây dựng theo kiến trúc container hóa 100% bằng Docker. Giai đoạn hiện tại thiết lập **Base Platform (v0.1.0)** nhằm cung cấp một nền móng kỹ thuật sạch, chuẩn hóa, bảo mật và dễ mở rộng cho toàn bộ các tính năng nghiệp vụ về sau.

### Mục tiêu cốt lõi của Base Platform
- **Zero Host Toolchain Dependency**: Toàn bộ quá trình cài đặt thư viện, phát triển, kiểm thử, build và triển khai đều chạy bên trong Docker container — lập trình viên không cần cài đặt Go, Node.js hay Nginx trực tiếp trên máy host.
- **Kiến trúc tách lớp rõ ràng**: Tách biệt rành mạch Edge Proxy (Nginx), Presentation/SSR Layer (Next.js) và API/Business Layer (Go 1.26).
- **Production-ready ngay từ đầu**: Multi-stage Docker builds, container chạy bằng non-root user, healthcheck tự động, structured JSON logging, graceful shutdown và request tracing (`X-Request-ID`) xuyên suốt hệ thống.

---

## 2. Công nghệ sử dụng (Tech Stack)

| Tầng (Layer) | Công nghệ | Phiên bản / Image | Lý do lựa chọn |
| :--- | :--- | :--- | :--- |
| **Edge / Reverse Proxy** | Nginx | `nginx:alpine` | Nhẹ, hiệu năng cao, gom toàn bộ traffic vào 1 entrypoint duy nhất, hỗ trợ Gzip, rate limiting, security headers và WebSocket/HMR. |
| **Backend API** | Golang | `golang:1.26-alpine` (build) / `alpine:3.21` (runtime) | Hiệu năng cao, tiêu thụ RAM thấp, binary tĩnh gọn nhẹ, standard library `net/http`, `log/slog` và driver `pgx/v5` hiệu suất cao. |
| **Python Runtime / Service** | Python + FastAPI | `python:3.13-slim` (Python `3.13.18`) | Sẵn sàng cho xử lý media, AI/ML, automation và script nền; đóng gói bằng `/opt/venv`, FastAPI + Uvicorn + `psycopg[binary]`. |
| **Database** | PostgreSQL | `postgres:18-alpine` (PostgreSQL `18.6`) | Cơ sở dữ liệu quan hệ mạnh mẽ, hỗ trợ JSONB, ACID chuẩn mực, lưu trữ bền vững bằng named volume và hệ thống migration tích hợp sẵn. |
| **Frontend Web** | Next.js + React + TypeScript | Next.js `16.4.6`, React `19.3.0`, TypeScript `5`, `node:22-alpine` | App Router + React Server Components (RSC), build `standalone` tối ưu cho Docker, kiểm soát kiểu dữ liệu chặt chẽ. |
| **UI System** | Shadcn UI + Tailwind CSS | Shadcn UI `4.23.1` (`base-nova`), Tailwind CSS `4.3.3`, `@base-ui/react` | Mã nguồn component thuộc sở hữu dự án, dễ tùy biến theo design system, hỗ trợ Dark/Light mode và chuẩn accessibility. |
| **Containerization** | Docker & Docker Compose | Docker Compose v2+ | Đảm bảo môi trường đồng nhất giữa local dev, staging và production. |

---

## 3. Phạm vi sản phẩm (Scope)

### 3.1. Trong phạm vi hiện tại (Phase 0 — Đã triển khai)
- Hạ tầng Docker Compose cho cả **Production-like** (`docker-compose.yml`) và **Development hot-reload** (`docker-compose.dev.yml`).
- Backend Go 1.26 với cấu trúc Clean/Layered (`cmd/`, `internal/`, `migrations/`, `pkg/`), kết nối **PostgreSQL 18** bằng `pgx/v5` connection pool, auto-migration giao dịch (`schema_migrations`), middleware chuẩn (`Recover`, `RequestID`, `Logger`, `CORS`), API response envelope thống nhất và hệ thống endpoint kiểm tra sức khỏe (`/health`, `/api/v1/health`, `/api/v1/ready`, `/api/v1/info`).
- Service **Python 3.13** (`python:3.13-slim`, `/opt/venv`) chạy nội bộ tại `python:8000` (FastAPI + Uvicorn + `psycopg[binary]`), hỗ trợ cả gọi API nội bộ từ Go Backend lẫn chạy CLI/scripts qua Docker (`make python-run`), kèm kiểm thử `ruff` và `pytest`.
- Cơ sở dữ liệu **PostgreSQL 18** (`postgres:18-alpine`) chạy cô lập trong Docker network `tiktokup-net` với persistent volume `postgres_data:/var/lib/postgresql` và healthcheck `pg_isready`.
- Frontend Next.js 16 + TypeScript + Tailwind CSS v4 + Shadcn UI (12 component nền tảng), hỗ trợ Light/Dark/System theme và trang **System Telemetry Dashboard** kiểm chứng kết nối live từ trình duyệt → Nginx → Go 1.26 Backend → PostgreSQL 18 & Python 3.13.
- Nginx Reverse Proxy định tuyến `/api/*` → Backend, `/*` → Frontend, tối ưu cache tài sản tĩnh `/_next/static/*`, hỗ trợ WebSocket HMR và rate limiting cơ bản.
- Bộ tài liệu kỹ thuật và sản phẩm hoàn chỉnh trong thư mục `docs/`.

### 3.2. Ngoài phạm vi hiện tại (Sẽ bổ sung ở các Phase tiếp theo)
- Cache/Queue phân tán (Redis) và Object Storage (MinIO/S3).
- Xác thực & phân quyền người dùng (Auth, JWT/Session, RBAC).
- Các module nghiệp vụ chuyên sâu của TikTokUp (quản lý kênh, upload/lập lịch nội dung, phân tích dữ liệu).
- Hệ thống CI/CD tự động, TLS/HTTPS certificate automation và stack giám sát tập trung (Prometheus/Grafana/OpenTelemetry).

---

## 4. Yêu cầu phi chức năng (Non-Functional Requirements — NFRs)

| Nhóm | Yêu cầu cụ thể | Cơ chế đáp ứng hiện tại |
| :--- | :--- | :--- |
| **Hiệu năng (Performance)** | API health/info phản hồi < 50ms ở local; trang chủ SSR phản hồi nhanh, kích thước image nhỏ gọn. | Go static binary (`-ldflags="-s -w"`), Next.js `output: "standalone"`, Nginx Gzip & upstream `keepalive 32`. |
| **Bảo mật (Security)** | Hạn chế lộ thông tin server, không chạy container dưới quyền root, kiểm soát CORS và rate limit. | `server_tokens off`, security headers trên Nginx, `limit_req` cho `/api/`, non-root user (`appuser` / `nextjs`), cô lập backend & frontend trong Docker bridge network nội bộ. |
| **Độ tin cậy (Reliability)** | Tự động khởi động lại khi lỗi, tắt ứng dụng an toàn không rớt request đang xử lý. | Docker `restart: unless-stopped`, `HEALTHCHECK` cho cả 3 service, `depends_on: condition: service_healthy`, Go `http.Server.Shutdown` theo `SIGINT`/`SIGTERM`. |
| **Khả năng quan sát (Observability)** | Truy vết được request từ edge đến backend dưới dạng JSON có cấu trúc. | Nginx tạo/truyền `X-Request-ID`, access log dạng JSON (`main_json`); Backend Go log bằng `log/slog` JSON kèm `request_id` và `duration_ms`. |
| **Khả năng bảo trì (Maintainability)** | Cấu trúc code chuẩn mực, có unit test, typecheck và linter chạy trong Docker. | `Makefile` tích hợp `go vet`, `go test`, `tsc --noEmit`, `eslint` chạy hoàn toàn qua Docker. |

---

## 5. Tiêu chí nghiệm thu (Acceptance Criteria — Base v0.1.0)

1. `docker compose up -d --build --wait` khởi động thành công cả 5 container (`tiktokup-postgres`, `tiktokup-python`, `tiktokup-backend`, `tiktokup-frontend`, `tiktokup-nginx`) với trạng thái `healthy`.
2. Truy cập `http://localhost:${NGINX_PORT}` trả về giao diện TikTokUp Base Platform (HTTP 200) và hiển thị trực tiếp thông số runtime từ Go 1.26 backend cùng trạng thái kết nối PostgreSQL 18 và Python 3.13.
3. Các endpoint `/nginx-health`, `/health`, `/api/v1/health`, `/api/v1/ready` (`checks.postgres = "ok"`, `checks.python = "ok"`), `/api/v1/info` trả về JSON chuẩn theo envelope và header `X-Request-ID`.
4. Endpoint không tồn tại (`/api/v1/unknown`) trả về HTTP 404 với JSON error envelope nhất quán.
5. Toàn bộ unit tests backend (`go vet`, `go test`) và kiểm tra tĩnh frontend (`tsc --noEmit`, `eslint`) đều pass 100% thông qua Docker.

---

## 6. Bảng thuật ngữ (Glossary)

- **RSC (React Server Components)**: Component render phía server trong Next.js App Router giúp giảm JavaScript gửi xuống trình duyệt.
- **Standalone Output**: Chế độ đóng gói của Next.js chỉ xuất các file tối thiểu cần thiết (`server.js`) để giảm kích thước Docker image.
- **Reverse Proxy**: Máy chủ trung gian (Nginx) tiếp nhận request từ client và điều phối tới các service nội bộ (`frontend`, `backend`).
- **Liveness / Readiness Probe**: Cơ chế kiểm tra container còn sống (`/health`) và đã sẵn sàng nhận traffic (`/api/v1/ready`).
- **Correlation ID (`X-Request-ID`)**: Mã định danh duy nhất cho mỗi HTTP request giúp liên kết log giữa Nginx và Backend.
