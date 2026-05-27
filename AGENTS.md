# MyInquisitor — AGENTS.md

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
- **Theming**: CSS variables only (e.g. `var(--color-accent)`). No hardcoded colors. Dark/light via `.theme-dark` / `.theme-light` on `<html>`.
- **CSS**: Tailwind v4 utilities. **Never** use unlayered `* { margin:0; padding:0 }` — it beats layered utility classes.
- **`cn()` utility**: simple class joiner (`frontend/src/utils/cn.ts`).
- **UI components** (`components/ui/`): pure presentational, no side effects. Props and CSS variables only.
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
| `docs/manual-tecnico.md` | Full 10-section technical manual (995 lines) |
