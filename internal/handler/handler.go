package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "simple_crud_go/docs"

	"simple_crud_go/internal/service"
)

type Handler struct {
	services service.UserService
}

func NewHandler(services service.UserService) *Handler {
	return &Handler{services: services}
}

// InitRouters инициализирует маршруты приложения
func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	// Роут для Swagger-документации
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Роуты для пользователя
	user := router.Group("/user")
	{
		user.POST("/", h.CreateUser)
		user.GET("/:id", h.GetUserByID)
		user.PUT("/:id", h.UpdateUser)
		user.DELETE("/:id", h.DeleteUser)
		user.GET("/", h.ListUser)
	}

	return router
}
