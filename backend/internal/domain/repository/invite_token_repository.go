package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
)

type InviteTokenRepository interface {
	Create(ctx context.Context, token *entity.InviteToken) error
	GetByToken(ctx context.Context, token string) (*entity.InviteToken, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
}
