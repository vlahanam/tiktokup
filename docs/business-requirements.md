# Business Requirements Document (BRD)

## 1. Mục tiêu nghiệp vụ (Business Objectives)

Dự án **TikTokUp** hướng tới xây dựng một hệ thống web application hiện đại, ổn định và dễ mở rộng. Ở giai đoạn khởi tạo nền tảng (Base Platform v0.1.0), các mục tiêu nghiệp vụ ưu tiên bao gồm:

1. **Rút ngắn thời gian onboarding & phát triển tính năng**: Chuẩn hóa toàn bộ môi trường qua Docker để bất kỳ thành viên nào cũng khởi chạy được dự án chỉ với một lệnh duy nhất.
2. **Đảm bảo tính nhất quán trải nghiệm người dùng (UX)**: Thiết lập hệ thống giao diện thống nhất (Shadcn UI + Tailwind CSS v4 + Dark/Light theme) sẵn sàng ghép nối các màn hình nghiệp vụ.
3. **Chuẩn hóa giao tiếp giữa Frontend và Backend**: Mọi API tuân theo hợp đồng JSON thống nhất (Envelope + Error Code + Request ID) đi qua một domain/origin duy nhất nhờ Nginx Reverse Proxy.
4. **Sẵn sàng mở rộng quy mô**: Cho phép bổ sung nhanh cơ sở dữ liệu, hàng đợi xử lý bất đồng bộ (media/upload) và các module nghiệp vụ mà không phải đập bỏ kiến trúc cốt lõi.

---

## 2. Đối tượng sử dụng & Vai trò (Stakeholders & Personas)

| Vai trò | Mô tả | Nhu cầu chính |
| :--- | :--- | :--- |
| **Developer (Backend / Frontend)** | Kỹ sư phát triển tính năng trên Go 1.26 và Next.js | Môi trường Docker đồng nhất, hot-reload nhanh, cấu trúc code rõ ràng, kiểm thử tự động dễ chạy. |
| **DevOps / System Admin** | Người vận hành và triển khai hệ thống | Healthcheck chuẩn, log JSON dễ thu thập, cấu hình qua `.env`, image tối ưu và bảo mật non-root. |
| **Operator / Creator *(Giai đoạn tiếp theo)*** | Người dùng vận hành nội dung trên nền tảng TikTokUp | Giao diện trực quan, tốc độ tải nhanh, phản hồi lỗi rõ ràng, hỗ trợ upload nội dung ổn định. |
| **Administrator *(Giai đoạn tiếp theo)*** | Quản trị viên hệ thống | Quản lý người dùng, phân quyền (RBAC), theo dõi trạng thái dịch vụ và báo cáo vận hành. |

---

## 3. Yêu cầu chức năng (Functional Requirements)

### 3.1. Nhóm FR-0: Nền tảng cốt lõi (Base Platform) — *Đã triển khai*

| Mã | Tên chức năng | Mô tả yêu cầu | Trạng thái |
| :--- | :--- | :--- | :--- |
| **FR-001** | Single Entrypoint Routing | Mọi truy cập người dùng và API đều đi qua Nginx Reverse Proxy trên một cổng cấu hình qua `NGINX_PORT`. `/api/*` định tuyến về Backend, `/*` định tuyến về Frontend. | Hoàn thành |
| **FR-002** | Service Health & Readiness | Hệ thống cung cấp các endpoint kiểm tra sức khỏe cho Nginx (`/nginx-health`) và Backend (`/health`, `/api/v1/health`, `/api/v1/ready` có live ping PostgreSQL, trả 503 nếu DB mất kết nối). | Hoàn thành |
| **FR-003** | Runtime Diagnostics API | Backend cung cấp `GET /api/v1/info` trả về metadata dịch vụ (tên service, phiên bản, môi trường, phiên bản Go 1.26, OS/Arch, số CPU, uptime và phiên bản PostgreSQL 18). | Hoàn thành |
| **FR-004** | Standardized API Contract | Toàn bộ phản hồi API bọc trong JSON Envelope (`success`, `data`, `error`, `meta`), kèm `X-Request-ID` và xử lý 404/500 dạng JSON nhất quán. | Hoàn thành |
| **FR-005** | System Telemetry Dashboard | Trang chủ Frontend hiển thị trạng thái kết nối thời gian thực của Backend API và PostgreSQL 18, độ trễ (latency ms), metadata runtime, danh sách route hệ thống và nút làm mới có thông báo toast. | Hoàn thành |
| **FR-006** | Theme & Responsive Shell | Giao diện hỗ trợ chuyển đổi Light / Dark / System theme, responsive trên desktop & mobile, có sẵn trang Loading, Error Boundary và 404 Not Found. | Hoàn thành |
| **FR-007** | Dockerized Dev & Prod Modes | Hỗ trợ vận hành ở 2 chế độ qua Docker Compose: Production (build tối ưu) và Development (hot-reload cho Go, Next.js và mở cổng PostgreSQL dev `5433`). | Hoàn thành |
| **FR-008** | PostgreSQL Persistence & Migrations | Tích hợp PostgreSQL 18 (`postgres:18-alpine`) với volume bền vững, `pgx/v5` connection pool và hệ thống SQL migration giao dịch (`schema_migrations`). | Hoàn thành |

