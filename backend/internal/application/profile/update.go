package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/application/dto"
	"github.com/myinquisitor/backend/internal/domain/repository"
)

type UpdateProfileUseCase struct {
	userRepo repository.UserRepository
}

func NewUpdateProfileUseCase(userRepo repository.UserRepository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{userRepo: userRepo}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID uuid.UUID, input dto.UpdateProfileInput) (*dto.UpdateProfileOutput, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if input.FullName != "" {
		user.FullName = input.FullName
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return &dto.UpdateProfileOutput{
		User: dto.UserDTO{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}, nil
}
