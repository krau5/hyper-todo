package task

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/task/mocks"
	userMocks "github.com/krau5/hyper-todo/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateTask(t *testing.T) {
	tasksRepo := mocks.NewTasksRepository(t)
	usersRepo := userMocks.NewUsersRepository(t)
	service := NewService(tasksRepo, usersRepo)

	ctx := context.TODO()
	name := "task name"
	description := "useful task description"
	deadline := time.Now()
	var userId int64 = 1

	t.Run("throws an error if name is empty", func(t *testing.T) {
		expectedErr := fmt.Errorf("field name is missing or empty")
		_, err := service.Create(ctx, "", description, deadline, userId)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("throws an error if description is empty", func(t *testing.T) {
		expectedErr := fmt.Errorf("field description is missing or empty")
		_, err := service.Create(ctx, name, "", deadline, userId)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("throws an error if the user does not exist", func(t *testing.T) {
		usersRepo.On("GetById", mock.Anything, userId).Return(domain.User{}, gorm.ErrRecordNotFound)

		_, err := service.Create(ctx, name, description, deadline, userId)

		assert.Error(t, err)
		assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	})
}

func TestGetTask(t *testing.T) {
	tasksRepo := mocks.NewTasksRepository(t)
	usersRepo := userMocks.NewUsersRepository(t)
	service := NewService(tasksRepo, usersRepo)

	ctx := context.TODO()
	mockTask := domain.Task{
		Name:        "eat",
		Description: "eat the pizza",
		Deadline:    time.Now(),
		UserId:      1,
	}

	t.Run("throws an error if id is zero", func(t *testing.T) {
		expectedErr := fmt.Errorf("field id is missing or empty")
		_, err := service.GetById(ctx, 0)

		assert.Error(t, err)
		assert.EqualError(t, err, expectedErr.Error())
	})

	t.Run("returns a task if it was found", func(t *testing.T) {
		var taskId int64 = 1
		tasksRepo.On("GetById", mock.Anything, taskId).Return(mockTask, nil)
		task, err := service.GetById(ctx, taskId)

		assert.Nil(t, err)
		assert.Equal(t, task, mockTask)
	})
}
