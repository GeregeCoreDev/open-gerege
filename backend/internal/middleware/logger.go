// Package middleware provides implementation for middleware
//
// File: logger.go
// Description: implementation for middleware
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package middleware нь HTTP middleware-уудыг агуулна.

Энэ файл нь request logging middleware-ийг тодорхойлно.
Бүх HTTP request-ийг structured log format-аар бичнэ.

Features:
  - Request/response metadata (method, path, status, latency)
  - PII scrubbing (Authorization, Cookie headers masked)
  - Log level by status (5xx=Error, 4xx=Warn, else=Info)
  - User tracking (user_id, request_id)

Log format (JSON):

	{
	    "level": "info",
	    "msg": "http_request",
	    "method": "GET",
	    "path": "/api/users",
	    "route": "/api/users",
	    "status": 200,
	    "latency": "15.234ms",
	    "req_size": 0,
	    "res_size": 1234,
	    "ip": "127.0.0.1",
	    "user_agent": "Mozilla/5.0...",
	    "request_id": "uuid",
	    "user_id": 123
	}
*/
package middleware

import (
	"context"
	"encoding/json"
	"strings" // String manipulation
	"sync"
	"time" // Duration

	"templatev25/internal/domain"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/ctx" // Context helpers

	"github.com/gofiber/fiber/v2" // Web framework
	"go.uber.org/zap"             // Structured logging
	"gorm.io/datatypes"
)

// ============================================================
// LOG WORKER POOL (Goroutine leak prevention)
// ============================================================

const (
	logWorkerCount   = 5     // Number of worker goroutines
	logQueueSize     = 1000  // Buffer size for log queue
	logWriteTimeout  = 5 * time.Second
)

var (
	logQueue     chan logEntry
	logQueueOnce sync.Once
	logLogger    *zap.Logger
)

type logEntry struct {
	repo   repository.APILogRepository
	apiLog domain.APILog
}

// initLogWorkers starts the worker pool for async log writing.
// Called once when first log repo is provided.
func initLogWorkers(log *zap.Logger) {
	logQueueOnce.Do(func() {
		logQueue = make(chan logEntry, logQueueSize)
		logLogger = log

		// Start worker goroutines
		for i := 0; i < logWorkerCount; i++ {
			go logWorker()
		}
	})
}

// logWorker processes log entries from the queue
func logWorker() {
	for entry := range logQueue {
		ctx, cancel := context.WithTimeout(context.Background(), logWriteTimeout)
		if err := entry.repo.Create(ctx, entry.apiLog); err != nil {
			if logLogger != nil {
				logLogger.Error("failed to save api log to database", zap.Error(err))
			}
		}
		cancel()
	}
}

// ============================================================
// REQUEST LOGGER
// ============================================================

