// Package telemetry provides implementation for telemetry
//
// File: metrics.go
// Description: Business and technical metrics using OpenTelemetry
//
// This file defines metrics for monitoring application health and business KPIs.
//
// Metrics categories:
//   - HTTP: Request count, latency, error rates
//   - Database: Query latency, connection pool
//   - Cache: Hit/miss ratios
//   - Business: Logins, user actions, etc.
package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Metrics holds all application metrics
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal   metric.Int64Counter
	HTTPRequestDuration metric.Float64Histogram
	HTTPActiveRequests  metric.Int64UpDownCounter

	// Database metrics
	DBQueryDuration   metric.Float64Histogram
	DBConnectionsOpen metric.Int64UpDownCounter
	DBQueryErrors     metric.Int64Counter

	// Cache metrics
	CacheHits   metric.Int64Counter
	CacheMisses metric.Int64Counter
	CacheSize   metric.Int64UpDownCounter

	// Business metrics
	UserLogins        metric.Int64Counter
	UserLogouts       metric.Int64Counter
	UserRegistrations metric.Int64Counter
	RoleChanges       metric.Int64Counter
	PermissionChecks  metric.Int64Counter
	OrgSwitches       metric.Int64Counter
	APICallsByUser    metric.Int64Counter

	// System metrics
	GoroutineCount metric.Int64UpDownCounter
	MemoryUsage    metric.Int64UpDownCounter

	meter metric.Meter
}

// NewMetrics creates and registers all metrics
func NewMetrics(ctx context.Context) (*Metrics, *prometheus.Exporter, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, nil, err
	}

	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	otel.SetMeterProvider(provider)

	meter := provider.Meter("templatev25")
	m := &Metrics{meter: meter}

	// Initialize HTTP metrics
	m.HTTPRequestsTotal, _ = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("1"),
	)

	m.HTTPRequestDuration, _ = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
	)

	m.HTTPActiveRequests, _ = meter.Int64UpDownCounter(
		"http_active_requests",
		metric.WithDescription("Number of active HTTP requests"),
		metric.WithUnit("1"),
	)

	// Initialize Database metrics
	m.DBQueryDuration, _ = meter.Float64Histogram(
		"db_query_duration_seconds",
		metric.WithDescription("Database query duration in seconds"),
		metric.WithUnit("s"),
	)

	m.DBConnectionsOpen, _ = meter.Int64UpDownCounter(
		"db_connections_open",
		metric.WithDescription("Number of open database connections"),
		metric.WithUnit("1"),
	)

	m.DBQueryErrors, _ = meter.Int64Counter(
		"db_query_errors_total",
		metric.WithDescription("Total number of database query errors"),
		metric.WithUnit("1"),
	)

	// Initialize Cache metrics
	m.CacheHits, _ = meter.Int64Counter(
		"cache_hits_total",
		metric.WithDescription("Total number of cache hits"),
		metric.WithUnit("1"),
	)

	m.CacheMisses, _ = meter.Int64Counter(
		"cache_misses_total",
		metric.WithDescription("Total number of cache misses"),
		metric.WithUnit("1"),
	)

	m.CacheSize, _ = meter.Int64UpDownCounter(
		"cache_size",
		metric.WithDescription("Current cache size"),
		metric.WithUnit("1"),
	)

	// Initialize Business metrics
	m.UserLogins, _ = meter.Int64Counter(
		"user_logins_total",
		metric.WithDescription("Total number of user logins"),
		metric.WithUnit("1"),
	)

	m.UserLogouts, _ = meter.Int64Counter(
		"user_logouts_total",
		metric.WithDescription("Total number of user logouts"),
		metric.WithUnit("1"),
	)

	m.UserRegistrations, _ = meter.Int64Counter(
		"user_registrations_total",
		metric.WithDescription("Total number of user registrations"),
		metric.WithUnit("1"),
	)

	m.RoleChanges, _ = meter.Int64Counter(
		"role_changes_total",
		metric.WithDescription("Total number of role changes"),
		metric.WithUnit("1"),
	)

	m.PermissionChecks, _ = meter.Int64Counter(
		"permission_checks_total",
		metric.WithDescription("Total number of permission checks"),
		metric.WithUnit("1"),
	)

	m.OrgSwitches, _ = meter.Int64Counter(
		"org_switches_total",
		metric.WithDescription("Total number of organization switches"),
		metric.WithUnit("1"),
	)

	m.APICallsByUser, _ = meter.Int64Counter(
		"api_calls_by_user_total",
		metric.WithDescription("Total API calls by user"),
		metric.WithUnit("1"),
	)

	// Initialize System metrics
	m.GoroutineCount, _ = meter.Int64UpDownCounter(
		"goroutine_count",
		metric.WithDescription("Current number of goroutines"),
		metric.WithUnit("1"),
	)

	m.MemoryUsage, _ = meter.Int64UpDownCounter(
		"memory_usage_bytes",
		metric.WithDescription("Current memory usage in bytes"),
		metric.WithUnit("By"),
	)

	return m, exporter, nil
}

