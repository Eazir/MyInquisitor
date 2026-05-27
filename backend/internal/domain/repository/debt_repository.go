package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
)

type DebtRepository interface {
	Create(ctx context.Context, debt *entity.Debt) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Debt, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Debt, error)
	ListActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Debt, error)
	Update(ctx context.Context, debt *entity.Debt) error
	UpdateCurrentInstallment(ctx context.Context, id uuid.UUID, installment int32) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type DebtMonthlyStatusRepository interface {
	Create(ctx context.Context, status *entity.DebtMonthlyStatus) error
	GetByDebtIDAndMonth(ctx context.Context, debtID uuid.UUID, month string) (*entity.DebtMonthlyStatus, error)
	ListByDebtID(ctx context.Context, debtID uuid.UUID) ([]entity.DebtMonthlyStatus, error)
	MarkAsPaid(ctx context.Context, debtID uuid.UUID, month string, amountPaid float64, notes *string) error
	Update(ctx context.Context, status *entity.DebtMonthlyStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}
