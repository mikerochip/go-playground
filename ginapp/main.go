package main

import (
	"context"
	"log/slog"

	gin "github.com/gin-gonic/gin"
	otgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	ottrace "go.opentelemetry.io/otel/trace"
)

const serviceName string = "ginapp"

func main() {
	c := context.Background()
	setupOtel(c)

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(otgin.Middleware(serviceName))
	r.Use(makeTestTracerMiddleware())

	slog.InfoContext(c, "Starting...")

	r.GET("/", func(ginCtx *gin.Context) {
		testTracer := ginCtx.MustGet(testTracerKey).(ottrace.Tracer)
		c := ginCtx.Request.Context()

		slog.InfoContext(c, "hit /")

		_, childSpan := testTracer.Start(c, "child_test")
		childSpan.SetAttributes(attribute.String("foo", "bar"))
		slog.InfoContext(c, "child span", slog.String("span_id", childSpan.SpanContext().SpanID().String()))
		childSpan.End()

		ginCtx.JSON(200, gin.H{"message": "hello"})
	})

	r.Run()

	slog.InfoContext(c, "Shutting down...")
}
