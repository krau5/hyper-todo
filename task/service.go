package task

import (
	"context"
	"fmt"
	"time"

	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/user"
)

//go:generate mockery --name TasksRepository
type TasksRepository interface {
	Create(ctx context.Context, name, description string, deadline time.Time, userId int64) (domain.Task, error)
	GetById(context.Context, int64) (domain.Task, error)
	GetByUser(context.Context, int64) ([]domain.Task, error)
	UpdateById(context.Context, int64, *domain.Task) (domain.Task, error)
	DeleteById(context.Context, int64) error
}

type Service struct {
	usersRepo user.UsersRepository
	tasksRepo TasksRepository
}

func NewService(tasksRepo TasksRepository, usersRepo user.UsersRepository) *Service {
	return &Service{
		tasksRepo: tasksRepo,
		usersRepo: usersRepo,
	}
}

func (s *Service) Create(ctx context.Context, name, description string, deadline time.Time, userId int64) (domain.Task, error) {
	if len(name) == 0 {
		return domain.Task{}, fmt.Errorf("field name is missing or empty")
	}

	if len(description) == 0 {
		return domain.Task{}, fmt.Errorf("field description is missing or empty")
	}

	_, err := s.usersRepo.GetById(ctx, userId)
	if err != nil {
		return domain.Task{}, err
	}

	task, err := s.tasksRepo.Create(ctx, name, description, deadline, userId)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *Service) GetById(ctx context.Context, id int64) (domain.Task, error) {
	if id == 0 {
		return domain.Task{}, fmt.Errorf("field id is missing or empty")
	}

	task, err := s.tasksRepo.GetById(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *Service) GetByUser(ctx context.Context, userId int64) ([]domain.Task, error) {
	return []domain.Task{}, nil
}

func (s *Service) UpdateById(ctx context.Context, id int64, data *domain.Task) (domain.Task, error) {
	return domain.Task{}, nil
}

func (s *Service) DeleteById(ctx context.Context, id int64) error {
	if id == 0 {
		return fmt.Errorf("field id is missing or empty")
	}

	return s.tasksRepo.DeleteById(ctx, id)
}
