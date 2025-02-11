package task

import (
	"context"
	"testing"
	"time"

	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/task/mocks"
	userMocks "github.com/krau5/hyper-todo/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
	ctx := context.TODO()
	name := "task name"
	description := "useful task description"
	deadline := time.Now()
	var userId int64 = 1

	t.Run("throws an error if name is invalid", func(t *testing.T) {
		service, _, _ := setupTest(t)

		_, err := service.Create(ctx, "", description, deadline, userId)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrInvalidName.Error())
	})

	t.Run("throws an error if description is invalid", func(t *testing.T) {
		service, _, _ := setupTest(t)

		_, err := service.Create(ctx, name, "", deadline, userId)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrInvalidDescription.Error())
	})

	t.Run("throws an error if the user was not found", func(t *testing.T) {
		service, _, usersRepo := setupTest(t)

		usersRepo.On("GetById", mock.Anything, userId).Return(domain.User{}, gorm.ErrRecordNotFound)

		_, err := service.Create(ctx, name, description, deadline, userId)

		assert.Error(t, err)
		assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	})
}

func TestGetById(t *testing.T) {
	ctx := context.TODO()

	t.Run("throws an error if id is invalid", func(t *testing.T) {
		service, _, _ := setupTest(t)

		_, err := service.GetById(ctx, 0)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrInvalidId.Error())
	})

	t.Run("returns a task if it was found", func(t *testing.T) {
		service, tasksRepo, _ := setupTest(t)

		mockTask := domain.Task{
			Name:        "eat",
			Description: "eat the pizza",
			Deadline:    time.Now(),
			UserId:      1,
		}
		var taskId int64 = 1

		tasksRepo.On("GetById", mock.Anything, taskId).Return(mockTask, nil)

		task, err := service.GetById(ctx, taskId)
		assert.Nil(t, err)
		assert.Equal(t, task, mockTask)
	})
}

func TestGetByUser(t *testing.T) {
	ctx := context.TODO()
	var userId int64 = 1

	t.Run("throws an error if userId is invalid", func(t *testing.T) {
		service, _, _ := setupTest(t)

		_, err := service.GetByUser(ctx, 0)
		assert.Error(t, err)
		assert.EqualError(t, err, ErrInvalidUserId.Error())
	})

	t.Run("throws an error if user was not found", func(t *testing.T) {
		service, _, usersRepo := setupTest(t)

		usersRepo.On("GetById", mock.Anything, userId).Return(domain.User{}, gorm.ErrRecordNotFound)

		_, err := service.GetByUser(ctx, userId)
		assert.Error(t, err)
		assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	})

	t.Run("retrieves and returns tasks if userId is correct", func(t *testing.T) {
		service, tasksRepo, usersRepo := setupTest(t)

		mockTasks := []domain.Task{
			{Name: "task 1", Description: "description 1", Deadline: time.Now()},
			{Name: "task 2", Description: "description 2", Deadline: time.Now()},
		}

		usersRepo.On("GetById", mock.Anything, userId).Return(domain.User{}, nil)
		tasksRepo.On("GetByUser", mock.Anything, userId).Return(mockTasks, nil)

		tasks, err := service.GetByUser(ctx, userId)
		assert.Nil(t, err)
		assert.Equal(t, mockTasks, tasks)
	})
}

func setupTest(t *testing.T) (*Service, *mocks.TasksRepository, *userMocks.UsersRepository) {
	tasksRepo := mocks.NewTasksRepository(t)
	usersRepo := userMocks.NewUsersRepository(t)
	service := NewService(tasksRepo, usersRepo)

	return service, tasksRepo, usersRepo
}
