package rest

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/task"
)

type TasksService interface {
	Create(ctx context.Context, name, description string, deadline time.Time, userId int64) (domain.Task, error)
	GetById(context.Context, int64) (domain.Task, error)
	GetByUser(context.Context, int64) ([]domain.Task, error)
	UpdateById(context.Context, int64, task.TaskUpdate) (domain.Task, error)
	DeleteById(context.Context, int64) error
}

type TasksHandler struct {
	tasksService TasksService
}

func NewTasksHandler(r *gin.Engine, tasksService TasksService) {
	h := &TasksHandler{
		tasksService: tasksService,
	}

	r.POST("/tasks", h.handleCreateTask)
}

func (h *TasksHandler) handleCreateTask(c *gin.Context) {}
