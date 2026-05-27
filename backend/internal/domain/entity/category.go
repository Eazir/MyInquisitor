package entity

import (
	"time"
	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Type      string
	Icon      *string
	Color     *string
	CreatedAt time.Time
}

type CreateCategoryInput struct {
	UserID uuid.UUID
	Name   string
	Type   string
	Icon   *string
	Color  *string
}
