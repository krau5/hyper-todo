package repository

import (
	"time"

	"github.com/krau5/hyper-todo/domain"
	"gorm.io/gorm"
)

type TaskModel struct {
	domain.Task
	gorm.Model
}

type tasksRepository struct {
	db *gorm.DB
}

func NewTasksRepository(db *gorm.DB) *tasksRepository {
	return &tasksRepository{db: db}
}

func (r *tasksRepository) Create(name, description string, deadline time.Time, userId int64) (domain.Task, error) {
	return domain.Task{}, nil
}

func (r *tasksRepository) GetById(id int64) (domain.Task, error) {
	return domain.Task{}, nil
}

func (r *tasksRepository) GetByUser(userId int64) ([]domain.Task, error) {
	return []domain.Task{}, nil
}

func (r *tasksRepository) UpdateById(id int64, data TaskUpdate) (domain.Task, error) {
	return domain.Task{}, nil
}

func (r *tasksRepository) DeleteById(id int64) error {
	return nil
}

type TaskUpdate struct {
	Name        *string
	Description *string
	Deadline    *time.Time
	Completed   *bool
}
