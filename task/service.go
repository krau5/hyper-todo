package task

import (
	"time"

	"github.com/krau5/hyper-todo/domain"
)

type TasksRepository interface {
	Create(name, description string, deadline time.Time, userId int64) (domain.Task, error)
	GetById(id int64) (domain.Task, error)
	GetByUser(userId int64) ([]domain.Task, error)
	UpdateById(id int64, data TaskUpdate) (domain.Task, error)
	DeleteById(id int64) error
}

type TaskUpdate struct {
	Name        *string
	Description *string
	Deadline    *time.Time
	Completed   *bool
}
