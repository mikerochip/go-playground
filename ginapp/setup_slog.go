package main

import (
	"log/slog"
	"os"

	otslog "go.opentelemetry.io/contrib/bridges/otelslog"
	otlogglobal "go.opentelemetry.io/otel/log/global"
)

var otSlogLogger *slog.Logger

func initSlog() {
	otLoggerProvider := otlogglobal.GetLoggerProvider()
	otSlogLogger = otslog.NewLogger(serviceName, otslog.WithLoggerProvider(otLoggerProvider))

	traceSlogHandler := traceSlogHandler{
		stdoutHandler: slog.NewTextHandler(os.Stdout, nil),
		otHandler:     otSlogLogger.Handler(),
	}
	slog.SetDefault(slog.New(traceSlogHandler))
}
