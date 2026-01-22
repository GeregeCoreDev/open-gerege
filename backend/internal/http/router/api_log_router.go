// Package router provides implementation for router
//
// File: api_log_router.go
// Description: API Log routes implementation
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-01-09
// Last Updated: 2025-01-09
package router

import (
	"time"

	"templatev25/internal/app"
	"templatev25/internal/auth"
	"templatev25/internal/http/handlers"
	"templatev25/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// MapAPILogRoutes нь API log route-уудыг бүртгэнэ.
func MapAPILogRoutes(v1 fiber.Router, d *app.Dependencies, requireAuth fiber.Handler) {
	// Permission checker (cache-тэй)
	perm := d.PermCache

	// ------------------------------------------------------------
	// API LOG ROUTES
	// ------------------------------------------------------------
	// API log-ийн list (paginated).
	v1.Group("/api-logs", requireAuth, middleware.Timeout(10*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewAPILogHandler(d)

		// List API logs (paginated) with permission check
		router.Get("/", auth.RequirePermission(perm, "admin.api-log.read"), h.List)
	})
}
