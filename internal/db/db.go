package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	logger "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"simple_crud_go/configs"
)

func ConnectPostgres(cfg *configs.PostgresConfig) (*pgxpool.Pool, error) {
	// Формируем строку подключения
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	// Настраиваем контекст для подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем пул соединений
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Errorf("Error parsing PostgreSQL DSN: %v", err)
		return nil, err
	}

	// Применяем дополнительные настройки пула
	poolConfig.MaxConns = 50                       // Максимальное количество соединений
	poolConfig.MinConns = 5                        // Минимальное количество соединений
	poolConfig.HealthCheckPeriod = 1 * time.Minute // Период проверки соединений

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Errorf("Failed to connect to PostgreSQL: %v", err)
		return nil, err
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		logger.Errorf("PostgreSQL ping failed: %v", err)
		return nil, err
	}

	logger.Infof("Successfully connected to PostgreSQL at %s:%d", cfg.Host, cfg.Port)
	return pool, nil
}

func ApplyMigrations() {
	m, err := migrate.New(
		"file:///app/internal/db/migrations",
		"postgres://postgres:postgres@db:5432/postgres?sslmode=disable",
	)
	if err != nil {
		logger.Fatalf("Could not initialize migrate: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatalf("Could not apply migrations: %v", err)
	}

	logger.Println("Migrations applied successfully!")
}
