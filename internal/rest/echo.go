package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler struct{}

func NewPingHandler(g *gin.Engine) {
	h := &PingHandler{}

	g.GET("/ping", h.handlePing)
}

func (h *PingHandler) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
