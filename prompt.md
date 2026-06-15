# MyInquisitor — Prompt de Desarrollo para Agente de Código

---

## 1. 🎯 Rol del Agente

Eres un **ingeniero de software senior** especializado en desarrollo full-stack con Go y React. Tu misión es construir **MyInquisitor**, una plataforma de control financiero personal, siguiendo estrictamente:

- **Arquitectura hexagonal** en el backend (Go + Gin)
- **Componentes plantilla** en el frontend (React + TypeScript + TailwindCSS)
- **Dockerización completa** (3 contenedores: backend, frontend, DB)
- **Buenas prácticas** de seguridad, testing y bajo acoplamiento

### Cómo debes trabajar:

1. **Fase por fase** — No avances a la siguiente fase hasta que la actual esté completa y verificada.
2. **Verifica siempre** — Al terminar cada fase ejecuta los pasos de verificación indicados (compilar, migrar, testear, inspeccionar).
3. **Pregunta si dudas** — Si encuentras una decisión no cubierta o un error inesperado, detente y consulta antes de continuar.
4. **Sé específico** — Crea archivos completos, no fragmentos. Da el código exacto con paths absolutos.
5. **No añadas dependencias no listadas** sin consultar.
6. **Reporta el progreso** — Al completar cada fase, muestra un resumen de lo creado y el resultado de la verificación.

---

## 2. 📋 Descripción del Proyecto

**MyInquisitor** es una plataforma web de control financiero personal que permite al usuario:

- Llevar un **control y registro de deudas** de forma general y mes a mes.
- Gestionar **gastos recurrentes** (suscripciones, alquiler, etc.).
- Registrar **contabilidad mensual**: ingresos, egresos, balance.
- Visualizar **flujo de caja** por semanas, meses y años.
- Ver **balances**: dinero total manejado, disponible, reservado.
- Marcar qué gastos/deudas del mes están **pagados o pendientes**.
- Hacer **estimaciones de ingresos futuros** (dado por el usuario).
- Autenticación mediante **JWT** (access + refresh tokens).
- **Panel de administración** para un super usuario que gestiona cuentas, permisos y datos de usuarios (con verificación de contraseña del usuario objetivo).
- **Datos de usuarios cifrados** (AES-256-GCM).
- **Base de datos PostgreSQL 16**.

### Principios de diseño:

- **Desacoplamiento total** — Cada funcionalidad (deudas, gastos, contabilidad, admin) debe poder modificarse, agregarse o eliminarse sin afectar a las demás.
- **Tema dinámico** — Toda la UI se maneja con componentes plantilla que leen variables CSS; cambiar el tema (claro/oscuro) o las fuentes impacta globalmente.
- **Responsividad** — Todos los componentes deben funcionar en móvil y escritorio.
- **Código compartido Front/Mobile** — La arquitectura de componentes plantilla debe permitir que en el futuro se cree una app móvil reutilizando la misma lógica de componentes y API.

---

## 3. 🏗️ Stack Tecnológico Detallado

| Capa         | Tecnología               | Versión   | Propósito                                        |
|-------------|--------------------------|-----------|--------------------------------------------------|
| Backend     | Go                       | 1.23+     | Lenguaje principal del servidor                  |
| Backend     | Gin                      | v1.10+    | Framework HTTP, routing, middleware              |
| Backend     | pgx                      | v5        | Driver nativo de PostgreSQL para Go             |
| Backend     | SQLC                     | latest    | Generación de código Go tipado desde SQL        |
| Backend     | golang-migrate           | v4        | Migraciones de base de datos                     |
| Backend     | golang-jwt               | v5        | Creación y validación de tokens JWT              |
| Backend     | bcrypt (x/crypto)        | -         | Hashing de contraseñas                           |
| Backend     | AES-256-GCM (crypto)     | -         | Cifrado de datos sensibles del usuario           |
| Backend     | go-playground/validator  | v10       | Validación de estructuras/requests               |
| Backend     | testify                  | v1        | Testing unitario y aserciones                    |
| Backend     | testcontainers-go        | latest    | Tests de integración con contenedores reales     |
| Frontend    | React                    | 19        | UI library                                       |
| Frontend    | TypeScript               | 5.x       | Tipado estático                                  |
| Frontend    | Vite                     | 6         | Build tool                                       |
| Frontend    | TailwindCSS              | 4         | Estilos utilitarios + CSS variables para temas   |
| Frontend    | axios                    | latest    | Cliente HTTP con interceptors                    |
| Frontend    | react-router-dom         | v7        | Routing SPA                                      |
| Frontend    | recharts                 | latest    | Gráficas de balance y flujo de caja              |
| Frontend    | Vitest                   | latest    | Testing unitario                                 |
| Frontend    | Testing Library          | latest    | Testing de componentes React                     |
| DB          | PostgreSQL               | 16        | Base de datos relacional                         |
| Infra       | Docker                   | latest    | Contenerización                                  |
| Infra       | Docker Compose           | v2        | Orquestación multi-contenedor                    |

### Reglas de dependencias:

- **Backend**: `domain/` no importa nada fuera de la stdlib. `application/` solo importa `domain/`. `infrastructure/` importa `domain/` y `application/`.
- **Frontend**: cada feature (deudas, gastos, etc.) es una carpeta autocontenida dentro de `pages/`. Los componentes compartidos van en `components/` y son puramente presentacionales (reciben props, no llaman APIs directamente).
- **Ninguna capa** debe tener dependencia cíclica.

---

## 4. 📐 Arquitectura Hexagonal (Backend)

```
┌─────────────────────────────────────────────────────────┐
│                    cmd/api/main.go                      │
│              (Punto de entrada + DI manual)             │
└───────────────┬─────────────────────────┬───────────────┘
                │                         │
    ┌───────────▼───────────┐   ┌─────────▼───────────┐
    │  internal/api/        │   │  internal/config/    │
    │  (Handlers Gin)       │   │  (Env loader)        │
    │  (Middleware)         │   │                      │
    │  (Router)             │   └──────────────────────┘
    └───────────┬───────────┘
                │ llama a casos de uso
    ┌───────────▼───────────┐
    │  internal/application/│
    │  (Casos de uso)       │
    │  (DTOs)               │
    │  (Puertos/inputs)     │
    └───────────┬───────────┘
                │ depende de interfaces (dominio)
    ┌───────────▼───────────┐
    │  internal/domain/     │
    │  (Entidades)          │
    │  (Value Objects)      │
    │  (Repository ifaces)  │
    └───────────┬───────────┘
                │ implementado por infraestructura
    ┌───────────▼───────────┐
    │  internal/infra/      │
    │  (persistence/SQLC)   │
    │  (auth/JWT)           │
    │  (security/encrypt)   │
    └───────────────────────┘
```

### Flujo de una petición:

```
Cliente HTTP
    │
    ▼
Gin Router → Middleware (auth, cors, rate-limit)
    │
    ▼
Handler (api/handler/user_handler.go)
    │  valida request → crea DTO de entrada
    ▼
UseCase (application/auth/login.go)
    │  orquesta lógica de negocio
    │  llama a interfaces del dominio
    ▼
Repository (domain/repository/user_repository.go)
    │  (es una interfaz)
    ▼
PostgresRepository (infrastructure/persistence/user_repo.go)
    │  usa código generado por SQLC
    │  encripta datos sensibles antes de insertar
    ▼
PostgreSQL 16
```

### Reglas de la arquitectura:

