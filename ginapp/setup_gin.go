package main

import (
	gin "github.com/gin-gonic/gin"
	otgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func initGin(serviceName string) *gin.Engine {
	engine := gin.Default()
	engine.Use(otgin.Middleware(serviceName))

	return engine
}
