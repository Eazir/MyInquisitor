package entity

import (
	"time"
	"github.com/google/uuid"
)

type InviteToken struct {
	ID        uuid.UUID
	Token     string
	CreatedBy uuid.UUID
	Used      bool
	ExpiresAt time.Time
	CreatedAt time.Time
}

type CreateInviteTokenInput struct {
	CreatedBy uuid.UUID
	ExpiresAt time.Time
}
