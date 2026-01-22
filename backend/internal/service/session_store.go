// Package service provides implementation for service
//
// File: session_store.go
// Description: Redis-backed session storage for local authentication
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"templatev25/internal/config"

	"github.com/redis/go-redis/v9"
)

// SessionStore defines the interface for session storage
type SessionStore interface {
	// Session management
	Create(ctx context.Context, session *SessionData) error
	Get(ctx context.Context, sessionID string) (*SessionData, error)
	Update(ctx context.Context, session *SessionData) error
	Delete(ctx context.Context, sessionID string) error
	Refresh(ctx context.Context, sessionID string, newExpiry time.Time) error

	// User session management
	GetUserSessions(ctx context.Context, userID int) ([]string, error)
	DeleteAllUserSessions(ctx context.Context, userID int) error

	// MFA token management (temporary storage during MFA flow)
	StoreMFAToken(ctx context.Context, token string, data *MFAPendingData, ttl time.Duration) error
	GetMFAToken(ctx context.Context, token string) (*MFAPendingData, error)
	DeleteMFAToken(ctx context.Context, token string) error

	// Health check
	Ping(ctx context.Context) error

	// Close connection
	Close() error
}

// SessionData represents a user session stored in Redis
type SessionData struct {
	SessionID      string    `json:"session_id"`
	UserID         int       `json:"user_id"`
	Email          string    `json:"email"`
	IPAddress      string    `json:"ip_address"`
	UserAgent      string    `json:"user_agent"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	LastActivityAt time.Time `json:"last_activity_at"`
}

// MFAPendingData represents temporary data during MFA verification
type MFAPendingData struct {
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RedisSessionStore implements SessionStore using Redis
type RedisSessionStore struct {
	client *redis.Client
	prefix string
}

// Redis key prefixes
const (
	sessionPrefix     = "session:"
	userSessionPrefix = "user:sessions:"
	mfaTokenPrefix    = "mfa:token:"
)

// NewRedisSessionStore creates a new Redis session store with a pre-created Redis client
func NewRedisSessionStore(client *redis.Client, prefix string, defaultTTL time.Duration) *RedisSessionStore {
	if prefix == "" {
		prefix = "auth:"
	}
	return &RedisSessionStore{
		client: client,
		prefix: prefix,
	}
}

// NewRedisSessionStoreFromConfig creates a new Redis session store from config
func NewRedisSessionStoreFromConfig(cfg *config.RedisConfig) (*RedisSessionStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisSessionStore{
		client: client,
		prefix: "auth:",
	}, nil
}

// ============================================================
// SESSION MANAGEMENT
// ============================================================

// Create stores a new session in Redis
func (s *RedisSessionStore) Create(ctx context.Context, session *SessionData) error {
	// Serialize session data
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Calculate TTL
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}

	// Store session
	key := s.sessionKey(session.SessionID)
	if err := s.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to store session: %w", err)
	}

	// Add to user's session set
	userKey := s.userSessionsKey(session.UserID)
	if err := s.client.SAdd(ctx, userKey, session.SessionID).Err(); err != nil {
		return fmt.Errorf("failed to add to user sessions: %w", err)
	}

	// Set expiry on user sessions set (longer than session TTL to handle cleanup)
	s.client.Expire(ctx, userKey, ttl+24*time.Hour)

	return nil
}

// Get retrieves a session from Redis
func (s *RedisSessionStore) Get(ctx context.Context, sessionID string) (*SessionData, error) {
	key := s.sessionKey(sessionID)
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Session not found
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session SessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// Update updates an existing session in Redis
func (s *RedisSessionStore) Update(ctx context.Context, session *SessionData) error {
	// Serialize session data
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Calculate remaining TTL
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}

	// Update session
	key := s.sessionKey(session.SessionID)
	if err := s.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// Delete removes a session from Redis
func (s *RedisSessionStore) Delete(ctx context.Context, sessionID string) error {
	// Get session first to get userID for cleanup
	session, err := s.Get(ctx, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return nil // Already deleted
	}

	// Delete session
	key := s.sessionKey(sessionID)
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Remove from user's session set
	userKey := s.userSessionsKey(session.UserID)
	s.client.SRem(ctx, userKey, sessionID)

	return nil
}

// Refresh extends the session expiry
func (s *RedisSessionStore) Refresh(ctx context.Context, sessionID string, newExpiry time.Time) error {
	session, err := s.Get(ctx, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found")
	}

	session.ExpiresAt = newExpiry
	session.LastActivityAt = time.Now()

	return s.Update(ctx, session)
}

// ============================================================
// USER SESSION MANAGEMENT
// ============================================================

// GetUserSessions returns all session IDs for a user
func (s *RedisSessionStore) GetUserSessions(ctx context.Context, userID int) ([]string, error) {
	key := s.userSessionsKey(userID)
	sessionIDs, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	// Filter out expired sessions
	var validSessions []string
	for _, sessionID := range sessionIDs {
		session, err := s.Get(ctx, sessionID)
		if err != nil {
			continue
		}
		if session != nil && time.Now().Before(session.ExpiresAt) {
			validSessions = append(validSessions, sessionID)
		} else {
			// Clean up expired session from set
			s.client.SRem(ctx, key, sessionID)
		}
	}

	return validSessions, nil
}

// DeleteAllUserSessions removes all sessions for a user
func (s *RedisSessionStore) DeleteAllUserSessions(ctx context.Context, userID int) error {
	// Get all session IDs
	sessionIDs, err := s.GetUserSessions(ctx, userID)
	if err != nil {
		return err
	}

	// Delete each session
	for _, sessionID := range sessionIDs {
		key := s.sessionKey(sessionID)
		s.client.Del(ctx, key)
	}

	// Delete user sessions set
	userKey := s.userSessionsKey(userID)
	s.client.Del(ctx, userKey)

	return nil
}

// ============================================================
// MFA TOKEN MANAGEMENT
// ============================================================

// StoreMFAToken stores temporary MFA verification data
func (s *RedisSessionStore) StoreMFAToken(ctx context.Context, token string, data *MFAPendingData, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal MFA data: %w", err)
	}

	key := s.mfaTokenKey(token)
	if err := s.client.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		return fmt.Errorf("failed to store MFA token: %w", err)
	}

	return nil
}

// GetMFAToken retrieves MFA verification data
func (s *RedisSessionStore) GetMFAToken(ctx context.Context, token string) (*MFAPendingData, error) {
	key := s.mfaTokenKey(token)
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Token not found or expired
		}
		return nil, fmt.Errorf("failed to get MFA token: %w", err)
	}

	var mfaData MFAPendingData
	if err := json.Unmarshal(data, &mfaData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal MFA data: %w", err)
	}

	return &mfaData, nil
}

// DeleteMFAToken removes MFA verification data
func (s *RedisSessionStore) DeleteMFAToken(ctx context.Context, token string) error {
	key := s.mfaTokenKey(token)
	return s.client.Del(ctx, key).Err()
}

// ============================================================
// UTILITY METHODS
// ============================================================

// Ping checks Redis connectivity
func (s *RedisSessionStore) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (s *RedisSessionStore) Close() error {
	return s.client.Close()
}

// Key helpers
func (s *RedisSessionStore) sessionKey(sessionID string) string {
	return s.prefix + sessionPrefix + sessionID
}

func (s *RedisSessionStore) userSessionsKey(userID int) string {
	return s.prefix + userSessionPrefix + strconv.Itoa(userID)
}

func (s *RedisSessionStore) mfaTokenKey(token string) string {
	return s.prefix + mfaTokenPrefix + token
}
