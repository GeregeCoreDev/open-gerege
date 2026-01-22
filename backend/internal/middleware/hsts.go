// Package middleware provides implementation for middleware
//
// File: hsts.go
// Description: HTTP Strict Transport Security middleware
package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// HSTSConfig holds HSTS middleware configuration
type HSTSConfig struct {
	// Enabled controls whether HSTS header is added
	Enabled bool
	// MaxAge is the time in seconds that the browser should remember HTTPS-only
	MaxAge int
	// IncludeSubDomains applies HSTS to all subdomains
	IncludeSubDomains bool
	// Preload allows submission to browser HSTS preload lists
	Preload bool
}

// DefaultHSTSConfig returns production-ready HSTS settings
//
// MaxAge: 1 year (31536000 seconds) - OWASP recommended minimum
// IncludeSubDomains: true - protects all subdomains
// Preload: true - enables browser preload list submission
func DefaultHSTSConfig() HSTSConfig {
	return HSTSConfig{
		Enabled:           true,
		MaxAge:            31536000, // 1 year
		IncludeSubDomains: true,
		Preload:           true,
	}
}

// HSTS returns an HTTP Strict Transport Security middleware
//
// The HSTS header tells browsers to only access the site via HTTPS.
// This prevents:
//   - Protocol downgrade attacks
//   - Cookie hijacking
//   - Man-in-the-middle attacks
//
// Usage:
//
//	app.Use(middleware.HSTS(middleware.DefaultHSTSConfig()))
//
// Response header:
//
//	Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
func HSTS(cfg HSTSConfig) fiber.Handler {
	if !cfg.Enabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// Pre-compute header value
	headerValue := fmt.Sprintf("max-age=%d", cfg.MaxAge)
	if cfg.IncludeSubDomains {
		headerValue += "; includeSubDomains"
	}
	if cfg.Preload {
		headerValue += "; preload"
	}

	return func(c *fiber.Ctx) error {
		c.Set("Strict-Transport-Security", headerValue)
		return c.Next()
	}
}
