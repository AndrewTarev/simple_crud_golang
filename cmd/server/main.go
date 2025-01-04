package main

import (
	"fmt"

	logger "github.com/sirupsen/logrus"

	"recipes/configs"
	"recipes/internal/db"
	"recipes/internal/handler"
	"recipes/internal/repository"
	"recipes/internal/service"
	"recipes/pkg/logging"
	"recipes/pkg/server"
)

// @title           User Management API
// @version         1.0
// @description     API для управления пользователями: создание, обновление, удаление и получение данных.

// @host      localhost:8080
// @BasePath  /
func main() {
	// Загружаем конфигурацию
	cfg, err := configs.LoadConfig("./configs")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}

	// Настройка логгера
	logging.SetupLogger(&cfg.Logging)

	// Подключение к базе данных
	dbConn, err := db.ConnectPostgres(&cfg.Database)
	if err != nil {
		logger.Fatalf("Database connection failed: %v", err)
	}
	defer dbConn.Close()

	repo := repository.NewUserRepository(dbConn)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	// Настройка и запуск сервера
	server.SetupAndRunServer(&cfg.Server, handlers.InitRouters())
	fmt.Println("Server started")
}
