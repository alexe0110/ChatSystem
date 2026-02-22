package repository

import (
	"context"
	"github.com/alexe0110/chat-system/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, login, name, hashedPassword string) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
}
