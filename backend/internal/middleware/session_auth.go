// Package middleware provides HTTP middlewares
//
// File: session_auth.go
// Description: Session-based authentication middleware for local auth
package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SessionData represents session information needed by the middleware
// This is a local interface to avoid import cycles with the service package
type SessionData struct {
	SessionID      string
	UserID         int
	Email          string
	IPAddress      string
	UserAgent      string
	CreatedAt      time.Time
	ExpiresAt      time.Time
	LastActivityAt time.Time
}

// SessionStore interface defines methods needed by the session auth middleware
// This is a local interface to avoid import cycles with the service package
type SessionStore interface {
	Get(ctx context.Context, sessionID string) (*SessionData, error)
	Update(ctx context.Context, session *SessionData) error
}

// SessionAuth creates a middleware that validates sessions from Redis
func SessionAuth(sessionStore SessionStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		token := extractBearerToken(c)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "missing authorization token",
			})
		}

		// Get session from Redis
		session, err := sessionStore.Get(c.UserContext(), token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "session lookup failed",
			})
		}

		if session == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "invalid or expired session",
			})
		}

		// Check expiry
		if time.Now().After(session.ExpiresAt) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "session expired",
			})
		}

		// Update last activity (async, don't block)
		go func() {
			session.LastActivityAt = time.Now()
			sessionStore.Update(c.UserContext(), session)
		}()

		// Set session info in context
		c.Locals("session_id", session.SessionID)
		c.Locals("user_id", session.UserID)
		c.Locals("email", session.Email)
		c.Locals("session", session)

		return c.Next()
	}
}

// OptionalSessionAuth creates a middleware that optionally validates sessions
// If no token is provided, the request continues without auth info
// If a token is provided but invalid, the request continues without auth info
func OptionalSessionAuth(sessionStore SessionStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		token := extractBearerToken(c)
		if token == "" {
			return c.Next()
		}

		// Get session from Redis
		session, err := sessionStore.Get(c.UserContext(), token)
		if err != nil || session == nil {
			return c.Next()
		}

		// Check expiry
		if time.Now().After(session.ExpiresAt) {
			return c.Next()
		}

		// Update last activity (async, don't block)
		go func() {
			session.LastActivityAt = time.Now()
			sessionStore.Update(c.UserContext(), session)
		}()

		// Set session info in context
		c.Locals("session_id", session.SessionID)
		c.Locals("user_id", session.UserID)
		c.Locals("email", session.Email)
		c.Locals("session", session)

		return c.Next()
	}
}

// RequireLocalAuth is a convenience middleware that ensures local auth session
// This can be used in combination with SSO auth to support both auth methods
func RequireLocalAuth(sessionStore SessionStore) fiber.Handler {
	return SessionAuth(sessionStore)
}

// Helper to extract bearer token from Authorization header
func extractBearerToken(c *fiber.Ctx) string {
	auth := c.Get("Authorization")
	if auth == "" {
		return ""
	}

	// Check for "Bearer " prefix
	if len(auth) > 7 && strings.ToLower(auth[:7]) == "bearer " {
		return auth[7:]
	}

	return ""
}