1. **domain/** no importa nada del proyecto. Solo stdlib.
2. **application/** solo importa domain/ y DTOs (que son planos).
3. **infrastructure/** importa domain/ y application/.
4. **Handlers** convierten requests HTTP → DTOs → llaman UseCase → convierten resultado → response HTTP.
5. **UseCases** no saben que existen HTTP, Gin, SQL, ni JWT. Solo trabajan con interfaces del dominio.
6. **main.go** hace la inyección de dependencias manual: crea instancias concretas de repos, servicios, casos de uso, handlers.

---

## 5. 🧩 Arquitectura de Componentes Plantilla (Frontend)

### Sistema de temas (claro/oscuro):

```css
/* themes/light.css */
:root {
  --color-bg-primary: #ffffff;
  --color-bg-secondary: #f8f9fa;
  --color-text-primary: #1a1a2e;
  --color-text-secondary: #6c757d;
  --color-accent: #4361ee;
  --color-success: #2ec4b6;
  --color-warning: #ff9f1c;
  --color-danger: #e71d36;
  --font-family: 'Inter', system-ui, sans-serif;
  --border-radius: 8px;
  --spacing-unit: 4px;
}
```

```css
/* themes/dark.css */
[data-theme="dark"] {
  --color-bg-primary: #1a1a2e;
  --color-bg-secondary: #16213e;
  --color-text-primary: #e0e0e0;
  /* ... */
}
```

### Componente plantilla (ejemplo):

```tsx
// components/ui/Card.tsx
interface CardProps {
  title?: string;
  subtitle?: string;
  children: React.ReactNode;
  variant?: 'default' | 'stats' | 'highlight';
  className?: string;
}

export function Card({ title, subtitle, children, variant = 'default', className }: CardProps) {
  return (
    <div
      className={cn(
        'rounded-[var(--border-radius)] bg-[var(--color-bg-primary)]',
        'border border-[var(--color-bg-secondary)] p-4',
        'shadow-sm transition-shadow hover:shadow-md',
        className
      )}
      data-variant={variant}
    >
      {title && (
        <h3 className="text-lg font-semibold text-[var(--color-text-primary)]">
          {title}
        </h3>
      )}
      {subtitle && (
        <p className="text-sm text-[var(--color-text-secondary)] mt-1">
          {subtitle}
        </p>
      )}
      <div className="mt-3">{children}</div>
    </div>
  );
}
```

**Principio**: Si cambias `Card.tsx`, TODAS las Cards del sistema se actualizan. Si cambias `--color-bg-primary` en el tema, TODOS los componentes que usen esa variable se actualizan.

### Catálogo de componentes plantilla:

| Componente   | Props clave                          | Variantes                     |
|-------------|--------------------------------------|-------------------------------|
| Button      | variant, size, loading, disabled, icon | primary, secondary, ghost, danger |
| Input       | label, error, helperText, icon       | default, withIcon, textarea    |
| Select      | label, options, error                | default, searchable            |
| Card        | title, subtitle, variant             | default, stats, highlight      |
| Table       | columns, data, loading, sortable     | default, compact, striped      |
| Modal       | title, size, closable                | sm, md, lg, fullscreen         |
| Badge       | variant, size, removable             | info, success, warning, danger |
| Tabs        | tabs, activeTab, onChange            | underline, pills, icons        |
| StatsCard   | label, value, trend, icon            | up, down, neutral              |
| PageHeader  | title, subtitle, actions             | (slot para botones)            |
| DataTable   | columns, data, pagination, loading   | server, client                 |
| EmptyState  | icon, title, description, action     | (slot para CTA)                |
| Loading     | size, text                           | spinner, skeleton, dots        |
| Switch      | label, checked, onChange             | default, withIcon              |

---

## 6. 🗄️ Diseño de Base de Datos

### Tablas y relaciones:

```sql
-- users: almacena cuentas de usuario (super admin incluido)
-- - password_hash: bcrypt
-- - email, full_name, phone: CIFRADOS con AES-256-GCM
-- - role: 'super_admin' | 'user'
-- - active: para desactivar cuentas sin borrarlas
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL UNIQUE,          -- cifrado
    password_hash   TEXT NOT NULL,                  -- bcrypt
    full_name       TEXT NOT NULL,                  -- cifrado
    phone           TEXT,                           -- cifrado
    role            TEXT NOT NULL DEFAULT 'user',
    active          BOOLEAN NOT NULL DEFAULT true,
    super_admin     BOOLEAN NOT NULL DEFAULT false, -- solo 1 puede ser true
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- categories: taxonomía de gastos/ingresos
-- - type: 'income' | 'expense' | 'debt'
-- - Ejemplos: "Salario", "Alquiler", "Comida", "Préstamo personal"
CREATE TABLE categories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    type            TEXT NOT NULL CHECK (type IN ('income', 'expense', 'debt')),
    icon            TEXT,
    color           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- debts: deudas generales (ej: préstamo bancario, deuda con amigo)
-- - total_amount: monto total original
-- - remaining_amount: monto restante por pagar
-- - status: 'active' | 'paid' | 'settled'
CREATE TABLE debts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id     UUID REFERENCES categories(id),
    name            TEXT NOT NULL,
    description     TEXT,
    total_amount    NUMERIC(14,2) NOT NULL,
    remaining_amount NUMERIC(14,2) NOT NULL,
    interest_rate   NUMERIC(5,2) DEFAULT 0,
    status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paid', 'settled')),
    start_date      DATE NOT NULL,
    end_date        DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- debt_monthly_status: seguimiento mes a mes de cada deuda
-- - paid: si ya se pagó la cuota de ese mes
-- - amount_paid: cuanto se pagó ese mes
CREATE TABLE debt_monthly_status (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    debt_id         UUID NOT NULL REFERENCES debts(id) ON DELETE CASCADE,
    month           DATE NOT NULL, -- primer día del mes (2026-05-01)
    amount_due      NUMERIC(14,2) NOT NULL,
    amount_paid     NUMERIC(14,2) DEFAULT 0,
    paid            BOOLEAN NOT NULL DEFAULT false,
    paid_at         TIMESTAMPTZ,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(debt_id, month)
);

-- recurring_expenses: gastos recurrentes (suscripciones, alquiler, etc.)
-- - frequency: 'monthly' | 'yearly' | 'weekly' | 'biweekly'
-- - status: 'active' | 'cancelled'
CREATE TABLE recurring_expenses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id     UUID REFERENCES categories(id),
    name            TEXT NOT NULL,
    description     TEXT,
    amount          NUMERIC(14,2) NOT NULL,
    frequency       TEXT NOT NULL CHECK (frequency IN ('monthly', 'yearly', 'weekly', 'biweekly')),
    due_day         INT, -- día del mes en que vence (1-31)
    status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'cancelled')),
    start_date      DATE NOT NULL,
    end_date        DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- expense_monthly_status: seguimiento mes a mes de gastos recurrentes
CREATE TABLE expense_monthly_status (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expense_id      UUID NOT NULL REFERENCES recurring_expenses(id) ON DELETE CASCADE,
    month           DATE NOT NULL, -- primer día del mes
    paid            BOOLEAN NOT NULL DEFAULT false,
    paid_at         TIMESTAMPTZ,
    amount_paid     NUMERIC(14,2),
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(expense_id, month)
);

-- transactions: registro contable mensual (ingresos y egresos puntuales)
-- - type: 'income' | 'expense' | 'transfer'
-- - Ej: "Venta de objeto", "Pago de factura extra", "Transferencia a ahorros"
CREATE TABLE transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id     UUID REFERENCES categories(id),
    type            TEXT NOT NULL CHECK (type IN ('income', 'expense', 'transfer')),
    amount          NUMERIC(14,2) NOT NULL,
    description     TEXT,
    reference_date  DATE NOT NULL DEFAULT CURRENT_DATE,
    is_recurring    BOOLEAN NOT NULL DEFAULT false,
    recurring_expense_id UUID REFERENCES recurring_expenses(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- monthly_summary: resumen mensual precalculado
CREATE TABLE monthly_summary (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    month           DATE NOT NULL, -- primer día del mes
    total_income    NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_expenses  NUMERIC(14,2) NOT NULL DEFAULT 0,
    total_debt_payments NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_balance     NUMERIC(14,2) NOT NULL DEFAULT 0,
    projected_income NUMERIC(14,2),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, month)
);

