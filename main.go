package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(gin.Recovery())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.Writer.Write([]byte("pong"))
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
