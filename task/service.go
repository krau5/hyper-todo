package task

import (
	"context"
	"errors"
	"time"

	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/user"
	"gorm.io/gorm"
)

//go:generate mockery --name TasksRepository
type TasksRepository interface {
	Create(ctx context.Context, name, description string, deadline time.Time, userId int64) (domain.Task, error)
	GetById(context.Context, int64) (domain.Task, error)
	GetByUser(context.Context, int64) ([]domain.Task, error)
	UpdateById(context.Context, int64, domain.UpdateTaskData) (domain.Task, error)
	DeleteById(context.Context, int64) error
}

type Service struct {
	usersRepo user.UsersRepository
	tasksRepo TasksRepository
}

var (
	ErrInvalidName        = errors.New("name is missing or empty")
	ErrInvalidDescription = errors.New("description is missing or empty")
	ErrInvalidId          = errors.New("id is missing or empty")
	ErrInvalidUserId      = errors.New("userId is missing or empty")
)

func NewService(tasksRepo TasksRepository, usersRepo user.UsersRepository) *Service {
	return &Service{
		tasksRepo: tasksRepo,
		usersRepo: usersRepo,
	}
}

func (s *Service) Create(ctx context.Context, name, description string, deadline time.Time, userId int64) (domain.Task, error) {
	if len(name) == 0 {
		return domain.Task{}, ErrInvalidName
	}

	if len(description) == 0 {
		return domain.Task{}, ErrInvalidDescription
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
		return domain.Task{}, ErrInvalidId
	}

	task, err := s.tasksRepo.GetById(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *Service) GetByUser(ctx context.Context, userId int64) ([]domain.Task, error) {
	if userId == 0 {
		return []domain.Task{}, ErrInvalidUserId
	}

	_, err := s.usersRepo.GetById(ctx, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []domain.Task{}, gorm.ErrRecordNotFound
	}

	tasks, err := s.tasksRepo.GetByUser(ctx, userId)
	if err != nil {
		return []domain.Task{}, err
	}

	return tasks, nil
}

func (s *Service) UpdateById(ctx context.Context, id int64, data domain.UpdateTaskData) (domain.Task, error) {
	if id == 0 {
		return domain.Task{}, ErrInvalidId
	}

	task, err := s.tasksRepo.UpdateById(ctx, id, data)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (s *Service) DeleteById(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrInvalidId
	}

	return s.tasksRepo.DeleteById(ctx, id)
}