-- Índices
CREATE INDEX idx_debts_user_id ON debts(user_id);
CREATE INDEX idx_debts_status ON debts(status);
CREATE INDEX idx_debt_monthly_status_debt_id ON debt_monthly_status(debt_id);
CREATE INDEX idx_debt_monthly_status_month ON debt_monthly_status(month);
CREATE INDEX idx_recurring_expenses_user_id ON recurring_expenses(user_id);
CREATE INDEX idx_expense_monthly_status_expense_id ON expense_monthly_status(expense_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_reference_date ON transactions(reference_date);
CREATE INDEX idx_monthly_summary_user_id ON monthly_summary(user_id);
CREATE INDEX idx_monthly_summary_month ON monthly_summary(month);
CREATE INDEX idx_categories_user_id ON categories(user_id);
```

### Política de cifrado:

- **Se cifra** (AES-256-GCM): `email`, `full_name`, `phone` en tabla `users`.
- **NO se cifra**: contraseñas (van con bcrypt), montos, fechas, estados, relaciones.
- **Clave de cifrado**: se pasa como variable de entorno `ENCRYPTION_KEY` (64 caracteres hex → 32 bytes para AES-256). Distinta por entorno.
- **El cifrado ocurre a nivel de aplicación** (en los repositorios), no en la base de datos.

---

## 7. 🛡️ Buenas Prácticas y Estándares

### Código (Backend)

- **SOLID**: cada struct/función tiene una responsabilidad única.
- **Clean Code**: nombres descriptivos en inglés. Errores siempre manejados con `if err != nil`.
- **Go idioms**: `context.Context` como primer parámetro, errores como valores, zero values significativos.
- **Errores**: usar `fmt.Errorf("context: %w", err)` para wrapping. Definir errores centinela en domain (`var ErrNotFound = errors.New("user not found")`).
- **DTOs**: nunca exponer entidades del dominio directamente en la API. Convertir entre entidad ↔ DTO en los handlers.
- **Comentarios**: solo en interfaces públicas exportadas. Nada de comentarios obvios (`// getUser returns user`).

### Código (Frontend)

- **Componentes puros**: los componentes en `components/ui/` deben ser puramente presentacionales. No llaman hooks de datos, no llaman a la API.
- **Composición**: los containers (páginas) orquestan datos y se los pasan a componentes presentacionales.
- **Estados**: manejar loading, empty, error, y éxito en cada página.
- **Hooks**: custom hooks para lógica reutilizable (useDebounce, usePagination).
- **TypeScript**: evitar `any`. Tipar props de componentes siempre.

### Seguridad

| Aspecto              | Implementación                                          |
|---------------------|----------------------------------------------------------|
| Autenticación       | JWT access token (15 min) + refresh token (7 días)      |
| Almacenamiento JWT  | Refresh token en httpOnly cookie seguro. Access en memoria (variable) o Authorization header |
| Contraseñas         | bcrypt con costo 12                                     |
| Cifrado datos       | AES-256-GCM con clave de 32 bytes por entorno           |
| SQL Injection       | Prevenido por SQLC (queries parametrizadas siempre)     |
| CORS                | Orígenes permitidos configurables por entorno            |
| Rate Limiting       | Endpoints de auth: 5 intentos por minuto por IP         |
| Headers             | Helmet-like: X-Frame-Options, X-Content-Type-Options, Strict-Transport-Security |
| Input validation    | Todos los endpoints validan el body/query con go-playground/validator |

### API REST

```
Formato request body (POST/PUT/PATCH):  application/json
Formato response éxito:                  { "data": ..., "meta": { "page", "limit", "total" } }
Formato response error:                  { "error": "mensaje", "code": "VALIDATION_ERROR", "details": [...] }
Paginación:                              ?page=1&limit=20
Fechas:                                  ISO 8601 (2026-05-25T10:30:00Z)
Moneda:                                  NUMERIC(14,2) en centésimas. Siempre en USD (o moneda única configurable).
```

### Base de Datos

- **Migraciones**: siempre archivos numerados (`001_`, `002_`) con `up` y `down`.
- **Transacciones**: usar `pgx.BeginTx` en operaciones multi-tabla (ej: crear deuda + crear primer debt_monthly_status).
- **Soft delete**: las tablas que lo requieren tienen `deleted_at TIMESTAMPTZ`. No se borra físicamente.
- **Índices**: crear siempre para claves foráneas y columnas de búsqueda frecuente.

### Testing

| Capa               | Herramienta              | Enfoque                                      |
|--------------------|--------------------------|----------------------------------------------|
| domain             | testing + testify        | Tests de value objects y lógica de entidades |
| application        | testing + testify + mock | Table-driven tests con repositorios mock     |
| infrastructure    | testcontainers-go + pgx  | Tests de integración contra PostgreSQL real  |
| handlers           | httptest + testify       | Tests de rutas Gin con responses JSON        |
| frontend hooks     | Vitest                   | Tests de hooks y utilidades                  |
| frontend components| Vitest + Testing Library | Render + interacción del usuario              |
| frontend pages     | Vitest + MSW             | Mock de API, render completo                 |

---

## 8. 📋 Fases de Desarrollo

---

### Fase 0: Preparación del Entorno

**🎯 Objetivo**: Crear la estructura del monorepo, inicializar módulos, instalar dependencias, preparar herramientas.

**📁 Archivos a crear:**

```
/mnt/datos/INFO-ARCH/documentos/Dev/MyInquisitor/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── entity/
│   │   │   ├── vo/
│   │   │   └── repository/
│   │   ├── application/
│   │   │   ├── auth/
│   │   │   ├── debt/
│   │   │   ├── expense/
│   │   │   ├── accounting/
│   │   │   ├── admin/
│   │   │   └── dto/
│   │   └── infrastructure/
│   │       ├── persistence/
│   │       │   ├── migrations/
│   │       │   ├── queries/
│   │       │   └── sqlc/
│   │       ├── auth/
│   │       ├── security/
│   │       ├── api/
│   │       │   ├── handler/
│   │       │   ├── middleware/
│   │       │   └── router/
│   │       └── config/
│   ├── go.mod
│   ├── sqlc.yaml
│   └── Makefile
├── frontend/
│   (generado por Vite + modificaciones)
├── docker/
│   ├── backend/
│   │   └── Dockerfile
│   ├── frontend/
│   │   └── Dockerfile
│   └── docker-compose.yml
├── .env.example
├── .gitignore
└── README.md
```

**🧩 Implementación:**

```bash
# 1. Crear estructura de directorios
mkdir -p backend/cmd/api
mkdir -p backend/internal/{domain/{entity,vo,repository},application/{auth,debt,expense,accounting,admin,dto},infrastructure/{persistence/{migrations,queries,sqlc},auth,security,api/{handler,middleware,router},config}}
mkdir -p docker/{backend,frontend}

# 2. Inicializar Go module
cd backend && go mod init github.com/myinquisitor/backend

# 3. Instalar dependencias Go
go get github.com/gin-gonic/gin@v1.10
go get github.com/jackc/pgx/v5@latest
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/golang-jwt/jwt/v5@latest
go get github.com/go-playground/validator/v10@latest
go get github.com/golang-migrate/migrate/v4@latest
go get github.com/golang-migrate/migrate/v4/database/pgx/v5
go get github.com/golang-migrate/migrate/v4/source/file
go get golang.org/x/crypto/bcrypt
go get github.com/stretchr/testify@latest
go get github.com/testcontainers/testcontainers-go@latest
go get github.com/google/uuid
go get github.com/joho/godotenv
go get github.com/rs/cors

# 4. Inicializar frontend con Vite
cd ../frontend && npm create vite@latest . -- --template react-ts
npm install react-router-dom axios recharts
npm install -D tailwindcss @tailwindcss/vite vitest @testing-library/react @testing-library/jest-dom msw

# 5. Crear .gitignore
```

**✅ Verificación:**
- `cd backend && go build ./...` → compila sin errores
- `cd frontend && npm run build` → compila sin errores
- `tree -L 3` muestra la estructura completa

---

### Fase 1: Diseño de Base de Datos

**🎯 Objetivo**: Crear migraciones SQL y queries SQLC, generar código Go.

**📁 Archivos a crear:**
- `backend/internal/infrastructure/persistence/migrations/001_init.up.sql`
- `backend/internal/infrastructure/persistence/migrations/001_init.down.sql`
- `backend/internal/infrastructure/persistence/queries/` (un archivo .sql por entidad)
- `backend/sqlc.yaml`

**🧩 Implementación:**

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "./internal/infrastructure/persistence/migrations/"
    queries: "./internal/infrastructure/persistence/queries/"
    gen:
      go:
        package: "sqlc"
        out: "./internal/infrastructure/persistence/sqlc/"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true
        json_tags_case_style: "snake"
        overrides:
          - db_type: "uuid"
            go_type: "github.com/google/uuid.UUID"
          - db_type: "timestamptz"
            go_type: "time.Time"
          - db_type: "numeric"
            go_type: "float64"
```

Crear archivos `.sql` para cada entidad en `queries/`:
- `users.sql`: Create, GetByID, GetByEmail, List, Update, Delete, SetActive
- `debts.sql`: Create, GetByID, GetByUserID, ListActive, Update, Delete, GetByMonth
- `debt_monthly_status.sql`: Create, GetByDebtIDAndMonth, MarkAsPaid, ListByMonth
- `recurring_expenses.sql`: Create, GetByID, GetByUserID, ListActive, Update, Delete
- `expense_monthly_status.sql`: Create, GetByExpenseIDAndMonth, MarkAsPaid, ListByMonth
- `transactions.sql`: Create, GetByID, GetByUserIDAndMonth, ListByDateRange, GetSummary
- `monthly_summary.sql`: Upsert, GetByUserIDAndMonth, ListByYear
- `categories.sql`: Create, ListByUserID, ListByType, Delete

**✅ Verificación:**
- `cd backend && sqlc generate` → código generado sin errores en `internal/infrastructure/persistence/sqlc/`
- Revisar que los structs generados tengan JSON tags correctos

---

### Fase 2: Esqueleto Backend — Domain + Config

**🎯 Objetivo**: Crear entidades, value objects, interfaces de repositorio, y sistema de configuración.

**📁 Archivos a crear/modificar:**
- `backend/internal/domain/entity/user.go`
- `backend/internal/domain/entity/debt.go`
- `backend/internal/domain/entity/recurring_expense.go`
- `backend/internal/domain/entity/transaction.go`
- `backend/internal/domain/entity/category.go`
- `backend/internal/domain/entity/monthly_summary.go`
- `backend/internal/domain/vo/money.go`
- `backend/internal/domain/vo/period.go`
- `backend/internal/domain/vo/debt_status.go`
- `backend/internal/domain/repository/user_repository.go`
- `backend/internal/domain/repository/debt_repository.go`
- `backend/internal/domain/repository/expense_repository.go`
- `backend/internal/domain/repository/accounting_repository.go`
- `backend/internal/domain/repository/category_repository.go`
- `backend/internal/infrastructure/config/config.go`
- `backend/cmd/api/main.go`

**🧩 Reglas de entidades:**

```go
// Ejemplo: User entity
package entity

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID
    Email        string    // Siempre se almacena cifrado en DB
    PasswordHash string    // bcrypt
    FullName     string    // siempre cifrado en DB
    Phone        *string   // nullable, siempre cifrado en DB
    Role         string    // "super_admin" | "user"
    Active       bool
    SuperAdmin   bool
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type CreateUserInput struct {
    Email    string
    Password string
    FullName string
    Phone    *string
}

type UpdateUserInput struct {
    ID       uuid.UUID
    Email    *string
    FullName *string
    Phone    *string
    Active   *bool
    Role     *string
}
```

```go
// Value Object: Money
package vo

import "errors"

type Money struct {
    amount float64
}

func NewMoney(amount float64) (Money, error) {
    if amount < 0 {
        return Money{}, errors.New("amount cannot be negative")
    }
    return Money{amount: amount}, nil
}

func (m Money) Amount() float64 { return m.amount }
func (m Money) Add(other Money) Money { return Money{amount: m.amount + other.amount} }
func (m Money) Subtract(other Money) (Money, error) {
    if other.amount > m.amount {
        return Money{}, errors.New("insufficient funds")
    }
    return Money{amount: m.amount - other.amount}, nil
}
```

```go
// Interfaces de repositorio
package repository

import (
    "context"
    "github.com/google/uuid"
    "github.com/myinquisitor/backend/internal/domain/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    List(ctx context.Context, page, limit int) ([]entity.User, int, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    SetActive(ctx context.Context, id uuid.UUID, active bool) error
}
```

```go
// Config
package config

import (
    "os"
    "time"
)

type Config struct {
    ServerPort     string
    DatabaseURL    string
    JWTSecret      string
    JWTExpiration  time.Duration // access: 15min
    RefreshExpiration time.Duration // refresh: 7d
    EncryptionKey  string // 32 bytes hex
    AllowedOrigins []string
    Environment    string // "development" | "production"
}

func Load() (*Config, error) {
    // Leer de variables de entorno con valores por defecto
    // Usar os.Getenv, si falta .env.example mostrar error claro
}
```

```go
// main.go (esqueleto)
package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/myinquisitor/backend/internal/infrastructure/config"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    r := gin.Default()

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    log.Printf("Starting server on port %s", cfg.ServerPort)
    r.Run(":" + cfg.ServerPort)
}
```

**✅ Verificación:**
- `cd backend && go build ./...` → compila
- `go vet ./...` → sin errores

---

### Fase 3: Infraestructura de Persistencia

**🎯 Objetivo**: Implementar los repositorios usando SQLC + pgx, con cifrado de campos sensibles.

**📁 Archivos a crear:**
- `backend/internal/infrastructure/persistence/postgres.go` (pool, transactions)
- `backend/internal/infrastructure/persistence/user_repo.go`
- `backend/internal/infrastructure/persistence/debt_repo.go`
- `backend/internal/infrastructure/persistence/expense_repo.go`
- `backend/internal/infrastructure/persistence/accounting_repo.go`
- `backend/internal/infrastructure/persistence/category_repo.go`

**🧩 Detalles de implementación:**

```go
// postgres.go
package persistence

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/myinquisitor/backend/internal/infrastructure/persistence/sqlc" // código generado
)

