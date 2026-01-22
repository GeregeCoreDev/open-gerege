// Package router provides implementation for router
//
// File: router.go
// Description: implementation for router
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package router нь application-ийн бүх HTTP route-уудыг тодорхойлно.

Route бүтэц:
  - Public routes: Authentication шаардлагагүй
  - Protected routes: Auth middleware шаардлагатай

Endpoint groups:

	/health              - Health check
	/docs/*              - Swagger UI
	/auth/*              - Authentication (login, logout, callback)
	/user/*              - User management
	/user-role/*         - User-Role assignments
	/system/*            - System management
	/module/*            - Module management
	/permission/*        - Permission management
	/role/*              - Role management
	/organization/*      - Organization management
	/terminal/*          - Terminal management
	/notification/*      - Notification management
	/news/*              - News management
	/verify/*            - Verification (DAN, email, phone)
	/room/*              - Video conference rooms
	/tpay/*              - Terminal payment
	/chat/*              - Chat items

Ашиглалт:

	deps := app.NewDependencies(...)
	router.MapV1(fiberApp, deps)
*/
package router

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"templatev25/internal/app"        // Dependency container
	"templatev25/internal/auth"       // Auth middleware
	"templatev25/internal/middleware" // Middleware

	"git.gerege.mn/backend-packages/resp" // Response helpers

	"github.com/gofiber/fiber/v2"        // Web framework
	swagger "github.com/gofiber/swagger" // Swagger UI middleware
	"gorm.io/gorm"                       // ORM (for health check)
)

// ============================================================
// HEALTH CHECK CACHE (Performance optimization)
// ============================================================

var (
	healthCache     atomic.Value  // Cached health result
	healthCacheTime atomic.Int64  // Last cache time (unix seconds)
	healthCacheMu   sync.Mutex    // Mutex for cache update
)

const healthCacheTTL = 5 // Cache TTL in seconds

// ============================================================
// MAIN ROUTE MAPPING FUNCTION
// ============================================================

// MapV1 нь application-ийн бүх route-уудыг бүртгэнэ.
// V1 нь API version 1 гэсэн утгатай.
//
// Parameters:
//   - app: Fiber application instance
//   - d: Dependencies container (repositories, services, config, etc.)
//
// Route structure:
//
//	┌──────────────────────────────────────────────────────────┐
//	│                     PUBLIC ROUTES                         │
//	├──────────────────────────────────────────────────────────┤
//	│  GET  /health     → Health check (DB ping)               │
//	│  GET  /docs/*     → Swagger UI                           │
//	└──────────────────────────────────────────────────────────┘
//	┌──────────────────────────────────────────────────────────┐
//	│                     AUTH ROUTES                           │
//	├──────────────────────────────────────────────────────────┤
//	│  GET  /auth/login      → SSO redirect                    │
//	│  GET  /auth/callback   → OAuth2 callback                 │
//	│  POST /auth/logout     → Logout                          │
//	│  POST /auth/google     → Google OAuth                    │
//	│  GET  /auth/verify     → Token verification              │
//	│  POST /auth/org/change → Change organization (protected) │
//	└──────────────────────────────────────────────────────────┘
//	┌──────────────────────────────────────────────────────────┐
//	│                   PROTECTED ROUTES                        │
//	├──────────────────────────────────────────────────────────┤
//	│  /user/*           → User CRUD + profile                 │
//	│  /user-role/*      → User-Role management                │
//	│  /system/*         → System CRUD                         │
//	│  /module/*         → Module CRUD                         │
//	│  /permission/*     → Permission CRUD                     │
//	│  /role/*           → Role CRUD + permissions             │
//	│  /organization/*   → Organization CRUD                   │
//	│  /terminal/*       → Terminal CRUD                       │
//	│  ... (more routes)                                       │
//	└──────────────────────────────────────────────────────────┘
func MapV1(app *fiber.App, d *app.Dependencies) {

	// ============================================================
	// PUBLIC ROUTES (Authentication шаардлагагүй)
	// ============================================================
	pub := app.Group("/")

	// Health check endpoint
	// Database connection-ийг шалгана (2 секундын timeout-тэй)
	// Response: {"code": "OK", "data": {"status": "ok"}}
	pub.Get("/health", healthHandler(d.DB))

	// Swagger UI (зөвхөн Docs.Enabled=true үед)
	// URL: /docs/index.html
	// Swagger JSON: /docs/doc.json
	if d.Cfg.Docs.Enabled {
		pub.Get("/docs/*", swagger.New(swagger.Config{
			Title: d.Cfg.Docs.Title,
		}))
	}

	// ============================================================
	// AUTH MIDDLEWARE
	// ============================================================
	// Protected route-уудад хэрэглэгчийн session-ийг шалгана.
	// Cookie-д "sid" байвал түүнийг validate хийнэ.
	// Session invalid бол 401 Unauthorized буцаана.
	requireAuth := auth.Require(d.Cfg, d.Log, d.AuthCache)

	// ============================================================
	// V1 API ROUTES
	// ============================================================
	// Pagination хязгаарлалт нэмэх (max 100 бичлэг)
	v1 := app.Group("/", middleware.PaginationLimit(100))

	// ------------------------------------------------------------
	// AUTH ROUTES
	// ------------------------------------------------------------
	// Authentication-тай холбоотой endpoint-ууд.
	MapAuthRoutes(v1, d, requireAuth)

	// ------------------------------------------------------------
	// ME ROUTES (Current User)
	// ------------------------------------------------------------
	MapMeRoutes(v1, d, requireAuth)

	// ------------------------------------------------------------
	// USER ROUTES
	// ------------------------------------------------------------
	MapUserRoutes(v1, d, requireAuth)

	// ------------------------------------------------------------
	// SYSTEM, MODULE, PERMISSION, ACTION, ROLE, CLIENT ROUTES
	// ------------------------------------------------------------
	MapSystemRoutes(v1, d, requireAuth)

	// ------------------------------------------------------------
	// ORGANIZATION, ORGUSER, ORGTYPE ROUTES
	// ------------------------------------------------------------
	MapOrganizationRoutes(v1, d, requireAuth)

	// ------------------------------------------------------------
	// APP SERVICE ICON, APP SERVICE GROUP ROUTES
	// ------------------------------------------------------------
	MapAppIconRoutes(v1, d, requireAuth)

	// ------------------------------------------------------------
	// FILE ROUTES
	// ------------------------------------------------------------
	MapFileRoutes(v1, d, requireAuth)

	// ------------------------------------------------------------
	// NOTIFICATION ROUTES
	// ------------------------------------------------------------
	MapNotificationRoutes(v1, d, requireAuth)

	// ------------------------------------------------------------
	// NEWS ROUTES
	// ------------------------------------------------------------
	MapNewsRoutes(v1, d, requireAuth)

	// ------------------------------------------------------------
	// ROOM, CHAT ROUTES
	// ------------------------------------------------------------
	MapChatRoutes(v1, d, requireAuth)

	// ------------------------------------------------------------
	// API LOG ROUTES
	// ------------------------------------------------------------
	MapAPILogRoutes(v1, d, requireAuth)

	// ------------------------------------------------------------
	// TPAY ROUTES (Terminal Payment)
	// ------------------------------------------------------------
	// Терминал төлбөрийн API-г me_router.go файлд шилжүүлсэн.

	// ============================================================
	// 404 HANDLER
	// ============================================================
	// Бүртгэгдээгүй route-д 404 буцаана.
	app.All("/*", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "resource not found")
	})
}

