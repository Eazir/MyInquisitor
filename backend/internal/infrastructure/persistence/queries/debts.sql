-- name: CreateDebt :one
INSERT INTO debts (user_id, category_id, name, description, total_amount, remaining_amount,
                   interest_rate, total_installments, current_installment, status, start_date, end_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetDebtByID :one
SELECT * FROM debts WHERE id = $1;

-- name: ListDebtsByUserID :many
SELECT * FROM debts WHERE user_id = $1 ORDER BY created_at DESC;

-- name: ListActiveDebtsByUserID :many
SELECT * FROM debts WHERE user_id = $1 AND status = 'active' ORDER BY created_at DESC;

-- name: UpdateDebt :one
UPDATE debts
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    total_amount = COALESCE($4, total_amount),
    remaining_amount = COALESCE($5, remaining_amount),
    interest_rate = COALESCE($6, interest_rate),
    total_installments = COALESCE($7, total_installments),
    current_installment = COALESCE($8, current_installment),
    status = COALESCE($9, status),
    end_date = COALESCE($10, end_date),
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateDebtCurrentInstallment :one
UPDATE debts SET current_installment = $2, updated_at = now() WHERE id = $1 RETURNING *;

-- name: DeleteDebt :exec
DELETE FROM debts WHERE id = $1;
