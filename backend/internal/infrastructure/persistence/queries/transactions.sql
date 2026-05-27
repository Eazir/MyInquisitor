-- name: CreateTransaction :one
INSERT INTO transactions (user_id, category_id, type, amount, description, source,
                          reference_date, is_recurring, recurring_expense_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transactions WHERE id = $1;

-- name: ListTransactionsByUserID :many
SELECT * FROM transactions WHERE user_id = $1 ORDER BY reference_date DESC, created_at DESC;

-- name: ListTransactionsByUserIDAndMonth :many
SELECT * FROM transactions
WHERE user_id = $1
  AND reference_date >= $2
  AND reference_date < $3
ORDER BY reference_date DESC;

-- name: ListTransactionsByDateRange :many
SELECT * FROM transactions
WHERE user_id = $1
  AND reference_date >= $2
  AND reference_date <= $3
ORDER BY reference_date DESC;

-- name: ListTransactionsByType :many
SELECT * FROM transactions
WHERE user_id = $1 AND type = $2
ORDER BY reference_date DESC;

-- name: UpdateTransaction :one
UPDATE transactions
SET category_id = COALESCE($3, category_id),
    amount = COALESCE($4, amount),
    description = COALESCE($5, description),
    source = COALESCE($6, source),
    reference_date = COALESCE($7, reference_date)
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions WHERE id = $1 AND user_id = $2;
