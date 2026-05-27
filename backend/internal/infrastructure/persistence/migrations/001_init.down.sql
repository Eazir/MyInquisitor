-- ============================================================
-- Migration 001: Rollback
-- ============================================================
DROP TABLE IF EXISTS monthly_summary CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS expense_monthly_status CASCADE;
DROP TABLE IF EXISTS recurring_expenses CASCADE;
DROP TABLE IF EXISTS debt_monthly_status CASCADE;
DROP TABLE IF EXISTS debts CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS users CASCADE;
