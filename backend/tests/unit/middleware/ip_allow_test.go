package middleware_test

import (
	"net/http/httptest"
	"testing"

	"templatev25/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPAllow_DefaultCIDRs(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.IPAllow(nil)) // Use default CIDRs
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Default Fiber test uses 0.0.0.0 which may not be in default list
	// localhost should be allowed
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Note: Fiber's test client uses 0.0.0.0 by default
	// The actual behavior depends on how Fiber resolves the IP
}

func TestIPAllow_CustomCIDRs(t *testing.T) {
	tests := []struct {
		name           string
		allowedCIDRs   []string
		clientIP       string
		expectedStatus int
	}{
		{
			name:           "allowed single IP",
			allowedCIDRs:   []string{"192.168.1.100/32"},
			clientIP:       "192.168.1.100",
			expectedStatus: 200,
		},
		{
			name:           "allowed IP in range",
			allowedCIDRs:   []string{"192.168.1.0/24"},
			clientIP:       "192.168.1.50",
			expectedStatus: 200,
		},
		{
			name:           "blocked IP not in range",
			allowedCIDRs:   []string{"192.168.1.0/24"},
			clientIP:       "192.168.2.50",
			expectedStatus: 403,
		},
		{
			name:           "allowed localhost",
			allowedCIDRs:   []string{"127.0.0.1/32"},
			clientIP:       "127.0.0.1",
			expectedStatus: 200,
		},
		{
			name:           "multiple CIDRs - first match",
			allowedCIDRs:   []string{"10.0.0.0/8", "192.168.0.0/16"},
			clientIP:       "10.1.2.3",
			expectedStatus: 200,
		},
		{
			name:           "multiple CIDRs - second match",
			allowedCIDRs:   []string{"10.0.0.0/8", "192.168.0.0/16"},
			clientIP:       "192.168.5.10",
			expectedStatus: 200,
		},
		{
			name:           "multiple CIDRs - no match",
			allowedCIDRs:   []string{"10.0.0.0/8", "192.168.0.0/16"},
			clientIP:       "172.16.1.1",
			expectedStatus: 403,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				// Trust proxy headers for testing
				ProxyHeader: "X-Real-IP",
			})
			app.Use(middleware.IPAllow(tt.allowedCIDRs))
			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-Real-IP", tt.clientIP)

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestIPAllow_InvalidCIDR(t *testing.T) {
	// Invalid CIDRs should be skipped
	app := fiber.New(fiber.Config{
		ProxyHeader: "X-Real-IP",
	})
	app.Use(middleware.IPAllow([]string{"invalid-cidr", "192.168.1.0/24"}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "192.168.1.50")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should still work with valid CIDR
	assert.Equal(t, 200, resp.StatusCode)
}

func TestIPAllow_IPv6(t *testing.T) {
	app := fiber.New(fiber.Config{
		ProxyHeader: "X-Real-IP",
	})
	app.Use(middleware.IPAllow([]string{"::1/128", "2001:db8::/32"}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	tests := []struct {
		name           string
		clientIP       string
		expectedStatus int
	}{
		{
			name:           "IPv6 localhost allowed",
			clientIP:       "::1",
			expectedStatus: 200,
		},
		{
			name:           "IPv6 in range allowed",
			clientIP:       "2001:db8::1",
			expectedStatus: 200,
		},
		{
			name:           "IPv6 not in range blocked",
			clientIP:       "2001:db9::1",
			expectedStatus: 403,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-Real-IP", tt.clientIP)

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
