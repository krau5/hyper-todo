package user

import (
	"context"
	"errors"

	"github.com/krau5/hyper-todo/domain"
)

//go:generate mockery --name UsersRepository
type UsersRepository interface {
	Create(ctx context.Context, name, email, password string) error
	GetByEmail(context.Context, string) (domain.User, error)
	GetById(context.Context, int64) (domain.User, error)
}

type Service struct {
	usersRepo UsersRepository
}

var (
	ErrInvalidName     = errors.New("name is missing or empty")
	ErrInvalidEmail    = errors.New("email is missing or empty")
	ErrInvalidPassword = errors.New("password is missing or empty")
	ErrInvalidId       = errors.New("id is missing or empty")
)

func NewService(usersRepo UsersRepository) *Service {
	return &Service{usersRepo: usersRepo}
}

func (s *Service) Create(ctx context.Context, name, email, password string) error {
	if len(name) == 0 {
		return ErrInvalidName
	}

	if len(email) == 0 {
		return ErrInvalidEmail
	}

	if len(password) == 0 {
		return ErrInvalidPassword
	}

	return s.usersRepo.Create(ctx, name, email, password)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if len(email) == 0 {
		return domain.User{}, ErrInvalidEmail
	}

	user, err := s.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s *Service) GetById(ctx context.Context, id int64) (domain.User, error) {
	if id == 0 {
		return domain.User{}, ErrInvalidId
	}

	user, err := s.usersRepo.GetById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
