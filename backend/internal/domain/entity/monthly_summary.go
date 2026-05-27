package entity

import (
	"time"
	"github.com/google/uuid"
)

type MonthlySummary struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Month            time.Time
	TotalIncome      float64
	IncomeBreakdown  map[string]float64
	TotalExpenses    float64
	ExpenseBreakdown map[string]float64
	TotalDebtPayments float64
	DebtBreakdown    map[string]float64
	TotalObligations float64
	NetBalance       float64
	ProjectedIncome  *float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
