// Package middleware provides implementation for middleware
//
// File: limiter.go
// Description: implementation for middleware
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package middleware нь HTTP middleware-уудыг агуулна.

Энэ файл нь rate limiting middleware-ийг тодорхойлно.
DDoS халдлага, brute force халдлагаас хамгаална.

Rate Limiting:
  - Sliding window algorithm
  - Per-user (authenticated) эсвэл per-IP (anonymous)
  - Configurable max requests болон window

Ашиглалт:

	// 1 минутад 100 request
	app.Use(middleware.RateLimiter(100, time.Minute))

	// Тодорхой route-д хатуу хязгаар
	app.Post("/login", middleware.RateLimiter(5, time.Minute), handler.Login)
*/
package middleware

import (
	"fmt"  // String formatting
	"time" // Duration

	"git.gerege.mn/backend-packages/sso-client" // Session ID авах

	"github.com/gofiber/fiber/v2"                    // Web framework
	"github.com/gofiber/fiber/v2/middleware/limiter" // Rate limiter middleware
)

// ============================================================
// RATE LIMITER
// ============================================================

// RateLimiter нь sliding window rate limiter middleware буцаана.
//
// Algorithm: Sliding window
//   - Window: Тодорхой хугацаа (жишээ: 1 минут)
//   - Max: Window-д зөвшөөрөгдөх хамгийн их request тоо
//   - Хэтэрвэл 429 Too Many Requests буцаана
//
// Key generation:
//   - Authenticated user: "sid:{session_id}"
//   - Anonymous user: "ip:{ip_address}"
//
// Parameters:
//   - max: Window-д зөвшөөрөгдөх хамгийн их request тоо
//   - window: Time window (жишээ: time.Minute, 10*time.Second)
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Response (хэтэрсэн бол):
//   - 429 Too Many Requests
//   - X-RateLimit-Limit header
//   - X-RateLimit-Remaining header
//   - X-RateLimit-Reset header
//
// Жишээ:
//
//	// Global: 1 минутад 100 request
//	app.Use(middleware.RateLimiter(100, time.Minute))
//
//	// Login: 1 минутад 5 request (brute force хамгаалалт)
//	app.Post("/login", middleware.RateLimiter(5, time.Minute), handler.Login)
//
//	// API: 10 секундад 50 request
//	api.Use(middleware.RateLimiter(50, 10*time.Second))
func RateLimiter(max int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		// Хамгийн их request тоо
		Max: max,

		// Time window
		Expiration: window,

		// Key үүсгэх функц
		// Authenticated бол session ID, anonymous бол IP ашиглана
		KeyGenerator: func(c *fiber.Ctx) string {
			// Session ID авах (authenticated user)
			if sid := ssoclient.GetSessionID(c); sid != "" {
				return fmt.Sprintf("sid:%s", sid)
			}
			// IP address (anonymous user)
			return "ip:" + c.IP()
		},
	})
}

// ============================================================
// ENDPOINT-SPECIFIC RATE LIMITERS
// ============================================================

// AuthRateLimiter returns a strict rate limiter for authentication endpoints.
// Protects against brute force attacks on login, password reset, etc.
//
// Default: 5 requests per minute per IP
//
// Usage:
//
//	auth.Post("/login", middleware.AuthRateLimiter(), handler.Login)
//	auth.Post("/reset-password", middleware.AuthRateLimiter(), handler.ResetPassword)
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "auth:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(
				fiber.StatusTooManyRequests,
				"too many authentication attempts, please try again later",
			)
		},
	})
}

// APIRateLimiter returns a moderate rate limiter for general API endpoints.
// Authenticated users get higher limits than anonymous users.
//
// Default: 100 requests per minute
//
// Usage:
//
//	api := app.Group("/api", middleware.APIRateLimiter())
func APIRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if sid := ssoclient.GetSessionID(c); sid != "" {
				return "api:sid:" + sid
			}
			return "api:ip:" + c.IP()
		},
	})
}

// StrictRateLimiter returns a very strict rate limiter for sensitive operations.
// Use for password changes, account deletion, etc.
//
// Default: 3 requests per 5 minutes per user/IP
//
// Usage:
//
//	user.Post("/change-password", middleware.StrictRateLimiter(), handler.ChangePassword)
func StrictRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        3,
		Expiration: 5 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if sid := ssoclient.GetSessionID(c); sid != "" {
				return "strict:sid:" + sid
			}
			return "strict:ip:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(
				fiber.StatusTooManyRequests,
				"rate limit exceeded for sensitive operation",
			)
		},
	})
}
