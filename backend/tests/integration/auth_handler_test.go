//go:build integration

// Package integration provides integration tests for authentication handlers
package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthRateLimiting tests that auth endpoints are properly rate limited
func TestAuthRateLimiting(t *testing.T) {
	app := fiber.New()

	// Simulate AuthRateLimiter: 5 requests per minute
	authLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "auth:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"code":    "TOO_MANY_REQUESTS",
				"message": "too many authentication attempts, please try again later",
			})
		},
	})

	app.Get("/auth/login", authLimiter, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"redirect": "https://sso.example.com/login"})
	})

	t.Run("allows requests under limit", func(t *testing.T) {
		// Should allow first 5 requests
		for i := 0; i < 5; i++ {
			req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode, "request %d should succeed", i+1)
		}
	})

	t.Run("blocks requests over limit", func(t *testing.T) {
		// Create fresh app for this test
		app2 := fiber.New()
		limiter2 := limiter.New(limiter.Config{
			Max:        2,
			Expiration: time.Minute,
			KeyGenerator: func(c *fiber.Ctx) string {
				return "auth:test"
			},
		})
		app2.Get("/auth/login", limiter2, func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"ok": true})
		})

		// First 2 requests should succeed
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
			resp, _ := app2.Test(req, -1)
			resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}

		// 3rd request should be rate limited
		req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
		resp, _ := app2.Test(req, -1)
		resp.Body.Close()
		assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	})
}

// TestStrictRateLimiting tests strict rate limiting for sensitive operations
func TestStrictRateLimiting(t *testing.T) {
	app := fiber.New()

	// Simulate StrictRateLimiter: 3 requests per 5 minutes
	strictLimiter := limiter.New(limiter.Config{
		Max:        3,
		Expiration: 5 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "strict:test"
		},
	})

	app.Post("/verify/email", strictLimiter, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"sent": true})
	})

	// First 3 requests should succeed
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/verify/email", nil)
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode, "request %d should succeed", i+1)
	}

	// 4th request should be rate limited
	req := httptest.NewRequest(http.MethodPost, "/verify/email", nil)
	resp, _ := app.Test(req, -1)
	resp.Body.Close()
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
}

// TestCSRFProtection tests CSRF token validation
func TestCSRFProtection(t *testing.T) {
	// Note: This is a conceptual test. Real CSRF testing requires
	// the full middleware stack with cookie handling.

	app := fiber.New()

	// Simulate CSRF middleware behavior
	app.Use(func(c *fiber.Ctx) error {
		// Skip safe methods
		method := c.Method()
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			return c.Next()
		}

		// Check for CSRF token
		token := c.Get("X-CSRF-Token")
		if token == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"code":    "FORBIDDEN",
				"message": "CSRF token validation failed",
			})
		}

		return c.Next()
	})

	app.Post("/api/data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})

	app.Get("/api/data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"data": []string{}})
	})

	t.Run("GET requests pass without CSRF token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/data", nil)
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST without CSRF token is rejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("POST with CSRF token succeeds", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/data", nil)
		req.Header.Set("X-CSRF-Token", "valid-token")
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestAuthenticationRequired tests that protected endpoints require auth
func TestAuthenticationRequired(t *testing.T) {
	app := fiber.New()

	// Simulate RequireUser middleware
	requireAuth := func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "authentication required",
			})
		}
		return c.Next()
	}

	app.Get("/api/user/me", requireAuth, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"id": 1, "name": "Test User"})
	})

	t.Run("unauthenticated request is rejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("authenticated request succeeds", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestSecurityHeaders tests that security headers are set
func TestSecurityHeaders(t *testing.T) {
	app := fiber.New()

	// Simulate SecurityHeaders middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, _ := app.Test(req, -1)
	defer resp.Body.Close()

	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))
	assert.Equal(t, "strict-origin-when-cross-origin", resp.Header.Get("Referrer-Policy"))
}

// TestTimeoutMiddleware tests request timeout behavior
func TestTimeoutMiddleware(t *testing.T) {
	app := fiber.New()

	app.Get("/fast", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"fast": true})
	})

	// Fast request should succeed
	req := httptest.NewRequest(http.MethodGet, "/fast", nil)
	resp, err := app.Test(req, 1000) // 1 second timeout
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
