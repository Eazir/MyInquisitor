package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTransactionInput struct {
	CategoryID         *uuid.UUID `json:"category_id"`
	Type               string     `json:"type" validate:"required,oneof=income expense transfer"`
	Amount             float64    `json:"amount" validate:"required,gt=0"`
	Description        *string    `json:"description,omitempty"`
	Source             *string    `json:"source,omitempty"`
	ReferenceDate      string     `json:"reference_date" validate:"required"`
	RecurringExpenseID *uuid.UUID `json:"recurring_expense_id,omitempty"`
}

type UpdateTransactionInput struct {
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
	Type          *string    `json:"type,omitempty"`
	Amount        *float64   `json:"amount,omitempty"`
	Description   *string    `json:"description,omitempty"`
	Source        *string    `json:"source,omitempty"`
	ReferenceDate *string    `json:"reference_date,omitempty"`
}

type TransactionOutput struct {
	ID                 uuid.UUID  `json:"id"`
	UserID             uuid.UUID  `json:"user_id"`
	CategoryID         *uuid.UUID `json:"category_id"`
	Type               string     `json:"type"`
	Amount             float64    `json:"amount"`
	Description        *string    `json:"description"`
	Source             *string    `json:"source"`
	ReferenceDate      string     `json:"reference_date"`
	IsRecurring        bool       `json:"is_recurring"`
	RecurringExpenseID *uuid.UUID `json:"recurring_expense_id"`
	CreatedAt          time.Time  `json:"created_at"`
}

type MonthlyBalanceOutput struct {
	Month            string             `json:"month"`
	TotalIncome      float64            `json:"total_income"`
	TotalExpenses    float64            `json:"total_expenses"`
	TotalDebtPayments float64           `json:"total_debt_payments"`
	TotalObligations float64            `json:"total_obligations"`
	NetBalance       float64            `json:"net_balance"`
	ProjectedIncome  *float64           `json:"projected_income"`
	IncomeBreakdown  map[string]float64 `json:"income_breakdown,omitempty"`
	ExpenseBreakdown map[string]float64 `json:"expense_breakdown,omitempty"`
}

type CashFlowOutput struct {
	Period   string               `json:"period"`
	Entries  []CashFlowEntry       `json:"entries"`
	TotalIn  float64              `json:"total_in"`
	TotalOut float64              `json:"total_out"`
	Balance  float64              `json:"balance"`
}

type CashFlowEntry struct {
	Period string  `json:"period"`
	In     float64 `json:"in"`
	Out    float64 `json:"out"`
}

type ExtraExpenseItem struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

type FixedExpenseItem struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Amount    float64 `json:"amount"`
	Frequency string  `json:"frequency"`
	DueDay    *int32  `json:"due_day"`
}

type ProjectionOutput struct {
	Month              string             `json:"month"`
	MonthLabel         string             `json:"month_label"`
	BaseIncome         float64            `json:"base_income"`
	IncomeModifier     float64            `json:"income_modifier"`
	ProjectedIncome    float64            `json:"projected_income"`
	FixedExpenses      float64            `json:"fixed_expenses"`
	FixedExpensesList  []FixedExpenseItem `json:"fixed_expenses_list"`
	ExtraBudgetaryAvg  float64            `json:"extra_budgetary_avg"`
	ExtraExpensesTotal float64            `json:"extra_expenses_total"`
	ExtraExpensesList  []ExtraExpenseItem `json:"extra_expenses_list"`
	DebtPayments       float64            `json:"debt_payments"`
	TotalExpenses      float64            `json:"total_expenses"`
	ProjectedBalance   float64            `json:"projected_balance"`
}

type CreateCategoryInput struct {
	Name  string  `json:"name" validate:"required"`
	Type  string  `json:"type" validate:"required,oneof=income expense debt"`
	Icon  *string `json:"icon,omitempty"`
	Color *string `json:"color,omitempty"`
}

type CategoryOutput struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Icon      *string   `json:"icon"`
	Color     *string   `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}
