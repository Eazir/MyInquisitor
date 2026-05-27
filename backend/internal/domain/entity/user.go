package entity

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string
	Phone        *string
	Role         string
	Active       bool
	SuperAdmin   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateUserInput struct {
	Email    string
	Password string
	FullName string
	Phone    *string
}

type UpdateUserInput struct {
	ID       uuid.UUID
	Email    *string
	FullName *string
	Phone    *string
	Active   *bool
	Role     *string
}
