// Package router provides implementation for router
//
// File: session_adapter.go
// Description: Adapter to bridge service.SessionStore and middleware.SessionStore
package router

import (
	"context"

	"templatev25/internal/middleware"
	"templatev25/internal/service"
)

// SessionStoreAdapter adapts service.SessionStore to middleware.SessionStore interface
// This is needed to avoid import cycles between service and middleware packages
type SessionStoreAdapter struct {
	store service.SessionStore
}

// NewSessionStoreAdapter creates a new adapter
func NewSessionStoreAdapter(store service.SessionStore) *SessionStoreAdapter {
	return &SessionStoreAdapter{store: store}
}

// Get retrieves a session and converts it to middleware.SessionData
func (a *SessionStoreAdapter) Get(ctx context.Context, sessionID string) (*middleware.SessionData, error) {
	session, err := a.store.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}

	return &middleware.SessionData{
		SessionID:      session.SessionID,
		UserID:         session.UserID,
		Email:          session.Email,
		IPAddress:      session.IPAddress,
		UserAgent:      session.UserAgent,
		CreatedAt:      session.CreatedAt,
		ExpiresAt:      session.ExpiresAt,
		LastActivityAt: session.LastActivityAt,
	}, nil
}

// Update converts middleware.SessionData back to service.SessionData and updates
func (a *SessionStoreAdapter) Update(ctx context.Context, session *middleware.SessionData) error {
	serviceSession := &service.SessionData{
		SessionID:      session.SessionID,
		UserID:         session.UserID,
		Email:          session.Email,
		IPAddress:      session.IPAddress,
		UserAgent:      session.UserAgent,
		CreatedAt:      session.CreatedAt,
		ExpiresAt:      session.ExpiresAt,
		LastActivityAt: session.LastActivityAt,
	}
	return a.store.Update(ctx, serviceSession)
}
