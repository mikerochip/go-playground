package main

import (
	"context"
	"log/slog"
	"os"

	gin "github.com/gin-gonic/gin"
	otgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	otattr "go.opentelemetry.io/otel/attribute"
	ottrace "go.opentelemetry.io/otel/trace"
)

const serviceName string = "ginapp"

func main() {
	c := context.Background()

	otelSetup, err := configureOtel(c)
	if err != nil {
		_, _ = os.Stderr.WriteString("otel setup failed: " + err.Error() + "\n")
		os.Exit(1)
	}

	configureSlog()

	r := gin.Default()
	r.Use(otgin.Middleware(serviceName))

	slog.InfoContext(c, "Starting...")

	r.GET("/", func(ctx *gin.Context) {
		requestCtx := ctx.Request.Context()
		tracer := ottrace.SpanFromContext(requestCtx).TracerProvider()
		testTracer := tracer.Tracer(serviceName)

		slog.InfoContext(requestCtx, "hit /")

		_, childSpan := testTracer.Start(requestCtx, "child_test")
		childSpan.SetAttributes(otattr.String("foo", "bar"))

		slog.InfoContext(requestCtx, "******child span log1", slog.String("stuff", "thing"))

		childSpan.End()

		ctx.JSON(200, gin.H{"message": "hello"})
	})

	r.Run()

	slog.InfoContext(c, "Shutting down...")
	shutdownOtel(c, otelSetup)
}
