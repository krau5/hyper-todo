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
)

func TestCreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	name := "eat"
	description := "eat the pizza"
	rawDeadline := "2025-01-01T22:22:22.220Z"
	var userId int64 = 1

	deadline, err := time.Parse(time.RFC3339, rawDeadline)
	if err != nil {
		t.Error(err)
	}
	mockTask := domain.Task{
		Name:        name,
		Description: description,
		Deadline:    deadline,
	}
	tasksService := mocks.NewTasksService(t)
	tasksService.On("Create", mock.Anything, name, description, mock.Anything, userId).Return(mockTask, nil)

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

	assert.Equal(t, 201, w.Code)
}
