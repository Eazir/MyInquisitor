# MyInquisitor — AGENTS.md

## Git workflow — MUST FOLLOW
- **`main`**: production-ready, stable. **NEVER** modify `main` unless explicitly asked.
- **All work on secondary branches**: `responsive`, `feature/*`, `fix/*`, etc.
- **Branch naming**: use descriptive kebab-case (`responsive-layout`, `fix-auth-timeout`).
- **Commits**: conventional commits (`feat:`, `fix:`, `chore:`, `refactor:`, `style:`).
- **Pull requests**: when done, the user will request a merge to `main`. Do NOT merge unless told.
- **Current branch**: the active branch is `responsive` (responsive layout changes).
- Do NOT switch branches unless asked.

## Stack (verified)
| Capa | Stack |
|------|-------|
| Backend | Go 1.26, Gin v1.12, pgx v5, SQLC, golang-jwt v5 |
| Frontend | React 19, TypeScript 6, Vite 8, TailwindCSS v4, Vitest v4 |
| DB | PostgreSQL 16, `pgx/v5/pgxpool` |
| Auth | JWT HS256 (access 30m + refresh 720h), bcrypt cost 12, AES-256-GCM (PII) |
| Infra | Docker Compose (3 containers), multi-stage builds |

## Architecture
- **Hexagonal**: `domain/` → `application/` → `infrastructure/`. Domain imports nothing from the project.
- **Manual DI** in `backend/cmd/api/main.go` — no DI framework. All wiring visible in one file.
- **Repository interfaces** in `domain/repository/`, implementations in `infrastructure/persistence/`.

## Backend commands
```bash
cd backend
go build -o bin/myinquisitor ./cmd/api  # build
go run ./cmd/api                         # run dev server
go test ./... -v -cover                  # all tests (17 unit tests, no DB needed)
go vet ./...                             # static analysis
sqlc generate                            # codegen from SQL queries
```

## Frontend commands
```bash
cd frontend
npm run dev        # Vite dev server on port 5173
npm run build      # tsc -b && vite build
npm run test       # vitest run (11 tests)
```

## Critical DB quirks
- Docker maps `5433:5432`. `.env` must use port **5433** for local dev (`localhost:5433`).
- Migrations: `.sql` files in `persistence/migrations/`. `RunMigrations()` in `migrator.go` runs on startup. Uses `schema_migrations` table to track what's been applied. Auto-detects manually-applied migrations on first run (checks if `users` table exists).
- To start only the DB for local dev: `sudo docker-compose up -d db` from `./docker/`.
- **pgx v5.9.2 `pgtype.Numeric.Scan()` does NOT accept `float64`**. Use `strconv.FormatFloat` to string first. All numeric conversions go through `toPGNumeric()` in `persistence/helpers.go` — the single source of truth for this mapping.

## Planning module (added Jun 2026)
- **Endpoint**: `GET /accounting/projections?months=1&history_months=12`
- **Interest rate is MONTHLY**: `monthlyRate := interestRate / 100.0` (NO `/12` division). User inputs monthly rate directly.
- **Debt projection**: from `debt_monthly_status` unpaid records via `ListUnpaidByUserIDAndDateRange` (JOIN to debts for user_id)
- **Averages**: `sum_of_months / count_of_months_with_data` (non-empty months only, avoids dilution)
- **Adjustments** stored in `localStorage['myinquisitor_planning']` per-month (no DB). Cache in `sessionStorage['myinquisitor_planning_cache']` for tab-switch persistence.
- **Month edit**: clicking the month label in the projections table opens a detail view styled like the monthly payments page (summary cards + itemized list). Supports income modifier + add/remove one-time expenses.
- **Month labels**: Spanish month names hardcoded in `balance.go`.
- **Defaults**: 1 month projection, 12 history months (1–24 and 1–36 respectively).

## Auth system — MUST KNOW
- **Auth middleware** sets `c.Set("user_id", ...)` — note the exact key name. All handlers read `c.Get("user_id")`. If these don't match, the handler gets zero UUID and FK constraints fail (500).
- **Access token**: JWT with `user_id` + `role` claims, 30m expiry.
- **Refresh token**: JWT with `user_id` claim (CRITICAL), 720h expiry.
- **Refresh token MUST include `user_id`** in the JWT claims (`jwtClaims{UserID: userID, ...}`). Using `jwt.RegisteredClaims{}` alone breaks session persistence (ValidateToken returns zero UUID → 401 → redirect to /login). Fixed in `jwt_service.go:47-59`.
- **Old tokens** in localStorage (generated before the fix) still lack `user_id`. User must re-login.
- **Token rotation**: `/auth/refresh` returns new access + new refresh token. Interceptor saves new refresh to localStorage.
- **Frontend interceptor** (`services/api.ts`): on 401, queues failed requests, tries refresh once. On refresh failure, clears tokens and redirects to `/login`.

