package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/application/dto"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidInviteToken = errors.New("invalid or expired invite token")
	ErrInviteTokenUsed    = errors.New("invite token already used")
)

type RegisterUseCase struct {
	userRepo   repository.UserRepository
	inviteRepo repository.InviteTokenRepository
	password   PasswordService
	token      TokenService
}

func NewRegisterUseCase(userRepo repository.UserRepository, inviteRepo repository.InviteTokenRepository, password PasswordService, token TokenService) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:   userRepo,
		inviteRepo: inviteRepo,
		password:   password,
		token:      token,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input dto.RegisterInput) (*dto.AuthOutput, error) {
	invite, err := uc.inviteRepo.GetByToken(ctx, input.InviteToken)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidInviteToken
		}
		return nil, fmt.Errorf("validate invite token: %w", err)
	}

	if invite.Used {
		return nil, ErrInviteTokenUsed
	}

	if time.Now().After(invite.ExpiresAt) {
		return nil, ErrInvalidInviteToken
	}

	existing, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	passwordHash, err := uc.password.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &entity.User{
		Email:        input.Email,
		PasswordHash: passwordHash,
		FullName:     input.FullName,
		Phone:        input.Phone,
		Role:         "user",
		Active:       true,
		SuperAdmin:   false,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := uc.inviteRepo.MarkAsUsed(ctx, invite.ID); err != nil {
		return nil, fmt.Errorf("mark invite as used: %w", err)
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
