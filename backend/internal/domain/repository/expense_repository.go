package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
)

type RecurringExpenseRepository interface {
	Create(ctx context.Context, expense *entity.RecurringExpense) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.RecurringExpense, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.RecurringExpense, error)
	ListActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entity.RecurringExpense, error)
	Update(ctx context.Context, expense *entity.RecurringExpense) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ExpenseMonthlyStatusRepository interface {
	Create(ctx context.Context, status *entity.ExpenseMonthlyStatus) error
	GetByExpenseIDAndMonth(ctx context.Context, expenseID uuid.UUID, month string) (*entity.ExpenseMonthlyStatus, error)
	ListByExpenseID(ctx context.Context, expenseID uuid.UUID) ([]entity.ExpenseMonthlyStatus, error)
	Upsert(ctx context.Context, status *entity.ExpenseMonthlyStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}
