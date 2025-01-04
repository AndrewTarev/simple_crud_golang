package repository

import (
	"context"

	"recipes/internal/db/models"
)

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (int, error) {
	var id int
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *userRepository) GetUserById(ctx context.Context, id int) (models.User, error) {
	var user models.User
	query := `SELECT id, username, email FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	if err := row.Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.UserUpdate) error {
	query := `UPDATE users SET username = $1, email = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, user.Username, user.Email, user.ID)
	return err
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *userRepository) ListUser(ctx context.Context) ([]models.UserResponse, error) {
	query := `SELECT id, username, email FROM users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
