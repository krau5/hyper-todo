package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/config"
	"github.com/krau5/hyper-todo/domain"
	appErrors "github.com/krau5/hyper-todo/internal/rest/errors"
	"github.com/krau5/hyper-todo/internal/utils"
	"gorm.io/gorm"
)

//go:generate mockery --name UsersService
type UsersService interface {
	Create(context context.Context, name, email, password string) error
	GetByEmail(context.Context, string) (domain.User, error)
	GetById(context.Context, int64) (domain.User, error)
}

type AuthHandler struct {
	usersService UsersService
}

type RegisterBody struct {
	Name     string `json:"name" binding:"required,min=4"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

var (
	ErrUserExists           = appErrors.NewResponseError(http.StatusConflict, "user with this email already exists")
	ErrUserNotFound         = appErrors.NewResponseError(http.StatusNotFound, "user was not found")
	ErrInvalidCredentials   = appErrors.NewResponseError(http.StatusBadRequest, "invalid email or password")
	ErrFailedToRetrieveUser = appErrors.NewResponseError(http.StatusInternalServerError, "failed to retrieve user")
	ErrFailedToCreateUser   = appErrors.NewResponseError(http.StatusInternalServerError, "failed to create user")
	ErrFailedToCreateToken  = appErrors.NewResponseError(http.StatusInternalServerError, "failed to create jwt token")
)

func NewAuthHandler(g *gin.Engine, usersService UsersService) {
	h := &AuthHandler{usersService: usersService}

	g.POST("/register", h.handleRegister)
	g.POST("/login", h.handleLogin)
}

func (h *AuthHandler) handleRegister(c *gin.Context) {
	var data RegisterBody

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(appErrors.ErrInvalidBody.Status, appErrors.ErrInvalidBody)
		return
	}

	err := h.usersService.Create(
		c.Request.Context(),
		data.Name,
		data.Email,
		data.Password,
	)

	if utils.IsErrDuplicatedKey(err) {
		c.JSON(ErrUserExists.Status, ErrUserExists)
		return
	}

	if err != nil {
		c.JSON(ErrFailedToCreateUser.Status, ErrFailedToCreateUser)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *AuthHandler) handleLogin(c *gin.Context) {
	var data LoginBody

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(appErrors.ErrInvalidBody.Status, appErrors.ErrInvalidBody)
		return
	}

	user, err := h.usersService.GetByEmail(c.Request.Context(), data.Email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(ErrUserNotFound.Status, ErrUserNotFound)
		return
	}

	if err != nil {
		c.JSON(ErrFailedToRetrieveUser.Status, ErrFailedToRetrieveUser)
		return
	}

	if ok := utils.VerifyPassword(data.Password, user.Password); !ok {
		c.JSON(ErrInvalidCredentials.Status, ErrInvalidCredentials)
		return
	}

	token, err := utils.CreateJwt(user.ID)
	if err != nil {
		c.JSON(ErrFailedToCreateToken.Status, ErrFailedToCreateToken)
		return
	}

	c.SetCookie("token", token, 3600, "/", config.Envs.CookieDomain, false, true)
	c.Status(http.StatusOK)
}
