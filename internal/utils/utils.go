package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsErrDuplicatedKey(err error) bool {
	var perr *pgconn.PgError
	errors.As(err, &perr)

	return perr.Code == "23505"
}