type PostgresDB struct {
    Pool    *pgxpool.Pool
    Queries *sqlc.Queries
}

func NewPostgresDB(ctx context.Context, databaseURL string) (*PostgresDB, error) {
    pool, err := pgxpool.New(ctx, databaseURL)
    if err != nil {
        return nil, fmt.Errorf("unable to create connection pool: %w", err)
    }
    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("unable to ping database: %w", err)
    }
    return &PostgresDB{
        Pool:    pool,
        Queries: sqlc.New(pool),
    }, nil
}

func (db *PostgresDB) Close() {
    db.Pool.Close()
}
```

```go
// user_repo.go (esquema)
package persistence

import (
    "context"
    "github.com/google/uuid"
    "github.com/myinquisitor/backend/internal/domain/entity"
    "github.com/myinquisitor/backend/internal/domain/repository"
    "github.com/myinquisitor/backend/internal/infrastructure/security"
    "github.com/myinquisitor/backend/internal/infrastructure/persistence/sqlc"
)

type userRepository struct {
    db       *PostgresDB
    encrypt  *security.EncryptionService
}

func NewUserRepository(db *PostgresDB, encrypt *security.EncryptionService) repository.UserRepository {
    return &userRepository{db: db, encrypt: encrypt}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
    // Cifrar email, full_name, phone
    encryptedEmail, err := r.encrypt.Encrypt(user.Email)
    // ... insertar con sqlc ...
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    // Obtener con sqlc, descifrar campos sensibles
}
```

**✅ Verificación:**
- `cd backend && go build ./...`
- Poder conectar a una DB PostgreSQL local (o testcontainers) y ejecutar consultas básicas

---

### Fase 4: Infraestructura de Seguridad

**🎯 Objetivo**: Implementar servicios de JWT, cifrado AES-256-GCM y bcrypt.

**📁 Archivos a crear:**
- `backend/internal/infrastructure/auth/jwt_service.go`
- `backend/internal/infrastructure/security/encryption.go`
- `backend/internal/infrastructure/security/password.go`

**🧩 Detalles:**

```go
// jwt_service.go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    Role   string    `json:"role"`
    jwt.RegisteredClaims
}

type JWTService struct {
    secret           []byte
    accessExpiration  time.Duration // 15min
    refreshExpiration time.Duration // 7d
}

func NewJWTService(secret string, accessExp, refreshExp time.Duration) *JWTService {
    return &JWTService{
        secret:           []byte(secret),
        accessExpiration:  accessExp,
        refreshExpiration: refreshExp,
    }
}

func (s *JWTService) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
    claims := Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secret)
}

