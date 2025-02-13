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

// TasksHandler handles task-related requests.
type TasksHandler struct {
	tasksService TasksService
}

// CreateTaskBody defines the request body for the /tasks endpoint.
type CreateTaskBody struct {
	Name        string `json:"name" example:"Eat"`                      // Name of the task
	Description string `json:"description" example:"Eat the pizza"`     // Description of the task
	Deadline    string `json:"deadline" example:"2023-12-31T23:59:59Z"` // Deadline for the task (RFC3339 format)
}

var (
	ErrInvalidDeadline       = appErrors.NewResponseError(http.StatusBadRequest, "failed to parse deadline")
	ErrFailedToCreateTask    = appErrors.NewResponseError(http.StatusBadRequest, "failed to create task")
	ErrInvalidTaskId         = appErrors.NewResponseError(http.StatusBadRequest, "task id is missing or invalid")
	ErrTaskNotFound          = appErrors.NewResponseError(http.StatusNotFound, "task was not found")
	ErrFailedToDeleteTask    = appErrors.NewResponseError(http.StatusInternalServerError, "failed to delete task")
	ErrFailedToRetrieveTasks = appErrors.NewResponseError(http.StatusInternalServerError, "failed to retrieve tasks")
)

// NewTasksHandler registers the task handler with the Gin engine.
func NewTasksHandler(r *gin.Engine, tasksService TasksService) {
	h := &TasksHandler{
		tasksService: tasksService,
	}

	r.GET("/tasks", middleware.AuthMiddleware, h.handleGetTasks)
	r.POST("/tasks", middleware.AuthMiddleware, h.handleCreateTask)
	r.DELETE("/tasks/:taskId", middleware.AuthMiddleware, h.handleDeleteTask)
}

// handleGetTasks retrieves all tasks for the authenticated user.
// @Summary Get all tasks for the current user
// @Description Retrieve a list of tasks for the currently authenticated user
// @Tags tasks
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} domain.Task "List of tasks"
// @Failure 404 {object} appErrors.ResponseError "User not found"
// @Failure 500 {object} appErrors.ResponseError "Failed to retrieve tasks"
// @Router /tasks [get]
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

// handleCreateTask creates a new task.
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body CreateTaskBody true "Task details"
// @Success 201 {object} domain.Task "Created task"
// @Failure 400 {object} appErrors.ResponseError "Invalid request body or deadline"
// @Failure 500 {object} appErrors.ResponseError "Failed to create task"
// @Router /tasks [post]
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

	userId := c.GetInt64("user-id")
	task, err := h.tasksService.Create(
		c.Request.Context(),
		data.Name,
		data.Description,
		deadline,
		userId,
	)
	if err != nil {
		c.JSON(ErrFailedToCreateTask.Status, ErrFailedToCreateTask)
		return
	}

	c.JSON(http.StatusCreated, task)
}

// handleDeleteTask deletes a task by ID.
// @Summary Delete a task
// @Description Delete a task by ID for the authenticated user
// @Tags tasks
// @Security ApiKeyAuth
// @Param taskId path int true "Task ID"
// @Success 200 "Task deleted successfully"
// @Failure 400 {object} appErrors.ResponseError "Invalid task ID"
// @Failure 404 {object} appErrors.ResponseError "Task not found"
// @Failure 403 "Forbidden if the task does not belong to the user"
// @Failure 500 {object} appErrors.ResponseError "Failed to delete task"
// @Router /tasks/{taskId} [delete]
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
