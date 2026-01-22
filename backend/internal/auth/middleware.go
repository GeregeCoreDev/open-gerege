// Package auth provides implementation for auth
//
// File: middleware.go
// Description: implementation for auth
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package auth нь SSO authentication болон session management-ийг хариуцна.

Энэ файл нь authentication middleware болон session ID extract хийх
функцуудыг агуулна.

Authentication flow:
 1. Cookie эсвэл Authorization header-ээс SID авах
 2. Cache-ээс Claims хайх (cache hit бол SSO руу явахгүй)
 3. Cache-д байхгүй бол SSO /whoami руу request илгээх
 4. Claims-ийг cache-д хадгалах
 5. Claims-ийг Fiber Locals болон context-д хадгалах
 6. Handler ажиллах

Session extraction:
  - Cookie: sid=xxx (тохиргооноос cookie нэр авна)
  - Authorization header:
  - Bearer xxx
  - sid=xxx
  - xxx (raw token)
*/
package auth

import (
	"context" // Timeout context
	"strings" // String manipulation
	"time"    // Timeout duration

	"git.gerege.mn/backend-packages/config"   // Configuration
	"git.gerege.mn/backend-packages/ctx"      // Context helpers
	"git.gerege.mn/backend-packages/sso-client" // SSO client

	"github.com/gofiber/fiber/v2" // Web framework
	"go.uber.org/zap"             // Structured logging
)

// ============================================================
// REQUIRE MIDDLEWARE
// ============================================================

// Require нь authentication middleware буцаана.
// Protected route-уудад ашиглана.
//
// Flow:
//  1. Cookie/Authorization-оос SID авах
//  2. Cache-ээс Claims хайх
//  3. Cache miss бол SSO /whoami руу request
//  4. Claims-ийг Locals/context-д хадгалах
//  5. Дараагийн handler руу шилжих
//
// Parameters:
//   - cfg: Application configuration
//   - log: Zap logger
//   - cache: Session cache
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Жишээ:
//
//	requireAuth := auth.Require(cfg, log, cache)
//	app.Get("/protected", requireAuth, handler.Protected)
func Require(cfg *config.Config, log *zap.Logger, cache *ssoclient.Cache) fiber.Handler {
	// Урьдчилсан шалгалт: Auth тохиргоо бүрэн байгаа эсэх
	if cfg.Auth.ClientID == "" || cfg.Auth.ClientSecret == "" || cfg.URLS.SSO == "" {
		// Тохиргоо дутуу бол бүх request-д 401 буцаах
		return func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusUnauthorized, "auth is not configured")
		}
	}

	// SSO HTTP client үүсгэх
	sso := ssoclient.NewSSOClient(cfg, log, cache)

	return func(c *fiber.Ctx) error {
		// ============================================================
		// STEP 1: Session ID олох
		// ============================================================
		// Cookie болон Authorization header-ээс SID хайна
		sid := ExtractSID(c, cfg)
		if sid == "" {
			// SID байхгүй бол 401 Unauthorized
			return fiber.NewError(fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message)
		}

		// ============================================================
		// STEP 2: Request ID авах
		// ============================================================
		// Logging, debugging-д ашиглагдана
		reqID := ctx.RequestID(c)

		// ============================================================
		// STEP 3: Claims авах (cache → SSO)
		// ============================================================
		// 3 секундын timeout тавих
		ctxTimeout, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
		defer cancel()

		// SSO client ашиглан Claims авах
		// Cache-д байвал SSO руу явахгүй
		claims, err := sso.GetClaims(ctxTimeout, sid, reqID)
		if err != nil {
			// SSO алдаа эсвэл session invalid
			return fiber.NewError(fiber.StatusUnauthorized, fiber.ErrUnauthorized.Message)
		}

		// ============================================================
		// STEP 4: Context/Locals-д хадгалах
		// ============================================================
		// Handler-ууд auth.GetUserID(c) гэх мэтээр авна
		attachToCtx(c, sid, &claims)

		// ============================================================
		// STEP 5: Дараагийн handler руу шилжих
		// ============================================================
		return c.Next()
	}
}

// ============================================================
// ATTACH TO CONTEXT
// ============================================================

