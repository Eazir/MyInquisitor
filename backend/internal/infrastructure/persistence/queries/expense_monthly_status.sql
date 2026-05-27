-- name: CreateExpenseMonthlyStatus :one
INSERT INTO expense_monthly_status (expense_id, month, amount_paid, notes)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetExpenseMonthlyStatus :one
SELECT * FROM expense_monthly_status WHERE expense_id = $1 AND month = $2;

-- name: ListExpenseMonthlyStatusByExpenseID :many
SELECT * FROM expense_monthly_status WHERE expense_id = $1 ORDER BY month DESC;

-- name: ListExpenseMonthlyStatusByMonth :many
SELECT * FROM expense_monthly_status WHERE month = $1 ORDER BY expense_id;

-- name: MarkExpenseMonthAsPaid :one
UPDATE expense_monthly_status
SET paid = true, paid_at = now(), amount_paid = COALESCE($3, amount_paid), notes = COALESCE($4, notes)
WHERE expense_id = $1 AND month = $2
RETURNING *;

-- name: UpsertExpenseMonthlyStatus :one
INSERT INTO expense_monthly_status (expense_id, month, paid, amount_paid, notes)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (expense_id, month)
DO UPDATE SET paid = COALESCE($3, expense_monthly_status.paid),
              amount_paid = COALESCE($4, expense_monthly_status.amount_paid),
              paid_at = CASE WHEN $3 THEN now() ELSE expense_monthly_status.paid_at END,
              notes = COALESCE($5, expense_monthly_status.notes)
RETURNING *;

-- name: DeleteExpenseMonthlyStatus :exec
DELETE FROM expense_monthly_status WHERE id = $1;
