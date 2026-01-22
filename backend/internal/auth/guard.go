// Package auth provides implementation for auth
//
// File: guard.go
// Description: implementation for auth
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package auth нь SSO authentication болон session management-ийг хариуцна.

Энэ файл нь authorization guard middleware-уудыг тодорхойлно.
Guard-ууд нь authenticated хэрэглэгчийн нэмэлт шаардлагуудыг шалгана.

Guard-ууд (Require middleware-ийн дараа ашиглана):
  - RequireUser: Хэрэглэгч нэвтэрсэн байх
  - RequireCitizen: Иргэн баталгаажсан байх
  - RequireOrg: Байгууллагад хамаарах
  - RequireTerminal: Терминал ID-тай байх
  - RequireApp: Тодорхой application-д хамаарах
  - RequireIsOrg: Байгууллагын горимд байх

Ашиглалт:

	// Route-д хэд хэдэн guard хэрэглэх
	app.Get("/org/data",
	    auth.Require(cfg, log, cache),  // Authentication
	    auth.RequireOrg(),               // Байгууллагад хамаарах
	    auth.RequireIsOrg(),             // Байгууллагын горимд байх
	    handler.OrgData,
	)
*/
package auth

import (
	"slices"  // Contains функц
	"strconv" // Int to string

	"git.gerege.mn/backend-packages/sso-client" // SSO client

	"github.com/gofiber/fiber/v2" // Web framework
)

// ============================================================
// REQUIRE USER
// ============================================================

// RequireUser нь authenticated хэрэглэгч байхыг шаардана.
// UserID == 0 бол 403 Forbidden буцаана.
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Ашиглалт:
//
//	app.Get("/profile", auth.Require(...), auth.RequireUser(), handler.Profile)
func RequireUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := ssoclient.GetUserID(c)
		if userId == 0 {
			return fiber.NewError(fiber.StatusForbidden, "user access required")
		}
		return c.Next()
	}
}

// ============================================================
// REQUIRE CITIZEN
// ============================================================

// RequireCitizen нь баталгаажсан иргэн байхыг шаардана.
// CitizenID == 0 бол 403 Forbidden буцаана.
//
// Иргэний баталгаажуулалт:
//   - ХУР (Хувийн мэдээллийн улсын регистр)-ээс шалгагдсан
//   - Регистрийн дугаар, нэр зэрэг мэдээлэл баталгаажсан
//
// Returns:
//   - fiber.Handler: Middleware function
func RequireCitizen() fiber.Handler {
	return func(c *fiber.Ctx) error {
		citizenId := ssoclient.GetCitizenID(c)
		if citizenId == 0 {
			return fiber.NewError(fiber.StatusForbidden, "citizen access required")
		}
		return c.Next()
	}
}

// ============================================================
// REQUIRE ORG
// ============================================================

// RequireOrg нь хэрэглэгч байгууллагад хамаарахыг шаардана.
// UserID == 0 бол 403 Forbidden буцаана.
//
// Тайлбар: Энэ guard нь OrgID биш UserID шалгана.
// Байгууллагын хэрэглэгч гэдгийг илэрхийлнэ.
//
// Returns:
//   - fiber.Handler: Middleware function
func RequireOrg() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl, ok := ssoclient.GetClaims(c)
		if !ok || cl.UserID == 0 {
			return fiber.NewError(fiber.StatusForbidden, "user access required")
		}
		return c.Next()
	}
}

// ============================================================
// REQUIRE TERMINAL
// ============================================================

// RequireTerminal нь терминал ID-тай байхыг шаардана.
// TerminalID == 0 бол 403 Forbidden буцаана.
//
// Терминал:
//   - Байгууллагын төхөөрөмж (POS, kiosk гэх мэт)
//   - Тодорхой байршилд суурилагдсан
//
// Returns:
//   - fiber.Handler: Middleware function
func RequireTerminal() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl, ok := ssoclient.GetClaims(c)
		if !ok || cl.TerminalID == 0 {
			return fiber.NewError(fiber.StatusForbidden, "terminal access required")
		}
		return c.Next()
	}
}

// ============================================================
// REQUIRE APP
// ============================================================

// RequireApp нь тодорхой application-д хамаарахыг шаардана.
// AppID жагсаалтад байхгүй бол 403 Forbidden буцаана.
//
// Parameters:
//   - appIDs: Зөвшөөрөгдсөн AppID-уудын жагсаалт (variadic)
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Ашиглалт:
//
//	// Зөвхөн AppID 1 эсвэл 2-т зөвшөөрнө
//	app.Get("/admin", auth.RequireApp(1, 2), handler.Admin)
//
//	// Хоосон жагсаалт = бүгдэд зөвшөөрнө
//	app.Get("/public", auth.RequireApp(), handler.Public)
func RequireApp(appIDs ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl, ok := ssoclient.GetClaims(c)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "no claims")
		}

		// Хоосон жагсаалт = бүгдэд зөвшөөрнө
		if len(appIDs) == 0 || slices.Contains(appIDs, cl.AppID) {
			return c.Next()
		}

		// AppID жагсаалтад байхгүй
		return fiber.NewError(fiber.StatusForbidden, "forbidden for app: "+strconv.FormatInt(int64(cl.AppID), 10))
	}
}

// ============================================================
// REQUIRE IS ORG
// ============================================================

// RequireIsOrg нь байгууллагын горимд байхыг шаардана.
// IsOrg == false бол 403 Forbidden буцаана.
//
// Байгууллагын горим:
//   - Хэрэглэгч байгууллагаар нэвтэрсэн (хувь хүнээр биш)
//   - OrgID сонгогдсон байх
//
// Returns:
//   - fiber.Handler: Middleware function
func RequireIsOrg() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl, ok := ssoclient.GetClaims(c)
		if !ok || !cl.IsOrg {
			return fiber.NewError(fiber.StatusForbidden, "org access required")
		}
		return c.Next()
	}
}