func (s *JWTService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
    claims := jwt.RegisteredClaims{
        Subject:   userID.String(),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpiration)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secret)
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return s.secret, nil
    })
    if err != nil {
        return nil, err
    }
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, jwt.ErrSignatureInvalid
    }
    return claims, nil
}
```

```go
// encryption.go (AES-256-GCM)
package security

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "io"
)

type EncryptionService struct {
    key []byte
}

func NewEncryptionService(hexKey string) (*EncryptionService, error) {
    key, err := hex.DecodeString(hexKey)
    if err != nil {
        return nil, errors.New("invalid encryption key: must be hex-encoded 32 bytes")
    }
    if len(key) != 32 {
        return nil, errors.New("encryption key must be exactly 32 bytes (64 hex chars)")
    }
    return &EncryptionService{key: key}, nil
}

func (s *EncryptionService) Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(s.key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return hex.EncodeToString(ciphertext), nil
}

func (s *EncryptionService) Decrypt(cipherHex string) (string, error) {
    ciphertext, err := hex.DecodeString(cipherHex)
    if err != nil {
        return "", err
    }
    block, err := aes.NewCipher(s.key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return "", errors.New("ciphertext too short")
    }
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }
    return string(plaintext), nil
}
```

```go
// password.go
package security

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12

type PasswordService struct{}

func NewPasswordService() *PasswordService { return &PasswordService{} }

func (s *PasswordService) Hash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
    return string(bytes), err
}

func (s *PasswordService) Verify(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
```

**✅ Verificación:**
- `cd backend && go test ./internal/infrastructure/auth/... -v`
- `cd backend && go test ./internal/infrastructure/security/... -v`

---

### Fase 5: Casos de Uso (Application Layer)

**🎯 Objetivo**: Implementar la lógica de negocio pura en casos de uso.

**📁 Archivos a crear:**
- `backend/internal/application/auth/register.go`
- `backend/internal/application/auth/login.go`
- `backend/internal/application/auth/refresh.go`
- `backend/internal/application/debt/create.go`
- `backend/internal/application/debt/list.go`
- `backend/internal/application/debt/mark_paid.go`
- `backend/internal/application/debt/get_monthly.go`
- `backend/internal/application/expense/create.go`
- `backend/internal/application/expense/list.go`
- `backend/internal/application/expense/toggle_paid.go`
- `backend/internal/application/accounting/record.go`
- `backend/internal/application/accounting/monthly_balance.go`
- `backend/internal/application/accounting/cash_flow.go`
- `backend/internal/application/accounting/projections.go`
- `backend/internal/application/admin/list_users.go`
- `backend/internal/application/admin/update_user.go`
- `backend/internal/application/admin/create_user.go`
- `backend/internal/application/admin/deactivate_user.go`
- `backend/internal/application/dto/` (todos los DTOs)

**🧩 Patrón de Use Case:**

```go
// application/auth/login.go
package auth

import (
    "context"
    "errors"
    "github.com/myinquisitor/backend/internal/domain/repository"
    "github.com/myinquisitor/backend/internal/infrastructure/auth"
    "github.com/myinquisitor/backend/internal/infrastructure/security"
)

var (
    ErrInvalidCredentials = errors.New("invalid email or password")
    ErrUserInactive       = errors.New("account is inactive")
)

type LoginInput struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type LoginOutput struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    User         UserDTO `json:"user"`
}

type LoginUseCase struct {
    userRepo repository.UserRepository
    jwtSvc   *auth.JWTService
    pwdSvc   *security.PasswordService
}

func NewLoginUseCase(
    userRepo repository.UserRepository,
    jwtSvc *auth.JWTService,
    pwdSvc *security.PasswordService,
) *LoginUseCase {
    return &LoginUseCase{
        userRepo: userRepo,
        jwtSvc:   jwtSvc,
        pwdSvc:   pwdSvc,
    }
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
    // 1. Buscar usuario por email
    user, err := uc.userRepo.GetByEmail(ctx, input.Email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    // 2. Verificar si está activo
    if !user.Active {
        return nil, ErrUserInactive
    }

    // 3. Verificar contraseña
    if !uc.pwdSvc.Verify(input.Password, user.PasswordHash) {
        return nil, ErrInvalidCredentials
    }

    // 4. Generar tokens
    accessToken, err := uc.jwtSvc.GenerateAccessToken(user.ID, user.Role)
    if err != nil {
        return nil, err
    }
    refreshToken, err := uc.jwtSvc.GenerateRefreshToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &LoginOutput{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        User:         UserDTO{ID: user.ID, Email: user.Email, FullName: user.FullName, Role: user.Role},
    }, nil
}
```

**✅ Verificación:**
- `cd backend && go test ./internal/application/...` (con mocks de repositorios)
- Cada caso de uso debe tener al menos: test de éxito, test de error (credenciales inválidas, usuario no encontrado, etc.)

---

### Fase 6: API REST

**🎯 Objetivo**: Implementar handlers HTTP con Gin, middleware, router, y tests de integración.

**📁 Archivos a crear:**
- `backend/internal/infrastructure/api/middleware/auth.go`
- `backend/internal/infrastructure/api/middleware/admin.go`
- `backend/internal/infrastructure/api/middleware/cors.go`
- `backend/internal/infrastructure/api/middleware/ratelimit.go`
- `backend/internal/infrastructure/api/handler/auth_handler.go`
- `backend/internal/infrastructure/api/handler/debt_handler.go`
- `backend/internal/infrastructure/api/handler/expense_handler.go`
- `backend/internal/infrastructure/api/handler/accounting_handler.go`
- `backend/internal/infrastructure/api/handler/admin_handler.go`
- `backend/internal/infrastructure/api/router/router.go`
- `backend/internal/infrastructure/api/response/response.go`

**🧩 Router structure:**

```go
// router.go
package router

import (
    "github.com/gin-gonic/gin"
    "github.com/myinquisitor/backend/internal/infrastructure/api/handler"
    "github.com/myinquisitor/backend/internal/infrastructure/api/middleware"
)

func Setup(
    r *gin.Engine,
    authH *handler.AuthHandler,
    debtH *handler.DebtHandler,
    expenseH *handler.ExpenseHandler,
    accH *handler.AccountingHandler,
    adminH *handler.AdminHandler,
    authMW *middleware.AuthMiddleware,
    adminMW *middleware.AdminMiddleware,
) {
    r.Use(middleware.CORS())
    r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

    v1 := r.Group("/api/v1")
    {
        // Públicas
        auth := v1.Group("/auth")
        {
            auth.POST("/register", authH.Register)
            auth.POST("/login", authH.Login)
            auth.POST("/refresh", authH.Refresh)
        }

        // Protegidas (requieren JWT)
        protected := v1.Group("")
        protected.Use(authMW.Authenticate())
        {
            // Deudas
            debts := protected.Group("/debts")
            {
                debts.POST("/", debtH.Create)
                debts.GET("/", debtH.List)
                debts.GET("/:id", debtH.GetByID)
                debts.PUT("/:id", debtH.Update)
                debts.DELETE("/:id", debtH.Delete)
                debts.GET("/monthly/:year/:month", debtH.GetMonthlyStatus)
                debts.PUT("/:id/monthly/:year/:month/pay", debtH.MarkAsPaid)
            }

            // Gastos recurrentes
            expenses := protected.Group("/expenses")
            {
                expenses.POST("/", expenseH.Create)
                expenses.GET("/", expenseH.List)
                expenses.GET("/:id", expenseH.GetByID)
                expenses.PUT("/:id", expenseH.Update)
                expenses.DELETE("/:id", expenseH.Delete)
                expenses.PUT("/:id/monthly/:year/:month/pay", expenseH.TogglePaid)
            }

            // Contabilidad
            accounting := protected.Group("/accounting")
            {
                accounting.POST("/transactions", accH.RecordTransaction)
                accounting.GET("/transactions", accH.ListTransactions)
                accounting.GET("/balance/:year/:month", accH.MonthlyBalance)
                accounting.GET("/cash-flow", accH.CashFlow)
                accounting.GET("/projections", accH.Projections)
            }

            // Categorías
            categories := protected.Group("/categories")
            {
                categories.POST("/", accH.CreateCategory)
                categories.GET("/", accH.ListCategories)
                categories.DELETE("/:id", accH.DeleteCategory)
            }

            // Admin (requiere ser super admin)
            admin := protected.Group("/admin")
            admin.Use(adminMW.RequireSuperAdmin())
            {
                admin.GET("/users", adminH.ListUsers)
                admin.POST("/users", adminH.CreateUser)
                admin.PUT("/users/:id", adminH.UpdateUser)
                admin.DELETE("/users/:id", adminH.DeleteUser)
                admin.PUT("/users/:id/activate", adminH.SetActive)
            }
        }
    }
}
```

**🧩 Handler pattern:**

```go
// auth_handler.go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/myinquisitor/backend/internal/application/auth"
)

