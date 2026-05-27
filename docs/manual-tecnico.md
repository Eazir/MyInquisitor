# Manual Técnico — MyInquisitor

Sistema de gestión financiera personal con registro por invitación, arquitectura hexagonal en el backend y React 19 + TypeScript 6 en el frontend.

---

## Índice

1. [Arquitectura General](#1-arquitectura-general)
2. [Backend (Go)](#2-backend-go)
   - 2.1 [Estructura Hexagonal](#21-estructura-hexagonal)
   - 2.2 [Domain](#22-domain)
   - 2.3 [Application](#23-application)
   - 2.4 [Infrastructure](#24-infrastructure)
   - 2.5 [Handler y Router](#25-handler-y-router)
   - 2.6 [DI Manual en main.go](#26-di-manual-en-maingo)
   - 2.7 [Configuración](#27-configuración)
3. [Frontend (React 19 + TypeScript 6)](#3-frontend-react-19--typescript-6)
   - 3.1 [Estructura de Carpetas](#31-estructura-de-carpetas)
   - 3.2 [Sistema de Ruteo](#32-sistema-de-ruteo)
   - 3.3 [Contextos Globales](#33-contextos-globales)
   - 3.4 [Servicios API](#34-servicios-api)
   - 3.5 [Interceptor de Axios](#35-interceptor-de-axios)
   - 3.6 [Componentes UI](#36-componentes-ui)
   - 3.7 [Sistema de Temas](#37-sistema-de-temas)
4. [Sistema de Autenticación](#4-sistema-de-autenticación)
   - 4.1 [Modelo de Tokens](#41-modelo-de-tokens)
   - 4.2 [Flujo de Sesión](#42-flujo-de-sesión)
   - 4.3 [Refresh con Rotación](#43-refresh-con-rotación)
   - 4.4 [Middleware de Autenticación](#44-middleware-de-autenticación)
   - 4.5 [Roles](#45-roles)
5. [Base de Datos](#5-base-de-datos)
   - 5.1 [Esquema](#51-esquema)
   - 5.2 [SQLC](#52-sqlc)
   - 5.3 [Migraciones](#53-migraciones)
6. [Seguridad](#6-seguridad)
   - 6.1 [Encriptación de PII (AES-256-GCM)](#61-encriptación-de-pii-aes-256-gcm)
   - 6.2 [Hashing de Contraseñas (Bcrypt)](#62-hashing-de-contraseñas-bcrypt)
   - 6.3 [CORS](#63-cors)
   - 6.4 [Invitaciones (One-Time Token)](#64-invitaciones-one-time-token)
7. [Internacionalización (i18n)](#7-internacionalización-i18n)
8. [Docker y Despliegue](#8-docker-y-despliegue)
9. [Mantenimiento](#9-mantenimiento)
   - 9.1 [Agregar un Endpoint Nuevo](#91-agregar-un-endpoint-nuevo)
   - 9.2 [Agregar una Página Nueva](#92-agregar-una-página-nueva)
   - 9.3 [Agregar Traducciones](#93-agregar-traducciones)
   - 9.4 [Comandos Útiles](#94-comandos-útiles)
   - 9.5 [Solución de Problemas Comunes](#95-solución-de-problemas-comunes)
10. [Referencia de API](#10-referencia-de-api)

---

## 1. Arquitectura General

```
┌──────────────────────────────────────────────────────┐
│                   Frontend (React)                     │
│  localhost:5173 (dev) / localhost:3000 (docker)        │
│                                                        │
│  Pages → Services (api.ts) ──axios──→ Backend API       │
│       ↑                              (REST JSON)       │
│  Contexts / Components                                   │
└──────────────────────────┬───────────────────────────┘
                           │ HTTP (JSON)
┌──────────────────────────▼───────────────────────────┐
│                   Backend (Go + Gin)                   │
│  localhost:8080                                        │
│                                                        │
│  Handler → UseCase → Repository → PostgreSQL            │
│       ↑                    ↑                           │
│  Middleware (Auth/CORS)   SQLC (queries generadas)      │
└──────────────────────────┬───────────────────────────┘
                           │ TCP :5433 (host) / :5432 (docker)
┌──────────────────────────▼───────────────────────────┐
│                  PostgreSQL 16                          │
│  Container: myinquisitor-db                             │
└──────────────────────────────────────────────────────┘
```

### Stack Tecnológico

| Capa | Tecnología | Versión |
|------|-----------|---------|
| Frontend | React + TypeScript | 19.2 + 6.0 |
| Build | Vite | 8.0 |
| CSS | TailwindCSS | 4.3 |
| Backend | Go + Gin | 1.24 + 1.10 |
| DB | PostgreSQL | 16 |
| ORM | SQLC (codegen) | — |
| Auth | JWT (golang-jwt) | v5 |
| Contenedores | Docker + compose | — |

---

## 2. Backend (Go)

### 2.1 Estructura Hexagonal

```
backend/
├── cmd/
│   ├── api/main.go          # Punto de entrada, DI manual
│   └── hashgen/main.go      # Generador de hash para seed de admin
├── internal/
│   ├── domain/               # Capa DOMAIN — núcleo puro
│   │   ├── entity/           #   Entidades (User, Debt, etc.)
│   │   ├── repository/       #   Interfaces de repositorio
│   │   └── vo/               #   Value Objects (Money, Period)
│   ├── application/          # Capa APPLICATION — casos de uso
│   │   ├── auth/             #   Login, Register, Refresh
│   │   ├── admin/            #   CRUD de usuarios, invites
│   │   ├── debt/             #   Gestión de deudas
│   │   ├── expense/          #   Gastos recurrentes
│   │   ├── accounting/       #   Transacciones, balances, proyecciones
│   │   ├── profile/          #   Perfil (update, cambio password)
│   │   └── dto/              #   Data Transfer Objects
│   └── infrastructure/       # Capa INFRASTRUCTURE — implements
│       ├── api/
│       │   ├── handler/      #   Handlers HTTP
│       │   ├── middleware/   #   Auth, Admin, CORS
│       │   ├── router/       #   Definición de rutas
│       │   └── response/     #   Helpers de respuesta JSON
│       ├── auth/             #   JWT Service
│       ├── config/           #   Carga de configuración
│       ├── persistence/      #   Repositorios PostgreSQL + SQLC
│       └── security/         #   Password (bcrypt) + Encryption (AES)
├── migrations/               # Migraciones SQL
└── Makefile
```

#### Reglas de Dependencia

```
domain → nada (capa más interna)
application → domain
infrastructure → domain + application (implementa interfaces)
cmd/api → todo (wiring)
```

**Nunca**:
- `domain` importa `application` o `infrastructure`
- `application` importa `infrastructure`
- Un paquete `internal` importa `cmd/`

### 2.2 Domain

Contiene las entidades del negocio y sus interfaces de repositorio. Sin dependencias externas (solo `uuid` y `time`).

**Entidades principales** (`internal/domain/entity/`):

| Archivo | Entidad | Campos clave |
|---------|---------|-------------|
| `user.go` | `User` | ID, Email, PasswordHash, FullName, Phone, Role, Active, SuperAdmin |
| `debt.go` | `Debt` | ID, UserID, Name, TotalAmount, RemainingAmount, InterestRate, Status |
| `recurring_expense.go` | `RecurringExpense` | ID, UserID, Name, Amount, Frequency, DueDay, Status |
| `transaction.go` | `Transaction` | ID, UserID, Type, Amount, Description, ReferenceDate |
| `category.go` | `Category` | ID, UserID, Name, Type |
| `monthly_summary.go` | `MonthlySummary` | ID, UserID, Year, Month, TotalIncome, TotalExpenses |
| `debt_monthly_status.go` | `DebtMonthlyStatus` | ID, DebtID, Month, AmountDue, AmountPaid, Paid |
| `invite_token.go` | `InviteToken` | ID, Token, CreatedBy, Used, ExpiresAt |

**Interfaces de repositorio** (`internal/domain/repository/`):

```go
type UserRepository interface {
    Create(ctx, user) error
    GetByID(ctx, id) (*User, error)
    GetByEmail(ctx, email) (*User, error)
    List(ctx, limit, offset) ([]User, int, error)
    Update(ctx, user) error
    UpdatePassword(ctx, id, hash) error
    SetActive(ctx, id, active) error
    Delete(ctx, id) error
}
```

Cada repositorio sigue el mismo patrón. Las implementaciones concretas están en `internal/infrastructure/persistence/`.

### 2.3 Application

Contiene los casos de uso (use cases). Cada uno es un struct con dependencias inyectadas (repositorios + servicios) y un método `Execute`.

**Patrón**:
```go
type CreateUseCase struct {
    repo repository.DebtRepository
}

func NewCreateUseCase(repo repository.DebtRepository) *CreateUseCase {
    return &CreateUseCase{repo: repo}
}

func (uc *CreateUseCase) Execute(ctx context.Context, input dto.CreateDebtInput) (*dto.DebtOutput, error) {
    // validaciones + lógica de negocio
    entity := &entity.Debt{...}
    if err := uc.repo.Create(ctx, entity); err != nil { ... }
    return output, nil
}
```

Los DTOs están en `internal/application/dto/`. Sirven como contratos de entrada/salida entre handlers y use cases.

**Puertos (interfaces de servicio)** en `internal/application/auth/ports.go`:

```go
type TokenService interface {
    GenerateAccessToken(userID, role) (string, error)
    GenerateRefreshToken(userID) (string, error)
    ValidateToken(token) (*Claims, error)
}

type PasswordService interface {
    Hash(password) (string, error)
    Verify(password, hash) bool
}
```

### 2.4 Infrastructure

Implementa las interfaces definidas en domain y application.

#### Persistencia (`internal/infrastructure/persistence/`)

- Usa SQLC para generar código Go a partir de queries SQL.
- Cada repositorio recibe `*sql.DB` y llama a las funciones generadas por SQLC.
- El `UserRepository` usa `EncryptionService` para encriptar/desencriptar PII (email, full_name) al leer/escribir.

#### JWT Service (`internal/infrastructure/auth/jwt_service.go`)

- **Access token**: incluye `user_id`, `role`, `sub`, `exp`, `iat` — expira en 30m.
- **Refresh token**: incluye `user_id`, `sub`, `exp`, `iat` — expira en 720h/30d.
- Ambos se firman con HMAC-SHA256 usando `JWT_SECRET`.

#### Security (`internal/infrastructure/security/`)

- `PasswordService`: bcrypt con costo 12.
- `EncryptionService`: AES-256-GCM, key de 32 bytes (64 chars hex). Se usa para email y full_name antes de insertar en DB.

### 2.5 Handler y Router

**Handlers** (`internal/infrastructure/api/handler/`):
- Reciben `*gin.Context`.
- Extraen parámetros (path, query, body).
- Llaman al use case correspondiente.
- Responden usando helpers de `internal/infrastructure/api/response/`.

```go
func (h *DebtHandler) Create(c *gin.Context) {
    var input dto.CreateDebtInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.ValidationError(c, "invalid body", err.Error())
        return
    }
    result, err := h.createUC.Execute(c.Request.Context(), input)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "msg")
        return
    }
    response.Success(c, http.StatusCreated, result)
}
```

**Router** (`internal/infrastructure/api/router/router.go`):
- Define grupos de rutas: `/api/v1/auth`, `/api/v1/debts`, etc.
- Las rutas protegidas usan `authMW.Authenticate()`.
- Las rutas de admin usan además `adminMW.RequireSuperAdmin()`.

### 2.6 DI Manual en main.go

No se usa framework de DI. Todo se instancia en `cmd/api/main.go` en orden:

```go
// 1. Config
cfg, err := config.Load()

// 2. DB
db, err := persistence.NewPostgresDB(ctx, cfg.DatabaseURL)

// 3. Servicios
encryptSvc := security.NewEncryptionService(cfg.EncryptionKey)
pwdSvc := security.NewPasswordService()
jwtSvc := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpiration, cfg.RefreshExpiration)

// 4. Repositorios
userRepo := persistence.NewUserRepository(db, encryptSvc)
debtRepo := persistence.NewDebtRepository(db)
// ...

// 5. Use cases
registerUC := appAuth.NewRegisterUseCase(userRepo, inviteRepo, pwdSvc, jwtSvc)
// ...

// 6. Handlers
authH := handler.NewAuthHandler(registerUC, loginUC, refreshUC)
// ...

// 7. Middleware
authMW := middleware.NewAuthMiddleware(jwtSvc)
// ...

// 8. Router
router.Setup(r, authH, profileH, ...)
```

Para agregar una nueva dependencia, se sigue este orden.

### 2.7 Configuración

Variables de entorno (`.env` en la raíz del proyecto):

| Variable | Default | Descripción |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Puerto del servidor HTTP |
| `DATABASE_URL` | — (requerida) | DSN de PostgreSQL |
| `JWT_SECRET` | — (requerido) | Clave secreta para firmar JWT (min 32 chars) |
| `JWT_ACCESS_EXPIRATION` | `30m` | Duración del access token |
| `JWT_REFRESH_EXPIRATION` | `720h` | Duración del refresh token (30 días) |
| `ENCRYPTION_KEY` | — (requerido) | Clave AES-256-GCM en hex (64 chars = 32 bytes) |
| `ALLOWED_ORIGINS` | `http://localhost:5173` | Orígenes CORS separados por coma |
| `ENVIRONMENT` | `development` | Modo de ejecución |

---

## 3. Frontend (React 19 + TypeScript 6)

### 3.1 Estructura de Carpetas

```
frontend/src/
├── components/
│   ├── layout/           # AppLayout, Sidebar, Header, PrivateRoute, AdminRoute, PageContainer
│   └── ui/               # Button, Input, Select, Table, Card, Modal, Badge, Loading, EmptyState
├── contexts/             # AuthContext, ThemeContext, LanguageContext
├── i18n/                 # en.ts, es.ts, index.ts (traducciones)
├── pages/
│   ├── auth/             # LoginPage, RegisterPage
│   ├── dashboard/        # DashboardPage
│   ├── debts/            # DebtsListPage, DebtDetailPage
│   ├── expenses/         # ExpensesPage
│   ├── accounting/       # AccountingPage
│   └── planning/         # PlanningPage
│   └── admin/            # AdminPage
│   └── settings/         # SettingsPage
├── services/             # api.ts, auth.ts, debts.ts, expenses.ts, accounting.ts, admin.ts, profile.ts
├── themes/               # theme.css (variables CSS para light/dark)
├── types/                # auth.ts (interfaces compartidas)
├── utils/                # cn.ts (clsx + tailwind-merge)
├── test/                 # setup.ts (config vitest + jsdom)
├── App.tsx               # Router principal
├── index.css             # Estilos globales + Tailwind
└── main.tsx              # Entry point
```

### 3.2 Sistema de Ruteo

Definido en `App.tsx`:

```tsx
<BrowserRouter>
  <ThemeProvider>
    <LanguageProvider>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register/:token" element={<RegisterPage />} />
          <Route element={<PrivateRoute />}>           ← requiere sesión
            <Route element={<AppLayout />}>             ← layout con sidebar + header
              <Route path="/dashboard" ... />
              <Route path="/debts" ... />
              <Route path="/settings" ... />
              <Route element={<AdminRoute />}>           ← requiere super_admin
                <Route path="/admin" ... />
              </Route>
            </Route>
          </Route>
        </Routes>
      </AuthProvider>
    </LanguageProvider>
  </ThemeProvider>
</BrowserRouter>
```

**PrivateRoute** (`components/layout/PrivateRoute.tsx`):
- Si `loading` → muestra `<Loading />`
- Si `!isAuthenticated` → `<Navigate to="/login" />`
- Si autenticado → `<Outlet />`

**AdminRoute** (`components/layout/AdminRoute.tsx`):
- Igual que PrivateRoute pero verifica `isAdmin`.

### 3.3 Contextos Globales

#### AuthContext (`contexts/AuthContext.tsx`)

```tsx
interface AuthContextType {
  user: User | null;            // datos del usuario logueado
  loading: boolean;             // true mientras se verifica sesión
  login(email, password);        // login + guarda tokens
  register(email, pass, name, inviteToken);  // registro + guarda tokens
  logout();                     // borra tokens, limpia estado
  isAuthenticated: boolean;     // !!user
  isAdmin: boolean;             // user?.role === 'super_admin'
  setUser(user);                // actualiza usuario (usado por SettingsPage)
}
```

En el mount, lee `refreshToken` de localStorage y llama a `/auth/refresh`. Si falla, borra el token. Si éxito, guarda el nuevo token y setea el usuario.

#### ThemeContext (`contexts/ThemeContext.tsx`)

```tsx
interface ThemeContextType {
  theme: 'light' | 'dark';
  toggleTheme();   // alterna + persiste en localStorage
}
```

#### LanguageContext (`contexts/LanguageContext.tsx`)

```tsx
interface LanguageContextType {
  language: 'en' | 'es';
  setLanguage(lang);           // cambia + persiste en localStorage
  t(key, params?): string;     // traduce una clave con interpolación {{var}}
}
```

**Uso en componentes**:
```tsx
const { t } = useLanguage();
return <h1>{t('dashboard.title')}</h1>;        // "Dashboard" / "Dashboard"
<p>{t('common.loading')}</p>;                    // "Loading..." / "Cargando..."
<p>{t('admin.editUser', { name: user.full_name })}</p>;  // "Edit User: Juan" / "Editar Usuario: Juan"
```

### 3.4 Servicios API

Cada servicio es un módulo que exporta funciones async que llaman a la instancia `api` de axios:

```tsx
// services/debts.ts
import api from './api';

export interface Debt { ... }

export const debtsApi = {
  list: () => api.get('/debts').then(r => r.data.data),
  create: (input) => api.post('/debts', input).then(r => r.data.data),
  getByID: (id) => api.get(`/debts/${id}`).then(r => r.data.data),
  getMonthlyStatus: (id) => api.get(`/debts/${id}/monthly`).then(r => r.data.data),
};
```

La instancia `api` en `services/api.ts` tiene `baseURL` configurada desde `VITE_API_URL` (default `http://localhost:8080/api/v1`).

### 3.5 Interceptor de Axios

**Request interceptor**: agrega `Authorization: Bearer <accessToken>` si hay token en memoria.

**Response interceptor**: maneja renovación automática de tokens:

```
Request → 401 → ¿Ya se intentó refresh?
  ├── No → ¿Otro refresh en curso?
  │   ├── Sí → encolar request, esperar nuevo token
  │   └── No → llamar a /auth/refresh
  │       ├── Éxito → guardar nuevo access + refresh, reintentar request original
  │       └── Falla → borrar tokens, window.location = '/login'
  └── Sí → reject original
```

El interceptor usa `axios.post` directo (no la instancia `api`) para evitar bucles infinitos.

### 3.6 Componentes UI

Todos en `components/ui/`. Siguen TailwindCSS v4 con variables CSS para colores.

| Componente | Props clave |
|-----------|------------|
| `Button` | `variant` (primary/secondary/danger/ghost), `loading`, `size` (sm/md) |
| `Input` | `label`, `error`, `type`, className adicional |
| `Select` | `label`, `options: {value, label}[]` |
| `Card` | `title`, `subtitle`, `variant` (default/stats), `className` |
| `Table` | `columns: Column<T>[]`, `data: T[]`, `variant` (striped/default), `onRowClick` |
| `Modal` | `isOpen`, `onClose`, `title`, children |
| `Badge` | `variant` (info/success/warning/danger) |
| `Loading` | `size` (sm/md/lg), `text` (mensaje opcional) |
| `EmptyState` | `title`, `description`, `action` (ReactNode opcional) |

### 3.7 Sistema de Temas

Definido en `themes/theme.css` con variables CSS:

```css
:root {
  --color-bg-primary: #ffffff;
  --color-bg-secondary: #f8fafc;
  --color-text-primary: #1e293b;
  --color-accent: #3b82f6;
  --sidebar-width: 260px;
  --header-height: 64px;
  /* ... más variables ... */
}

.dark {
  --color-bg-primary: #0f172a;
  --color-bg-secondary: #1e293b;
  /* ... */
}
```

---

## 4. Sistema de Autenticación

### 4.1 Modelo de Tokens

| Propiedad | Access Token | Refresh Token |
|-----------|-------------|---------------|
| **Propósito** | Autenticar cada request (header) | Obtener nuevos access tokens |
| **Duración** | 30 minutos | 720 horas (30 días) |
| **Claims** | `user_id`, `role`, `sub`, `exp`, `iat` | `user_id`, `sub`, `exp`, `iat` |
| **Storage** | Memoria (variable JS) | `localStorage` |
| **Se envía en** | `Authorization: Bearer <token>` | Body de `POST /auth/refresh` |
| **Rota** | No (se genera uno nuevo en cada refresh) | Sí (cada refresh genera uno nuevo) |

### 4.2 Flujo de Sesión

```
LOGIN:
  POST /auth/login { email, password }
    → backend valida credenciales
    → genera access_token (30m) + refresh_token (30d)
    → frontend: access → memoria, refresh → localStorage

CADA REQUEST:
  api.get('/debts') con Authorization: Bearer <access>
    → middleware valida access_token (firma + exp)
    → extrae user_id, lo pasa al handler

ACCESS EXPIRA (30m):
  api.get('/debts') → 401
    → interceptor captura 401
    → POST /auth/refresh { refresh_token }
    → backend valida refresh_token
    → genera NUEVO access_token + NUEVO refresh_token
    → frontend: actualiza ambos, reintenta request original

RECARGA DE PÁGINA:
  AuthContext.mount()
    → lee refresh_token de localStorage
    → POST /auth/refresh { refresh_token }
    → si éxito: restaura sesión completa
    → si falla: redirige a /login

LOGOUT:
  logout()
    → borra access_token de memoria
    → borra refresh_token de localStorage
    → setUser(null)
    → (backend: no hace nada, los tokens siguen siendo válidos hasta expirar)
```

### 4.3 Refresh con Rotación

Cada vez que se llama a `/auth/refresh`:
1. Se valida el refresh token actual.
2. Se generan UN NUEVO access token y UN NUEVO refresh token.
3. El refresh token anterior queda huérfano (sigue siendo válido hasta su expiración, pero el frontend ya no lo tiene).

**Importante**: No hay blacklist de tokens en DB. Si un token es robado, puede ser usado hasta que expire. Por eso el access token tiene expiración corta (30m) y el refresh se rota.

### 4.4 Middleware de Autenticación

**AuthMiddleware** (`middleware/auth.go`):
- Extrae `Authorization: Bearer <token>`.
- Llama a `tokenService.ValidateToken()`.
- Si es válido: setea `user_id` y `role` en el contexto de Gin.
- Si no: responde 401.

**AdminMiddleware** (`middleware/admin.go`):
- Lee `role` del contexto (seteado por AuthMiddleware).
- Si no es `super_admin`: responde 403.

**CORSMiddleware** (`middleware/cors.go`):
- Verifica el header `Origin` contra `ALLOWED_ORIGINS`.
- Setea headers CORS.
- Responde 204 a preflight (OPTIONS).

### 4.5 Roles

| Rol | Acceso |
|-----|--------|
| `user` | CRUD de debts, expenses, accounting, planning, settings |
| `super_admin` | Todo lo anterior + `/admin/*` (gestionar usuarios, invites) |

---

## 5. Base de Datos

### 5.1 Esquema

El esquema se define en migraciones SQL en `backend/migrations/`. Tablas principales:

- `users` — id (uuid), email_hash (text UNIQUE), email (text encrypted), full_name (text encrypted), password_hash, phone, role, active, super_admin, created_at, updated_at
- `debts` — id, user_id, name, total_amount, remaining_amount, interest_rate, total_installments, current_installment, start_date, status
- `debt_monthly_status` — id, debt_id, month (YYYY-MM), amount_due, amount_paid, paid
- `recurring_expenses` — id, user_id, name, amount, frequency, due_day, start_date, status
- `expense_monthly_status` — id, expense_id, month, amount, paid
- `transactions` — id, user_id, type (income/expense/transfer), amount, description, reference_date, category_id
- `monthly_summaries` — id, user_id, year, month, total_income, total_expenses, net_balance
- `categories` — id, user_id, name, type
- `invite_tokens` — id, token, created_by, used, expires_at, created_at

### 5.2 SQLC

Las queries SQL están en `backend/internal/infrastructure/persistence/sqlc/`. SQLC genera código Go tipado a partir de archivos `.sql`.

Para regenerar después de cambiar una query:
```bash
cd backend && sqlc generate
```

### 5.3 Migraciones

Las migraciones están en `backend/migrations/`. Se ejecutan en orden numérico en el startup del backend.

---

## 6. Seguridad

### 6.1 Encriptación de PII (AES-256-GCM)

- **Qué se encripta**: `email` y `full_name` en la tabla `users`.
- **Cómo**: El `EncryptionService` en `internal/infrastructure/security/encryption.go` usa AES-256 en modo GCM.
- **Key**: 32 bytes (64 caracteres hex), configurable vía `ENCRYPTION_KEY`.
- **Para qué**: Proteger datos personales en reposo. Si alguien accede a la DB, no puede leer emails ni nombres en texto plano.
- **Búsqueda por email**: Se usa una columna separada `email_hash` (SHA-256 del email en minúsculas) para lookup UNIQUE sin desencriptar.

### 6.2 Hashing de Contraseñas (Bcrypt)

- Algoritmo: bcrypt con costo 12 (`internal/infrastructure/security/password.go`).
- Se almacena solo el hash, nunca la contraseña en texto plano.
- La verificación se hace comparando el hash contra la contraseña ingresada.

### 6.3 CORS

El `CORSMiddleware` en `middleware/cors.go`:
- Solo permite orígenes listados en `ALLOWED_ORIGINS`.
- Headers permitidos: `Authorization`, `Content-Type`.
- Métodos: `GET, POST, PUT, DELETE, OPTIONS`.
- Max age: 86400s (24h) para cachear preflight.
- Credentials: true.

### 6.4 Invitaciones (One-Time Token)

- Se generan tokens de 32 bytes aleatorios (hex) con expiración de 72h.
- Solo el super admin puede generar invitaciones (`POST /admin/invite`).
- Al registrarse, el token se marca como usado (`used = true`).
- No se puede reusar el mismo token.

---

## 7. Internacionalización (i18n)

### Archivos

| Archivo | Contenido |
|---------|-----------|
| `frontend/src/i18n/en.ts` | Traducciones en inglés (~170 claves) |
| `frontend/src/i18n/es.ts` | Traducciones en español (~170 claves) |
| `frontend/src/i18n/index.ts` | `flatten()` + export de `translations` por idioma |

### Estructura de claves

Las traducciones están organizadas por módulos con nombres anidados:

```typescript
export const en = {
  common: { loading: 'Loading...', save: 'Save', cancel: 'Cancel' },
  dashboard: { title: 'Dashboard', description: 'Overview...', totalManaged: 'Total Managed' },
  auth: { signIn: 'Sign In', email: 'Email', password: 'Password' },
  debts: { title: 'Debts', addDebt: 'Add Debt', active: 'Active' },
  settings: { title: 'Settings', profile: 'Profile', language: 'Language' },
  // ... más módulos
};
```

El `flatten()` convierte `{ common: { loading: '...' } }` en `{ 'common.loading': '...' }`.

### Cómo agregar una nueva clave

1. Agregar la clave en `en.ts` y `es.ts` dentro del módulo correspondiente.
2. Usarla en el componente con `{t('modulo.clave')}`.
3. Si tiene interpolación: `{t('modulo.clave', { nombre: valor })}` y en las traducciones usar `{{nombre}}`.

### Fallback

Si una clave no existe en el idioma activo, se busca en inglés. Si tampoco existe, se devuelve la clave literal.

---

## 8. Docker y Despliegue

### Estructura

```
docker/
├── docker-compose.yml        # Orquestación de 3 servicios
├── backend/Dockerfile         # Multi-stage build (scratch, ~15MB)
├── frontend/
│   ├── Dockerfile             # Build + nginx (SPA)
│   └── nginx.conf             # Proxy reverso + SPA fallback
└── db/
    └── init/01-setup.sql      # Inicialización de PostgreSQL
```

### Comandos

```bash
# Iniciar solo DB (desarrollo)
sudo docker compose -f docker/docker-compose.yml up -d db

# Iniciar todo
sudo docker compose -f docker/docker-compose.yml up -d --build

# Ver logs
sudo docker compose -f docker/docker-compose.yml logs -f

# Detener
sudo docker compose -f docker/docker-compose.yml down
```

### Puertos

| Servicio | Puerto host | Puerto contenedor |
|----------|------------|------------------|
| DB | 5433 | 5432 |
| Backend | 8080 | 8080 |
| Frontend | 3000 | 80 |

### nginx (frontend)

```nginx
location /api/ {
    proxy_pass http://backend:8080/;    # proxy reverso a la API
}
location / {
    try_files $uri $uri/ /index.html;    # SPA fallback
}
```

---

## 9. Mantenimiento

### 9.1 Agregar un Endpoint Nuevo

Ejemplo: agregar `PUT /profile/avatar` (backend + frontend).

**Backend:**

1. **DTO** en `internal/application/dto/profile.go`:
   ```go
   type UpdateAvatarInput struct {
       AvatarURL string `json:"avatar_url"`
   }
   ```

2. **Use case** en `internal/application/profile/avatar.go`:
   ```go
   type UpdateAvatarUseCase struct { userRepo repository.UserRepository }
   func (uc *UpdateAvatarUseCase) Execute(ctx, userID, input) (*dto.UpdateProfileOutput, error) { ... }
   ```

3. **Handler** en `internal/infrastructure/api/handler/profile_handler.go`:
   ```go
   func (h *ProfileHandler) UpdateAvatar(c *gin.Context) { ... }
   ```

4. **Router** en `router.go`:
   ```go
   profile.PUT("/avatar", profileH.UpdateAvatar)
   ```

5. **Wiring** en `cmd/api/main.go`:
   ```go
   profileAvatarUC := appProfile.NewUpdateAvatarUseCase(userRepo)
   profileH := handler.NewProfileHandler(profileUpdateUC, profileChangePasswordUC, profileAvatarUC)
   ```

**Frontend:**

6. **Servicio** en `services/profile.ts`:
   ```tsx
   export async function updateAvatar(url: string): Promise<UpdateProfileResponse> {
       const { data } = await api.put('/profile/avatar', { avatar_url: url });
       return data.data;
   }
   ```

7. **Llamada** desde el componente:
   ```tsx
   import { updateAvatar } from '../../services/profile';
   // ...
   const result = await updateAvatar(url);
   setUser(result.user);
   ```

### 9.2 Agregar una Página Nueva

1. Crear `pages/mimodulo/MiPagina.tsx`:
   ```tsx
   import { useLanguage } from '../../contexts/LanguageContext';
   import { PageContainer } from '../../components/layout/PageContainer';

   export function MiPagina() {
       const { t } = useLanguage();
       return (
           <PageContainer>
               <h2>{t('mimodulo.title')}</h2>
           </PageContainer>
       );
   }
   ```

2. Agregar ruta en `App.tsx`:
   ```tsx
   import { MiPagina } from './pages/mimodulo/MiPagina';
   // ...
   <Route path="/mi-ruta" element={<MiPagina />} />
   ```

3. Agregar a `routeTitles` en `AppLayout.tsx`.

4. Agregar al `Sidebar.tsx`.

5. Agregar traducciones en `en.ts` y `es.ts`.

### 9.3 Agregar Traducciones

1. En `i18n/en.ts`:
   ```typescript
   mimodulo: {
       title: 'My Module',
       description: 'This is my new module',
   }
   ```

2. En `i18n/es.ts`:
   ```typescript
   mimodulo: {
       title: 'Mi Módulo',
       description: 'Este es mi nuevo módulo',
   }
   ```

3. Usar en el componente:
   ```tsx
   const { t } = useLanguage();
   t('mimodulo.title')    // → "My Module" o "Mi Módulo"
   t('mimodulo.description')
   ```

### 9.4 Comandos Útiles

```bash
# ─── Backend ───
cd backend

go run ./cmd/api/...              # iniciar servidor
go build ./cmd/api/...            # compilar
go test ./...                     # ejecutar tests
go run ./cmd/hashgen "email" "password" "name"  # generar seed admin

# ─── Frontend ───
cd frontend

npm run dev                       # servidor de desarrollo (Vite)
npm run build                     # build producción
npm test                          # tests
npx tsc -b --noEmit              # type checking
npm run lint                      # ESLint

# ─── Docker ───
sudo docker compose -f docker/docker-compose.yml up -d db    # solo DB
sudo docker compose -f docker/docker-compose.yml up -d       # todo
sudo docker compose -f docker/docker-compose.yml down        # detener

# ─── SQLC (regenerar queries) ───
cd backend && sqlc generate

# ─── Seed de super admin ───
go run ./cmd/hashgen "admin@email.com" "password" "Admin Name"
# → copiar el INSERT, ejecutar contra la DB
```

### 9.5 Solución de Problemas Comunes

| Problema | Causa | Solución |
|----------|-------|----------|
| **Al refrescar la página redirige a /login** | Refresh token no tiene `user_id` (tokens viejos) | Hacer login de nuevo (genera token con `user_id`). Verificar que el backend tenga el fix en `GenerateRefreshToken`. |
| **Los márgenes/paddings no se ven** | CSS Cascade Layers: el `* { margin:0; padding:0 }` en `index.css` sin layer sobreescribe Tailwind | Asegurarse de NO tener `* { margin:0 }` suelto. Tailwind v4 pone utilidades en `@layer utilities`. |
| **Error de CORS** | `ALLOWED_ORIGINS` no incluye el origen del frontend | Agregar el origen (ej. `http://localhost:5173`) a `ALLOWED_ORIGINS`. |
| **"invalid or expired refresh token"** | Token expirado o generado sin `user_id` | Refrescar haciendo login de nuevo. Verificar que `JWT_SECRET` no haya cambiado. |
| **Login fails con "invalid email or password"** | Hashgen mal ejecutado o credenciales incorrectas | Revisar el INSERT generado: `email_hash` debe ser SHA-256 del email, no AES. |
| **El frontend no se conecta al backend** | `VITE_API_URL` incorrecto o backend no corriendo | Verificar que backend esté en `http://localhost:8080`. En Docker se usa el proxy de nginx. |
| **Error de compilación Go: undefined reference** | Falta algún parámetro en la DI de `main.go` | Revisar `router.Setup()`, `handler.NewXxxHandler()`, `appXxx.NewXxxUseCase()` — todos los parámetros deben coincidir. |

---

## 10. Referencia de API

### Autenticación

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| POST | `/api/v1/auth/register` | No | Registro con invite token |
| POST | `/api/v1/auth/login` | No | Login |
| POST | `/api/v1/auth/refresh` | No | Refrescar tokens |

### Perfil

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| PUT | `/api/v1/profile` | Sí | Actualizar nombre/email |
| PUT | `/api/v1/profile/password` | Sí | Cambiar contraseña |

### Deudas

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| POST | `/api/v1/debts` | Sí | Crear deuda |
| GET | `/api/v1/debts` | Sí | Listar deudas |
| GET | `/api/v1/debts/:id` | Sí | Detalle de deuda |
| PUT | `/api/v1/debts/:id` | Sí | Actualizar deuda |
| DELETE | `/api/v1/debts/:id` | Sí | Eliminar deuda |
| GET | `/api/v1/debts/:id/monthly` | Sí | Estado mensual |
| PUT | `/api/v1/debts/:id/monthly/:year/:month/pay` | Sí | Marcar cuota como pagada |

### Gastos Recurrentes

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| POST | `/api/v1/expenses` | Sí | Crear gasto |
| GET | `/api/v1/expenses` | Sí | Listar gastos |
| GET | `/api/v1/expenses/:id` | Sí | Detalle de gasto |
| PUT | `/api/v1/expenses/:id` | Sí | Actualizar gasto |
| DELETE | `/api/v1/expenses/:id` | Sí | Eliminar gasto |
| PUT | `/api/v1/expenses/:id/monthly/:year/:month/toggle` | Sí | Marcar mes como pagado/no pagado |

### Contabilidad

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| POST | `/api/v1/accounting/transactions` | Sí | Registrar transacción |
| GET | `/api/v1/accounting/transactions` | Sí | Listar transacciones |
| GET | `/api/v1/accounting/balance/:year/:month` | Sí | Balance mensual |
| GET | `/api/v1/accounting/cash-flow` | Sí | Flujo de caja |
| GET | `/api/v1/accounting/projections` | Sí | Proyecciones |

### Categorías

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| POST | `/api/v1/categories` | Sí | Crear categoría |
| GET | `/api/v1/categories` | Sí | Listar categorías |
| DELETE | `/api/v1/categories/:id` | Sí | Eliminar categoría |

### Admin (solo super_admin)

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| GET | `/api/v1/admin/users` | Sí+Admin | Listar usuarios |
| POST | `/api/v1/admin/users` | Sí+Admin | Crear usuario |
| PUT | `/api/v1/admin/users/:id` | Sí+Admin | Actualizar usuario (requiere admin_password) |
| PUT | `/api/v1/admin/users/:id/activate/:active` | Sí+Admin | Activar/desactivar usuario |
| POST | `/api/v1/admin/invite` | Sí+Admin | Generar invitación |

### Formato de Respuesta

```json
// Éxito
{ "data": { ... } }

// Con paginación
{ "data": [...], "meta": { "page": 1, "limit": 20, "total": 42 } }

// Error
{ "error": { "code": "AUTH_ERROR", "message": "invalid email or password", "details": "..." } }

// Error de validación (422)
{ "error": { "code": "VALIDATION_ERROR", "message": "invalid request body", "details": "..." } }
```
