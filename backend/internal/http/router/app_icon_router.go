// Package router provides implementation for router
//
// File: app_icon_router.go
// Description: App Service Icon and App Service Group routes implementation
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

// MapAppIconRoutes нь app-service-icon болон app-service-group route-уудыг бүртгэнэ.
func MapAppIconRoutes(v1 fiber.Router, d *app.Dependencies, requireAuth fiber.Handler) {
	// Permission checker (cache-тэй)
	perm := d.PermCache

	// ------------------------------------------------------------
	// APP SERVICE ICON ROUTES
	// ------------------------------------------------------------
	// App service icon CRUD.
	v1.Group("/app-service-icon", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewAppServiceIconHandler(d)

		router.Get("/", auth.RequirePermission(perm, "admin.app-icon.read"), h.List)
		router.Post("/", auth.RequirePermission(perm, "admin.app-icon.create"), h.Create)
		router.Put("/:id", auth.RequirePermission(perm, "admin.app-icon.update"), h.Update)
		router.Delete("/:id", auth.RequirePermission(perm, "admin.app-icon.delete"), h.Delete)
	})

	// ------------------------------------------------------------
	// APP SERVICE GROUP ROUTES
	// ------------------------------------------------------------
	// App service icon group CRUD.
	v1.Group("/app-service-group", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewAppServiceGroupHandler(d)

		router.Get("/", auth.RequirePermission(perm, "admin.app-group.read"), h.List)
		router.Get("/with-icons", auth.RequirePermission(perm, "admin.app-group.read"), h.ListGroupsWithIcons)
		router.Post("/", auth.RequirePermission(perm, "admin.app-group.create"), h.Create)
		router.Put("/:id", auth.RequirePermission(perm, "admin.app-group.update"), h.Update)
		router.Delete("/:id", auth.RequirePermission(perm, "admin.app-group.delete"), h.Delete)
	})
}

