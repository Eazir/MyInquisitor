-- ============================================================
-- Migration 001: Initial schema for MyInquisitor
-- ============================================================

-- 1. USERS
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL,
    email_hash      TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    full_name       TEXT NOT NULL,
    phone           TEXT,
    role            TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'super_admin')),
    active          BOOLEAN NOT NULL DEFAULT true,
    super_admin     BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 2. CATEGORIES
CREATE TABLE categories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    type            TEXT NOT NULL CHECK (type IN ('income', 'expense', 'debt')),
    icon            TEXT,
    color           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 3. DEBTS
CREATE TABLE debts (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id         UUID REFERENCES categories(id),
    name                TEXT NOT NULL,
    description         TEXT,
    total_amount        NUMERIC(14,2) NOT NULL,
    remaining_amount    NUMERIC(14,2) NOT NULL,
    interest_rate       NUMERIC(5,2) DEFAULT 0,
    total_installments  INT NOT NULL DEFAULT 1,
    current_installment INT NOT NULL DEFAULT 1,
    status              TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paid', 'settled')),
    start_date          DATE NOT NULL,
    end_date            DATE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 4. DEBT_MONTHLY_STATUS
CREATE TABLE debt_monthly_status (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    debt_id             UUID NOT NULL REFERENCES debts(id) ON DELETE CASCADE,
    month               DATE NOT NULL,
    installment_num     INT NOT NULL,
    total_installments  INT NOT NULL,
    amount_due          NUMERIC(14,2) NOT NULL,
    interest_amount     NUMERIC(14,2) DEFAULT 0,
    principal_amount    NUMERIC(14,2) DEFAULT 0,
    amount_paid         NUMERIC(14,2) DEFAULT 0,
    paid                BOOLEAN NOT NULL DEFAULT false,
    paid_at             TIMESTAMPTZ,
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(debt_id, month)
);

-- 5. RECURRING_EXPENSES
CREATE TABLE recurring_expenses (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id     UUID REFERENCES categories(id),
    name            TEXT NOT NULL,
    description     TEXT,
    amount          NUMERIC(14,2) NOT NULL,
    frequency       TEXT NOT NULL CHECK (frequency IN ('monthly', 'yearly', 'weekly', 'biweekly')),
    due_day         INT,
    status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'cancelled')),
    start_date      DATE NOT NULL,
    end_date        DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 6. EXPENSE_MONTHLY_STATUS
CREATE TABLE expense_monthly_status (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expense_id      UUID NOT NULL REFERENCES recurring_expenses(id) ON DELETE CASCADE,
    month           DATE NOT NULL,
    paid            BOOLEAN NOT NULL DEFAULT false,
    paid_at         TIMESTAMPTZ,
    amount_paid     NUMERIC(14,2),
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(expense_id, month)
);

-- 7. TRANSACTIONS
CREATE TABLE transactions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id         UUID REFERENCES categories(id),
    type                TEXT NOT NULL CHECK (type IN ('income', 'expense', 'transfer')),
    amount              NUMERIC(14,2) NOT NULL,
    description         TEXT,
    source              TEXT,
    reference_date      DATE NOT NULL DEFAULT CURRENT_DATE,
    is_recurring        BOOLEAN NOT NULL DEFAULT false,
    recurring_expense_id UUID REFERENCES recurring_expenses(id),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 8. MONTHLY_SUMMARY
CREATE TABLE monthly_summary (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    month               DATE NOT NULL,
    total_income        NUMERIC(14,2) NOT NULL DEFAULT 0,
    income_breakdown    JSONB,
    total_expenses      NUMERIC(14,2) NOT NULL DEFAULT 0,
    expense_breakdown   JSONB,
    total_debt_payments NUMERIC(14,2) NOT NULL DEFAULT 0,
    debt_breakdown      JSONB,
    total_obligations   NUMERIC(14,2) NOT NULL DEFAULT 0,
    net_balance         NUMERIC(14,2) NOT NULL DEFAULT 0,
    projected_income    NUMERIC(14,2),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, month)
);

-- ============================================================
-- ÍNDICES
-- ============================================================
CREATE INDEX idx_debts_user_id ON debts(user_id);
CREATE INDEX idx_debts_status ON debts(status);
CREATE INDEX idx_debt_monthly_status_debt_id ON debt_monthly_status(debt_id);
CREATE INDEX idx_debt_monthly_status_month ON debt_monthly_status(month);
CREATE INDEX idx_recurring_expenses_user_id ON recurring_expenses(user_id);
CREATE INDEX idx_expense_monthly_status_expense_id ON expense_monthly_status(expense_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_reference_date ON transactions(reference_date);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_monthly_summary_user_id ON monthly_summary(user_id);
CREATE INDEX idx_monthly_summary_month ON monthly_summary(month);
CREATE INDEX idx_categories_user_id ON categories(user_id);
