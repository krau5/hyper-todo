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

func TestMeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var userId int64 = 1
	mockUser := domain.User{Name: "user", Email: "user@example.com"}

	usersService := mocks.NewUsersService(t)
	usersService.On("GetById", mock.Anything, userId).Return(mockUser, nil)

	h := &UsersHandler{usersService: usersService}
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("user-id", userId)
		c.Next()
	})
	r.GET("/me", h.handleMe)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	r.ServeHTTP(w, req)

	expectedBody, _ := json.Marshal(mockUser)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedBody), w.Body.String())
}

func TestMeHandlerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var userId int64 = 1
	usersService := mocks.NewUsersService(t)
	usersService.On("GetById", mock.Anything, userId).Return(domain.User{}, gorm.ErrRecordNotFound)

	h := &UsersHandler{usersService: usersService}
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("user-id", userId)
		c.Next()
	})
	r.GET("/me", h.handleMe)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me", nil)
	r.ServeHTTP(w, req)

	expectedBody, _ := json.Marshal(ErrUserNotFound)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, string(expectedBody), w.Body.String())
}
