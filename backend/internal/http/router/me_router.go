// Package router provides implementation for router
//
// File: me_router.go
// Description: Current user (me) routes implementation
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package router

import (
	"time"

	"templatev25/internal/app"
	"templatev25/internal/http/handlers"
	"templatev25/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// MapMeRoutes нь current user (me)-тэй холбоотой route-уудыг бүртгэнэ.
//
// Routes:
//   Profile:
//   - GET  /me           → Current user info
//   - GET  /me/profile   → Full profile
//   - GET  /me/profile/sso → SSO profile
//   - GET  /me/organizations → User organizations
//
//   Security (Local Auth) - Path: /auth/local/me/*
//   - GET    /auth/local/me/sessions         → List active sessions
//   - DELETE /auth/local/me/sessions/:id     → Revoke specific session
//   - POST   /auth/local/me/password         → Change password
//   - GET    /auth/local/me/mfa              → Get MFA status
//   - POST   /auth/local/me/mfa/totp/setup   → Setup TOTP
//   - POST   /auth/local/me/mfa/totp/confirm → Confirm TOTP setup
//   - DELETE /auth/local/me/mfa/totp         → Disable TOTP
//   - POST   /auth/local/me/mfa/backup-codes → Generate backup codes
//   - GET    /auth/local/me/login-history    → Login history
//   - GET    /auth/local/me/security-audit   → Security audit trail
//
//   Payment:
//   - /me/accounts/*          → Account management
//   - /me/card/*              → Card management
//   - /me/tpay/transaction/*  → Payment transactions
func MapMeRoutes(v1 fiber.Router, d *app.Dependencies, requireAuth fiber.Handler) {
	// ------------------------------------------------------------
	// ME ROUTES (Current User)
	// ------------------------------------------------------------
	// Current user-тэй холбоотой endpoint-ууд.
	v1.Group("/me", requireAuth).Route("", func(router fiber.Router) {
		userHandler := handlers.NewUserHandler(d)
		tpayHandler := handlers.NewTpayHandler(d)

		// Current user info (from session)
		router.Get("/", middleware.Timeout(5*time.Second), userHandler.Me)

		// Profile & organizations
		router.Get("/profile", middleware.Timeout(5*time.Second), userHandler.Profile)
		router.Get("/profile/sso", middleware.Timeout(5*time.Second), userHandler.ProfileSSO)
		router.Get("/organizations", middleware.Timeout(5*time.Second), userHandler.Organizations)

		// Account management
		accr := router.Group("/accounts")
		accr.Get("/", middleware.Timeout(5*time.Second), tpayHandler.Account.GetMyAccounts)
		accr.Put("/default", middleware.Timeout(5*time.Second), tpayHandler.Account.SetDefaultAccount)
		accr.Get("/statement", middleware.Timeout(5*time.Second), tpayHandler.Account.GetStatement)
		accr.Post("/:account_id/qr", middleware.Timeout(5*time.Second), tpayHandler.Account.GenerateQR)

		// Card management
		cardr := router.Group("/card")
		cardr.Get("/list", middleware.Timeout(5*time.Second), tpayHandler.Card.CardList)
		cardr.Post("/create", middleware.Timeout(5*time.Second), tpayHandler.Card.AddCard)
		cardr.Post("/confirm", middleware.Timeout(5*time.Second), tpayHandler.Card.Confirm)
		cardr.Get("/otp", middleware.Timeout(5*time.Second), tpayHandler.Card.SendOtp)
		cardr.Post("/verify", middleware.Timeout(5*time.Second), tpayHandler.Card.VerifyCard)

		// TPAY Payment transactions
		payr := router.Group("/tpay/transaction")
		payr.Post("/qr-pay", middleware.Timeout(5*time.Second), tpayHandler.Payment.QrPay)
		payr.Post("/p2p", middleware.Timeout(5*time.Second), tpayHandler.Payment.P2PTransfer)

	})

	// ------------------------------------------------------------
	// ME SECURITY ROUTES (Local Auth - Session Protected)
	// ------------------------------------------------------------
	// Local authentication-тай холбоотой user management endpoint-ууд
	// Session auth middleware ашиглана
	// Use adapter to bridge service.SessionStore to middleware.SessionStore interface
	sessionStoreAdapter := NewSessionStoreAdapter(d.Service.SessionStore)
	sessionAuth := middleware.SessionAuth(sessionStoreAdapter)

	v1.Group("/auth/local/me", sessionAuth).Route("", func(router fiber.Router) {
		userMgmtHandler := handlers.NewUserManagementHandler(d.Service.Auth)
		strictLimiter := middleware.StrictRateLimiter()

		// Session management
		// GET  /me/sessions     → List all active sessions
		// DELETE /me/sessions/:id → Revoke specific session
		router.Get("/sessions", middleware.Timeout(5*time.Second), userMgmtHandler.ListSessions)
		router.Delete("/sessions/:id", middleware.Timeout(5*time.Second), userMgmtHandler.RevokeSession)

		// Password management (rate limited)
		// POST /me/password → Change password
		router.Post("/password", strictLimiter, middleware.Timeout(5*time.Second), userMgmtHandler.ChangePassword)

		// MFA management
		mfar := router.Group("/mfa")
		// GET /me/mfa → Get MFA status
		mfar.Get("/", middleware.Timeout(5*time.Second), userMgmtHandler.GetMFAStatus)

		// TOTP setup (rate limited)
		// POST /me/mfa/totp/setup   → Initiate TOTP setup
		// POST /me/mfa/totp/confirm → Confirm TOTP with code
		// DELETE /me/mfa/totp       → Disable TOTP
		mfar.Post("/totp/setup", strictLimiter, middleware.Timeout(5*time.Second), userMgmtHandler.SetupTOTP)
		mfar.Post("/totp/confirm", strictLimiter, middleware.Timeout(5*time.Second), userMgmtHandler.ConfirmTOTP)
		mfar.Delete("/totp", strictLimiter, middleware.Timeout(5*time.Second), userMgmtHandler.DisableTOTP)

		// Backup codes (rate limited)
		// POST /me/mfa/backup-codes → Generate new backup codes
		mfar.Post("/backup-codes", strictLimiter, middleware.Timeout(5*time.Second), userMgmtHandler.GenerateBackupCodes)

		// Login history & audit
		// GET /me/login-history  → Login attempts history
		// GET /me/security-audit → Security audit trail
		router.Get("/login-history", middleware.Timeout(5*time.Second), userMgmtHandler.GetLoginHistory)
		router.Get("/security-audit", middleware.Timeout(5*time.Second), userMgmtHandler.GetSecurityAudit)
	})
}
