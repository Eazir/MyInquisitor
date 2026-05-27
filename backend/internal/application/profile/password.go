package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	appAuth "github.com/myinquisitor/backend/internal/application/auth"
	"github.com/myinquisitor/backend/internal/application/dto"
	"github.com/myinquisitor/backend/internal/domain/repository"
)

var (
	ErrInvalidPassword = errors.New("current password is incorrect")
)

type ChangePasswordUseCase struct {
	userRepo repository.UserRepository
	password appAuth.PasswordService
}

func NewChangePasswordUseCase(userRepo repository.UserRepository, password appAuth.PasswordService) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{userRepo: userRepo, password: password}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, userID uuid.UUID, input dto.ChangePasswordInput) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if !uc.password.Verify(input.CurrentPassword, user.PasswordHash) {
		return ErrInvalidPassword
	}

	hash, err := uc.password.Hash(input.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := uc.userRepo.UpdatePassword(ctx, userID, hash); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}