## Security & PII
- `email` + `full_name` are **AES-256-GCM encrypted** before storage.
- `email_hash` = SHA-256 hex (for UNIQUE lookup via index).
- `password_hash` = bcrypt cost 12.
- Admin edit user requires admin's own password verified against bcrypt hash.

## Registration flow
- **Invite-only**: one-time token, 32-byte random hex (64 hex chars), 72h expiry.
- `POST /admin/invite` (super admin only) → auto-copies URL to clipboard + toast notification.
- `GET /register?token=...` validates token, `POST /auth/register` consumes it (marks `used=true`).
- Token is consumed by `RegisterUseCase` (`application/auth/register.go`), which checks `inviteRepo`.

## i18n
- Translations in `frontend/src/i18n/{en,es}.ts`, flat via `flatten()`. Keys use dot notation: `auth.signInTitle`.
- `useLanguage().t('key', { param: val })` — template vars via `{{param}}`.
- Default: `es`. Persisted in `localStorage` key `language`.
- `LanguageProvider` wraps `AuthProvider` in `App.tsx`.

## Frontend patterns
- **Theming**: CSS variables only (e.g. `var(--color-accent)`). No hardcoded colors. Dark/light via `data-theme` attribute on `<html>`.
- **Super admin purple theme**: `AuthContext` toggles `.super-admin` class on `<html>` when `user?.role === 'super_admin'`. `theme.css` overrides `--color-accent` to purple (`#7c3aed` light / `#a78bfa` dark) and `--color-sidebar-active` accordingly. Instant switch on login/logout, no page reload.
- **Responsive layout**: sidebar is a drawer overlay on mobile/tablet (< 1024px), always visible on desktop. Header has hamburger button on mobile/tablet. Breakpoint: `lg:` (1024px).
- **CSS**: Tailwind v4 utilities. **Never** use unlayered `* { margin:0; padding:0 }` — it beats layered utility classes.
- **`cn()` utility**: simple class joiner (`frontend/src/utils/cn.ts`).
- **UI components** (`components/ui/`): pure presentational, no side effects. Props and CSS variables only.
- **PageContainer**: responsive padding `p-4 md:p-6 lg:p-8`.
- **AuthContext** exposes `setUser` publicly for SettingsPage after profile update.
- **Toast**: `toast(message, 'success' | 'error' | 'info')` from `components/ui/Toast`. `<ToastContainer />` already in `App.tsx`.
- **Routes**: `/register/:token` (public), `/settings` (protected), `/admin` (super_admin only via `<AdminRoute>`).

## Testing quirks
- **Backend**: 17 tests, all unit-level. Use Go `testing` package only (no testify). Hand-rolled mock structs implementing interfaces. DB not required.
- **Frontend**: 11 tests across 3 files, Vitest + Testing Library. `jsdom` environment. MSW available but not widely used yet.

## API conventions
- Base: `/api/v1/`
- Standard response: `{ data: ..., meta?: { page, limit, total } }`
- Error response: defined in `infrastructure/api/response/`
- Pagination: `?page=1&limit=20`
- Protected routes use `authMW.Authenticate()` middleware.
- Admin routes additionally use `adminMW.RequireSuperAdmin()`.

## Important files
| File | Purpose |
|------|---------|
| `backend/cmd/api/main.go` | DI wiring, server entrypoint |
| `backend/internal/infrastructure/auth/jwt_service.go` | JWT token generation/validation |
| `backend/internal/infrastructure/persistence/migrator.go` | Auto-migration runner (uses `schema_migrations` table for tracking) |
| `frontend/src/services/api.ts` | Axios instance with 401 interceptor + refresh |
| `frontend/src/contexts/AuthContext.tsx` | Auth state, login/register/logout, session refresh on load |
| `frontend/src/contexts/LanguageContext.tsx` | i18n `t()` hook |
| `frontend/src/components/ui/Toast.tsx` | Global toast notifications |
| `frontend/src/themes/theme.css` | CSS variables for light/dark + super admin purple palette |
| `frontend/src/components/layout/AppLayout.tsx` | Responsive layout: sidebar drawer, header, backdrop |
| `frontend/src/components/layout/Sidebar.tsx` | Sidebar with open/close props, transform animation |
| `frontend/src/components/layout/Header.tsx` | Header with hamburger button, responsive full width |
| `backend/internal/application/debt/installments.go` | Interest calc: `monthlyRate = rate / 100` (no `/12`) |
| `backend/internal/application/accounting/balance.go` | `GetProjectionsUseCase` — full monthly breakdown + averages |
| `backend/internal/infrastructure/persistence/queries/debt_monthly_status.sql` | `ListUnpaidByUserIDAndDateRange` query |
| `backend/internal/infrastructure/persistence/debt_repo.go` | SQL impl of the above |
| `frontend/src/pages/planning/PlanningPage.tsx` | Planning page: table, detail view, localStorage + sessionStorage |
| `frontend/src/services/accounting.ts` | API client with `Projection` / `ExtraExpenseItem` types |
| `docs/manual-tecnico.md` | Full 10-section technical manual (995 lines) |
