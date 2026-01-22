// Package auth provides implementation for auth
//
// File: permission.go
// Description: Permission-based authorization middleware
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package auth нь SSO authentication болон authorization-ийг хариуцна.

Энэ файл нь permission-based authorization middleware-уудыг тодорхойлно.
Permission middleware нь RBAC системийн гол хэсэг бөгөөд
хэрэглэгчийн эрхийг runtime-д шалгана.

Permission шалгах flow:
 1. Cache-ээс permission хайх (хурдан)
 2. Cache-д байхгүй бол DB-ээс авах
 3. Cache-д хадгалах
 4. Permission байвал c.Next(), байхгүй бол 403

Middleware-ууд:
  - RequirePermission: Нэг permission шалгах
  - RequireAnyPermission: Аль нэг permission байвал болно
  - RequireAllPermissions: Бүх permission байх шаардлагатай

Ашиглалт:

	// Route-д permission шалгах
	app.Post("/role",
	    auth.Require(cfg, log, cache),
	    auth.RequirePermission(permChecker, "admin.role.create"),
	    handler.RoleCreate,
	)

	// Олон permission-ийн аль нэг нь байвал болно
	app.Get("/data",
	    auth.Require(cfg, log, cache),
	    auth.RequireAnyPermission(permChecker, "admin.data.read", "user.data.read"),
	    handler.GetData,
	)
*/
package auth

import (
	"context"
	"slices"
	"time"

	"git.gerege.mn/backend-packages/sso-client"

	"github.com/gofiber/fiber/v2"
)

// ============================================================
// PERMISSION CHECKER INTERFACE
// ============================================================

// PermissionChecker нь permission шалгах интерфейс.
// Service эсвэл cached service энэ интерфейсийг implement хийнэ.
type PermissionChecker interface {
	// HasPermission нь хэрэглэгч тодорхой permission-тэй эсэхийг шалгана
	HasPermission(ctx context.Context, userID int, permissionCode string) (bool, error)
	// GetUserPermissions нь хэрэглэгчийн бүх permission-уудыг буцаана
	GetUserPermissions(ctx context.Context, userID int) ([]string, error)
}

// ============================================================
// REQUIRE PERMISSION
// ============================================================

// RequirePermission нь нэг тодорхой permission шаардана.
// Permission байхгүй бол 403 Forbidden буцаана.
//
// Parameters:
//   - checker: Permission шалгах service
//   - permissionCode: Шаардлагатай permission код (жишээ: "admin.role.create")
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Ашиглалт:
//
//	app.Post("/role", auth.RequirePermission(checker, "admin.role.create"), handler.Create)
func RequirePermission(checker PermissionChecker, permissionCode string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ============================================================
		// STEP 1: User ID авах
		// ============================================================
		userID := ssoclient.GetUserID(c)
		if userID == 0 {
			return fiber.NewError(fiber.StatusForbidden, "user not authenticated")
		}

		// ============================================================
		// STEP 2: Permission шалгах
		// ============================================================
		ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
		defer cancel()

		hasPermission, err := checker.HasPermission(ctx, userID, permissionCode)
		if err != nil {
			// DB алдаа - internal error биш 403 буцаах (security)
			return fiber.NewError(fiber.StatusForbidden, "permission check failed")
		}

		if !hasPermission {
			return fiber.NewError(fiber.StatusForbidden, "insufficient permissions: "+permissionCode)
		}

		// ============================================================
		// STEP 3: Дараагийн handler руу шилжих
		// ============================================================
		return c.Next()
	}
}

// ============================================================
// REQUIRE ANY PERMISSION
// ============================================================

// RequireAnyPermission нь өгөгдсөн permission-уудын аль нэг нь байхыг шаардана.
// Бүгд байхгүй бол 403 Forbidden буцаана.
//
// Parameters:
//   - checker: Permission шалгах service
//   - permissionCodes: Зөвшөөрөгдсөн permission кодуудын жагсаалт
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Ашиглалт:
//
//	// Admin эсвэл user аль нэг нь data унших эрхтэй байвал болно
//	app.Get("/data",
//	    auth.RequireAnyPermission(checker, "admin.data.read", "user.data.read"),
//	    handler.GetData,
//	)
func RequireAnyPermission(checker PermissionChecker, permissionCodes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// User ID авах
		userID := ssoclient.GetUserID(c)
		if userID == 0 {
			return fiber.NewError(fiber.StatusForbidden, "user not authenticated")
		}

		// Хоосон жагсаалт = бүгдэд зөвшөөрнө
		if len(permissionCodes) == 0 {
			return c.Next()
		}

		// Хэрэглэгчийн бүх permission авах
		ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
		defer cancel()

		userPerms, err := checker.GetUserPermissions(ctx, userID)
		if err != nil {
			return fiber.NewError(fiber.StatusForbidden, "permission check failed")
		}

		// Аль нэг нь байвал зөвшөөрнө
		for _, code := range permissionCodes {
			if slices.Contains(userPerms, code) {
				return c.Next()
			}
		}

		return fiber.NewError(fiber.StatusForbidden, "insufficient permissions")
	}
}

// ============================================================
// REQUIRE ALL PERMISSIONS
// ============================================================

// RequireAllPermissions нь өгөгдсөн бүх permission-ууд байхыг шаардана.
// Аль нэг нь байхгүй бол 403 Forbidden буцаана.
//
// Parameters:
//   - checker: Permission шалгах service
//   - permissionCodes: Шаардлагатай permission кодуудын жагсаалт
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Ашиглалт:
//
//	// Role удирдах болон хэрэглэгч удирдах хоёулаа шаардлагатай
//	app.Post("/admin/assign-role",
//	    auth.RequireAllPermissions(checker, "admin.role.update", "admin.user.update"),
//	    handler.AssignRole,
//	)
func RequireAllPermissions(checker PermissionChecker, permissionCodes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// User ID авах
		userID := ssoclient.GetUserID(c)
		if userID == 0 {
			return fiber.NewError(fiber.StatusForbidden, "user not authenticated")
		}

		// Хоосон жагсаалт = бүгдэд зөвшөөрнө
		if len(permissionCodes) == 0 {
			return c.Next()
		}

		// Хэрэглэгчийн бүх permission авах
		ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
		defer cancel()

		userPerms, err := checker.GetUserPermissions(ctx, userID)
		if err != nil {
			return fiber.NewError(fiber.StatusForbidden, "permission check failed")
		}

		// Бүгд байх ёстой
		for _, code := range permissionCodes {
			if !slices.Contains(userPerms, code) {
				return fiber.NewError(fiber.StatusForbidden, "insufficient permissions: "+code)
			}
		}

		return c.Next()
	}
}
