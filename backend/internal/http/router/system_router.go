// Package router provides implementation for router
//
// File: system_router.go
// Description: System, Module, Permission, Action, Role, and Client routes implementation
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

// MapSystemRoutes нь system, module, permission, action, role, client route-уудыг бүртгэнэ.
func MapSystemRoutes(v1 fiber.Router, d *app.Dependencies, requireAuth fiber.Handler) {
	// Permission checker (cache-тэй)
	perm := d.PermCache

	// ------------------------------------------------------------
	// SYSTEM ROUTES
	// ------------------------------------------------------------
	// Системийн CRUD (app groups).
	v1.Group("/system", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewSystemHandler(d)

		// CRUD operations with permission checks
		router.Get("/", auth.RequirePermission(perm, "admin.system.read"), h.List)
		router.Get("/:id", auth.RequirePermission(perm, "admin.system.read"), h.Get)
		router.Post("/", auth.RequirePermission(perm, "admin.system.create"), h.Create)
		router.Put("/:id", auth.RequirePermission(perm, "admin.system.update"), h.Update)
		router.Delete("/:id", auth.RequirePermission(perm, "admin.system.delete"), h.Delete)
	})

	// ------------------------------------------------------------
	// MODULE ROUTES
	// ------------------------------------------------------------
	// Модулийн (menu) CRUD болон access control.
	v1.Group("/module", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(r fiber.Router) {
		h := handlers.NewModuleHandler(d)

		// CRUD operations with permission checks
		r.Get("/", auth.RequirePermission(perm, "admin.module.read"), h.List)
		r.Post("/", auth.RequirePermission(perm, "admin.module.create"), h.Create)
		r.Put("/:id", auth.RequirePermission(perm, "admin.module.update"), h.Update)
		r.Delete("/:id", auth.RequirePermission(perm, "admin.module.delete"), h.Delete)
	})

	// ------------------------------------------------------------
	// PERMISSION ROUTES
	// ------------------------------------------------------------
	// Зөвшөөрлийн CRUD.
	v1.Group("/permission", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewPermissionHandler(d)

		// CRUD operations with permission checks
		router.Get("/", auth.RequirePermission(perm, "admin.permission.read"), h.List)
		router.Post("/", auth.RequirePermission(perm, "admin.permission.create"), h.Create)
		router.Put("/:id", auth.RequirePermission(perm, "admin.permission.update"), h.Update)
		router.Delete("/:id", auth.RequirePermission(perm, "admin.permission.delete"), h.Delete)
	})

	// ------------------------------------------------------------
	// ACTION ROUTES
	// ------------------------------------------------------------
	// Action-ийн CRUD (Permission-тэй ижил логик).
	v1.Group("/action", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewActionHandler(d)

		// CRUD operations with permission checks
		router.Get("/", auth.RequirePermission(perm, "admin.action.read"), h.List)
		router.Post("/", auth.RequirePermission(perm, "admin.action.create"), h.Create)
		router.Put("/:id", auth.RequirePermission(perm, "admin.action.update"), h.Update)
		router.Delete("/:id", auth.RequirePermission(perm, "admin.action.delete"), h.Delete)
	})

	// ------------------------------------------------------------
	// ROLE ROUTES
	// ------------------------------------------------------------
	// Эрхийн CRUD болон permission assignment.
	v1.Group("/role", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		role := handlers.NewRoleHandler(d)

		// CRUD operations with permission checks
		router.Get("/", auth.RequirePermission(perm, "admin.role.read"), role.List)
		router.Post("/", auth.RequirePermission(perm, "admin.role.create"), role.Create)
		router.Put("/:id", auth.RequirePermission(perm, "admin.role.update"), role.Update)
		router.Delete("/:id", auth.RequirePermission(perm, "admin.role.delete"), role.Delete)

		// Permission management with permission checks
		// GET  /role/permissions?role_id=1 → Role's permissions
		// POST /role/permissions {role_id, permission_ids} → Set permissions
		router.Get("/permissions", auth.RequirePermission(perm, "admin.role.read"), role.GetRolePermissions)
		router.Post("/permissions", auth.RequirePermission(perm, "admin.role.update"), role.SetRolePermissions)
	})

	// ------------------------------------------------------------
	// CLIENT ROUTES
	// ------------------------------------------------------------
	// OAuth client management.
	v1.Group("/client", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		handler := handlers.NewClientHandler(d)

		// CRUD operations with permission checks
		router.Get("", auth.RequirePermission(perm, "admin.client.read"), handler.List)
		router.Get("/scope", auth.RequirePermission(perm, "admin.client.read"), handler.ScopeList)
		router.Post("/scope", auth.RequirePermission(perm, "admin.client.create"), handler.ScopeCreate)
		router.Delete("/scope", auth.RequirePermission(perm, "admin.client.delete"), handler.ScopeDelete)
	})

	// ------------------------------------------------------------
	// MENU ROUTES
	// ------------------------------------------------------------
	// Menu CRUD.
	v1.Group("/menu", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewMenuHandler(d)

		// CRUD operations with permission checks
		// /my endpoint нь хэрэглэгчийн өөрийн menu-г буцаадаг тул permission шаардахгүй
		router.Get("/", auth.RequirePermission(perm, "admin.menu.read"), h.List)
		router.Get("/my", h.ListByRole) // Get menus by current user's roles (no permission required)
		router.Get("/:id", auth.RequirePermission(perm, "admin.menu.read"), h.Get)
		router.Post("/", auth.RequirePermission(perm, "admin.menu.create"), h.Create)
		router.Put("/:id", auth.RequirePermission(perm, "admin.menu.update"), h.Update)
		router.Delete("/:id", auth.RequirePermission(perm, "admin.menu.delete"), h.Delete)
	})

	// ------------------------------------------------------------
	// ROLE-MATRIX ROUTES
	// ------------------------------------------------------------
	// Хэрэглэгч-эрхийн холбоосыг удирдах.
	v1.Group("/role-matrix", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(g fiber.Router) {
		h := handlers.NewUserRoleHandler(d)

		// List users by role with permission checks
		// GET /role-matrix/users?role_id=1 → Users with specific role
		g.Get("/users", auth.RequirePermission(perm, "admin.user-role.read"), h.UsersByRole)

		// List roles by user with permission checks
		// GET /role-matrix/roles?user_id=1 → Roles of specific user
		g.Get("/roles", auth.RequirePermission(perm, "admin.user-role.read"), h.RolesByUser)

		// Assign role to user with permission checks
		// POST /role-matrix {user_id, role_id}
		g.Post("/", auth.RequirePermission(perm, "admin.user-role.create"), h.Create)

		// Remove role from user with permission checks
		// DELETE /role-matrix {user_id, role_id}
		g.Delete("/", auth.RequirePermission(perm, "admin.user-role.delete"), h.Delete)
	})
}
