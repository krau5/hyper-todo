package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/krau5/hyper-todo/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUsersService(t *testing.T) {
	usersRepo := mocks.NewUsersRepository(t)
	service := NewService(usersRepo)

	ctx := context.TODO()
	name := "user"
	email := "user@example.com"
	password := "password123"

	t.Run("throws an error if name is empty", func(t *testing.T) {
		expectedErr := fmt.Errorf("field name is missing or empty")
		err := service.Create(ctx, "", email, password)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("throws an error if email is empty", func(t *testing.T) {
		expectedErr := fmt.Errorf("field email is missing or empty")
		err := service.Create(ctx, name, "", password)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("throws an error if password is empty", func(t *testing.T) {
		expectedErr := fmt.Errorf("field password is missing or empty")
		err := service.Create(ctx, name, email, "")

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("returns nil if user was successfully created", func(t *testing.T) {
		usersRepo.On("Create", mock.Anything, name, email, password).Return(nil)
		err := service.Create(ctx, name, email, password)

		assert.Nil(t, err)
	})
}
