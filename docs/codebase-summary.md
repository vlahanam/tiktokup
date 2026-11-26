# Codebase Summary

Tài liệu này mô tả cấu trúc mã nguồn thực tế, trách nhiệm của từng module, danh mục API và các biến môi trường trong repository **TikTokUp**.

---

## 1. Cấu trúc thư mục tổng thể

```text
tiktokup/
├── .env                          # Biến môi trường local (không commit lên Git)
├── .env.example                  # Mẫu cấu hình biến môi trường
├── .gitignore                    # Quy tắc loại trừ file cho Go, Next.js, Docker
├── AGENTS.md                     # Hướng dẫn làm việc cho AI Coding Agent
├── Makefile                      # Các lệnh tắt chạy 100% qua Docker
├── README.md                     # Tài liệu khởi chạy nhanh
├── docker-compose.yml            # Stack production-ready (postgres, backend, frontend, nginx)
├── docker-compose.dev.yml        # Override cho dev hot-reload (Air + Next.js dev + Postgres :5433)
├── backend/                      # Golang 1.26 API Service
│   ├── .air.toml                 # Cấu hình hot-reload Air cho môi trường dev
│   ├── .dockerignore
│   ├── Dockerfile                # Multi-stage: base, dev, builder (server+migrate), production
│   ├── go.mod                    # module tiktokup/backend, go 1.26, pgx/v5 v5.13.0
│   ├── go.sum                    # Checksums dependencies (sinh trong Docker)
│   ├── cmd/
│   │   ├── migrate/
│   │   │   └── main.go           # CLI binary chạy SQL migrations độc lập
│   │   └── server/
│   │       └── main.go           # Entrypoint kết nối Postgres, auto-migrate, khởi chạy HTTP server
│   ├── internal/
│   │   ├── config/
│   │   │   ├── config.go         # Nạp & validate cấu hình app + DatabaseConfig/DSN
│   │   │   └── config_test.go    # Unit tests cho config & DSN
│   │   ├── database/
│   │   │   ├── migrate.go        # Transactional SQL migration runner (schema_migrations)
│   │   │   ├── migrate_test.go   # Unit tests kiểm tra embedded SQL migrations
│   │   │   └── postgres.go       # pgxpool.Pool connect (có retry), Ping, ServerVersion
│   │   ├── handler/
│   │   │   └── health.go         # SystemHandler: Health, Ready (ping Postgres), Info (DB metadata)
│   │   ├── middleware/
│   │   │   ├── cors.go           # CORS & OPTIONS preflight
│   │   │   ├── logger.go         # Structured slog access log
│   │   │   ├── middleware_test.go # Unit tests cho RequestID, CORS, Recover
│   │   │   ├── recover.go        # Bắt panic và trả 500 JSON
│   │   │   └── request_id.go     # Tạo/truyền header & context X-Request-ID
│   │   ├── router/
│   │   │   ├── router.go         # Đăng ký route & nối chuỗi middleware
│   │   │   └── router_test.go    # Tests endpoints (healthy DB + 503 khi DB lỗi)
│   │   └── server/
│   │       └── server.go         # Cấu hình http.Server & graceful shutdown
│   ├── migrations/
│   │   ├── 000001_init.sql       # Migration khởi tạo bảng baseline app_settings (JSONB)
│   │   └── migrations.go         # //go:embed *.sql
│   └── pkg/
│       └── response/
│           ├── json.go           # Envelope chuẩn (OK, Created, Error, RequestID context)
│           └── json_test.go      # Unit tests cho response helpers
├── frontend/                     # Next.js 16 + TypeScript + Shadcn UI
│   ├── .dockerignore
│   ├── Dockerfile                # Multi-stage: base, deps, dev, builder, production
│   ├── components.json           # Cấu hình Shadcn UI (style: base-nova, lucide)
│   ├── eslint.config.mjs         # Cấu hình ESLint 9 + eslint-config-next
│   ├── next.config.ts            # output: "standalone", cacheComponents, turbopack
│   ├── package.json              # Scripts & dependencies (tiktokup-frontend)
│   ├── package-lock.json         # Lockfile npm được sinh trong Docker
│   ├── postcss.config.mjs        # Cấu hình PostCSS
│   ├── tsconfig.json             # Cấu hình TypeScript strict + alias @/*
│   ├── public/                   # Static assets
│   └── src/
│       ├── app/
│       │   ├── error.tsx         # Global client error boundary
│       │   ├── globals.css       # Tailwind CSS v4 + Shadcn OKLCH theme tokens
│       │   ├── layout.tsx        # RootLayout (Geist fonts, ThemeProvider, Header, Footer, Toaster)
│       │   ├── loading.tsx       # Suspense skeleton fallback
│       │   ├── not-found.tsx     # Trang 404 tùy chỉnh
│       │   └── page.tsx          # Trang chủ TikTokUp Base Platform
│       ├── components/
│       │   ├── layout/
│       │   │   ├── footer.tsx    # Footer hiển thị stack info
│       │   │   └── header.tsx    # Sticky header + brand + ThemeToggle
│       │   ├── ui/               # 12 Shadcn UI components (Base UI)
│       │   │   ├── alert.tsx
│       │   │   ├── badge.tsx
│       │   │   ├── button.tsx
│       │   │   ├── card.tsx
│       │   │   ├── dialog.tsx
│       │   │   ├── dropdown-menu.tsx
│       │   │   ├── input.tsx
│       │   │   ├── label.tsx
│       │   │   ├── separator.tsx
│       │   │   ├── skeleton.tsx
│       │   │   ├── sonner.tsx
│       │   │   └── tabs.tsx
│       │   ├── system-status-dashboard.tsx # Live telemetry dashboard gọi Backend qua Nginx
│       │   ├── theme-provider.tsx # Wrapper next-themes ThemeProvider
│       │   └── theme-toggle.tsx  # Menu chuyển đổi Light / Dark / System
│       ├── lib/
│       │   ├── api.ts            # Typed API fetch client (dual browser/SSR base URL)
│       │   └── utils.ts          # Helper cn() của Shadcn UI
│       └── types/
│           └── api.ts            # Interface khớp với Go response Envelope & payloads
├── nginx/                        # Reverse Proxy Service
│   ├── Dockerfile                # nginx:alpine + HEALTHCHECK
│   ├── nginx.conf                # Main config: worker, gzip, JSON logs, rate limit
│   └── conf.d/
│       └── default.conf          # Upstreams & location routing rules
└── docs/                         # Tài liệu dự án
    ├── business-requirements.md
    ├── codebase-summary.md
    ├── deployment-guide.md
    ├── design-guidelines.md
    ├── project-overview-pdr.md
    ├── project-roadmap.md
    └── system-architecture.md
```

