package main

import (
	"context"
	"log/slog"

	gin "github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())

	c := context.Background()
	slog.InfoContext(c, "Starting...")

	r.GET("/", func(c *gin.Context) {
		slog.InfoContext(c.Request.Context(), "hit /")
		c.JSON(200, gin.H{"message": "hello"})
	})

	r.Run()

	slog.InfoContext(c, "Shutting down...")
}