type AuthHandler struct {
    loginUseCase    *auth.LoginUseCase
    registerUseCase *auth.RegisterUseCase
    refreshUseCase  *auth.RefreshUseCase
}

func NewAuthHandler(login, register, refresh *auth.LoginUseCase, *auth.RegisterUseCase, *auth.RefreshUseCase) *AuthHandler {
    return &AuthHandler{
        loginUseCase:    login,
        registerUseCase: register,
        refreshUseCase:  refresh,
    }
}

func (h *AuthHandler) Login(c *gin.Context) {
    var input auth.LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
        return
    }
    // Validar con go-playground/validator
    if err := validate.Struct(input); err != nil {
        response.ValidationError(c, err)
        return
    }
    result, err := h.loginUseCase.Execute(c.Request.Context(), input)
    if err != nil {
        if errors.Is(err, auth.ErrInvalidCredentials) || errors.Is(err, auth.ErrUserInactive) {
            response.Error(c, http.StatusUnauthorized, "AUTH_ERROR", err.Error())
            return
        }
        response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
        return
    }
    response.Success(c, http.StatusOK, result)
}
```

```go
// response/response.go
package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
    Data  interface{} `json:"data,omitempty"`
    Meta  *Meta       `json:"meta,omitempty"`
    Error *APIError   `json:"error,omitempty"`
}

type Meta struct {
    Page  int   `json:"page"`
    Limit int   `json:"limit"`
    Total int64 `json:"total"`
}

type APIError struct {
    Message string      `json:"message"`
    Code    string      `json:"code"`
    Details interface{} `json:"details,omitempty"`
}

func Success(c *gin.Context, status int, data interface{}) {
    c.JSON(status, APIResponse{Data: data})
}

func Error(c *gin.Context, status int, code, message string) {
    c.JSON(status, APIResponse{Error: &APIError{Code: code, Message: message}})
}
```

**✅ Verificación:**
- `cd backend && go build ./...`
- Levantar servidor, probar con curl:
  ```bash
  curl -X POST localhost:8080/api/v1/auth/register \
    -H 'Content-Type: application/json' \
    -d '{"email":"test@test.com","password":"12345678","full_name":"Test User"}'
  curl -X POST localhost:8080/api/v1/auth/login \
    -H 'Content-Type: application/json' \
    -d '{"email":"test@test.com","password":"12345678"}'
  ```

---

### Fase 7: Frontend Base

**🎯 Objetivo**: Inicializar frontend con Vite + Tailwind + sistema de temas + componentes plantilla base.

**📁 Archivos a crear:**
- `frontend/src/themes/light.css`
- `frontend/src/themes/dark.css`
- `frontend/src/themes/theme.css` (variables base)
- `frontend/src/contexts/ThemeContext.tsx`
- `frontend/src/components/ui/Button.tsx`
- `frontend/src/components/ui/Input.tsx`
- `frontend/src/components/ui/Card.tsx`
- `frontend/src/components/ui/Modal.tsx`
- `frontend/src/components/ui/Table.tsx`
- `frontend/src/components/ui/Badge.tsx`
- `frontend/src/components/ui/Select.tsx`
- `frontend/src/components/ui/Loading.tsx`
- `frontend/src/components/ui/EmptyState.tsx`
- `frontend/src/components/layout/Sidebar.tsx`
- `frontend/src/components/layout/Header.tsx`
- `frontend/src/components/layout/PageContainer.tsx`
- `frontend/src/components/layout/AppLayout.tsx`
- `frontend/src/utils/cn.ts` (utilidad para clases condicionales)

**🧩 Theme system:**

```tsx
// contexts/ThemeContext.tsx
import { createContext, useContext, useEffect, useState, ReactNode } from 'react';

type Theme = 'light' | 'dark';

interface ThemeContextType {
  theme: Theme;
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType>({
  theme: 'light',
  toggleTheme: () => {},
});

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setTheme] = useState<Theme>(() => {
    const saved = localStorage.getItem('theme') as Theme;
    if (saved) return saved;
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  });

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
  }, [theme]);

  const toggleTheme = () => setTheme(prev => prev === 'light' ? 'dark' : 'light');

  return (
    <ThemeContext.Provider value={{ theme, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export const useTheme = () => useContext(ThemeContext);
```

```css
/* themes/theme.css */
:root {
  --color-bg-primary: #ffffff;
  --color-bg-secondary: #f1f5f9;
  --color-bg-tertiary: #e2e8f0;
  --color-text-primary: #0f172a;
  --color-text-secondary: #64748b;
  --color-text-muted: #94a3b8;
  --color-accent: #4361ee;
  --color-accent-hover: #3651d4;
  --color-success: #10b981;
  --color-warning: #f59e0b;
  --color-danger: #ef4444;
  --color-info: #3b82f6;
  --color-border: #e2e8f0;
  --color-sidebar-bg: #1e293b;
  --color-sidebar-text: #cbd5e1;
  --color-sidebar-active: #4361ee;
  --font-family: 'Inter', system-ui, -apple-system, sans-serif;
  --font-size-xs: 0.75rem;
  --font-size-sm: 0.875rem;
  --font-size-base: 1rem;
  --font-size-lg: 1.125rem;
  --font-size-xl: 1.25rem;
  --font-size-2xl: 1.5rem;
  --font-size-3xl: 1.875rem;
  --radius-sm: 0.375rem;
  --radius-md: 0.5rem;
  --radius-lg: 0.75rem;
  --radius-xl: 1rem;
  --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
  --shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1);
  --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1);
  --spacing-1: 0.25rem;
  --spacing-2: 0.5rem;
  --spacing-3: 0.75rem;
  --spacing-4: 1rem;
  --spacing-6: 1.5rem;
  --spacing-8: 2rem;
  --sidebar-width: 260px;
  --header-height: 64px;
}

[data-theme="dark"] {
  --color-bg-primary: #0f172a;
  --color-bg-secondary: #1e293b;
  --color-bg-tertiary: #334155;
  --color-text-primary: #f1f5f9;
  --color-text-secondary: #94a3b8;
  --color-text-muted: #64748b;
  --color-accent: #6366f1;
  --color-accent-hover: #818cf8;
  --color-border: #334155;
  --color-sidebar-bg: #020617;
  --color-sidebar-text: #94a3b8;
  --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.3);
  --shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.4);
  --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.4);
}
```

```tsx
// utils/cn.ts
export function cn(...classes: (string | boolean | undefined | null)[]): string {
  return classes.filter(Boolean).join(' ');
}
```

```tsx
// components/ui/Button.tsx
import { ButtonHTMLAttributes, ReactNode } from 'react';
import { cn } from '../../utils/cn';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  icon?: ReactNode;
}

