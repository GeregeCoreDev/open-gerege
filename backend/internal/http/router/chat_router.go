// Package router provides implementation for router
//
// File: chat_router.go
// Description: Room (Video Conference) and Chat routes implementation
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

// MapChatRoutes нь room болон chat route-уудыг бүртгэнэ.
func MapChatRoutes(v1 fiber.Router, d *app.Dependencies, requireAuth fiber.Handler) {
	// Permission checker (cache-тэй)
	perm := d.PermCache

	// ------------------------------------------------------------
	// ROOM ROUTES (Video Conference)
	// ------------------------------------------------------------
	// Видео хурлын өрөө удирдах.
	v1.Group("/room", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewRoomHandler(d)

		// List rooms (user's own rooms - no admin permission required)
		router.Get("/", h.List)

		// Generate join token (user's own token - no admin permission required)
		router.Get("/token", h.GenerateToken)

		// Create room with permission check
		router.Post("/", auth.RequirePermission(perm, "admin.room.create"), h.Create)

		// Join room (user action - no admin permission required)
		router.Post("/join", h.Join)

		// Add users to room with permission check
		router.Post("/:id/users", auth.RequirePermission(perm, "admin.room.update"), h.AddUsers)

		// Delete room with permission check
		router.Delete("/:id", auth.RequirePermission(perm, "admin.room.delete"), h.Delete)

		// Remove user from room with permission check
		router.Delete("/:id/users/:user_id", auth.RequirePermission(perm, "admin.room.update"), h.RemoveUser)
	})

	// ------------------------------------------------------------
	// CHAT ROUTES
	// ------------------------------------------------------------
	// Chat item CRUD.
	v1.Group("/chat", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(r fiber.Router) {
		h := handlers.NewChatItemHandler(d)

		r.Get("/", auth.RequirePermission(perm, "admin.chat.read"), h.List)
		r.Post("/", auth.RequirePermission(perm, "admin.chat.create"), h.Create)
		r.Put("/:id", auth.RequirePermission(perm, "admin.chat.update"), h.Update)
		r.Delete("/:id", auth.RequirePermission(perm, "admin.chat.delete"), h.Delete)
		r.Post("/key", h.GetByKey) // Public endpoint for chat bot
	})
}

