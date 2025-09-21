package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const defaultPort string = "8080"
const serviceName string = "ginapp"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Kill)
	defer stop()

	srv, otelProviders, err := initialize(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "setup failed: %s\n", err.Error())
		panic(err)
	}

	slog.InfoContext(ctx, "Starting...")

	go run(srv)
	<-ctx.Done()

	slog.InfoContext(ctx, "Shutting down...")

	shutdown(ctx, srv, otelProviders)
}

func initialize(ctx context.Context) (*http.Server, *installedOtelProviders, error) {
	providers, err := initOtel(ctx)
	if err != nil {
		return nil, nil, err
	}

	initSlog()

	ginEngine := initGin(serviceName)
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: ginEngine,
	}

	initGinHandlers(ginEngine)

	return srv, providers, nil
}

func run(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "err", err)
	}
}

func shutdown(ctx context.Context, srv *http.Server, otelProviders *installedOtelProviders) {
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "err", err)
		panic(err)
	}

	shutdownOtel(ctx, otelProviders)
}
