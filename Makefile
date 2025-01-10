migrateup:
	migrate -path ./internal/db/migrations -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up

migratedown:
	migrate -path ./internal/db/migrationsy -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' down