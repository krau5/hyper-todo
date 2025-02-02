package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/internal/utils"
)

func AuthMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.VerifyJwt(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Set("user-email", sub)
	c.Next()
}