---

## 2. Chi tiết Backend (`backend/`)

- **Module**: `tiktokup/backend` (`go 1.26`, build & kiểm thử trên `go1.26.11`).
- **Dependencies**: Thư viện chuẩn Go 1.26 (`net/http`, `log/slog`, `encoding/json`, `crypto/rand`, `embed`, `slices`) và `github.com/jackc/pgx/v5 v5.13.0` (`pgxpool`) cho PostgreSQL 18.
- **Tách trách nhiệm**:
  - `cmd/server`: Wiring khởi tạo config, kết nối PostgreSQL, chạy auto-migration (nếu bật), lắng nghe `SIGINT`/`SIGTERM` và đóng pool khi shutdown.
  - `cmd/migrate`: Binary CLI chạy SQL migrations độc lập (`/app/migrate` trong production container).
  - `internal/config`: Tập trung toàn bộ biến môi trường vào `Config` & `DatabaseConfig` (hỗ trợ `DATABASE_URL` hoặc từng biến `DB_*`).
  - `internal/database`: Quản lý `pgxpool.Pool` (có retry lúc boot, `Ping`, `ServerVersion`) và transactional SQL migrator (`schema_migrations`).
  - `migrations/`: Chứa file SQL có đánh số thứ tự (`000001_init.sql`) được nhúng trực tiếp vào binary bằng `embed.FS`.
  - `internal/middleware`: Xử lý cross-cutting concerns độc lập (`Recover`, `RequestID`, `Logger`, `CORS`).
  - `internal/handler`: Nhận HTTP request, kiểm tra DB qua interface `DatabaseInspector`, trả kết quả qua `pkg/response`.
  - `internal/router`: Tập trung khai báo route bằng Go `ServeMux` pattern (`GET /api/v1/...`).
  - `internal/server`: Đóng gói vòng đời `http.Server` (timeouts và graceful shutdown).
  - `pkg/response`: Định nghĩa contract JSON dùng chung (`Envelope`, `ErrorInfo`, `Meta`).

