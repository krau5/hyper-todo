package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/internal/rest/errors"
	"github.com/krau5/hyper-todo/internal/utils"
)

var (
	errMissingToken   = errors.NewResponseError(http.StatusUnauthorized, "missing or invalid token")
	errInvalidToken   = errors.NewResponseError(http.StatusUnauthorized, "invalid token")
	errExtractSubject = errors.NewResponseError(http.StatusBadRequest, "failed to extract subject from token")
	errParseUserID    = errors.NewResponseError(http.StatusBadRequest, "failed to parse user ID from token")
)

func validateToken(c *gin.Context) (int64, *errors.ResponseError) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		return 0, errMissingToken
	}

	token, err := utils.VerifyJwt(tokenString)
	if err != nil {
		return 0, errInvalidToken
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return 0, errExtractSubject
	}

	userId, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return 0, errParseUserID
	}

	return userId, nil
}

func AuthMiddleware(c *gin.Context) {
	userId, err := validateToken(c)
	if err != nil {
		c.JSON(err.Status, err)
		c.Abort()
		return
	}

	c.Set("user-id", userId)
	c.Next()
}
