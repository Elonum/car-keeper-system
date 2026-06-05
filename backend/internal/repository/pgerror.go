package repository

import (
	"errors"

	"github.com/carkeeper/backend/internal/apperr"
	"github.com/jackc/pgx/v5/pgconn"
)

// mapUniqueViolation returns a conflict API error when err is a PostgreSQL unique violation.
func mapUniqueViolation(err error, message string) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return apperr.Conflict(message)
	}
	return nil
}
