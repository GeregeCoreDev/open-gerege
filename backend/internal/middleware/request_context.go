// Package middleware provides implementation for middleware
//
// File: request_context.go
// Description: Request context propagation middleware
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package middleware нь HTTP middleware-уудыг агуулна.

Энэ файл нь request context propagation middleware-ийг тодорхойлно.
Request ID болон бусад мэдээллийг context-д хадгалж, service layer руу дамжуулна.

Features:
  - Request ID context-д хадгалах
  - Logger-д request ID нэмэх
  - Service layer-д context-ээс request ID авах боломж
*/
package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ============================================================
// CONTEXT KEYS
// ============================================================

// ContextKey нь context-д хадгалах түлхүүрийн төрөл
type ContextKey string

const (
	// KeyRequestID нь request ID-ийн context key
	KeyRequestID ContextKey = "request_id"
	// KeyLogger нь logger-ийн context key
	KeyLogger ContextKey = "logger"
)

// ============================================================
// REQUEST CONTEXT MIDDLEWARE
// ============================================================

// RequestContext нь request ID болон logger-ийг context-д хадгалах middleware буцаана.
// Service layer-д context-ээс request ID авч log-д нэмэх боломжтой болно.
//
// Parameters:
//   - log: Base Zap logger
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Context values:
//   - request_id: Unique request identifier (X-Request-Id header-ээс)
//   - logger: Request-specific logger with request_id field
//
// Ашиглалт:
//
//	app.Use(middleware.RequestContext(baseLogger))
//
//	// Service layer-д:
//	log := middleware.GetLogger(ctx)
//	log.Info("operation_completed")
func RequestContext(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Request ID авах (fbrequestid middleware-ээс)
		reqID := c.Locals("requestid")
		if reqID == nil {
			reqID = c.Get("X-Request-Id", "")
		}

		reqIDStr, ok := reqID.(string)
		if !ok || reqIDStr == "" {
			reqIDStr = "unknown"
		}

		// Request-specific logger үүсгэх
		reqLog := log.With(zap.String("request_id", reqIDStr))

		// Context-д хадгалах
		ctx := c.UserContext()
		ctx = context.WithValue(ctx, KeyRequestID, reqIDStr)
		ctx = context.WithValue(ctx, KeyLogger, reqLog)

		// Fiber context-д буцаах
		c.SetUserContext(ctx)

		return c.Next()
	}
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

// GetRequestID нь context-ээс request ID авна.
//
// Parameters:
//   - ctx: Context (Fiber context.UserContext() ашиглан авна)
//
// Returns:
//   - string: Request ID (байхгүй бол "unknown")
//
// Ашиглалт:
//
//	reqID := middleware.GetRequestID(ctx)
func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(KeyRequestID).(string); ok {
		return reqID
	}
	return "unknown"
}

// GetLogger нь context-ээс request-specific logger авна.
// Logger нь request_id талбартай тул бүх log-д request ID орно.
//
// Parameters:
//   - ctx: Context (Fiber context.UserContext() ашиглан авна)
//
// Returns:
//   - *zap.Logger: Request-specific logger (байхгүй бол nil)
//
// Ашиглалт:
//
//	log := middleware.GetLogger(ctx)
//	if log != nil {
//	    log.Info("operation_completed", zap.String("user_id", userID))
//	}
func GetLogger(ctx context.Context) *zap.Logger {
	if log, ok := ctx.Value(KeyLogger).(*zap.Logger); ok {
		return log
	}
	return nil
}

// LoggerOrDefault нь context-ээс logger авах, байхгүй бол default logger буцаана.
//
// Parameters:
//   - ctx: Context
//   - defaultLog: Default logger (context-д байхгүй бол ашиглана)
//
// Returns:
//   - *zap.Logger: Logger instance
//
// Ашиглалт:
//
//	log := middleware.LoggerOrDefault(ctx, s.log)
//	log.Info("operation_completed")
func LoggerOrDefault(ctx context.Context, defaultLog *zap.Logger) *zap.Logger {
	if log := GetLogger(ctx); log != nil {
		return log
	}
	return defaultLog
}
