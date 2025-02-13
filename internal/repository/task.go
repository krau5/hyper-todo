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
	taskModel := TaskModel{
		Task: domain.Task{Name: name, Description: description, Deadline: deadline, UserId: userId},
	}

	result := r.db.WithContext(ctx).Create(&taskModel)
	if result.Error != nil {
		return domain.Task{}, result.Error
	}

	taskModel.Task.ID = int64(taskModel.Model.ID)

	return taskModel.Task, nil
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
	result := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&rawTasks)
	if result.Error != nil {
		return []domain.Task{}, result.Error
	}

	tasks := make([]domain.Task, len(rawTasks))
	for i, taskModel := range rawTasks {
		tasks[i] = taskModel.Task
	}

	return tasks, nil
}

func (r *tasksRepository) UpdateById(ctx context.Context, id int64, data domain.UpdateTaskData) (domain.Task, error) {
	taskModel := TaskModel{}

	result := r.db.WithContext(ctx).First(&taskModel, id)
	if result.Error != nil {
		return domain.Task{}, result.Error
	}

	updates := make(map[string]interface{})
	if data.Name != nil && len(*data.Name) != 0 {
		updates["name"] = *data.Name
	}
	if data.Description != nil && len(*data.Description) != 0 {
		updates["description"] = *data.Description
	}
	if data.Deadline != nil && !(*data.Deadline).IsZero() {
		updates["deadline"] = *data.Deadline
	}
	if data.Completed != nil {
		updates["completed"] = *data.Completed
	}

	result = r.db.WithContext(ctx).Model(&taskModel).Updates(updates)
	if result.Error != nil {
		return domain.Task{}, result.Error
	}

	return taskModel.Task, nil
}

func (r *tasksRepository) DeleteById(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&TaskModel{}, id)
	return result.Error
}
