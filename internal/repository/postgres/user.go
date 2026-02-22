package postgres

import (
	"context"
	"database/sql"
	"github.com/alexe0110/chat-system/internal/model"
	"github.com/google/uuid"
)

const queryCreateUser = `INSERT INTO users (id, login, name, hashed_password) 
          VALUES ($1, $2, $3, $4) 
          RETURNING id, login, name, created_at, updated_at`

const getUserByID = `SELECT id, login, name, created_at, updated_at FROM users WHERE id=$1`
const getUserByLogin = `SELECT id, login, name, hashed_password, created_at, updated_at FROM users WHERE login=$1`

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, login, name, hashedPassword string) (*model.User, error) {
	id := uuid.New()
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, queryCreateUser, id, login, name, hashedPassword).Scan(&user.ID, &user.Login, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, getUserByID, id).Scan(&user.ID, &user.Login, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, getUserByLogin, login).Scan(&user.ID, &user.Login, &user.Name, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}
