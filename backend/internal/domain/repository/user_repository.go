package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
)

var ErrNotFound = errors.New("resource not found")

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	List(ctx context.Context, limit, offset int) ([]entity.User, int, error)
	Update(ctx context.Context, user *entity.User) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	SetActive(ctx context.Context, id uuid.UUID, active bool) error
	Delete(ctx context.Context, id uuid.UUID) error
}
