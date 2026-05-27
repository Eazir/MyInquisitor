package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *entity.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Transaction, error)
	ListByUserIDAndMonth(ctx context.Context, userID uuid.UUID, startMonth, endMonth string) ([]entity.Transaction, error)
	ListByDateRange(ctx context.Context, userID uuid.UUID, start, end string) ([]entity.Transaction, error)
	Update(ctx context.Context, tx *entity.Transaction) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type MonthlySummaryRepository interface {
	Upsert(ctx context.Context, summary *entity.MonthlySummary) error
	GetByUserIDAndMonth(ctx context.Context, userID uuid.UUID, month string) (*entity.MonthlySummary, error)
	ListByYear(ctx context.Context, userID uuid.UUID, startYear, endYear string) ([]entity.MonthlySummary, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
