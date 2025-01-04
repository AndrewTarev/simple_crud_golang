package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword принимает строку пароля и возвращает его хэш.
func HashPassword(password string) (string, error) {
	// Генерируем хэш пароля с использованием bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword сравнивает предоставленный пароль с хэшированным значением.
func CheckPassword(password, hashedPassword string) bool {
	// Проверяем пароль
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
