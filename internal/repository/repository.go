package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"recipes/internal/db/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (int, error)
	GetUserById(ctx context.Context, id int) (models.User, error)
	UpdateUser(ctx context.Context, user *models.UserUpdate) error
	DeleteUser(ctx context.Context, id int) error
	ListUser(ctx context.Context) ([]models.UserResponse, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}
