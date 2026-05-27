package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateExpenseInput struct {
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        string     `json:"name" validate:"required"`
	Description *string    `json:"description,omitempty"`
	Amount      float64    `json:"amount" validate:"required,gt=0"`
	Frequency   string     `json:"frequency" validate:"required,oneof=monthly yearly weekly biweekly"`
	DueDay      *int32     `json:"due_day,omitempty"`
	StartDate   string     `json:"start_date" validate:"required"`
	EndDate     *string    `json:"end_date,omitempty"`
}

type UpdateExpenseInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Amount      *float64 `json:"amount,omitempty"`
	Frequency   *string  `json:"frequency,omitempty"`
	DueDay      *int32   `json:"due_day,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

type ExpenseOutput struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Amount      float64    `json:"amount"`
	Frequency   string     `json:"frequency"`
	DueDay      *int32     `json:"due_day"`
	Status      string     `json:"status"`
	StartDate   string     `json:"start_date"`
	EndDate     *string    `json:"end_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ToggleExpensePaidInput struct {
	AmountPaid *float64 `json:"amount_paid,omitempty"`
	Notes      *string  `json:"notes,omitempty"`
}

type ExpenseMonthlyStatusOutput struct {
	ID        uuid.UUID  `json:"id"`
	ExpenseID uuid.UUID  `json:"expense_id"`
	Month     string     `json:"month"`
	Paid      bool       `json:"paid"`
	PaidAt    *time.Time `json:"paid_at"`
	AmountPaid *float64  `json:"amount_paid"`
	Notes     *string    `json:"notes"`
}
