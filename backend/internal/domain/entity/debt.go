package entity

import (
	"time"
	"github.com/google/uuid"
)

type Debt struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	CategoryID         *uuid.UUID
	Name               string
	Description        *string
	TotalAmount        float64
	RemainingAmount    float64
	InterestRate       float64
	TotalInstallments  int32
	CurrentInstallment int32
	Status             string
	StartDate          time.Time
	EndDate            *time.Time
	DueDay             *int32
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type CreateDebtInput struct {
	UserID            uuid.UUID
	CategoryID        *uuid.UUID
	Name              string
	Description       *string
	TotalAmount       float64
	InterestRate      float64
	TotalInstallments int32
	StartDate         time.Time
	EndDate           *time.Time
	DueDay             *int32
}
