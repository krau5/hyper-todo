package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	Message string `json:"message"`
}

func TestPingHandler(t *testing.T) {
	r := gin.Default()
	NewPingHandler(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	body := Response{Message: "pong"}
	expectedBody, _ := json.Marshal(body)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, string(expectedBody), w.Body.String())
}
