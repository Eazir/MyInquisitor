package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Category, error)
	ListByType(ctx context.Context, userID uuid.UUID, categoryType string) ([]entity.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
