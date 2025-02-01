package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type UsersService interface {
	Create(context context.Context, name, email, password string) error
}

type UsersHandler struct {
	usersService UsersService
}

type CreateUserBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(g *gin.Engine, usersService UsersService) {
	h := &UsersHandler{usersService: usersService}

	g.POST("/register", h.handleCreate)
}

func (h *UsersHandler) handleCreate(c *gin.Context) {
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

	if isDublicatedKeyErr(err) {
		c.JSON(http.StatusConflict, gin.H{"error": "user with this email already exists"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func isDublicatedKeyErr(err error) bool {
	var perr *pgconn.PgError
	errors.As(err, &perr)

	return perr.Code == "23505"
}
