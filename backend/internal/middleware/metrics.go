// Package middleware provides implementation for middleware
//
// File: metrics.go
// Description: HTTP metrics middleware for request tracking
//
// This middleware records HTTP request metrics including:
//   - Request count by method, path, and status
//   - Request duration histogram
//   - Active request gauge
package middleware

import (
	"time"

	"templatev25/internal/telemetry"

	"github.com/gofiber/fiber/v2"
)

// MetricsConfig holds metrics middleware configuration
type MetricsConfig struct {
	// SkipPaths are paths that skip metrics recording
	SkipPaths []string
}

// DefaultMetricsConfig returns default metrics configuration
func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/favicon.ico",
		},
	}
}

// HTTPMetrics returns a metrics middleware with default configuration
func HTTPMetrics(m *telemetry.Metrics) fiber.Handler {
	return HTTPMetricsWithConfig(m, DefaultMetricsConfig())
}

// HTTPMetricsWithConfig returns a metrics middleware with custom configuration
func HTTPMetricsWithConfig(m *telemetry.Metrics, cfg MetricsConfig) fiber.Handler {
	if m == nil {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// Build skip paths map for O(1) lookup
	skipMap := make(map[string]bool, len(cfg.SkipPaths))
	for _, p := range cfg.SkipPaths {
		skipMap[p] = true
	}

	return func(c *fiber.Ctx) error {
		// Skip metrics for configured paths
		if skipMap[c.Path()] {
			return c.Next()
		}

		ctx := c.UserContext()
		start := time.Now()

		// Track active requests
		m.HTTPActiveRequests.Add(ctx, 1)
		defer m.HTTPActiveRequests.Add(ctx, -1)

		// Execute next handler
		err := c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		path := c.Route().Path
		if path == "" {
			path = c.Path()
		}
		m.RecordHTTPRequest(ctx, c.Method(), path, c.Response().StatusCode(), duration)

		return err
	}
}
