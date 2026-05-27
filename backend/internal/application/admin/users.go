package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/myinquisitor/backend/internal/domain/entity"
	"github.com/myinquisitor/backend/internal/domain/repository"
	"github.com/myinquisitor/backend/internal/application/auth"
	"github.com/myinquisitor/backend/internal/application/dto"
)

type ListUsersUseCase struct {
	userRepo repository.UserRepository
}

func NewListUsersUseCase(userRepo repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{userRepo: userRepo}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, page, limit int) ([]dto.AdminUserOutput, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	users, total, err := uc.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	output := make([]dto.AdminUserOutput, len(users))
	for i, u := range users {
		output[i] = *userToAdminOutput(&u)
	}

	return output, total, nil
}

type CreateUserUseCase struct {
	userRepo repository.UserRepository
	password auth.PasswordService
}

func NewCreateUserUseCase(userRepo repository.UserRepository, password auth.PasswordService) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo, password: password}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input dto.AdminCreateUserInput) (*dto.AdminUserOutput, error) {
	passwordHash, err := uc.password.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	role := input.Role
	if role == "" {
		role = "user"
	}

	user := &entity.User{
		Email:        input.Email,
		PasswordHash: passwordHash,
		FullName:     input.FullName,
		Phone:        input.Phone,
		Role:         role,
		Active:       input.Active,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return userToAdminOutput(user), nil
}

type UpdateUserUseCase struct {
	userRepo repository.UserRepository
	password auth.PasswordService
}

func NewUpdateUserUseCase(userRepo repository.UserRepository, password auth.PasswordService) *UpdateUserUseCase {
	return &UpdateUserUseCase{userRepo: userRepo, password: password}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, id uuid.UUID, adminID uuid.UUID, input dto.AdminUpdateUserInput) (*dto.AdminUserOutput, error) {
	admin, err := uc.userRepo.GetByID(ctx, adminID)
	if err != nil {
		return nil, fmt.Errorf("get admin: %w", err)
	}

	if !uc.password.Verify(input.AdminPassword, admin.PasswordHash) {
		return nil, fmt.Errorf("invalid admin password")
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.FullName != nil {
		user.FullName = *input.FullName
	}
	if input.Phone != nil {
		user.Phone = input.Phone
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.Active != nil {
		user.Active = *input.Active
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return userToAdminOutput(user), nil
}

type DeactivateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewDeactivateUserUseCase(userRepo repository.UserRepository) *DeactivateUserUseCase {
	return &DeactivateUserUseCase{userRepo: userRepo}
}

func (uc *DeactivateUserUseCase) Execute(ctx context.Context, id uuid.UUID, active bool) (*dto.AdminUserOutput, error) {
	if err := uc.userRepo.SetActive(ctx, id, active); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("set user active: %w", err)
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user after update: %w", err)
	}

	return userToAdminOutput(user), nil
}

func userToAdminOutput(u *entity.User) *dto.AdminUserOutput {
	return &dto.AdminUserOutput{
		ID:         u.ID,
		Email:      u.Email,
		FullName:   u.FullName,
		Phone:      u.Phone,
		Role:       u.Role,
		Active:     u.Active,
		SuperAdmin: u.SuperAdmin,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
