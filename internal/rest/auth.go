package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/internal/utils"
)

//go:generate mockery --name UsersService
type UsersService interface {
	Create(context context.Context, name, email, password string) error
}

type AuthHandler struct {
	usersService UsersService
}

type CreateUserBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(g *gin.Engine, usersService UsersService) {
	h := &AuthHandler{usersService: usersService}

	g.POST("/register", h.handleRegister)
}

func (h *AuthHandler) handleRegister(c *gin.Context) {
	var data CreateUserBody

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.usersService.Create(
		c.Request.Context(),
		data.Name,
		data.Email,
		data.Password,
	)

	if utils.IsErrDuplicatedKey(err) {
		c.JSON(http.StatusConflict, gin.H{"error": "user with this email already exists"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
