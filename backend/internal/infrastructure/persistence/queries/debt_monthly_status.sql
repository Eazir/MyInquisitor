-- name: CreateDebtMonthlyStatus :one
INSERT INTO debt_monthly_status (debt_id, month, installment_num, total_installments,
                                 amount_due, interest_amount, principal_amount)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetDebtMonthlyStatus :one
SELECT * FROM debt_monthly_status WHERE debt_id = $1 AND month = $2;

-- name: ListDebtMonthlyStatusByDebtID :many
SELECT * FROM debt_monthly_status WHERE debt_id = $1 ORDER BY month DESC;

-- name: ListDebtMonthlyStatusByMonth :many
SELECT * FROM debt_monthly_status WHERE month = $1 ORDER BY debt_id;

-- name: MarkDebtMonthAsPaid :one
UPDATE debt_monthly_status
SET paid = true, paid_at = now(), amount_paid = $3, notes = COALESCE($4, notes), updated_at = now()
WHERE debt_id = $1 AND month = $2
RETURNING *;

-- name: UpdateDebtMonthlyStatus :one
UPDATE debt_monthly_status
SET amount_paid = COALESCE($3, amount_paid),
    paid = $4,
    notes = COALESCE($5, notes),
    updated_at = now()
WHERE debt_id = $1 AND month = $2
RETURNING *;

-- name: DeleteDebtMonthlyStatus :exec
DELETE FROM debt_monthly_status WHERE id = $1;
