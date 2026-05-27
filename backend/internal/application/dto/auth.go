package dto

import "github.com/google/uuid"

type RegisterInput struct {
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8"`
	FullName    string  `json:"full_name" validate:"required,min=2"`
	Phone       *string `json:"phone,omitempty"`
	InviteToken string  `json:"invite_token" validate:"required"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AuthOutput struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         UserDTO  `json:"user"`
}

type UserDTO struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Role     string    `json:"role"`
}
