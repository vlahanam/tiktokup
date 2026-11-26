# Deployment Guide

Hướng dẫn cài đặt, phát triển, build, triển khai và vận hành dự án **TikTokUp** hoàn toàn thông qua **Docker & Docker Compose**.

---

## 1. Yêu cầu hệ thống (Prerequisites)

- **Docker Engine** `>= 24.x` (đã kiểm thử trên Docker `29.5.2`)
- **Docker Compose Plugin** `v2+` (đã kiểm thử trên `v5.1.4`)
- **GNU Make** *(tùy chọn — dùng cho các lệnh rút gọn trong `Makefile`)*
- **Không cần cài đặt** Go, Node.js, npm hay Nginx trên máy host.

Kiểm tra nhanh môi trường:
```bash
docker --version
docker compose version
```

---

## 2. Cấu hình môi trường (Environment Setup)

1. Tạo file `.env` từ mẫu `.env.example` (nếu chưa có):
   ```bash
   cp .env.example .env
   ```
2. Kiểm tra và điều chỉnh các biến quan trọng trong `.env`:
   ```env
   APP_ENV=development
   LOG_LEVEL=info

   # Cổng public của Nginx trên host (mặc định 8000 để tránh trùng cổng 80)
   NGINX_PORT=8000

   # Chỉ dùng khi chạy docker-compose.dev.yml
   BACKEND_PORT=8081
   FRONTEND_PORT=3000
   POSTGRES_PORT=5433
   PYTHON_PORT=8002

   # Python 3.13 Service
   PYTHON_APP_NAME=tiktokup-python
   PYTHON_SERVICE_URL=http://python:8000

   # PostgreSQL 18
   POSTGRES_USER=tiktokup
   POSTGRES_PASSWORD=tiktokup_secret
   POSTGRES_DB=tiktokup
   POSTGRES_SSLMODE=disable
   DB_MAX_CONNS=10
   DB_MIN_CONNS=2
   DB_AUTO_MIGRATE=true

   # Danh sách origin hợp lệ cho CORS
   CORS_ALLOWED_ORIGINS=http://localhost:8000,http://localhost:3000

   # Frontend API URLs
   NEXT_PUBLIC_APP_NAME=TikTokUp
   NEXT_PUBLIC_API_URL=/api
   INTERNAL_API_URL=http://backend:8080
   ```

> [!IMPORTANT]
> Biến bắt đầu bằng `NEXT_PUBLIC_` được Next.js nhúng vào bundle tại thời điểm **build Docker image** (truyền qua `build.args` trong `docker-compose.yml`). Khi thay đổi các biến này, hãy chạy lại `docker compose build frontend`.

---

## 3. Chế độ Production / Staging (Mặc định)

Sử dụng `docker-compose.yml` với multi-stage images tối ưu (`postgres:18-alpine` với persistent volume `postgres_data:/var/lib/postgresql`, `python:3.13-slim` với `/opt/venv`, `alpine:3.21` cho Backend, `node:22-alpine` standalone cho Frontend, `nginx:alpine` cho Edge Proxy).

### Khởi chạy & chờ toàn bộ service healthy
```bash
docker compose up -d --build --wait
# hoặc
make up
```

### Kiểm tra trạng thái container
```bash
docker compose ps
```
Kết quả mong đợi: cả 5 container `tiktokup-postgres`, `tiktokup-python`, `tiktokup-backend`, `tiktokup-frontend`, `tiktokup-nginx` đều ở trạng thái `Up (...) (healthy)`.

### Kiểm tra kết nối End-to-End
```bash
# Kiểm tra Nginx edge health
curl -i http://localhost:8000/nginx-health

# Kiểm tra Go 1.26 Backend qua Nginx
curl -s http://localhost:8000/api/v1/health
curl -s http://localhost:8000/api/v1/ready
curl -s http://localhost:8000/api/v1/info

# Kiểm tra Frontend SSR qua Nginx
curl -I http://localhost:8000/
```

Mở trình duyệt tại: `http://localhost:8000` (hoặc cổng `NGINX_PORT` bạn đã cấu hình).

### Dừng hệ thống
```bash
docker compose down
# hoặc
make down
```

---

## 4. Chế độ Development (Hot Reload qua Docker)

File `docker-compose.dev.yml` ghi đè `docker-compose.yml` để bật:
- **PostgreSQL 18**: mở cổng `${POSTGRES_PORT:-5433}:5432` để kết nối trực tiếp từ các công cụ như DBeaver, TablePlus, DataGrip hoặc `psql` mà không xung đột cổng `5432` trên host.
- **Python 3.13**: target `dev` chạy `uvicorn --reload`, mount `./python:/app`, mở cổng `${PYTHON_PORT:-8002}:8000` (có thể mở Swagger UI tại `http://localhost:8002/docs`).
- **Backend**: target `dev` sử dụng **Air** (`.air.toml`) tự động rebuild Go 1.26 khi lưu file, mở cổng `${BACKEND_PORT:-8081}`.
- **Frontend**: target `dev` chạy `next dev` với `WATCHPACK_POLLING=true` và WebSocket HMR đi qua Nginx, mở cổng `${FRONTEND_PORT:-3000}`.
- **Nginx**: mount trực tiếp `./nginx/nginx.conf` và `./nginx/conf.d` (read-only).

Khởi chạy dev stack:
```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build
# hoặc
make dev
```

Dừng dev stack:
```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml down
# hoặc
make dev-down
```

---

## 5. Kiểm thử & Bảo trì hoàn toàn bằng Docker

