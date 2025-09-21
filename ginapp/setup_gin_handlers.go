package main

import (
	"log/slog"
	"net/http"

	gin "github.com/gin-gonic/gin"
	otattr "go.opentelemetry.io/otel/attribute"
	ottrace "go.opentelemetry.io/otel/trace"
)

func initGinHandlers(engine *gin.Engine) {
	engine.Handle(http.MethodGet, "/", func(ctx *gin.Context) {
		requestCtx := ctx.Request.Context()

		slog.InfoContext(requestCtx, "hit /")

		tracer := ottrace.SpanFromContext(requestCtx).TracerProvider().Tracer(serviceName)
		_, childSpan := tracer.Start(requestCtx, "child_test")
		{
			childSpan.SetAttributes(otattr.String("foo", "bar"))
			slog.InfoContext(requestCtx, "******child span log1", slog.String("stuff", "thing"))
		}
		childSpan.End()

		ctx.JSON(200, gin.H{"message": "hello"})
	})
}
