migrateup:
	migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up

migratedown:
	migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' down