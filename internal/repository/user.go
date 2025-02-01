package repository

import (
	"context"

	"github.com/krau5/hyper-todo/domain"
	"gorm.io/gorm"
)

type UserModel struct {
	domain.User
	gorm.Model
}

// Repository implements UsersRepository interface
type Repository struct {
	db *gorm.DB
}

// NewUserRepository returns the implementation of UsersRepository interface
func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, name, email, password string) error {
	user := UserModel{
		User: domain.User{Name: name, Email: email, Password: password},
	}

	result := r.db.Create(&user)
	return result.Error
}

func (r *Repository) FindByEmail(email string) (domain.User, error) {
	return domain.User{}, nil
}
