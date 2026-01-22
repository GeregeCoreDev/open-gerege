// Package router provides implementation for router
//
// File: file_router.go
// Description: File upload/download routes implementation
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

// MapFileRoutes нь file upload/download route-уудыг бүртгэнэ.
func MapFileRoutes(v1 fiber.Router, d *app.Dependencies, requireAuth fiber.Handler) {
	// Permission checker (cache-тэй)
	perm := d.PermCache

	// ------------------------------------------------------------
	// FILE ROUTES
	// ------------------------------------------------------------
	// Файл upload/download.
	v1.Group("/file", middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewFileHandler(d)

		// Protected file management with permission checks
		router.Get("/list", requireAuth, auth.RequirePermission(perm, "admin.file.read"), h.GetPublicFileList)
		router.Post("/upload", requireAuth, auth.RequirePermission(perm, "admin.file.create"), h.Upload)
		router.Delete("/", requireAuth, auth.RequirePermission(perm, "admin.file.delete"), h.DeletePublicFile)

		// Public file download (auth хэрэггүй)
		// GET /file/:uuid → Download file by UUID
		router.Get("/:uuid", h.GetFile)
	})
}