export function Button({
  variant = 'primary',
  size = 'md',
  loading = false,
  icon,
  children,
  className,
  disabled,
  ...props
}: ButtonProps) {
  const baseStyles = 'inline-flex items-center justify-center font-medium rounded-[var(--radius-md)] transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed';

  const variants = {
    primary: 'bg-[var(--color-accent)] text-white hover:bg-[var(--color-accent-hover)] focus:ring-[var(--color-accent)]',
    secondary: 'bg-[var(--color-bg-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] hover:bg-[var(--color-bg-tertiary)] focus:ring-[var(--color-accent)]',
    ghost: 'text-[var(--color-text-secondary)] hover:bg-[var(--color-bg-secondary)] hover:text-[var(--color-text-primary)] focus:ring-[var(--color-accent)]',
    danger: 'bg-[var(--color-danger)] text-white hover:opacity-90 focus:ring-[var(--color-danger)]',
  };

  const sizes = {
    sm: 'px-3 py-1.5 text-[var(--font-size-sm)] gap-1.5',
    md: 'px-4 py-2 text-[var(--font-size-base)] gap-2',
    lg: 'px-6 py-3 text-[var(--font-size-lg)] gap-2',
  };

  return (
    <button
      className={cn(baseStyles, variants[variant], sizes[size], className)}
      disabled={disabled || loading}
      {...props}
    >
      {loading ? (
        <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24" fill="none">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
      ) : icon ? (
        <span className="flex-shrink-0">{icon}</span>
      ) : null}
      {children}
    </button>
  );
}
```

**✅ Verificación:**
- `cd frontend && npm run build` → sin errores
- `cd frontend && npm run dev` → ver componentes en navegador

---

### Fase 8: Frontend — Auth + API Layer

**🎯 Objetivo**: Implementar cliente API con axios, contextos de autenticación, rutas protegidas y páginas de login/register.

**📁 Archivos a crear:**
- `frontend/src/services/api.ts` (axios instance + interceptors)
- `frontend/src/contexts/AuthContext.tsx`
- `frontend/src/hooks/useAuth.ts`
- `frontend/src/services/auth.ts` (API functions)
- `frontend/src/pages/auth/LoginPage.tsx`
- `frontend/src/pages/auth/RegisterPage.tsx`
- `frontend/src/components/layout/PrivateRoute.tsx`
- `frontend/src/components/layout/AdminRoute.tsx`
- `frontend/src/types/auth.ts`

**🧩 API Client con interceptor de refresh:**

```tsx
// services/api.ts
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  headers: { 'Content-Type': 'application/json' },
});

let accessToken: string | null = null;

export function setAccessToken(token: string | null) {
  accessToken = token;
}

export function getAccessToken(): string | null {
  return accessToken;
}

// Request interceptor: añade token
api.interceptors.request.use((config) => {
  if (accessToken) {
    config.headers.Authorization = `Bearer ${accessToken}`;
  }
  return config;
});

// Response interceptor: maneja 401 y refresh
let isRefreshing = false;
let failedQueue: Array<{ resolve: (token: string) => void; reject: (err: any) => void }> = [];

const processQueue = (error: any, token: string | null = null) => {
  failedQueue.forEach(prom => {
    if (error) prom.reject(error);
    else prom.resolve(token!);
  });
  failedQueue = [];
};

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(token => {
          originalRequest.headers.Authorization = `Bearer ${token}`;
          return api(originalRequest);
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const refreshToken = localStorage.getItem('refreshToken');
        if (!refreshToken) throw new Error('No refresh token');

        const { data } = await axios.post(
          `${import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'}/auth/refresh`,
          { refresh_token: refreshToken }
        );

        setAccessToken(data.data.access_token);
        processQueue(null, data.data.access_token);
        originalRequest.headers.Authorization = `Bearer ${data.data.access_token}`;
        return api(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        setAccessToken(null);
        localStorage.removeItem('refreshToken');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

export default api;
```

```tsx
// contexts/AuthContext.tsx
import { createContext, useState, useEffect, ReactNode } from 'react';
import api, { setAccessToken } from '../services/api';

interface User {
  id: string;
  email: string;
  full_name: string;
  role: string;
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, fullName: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
  isAdmin: boolean;
}

export const AuthContext = createContext<AuthContextType>({} as AuthContextType);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const refreshToken = localStorage.getItem('refreshToken');
    if (refreshToken) {
      api.post('/auth/refresh', { refresh_token: refreshToken })
        .then(({ data }) => {
          setAccessToken(data.data.access_token);
          setUser(data.data.user);
        })
        .catch(() => {
          localStorage.removeItem('refreshToken');
        })
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  const login = async (email: string, password: string) => {
    const { data } = await api.post('/auth/login', { email, password });
    setAccessToken(data.data.access_token);
    localStorage.setItem('refreshToken', data.data.refresh_token);
    setUser(data.data.user);
  };

  const register = async (email: string, password: string, fullName: string) => {
    const { data } = await api.post('/auth/register', { email, password, full_name: fullName });
    setAccessToken(data.data.access_token);
    localStorage.setItem('refreshToken', data.data.refresh_token);
    setUser(data.data.user);
  };

  const logout = () => {
    setAccessToken(null);
    localStorage.removeItem('refreshToken');
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{
      user,
      loading,
      login,
      register,
      logout,
      isAuthenticated: !!user,
      isAdmin: user?.role === 'super_admin',
    }}>
      {children}
    </AuthContext.Provider>
  );
}
```

**✅ Verificación:**
- `cd frontend && npm run build` → sin errores
- Probar flujo completo: register → login → refresh → logout en navegador

---

### Fase 9: Frontend — Funcionalidades

**🎯 Objetivo**: Implementar todas las páginas funcionales.

**📁 Archivos a crear:**

**Dashboard (`/dashboard`):**
- `frontend/src/pages/dashboard/DashboardPage.tsx`
- Cards de resumen: Total manejado, Disponible, Reservado
- Gráfica de balance de últimos 6 meses (recharts)
- Lista de próximos pagos
- Hook: `useDashboard`

**Deudas (`/debts`):**
- `frontend/src/pages/debts/DebtsListPage.tsx`
- `frontend/src/pages/debts/DebtDetailPage.tsx`
- `frontend/src/pages/debts/DebtForm.tsx` (modal)
- `frontend/src/services/debts.ts`
- Tabla de deudas con estado, monto restante, progreso
- Vista mensual: qué se pagó y qué falta
- CRUD completo con modales

**Gastos Recurrentes (`/expenses`):**
- `frontend/src/pages/expenses/ExpensesPage.tsx`
- `frontend/src/services/expenses.ts`
- Lista de suscripciones/gastos fijos
- Toggle pagado/no pagado por mes
- Proyección mensual de gastos

**Contabilidad (`/accounting`):**
- `frontend/src/pages/accounting/AccountingPage.tsx`
- `frontend/src/services/accounting.ts`
- Tabla de transacciones con filtro por mes
- Formulario para agregar ingreso/gasto
- Balance mensual
- Flujo de caja semanal/mensual/anual

**Planificación (`/planning`):**
- `frontend/src/pages/planning/PlanningPage.tsx`
- Estimaciones de ingresos
- Proyecciones basadas en gastos recurrentes + deudas

**Admin (`/admin`):**
- `frontend/src/pages/admin/AdminPage.tsx`
- `frontend/src/services/admin.ts`
- Tabla de usuarios
- Modal para crear/editar usuario
- Botón para desactivar/activar
- Solo visible si role = super_admin

**Main App:**
- `frontend/src/App.tsx` (router setup)
- `frontend/src/main.tsx`

```tsx
// App.tsx (router)
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from './contexts/ThemeContext';
import { AuthProvider } from './contexts/AuthContext';
import { AppLayout } from './components/layout/AppLayout';
import { PrivateRoute } from './components/layout/PrivateRoute';
import { AdminRoute } from './components/layout/AdminRoute';
import { LoginPage } from './pages/auth/LoginPage';
import { RegisterPage } from './pages/auth/RegisterPage';
import { DashboardPage } from './pages/dashboard/DashboardPage';
import { DebtsListPage } from './pages/debts/DebtsListPage';
import { ExpensesPage } from './pages/expenses/ExpensesPage';
import { AccountingPage } from './pages/accounting/AccountingPage';
import { PlanningPage } from './pages/planning/PlanningPage';
import { AdminPage } from './pages/admin/AdminPage';

export default function App() {
  return (
    <BrowserRouter>
      <ThemeProvider>
        <AuthProvider>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route element={<PrivateRoute />}>
              <Route element={<AppLayout />}>
                <Route path="/" element={<Navigate to="/dashboard" replace />} />
                <Route path="/dashboard" element={<DashboardPage />} />
                <Route path="/debts" element={<DebtsListPage />} />
                <Route path="/expenses" element={<ExpensesPage />} />
                <Route path="/accounting" element={<AccountingPage />} />
                <Route path="/planning" element={<PlanningPage />} />
                <Route element={<AdminRoute />}>
                  <Route path="/admin" element={<AdminPage />} />
                </Route>
              </Route>
            </Route>
          </Routes>
        </AuthProvider>
      </ThemeProvider>
    </BrowserRouter>
  );
}
```

**✅ Verificación:**
- `cd frontend && npm run build` → sin errores
- Navegar por todas las rutas en desarrollo
- Probar CRUD de deudas, toggle de gastos, registro de transacciones

---

### Fase 10: Dockerización

**🎯 Objetivo**: Crear Dockerfiles y docker-compose para los 3 servicios.

**📁 Archivos a crear/modificar:**
- `docker/backend/Dockerfile`
- `docker/frontend/Dockerfile`
- `docker/docker-compose.yml`
- `.env.example` (completo)

```dockerfile
# docker/backend/Dockerfile
# Stage 1: Build
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# Stage 2: Runtime
FROM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/internal/infrastructure/persistence/migrations ./migrations
EXPOSE 8080
CMD ["./server"]
```

```dockerfile
# docker/frontend/Dockerfile
# Stage 1: Build
FROM node:22-alpine AS builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Stage 2: Serve with nginx
FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY docker/frontend/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

```nginx
# docker/frontend/nginx.conf
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

```yaml
# docker/docker-compose.yml
version: '3.9'

services:
  db:
    image: postgres:16-alpine
    container_name: myinquisitor-db
    environment:
      POSTGRES_DB: ${POSTGRES_DB:-myinquisitor}
      POSTGRES_USER: ${POSTGRES_USER:-myinquisitor}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-myinquisitor_secret}
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-myinquisitor}"]
      interval: 5s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: ..
      dockerfile: docker/backend/Dockerfile
    container_name: myinquisitor-backend
    environment:
      SERVER_PORT: ${SERVER_PORT:-8080}
      DATABASE_URL: postgres://${POSTGRES_USER:-myinquisitor}:${POSTGRES_PASSWORD:-myinquisitor_secret}@db:5432/${POSTGRES_DB:-myinquisitor}?sslmode=disable
      JWT_SECRET: ${JWT_SECRET:-super-secret-jwt-key-change-in-production}
      ENCRYPTION_KEY: ${ENCRYPTION_KEY:-0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef}
      ENVIRONMENT: ${ENVIRONMENT:-development}
      ALLOWED_ORIGINS: ${ALLOWED_ORIGINS:-http://localhost:5173,http://localhost}
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped

  frontend:
    build:
      context: ..
      dockerfile: docker/frontend/Dockerfile
    container_name: myinquisitor-frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  pgdata:
```

**✅ Verificación:**
```bash
cd docker && docker-compose up --build
curl http://localhost:8080/health
curl http://localhost/api/v1/auth/login
```

---

### Fase 11: Testing Integral

**🎯 Objetivo**: Tests de integración backend con testcontainers, tests frontend con Vitest.

**📁 Archivos a crear:**
- `backend/internal/infrastructure/persistence/user_repo_test.go`
- `backend/internal/infrastructure/persistence/debt_repo_test.go`
- `backend/internal/infrastructure/api/handler/auth_handler_test.go`
- `backend/internal/infrastructure/api/handler/debt_handler_test.go`
- `frontend/src/components/ui/Button.test.tsx`
- `frontend/src/contexts/AuthContext.test.tsx`
- `frontend/src/pages/auth/LoginPage.test.tsx`

**🧩 Patrón testcontainers:**

```go
// user_repo_test.go
package persistence_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    "github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
    ctx := context.Background()
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image: "postgres:16-alpine",
            Env: map[string]string{
                "POSTGRES_DB":       "testdb",
                "POSTGRES_USER":     "test",
                "POSTGRES_PASSWORD": "test",
            },
            ExposedPorts: []string{"5432/tcp"},
            WaitingFor:   wait.ForLog("database system is ready to accept connections"),
        },
        Started: true,
    })
    require.NoError(t, err)

    host, _ := container.Host(ctx)
    port, _ := container.MappedPort(ctx, "5432")

    dsn := "postgres://test:test@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
    pool, err := pgxpool.New(ctx, dsn)
    require.NoError(t, err)

    // Ejecutar migraciones
    // ... run migrate up ...

    cleanup := func() {
        pool.Close()
        container.Terminate(ctx)
    }

    return pool, cleanup
}

