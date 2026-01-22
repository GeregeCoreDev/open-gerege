// Package router provides implementation for router
//
// File: notification_router.go
// Description: Notification routes implementation
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package router

import (
	"time"

	"templatev25/internal/app"
	"templatev25/internal/auth"
	"templatev25/internal/http/handlers"
	"templatev25/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// MapNotificationRoutes нь notification route-уудыг бүртгэнэ.
func MapNotificationRoutes(v1 fiber.Router, d *app.Dependencies, requireAuth fiber.Handler) {
	// Permission checker (cache-тэй)
	perm := d.PermCache

	// ------------------------------------------------------------
	// NOTIFICATION ROUTES
	// ------------------------------------------------------------
	// Мэдэгдэл илгээх, унших.
	v1.Group("/notification", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewNotificationHandler(d)

		// List notifications (user's own notifications - no admin permission required)
		router.Get("/", h.List)

		// Get notification groups (user's own groups - no admin permission required)
		router.Get("/groups", h.Groups)

		// Send notification (requires admin permission)
		router.Post("/", auth.RequirePermission(perm, "admin.notification.create"), h.Send)

		// Mark as read (user's own notifications - no admin permission required)
		router.Post("/read", h.Read)
		router.Post("/read-all", h.ReadAll)
	})
}

