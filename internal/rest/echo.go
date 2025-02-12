package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler handles ping requests.
type PingHandler struct{}

// PingResponse defines the response structure for the /ping endpoint.
type PingResponse struct {
	Message string `json:"message"`
}

// NewPingHandler registers the ping handler with the Gin engine.
func NewPingHandler(g *gin.Engine) {
	h := &PingHandler{}

	g.GET("/ping", h.handlePing)
}

// handlePing responds with a "pong" message.
// @Summary Ping the server
// @Description Get a "pong" response from the server
// @Tags Echo
// @Accept json
// @Produce json
// @Success 200 {object} PingResponse "Returns a pong message"
// @Router /ping [get]
func (h *PingHandler) handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, PingResponse{Message: "pong"})
}
