package task

import (
	"context"
	"time"

	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/user"
)

type TasksRepository interface {
	Create(ctx context.Context, name, description string, deadline time.Time, userId int64) (domain.Task, error)
	GetById(context.Context, int64) (domain.Task, error)
	GetByUser(context.Context, int64) ([]domain.Task, error)
	UpdateById(context.Context, int64, TaskUpdate) (domain.Task, error)
	DeleteById(context.Context, int64) error
}

type TaskUpdate struct {
	Name        *string
	Description *string
	Deadline    *time.Time
	Completed   *bool
}

type Service struct {
	usersRepo user.UsersRepository
	tasksRepo TasksRepository
}

func NewTasksService(tasksRepo TasksRepository, usersRepo user.UsersRepository) *Service {
	return &Service{
		tasksRepo: tasksRepo,
		usersRepo: usersRepo,
	}
}

func (s *Service) Create(ctx context.Context, name, description string, deadline time.Time, userId int64) (domain.Task, error) {
	return domain.Task{}, nil
}

func (s *Service) GetById(ctx context.Context, id int64) (domain.Task, error) {
	return domain.Task{}, nil
}

func (s *Service) GetByUser(ctx context.Context, userId int64) ([]domain.Task, error) {
	return []domain.Task{}, nil
}

func (s *Service) UpdateById(ctx context.Context, id int64, data TaskUpdate) (domain.Task, error) {
	return domain.Task{}, nil
}

func (s *Service) DeleteById(ctx context.Context, id int64) error {
	return nil
}
