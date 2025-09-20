package main

import (
	"context"
	"log/slog"

	otslog "go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	otstdoutlog "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	otstdouttrace "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	otsdklogglobal "go.opentelemetry.io/otel/log/global"
	otsdklog "go.opentelemetry.io/otel/sdk/log"
	otsdkresource "go.opentelemetry.io/otel/sdk/resource"
	otsdktrace "go.opentelemetry.io/otel/sdk/trace"
	otsemconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func makeOtelResource(c context.Context) (*otsdkresource.Resource, error) {
	// Compose custom + built-in detectors
	res, err := otsdkresource.New(c,
		otsdkresource.WithSchemaURL(otsemconv.SchemaURL),
		otsdkresource.WithFromEnv(),      // OTEL_RESOURCE_ATTRIBUTES, OTEL_SERVICE_NAME, etc.
		otsdkresource.WithTelemetrySDK(), // telemetry.sdk.*
		otsdkresource.WithHost(),         // host.*
		otsdkresource.WithOS(),           // os.*
		otsdkresource.WithProcess(),      // process.*
		otsdkresource.WithContainer(),    // container.*
	)
	return res, err
}

func setupOtel(c context.Context) error {
	//res, err := otsdkresource.Merge(otsdkresource.Default(), otsdkresource.Empty())
	res, err := makeOtelResource(c)
	if err != nil {
		return err
	}

	stdExporter, err := otstdouttrace.New(otstdouttrace.WithPrettyPrint())
	if err != nil {
		return err
	}

	tracerProvider := otsdktrace.NewTracerProvider(otsdktrace.WithBatcher(stdExporter), otsdktrace.WithResource(res))
	otel.SetTracerProvider(tracerProvider)

	logExporter, err := otstdoutlog.New(otstdoutlog.WithPrettyPrint())
	if err != nil {
		return err
	}
	loggerProvider := otsdklog.NewLoggerProvider(otsdklog.WithProcessor(otsdklog.NewBatchProcessor(logExporter)), otsdklog.WithResource(res))
	otsdklogglobal.SetLoggerProvider(loggerProvider)

	// slog -> OTel bridge (sends slog records into OTel logs; includes trace/span context)
	slog.SetDefault(otslog.NewLogger(serviceName, otslog.WithLoggerProvider(loggerProvider)))

	return nil
}
