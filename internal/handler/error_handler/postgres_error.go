package error_handler

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// ErrorCode - обработчик ошибок для pgx драйвера
func ErrorCode(err error) (string, string) {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code, pgErr.ConstraintName
	}

	return "", ""
}
