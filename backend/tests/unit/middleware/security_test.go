package middleware_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"templatev25/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPSRedirect(t *testing.T) {
	tests := []struct {
		name           string
		enabled        bool
		protocol       string
		expectedStatus int
		shouldRedirect bool
	}{
		{
			name:           "disabled - no redirect",
			enabled:        false,
			protocol:       "http",
			expectedStatus: 200,
			shouldRedirect: false,
		},
		{
			name:           "enabled with http - should redirect",
			enabled:        true,
			protocol:       "http",
			expectedStatus: 301,
			shouldRedirect: true,
		},
		{
			name:           "enabled with https - no redirect",
			enabled:        true,
			protocol:       "https",
			expectedStatus: 200,
			shouldRedirect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(middleware.HTTPSRedirect(tt.enabled))
			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-Forwarded-Proto", tt.protocol)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.shouldRedirect {
				location := resp.Header.Get("Location")
				assert.True(t, strings.HasPrefix(location, "https://"))
			}
		})
	}
}

func TestSecurityHeaders(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.SecurityHeaders())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	app.Get("/docs/swagger", func(c *fiber.Ctx) error {
		return c.SendString("Swagger")
	})

	t.Run("basic security headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
		assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
		assert.Equal(t, "no-referrer", resp.Header.Get("Referrer-Policy"))
		assert.NotEmpty(t, resp.Header.Get("Content-Security-Policy"))
		assert.NotEmpty(t, resp.Header.Get("Permissions-Policy"))
	})

	t.Run("swagger docs has relaxed CSP", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/docs/swagger", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		csp := resp.Header.Get("Content-Security-Policy")
		assert.Contains(t, csp, "validator.swagger.io")
		assert.Contains(t, csp, "fonts.googleapis.com")
	})
}

func TestBodySizeLimit(t *testing.T) {
	tests := []struct {
		name           string
		maxBytes       int
		bodySize       int
		expectedStatus int
	}{
		{
			name:           "body within limit",
			maxBytes:       1024,
			bodySize:       100,
			expectedStatus: 200,
		},
		{
			name:           "body at limit",
			maxBytes:       100,
			bodySize:       100,
			expectedStatus: 200,
		},
		{
			name:           "body exceeds limit",
			maxBytes:       100,
			bodySize:       200,
			expectedStatus: 413,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(middleware.BodySizeLimit(tt.maxBytes))
			app.Post("/", func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})

			body := strings.NewReader(strings.Repeat("x", tt.bodySize))
			req := httptest.NewRequest("POST", "/", body)
			req.Header.Set("Content-Type", "text/plain")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestPaginationLimit(t *testing.T) {
	tests := []struct {
		name           string
		maxSize        int
		query          string
		expectedStatus int
	}{
		{
			name:           "size within limit",
			maxSize:        100,
			query:          "?size=50",
			expectedStatus: 200,
		},
		{
			name:           "size at limit",
			maxSize:        100,
			query:          "?size=100",
			expectedStatus: 200,
		},
		{
			name:           "size exceeds limit",
			maxSize:        100,
			query:          "?size=200",
			expectedStatus: 400,
		},
		{
			name:           "pageSize parameter",
			maxSize:        100,
			query:          "?pageSize=50",
			expectedStatus: 200,
		},
		{
			name:           "negative size",
			maxSize:        100,
			query:          "?size=-1",
			expectedStatus: 400,
		},
		{
			name:           "negative page",
			maxSize:        100,
			query:          "?page=-1",
			expectedStatus: 400,
		},
		{
			name:           "no pagination params",
			maxSize:        100,
			query:          "",
			expectedStatus: 200,
		},
		{
			name:           "default max size",
			maxSize:        0, // Use default
			query:          "?size=50",
			expectedStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			if tt.maxSize > 0 {
				app.Use(middleware.PaginationLimit(tt.maxSize))
			} else {
				app.Use(middleware.PaginationLimit())
			}
			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})

			req := httptest.NewRequest("GET", "/"+tt.query, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestTimeout(t *testing.T) {
	t.Run("context has timeout", func(t *testing.T) {
		app := fiber.New()
		app.Use(middleware.Timeout(5 * time.Second))
		app.Get("/", func(c *fiber.Ctx) error {
			ctx := c.UserContext()
			deadline, ok := ctx.Deadline()
			if !ok {
				return fiber.NewError(500, "no deadline set")
			}
			// Deadline should be within 5 seconds from now
			remaining := time.Until(deadline)
			if remaining > 5*time.Second || remaining < 0 {
				return fiber.NewError(500, "invalid deadline")
			}
			return c.SendString("OK")
		})

		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, 200, resp.StatusCode, string(body))
	})
}

func TestDefaultConstants(t *testing.T) {
	assert.Equal(t, 100, middleware.DefaultMaxPageSize)
	assert.Equal(t, 1, middleware.DefaultMinPageSize)
}
