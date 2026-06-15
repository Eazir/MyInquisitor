package admin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/application/dto"
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

type ListInvitesUseCase struct {
	inviteRepo repository.InviteTokenRepository
}

func NewListInvitesUseCase(inviteRepo repository.InviteTokenRepository) *ListInvitesUseCase {
	return &ListInvitesUseCase{inviteRepo: inviteRepo}
}

func (uc *ListInvitesUseCase) Execute(ctx context.Context) ([]dto.InviteTokenOutput, error) {
	items, err := uc.inviteRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list invites: %w", err)
	}

	now := time.Now()
	output := make([]dto.InviteTokenOutput, len(items))
	for i, item := range items {
		output[i] = dto.InviteTokenOutput{
			ID:          item.InviteToken.ID,
			Token:       item.InviteToken.Token,
			CreatedBy:   item.InviteToken.CreatedBy,
			CreatorName: item.CreatorName,
			ExpiresAt:   item.InviteToken.ExpiresAt,
			URL:         "/register/" + item.InviteToken.Token,
			Used:        item.InviteToken.Used,
			CreatedAt:   item.InviteToken.CreatedAt,
			Expired:     now.After(item.InviteToken.ExpiresAt),
		}
	}
	return output, nil
}

type DeleteInviteUseCase struct {
	inviteRepo repository.InviteTokenRepository
}

func NewDeleteInviteUseCase(inviteRepo repository.InviteTokenRepository) *DeleteInviteUseCase {
	return &DeleteInviteUseCase{inviteRepo: inviteRepo}
}

func (uc *DeleteInviteUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	if err := uc.inviteRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("invite token not found: %w", err)
		}
		return fmt.Errorf("delete invite token: %w", err)
	}
	return nil
}
