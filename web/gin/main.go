package main

import (
	"github.com/gin-gonic/gin"
)

func InitMux() *gin.Engine {
	router := gin.Default()

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.Writer.WriteString("Hello world!")
	})

	return router
}

func main() {
	router := InitMux()
	router.Run(":8080")
}
