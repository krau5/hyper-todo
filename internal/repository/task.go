package repository

import (
	"context"
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

func (r *tasksRepository) Create(ctx context.Context, name, description string, deadline time.Time, userId int64) (domain.Task, error) {
	task := domain.Task{Name: name, Description: description, Deadline: deadline, UserId: userId}
	val := TaskModel{
		Task: task,
	}

	result := r.db.WithContext(ctx).Create(&val)
	if result.Error != nil {
		return domain.Task{}, result.Error
	}

	return task, nil
}

func (r *tasksRepository) GetById(ctx context.Context, id int64) (domain.Task, error) {
	task := TaskModel{}

	result := r.db.WithContext(ctx).First(&task, id)
	if result.Error != nil {
		return domain.Task{}, result.Error
	}

	return task.Task, nil
}

func (r *tasksRepository) GetByUser(ctx context.Context, userId int64) ([]domain.Task, error) {
	rawTasks := []TaskModel{}
	result := r.db.WithContext(ctx).Where("userId = ?", userId).Find(&rawTasks)
	if result.Error != nil {
		return []domain.Task{}, result.Error
	}

	tasks := make([]domain.Task, len(rawTasks))
	for i, taskModel := range rawTasks {
		tasks[i] = taskModel.Task
	}

	return tasks, nil
}

func (r *tasksRepository) UpdateById(ctx context.Context, id int64, data *domain.Task) (domain.Task, error) {
	return domain.Task{}, nil
}

func (r *tasksRepository) DeleteById(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&TaskModel{}, id)
	return result.Error
}
