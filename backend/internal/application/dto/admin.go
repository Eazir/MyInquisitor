package dto

import (
	"time"

	"github.com/google/uuid"
)

type AdminCreateUserInput struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
	FullName string  `json:"full_name" validate:"required"`
	Phone    *string `json:"phone,omitempty"`
	Role     string  `json:"role" validate:"required,oneof=user super_admin"`
	Active   bool    `json:"active"`
}

type AdminUpdateUserInput struct {
	Email        *string `json:"email,omitempty"`
	FullName     *string `json:"full_name,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Role         *string `json:"role,omitempty"`
	Active       *bool   `json:"active,omitempty"`
	AdminPassword string `json:"admin_password" validate:"required"`
}

type InviteTokenOutput struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	URL       string    `json:"url"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminUserOutput struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	FullName   string    `json:"full_name"`
	Phone      *string   `json:"phone"`
	Role       string    `json:"role"`
	Active     bool      `json:"active"`
	SuperAdmin bool      `json:"super_admin"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
