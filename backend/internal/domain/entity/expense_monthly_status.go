package entity

import (
	"time"
	"github.com/google/uuid"
)

type ExpenseMonthlyStatus struct {
	ID        uuid.UUID
	ExpenseID uuid.UUID
	Month     time.Time
	Paid      bool
	PaidAt    *time.Time
	AmountPaid *float64
	Notes     *string
	CreatedAt time.Time
}
