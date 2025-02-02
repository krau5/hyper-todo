package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func IsErrDuplicatedKey(err error) bool {
	if err == nil {
		return false
	}

	var perr *pgconn.PgError
	errors.As(err, &perr)

	return perr.Code == "23505"
}

// HashPassword generates a bcrypt hash for the given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
// Returns true if password matches, false otherwise.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
