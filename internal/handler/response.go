package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

const (
	StatusSuccess   = "success"
	StatusError     = "failed"
	UniqueViolation = "23505"
)

// SuccessResponse - схема успешного ответа
type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// NewErrorResponse обрабатывает ошибки и формирует JSON-ответ с ошибкой.
func NewErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	logger.Error(err)
	// Формируем ответ, передавая объект ErrorResponse
	c.AbortWithStatusJSON(statusCode, gin.H{
		"status": StatusError,
		"error": ErrorResponse{
			Message: message,
		},
	})
}