Tất cả lệnh kiểm thử và quản lý package đều chạy trong container tạm thời (`--rm`):

```bash
# Chạy toàn bộ: Go vet/tests + Python ruff/pytest + TypeScript typecheck + ESLint
make test

# Hoặc chạy riêng lẻ:
make test-backend
make test-python
make typecheck-frontend
make lint-frontend

# Làm việc với Python trong Docker (không cần cài Python trên host)
make python-run SCRIPT=scripts/example.py   # Chạy script mẫu kiểm tra Python + Postgres
make python-shell                           # Mở bash shell vào container tiktokup-python
# Thêm thư viện Python mới: bổ sung vào python/pyproject.toml rồi chạy:
docker compose up -d --build python

# Thêm Shadcn UI component mới (ví dụ: table)
make shadcn-add COMPONENT=table

# Thao tác PostgreSQL & SQL Migrations qua Docker
make db-status     # Kiểm tra pg_isready & danh sách schema_migrations đã áp dụng
make db-migrate    # Chạy các SQL migrations còn pending
make db-shell      # Mở psql shell tương tác trong container tiktokup-postgres

# Backup & Restore nhanh bằng pg_dump/psql qua Docker
docker compose exec -T postgres sh -c 'pg_dump -U "$POSTGRES_USER" "$POSTGRES_DB"' > backup.sql
docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"' < backup.sql

# Kiểm tra cú pháp cấu hình Nginx đang chạy
docker compose exec -T nginx nginx -t
```

---

## 6. Triển khai Production & Cấu hình HTTPS / TLS

### 6.1. Điều chỉnh biến môi trường cho Production
Trong `.env` trên server production:
- Đặt `APP_ENV=production`, `LOG_LEVEL=info`.
- Đặt `NGINX_PORT=80` (nếu Nginx trực tiếp nhận traffic public) hoặc cổng nội bộ (nếu đặt sau Cloud Load Balancer / Traefik / Cloudflare Tunnel).
- Cập nhật `CORS_ALLOWED_ORIGINS=https://your-domain.com`.

### 6.2. Bật HTTPS / TLS
Có 2 mô hình triển khai khuyến nghị:
1. **TLS Termination ở lớp ngoài (Khuyến nghị cho Cloud/VPS có Reverse Proxy chung)**: Đặt Cloudflare, AWS ALB hoặc Traefik/Caddy phía trước trỏ vào `http://<host>:${NGINX_PORT}` và giữ nguyên Nginx của dự án xử lý định tuyến nội bộ.
2. **TLS trực tiếp trên Nginx container**:
   - Mở thêm cổng `443:443` trong `docker-compose.yml` cho service `nginx`.
   - Mount thư mục chứng chỉ SSL (ví dụ Let's Encrypt `/etc/letsencrypt:/etc/letsencrypt:ro`).
   - Thêm server block `listen 443 ssl http2;` kèm redirect HTTP `80 → 443` và header `Strict-Transport-Security` trong `nginx/conf.d/default.conf`.

---

## 7. Giám sát Logs & Troubleshooting

### Xem logs cấu trúc JSON
```bash
# Toàn bộ stack
docker compose logs -f

# Riêng từng service
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f nginx
```

### Các lỗi thường gặp & Cách khắc phục

| Hiện tượng | Nguyên nhân phổ biến | Cách xử lý |
| :--- | :--- | :--- |
| `Bind for 0.0.0.0:8000 failed: port is already allocated` | Cổng `NGINX_PORT` đang bị dịch vụ khác chiếm | Đổi `NGINX_PORT` trong `.env` sang cổng trống (ví dụ `8085`) và chạy lại `docker compose up -d`. |
| Frontend báo `Backend API Unreachable` | Backend chưa sẵn sàng hoặc `NEXT_PUBLIC_API_URL` sai | Kiểm tra `docker compose ps` và `docker compose logs backend`; đảm bảo `NEXT_PUBLIC_API_URL=/api`. |
| Thay đổi `NEXT_PUBLIC_*` không có tác dụng | Giá trị được bake lúc build Next.js | Chạy `docker compose build --no-cache frontend && docker compose up -d`. |
| Quyền file trong `frontend/` bị thuộc `root` khi cài npm qua Docker | Chạy container không truyền UID/GID hoặc không `chown` | Dùng lệnh `make shadcn-add` hoặc truyền `-u "$(id -u):$(id -g)"` / `chown -R` sau khi chạy tool trong container. |
| Thay đổi `nginx/conf.d/default.conf` chưa nhận | Đang chạy production image (COPY file vào image) | Chạy `docker compose up -d --build nginx` (hoặc dùng dev mode có bind-mount và `docker compose exec nginx nginx -s reload`). |

---

## 8. Production Release & Rollback Checklist

- [ ] Toàn bộ `make test` (Go vet, Go tests, TypeScript typecheck, ESLint) pass 100%.
- [ ] Đã cập nhật `.env` production (`APP_ENV=production`, `CORS_ALLOWED_ORIGINS` đúng domain thật, không dùng `*`).
- [ ] Đã kiểm tra `docker compose config` hợp lệ và `nginx -t` không có cảnh báo.
- [ ] Đã xác minh toàn bộ container đạt trạng thái `healthy` sau khi `up -d --build --wait`.
- [ ] **Rollback nhanh**: Gắn tag version image trước khi deploy (`docker tag tiktokup-backend:latest tiktokup-backend:v0.1.0`), nếu phát sinh lỗi chỉ cần trỏ lại tag cũ và `docker compose up -d`.
