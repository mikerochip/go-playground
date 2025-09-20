package main

import (
	gin "github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

const testTracerKey string = "test-tracer"

func makeTestTracerMiddleware() gin.HandlerFunc {
	tracer := otel.Tracer("test")

	return func(c *gin.Context) {
		c.Set(testTracerKey, tracer)
		c.Next()
	}
}
