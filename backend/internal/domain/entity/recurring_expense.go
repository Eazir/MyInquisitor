package entity

import (
	"time"
	"github.com/google/uuid"
)

type RecurringExpense struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	CategoryID    *uuid.UUID
	Name          string
	Description   *string
	Amount        float64
	Frequency     string
	DueDay        *int32
	BillingMonth  *int32
	Status        string
	StartDate     time.Time
	EndDate       *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateRecurringExpenseInput struct {
	UserID       uuid.UUID
	CategoryID   *uuid.UUID
	Name         string
	Description  *string
	Amount       float64
	Frequency    string
	DueDay       *int32
	BillingMonth *int32
	StartDate    time.Time
	EndDate      *time.Time
}
