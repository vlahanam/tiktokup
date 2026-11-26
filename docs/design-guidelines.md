# Design Guidelines

Tài liệu này quy định hệ thống thiết kế (Design System), quy chuẩn UI/UX và cách sử dụng **Shadcn UI + Tailwind CSS v4** trong dự án **TikTokUp**.

---

## 1. Triết lý thiết kế (Design Principles)

1. **Rõ ràng & Trực diện (Clarity First)**: Ưu tiên thứ bậc thông tin rành mạch, khoảng trắng hợp lý và độ tương phản cao thay vì trang trí thừa.
2. **Nhất quán bằng Design Tokens**: Toàn bộ màu sắc, bo góc và typography đều tham chiếu CSS variables trong `frontend/src/app/globals.css` — không hardcode mã màu hex/rgb rải rác trong component.
3. **Hỗ trợ Dark/Light Mode mặc định**: Mọi màn hình mới phải hiển thị sắc nét ở cả chế độ sáng (`:root`) và chế độ tối (`.dark`) thông qua `next-themes`.
4. **Accessibility là yêu cầu mặc định**: Tuân thủ chuẩn WCAG 2.1 AA, đảm bảo điều hướng bằng bàn phím, focus ring rõ ràng và semantic HTML.

---

## 2. Hệ thống Design System & Tokens

Dự án sử dụng **Shadcn UI `4.23.1`** (`style: "base-nova"`, base color `neutral`, primitives `@base-ui/react`) kết hợp **Tailwind CSS `4.3.3`**.

### 2.1. Bảng màu Semantic (OKLCH Tokens trong `globals.css`)

| Token CSS | Utility Tailwind | Mục đích sử dụng | Light (`:root`) | Dark (`.dark`) |
| :--- | :--- | :--- | :--- | :--- |
| `--background` | `bg-background` | Nền trang chính | `oklch(1 0 0)` | `oklch(0.145 0 0)` |
| `--foreground` | `text-foreground` | Chữ nội dung chính | `oklch(0.145 0 0)` | `oklch(0.985 0 0)` |
| `--card` / `--card-foreground` | `bg-card text-card-foreground` | Bề mặt thẻ (Card) | `oklch(1 0 0)` | `oklch(0.205 0 0)` |
| `--primary` / `--primary-foreground` | `bg-primary text-primary-foreground` | Hành động chính (Primary CTA) | `oklch(0.205 0 0)` | `oklch(0.922 0 0)` |
| `--secondary` | `bg-secondary text-secondary-foreground` | Hành động phụ, nhãn phụ | `oklch(0.97 0 0)` | `oklch(0.269 0 0)` |
| `--muted` / `--muted-foreground` | `bg-muted text-muted-foreground` | Nền phụ tinh tế, chú thích, metadata | `oklch(0.97 0 0)` / `oklch(0.556 0 0)` | `oklch(0.269 0 0)` / `oklch(0.708 0 0)` |
| `--accent` | `bg-accent text-accent-foreground` | Trạng thái hover / active của menu | `oklch(0.97 0 0)` | `oklch(0.269 0 0)` |
| `--destructive` | `text-destructive bg-destructive/10` | Cảnh báo lỗi, hành động xóa/nguy hiểm | `oklch(0.577 0.245 27.325)` | `oklch(0.704 0.191 22.216)` |
| `--border` / `--input` / `--ring` | `border-border`, `ring-ring` | Viền card, viền ô nhập liệu, focus ring | `oklch(0.922 0 0)` | `oklch(1 0 0 / 10%)` |

> [!TIP]
> Khi mở rộng thương hiệu riêng cho TikTokUp, chỉ cần cập nhật token `--primary`, `--accent` hoặc `--chart-1..5` tại một nơi duy nhất trong `frontend/src/app/globals.css`.

### 2.2. Typography

| Vai trò | Font Family | CSS Variable | Quy chuẩn sử dụng |
| :--- | :--- | :--- | :--- |
| **Sans / Heading** | `Geist` (`next/font/google`) | `--font-geist-sans` → `--font-sans` | Tiêu đề, đoạn văn, nút bấm, nhãn điều hướng |
| **Monospace** | `Geist Mono` (`next/font/google`) | `--font-geist-mono` → `--font-mono` | URL endpoint, Request ID, phiên bản runtime, số liệu latency/uptime |

**Thang kích thước chữ (Type Scale)**:
- **Page Title (`h1`)**: `text-3xl sm:text-4xl font-bold tracking-tight`
- **Section Heading (`h2`)**: `text-xl font-semibold tracking-tight`
- **Card Title**: `font-heading text-base font-medium leading-snug`
- **Body / Description**: `text-sm text-muted-foreground`
- **Metadata / Caption**: `text-xs text-muted-foreground`

### 2.3. Spacing, Border Radius & Layout Grid

- **Base Radius**: `--radius: 0.625rem` (`10px`), dẫn xuất tự động:
  - `rounded-sm`: `calc(var(--radius) * 0.6)`
  - `rounded-md`: `calc(var(--radius) * 0.8)`
  - `rounded-lg`: `var(--radius)` (Buttons, Inputs, Alerts)
  - `rounded-xl`: `calc(var(--radius) * 1.4)` (Cards, Dialogs)