---

## 3. Chi tiết Frontend (`frontend/`)

### Dependencies chính (`frontend/package.json`)
| Package | Phiên bản | Vai trò |
| :--- | :--- | :--- |
| `next` | `16.4.6` | Full-stack React framework (App Router, RSC, Standalone output) |
| `react` / `react-dom` | `19.3.0` | UI runtime |
| `typescript` | `^5` | Kiểm tra kiểu tĩnh |
| `tailwindcss` / `@tailwindcss/turbopack` | `4.3.3` | Utility-first CSS engine tích hợp Turbopack |
| `shadcn` | `^4.23.1` | Hệ thống UI component (`style: base-nova`) |
| `@base-ui/react` | `^1.9.0` | Headless accessible UI primitives của Shadcn v4 |
| `lucide-react` | `^1.81.0` | Bộ icon vector |
| `next-themes` | `^0.4.6` | Quản lý Light / Dark / System theme không flash |
| `sonner` | `^2.0.8` | Toast notifications |
| `class-variance-authority` & `cn` | `^0.7.1` / `^0.5.2` | Quản lý variants và gộp class Tailwind |

### Scripts (`frontend/package.json`)
- `npm run dev`: Chạy Next.js dev server.
- `npm run build`: Build production bundle (`.next/standalone`).
- `npm run start`: Khởi chạy Next.js production server.
- `npm run lint`: Kiểm tra mã nguồn bằng ESLint 9 (`eslint-config-next`).
- `npm run typecheck`: Kiểm tra kiểu TypeScript (`tsc --noEmit`).

---

## 4. Danh mục API & Hợp đồng dữ liệu (API Contract)

### 4.1. Endpoints hiện có

| Method | Path (qua Nginx) | Đích xử lý | Mô tả |
| :--- | :--- | :--- | :--- |
| `GET` | `/nginx-health` | Nginx | Kiểm tra Nginx edge proxy còn sống |
| `GET` | `/health` | Backend `:8080` | Liveness probe cho Docker healthcheck |
| `GET` | `/api/v1/health` | Backend `:8080` | Trạng thái hoạt động và uptime của API |
| `GET` | `/api/v1/ready` | Backend + Postgres | Readiness probe kiểm tra `http_server` & live ping `postgres` (trả `503` nếu DB lỗi) |
| `GET` | `/api/v1/info` | Backend + Postgres | Thông tin phiên bản Go 1.26, OS/Arch, CPU và metadata PostgreSQL 18 |

### 4.2. JSON Response Envelope

Phản hồi thành công (`200 OK`):
```json
{
  "success": true,
  "data": {
    "service": "tiktokup-backend",
    "version": "0.1.0",
    "environment": "development",
    "go_version": "go1.26.11",
    "os": "linux",
    "arch": "amd64",
    "num_cpu": 12,
    "started_at": "2026-11-26T09:03:00Z",
    "uptime_seconds": 42,
    "database": {
      "engine": "postgresql",
      "connected": true,
      "version": "PostgreSQL 18.6"
    }
  },
  "meta": {
    "request_id": "a1b2c3d4e5f6...",
    "timestamp": "2026-11-26T09:03:42Z"
  }
}
```

