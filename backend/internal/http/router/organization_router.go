// Package router provides implementation for router
//
// File: organization_router.go
// Description: Organization, Organization User, and Organization Type routes implementation
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

// MapOrganizationRoutes нь organization, orguser, orgtype route-уудыг бүртгэнэ.
func MapOrganizationRoutes(v1 fiber.Router, d *app.Dependencies, requireAuth fiber.Handler) {
	// Permission checker (cache-тэй)
	perm := d.PermCache

	// ------------------------------------------------------------
	// ORGANIZATION ROUTES
	// ------------------------------------------------------------
	// Байгууллагын CRUD.
	v1.Group("/organization", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewOrganizationHandler(d)

		// Find organization from Core system
		router.Get("/find", auth.RequirePermission(perm, "admin.organization.read"), h.FindFromCore)

		// CRUD operations with permission checks
		router.Get("/", auth.RequirePermission(perm, "admin.organization.read"), h.List)
		router.Post("/", auth.RequirePermission(perm, "admin.organization.create"), h.Create)
		router.Put("/:id", auth.RequirePermission(perm, "admin.organization.update"), h.Update)
		router.Delete("/:id", auth.RequirePermission(perm, "admin.organization.delete"), h.Delete)

		// Get organization tree (hierarchical structure)
		router.Get("/tree", auth.RequirePermission(perm, "admin.organization.read"), h.Tree)
	})

	// ------------------------------------------------------------
	// ORGANIZATION USER ROUTES
	// ------------------------------------------------------------
	// Байгууллага-хэрэглэгчийн холбоос.
	v1.Group("/orguser", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewOrgUserHandler(d)

		// List all org-user relations
		router.Get("/", auth.RequirePermission(perm, "admin.orguser.read"), h.List)

		// Get users of organization
		router.Get("/users", auth.RequirePermission(perm, "admin.orguser.read"), h.Users)

		// Get organizations of user
		router.Get("/organizations", auth.RequirePermission(perm, "admin.orguser.read"), h.Orgs)

		// Add user to organization
		router.Post("/", auth.RequirePermission(perm, "admin.orguser.create"), h.Add)

		// Remove user from organization
		router.Delete("/", auth.RequirePermission(perm, "admin.orguser.delete"), h.Remove)
	})

	// ------------------------------------------------------------
	// ORGANIZATION TYPE ROUTES
	// ------------------------------------------------------------
	// Байгууллагын төрлийн CRUD.
	v1.Group("/orgtype", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewOrganizationTypeHandler(d)

		// CRUD operations with permission checks
		router.Get("/", auth.RequirePermission(perm, "admin.orgtype.read"), h.List)
		router.Post("/", auth.RequirePermission(perm, "admin.orgtype.create"), h.Create)
		router.Put("/:id", auth.RequirePermission(perm, "admin.orgtype.update"), h.Update)
		router.Delete("/:id", auth.RequirePermission(perm, "admin.orgtype.delete"), h.Delete)

		// System assignment with permission checks
		// GET  /orgtype/system?type_id=1 → Systems for org type
		// POST /orgtype/system {type_id, system_ids} → Add systems
		router.Get("/system", auth.RequirePermission(perm, "admin.orgtype.read"), h.Systems)
		router.Post("/system", auth.RequirePermission(perm, "admin.orgtype.update"), h.AddSystems)

		// Role assignment with permission checks
		// GET  /orgtype/role?type_id=1 → Roles for org type
		// POST /orgtype/role {type_id, role_ids} → Add roles
		router.Get("/role", auth.RequirePermission(perm, "admin.orgtype.read"), h.Roles)
		router.Post("/role", auth.RequirePermission(perm, "admin.orgtype.update"), h.AddRoles)
	})

	// ------------------------------------------------------------
	// TERMINAL ROUTES
	// ------------------------------------------------------------
	// Терминалын CRUD.
	v1.Group("/terminal", requireAuth, middleware.Timeout(5*time.Second)).Route("", func(router fiber.Router) {
		h := handlers.NewTerminalHandler(d)

		router.Get("/", auth.RequirePermission(perm, "admin.terminal.read"), h.List)
		router.Post("/", auth.RequirePermission(perm, "admin.terminal.create"), h.Create)
		router.Put("/:id", auth.RequirePermission(perm, "admin.terminal.update"), h.Update)
		router.Delete("/:id", auth.RequirePermission(perm, "admin.terminal.delete"), h.Delete)
	})
}

