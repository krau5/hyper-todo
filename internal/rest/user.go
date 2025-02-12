package rest

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/krau5/hyper-todo/internal/rest/errors"
	"github.com/krau5/hyper-todo/internal/rest/middleware"
	"gorm.io/gorm"
)

// UsersHandler handles user-related requests.
type UsersHandler struct {
	usersService UsersService
}

// NewUsersHandler registers the user handler with the Gin engine.
func NewUsersHandler(r *gin.Engine, usersService UsersService) {
	h := &UsersHandler{usersService: usersService}

	r.GET("/me", middleware.AuthMiddleware, h.handleMe)
}

// handleMe retrieves details of the currently authenticated user.
// @Summary Get current user details
// @Description Retrieve details of the currently authenticated user
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} domain.User "User details"
// @Failure 404 {object} errors.ResponseError "User not found"
// @Failure 401 {object} errors.ResponseError "Unauthorized"
// @Router /me [get]
func (h *UsersHandler) handleMe(c *gin.Context) {
	userId := c.GetInt64("user-id")
	user, err := h.usersService.GetById(c.Request.Context(), userId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(ErrUserNotFound.Status, ErrUserNotFound)
		return
	}

	c.JSON(http.StatusOK, user)
}
