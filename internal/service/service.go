package service

import (
	"context"

	"recipes/internal/db/models"
	"recipes/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (int, error)
	GetUserById(ctx context.Context, id int) (models.User, error)
	UpdateUser(ctx context.Context, user *models.UserUpdate) error
	DeleteUser(ctx context.Context, id int) error
	ListUser(ctx context.Context) ([]models.UserResponse, error)
}

type Service struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) UserService {
	return &Service{repo: repo}
}
