// Package middleware provides implementation for middleware
//
// File: csrf.go
// Description: CSRF protection middleware
package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

// CSRFConfig holds CSRF middleware configuration
type CSRFConfig struct {
	// Enabled controls whether CSRF protection is active
	Enabled bool
	// CookieSecure sets the Secure flag on the CSRF cookie (true for HTTPS)
	CookieSecure bool
	// CookieSameSite sets SameSite attribute (Strict recommended)
	CookieSameSite string
	// Expiration is the CSRF token lifetime
	Expiration time.Duration
	// SkipPaths are paths that skip CSRF validation (e.g., webhooks)
	SkipPaths []string
}

// DefaultCSRFConfig returns production-ready CSRF settings
func DefaultCSRFConfig(isProduction bool) CSRFConfig {
	return CSRFConfig{
		Enabled:        true,
		CookieSecure:   isProduction,
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/docs",
			"/auth/callback", // OAuth callbacks
			"/auth/local/",   // Local auth API (stateless)
		},
	}
}

// CSRF returns a CSRF protection middleware
//
// Headers:
//   - X-CSRF-Token: Client sends this header with the token from cookie
//
// Cookie:
//   - csrf_token: Contains the CSRF token
//
// Usage:
//
//	app.Use(middleware.CSRF(middleware.DefaultCSRFConfig(true)))
//
// Client-side:
//  1. Read csrf_token from cookie
//  2. Send X-CSRF-Token header with requests
func CSRF(cfg CSRFConfig) fiber.Handler {
	if !cfg.Enabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// Build skip paths map for O(1) lookup
	skipMap := make(map[string]bool, len(cfg.SkipPaths))
	for _, p := range cfg.SkipPaths {
		skipMap[p] = true
	}

	return csrf.New(csrf.Config{
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "csrf_token",
		CookieSameSite: cfg.CookieSameSite,
		CookieSecure:   cfg.CookieSecure,
		CookieHTTPOnly: true,
		Expiration:     cfg.Expiration,
		// Skip GET, HEAD, OPTIONS (safe methods)
		// Skip specific paths
		Next: func(c *fiber.Ctx) bool {
			// Safe methods don't need CSRF
			method := c.Method()
			if method == fiber.MethodGet || method == fiber.MethodHead || method == fiber.MethodOptions {
				return true
			}
			// Skip configured paths
			path := c.Path()
			if skipMap[path] {
				return true
			}
			// Skip paths starting with skip prefixes
			for p := range skipMap {
				if len(path) >= len(p) && path[:len(p)] == p {
					return true
				}
			}
			return false
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusForbidden, "CSRF token validation failed")
		},
	})
}
