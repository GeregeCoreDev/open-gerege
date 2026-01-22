// Package telemetry provides OpenTelemetry observability setup
//
// File: tracer.go
// Description: Distributed tracing configuration
//
// This package configures OpenTelemetry tracing for distributed tracing
// across services. It supports multiple exporters:
//   - OTLP (Jaeger, Tempo, etc.)
//   - Stdout (for development/debugging)
//
// Usage:
//
//	shutdown, err := telemetry.InitTracer(ctx, cfg, serviceName, serviceVersion)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer shutdown(ctx)
package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// TracerConfig holds configuration for the tracer
type TracerConfig struct {
	// Enabled controls whether tracing is active
	Enabled bool
	// Endpoint is the OTLP collector endpoint (e.g., "localhost:4317")
	Endpoint string
	// Insecure disables TLS for OTLP connection
	Insecure bool
	// SampleRate is the sampling rate (0.0 to 1.0)
	SampleRate float64
	// UseStdout enables stdout exporter (for development)
	UseStdout bool
}

// DefaultTracerConfig returns sensible defaults
func DefaultTracerConfig() TracerConfig {
	return TracerConfig{
		Enabled:    true,
		Endpoint:   "localhost:4317",
		Insecure:   true,
		SampleRate: 1.0, // Sample all traces in dev
		UseStdout:  false,
	}
}

// ShutdownFunc is the function to call for graceful shutdown
type ShutdownFunc func(context.Context) error

// InitTracer initializes OpenTelemetry tracing
//
// Parameters:
//   - ctx: Context for initialization
//   - cfg: Tracer configuration
//   - serviceName: Name of this service
//   - serviceVersion: Version of this service
//
// Returns:
//   - ShutdownFunc: Function to call for graceful shutdown
//   - error: Initialization error
//
// Example:
//
//	shutdown, err := telemetry.InitTracer(ctx, telemetry.TracerConfig{
//	    Enabled:    true,
//	    Endpoint:   "jaeger:4317",
//	    Insecure:   true,
//	    SampleRate: 0.1, // Sample 10% in production
//	}, "my-service", "1.0.0")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer shutdown(ctx)
func InitTracer(ctx context.Context, cfg TracerConfig, serviceName, serviceVersion string) (ShutdownFunc, error) {
	if !cfg.Enabled {
		// Return no-op shutdown function
		return func(context.Context) error { return nil }, nil
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
			attribute.String("environment", "production"),
		),
		resource.WithHost(),
		resource.WithProcess(),
	)
	if err != nil {
		return nil, err
	}

	// Create exporter based on configuration
	var exporter sdktrace.SpanExporter
	if cfg.UseStdout {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	} else {
		// OTLP exporter options
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(cfg.Endpoint),
		}
		if cfg.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}

		exporter, err = otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
	}
	if err != nil {
		return nil, err
	}

	// Create sampler based on sample rate
	var sampler sdktrace.Sampler
	if cfg.SampleRate >= 1.0 {
		sampler = sdktrace.AlwaysSample()
	} else if cfg.SampleRate <= 0 {
		sampler = sdktrace.NeverSample()
	} else {
		sampler = sdktrace.TraceIDRatioBased(cfg.SampleRate)
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(512),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global text map propagator (W3C Trace Context + Baggage)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Return shutdown function
	return tp.Shutdown, nil
}

// Tracer returns a named tracer from the global provider
func Tracer(name string) interface{ /* otel trace.Tracer */ } {
	return otel.Tracer(name)
}
