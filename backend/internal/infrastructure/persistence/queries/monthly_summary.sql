-- name: UpsertMonthlySummary :one
INSERT INTO monthly_summary (user_id, month, total_income, income_breakdown,
                             total_expenses, expense_breakdown,
                             total_debt_payments, debt_breakdown,
                             total_obligations, net_balance, projected_income)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (user_id, month)
DO UPDATE SET total_income = COALESCE($3, monthly_summary.total_income),
              income_breakdown = COALESCE($4, monthly_summary.income_breakdown),
              total_expenses = COALESCE($5, monthly_summary.total_expenses),
              expense_breakdown = COALESCE($6, monthly_summary.expense_breakdown),
              total_debt_payments = COALESCE($7, monthly_summary.total_debt_payments),
              debt_breakdown = COALESCE($8, monthly_summary.debt_breakdown),
              total_obligations = COALESCE($9, monthly_summary.total_obligations),
              net_balance = COALESCE($10, monthly_summary.net_balance),
              projected_income = COALESCE($11, monthly_summary.projected_income),
              updated_at = now()
RETURNING *;

-- name: GetMonthlySummary :one
SELECT * FROM monthly_summary WHERE user_id = $1 AND month = $2;

-- name: ListMonthlySummariesByYear :many
SELECT * FROM monthly_summary
WHERE user_id = $1
  AND month >= $2
  AND month < $3
ORDER BY month DESC;

-- name: DeleteMonthlySummary :exec
DELETE FROM monthly_summary WHERE id = $1;
