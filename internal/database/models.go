// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Advertisement struct {
	ID           uuid.UUID
	Title        string
	Description  string
	ImageAddress string
	Price        int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserID       uuid.UUID
}

type User struct {
	ID             uuid.UUID
	Login          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
