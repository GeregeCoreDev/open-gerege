// Package main provides implementation for main
//
// File: main.go
// Description: implementation for main
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package main нь Template Backend API-ийн entry point юм.

Энэ файл нь:
  - Тохиргоо ачаалах (config.Load)
  - Өгөгдлийн сантай холбогдох (PostgreSQL + GORM)
  - Fiber веб серверийг эхлүүлэх
  - Middleware-үүдийг идэвхжүүлэх
  - Route-уудыг бүртгэх
  - Graceful shutdown хийх

Ажиллуулах:

	go run cmd/server/main.go

Эсвэл:

	make run
*/
package main

// Swagger API documentation metadata.
// swag init командаар docs/ хавтаст swagger.json, swagger.yaml үүснэ.
//
// @title           Template Backend API
// @version         1.0
// @description     Fiber + Gorm + Auth + RBAC demo
// @BasePath        /api/v1
// @schemes         http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Internal packages
	appdep "templatev25/internal/app"  // Dependency injection container
	"templatev25/internal/db"          // Database connection (GORM + PostgreSQL)
	"templatev25/internal/http/router" // HTTP route definitions
	"templatev25/internal/middleware"  // HTTP middlewares
	"templatev25/internal/repository"  // Repository layer

	// External packages
	"git.gerege.mn/backend-packages/config"               // Configuration loading (Viper)
	"git.gerege.mn/backend-packages/logger"               // Structured logging (Zap)
	ssoclient "git.gerege.mn/backend-packages/sso-client" // SSO client

	// Swagger generated docs
	docs "templatev25/docs"

	// Internal HTTP package (middlewares apply)
	ihttp "templatev25/internal/http"

	// External packages
	"github.com/ansrivas/fiberprometheus/v2" // Prometheus middleware
	"github.com/gofiber/fiber/v2"            // Web framework
	"go.uber.org/zap"                        // Structured logging

	// Observability
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

// main нь application-ийн entry point функц.
// Дараах алхмуудыг гүйцэтгэнэ:
//  1. Configuration ачаалах (.env файл эсвэл environment variables)
//  2. Logger үүсгэх (development/production mode)
//  3. Observability init (Prometheus)
//  4. Database холболт
//  5. Swagger setup
//  6. Fiber app setup
//  7. Middlewares setup
//  8. App logic (Service/Repository) setup
//  9. Server start
//  10. Graceful shutdown
func main() {
	// ============================================================
	// STEP 1: Configuration ачаалах
	// ============================================================
	cfg := config.Load(".")
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	// ============================================================
	// STEP 2: Logger үүсгэх
	// ============================================================
	logg := logger.New(cfg.Server.ENV)

	// ============================================================
	// STEP 3: Observability (Prometheus)
	// ============================================================
	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	promExporter, err := prometheus.New()
	if err != nil {
		logg.Fatal("failed to initialize prometheus exporter", zap.Error(err))
	}
	provider := metric.NewMeterProvider(metric.WithReader(promExporter))
	otel.SetMeterProvider(provider)

	// ============================================================
	// STEP 4: Database холболт
	// ============================================================
	gormDB, err := db.NewPostgres(cfg)
	if err != nil {
		logg.Fatal("db init failed", zap.Error(err))
	}

	// ============================================================
	// STEP 5: Swagger documentation тохируулах
	// ============================================================
	if cfg.Docs.Enabled {
		docs.SwaggerInfo.Title = cfg.Docs.Title
		docs.SwaggerInfo.Version = cfg.Docs.Version
		docs.SwaggerInfo.BasePath = cfg.Docs.BasePath
	}

	// ============================================================
	// STEP 6: Fiber application үүсгэх
	// ============================================================
	app := fiber.New(fiber.Config{
		AppName:      cfg.Server.Name,
		ErrorHandler: middleware.ErrorHandler(logg),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	})

	// Add Prometheus middleware
	prometheusMiddleware := fiberprometheus.New(cfg.Server.Name)
	prometheusMiddleware.RegisterAt(app, "/metrics")
	app.Use(prometheusMiddleware.Middleware)

	// ============================================================
	// STEP 7: Middlewares идэвхжүүлэх
	// ============================================================
	apiLogRepo := repository.NewAPILogRepositoryWithConfig(gormDB, &cfg)
	ihttp.ApplyMiddlewares(app, &cfg, logg, apiLogRepo)

	// ============================================================
	// STEP 8: Auth cache үүсгэх
	// ============================================================
	authCache := ssoclient.NewCache(cfg.Auth.CacheTTL, cfg.Auth.CacheMax)

	// ============================================================
	// STEP 9: Dependencies inject хийх
	// ============================================================
	deps := appdep.NewDependencies(gormDB, &cfg, logg, authCache)

	// ============================================================
	// STEP 10: Routes бүртгэх
	// ============================================================
	router.MapV1(app, deps)

	// ============================================================
	// STEP 11: Server эхлүүлэх (non-blocking)
	// ============================================================
	go func() {
		addr := cfg.Server.Addr()
		logg.Info("starting server", zap.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			logg.Error("fiber.listen stopped", zap.Error(err))
		}
	}()

	// ============================================================
	// STEP 12: Graceful shutdown хүлээх
	// ============================================================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// ============================================================
	// STEP 13: Server зогсоох
	// ============================================================
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Println("server shutdown error:", err)
	}

	// ============================================================
	// STEP 14: Resources cleanup
	// ============================================================
	if sqlDB, err := gormDB.DB(); err == nil {
		_ = sqlDB.Close()
	}
	authCache.Stop()
}
