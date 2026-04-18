// Package tracing wires the OpenTelemetry tracer provider for the service.
package tracing

import (
	"context"
	"fmt"
	"reference-service-go/internal/build"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

var shutdown func(context.Context) error

// Setup installs a global tracer provider and W3C trace-context propagator.
// When enabled is false, propagation still works but no exporter is registered,
// so there is no shutdown stall when no collector is reachable.
func Setup(ctx context.Context, enabled bool, endpoint string) error {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(build.ServiceName),
			semconv.ServiceVersion(build.Version()),
		),
	)
	if err != nil {
		return fmt.Errorf("create resource: %w", err)
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	}

	if enabled {
		exp, expErr := otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
		)
		if expErr != nil {
			return fmt.Errorf("create OTLP exporter: %w", expErr)
		}

		opts = append(opts, sdktrace.WithBatcher(exp))
	}

	tp := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	shutdown = tp.Shutdown
	if !enabled {
		shutdown = func(context.Context) error { return nil }
	}

	return nil
}

// Shutdown flushes pending spans and releases exporter resources.
func Shutdown(ctx context.Context) error {
	if shutdown == nil {
		return nil
	}

	return shutdown(ctx)
}
