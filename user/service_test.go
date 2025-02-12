package user

import (
	"context"
	"testing"

	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	usersRepo := mocks.NewUsersRepository(t)
	service := NewService(usersRepo)

	ctx := context.TODO()
	name := "user"
	email := "user@example.com"
	password := "password123"

	t.Run("throws an error if name is invalid", func(t *testing.T) {
		err := service.Create(ctx, "", email, password)
		assert.EqualError(t, err, ErrInvalidName.Error())
	})

	t.Run("throws an error if email is invalid", func(t *testing.T) {
		err := service.Create(ctx, name, "", password)
		assert.EqualError(t, err, ErrInvalidEmail.Error())
	})

	t.Run("throws an error if password is invalid", func(t *testing.T) {
		err := service.Create(ctx, name, email, "")
		assert.EqualError(t, err, ErrInvalidPassword.Error())
	})

	t.Run("returns nil if user was successfully created", func(t *testing.T) {
		usersRepo.On("Create", mock.Anything, name, email, password).Return(nil)
		err := service.Create(ctx, name, email, password)

		assert.Nil(t, err)
	})
}

func TestGetByEmail(t *testing.T) {
	usersRepo := mocks.NewUsersRepository(t)
	service := NewService(usersRepo)

	ctx := context.TODO()
	name := "user"
	email := "user@example.com"
	password := "password123"

	t.Run("throws an error if email is invalid", func(t *testing.T) {
		user, err := service.GetByEmail(ctx, "")

		assert.Equal(t, domain.User{}, user)
		assert.EqualError(t, err, ErrInvalidEmail.Error())
	})

	t.Run("returns a user if user was successfully retrieved", func(t *testing.T) {
		mockUser := domain.User{Name: name, Email: email, Password: password}

		usersRepo.On("GetByEmail", mock.Anything, email).Return(mockUser, nil)
		user, err := service.GetByEmail(ctx, email)

		assert.Equal(t, mockUser, user)
		assert.Nil(t, err)
	})
}

func TestGetById(t *testing.T) {
	usersRepo := mocks.NewUsersRepository(t)
	service := NewService(usersRepo)

	ctx := context.TODO()
	name := "user"
	email := "user@example.com"
	password := "password123"
	var userId int64 = 1

	t.Run("throws an error if id is invalid", func(t *testing.T) {
		user, err := service.GetById(ctx, 0)

		assert.Equal(t, domain.User{}, user)
		assert.EqualError(t, err, ErrInvalidId.Error())
	})

	t.Run("returns a user if he was successfully retrieved", func(t *testing.T) {
		mockUser := domain.User{Name: name, Email: email, Password: password}

		usersRepo.On("GetById", mock.Anything, userId).Return(mockUser, nil)
		user, err := service.GetById(ctx, userId)

		assert.Equal(t, mockUser, user)
		assert.Nil(t, err)
	})
}
