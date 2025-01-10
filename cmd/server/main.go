package main

import (
	logger "github.com/sirupsen/logrus"

	_ "simple_crud_go/docs"

	"simple_crud_go/configs"
	"simple_crud_go/internal/db"
	"simple_crud_go/internal/handler"
	"simple_crud_go/internal/repository"
	"simple_crud_go/internal/service"
	"simple_crud_go/pkg/logging"
	"simple_crud_go/pkg/server"
)

// @title           User Management API
// @version         1.0
// @description     API для управления пользователями: создание, обновление, удаление и получение данных.

// @host      localhost:8000
// @BasePath  /user
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

	db.ApplyMigrations()

	repo := repository.NewUserRepository(dbConn)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	// Настройка и запуск сервера
	server.SetupAndRunServer(&cfg.Server, handlers.InitRouters())
}
