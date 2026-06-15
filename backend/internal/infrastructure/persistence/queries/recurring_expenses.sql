-- name: CreateRecurringExpense :one
INSERT INTO recurring_expenses (user_id, category_id, name, description, amount,
                                frequency, due_day, billing_month, status, start_date, end_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetRecurringExpenseByID :one
SELECT * FROM recurring_expenses WHERE id = $1;

-- name: ListRecurringExpensesByUserID :many
SELECT * FROM recurring_expenses WHERE user_id = $1 ORDER BY name;

-- name: ListActiveRecurringExpensesByUserID :many
SELECT * FROM recurring_expenses WHERE user_id = $1 AND status = 'active' ORDER BY name;

-- name: UpdateRecurringExpense :one
UPDATE recurring_expenses
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    amount = COALESCE($4, amount),
    frequency = COALESCE($5, frequency),
    due_day = COALESCE($6, due_day),
    billing_month = COALESCE($7, billing_month),
    status = COALESCE($8, status),
    end_date = COALESCE($9, end_date),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteRecurringExpense :exec
DELETE FROM recurring_expenses WHERE id = $1;
