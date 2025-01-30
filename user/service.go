package user

import "github.com/krau5/hyper-todo/domain"

type UsersRepository interface {
	Create(name, email, password string) (domain.User, error)
	FindByEmail(string) (domain.User, error)
}
