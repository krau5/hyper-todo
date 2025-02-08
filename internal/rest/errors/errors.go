package errors

import "fmt"

type ResponseError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Status, e.Message)
}