// ============================================================
// HEALTH CHECK HANDLER
// ============================================================

// serverStartTime нь server эхэлсэн хугацаа (uptime тооцоолоход хэрэглэнэ)
var serverStartTime = time.Now()

// healthHandler нь database connection-ийг шалгаж, server-ийн төлөвийг буцаана.
//
// Returns:
//   - 200 OK: {"code": "OK", "data": {...}}
//   - 500 Error: {"code": "INTERNAL_ERROR", "message": "db_down"}
//
// Response data includes:
//   - status: "ok" or "degraded"
//   - uptime: Server uptime in seconds
//   - database: Database connection status
//   - timestamp: Current server time (RFC3339)
//
// Database ping timeout: 2 секунд
// Cached for 5 seconds to reduce database load under high traffic
func healthHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		now := time.Now().Unix()

		// Check cache first (fast path)
		if cached := healthCache.Load(); cached != nil {
			if now-healthCacheTime.Load() < healthCacheTTL {
				return resp.OK(c, cached)
			}
		}

		// Cache miss or expired - compute new result
		healthCacheMu.Lock()
		defer healthCacheMu.Unlock()

		// Double-check after acquiring lock
		if cached := healthCache.Load(); cached != nil {
			if now-healthCacheTime.Load() < healthCacheTTL {
				return resp.OK(c, cached)
			}
		}

		// 2 секундын timeout-тэй context үүсгэх
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		// Health check result
		result := fiber.Map{
			"status":    "ok",
			"uptime":    int64(time.Since(serverStartTime).Seconds()),
			"timestamp": time.Now().Format(time.RFC3339),
		}

		// GORM-оос underlying *sql.DB авах
		sqlDB, err := db.DB()
		if err != nil {
			result["status"] = "degraded"
			result["database"] = fiber.Map{
				"status": "error",
				"error":  "db_connection_error",
			}
			// Cache error result too (avoid DB hammering)
			healthCache.Store(result)
			healthCacheTime.Store(now)
			return resp.OK(c, result)
		}

		// Database ping хийх
		if err := sqlDB.PingContext(ctx); err != nil {
			result["status"] = "degraded"
			result["database"] = fiber.Map{
				"status": "error",
				"error":  "db_unreachable",
			}
			healthCache.Store(result)
			healthCacheTime.Store(now)
			return resp.OK(c, result)
		}

		// Database stats авах
		stats := sqlDB.Stats()
		result["database"] = fiber.Map{
			"status":      "ok",
			"open_conns":  stats.OpenConnections,
			"in_use":      stats.InUse,
			"idle":        stats.Idle,
			"max_open":    stats.MaxOpenConnections,
			"wait_count":  stats.WaitCount,
			"wait_time":   stats.WaitDuration.String(),
		}

		// Cache the result
		healthCache.Store(result)
		healthCacheTime.Store(now)

		return resp.OK(c, result)
	}
}
