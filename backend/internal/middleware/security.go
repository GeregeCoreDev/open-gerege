// Package middleware provides implementation for middleware
//
// File: security.go
// Description: implementation for middleware
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package middleware нь HTTP middleware-уудыг агуулна.

Энэ файл нь security-тэй холбоотой middleware-уудыг тодорхойлно:
  - SecurityHeaders: Аюулгүй байдлын HTTP header-үүд
  - BodySizeLimit: Request body хэмжээний хязгаар
  - Timeout: Request timeout

Security headers:
  - X-Content-Type-Options: MIME sniffing хамгаалалт
  - X-Frame-Options: Clickjacking хамгаалалт
  - Content-Security-Policy: XSS, injection хамгаалалт
  - Referrer-Policy: Referrer мэдээлэл хязгаарлах
  - Permissions-Policy: Browser features хязгаарлах
*/
package middleware

import (
	"context" // Context with timeout
	"fmt"     // Format strings
	"strings" // String operations
	"time"    // Duration

	"github.com/gofiber/fiber/v2" // Web framework
)

// ============================================================
// PRE-COMPUTED SECURITY HEADERS (Performance optimization)
// ============================================================

// Pre-computed CSP headers (computed once at startup)
var (
	cspSwagger = strings.Join([]string{
		"default-src 'self'",
		"img-src 'self' data: https://validator.swagger.io",
		"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
		"font-src 'self' https://fonts.gstatic.com",
		"script-src 'self' 'unsafe-inline'",
		"connect-src 'self'",
	}, "; ")

	cspDefault = strings.Join([]string{
		"default-src 'self'",
		"img-src 'self' data:",
		"style-src 'self' 'unsafe-inline'",
		"script-src 'self'",
		"connect-src 'self'",
		"frame-ancestors 'none'",
	}, "; ")

	permissionsPolicy = "geolocation=(), microphone=(), camera=()"
)

// ============================================================
// HTTPS REDIRECT
// ============================================================

// HTTPSRedirect нь HTTP request-ийг HTTPS руу redirect хийх middleware буцаана.
// Production орчинд HTTPS enforce хийхэд ашиглана.
//
// Parameters:
//   - enabled: Redirect идэвхтэй эсэх
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Ашиглалт:
//
//	app.Use(middleware.HTTPSRedirect(cfg.Server.ForceHTTPS))
func HTTPSRedirect(enabled bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !enabled {
			return c.Next()
		}

		// X-Forwarded-Proto header шалгах (reverse proxy-ийн ард байвал)
		proto := c.Get("X-Forwarded-Proto", "")
		if proto == "" {
			proto = c.Protocol()
		}

		// HTTPS биш бол redirect хийх
		if proto != "https" {
			// 301 Permanent Redirect
			return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), fiber.StatusMovedPermanently)
		}

		return c.Next()
	}
}

// ============================================================
// SECURITY HEADERS
// ============================================================

// SecurityHeaders нь аюулгүй байдлын HTTP header-үүд тохируулах middleware буцаана.
//
// Тохируулах header-үүд:
//
//	X-Content-Type-Options: nosniff
//	  - Browser-ийн MIME type sniffing-ийг хориглоно
//	  - Content-Type header-д итгэхийг шаардана
//
//	X-Frame-Options: DENY
//	  - Энэ хуудсыг iframe-д оруулахыг хориглоно
//	  - Clickjacking халдлагаас хамгаална
//
//	Referrer-Policy: no-referrer
//	  - Бусад сайт руу Referrer header илгээхгүй
//	  - Хэрэглэгчийн privacy хамгаална
//
//	Content-Security-Policy (CSP):
//	  - XSS, data injection халдлагаас хамгаална
//	  - Swagger docs-д тусгай CSP (inline script, Google Fonts зөвшөөрнө)
//
//	Permissions-Policy:
//	  - Browser features (geolocation, microphone, camera) хориглоно
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Ашиглалт:
//
//	app.Use(middleware.SecurityHeaders())
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ============================================================
		// BASIC SECURITY HEADERS
		// ============================================================

		// MIME sniffing хамгаалалт
		c.Set("X-Content-Type-Options", "nosniff")

		// Clickjacking хамгаалалт
		c.Set("X-Frame-Options", "DENY")

		// Referrer хязгаарлалт
		c.Set("Referrer-Policy", "no-referrer")

		// ============================================================
		// CSP - Use pre-computed headers for performance
		// ============================================================
		if strings.HasPrefix(c.Path(), "/docs") {
			// Swagger UI needs relaxed CSP
			c.Set("Content-Security-Policy", cspSwagger)
			return c.Next()
		}

		// Default strict CSP
		if c.GetRespHeader("Content-Security-Policy") == "" {
			c.Set("Content-Security-Policy", cspDefault)
		}

		// Permissions policy (pre-computed)
		c.Set("Permissions-Policy", permissionsPolicy)

		return c.Next()
	}
}

