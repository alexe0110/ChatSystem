package service

import (
	"context"
	"github.com/alexe0110/chat-system/internal/model"
	"github.com/alexe0110/chat-system/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		repo: userRepository,
	}
}

func (r *UserService) Register(ctx context.Context, login, name, password string) (*model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := r.repo.CreateUser(ctx, login, name, string(passwordHash))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserService) Login(ctx context.Context, login, password string) (*model.User, error) {
	user, err := r.repo.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}
