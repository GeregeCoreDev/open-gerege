// Package middleware provides implementation for middleware
//
// File: error.go
// Description: implementation for middleware
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package middleware нь HTTP middleware-уудыг агуулна.

Энэ файл нь global error handler-ийг тодорхойлно.
Бүх unhandled error болон panic-ийг барьж, стандарт JSON response буцаана.

Error handling flow:
 1. Handler error буцаана (return err)
 2. ErrorHandler барьж авна
 3. Error code, message тодорхойлно
 4. Log бичнэ
 5. JSON response буцаана

Response format:

	{
	    "code": "BAD_REQUEST",
	    "request_id": "uuid",
	    "message": "validation failed"
	}
*/
package middleware

import (
	"errors" // Error type checking

	"git.gerege.mn/backend-packages/ctx"  // Request ID helper
	"git.gerege.mn/backend-packages/resp" // Response struct

	"github.com/gofiber/fiber/v2" // Web framework
	"go.uber.org/zap"             // Structured logging
)

// ============================================================
// ERROR HANDLER
// ============================================================

// ErrorHandler нь global error handler буцаана.
// Fiber app-д ErrorHandler тохируулахад ашиглана.
//
// Parameters:
//   - log: Zap logger
//
// Returns:
//   - fiber.ErrorHandler: Error handler function
//
// Ашиглалт:
//
//	app := fiber.New(fiber.Config{
//	    ErrorHandler: middleware.ErrorHandler(log),
//	})
//
// Log format:
//
//	{
//	    "level": "error",
//	    "msg": "http_error",
//	    "status": 400,
//	    "method": "POST",
//	    "path": "/api/user",
//	    "ip": "127.0.0.1",
//	    "error": "validation failed",
//	    "req_id": "uuid",
//	    "user_id": 123
//	}
func ErrorHandler(log *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Default values (500 Internal Server Error)
		code := fiber.StatusInternalServerError
		msg := "internal server error"

		// ============================================================
		// STEP 1: Fiber error шалгах
		// ============================================================
		// Fiber error бол түүний code, message авна
		// Жишээ: fiber.NewError(400, "bad request")
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
			msg = e.Message
		}

		// ============================================================
		// STEP 2: Request metadata авах
		// ============================================================
		// Logging-д ашиглана
		reqID := ctx.RequestID(c)

		// User ID авах (authenticated бол)
		userID := int(0)
		if v, ok := ctx.GetValue[int](c.UserContext(), ctx.KeyUserID); ok {
			userID = v
		}

		// ============================================================
		// STEP 3: Log бичих
		// ============================================================
		// Structured log: JSON format-аар гарна
		log.Error("http_error",
			zap.Int("status", code),
			zap.String("method", c.Method()),
			zap.String("path", c.OriginalURL()),
			zap.String("ip", c.IP()),
			zap.String("error", err.Error()),
			zap.String("req_id", reqID),
			zap.Int("user_id", userID),
		)

		// ============================================================
		// STEP 4: JSON response буцаах
		// ============================================================
		return c.Status(code).JSON(resp.APIResponse{
			Code:      httpStatusToCode(code),
			RequestID: reqID,
			Message:   msg,
		})
	}
}

// ============================================================
// HTTP STATUS TO CODE
// ============================================================

// httpStatusToCode нь HTTP status code-ийг API response code руу хөрвүүлнэ.
//
// Parameters:
//   - status: HTTP status code (400, 401, 403, 404, 500, etc.)
//
// Returns:
//   - string: API response code (BAD_REQUEST, UNAUTHORIZED, etc.)
//
// Mapping:
//   - 400 → BAD_REQUEST
//   - 401 → UNAUTHORIZED
//   - 403 → FORBIDDEN
//   - 404 → NOT_FOUND
//   - 405 → METHOD_NOT_ALLOWED
//   - 408 → REQUEST_TIMEOUT
//   - 409 → CONFLICT
//   - 413 → PAYLOAD_TOO_LARGE
//   - 422 → VALIDATION_ERROR
//   - 429 → TOO_MANY_REQUESTS
//   - 500 → INTERNAL_ERROR
//   - 502 → BAD_GATEWAY
//   - 503 → SERVICE_UNAVAILABLE
//   - 504 → GATEWAY_TIMEOUT
func httpStatusToCode(status int) string {
	switch status {
	// 4xx Client Errors
	case 400:
		return "BAD_REQUEST"
	case 401:
		return "UNAUTHORIZED"
	case 403:
		return "FORBIDDEN"
	case 404:
		return "NOT_FOUND"
	case 405:
		return "METHOD_NOT_ALLOWED"
	case 408:
		return "REQUEST_TIMEOUT"
	case 409:
		return "CONFLICT"
	case 413:
		return "PAYLOAD_TOO_LARGE"
	case 422:
		return "VALIDATION_ERROR"
	case 429:
		return "TOO_MANY_REQUESTS"
	// 5xx Server Errors
	case 500:
		return "INTERNAL_ERROR"
	case 502:
		return "BAD_GATEWAY"
	case 503:
		return "SERVICE_UNAVAILABLE"
	case 504:
		return "GATEWAY_TIMEOUT"
	default:
		// Categorize unknown status codes
		if status >= 400 && status < 500 {
			return "CLIENT_ERROR"
		}
		return "INTERNAL_ERROR"
	}
}
