package errors

import (
	"fmt"
	"net/http"
)

// ResponseError defines a standard error response.
// @name ResponseError
type ResponseError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Status, e.Message)
}

func NewResponseError(status int, message string) *ResponseError {
	return &ResponseError{Status: status, Message: message}
}

var (
	ErrInvalidBody = NewResponseError(http.StatusBadRequest, "invalid request body")
)
