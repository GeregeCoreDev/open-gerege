// Package http provides implementation for http
//
// File: wire_security.go
// Description: implementation for http
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package http

import (
	"time"

	"templatev25/internal/middleware"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"

	fiberprometheus "github.com/ansrivas/fiberprometheus/v2"
	fbhelmet "github.com/gofiber/fiber/v2/middleware/helmet"
	fbrecover "github.com/gofiber/fiber/v2/middleware/recover"
	fbrequestid "github.com/gofiber/fiber/v2/middleware/requestid"
)

// ApplyMiddlewares wires common middlewares.
func ApplyMiddlewares(app *fiber.App, cfg *config.Config, logg *zap.Logger, apiLogRepo ...interface{}) {
	var repo interface{}
	if len(apiLogRepo) > 0 {
		repo = apiLogRepo[0]
	}

	isProduction := cfg.Server.ENV == "production" || cfg.Server.ENV == "prod"

	// ---- Core Recovery & Request ID ----
	app.Use(fbrecover.New())
	app.Use(fbrequestid.New())
	app.Use(fbhelmet.New())

	// ---- Distributed Tracing (OpenTelemetry) ----
	// Creates spans for each request with trace context propagation
	app.Use(middleware.Tracing())

	// ---- HSTS (Production only) ----
	// Forces HTTPS for all future requests
	if isProduction {
		app.Use(middleware.HSTS(middleware.DefaultHSTSConfig()))
	}

	// ---- HTTPS Redirect (Production only) ----
	// Enable via environment: APP_FORCE_HTTPS=true
	if isProduction {
		app.Use(middleware.HTTPSRedirect(true))
	}

	// CORS (cookie-compatible)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization,X-CSRF-Token",
		AllowCredentials: cfg.CORS.AllowCredentials,
	}))

	// ---- CSRF Protection ----
	// Protects against Cross-Site Request Forgery attacks
	csrfConfig := middleware.DefaultCSRFConfig(isProduction)
	app.Use(middleware.CSRF(csrfConfig))

	// Security headers
	app.Use(middleware.SecurityHeaders())

	// Body size limit ~2MB (adjust via env if you want)
	app.Use(middleware.BodySizeLimit(2 * 1024 * 1024))

	// Rate limiter: 100 req/min per user/IP
	app.Use(middleware.RateLimiter(100, time.Minute))

	// Response compression (gzip, deflate, brotli)
	// Reduces response size by 50-80% for JSON/text responses
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // Fast compression, good for API responses
	}))

	// Prometheus metrics
	p := fiberprometheus.New(cfg.Server.Name)
	p.RegisterAt(app, "/metrics")
	app.Use(p.Middleware)

	// Protect /metrics
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/metrics" {
			return middleware.IPAllow(cfg.Server.MetricsAllowCIDRs)(c)
		}
		return c.Next()
	})

	// Request context propagation (request_id, logger-ийг context руу дамжуулна)
	app.Use(middleware.RequestContext(logg))

	// Access logger
	if repo != nil {
		app.Use(middleware.RequestLogger(logg, repo.(repository.APILogRepository)))
	} else {
		app.Use(middleware.RequestLogger(logg))
	}

}