### 3.2. Nhóm FR-1+: Backlog nghiệp vụ mở rộng — *Dự kiến cho các Phase tiếp theo*

> [!NOTE]
> Các yêu cầu dưới đây đóng vai trò khung định hướng (Business Backlog Template) cho các giai đoạn sau khi đội sản phẩm chốt đặc tả chi tiết; **chưa được triển khai** trong Base v0.1.0.

| Mã | Nhóm tính năng dự kiến | Mô tả sơ bộ | Phase dự kiến |
| :--- | :--- | :--- | :--- |
| **FR-101** | Authentication & Authorization | Đăng nhập, đăng xuất, quản lý phiên/JWT và phân quyền RBAC (Admin / Operator / Viewer). | Phase 1 |
| **FR-102** | Domain Data Schemas & Repositories | Thiết kế bảng nghiệp vụ chi tiết trên PostgreSQL 18 (users, workspaces, audit logs) kèm lớp `internal/repository`. | Phase 1 |
| **FR-201** | Account / Channel Management | Quản lý danh sách kênh/tài khoản liên kết, trạng thái kết nối và cấu hình theo workspace. | Phase 2 |
| **FR-202** | Media Upload & Processing | Tải lên tệp đa phương tiện (Nginx đã cấu hình sẵn `client_max_body_size 50m`), lưu trữ trên Object Storage và xử lý bất đồng bộ qua queue. | Phase 2 |
| **FR-203** | Publishing & Scheduling | Lập lịch đăng tải, theo dõi trạng thái job (pending/processing/published/failed) và cơ chế retry. | Phase 2 |
| **FR-204** | Analytics & Reporting | Bảng tổng hợp hiệu suất, bộ lọc thời gian và xuất báo cáo vận hành. | Phase 2 |

---

## 4. Tiêu chí chấp nhận (Acceptance Criteria)

### AC-01: Kết nối End-to-End qua Nginx
- **Given** toàn bộ stack đang chạy qua Docker Compose,
- **When** người dùng mở trình duyệt tại `http://localhost:${NGINX_PORT}`,
- **Then** trang chủ tải thành công (HTTP 200), tự động gọi `/api/v1/health`, `/api/v1/ready`, `/api/v1/info` qua cùng origin và hiển thị trạng thái `Healthy` kèm phiên bản runtime `go1.26.x`.

### AC-02: Truy vết Request và Xử lý Lỗi
- **Given** client gửi request bất kỳ tới `/api/*`,
- **When** request thành công hoặc gặp lỗi 404/500,
- **Then** response luôn chứa header `X-Request-ID` khớp với `meta.request_id` trong JSON body và định dạng lỗi luôn chứa `error.code` và `error.message`.

### AC-03: Tính độc lập môi trường
- **Given** máy host chỉ cài Docker và Docker Compose (không cài Go/Node),
- **When** chạy `docker compose up -d --build` và `make test`,
- **Then** hệ thống build, khởi chạy và hoàn thành toàn bộ unit test, typecheck, lint mà không có lỗi.

---

## 5. Giả định, Ràng buộc và Rủi ro

- **Giả định**: Người vận hành có quyền chạy Docker daemon và truy cập mạng để tải base images (`golang:1.26-alpine`, `node:22-alpine`, `nginx:alpine`, `alpine:3.21`) trong lần build đầu tiên.
- **Ràng buộc kỹ thuật**:
  - Máy host có thể đang dùng cổng `80` cho dịch vụ khác, do đó cổng mặc định của Nginx trên host được thiết lập linh hoạt qua biến `NGINX_PORT` (mặc định `8000`).
  - Biến `NEXT_PUBLIC_*` của Next.js được nhúng ở bước build Docker image; khi thay đổi giá trị này cho môi trường mới cần build lại image frontend.
- **Rủi ro & giảm thiểu**:
  - *Xung đột cổng trên môi trường dev*: Được xử lý bằng biến môi trường trong `.env` và tách cổng dev sang `docker-compose.dev.yml`.
  - *Tăng kích thước file upload sau này*: Đã đặt sẵn `client_max_body_size 50m` trên Nginx và có thể nâng cấp sang direct-to-S3 presigned URL ở Phase 2.
