package admin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
)

type GenerateInviteUseCase struct {
	inviteRepo repository.InviteTokenRepository
}

func NewGenerateInviteUseCase(inviteRepo repository.InviteTokenRepository) *GenerateInviteUseCase {
	return &GenerateInviteUseCase{inviteRepo: inviteRepo}
}

func (uc *GenerateInviteUseCase) Execute(ctx context.Context, adminID uuid.UUID) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	token := hex.EncodeToString(raw)

	invite := &entity.InviteToken{
		Token:     token,
		CreatedBy: adminID,
		ExpiresAt: time.Now().Add(72 * time.Hour),
	}

	if err := uc.inviteRepo.Create(ctx, invite); err != nil {
		return "", fmt.Errorf("save invite token: %w", err)
	}

	return token, nil
}
