// Package middleware provides implementation for middleware
//
// File: tracing.go
// Description: OpenTelemetry distributed tracing middleware
//
// This middleware adds distributed tracing to HTTP requests.
// It creates spans for each request and propagates trace context.
//
// Span attributes:
//   - http.method: HTTP method (GET, POST, etc.)
//   - http.url: Request URL
//   - http.route: Route pattern
//   - http.status_code: Response status code
//   - http.user_agent: User-Agent header
//   - net.peer.ip: Client IP address
//
// Usage:
//
//	app.Use(middleware.Tracing())
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "fiber-http"
)

// TracingConfig holds tracing middleware configuration
type TracingConfig struct {
	// TracerName is the name of the tracer
	TracerName string
	// SkipPaths are paths that skip tracing
	SkipPaths []string
}

// DefaultTracingConfig returns default tracing configuration
func DefaultTracingConfig() TracingConfig {
	return TracingConfig{
		TracerName: tracerName,
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/favicon.ico",
		},
	}
}

// Tracing returns a tracing middleware with default configuration
func Tracing() fiber.Handler {
	return TracingWithConfig(DefaultTracingConfig())
}

// TracingWithConfig returns a tracing middleware with custom configuration
func TracingWithConfig(cfg TracingConfig) fiber.Handler {
	tracer := otel.Tracer(cfg.TracerName)

	// Build skip paths map for O(1) lookup
	skipMap := make(map[string]bool, len(cfg.SkipPaths))
	for _, p := range cfg.SkipPaths {
		skipMap[p] = true
	}

	return func(c *fiber.Ctx) error {
		// Skip tracing for configured paths
		if skipMap[c.Path()] {
			return c.Next()
		}

		// Extract trace context from incoming request headers
		ctx := otel.GetTextMapPropagator().Extract(
			c.UserContext(),
			propagation.HeaderCarrier(c.GetReqHeaders()),
		)

		// Create span name from method and route
		spanName := c.Method() + " " + c.Route().Path
		if c.Route().Path == "" {
			spanName = c.Method() + " " + c.Path()
		}

		// Start span
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(c.Method()),
				semconv.HTTPURLKey.String(c.OriginalURL()),
				semconv.HTTPRouteKey.String(c.Route().Path),
				semconv.UserAgentOriginalKey.String(c.Get("User-Agent")),
				attribute.String("net.peer.ip", c.IP()),
				semconv.HTTPSchemeKey.String(c.Protocol()),
			),
		)
		defer span.End()

		// Set context with span
		c.SetUserContext(ctx)

		// Add trace ID to response headers for debugging
		if span.SpanContext().HasTraceID() {
			c.Set("X-Trace-ID", span.SpanContext().TraceID().String())
		}

		// Execute next handler
		err := c.Next()

		// Record response status
		status := c.Response().StatusCode()
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(status))

		// Set span status based on HTTP status code
		if status >= 500 {
			span.SetStatus(codes.Error, "server error")
		} else if status >= 400 {
			span.SetStatus(codes.Error, "client error")
		} else {
			span.SetStatus(codes.Ok, "")
		}

		// Record error if any
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		return err
	}
}

// SpanFromContext returns the current span from context
// Use this in handlers to add custom attributes or events
//
// Example:
//
//	span := middleware.SpanFromContext(c)
//	span.SetAttributes(attribute.String("user.id", userID))
//	span.AddEvent("user_action", trace.WithAttributes(
//	    attribute.String("action", "login"),
//	))
func SpanFromContext(c *fiber.Ctx) trace.Span {
	return trace.SpanFromContext(c.UserContext())
}

// AddSpanAttribute adds an attribute to the current span
//
// Example:
//
//	middleware.AddSpanAttribute(c, "user.id", userID)
func AddSpanAttribute(c *fiber.Ctx, key string, value interface{}) {
	span := trace.SpanFromContext(c.UserContext())
	switch v := value.(type) {
	case string:
		span.SetAttributes(attribute.String(key, v))
	case int:
		span.SetAttributes(attribute.Int(key, v))
	case int64:
		span.SetAttributes(attribute.Int64(key, v))
	case float64:
		span.SetAttributes(attribute.Float64(key, v))
	case bool:
		span.SetAttributes(attribute.Bool(key, v))
	}
}

// AddSpanEvent adds an event to the current span
//
// Example:
//
//	middleware.AddSpanEvent(c, "user_logged_in")
func AddSpanEvent(c *fiber.Ctx, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(c.UserContext())
	span.AddEvent(name, trace.WithAttributes(attrs...))
}
