package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/application/dto"
)

var (
	ErrInvalidCredentials = errors.New("invalid password")
	ErrEmailNotFound      = errors.New("no account found with this email")
	ErrUserInactive       = errors.New("account is inactive")
	ErrAccountIssue       = errors.New("account error, please contact support")
)

type LoginUseCase struct {
	userRepo repository.UserRepository
	password PasswordService
	token    TokenService
}

func NewLoginUseCase(userRepo repository.UserRepository, password PasswordService, token TokenService) *LoginUseCase {
	return &LoginUseCase{
		userRepo: userRepo,
		password: password,
		token:    token,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input dto.LoginInput) (*dto.AuthOutput, error) {
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrEmailNotFound
		}
		if strings.Contains(err.Error(), "decrypt") {
			return nil, fmt.Errorf("%w: data decryption failed", ErrAccountIssue)
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if !user.Active {
		return nil, ErrUserInactive
	}

	if !uc.password.Verify(input.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
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
