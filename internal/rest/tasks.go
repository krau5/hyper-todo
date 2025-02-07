package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/internal/rest/middleware"
)

//go:generate mockery --name TasksService
type TasksService interface {
	Create(ctx context.Context, name, description string, deadline time.Time, userId int64) (domain.Task, error)
	GetById(context.Context, int64) (domain.Task, error)
	GetByUser(context.Context, int64) ([]domain.Task, error)
	UpdateById(context.Context, int64, *domain.Task) (domain.Task, error)
	DeleteById(context.Context, int64) error
}

type TasksHandler struct {
	tasksService TasksService
}

type CreateTaskBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	UserId      int64  `json:"userId"`
}

func NewTasksHandler(r *gin.Engine, tasksService TasksService) {
	h := &TasksHandler{
		tasksService: tasksService,
	}

	r.POST("/tasks", middleware.AuthMiddleware, h.handleCreateTask)
}

func (h *TasksHandler) handleCreateTask(c *gin.Context) {
	var data CreateTaskBody

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deadline, err := time.Parse(time.RFC3339, data.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.tasksService.Create(
		c.Request.Context(),
		data.Name,
		data.Description,
		deadline,
		data.UserId,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}
