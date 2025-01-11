migrateup:
	migrate -path ./internal/db/migrations -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up

migratedown:
	migrate -path ./internal/db/migrationsy -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' down

test-mock:
	mockgen -source=internal/service/service.go -destination=internal/service/mocks/mock_user_service.go -package=mocks