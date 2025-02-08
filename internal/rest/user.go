package rest

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/internal/rest/middleware"
	"gorm.io/gorm"
)

type UsersHandler struct {
	usersService UsersService
}

func NewUsersHandler(r *gin.Engine, usersService UsersService) {
	h := &UsersHandler{usersService: usersService}

	r.GET("/me", middleware.AuthMiddleware, h.handleMe)
}

func (h *UsersHandler) handleMe(c *gin.Context) {
	userId := c.GetInt64("user-id")
	user, err := h.usersService.GetById(c.Request.Context(), userId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(ErrUserNotFound.Status, ErrUserNotFound)
		return
	}

	c.JSON(http.StatusOK, user)
}
