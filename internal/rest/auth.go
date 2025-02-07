package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/config"
	"github.com/krau5/hyper-todo/domain"
	"github.com/krau5/hyper-todo/internal/utils"
	"gorm.io/gorm"
)

//go:generate mockery --name UsersService
type UsersService interface {
	Create(context context.Context, name, email, password string) error
	GetByEmail(context.Context, string) (domain.User, error)
}

type AuthHandler struct {
	usersService UsersService
}

type RegisterBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(g *gin.Engine, usersService UsersService) {
	h := &AuthHandler{usersService: usersService}

	g.POST("/register", h.handleRegister)
	g.POST("/login", h.handleLogin)
}

func (h *AuthHandler) handleRegister(c *gin.Context) {
	var data RegisterBody

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

func (h *AuthHandler) handleLogin(c *gin.Context) {
	var data LoginBody

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.usersService.GetByEmail(c.Request.Context(), data.Email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user with this email does not exist"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if ok := utils.VerifyPassword(data.Password, user.Password); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	token, err := utils.CreateJwt(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 3600, "/", config.Envs.CookieDomain, false, true)
	c.Status(http.StatusOK)
}