// RequestLogger нь HTTP request-ийг structured log format-аар бичих middleware буцаана.
//
// Log fields:
//   - method: HTTP method (GET, POST, etc.)
//   - path: Full URL path with query string
//   - route: Route template (/users/:id)
//   - status: HTTP status code
//   - latency: Request duration
//   - req_size: Request body size (bytes)
//   - res_size: Response body size (bytes)
//   - ip: Client IP address
//   - user_agent: Browser user agent
//   - request_id: Unique request identifier
//   - user_id: Authenticated user ID (if available)
//
// Additional fields on 4xx/5xx:
//   - authorization: Masked Authorization header
//   - cookie: Masked Cookie header
//
// Log levels:
//   - 5xx → Error
//   - 4xx → Warn
//   - 2xx/3xx → Info
//
// Parameters:
//   - log: Zap logger
//   - apiLogRepo: Optional APILog repository for database logging
//
// Returns:
//   - fiber.Handler: Middleware function
//
// Ашиглалт:
//
//	app.Use(middleware.RequestLogger(log))
//	app.Use(middleware.RequestLogger(log, apiLogRepo))
func RequestLogger(log *zap.Logger, apiLogRepo ...repository.APILogRepository) fiber.Handler {
	var repo repository.APILogRepository
	if len(apiLogRepo) > 0 {
		repo = apiLogRepo[0]
		// Initialize worker pool (only once)
		initLogWorkers(log)
	}
	return func(c *fiber.Ctx) error {
		// Request эхлэх цаг
		start := time.Now()

		// Handler chain ажиллуулах
		err := c.Next()

		// Latency тооцоолох
		lat := time.Since(start)

		// ============================================================
		// REQUEST METADATA
		// ============================================================
		method := c.Method()

		// Route template (параметртэй path)
		routePath := ""
		if r := c.Route(); r != nil {
			routePath = r.Path
		}

		// Full URL path
		path := c.OriginalURL()

		// Client IP
		ip := c.IP()

		// User agent
		ua := string(c.Request().Header.UserAgent())

		// Response status
		status := c.Response().StatusCode()

		// ============================================================
		// REQUEST/RESPONSE SIZES
		// ============================================================
		// Request body size
		reqSize := int64(len(c.Request().Body()))
		if reqSize == 0 {
			reqSize = int64(c.Request().Header.ContentLength())
		}

		// Response body size
		resSize := int64(c.Response().Header.ContentLength())
		if resSize <= 0 {
			// Fallback: actual body length
			if b := c.Response().Body(); b != nil {
				resSize = int64(len(b))
			}
		}

		// ============================================================
		// CONTEXT VALUES
		// ============================================================
		// Request ID (header эсвэл locals-оос)
		reqID := headerOrLocal(c, "X-Request-ID", "requestid")

		// User ID (authenticated бол)
		userID, _ := ctx.GetValue[int](c.UserContext(), ctx.KeyUserID)

		// ============================================================
		// MASKED HEADERS (PII PROTECTION)
		// ============================================================
		// Authorization, Cookie header-ийг маскална
		// Алдааны үед л хавсаргана (debugging-д тусална)
		reqAuth := maskHeader(string(c.Request().Header.Peek("Authorization")))
		reqCookie := maskHeader(string(c.Request().Header.Peek("Cookie")))

		// ============================================================
		// BUILD LOG FIELDS
		// ============================================================
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("route", routePath),
			zap.Int("status", status),
			zap.Duration("latency", lat),
			zap.Int64("req_size", reqSize),
			zap.Int64("res_size", resSize),
			zap.String("ip", ip),
			zap.String("user_agent", ua),
		}

		// Optional fields
		if reqID != "" {
			fields = append(fields, zap.String("request_id", reqID))
		}
		if userID != 0 {
			fields = append(fields, zap.Int("user_id", userID))
		}

		// 4xx/5xx үед masked headers нэмэх
		if status >= 400 {
			fields = append(fields, zap.String("authorization", reqAuth))
			fields = append(fields, zap.String("cookie", reqCookie))
		}

		// ============================================================
		// LOG BY LEVEL
		// ============================================================
		switch {
		case status >= 500:
			log.Error("http_request", fields...)
		case status >= 400:
			log.Warn("http_request", fields...)
		default:
			log.Info("http_request", fields...)
		}

		// ============================================================
		// DATABASE LOGGING (if repository provided)
		// ============================================================
		if repo != nil {
			// Prepare request body (if available)
			// Optimized: Use raw bytes directly, avoid double JSON serialization
			var reqBody datatypes.JSON
			if body := c.Body(); len(body) > 0 && len(body) < 10000 {
				// Check if it's valid JSON
				if json.Valid(body) {
					reqBody = body
				} else {
					// Wrap non-JSON as string
					if bodyBytes, err := json.Marshal(string(body)); err == nil {
						reqBody = bodyBytes
					}
				}
			}

			// Prepare response body (only for error responses 4xx/5xx)
			// Амжилттай API хариу үед response body бичихгүй (хурд, багтаамж хэмнэнэ)
			var resBody datatypes.JSON

			// Зөвхөн алдаатай үед response body-г бичнэ
			if status >= 400 {
				var responseBodyBytes []byte

				// Эхлээд locals-оос авах (хэрэв handler-ууд хадгалсан бол)
				if responseBodyVal, ok := c.Locals("response_body").([]byte); ok && len(responseBodyVal) > 0 {
					responseBodyBytes = responseBodyVal
				} else {
					// Fallback: Response().Body() ашиглах
					responseBodyBytes = c.Response().Body()
				}

				// Optimized: Use raw bytes directly, avoid double JSON serialization
				if len(responseBodyBytes) > 0 && len(responseBodyBytes) < 10000 {
					if json.Valid(responseBodyBytes) {
						resBody = responseBodyBytes
					} else {
						// Wrap non-JSON as string
						if bodyBytes, err := json.Marshal(string(responseBodyBytes)); err == nil {
							resBody = bodyBytes
						}
					}
				}
			}

			// Prepare query parameters
			var queries datatypes.JSON
			if len(c.Queries()) > 0 {
				if queryBytes, err := json.Marshal(c.Queries()); err == nil {
					queries = queryBytes
				}
			}

			// Prepare path parameters
			var params datatypes.JSON
			if len(c.AllParams()) > 0 {
				if paramBytes, err := json.Marshal(c.AllParams()); err == nil {
					params = paramBytes
				}
			}

			// Get username from context (if available)
			username := ""
			// Try to get username from context or locals
			if usernameVal, ok := c.Locals("username").(string); ok {
				username = usernameVal
			}

			// Get org_id from context (if available)
			var orgID *int64
			if orgIDVal, ok := ctx.GetValue[int](c.UserContext(), ctx.KeyOrgID); ok {
				orgIDVal64 := int64(orgIDVal)
				orgID = &orgIDVal64
			}

			// Create APILog entry
			apiLog := domain.APILog{
				OrgId: orgID,
				UserId: func() *int64 {
					if userID != 0 {
						userID64 := int64(userID)
						return &userID64
					} else {
						return nil
					}
				}(),
				Username:    username,
				Path:        path,
				Method:      method,
				Params:      params,
				Queries:     queries,
				Body:        reqBody,
				StatusCode:  status,
				Response:    resBody,
				LatencyMs:   lat.Milliseconds(),
				ReqSize:     reqSize,
				ResSize:     resSize,
				IP:          ip,
				CreatedDate: time.Now(),
			}

			// Save to database asynchronously via worker pool (don't block response)
			// Non-blocking send - if queue is full, log warning and drop
			select {
			case logQueue <- logEntry{repo: repo, apiLog: apiLog}:
				// Successfully queued
			default:
				// Queue full, log warning
				log.Warn("api log queue full, dropping log entry",
					zap.String("path", path),
					zap.String("method", method))
			}
		}

		return err
	}
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