// attachToCtx нь claims-ийг Fiber Locals болон stdlib context-д хадгална.
//
// Хадгалах газрууд:
//  1. Fiber Locals: Handler-ууд auth.GetUserID(c) гэх мэтээр авна
//  2. stdlib context: Service layer-т c.UserContext()-оор дамжуулагдана
//
// Parameters:
//   - c: Fiber context
//   - sid: Session ID
//   - claims: Session claims
func attachToCtx(c *fiber.Ctx, sid string, claims *ssoclient.Claims) {
	// ============================================================
	// FIBER LOCALS
	// ============================================================
	// Handler-ууд шууд авах боломжтой
	// Жишээ: ssoclient.GetUserID(c)
	c.Locals(ssoclient.LocalsSID, sid)
	c.Locals(ssoclient.LocalsClaims, claims)

	// ============================================================
	// STDLIB CONTEXT
	// ============================================================
	// Service layer-т дамжуулагдана
	// Жишээ: userID := ctx.GetUserID(c.UserContext())
	uc := c.UserContext()

	// Session ID
	if sid != "" {
		uc = ctx.WithValue(uc, ctx.KeySID, sid)
	}

	// Claims-ийн талбаруудыг context-д нэмэх
	if claims != nil {
		if claims.UserID != 0 {
			uc = ctx.WithValue(uc, ctx.KeyUserID, claims.UserID)
		}
		if claims.Username != "" {
			uc = ctx.WithValue(uc, ctx.KeyUsername, claims.Username)
		}
		if claims.CitizenID != 0 {
			uc = ctx.WithValue(uc, ctx.KeyCitizenID, claims.CitizenID)
		}
		if claims.OrgID != 0 {
			uc = ctx.WithValue(uc, ctx.KeyOrgID, claims.OrgID)
		}
		if claims.TerminalID != 0 {
			uc = ctx.WithValue(uc, ctx.KeyTerminalID, claims.TerminalID)
		}
		if claims.AppID != 0 {
			uc = ctx.WithValue(uc, ctx.KeyAppID, claims.AppID)
		}
		uc = ctx.WithValue(uc, ctx.KeyIsOrg, claims.IsOrg)
	}

	// Request ID
	if rid := ctx.RequestID(c); rid != "" {
		uc = ctx.WithValue(uc, ctx.KeyRequestID, rid)
	}

	// Locals дахь Request ID (backup)
	if rid, ok := c.Locals(ctx.KeyRequestID).(string); ok && rid != "" {
		uc = context.WithValue(uc, ctx.KeyRequestID, rid)
	}

	// Шинэ context-ийг Fiber-д буцаах
	c.SetUserContext(uc)
}

// ============================================================
// SESSION ID EXTRACTION
// ============================================================

// ExtractSID нь session ID-г cookie эсвэл Authorization header-ээс авна.
//
// Хайх дараалал:
//  1. Cookie (тохиргоогоор нэрийг авна, default: "sid")
//  2. Authorization header (Bearer xxx, sid=xxx, эсвэл raw)
//
// Parameters:
//   - c: Fiber context
//   - cfg: Configuration (cookie нэр авах)
//
// Returns:
//   - string: Session ID (хоосон бол "")
func ExtractSID(c *fiber.Ctx, cfg *config.Config) string {
	// 1) Cookie-ээс хайх (давуу эрхтэй)
	if sid := extractFromCookie(c, cfg.Cookie.Name); sid != "" {
		return sid
	}
	// 2) Authorization header-ээс хайх
	return extractFromAuthHeader(c)
}

// extractFromCookie нь cookie-оос SID авна.
//
// Parameters:
//   - c: Fiber context
//   - cookieName: Cookie-ийн нэр (жишээ: "sid")
//
// Returns:
//   - string: Cookie value (trimmed)
func extractFromCookie(c *fiber.Ctx, cookieName string) string {
	return strings.TrimSpace(c.Cookies(cookieName))
}

// Authorization header prefix-ууд
const (
	bearerPrefix = "bearer " // "Bearer xxx" format
	sidPrefix    = "sid="    // "sid=xxx" format
)

// extractFromAuthHeader нь Authorization header-ээс SID авна.
//
// Дэмжлэг:
//   - "Bearer xxx" → "xxx"
//   - "bearer xxx" → "xxx" (case insensitive)
//   - "sid=xxx" → "xxx"
//   - "xxx" → "xxx" (raw token)
//
// Parameters:
//   - c: Fiber context
//
// Returns:
//   - string: Extracted session ID (хоосон бол "")
func extractFromAuthHeader(c *fiber.Ctx) string {
	// Authorization header авах
	authHeader := strings.TrimSpace(c.Get(fiber.HeaderAuthorization))
	if authHeader == "" {
		return ""
	}

	// Lowercase хувилбар (case insensitive matching-д)
	lowerAuth := strings.ToLower(authHeader)

	// "Bearer xxx" format
	if strings.HasPrefix(lowerAuth, bearerPrefix) {
		return strings.TrimSpace(authHeader[len(bearerPrefix):])
	}

	// "sid=xxx" format
	if strings.HasPrefix(lowerAuth, sidPrefix) {
		return strings.TrimSpace(authHeader[len(sidPrefix):])
	}

	// Raw token (хэрэв prefix байхгүй бол)
	return authHeader
}
