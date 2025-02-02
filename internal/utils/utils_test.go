package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "Password_123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}

	assert.True(t, VerifyPassword(password, hash))
	assert.False(t, VerifyPassword("Password_321", hash))
}