// headerOrLocal нь header эсвэл locals-оос утга авна.
//
// Parameters:
//   - c: Fiber context
//   - header: HTTP header name
//   - local: Fiber locals key
//
// Returns:
//   - string: Value (хоосон бол "")
func headerOrLocal(c *fiber.Ctx, header, local string) string {
	// Header-ээс хайх
	if v := c.Get(header); v != "" {
		return v
	}
	// Locals-оос хайх
	if v, ok := c.Locals(local).(string); ok && v != "" {
		return v
	}
	return ""
}

// maskHeader нь мэдрэмтгий header утгыг маскална.
// PII хамгаалалт: log-д бүтэн token харагдахгүй.
//
// Parameters:
//   - v: Header value
//
// Returns:
//   - string: Masked value
//
// Жишээ:
//
//	maskHeader("Bearer abcdefghijkl")  // "Bearer abcd•••jkl"
//	maskHeader("sid=abc123")           // "sid=•••"
func maskHeader(v string) string {
	if v == "" {
		return ""
	}

	lower := strings.ToLower(v)

	// Bearer token
	if strings.HasPrefix(lower, "bearer ") {
		return "Bearer " + maskString(v[7:])
	}

	// Cookie, Basic, бусад
	return maskString(v)
}

// maskString нь string-ийн дундыг маскална.
//
// Parameters:
//   - s: Original string
//
// Returns:
//   - string: Masked string
//
// Rules:
//   - 6 тэмдэгтээс бага: "•••"
//   - Бусад: эхний 4 + "•••" + сүүлийн 3
//
// Жишээ:
//
//	maskString("abc")               // "•••"
//	maskString("abcdefghij")        // "abcd•••hij"
func maskString(s string) string {
	n := len(s)
	if n <= 6 {
		return "•••"
	}
	return s[:4] + "•••" + s[n-3:]
}