- **Container chính**: `mx-auto w-full max-w-6xl px-4 sm:px-6`
- **Khoảng cách chuẩn**:
  - Giữa các section lớn: `space-y-8` (`32px`)
  - Giữa các nhóm nội dung trong section: `space-y-6` (`24px`)
  - Grid các thẻ: `grid gap-4 md:grid-cols-3` (`16px` gap)
- **Breakpoints chuẩn (Tailwind CSS)**:
  - Mobile: `< 640px` (1 cột)
  - Tablet (`sm` / `md`): `640px – 1024px` (2–3 cột)
  - Desktop (`lg` / `xl`): `>= 1024px` (3 cột, khung tối đa `1152px` / `max-w-6xl`)

---

## 3. Thư viện Components & Quy tắc sử dụng

### 3.1. Các Shadcn UI Components đã cài đặt (`frontend/src/components/ui/`)

| Component | File | Trường hợp sử dụng |
| :--- | :--- | :--- |
| `Button` | `button.tsx` | Hành động chính/phụ (`default`, `outline`, `secondary`, `ghost`, `destructive`, `link`) |
| `Card` | `card.tsx` | Nhóm thông tin/metric có cấu trúc (`CardHeader`, `CardAction`, `CardContent`, `CardFooter`) |
| `Badge` | `badge.tsx` | Trạng thái (`Healthy`, `Offline`), nhãn công nghệ, HTTP method |
| `Tabs` | `tabs.tsx` | Chuyển đổi góc nhìn nội dung trong cùng một màn hình |
| `Alert` | `alert.tsx` | Thông báo trạng thái lỗi hoặc cảnh báo nổi bật trên trang |
| `Skeleton` | `skeleton.tsx` | Trạng thái tải dữ liệu tránh giật bố cục (CLS) |
| `Separator` | `separator.tsx` | Đường kẻ phân tách rõ ràng giữa các khu vực |
| `Dialog` | `dialog.tsx` | Modal xác nhận hoặc form thao tác nhanh |
| `DropdownMenu` | `dropdown-menu.tsx` | Menu hành động ngữ cảnh, chọn theme |
| `Input` & `Label` | `input.tsx`, `label.tsx` | Form nhập liệu chuẩn có liên kết `htmlFor` |
| `Toaster` (Sonner) | `sonner.tsx` | Thông báo phản hồi không chặn thao tác (góc dưới phải `bottom-right`) |

### 3.2. Cách thêm Shadcn UI Component mới (chạy qua Docker)

Không chạy `npx` trên host; sử dụng lệnh `Makefile`:
```bash
make shadcn-add COMPONENT=table
```

### 3.3. Lưu ý quan trọng về Shadcn v4 (`@base-ui/react`)

Khác với Shadcn thế hệ Radix dùng `asChild`, phiên bản `base-nova` trong dự án dùng prop **`render`** của Base UI:

```tsx
// ĐÚNG — Render Next.js Link bằng Button styling:
<Button nativeButton={false} render={<Link href="/" />}>
  Return to Dashboard
</Button>

// ĐÚNG — DropdownMenuTrigger bọc Button:
<DropdownMenuTrigger
  render={<Button variant="outline" size="icon" aria-label="Toggle theme" />}
>
  <SunIcon className="size-4" />
</DropdownMenuTrigger>
```

---

## 4. Mẫu tương tác chuẩn (Interaction & Feedback Patterns)

- **Icons**: Dùng thống nhất bộ `lucide-react` với kích thước `size-4` (`16px`) cho nút/tiêu đề nhỏ và `size-3` (`12px`) trong `Badge`.
- **Loading State**: Luôn dùng `<Skeleton />` khớp hình dạng nội dung thật thay vì spinner toàn màn hình; cấp route có `src/app/loading.tsx`.
- **Error Handling**:
  - Lỗi cấp route xử lý tại `src/app/error.tsx` kèm nút `Try again`.
  - Lỗi gọi API cục bộ hiển thị `<Alert variant="destructive">` trong khu vực bị ảnh hưởng và thông báo ngắn qua `toast.error(...)`.
- **Success Feedback**: Hành động do người dùng chủ động kích hoạt (như bấm nút *Refresh Status*, lưu form) phản hồi bằng `toast.success(...)`.
- **Theme Switching**: Nút `ThemeToggle` đặt cố định ở góc phải `Header`, đồng bộ giữa `ThemeProvider` và `Sonner` toast.

---

## 5. Tiêu chuẩn Accessibility (WCAG 2.1 AA)

- Mọi nút chỉ có icon (như `ThemeToggle`) bắt buộc có `aria-label` hoặc `<span className="sr-only">`.
- Không loại bỏ `focus-visible:ring-*` trên các phần tử tương tác.
- Luôn gắn `Label` tương ứng cho mọi `Input` trong form.
- Sử dụng đúng thứ tự heading (`h1` duy nhất mỗi trang → `h2` → `h3`) và các thẻ landmark (`header`, `main`, `footer`).
