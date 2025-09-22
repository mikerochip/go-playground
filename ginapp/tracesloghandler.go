package main

import (
	"context"
	"log/slog"

	ottrace "go.opentelemetry.io/otel/trace"
)

// This slog.Handler always logs to stdout, but only logs to an otelslog handler iff
// inside a span. Reason for this is so that app stdout logs are clean e.g. in kube.
// The assumption is that a collector agent can be configured to handle the verbose
// otel stuff.
type traceSlogHandler struct {
	stdoutHandler slog.Handler
	otHandler     slog.Handler
}

func (t traceSlogHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return t.stdoutHandler.Enabled(ctx, lvl) || t.otHandler.Enabled(ctx, lvl)
}

func (t traceSlogHandler) WithAttrs(a []slog.Attr) slog.Handler {
	return traceSlogHandler{
		stdoutHandler: t.stdoutHandler.WithAttrs(a),
		otHandler:     t.otHandler.WithAttrs(a)}
}

func (t traceSlogHandler) WithGroup(g string) slog.Handler {
	return traceSlogHandler{
		stdoutHandler: t.stdoutHandler.WithGroup(g),
		otHandler:     t.otHandler.WithGroup(g)}
}

func (t traceSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	spanCtx := ottrace.SpanContextFromContext(ctx)

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
