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

type installedOtelProviders struct {
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

func initOtel(c context.Context) (*installedOtelProviders, error) {
	otelProviders := installedOtelProviders{}

	// setup otel resource
	res, err := makeOtelResource(c)
	if err != nil {
		return &otelProviders, err
	}

	// configure otel exporters
	traceExporter, err := otxtrace.New(c)
	if err != nil {
		return &otelProviders, err
	}

	logExporter, err := otxlog.New(c)
	if err != nil {
		return &otelProviders, err
	}

	// configure otel sdks
	otelProviders.tracerProvider = otsdktrace.NewTracerProvider(
		otsdktrace.WithBatcher(traceExporter),
		otsdktrace.WithResource(res))
	otel.SetTracerProvider(otelProviders.tracerProvider)

	otelProviders.loggerProvider = otsdklog.NewLoggerProvider(
		otsdklog.WithProcessor(otsdklog.NewBatchProcessor(logExporter)),
		otsdklog.WithResource(res))
	otlogglobal.SetLoggerProvider(otelProviders.loggerProvider)

	return &otelProviders, nil
}

func shutdownOtel(c context.Context, otelProviders *installedOtelProviders) error {
	err := otelProviders.tracerProvider.Shutdown(c)
	if err != nil {
		return err
	}
	err = otelProviders.loggerProvider.Shutdown(c)
	if err != nil {
		return err
	}
	return nil
}