// RecordHTTPRequest records an HTTP request metric
func (m *Metrics) RecordHTTPRequest(ctx context.Context, method, path string, statusCode int, durationSec float64) {
	attrs := []attribute.KeyValue{
		attribute.String("method", method),
		attribute.String("path", path),
		attribute.Int("status_code", statusCode),
		attribute.String("status_class", statusClass(statusCode)),
	}
	m.HTTPRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.HTTPRequestDuration.Record(ctx, durationSec, metric.WithAttributes(attrs...))
}

// RecordDBQuery records a database query metric
func (m *Metrics) RecordDBQuery(ctx context.Context, operation, table string, durationSec float64, err error) {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
		attribute.String("table", table),
	}
	m.DBQueryDuration.Record(ctx, durationSec, metric.WithAttributes(attrs...))
	if err != nil {
		m.DBQueryErrors.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

// RecordCacheHit records a cache hit
func (m *Metrics) RecordCacheHit(ctx context.Context, cache string) {
	attrs := []attribute.KeyValue{
		attribute.String("cache", cache),
	}
	m.CacheHits.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordCacheMiss records a cache miss
func (m *Metrics) RecordCacheMiss(ctx context.Context, cache string) {
	attrs := []attribute.KeyValue{
		attribute.String("cache", cache),
	}
	m.CacheMisses.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordUserLogin records a user login event
func (m *Metrics) RecordUserLogin(ctx context.Context, provider string) {
	attrs := []attribute.KeyValue{
		attribute.String("provider", provider),
	}
	m.UserLogins.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordUserLogout records a user logout event
func (m *Metrics) RecordUserLogout(ctx context.Context) {
	m.UserLogouts.Add(ctx, 1)
}

// RecordUserRegistration records a user registration event
func (m *Metrics) RecordUserRegistration(ctx context.Context, source string) {
	attrs := []attribute.KeyValue{
		attribute.String("source", source),
	}
	m.UserRegistrations.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordRoleChange records a role change event
func (m *Metrics) RecordRoleChange(ctx context.Context, action string) {
	attrs := []attribute.KeyValue{
		attribute.String("action", action),
	}
	m.RoleChanges.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordPermissionCheck records a permission check event
func (m *Metrics) RecordPermissionCheck(ctx context.Context, permission string, granted bool) {
	attrs := []attribute.KeyValue{
		attribute.String("permission", permission),
		attribute.Bool("granted", granted),
	}
	m.PermissionChecks.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordOrgSwitch records an organization switch event
func (m *Metrics) RecordOrgSwitch(ctx context.Context) {
	m.OrgSwitches.Add(ctx, 1)
}

// RecordAPICall records an API call by user
func (m *Metrics) RecordAPICall(ctx context.Context, userID int, endpoint string) {
	attrs := []attribute.KeyValue{
		attribute.Int("user_id", userID),
		attribute.String("endpoint", endpoint),
	}
	m.APICallsByUser.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// statusClass returns the status class (2xx, 3xx, 4xx, 5xx)
func statusClass(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "2xx"
	case code >= 300 && code < 400:
		return "3xx"
	case code >= 400 && code < 500:
		return "4xx"
	case code >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}
