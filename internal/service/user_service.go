package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	logger "github.com/sirupsen/logrus"

	"simple_crud_go/internal/db/models"
	"simple_crud_go/pkg/utils"
)

func (s *Service) CreateUser(ctx context.Context, user *models.User) (int, error) {
	// Хэшируем пароль пользователя
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		// Логируем ошибку (если у вас есть настроенный логгер)
		logger.Errorf("Ошибка хэширования пароля: %v", err)
		// Возвращаем ошибку вверх по цепочке
		return 0, fmt.Errorf("не удалось хэшировать пароль: %w", err)
	}

	// Заменяем пароль пользователя на хэш
	user.Password = hashedPassword

	// Сохраняем пользователя в репозитории
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) GetUserById(ctx context.Context, id int) (models.User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *Service) UpdateUser(ctx context.Context, user *models.UserUpdate) error {
	// Проверяем, существует ли пользователь
	existingUser, err := s.repo.GetUserById(ctx, user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	// Если данные не предоставлены, оставляем старые
	if user.Username == "" {
		user.Username = existingUser.Username
	}
	if user.Email == "" {
		user.Email = existingUser.Email
	}

	// Обновляем данные
	return s.repo.UpdateUser(ctx, user)
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *Service) ListUser(ctx context.Context) ([]models.UserResponse, error) {
	return s.repo.ListUser(ctx)
}
