package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/internal/rest/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var userId int64 = 1

func TestMeHandler(t *testing.T) {
	mockUser := domain.User{Name: "user", Email: "user@example.com"}

	r, usersService := setupTest(t)
	usersService.On("GetById", mock.Anything, userId).Return(mockUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	r.ServeHTTP(w, req)

	expectedBody, _ := json.Marshal(mockUser)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedBody), w.Body.String())
}

func TestMeHandlerError(t *testing.T) {
	r, usersService := setupTest(t)
	usersService.On("GetById", mock.Anything, userId).Return(domain.User{}, gorm.ErrRecordNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	r.ServeHTTP(w, req)

	expectedBody, _ := json.Marshal(ErrUserNotFound)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, string(expectedBody), w.Body.String())
}

func setupTest(t *testing.T) (*gin.Engine, *mocks.UsersService) {
	gin.SetMode(gin.TestMode)

	usersService := mocks.NewUsersService(t)
	h := &UsersHandler{usersService: usersService}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user-id", userId)
		c.Next()
	})
	r.GET("/me", h.handleMe)

	return r, usersService
}