Phản hồi lỗi (`404 Not Found` / `500 Internal Server Error`):
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "endpoint not found"
  },
  "meta": {
    "request_id": "a1b2c3d4e5f6...",
    "timestamp": "2026-11-26T08:50:42Z"
  }
}
```

---

## 5. Bảng tham chiếu Environment Variables (`.env.example`)

| Biến | Mặc định | Dịch vụ sử dụng | Ý nghĩa |
| :--- | :--- | :--- | :--- |
| `PROJECT_NAME` | `tiktokup` | Compose | Tên định danh dự án |
| `APP_NAME` | `tiktokup-backend` | Backend | Tên dịch vụ hiển thị trong log & `/api/v1/info` |
| `APP_VERSION` | `0.1.0` | Backend | Phiên bản ứng dụng |
| `APP_ENV` | `development` | Backend | Môi trường chạy (`development` / `staging` / `production`) |
| `LOG_LEVEL` | `info` | Backend | Mức log `slog` (`debug`, `info`, `warn`, `error`) |
| `NGINX_PORT` | `8000` | Nginx | Cổng public trên máy host ánh xạ vào `:80` của Nginx |
| `BACKEND_PORT` | `8081` | Backend (dev) | Cổng host ánh xạ tới `:8080` khi chạy `docker-compose.dev.yml` |
| `FRONTEND_PORT` | `3000` | Frontend (dev) | Cổng host ánh xạ tới `:3000` khi chạy `docker-compose.dev.yml` |
| `POSTGRES_PORT` | `5433` | PostgreSQL (dev) | Cổng host ánh xạ tới `:5432` khi chạy `docker-compose.dev.yml` |
| `POSTGRES_USER` | `tiktokup` | Postgres & Backend | Tên user quản trị cơ sở dữ liệu |
| `POSTGRES_PASSWORD` | `tiktokup_secret` | Postgres & Backend | Mật khẩu kết nối PostgreSQL (bắt buộc đổi trên production) |
| `POSTGRES_DB` | `tiktokup` | Postgres & Backend | Tên database chính |
| `POSTGRES_SSLMODE` | `disable` | Backend | Chế độ SSL kết nối Postgres (`disable` nội bộ / `require` cho managed DB) |
| `DB_MAX_CONNS` | `10` | Backend | Số lượng kết nối tối đa trong `pgxpool` |
| `DB_MIN_CONNS` | `2` | Backend | Số lượng kết nối tối thiểu giữ ấm trong `pgxpool` |
| `DB_AUTO_MIGRATE` | `true` | Backend | Tự động chạy pending SQL migrations khi khởi động server |
| `DATABASE_URL` | *(tùy chọn)* | Backend | Chuỗi URI đầy đủ (`postgres://...`), nếu đặt sẽ ưu tiên hơn `DB_*` |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:8000,http://localhost:3000` | Backend | Danh sách origin được phép gọi API (phân tách bằng dấu phẩy) |
| `SERVER_READ_TIMEOUT_SECONDS` | `15` | Backend | Timeout đọc toàn bộ HTTP request |
| `SERVER_WRITE_TIMEOUT_SECONDS` | `30` | Backend | Timeout ghi HTTP response |
| `SERVER_IDLE_TIMEOUT_SECONDS` | `60` | Backend | Timeout giữ kết nối keep-alive |
| `SERVER_SHUTDOWN_TIMEOUT_SECONDS` | `10` | Backend | Thời gian chờ tối đa khi graceful shutdown |
| `NEXT_PUBLIC_APP_NAME` | `TikTokUp` | Frontend (build & runtime) | Tên ứng dụng hiển thị trên UI |
| `NEXT_PUBLIC_API_URL` | `/api` | Frontend (client) | Base path gọi API từ trình duyệt qua Nginx |
| `INTERNAL_API_URL` | `http://backend:8080` | Frontend (SSR) | URL gọi API nội bộ giữa container Frontend và Backend |

---

## 6. Quy ước mã nguồn (Coding Conventions)

- **Go (Backend)**:
  - Format chuẩn `gofmt`, kiểm tra tĩnh bằng `go vet ./...`.
  - Mọi handler đều trả response thông qua `pkg/response` (`OK`, `Created`, `Error`).
  - Luôn dùng structured logging qua `slog` kèm `request_id` từ context; không dùng `fmt.Println` hay `log.Fatal` bên trong handler.
- **TypeScript / Next.js (Frontend)**:
  - Bật TypeScript strict; định nghĩa type payload API tập trung tại `src/types/api.ts`.
  - Ưu tiên Server Components mặc định; chỉ gắn `"use client"` cho component cần state, event hoặc browser API.
  - Component Shadcn UI v4 dùng Base UI (`@base-ui/react`), sử dụng prop `render={<... />}` thay cho Radix `asChild` (kèm `nativeButton={false}` khi render `<Link>` bên trong `<Button>`).