// ============================================================
// BODY SIZE LIMIT
// ============================================================

// BodySizeLimit нь request body хэмжээг хязгаарлах middleware буцаана.
// DDoS, memory exhaustion халдлагаас хамгаална.
//
// Parameters:
//   - maxBytes: Хамгийн их byte хэмжээ
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Response:
//   - 413 Request Entity Too Large (хэтэрсэн бол)
//
// Ашиглалт:
//
//	// 1MB хязгаар
//	app.Use(middleware.BodySizeLimit(1 * 1024 * 1024))
//
//	// 10MB хязгаар (file upload)
//	app.Post("/upload", middleware.BodySizeLimit(10*1024*1024), handler.Upload)
func BodySizeLimit(maxBytes int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Request body-ийн хэмжээ шалгах
		if len(c.BodyRaw()) > maxBytes {
			return fiber.NewError(fiber.StatusRequestEntityTooLarge, "request body too large")
		}
		return c.Next()
	}
}

// ============================================================
// PAGINATION VALIDATION
// ============================================================

// DefaultMaxPageSize нь нэг хуудсанд хамгийн их бичлэгийн тоо
const DefaultMaxPageSize = 100

// DefaultMinPageSize нь нэг хуудсанд хамгийн бага бичлэгийн тоо
const DefaultMinPageSize = 1

// PaginationLimit нь pagination параметрүүдийг хязгаарлах middleware буцаана.
// Хэт их мэдээлэл татаж авахаас сэргийлнэ.
//
// Parameters:
//   - maxSize: Нэг хуудсанд хамгийн их бичлэг (default: 100)
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Query parameters:
//   - size/pageSize: Нэг хуудсанд хэдэн бичлэг
//   - page: Хуудасны дугаар (1-ээс эхэлнэ)
//
// Response:
//   - 400 Bad Request (хязгаар хэтэрсэн бол)
//
// Ашиглалт:
//
//	app.Use(middleware.PaginationLimit(100))
func PaginationLimit(maxSize ...int) fiber.Handler {
	max := DefaultMaxPageSize
	if len(maxSize) > 0 && maxSize[0] > 0 {
		max = maxSize[0]
	}

	return func(c *fiber.Ctx) error {
		// Size параметр шалгах (size эсвэл pageSize)
		size := c.QueryInt("size", 0)
		if size == 0 {
			size = c.QueryInt("pageSize", 0)
		}

		// Size хязгаар шалгах
		if size > max {
			return fiber.NewError(
				fiber.StatusBadRequest,
				fmt.Sprintf("page size too large, maximum is %d", max),
			)
		}

		// Size сөрөг тоо байх ёсгүй
		if size < 0 {
			return fiber.NewError(
				fiber.StatusBadRequest,
				"page size must be positive",
			)
		}

		// Page параметр шалгах
		page := c.QueryInt("page", 0)
		if page < 0 {
			return fiber.NewError(
				fiber.StatusBadRequest,
				"page must be positive",
			)
		}

		return c.Next()
	}
}

// ============================================================
// TIMEOUT
// ============================================================

// Timeout нь request-д хатуу хугацааны хязгаар тавих middleware буцаана.
// Удаан ажиллаж байгаа request-үүдийг таслана.
//
// Parameters:
//   - d: Timeout хугацаа (жишээ: 5*time.Second)
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Хэрхэн ажиллах:
//  1. Context-д timeout нэмнэ
//  2. Service layer c.UserContext() авахад timeout-тэй context ирнэ
//  3. Timeout дуусвал context.DeadlineExceeded алдаа буцаана
//
// Ашиглалт:
//
//	// Route-д 5 секундын timeout
//	app.Get("/slow", middleware.Timeout(5*time.Second), handler.Slow)
//
//	// Route group-д timeout
//	api := app.Group("/api", middleware.Timeout(10*time.Second))
func Timeout(d time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Context-д timeout нэмэх
		ctx, cancel := context.WithTimeout(c.UserContext(), d)
		defer cancel() // Resource cleanup

		// Шинэ context-ийг Fiber-д буцаах
		c.SetUserContext(ctx)

		return c.Next()
	}
}
