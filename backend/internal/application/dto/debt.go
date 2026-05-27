package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDebtInput struct {
	CategoryID        *uuid.UUID `json:"category_id"`
	Name              string     `json:"name" validate:"required"`
	Description       *string    `json:"description,omitempty"`
	TotalAmount       float64    `json:"total_amount" validate:"required,gt=0"`
	InterestRate      float64    `json:"interest_rate"`
	TotalInstallments int32      `json:"total_installments" validate:"required,gt=0"`
	StartDate         string     `json:"start_date" validate:"required"`
	EndDate           *string    `json:"end_date,omitempty"`
}

type UpdateDebtInput struct {
	Name              *string   `json:"name,omitempty"`
	Description       *string   `json:"description,omitempty"`
	TotalAmount       *float64  `json:"total_amount,omitempty"`
	RemainingAmount   *float64  `json:"remaining_amount,omitempty"`
	InterestRate      *float64  `json:"interest_rate,omitempty"`
	TotalInstallments *int32    `json:"total_installments,omitempty"`
	Status            *string   `json:"status,omitempty"`
}

type DebtOutput struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	CategoryID         *uuid.UUID `json:"category_id"`
	Name               string     `json:"name"`
	Description        *string    `json:"description"`
	TotalAmount        float64    `json:"total_amount"`
	RemainingAmount    float64    `json:"remaining_amount"`
	InterestRate       float64    `json:"interest_rate"`
	TotalInstallments  int32      `json:"total_installments"`
	CurrentInstallment int32      `json:"current_installment"`
	Status             string     `json:"status"`
	StartDate          string     `json:"start_date"`
	EndDate            *string    `json:"end_date"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type DebtMonthlyStatusOutput struct {
	ID                uuid.UUID  `json:"id"`
	DebtID            uuid.UUID  `json:"debt_id"`
	Month             string     `json:"month"`
	InstallmentNum    int32      `json:"installment_num"`
	AmountDue         float64    `json:"amount_due"`
	InterestAmount    float64    `json:"interest_amount"`
	PrincipalAmount   float64    `json:"principal_amount"`
	AmountPaid        float64    `json:"amount_paid"`
	Paid              bool       `json:"paid"`
	PaidAt            *time.Time `json:"paid_at"`
	Notes             *string    `json:"notes"`
}

type MarkDebtPaidInput struct {
	AmountPaid float64 `json:"amount_paid" validate:"required,gt=0"`
	Notes      *string `json:"notes,omitempty"`
}
