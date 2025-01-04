package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

const (
	StatusOK        = "success"
	StatusError     = "failed"
	UniqueViolation = "23505"
)

// SuccessResponse - схема успешного ответа
type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// NewErrorResponse - обработчик ошибок
func NewErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	logger.Error(err)
	// Формируем ответ, но не передаем детали ошибки пользователю
	c.AbortWithStatusJSON(statusCode, gin.H{
		"status": StatusError,
		"error":  message,
	})
}
