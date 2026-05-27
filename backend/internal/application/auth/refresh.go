package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/application/dto"
)

var ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")

type RefreshUseCase struct {
	userRepo repository.UserRepository
	token    TokenService
}

func NewRefreshUseCase(userRepo repository.UserRepository, token TokenService) *RefreshUseCase {
	return &RefreshUseCase{
		userRepo: userRepo,
		token:    token,
	}
}

func (uc *RefreshUseCase) Execute(ctx context.Context, input dto.RefreshInput) (*dto.AuthOutput, error) {
	claims, err := uc.token.ValidateToken(input.RefreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if !user.Active {
		return nil, ErrUserInactive
	}

	accessToken, err := uc.token.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := uc.token.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &dto.AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserDTO{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}, nil
}
