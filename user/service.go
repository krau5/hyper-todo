package user

import (
	"context"
	"fmt"

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

func NewService(usersRepo UsersRepository) *Service {
	return &Service{usersRepo: usersRepo}
}

func (s *Service) Create(ctx context.Context, name, email, password string) error {
	if len(name) == 0 {
		return fmt.Errorf("field name is missing or empty")
	}

	if len(email) == 0 {
		return fmt.Errorf("field email is missing or empty")
	}

	if len(password) == 0 {
		return fmt.Errorf("field password is missing or empty")
	}

	return s.usersRepo.Create(ctx, name, email, password)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if len(email) == 0 {
		return domain.User{}, fmt.Errorf("field email is missing or empty")
	}

	user, err := s.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
