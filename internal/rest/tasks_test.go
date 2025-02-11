package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/internal/rest/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateTaskHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	name := "eat"
	description := "eat the pizza"
	rawDeadline := "2025-01-01T22:22:22.220Z"
	var userId int64 = 1

	deadline, err := time.Parse(time.RFC3339, rawDeadline)
	if err != nil {
		t.Fatalf("Failed to parse deadline: %v", err)
	}

	mockTask := domain.Task{
		ID:          1,
		Name:        name,
		Description: description,
		Deadline:    deadline,
		Completed:   false,
		UserId:      userId,
	}

	tasksService := mocks.NewTasksService(t)
	tasksService.On("Create", mock.Anything, name, description, deadline, userId).Return(mockTask, nil)

	h := &TasksHandler{
		tasksService: tasksService,
	}

	r := gin.New()
	r.POST("/tasks", h.handleCreateTask)

	body := CreateTaskBody{
		Name:        name,
		Description: description,
		Deadline:    rawDeadline,
		UserId:      userId,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", &buf)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	expectedBody, _ := json.Marshal(mockTask)
	assert.Equal(t, string(expectedBody), w.Body.String())
}

func TestGetTasksHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var userId int64 = 1
	mockTasks := []domain.Task{
		{
			ID:          1,
			Name:        "eat",
			Description: "eat the pizza",
			Deadline:    time.Now(),
		},
		{
			ID:          2,
			Name:        "drink",
			Description: "drink the coke",
			Deadline:    time.Now(),
		},
	}

	tasksService := mocks.NewTasksService(t)
	tasksService.On("GetByUser", mock.Anything, userId).Return(mockTasks, nil)

	h := &TasksHandler{
		tasksService: tasksService,
	}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user-id", userId)
		c.Next()
	})
	r.GET("/tasks", h.handleGetTasks)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody, _ := json.Marshal(mockTasks)
	assert.Equal(t, string(expectedBody), w.Body.String())
}

func TestGetTasksHandler_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var userId int64 = 1

	tasksService := mocks.NewTasksService(t)
	tasksService.On("GetByUser", mock.Anything, userId).Return([]domain.Task{}, gorm.ErrRecordNotFound)

	h := &TasksHandler{
		tasksService: tasksService,
	}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user-id", userId)
		c.Next()
	})
	r.GET("/tasks", h.handleGetTasks)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, ErrUserNotFound.Status, w.Code)

	expectedBody, _ := json.Marshal(ErrUserNotFound)
	assert.Equal(t, string(expectedBody), w.Body.String())
}

func TestGetTasksHandler_FailedToRetrieveTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var userId int64 = 1

	tasksService := mocks.NewTasksService(t)
	tasksService.On("GetByUser", mock.Anything, userId).Return([]domain.Task{}, assert.AnError)

	h := &TasksHandler{
		tasksService: tasksService,
	}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user-id", userId)
		c.Next()
	})
	r.GET("/tasks", h.handleGetTasks)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, ErrFailedToRetrieveTasks.Status, w.Code)

	expectedBody, _ := json.Marshal(ErrFailedToRetrieveTasks)
	assert.Equal(t, string(expectedBody), w.Body.String())
}
