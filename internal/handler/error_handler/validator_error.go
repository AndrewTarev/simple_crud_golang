package error_handler

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func ParseValidationErrors(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, ve := range validationErrors {
			field := ve.Field() // Имя поля
			tag := ve.Tag()     // Тег валидации, например, "required" или "min"
			switch tag {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", field))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be a valid email", field))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters", field, ve.Param()))
			case "max":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must not exceed %s characters", field, ve.Param()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("%s is invalid", field))
			}
		}
		return strings.Join(errorMessages, "; ")
	}
	return "Invalid input data"
}
