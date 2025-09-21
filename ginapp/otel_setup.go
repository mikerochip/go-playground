package main

import (
	"context"

	"go.opentelemetry.io/otel"
	otlogglobal "go.opentelemetry.io/otel/log/global"
	otsdklog "go.opentelemetry.io/otel/sdk/log"
	otsdkresource "go.opentelemetry.io/otel/sdk/resource"
	otsdktrace "go.opentelemetry.io/otel/sdk/trace"
	otsemconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	// OTLP exporters
	otxlog "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	otxtrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

type otelSetup struct {
	tracerProvider *otsdktrace.TracerProvider
	loggerProvider *otsdklog.LoggerProvider
}

func makeOtelResource(c context.Context) (*otsdkresource.Resource, error) {
	return otsdkresource.New(c,
		otsdkresource.WithSchemaURL(otsemconv.SchemaURL),
		otsdkresource.WithFromEnv(),      // OTEL_RESOURCE_ATTRIBUTES, OTEL_SERVICE_NAME, etc.
		otsdkresource.WithTelemetrySDK(), // telemetry.sdk.*
		otsdkresource.WithHost(),         // host.*
		otsdkresource.WithOS(),           // os.*
		otsdkresource.WithProcess(),      // process.*
		otsdkresource.WithContainer(),    // container.*
	)
}

func configureOtel(c context.Context) (*otelSetup, error) {
	otelSetup := otelSetup{}

	// setup otel resource
	res, err := makeOtelResource(c)
	if err != nil {
		return &otelSetup, err
	}

	// configure otel exporters
	traceExporter, err := otxtrace.New(c)
	if err != nil {
		return &otelSetup, err
	}

	logExporter, err := otxlog.New(c)
	if err != nil {
		return &otelSetup, err
	}

	// configure otel sdks
	otelSetup.tracerProvider = otsdktrace.NewTracerProvider(
		otsdktrace.WithBatcher(traceExporter),
		otsdktrace.WithResource(res))
	otel.SetTracerProvider(otelSetup.tracerProvider)

	otelSetup.loggerProvider = otsdklog.NewLoggerProvider(
		otsdklog.WithProcessor(otsdklog.NewBatchProcessor(logExporter)),
		otsdklog.WithResource(res))
	otlogglobal.SetLoggerProvider(otelSetup.loggerProvider)

	return &otelSetup, nil
}

func shutdownOtel(c context.Context, otelSetup *otelSetup) error {
	err := otelSetup.tracerProvider.Shutdown(c)
	if err != nil {
		return err
	}
	err = otelSetup.loggerProvider.Shutdown(c)
	if err != nil {
		return err
	}
	return nil
}
