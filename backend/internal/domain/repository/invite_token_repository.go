package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
)

type InviteTokenWithCreator struct {
	InviteToken entity.InviteToken
	CreatorName *string
}

type InviteTokenRepository interface {
	Create(ctx context.Context, token *entity.InviteToken) error
	GetByToken(ctx context.Context, token string) (*entity.InviteToken, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.InviteToken, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
	ListAll(ctx context.Context) ([]InviteTokenWithCreator, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
