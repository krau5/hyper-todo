package user

import (
	"context"
	"fmt"

	"github.com/krau5/hyper-todo/domain"
)

type UsersRepository interface {
	Create(ctx context.Context, name, email, password string) error
	FindByEmail(string) (domain.User, error)
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
