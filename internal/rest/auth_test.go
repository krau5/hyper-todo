package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/krau5/hyper-todo/internal/rest/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	name := "user"
	email := "user@example.com"
	password := "password123"

	usersService := mocks.NewUsersService(t)
	usersService.On("Create", mock.Anything, name, email, password).Return(nil)

	r := gin.New()
	NewAuthHandler(r, usersService)

	body := RegisterBody{
		Name:     name,
		Email:    email,
		Password: password,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", &buf)
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestRegisterError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	name := "user"
	email := "user@example.com"
	password := "password123"

	usersService := mocks.NewUsersService(t)
	usersService.On("Create", mock.Anything, name, email, password).Return(mockDuplicatedError())

	r := gin.New()
	NewAuthHandler(r, usersService)

	body := RegisterBody{
		Name:     name,
		Email:    email,
		Password: password,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", &buf)
	r.ServeHTTP(w, req)

	expectedBody, _ := json.Marshal(ErrUserExists)

	assert.Equal(t, 409, w.Code)
	assert.Equal(t, string(expectedBody), w.Body.String())
}

func mockDuplicatedError() *pgconn.PgError {
	return &pgconn.PgError{
		Code:    "23505",
		Message: "duplicate key violates unique constraint",
	}
}
