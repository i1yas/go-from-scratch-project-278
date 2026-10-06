package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func isConstraintError(err error, constraintName string) bool {
	pgErr, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return false
	}

	return pgErr.ConstraintName == constraintName
}
