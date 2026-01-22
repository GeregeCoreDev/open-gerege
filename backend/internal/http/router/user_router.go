// Package router provides implementation for router
//
// File: user_router.go
// Description: User CRUD routes implementation
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

// MapUserRoutes нь user CRUD route-уудыг бүртгэнэ.
func MapUserRoutes(v1 fiber.Router, d *app.Dependencies, requireAuth fiber.Handler) {
	// ------------------------------------------------------------
	// USER ROUTES
	// ------------------------------------------------------------
	// Хэрэглэгчийн CRUD.
	v1.Group("/user", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		handler := handlers.NewUserHandler(d)

		// Find user from Core system
		// POST /user/find-from-core → Search user in Core database
		router.Post("/find-from-core", auth.RequirePermission(d.PermCache, "admin.user.read"), handler.FindFromCore)

		// User CRUD
		// GET    /user       → List users (paginated)
		// POST   /user       → Create user
		// PUT    /user/:id   → Update user
		// DELETE /user/:id   → Delete user
		router.Get("/", auth.RequirePermission(d.PermCache, "admin.user.read"), handler.List)
		router.Post("/", auth.RequirePermission(d.PermCache, "admin.user.create"), handler.Create)
		router.Put("/:id", auth.RequirePermission(d.PermCache, "admin.user.update"), handler.Update)
		router.Delete("/:id", auth.RequirePermission(d.PermCache, "admin.user.delete"), handler.Delete)
	})
}