func TestUserRepository_Create(t *testing.T) {
    pool, cleanup := setupTestDB(t)
    defer cleanup()

    // ... test ...
}
```

**✅ Verificación:**
- `cd backend && go test ./... -v -count=1`
- `cd frontend && npx vitest run`
- Cobertura mínima: 80% en application layer

---

### Fase 12: Despliegue

**🎯 Objetivo**: Documentar y configurar el despliegue en producción.

**📁 Archivos a crear/modificar:**
- `docker/docker-compose.prod.yml`
- Scripts de backup/restore
- Health check endpoints

**Consideraciones de producción:**
- Cambiar todas las contraseñas/secretos en producción
- Usar redes separadas en Docker Compose
- Configurar logging estructurado
- Health checks en cada servicio
- Backup automático de la DB (cron + pg_dump)
- SSL/TLS (reverse proxy con Caddy o Traefik)

```yaml
# docker-compose.prod.yml
version: '3.9'

services:
  db:
    image: postgres:16-alpine
    restart: always
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./backups:/backups
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER}"]
    networks:
      - internal

  backend:
    build:
      context: ..
      dockerfile: docker/backend/Dockerfile
    restart: always
    environment:
      SERVER_PORT: "8080"
      DATABASE_URL: postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable
      JWT_SECRET: ${JWT_SECRET}
      ENCRYPTION_KEY: ${ENCRYPTION_KEY}
      ENVIRONMENT: production
      ALLOWED_ORIGINS: ${FRONTEND_URL}
    depends_on:
      db:
        condition: service_healthy
    networks:
      - internal
      - web

  frontend:
    build:
      context: ..
      dockerfile: docker/frontend/Dockerfile
    restart: always
    depends_on:
      - backend
    networks:
      - web

  # Opcional: reverse proxy para SSL
  caddy:
    image: caddy:2-alpine
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy_data:/data
    depends_on:
      - frontend
    networks:
      - web

networks:
  internal:
    driver: bridge
    internal: true
  web:
    driver: bridge

volumes:
  pgdata:
  caddy_data:
```

**✅ Verificación final:**
```bash
# Build y despliegue completo
cd docker && docker-compose -f docker-compose.prod.yml up --build -d

# Health checks
curl -f http://localhost/health
curl -f http://localhost/api/v1/auth/login -X POST -d '{"email":"admin@myinquisitor.com","password":"..."}'

# Logs
docker-compose logs -f
```

---

## 9. 📁 Variables de Entorno

```bash
# .env.example

# Server
SERVER_PORT=8080
ENVIRONMENT=development

# Database
POSTGRES_DB=myinquisitor
POSTGRES_USER=myinquisitor
POSTGRES_PASSWORD=myinquisitor_secret
DATABASE_URL=postgres://myinquisitor:myinquisitor_secret@localhost:5432/myinquisitor?sslmode=disable

# JWT
JWT_SECRET=change-this-to-a-random-64-char-string
JWT_ACCESS_EXPIRATION=15m
JWT_REFRESH_EXPIRATION=720h

# Encryption (64 hex chars = 32 bytes)
ENCRYPTION_KEY=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef

# CORS
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# Frontend (Vite)
VITE_API_URL=http://localhost:8080/api/v1
```

---

## 10. 📝 Notas Finales

1. **Cada fase debe producir código compilable y funcional.** Si algo queda broken, detente y arréglalo antes de continuar.
2. **Mantén el bajo acoplamiento.** No importes cosas de una feature dentro de otra. Si dos features necesitan compartir lógica, muévela a `domain/` o `components/shared/`.
3. **El sistema de temas es global.** Si necesitas un nuevo color, agrégalo a las variables CSS. Nunca uses colores hardcodeados en componentes.
4. **Seguridad: revisa siempre** que los datos sensibles estén cifrados antes de insertarlos y descifrados después de leerlos.
5. **No generes datos mock ni seeds automáticos** a menos que se indique explícitamente.
6. **Los mensajes de error de la API deben ser informativos pero sin revelar detalles internos.** Nunca expongas stack traces en producción.
