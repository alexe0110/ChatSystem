package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Login          string    `json:"login"`
	Name           string    `json:"name"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
