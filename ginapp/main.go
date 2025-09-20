package main

import (
	"context"
	"log/slog"

	gin "github.com/gin-gonic/gin"
	otslog "go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	otstdoutlog "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	otstdouttrace "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	otlogglobal "go.opentelemetry.io/otel/log/global"
	otlog "go.opentelemetry.io/otel/sdk/log"
	otresource "go.opentelemetry.io/otel/sdk/resource"
	ottrace "go.opentelemetry.io/otel/sdk/trace"
)

const serviceName string = "ginapp"

func setupOtel() error {
	res, err := otresource.Merge(otresource.Default(), otresource.Empty())
	if err != nil {
		return err
	}

	stdExporter, err := otstdouttrace.New(otstdouttrace.WithPrettyPrint())
	if err != nil {
		return err
	}

	tracerProvider := ottrace.NewTracerProvider(ottrace.WithBatcher(stdExporter), ottrace.WithResource(res))
	otel.SetTracerProvider(tracerProvider)

	logExporter, err := otstdoutlog.New(otstdoutlog.WithPrettyPrint())
	if err != nil {
		return err
	}
	loggerProvider := otlog.NewLoggerProvider(otlog.WithProcessor(otlog.NewBatchProcessor(logExporter)), otlog.WithResource(res))
	otlogglobal.SetLoggerProvider(loggerProvider)

	// slog → OTel bridge (sends slog records into OTel logs; includes trace/span context)
	slog.SetDefault(otslog.NewLogger(serviceName, otslog.WithLoggerProvider(loggerProvider)))

	return nil
}

func main() {
	setupOtel()

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
