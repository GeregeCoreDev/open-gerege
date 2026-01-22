// Package router provides implementation for router
//
// File: news_router.go
// Description: News routes implementation
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

// MapNewsRoutes нь news route-уудыг бүртгэнэ.
func MapNewsRoutes(v1 fiber.Router, d *app.Dependencies, requireAuth fiber.Handler) {
	// Permission checker (cache-тэй)
	perm := d.PermCache

	// ------------------------------------------------------------
	// NEWS ROUTES
	// ------------------------------------------------------------
	// Мэдээний CRUD (List, Get нь public).
	v1.Group("/news", middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewNewsHandler(d)

		// Public read (no permission required)
		router.Get("/", h.List)
		router.Get("/:id", h.Get)

		// Protected write with permission checks
		router.Post("/", requireAuth, auth.RequirePermission(perm, "admin.news.create"), h.Create)
		router.Put("/:id", requireAuth, auth.RequirePermission(perm, "admin.news.update"), h.Update)
		router.Delete("/:id", requireAuth, auth.RequirePermission(perm, "admin.news.delete"), h.Delete)
	})
}

