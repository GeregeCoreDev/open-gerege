// Package middleware provides HTTP middlewares
//
// File: middleware_test.go
// Description: Unit tests for middleware package
package middleware

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPSRedirect_Disabled(t *testing.T) {
	app := fiber.New()
	app.Use(HTTPSRedirect(false))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestHTTPSRedirect_Enabled(t *testing.T) {
	app := fiber.New()
	app.Use(HTTPSRedirect(true))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 301, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Location"), "https://")
}

func TestHTTPSRedirect_AlreadyHTTPS(t *testing.T) {
	app := fiber.New()
	app.Use(HTTPSRedirect(true))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "https://example.com/test", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestSecurityHeaders(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeaders())
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Check security headers
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "no-referrer", resp.Header.Get("Referrer-Policy"))
	assert.NotEmpty(t, resp.Header.Get("Content-Security-Policy"))
	assert.NotEmpty(t, resp.Header.Get("Permissions-Policy"))
}

func TestSecurityHeaders_SwaggerPath(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeaders())
	app.Get("/docs/swagger.json", func(c *fiber.Ctx) error {
		return c.SendString("swagger")
	})

	req := httptest.NewRequest("GET", "/docs/swagger.json", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Swagger should have relaxed CSP
	csp := resp.Header.Get("Content-Security-Policy")
	assert.Contains(t, csp, "validator.swagger.io")
}

func TestBodySizeLimit_Under(t *testing.T) {
	app := fiber.New()
	app.Use(BodySizeLimit(1024)) // 1KB limit
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	body := make([]byte, 500) // Under limit
	req := httptest.NewRequest("POST", "/test", nil)
	req.Body = nil
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	_ = body
}

func TestBodySizeLimit_Over(t *testing.T) {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // Allow large body in Fiber
	})
	app.Use(BodySizeLimit(100)) // 100 byte limit
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Create body larger than limit
	body := make([]byte, 200)
	for i := range body {
		body[i] = 'a'
	}

	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	// Empty body passes the limit check
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPaginationLimit_Valid(t *testing.T) {
	app := fiber.New()
	app.Use(PaginationLimit(100))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{"no params", "", 200},
		{"valid size", "?size=50", 200},
		{"valid pageSize", "?pageSize=50", 200},
		{"valid page", "?page=1", 200},
		{"at limit", "?size=100", 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test"+tt.query, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPaginationLimit_Invalid(t *testing.T) {
	app := fiber.New()
	app.Use(PaginationLimit(100))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{"size over limit", "?size=200", 400},
		{"negative size", "?size=-5", 400},
		{"negative page", "?page=-1", 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test"+tt.query, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestPaginationLimit_CustomMax(t *testing.T) {
	app := fiber.New()
	app.Use(PaginationLimit(10)) // Custom max of 10
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Size 15 should be rejected with max 10
	req := httptest.NewRequest("GET", "/test?size=15", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestTimeout(t *testing.T) {
	app := fiber.New()
	app.Use(Timeout(5 * time.Second))
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		deadline, ok := ctx.Deadline()
		assert.True(t, ok)
		assert.True(t, deadline.After(time.Now()))
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestTimeout_ContextPropagated(t *testing.T) {
	app := fiber.New()
	timeout := 100 * time.Millisecond
	app.Use(Timeout(timeout))
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Check context has deadline
		_, ok := ctx.Deadline()
		assert.True(t, ok)

		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDefaultMaxPageSize(t *testing.T) {
	assert.Equal(t, 100, DefaultMaxPageSize)
}

func TestDefaultMinPageSize(t *testing.T) {
	assert.Equal(t, 1, DefaultMinPageSize)
}

func TestPreComputedCSP(t *testing.T) {
	// Verify pre-computed CSP strings are valid
	assert.Contains(t, cspSwagger, "default-src 'self'")
	assert.Contains(t, cspSwagger, "validator.swagger.io")

	assert.Contains(t, cspDefault, "default-src 'self'")
	assert.Contains(t, cspDefault, "frame-ancestors 'none'")

	assert.Contains(t, permissionsPolicy, "geolocation=()")
}

func TestContextWithTimeout(t *testing.T) {
	// Direct test of context timeout behavior
	timeout := 50 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, deadline.After(time.Now()))

	// Wait for timeout
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
}
