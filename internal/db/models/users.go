package models

import "github.com/go-playground/validator/v10"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Метод для валидации данных
func (u *User) Validate() error {
	return validate.Struct(u)
}

// Структура ответа с данными пользователя
type UserIdResponse struct {
	ID int `json:"id"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserUpdate struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"omitempty,min=3,max=20"`
	Email    string `json:"email" validate:"omitempty,email"`
}

func (u *UserUpdate) Validate() error {
	return validate.Struct(u)
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}
