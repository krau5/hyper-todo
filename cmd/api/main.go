package main

import (
	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/internal/rest"
)

func main() {
	r := gin.Default()

	rest.NewPingHandler(r)

	r.Run()
}
