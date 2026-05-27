package entity

import (
	"time"
	"github.com/google/uuid"
)

type Transaction struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	CategoryID         *uuid.UUID
	Type               string
	Amount             float64
	Description        *string
	Source             *string
	ReferenceDate      time.Time
	IsRecurring        bool
	RecurringExpenseID *uuid.UUID
	CreatedAt          time.Time
}

type CreateTransactionInput struct {
	UserID             uuid.UUID
	CategoryID         *uuid.UUID
	Type               string
	Amount             float64
	Description        *string
	Source             *string
	ReferenceDate      time.Time
	RecurringExpenseID *uuid.UUID
}
