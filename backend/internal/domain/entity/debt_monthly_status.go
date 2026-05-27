package entity

import (
	"time"
	"github.com/google/uuid"
)

type DebtMonthlyStatus struct {
	ID                uuid.UUID
	DebtID            uuid.UUID
	Month             time.Time
	InstallmentNum    int32
	TotalInstallments int32
	AmountDue         float64
	InterestAmount    float64
	PrincipalAmount   float64
	AmountPaid        float64
	Paid              bool
	PaidAt            *time.Time
	Notes             *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
