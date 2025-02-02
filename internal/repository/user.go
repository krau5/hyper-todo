package repository

import (
	"context"

	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/internal/utils"
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
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := UserModel{
		User: domain.User{Name: name, Email: email, Password: hash},
	}

	result := r.db.WithContext(ctx).Create(&user)
	return result.Error
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user := UserModel{}

	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return domain.User{}, result.Error
	}

	return user.User, nil
}
