package main

import (
	"context"
	"log/slog"
	"os"

	otslog "go.opentelemetry.io/contrib/bridges/otelslog"
	otlogglobal "go.opentelemetry.io/otel/log/global"
	ottrace "go.opentelemetry.io/otel/trace"
)

type traceAwareSlogHandler struct {
	stdoutHandler slog.Handler
	otHandler     slog.Handler
}

var otSlogLogger *slog.Logger

func configureSlog() {
	// configure slog to send otel logs to the spec otlp port while default logs go to
	// stdout. reason for this setup is so that app stdout logs are clean. collector can
	// send logs to stdout if configured to do so
	otLoggerProvider := otlogglobal.GetLoggerProvider()
	otSlogLogger = otslog.NewLogger(serviceName, otslog.WithLoggerProvider(otLoggerProvider))

	slog.SetDefault(slog.New(traceAwareSlogHandler{
		stdoutHandler: slog.NewTextHandler(os.Stdout, nil),
		otHandler:     otSlogLogger.Handler()}))
}

func (t traceAwareSlogHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return t.stdoutHandler.Enabled(ctx, lvl) || t.otHandler.Enabled(ctx, lvl)
}
func (t traceAwareSlogHandler) WithAttrs(a []slog.Attr) slog.Handler {
	return traceAwareSlogHandler{
		stdoutHandler: t.stdoutHandler.WithAttrs(a),
		otHandler:     t.otHandler.WithAttrs(a)}
}
func (t traceAwareSlogHandler) WithGroup(g string) slog.Handler {
	return traceAwareSlogHandler{
		stdoutHandler: t.stdoutHandler.WithGroup(g),
		otHandler:     t.otHandler.WithGroup(g)}
}
func (t traceAwareSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	spanCtx := ottrace.SpanFromContext(ctx).SpanContext()

	if spanCtx.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)
	}

	err := t.stdoutHandler.Handle(ctx, r)
	if err != nil {
		return err
	}

	if !spanCtx.IsValid() {
		return nil
	}
	return t.otHandler.Handle(ctx, r)
}
