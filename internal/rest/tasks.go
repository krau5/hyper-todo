package rest

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/domain"
	appErrors "github.com/krau5/hyper-todo/internal/rest/errors"
	"github.com/krau5/hyper-todo/internal/rest/middleware"
	"gorm.io/gorm"
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

var (
	ErrInvalidDeadline       = appErrors.NewResponseError(http.StatusBadRequest, "failed to parse deadline")
	ErrFailedToCreateTask    = appErrors.NewResponseError(http.StatusBadRequest, "failed to create task")
	ErrInvalidTaskId         = appErrors.NewResponseError(http.StatusBadRequest, "task id is missing or invalid")
	ErrTaskNotFound          = appErrors.NewResponseError(http.StatusNotFound, "task was not found")
	ErrFailedToDeleteTask    = appErrors.NewResponseError(http.StatusInternalServerError, "failed to delete task")
	ErrFailedToRetrieveTasks = appErrors.NewResponseError(http.StatusInternalServerError, "failed to retrieve tasks")
)

func NewTasksHandler(r *gin.Engine, tasksService TasksService) {
	h := &TasksHandler{
		tasksService: tasksService,
	}

	r.GET("/tasks", middleware.AuthMiddleware, h.handleGetTasks)
	r.POST("/tasks", middleware.AuthMiddleware, h.handleCreateTask)
	r.DELETE("/tasks/:taskId", middleware.AuthMiddleware, h.handleDeleteTask)
}

func (h *TasksHandler) handleGetTasks(c *gin.Context) {
	tasks, err := h.tasksService.GetByUser(c.Request.Context(), c.GetInt64("user-id"))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(ErrUserNotFound.Status, ErrUserNotFound)
		return
	}

	if err != nil {
		c.JSON(ErrFailedToRetrieveTasks.Status, ErrFailedToRetrieveTasks)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TasksHandler) handleCreateTask(c *gin.Context) {
	var data CreateTaskBody

	if err := c.BindJSON(&data); err != nil {
		c.JSON(appErrors.ErrInvalidBody.Status, appErrors.ErrInvalidBody)
		return
	}

	deadline, err := time.Parse(time.RFC3339, data.Deadline)
	if err != nil {
		c.JSON(ErrInvalidDeadline.Status, ErrInvalidDeadline)
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
		c.JSON(ErrFailedToCreateTask.Status, ErrFailedToCreateTask)
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TasksHandler) handleDeleteTask(c *gin.Context) {
	rawTaskId := c.Param("taskId")
	taskId, err := strconv.ParseInt(rawTaskId, 10, 64)
	if err != nil {
		c.JSON(ErrInvalidTaskId.Status, ErrInvalidTaskId)
		return
	}

	task, err := h.tasksService.GetById(c.Request.Context(), taskId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(ErrTaskNotFound.Status, ErrTaskNotFound)
		return
	}

	if task.UserId != c.GetInt64("user-id") {
		c.Status(http.StatusForbidden)
		return
	}

	err = h.tasksService.DeleteById(c.Request.Context(), taskId)
	if err != nil {
		c.JSON(ErrFailedToDeleteTask.Status, ErrFailedToDeleteTask)
		return
	}

	c.Status(http.StatusOK)
}
