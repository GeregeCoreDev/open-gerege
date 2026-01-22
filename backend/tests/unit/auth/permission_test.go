// Package auth provides implementation for auth
//
// File: permission_test.go
// Description: Tests for permission middleware
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package auth_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"templatev25/internal/auth"

	ssoclient "git.gerege.mn/backend-packages/sso-client"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// MOCK PERMISSION CHECKER
// ============================================================

type mockPermissionChecker struct {
	mock.Mock
}

func (m *mockPermissionChecker) HasPermission(ctx context.Context, userID int, permissionCode string) (bool, error) {
	args := m.Called(ctx, userID, permissionCode)
	return args.Bool(0), args.Error(1)
}

func (m *mockPermissionChecker) GetUserPermissions(ctx context.Context, userID int) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// ============================================================
// TEST REQUIRE PERMISSION
// ============================================================

func TestRequirePermission(t *testing.T) {
	tests := []struct {
		name           string
		permissionCode string
		userID         int
		mockSetup      func(*mockPermissionChecker)
		wantStatus     int
	}{
		{
			name:           "success - user has permission",
			permissionCode: "admin.role.create",
			userID:         1,
			mockSetup: func(m *mockPermissionChecker) {
				m.On("HasPermission", mock.Anything, 1, "admin.role.create").Return(true, nil)
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name:           "forbidden - user lacks permission",
			permissionCode: "admin.role.create",
			userID:         1,
			mockSetup: func(m *mockPermissionChecker) {
				m.On("HasPermission", mock.Anything, 1, "admin.role.create").Return(false, nil)
			},
			wantStatus: fiber.StatusForbidden,
		},
		{
			name:           "forbidden - no user ID",
			permissionCode: "admin.role.create",
			userID:         0,
			mockSetup:      func(m *mockPermissionChecker) {},
			wantStatus:     fiber.StatusForbidden,
		},
		{
			name:           "forbidden - permission check error",
			permissionCode: "admin.role.create",
			userID:         1,
			mockSetup: func(m *mockPermissionChecker) {
				m.On("HasPermission", mock.Anything, 1, "admin.role.create").Return(false, errors.New("db error"))
			},
			wantStatus: fiber.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChecker := &mockPermissionChecker{}
			tt.mockSetup(mockChecker)

			app := fiber.New()
			app.Use(func(c *fiber.Ctx) error {
				// Simulate SSO claims being set
				if tt.userID > 0 {
					c.Locals(ssoclient.LocalsClaims, &ssoclient.Claims{UserID: tt.userID})
				}
				return c.Next()
			})
			app.Get("/test", auth.RequirePermission(mockChecker, tt.permissionCode), func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockChecker.AssertExpectations(t)
		})
	}
}

// ============================================================
// TEST REQUIRE ANY PERMISSION
// ============================================================

func TestRequireAnyPermission(t *testing.T) {
	tests := []struct {
		name            string
		permissionCodes []string
		userID          int
		mockSetup       func(*mockPermissionChecker)
		wantStatus      int
	}{
		{
			name:            "success - user has one of the permissions",
			permissionCodes: []string{"admin.role.create", "admin.role.read"},
			userID:          1,
			mockSetup: func(m *mockPermissionChecker) {
				m.On("GetUserPermissions", mock.Anything, 1).Return([]string{"admin.role.read"}, nil)
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name:            "success - empty permission list allows all",
			permissionCodes: []string{},
			userID:          1,
			mockSetup:       func(m *mockPermissionChecker) {},
			wantStatus:      fiber.StatusOK,
		},
		{
			name:            "forbidden - user has none of the permissions",
			permissionCodes: []string{"admin.role.create", "admin.role.delete"},
			userID:          1,
			mockSetup: func(m *mockPermissionChecker) {
				m.On("GetUserPermissions", mock.Anything, 1).Return([]string{"admin.user.read"}, nil)
			},
			wantStatus: fiber.StatusForbidden,
		},
		{
			name:            "forbidden - permission check error",
			permissionCodes: []string{"admin.role.create"},
			userID:          1,
			mockSetup: func(m *mockPermissionChecker) {
				m.On("GetUserPermissions", mock.Anything, 1).Return(nil, errors.New("db error"))
			},
			wantStatus: fiber.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChecker := &mockPermissionChecker{}
			tt.mockSetup(mockChecker)

			app := fiber.New()
			app.Use(func(c *fiber.Ctx) error {
				if tt.userID > 0 {
					c.Locals(ssoclient.LocalsClaims, &ssoclient.Claims{UserID: tt.userID})
				}
				return c.Next()
			})
			app.Get("/test", auth.RequireAnyPermission(mockChecker, tt.permissionCodes...), func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockChecker.AssertExpectations(t)
		})
	}
}

// ============================================================
// TEST REQUIRE ALL PERMISSIONS
// ============================================================

func TestRequireAllPermissions(t *testing.T) {
	tests := []struct {
		name            string
		permissionCodes []string
		userID          int
		mockSetup       func(*mockPermissionChecker)
		wantStatus      int
	}{
		{
			name:            "success - user has all permissions",
			permissionCodes: []string{"admin.role.create", "admin.role.read"},
			userID:          1,
			mockSetup: func(m *mockPermissionChecker) {
				m.On("GetUserPermissions", mock.Anything, 1).Return([]string{"admin.role.create", "admin.role.read", "admin.role.delete"}, nil)
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name:            "success - empty permission list allows all",
			permissionCodes: []string{},
			userID:          1,
			mockSetup:       func(m *mockPermissionChecker) {},
			wantStatus:      fiber.StatusOK,
		},
		{
			name:            "forbidden - user missing one permission",
			permissionCodes: []string{"admin.role.create", "admin.role.delete"},
			userID:          1,
			mockSetup: func(m *mockPermissionChecker) {
				m.On("GetUserPermissions", mock.Anything, 1).Return([]string{"admin.role.create"}, nil)
			},
			wantStatus: fiber.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockChecker := &mockPermissionChecker{}
			tt.mockSetup(mockChecker)

			app := fiber.New()
			app.Use(func(c *fiber.Ctx) error {
				if tt.userID > 0 {
					c.Locals(ssoclient.LocalsClaims, &ssoclient.Claims{UserID: tt.userID})
				}
				return c.Next()
			})
			app.Get("/test", auth.RequireAllPermissions(mockChecker, tt.permissionCodes...), func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockChecker.AssertExpectations(t)
		})
	}
}

// ============================================================
// TEST PERMISSION CACHE
// ============================================================

func TestPermissionCache_HasPermission(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		permissionCode string
		mockSetup      func(*mockPermissionChecker)
		want           bool
		wantErr        bool
	}{
		{
			name:           "cache miss - fetches from service and caches",
			userID:         1,
			permissionCode: "admin.role.create",
			mockSetup: func(m *mockPermissionChecker) {
				m.On("GetUserPermissions", mock.Anything, 1).Return([]string{"admin.role.create", "admin.role.read"}, nil).Once()
			},
			want:    true,
			wantErr: false,
		},
		{
			name:           "permission not found",
			userID:         1,
			permissionCode: "admin.role.delete",
			mockSetup: func(m *mockPermissionChecker) {
				m.On("GetUserPermissions", mock.Anything, 1).Return([]string{"admin.role.create"}, nil).Once()
			},
			want:    false,
			wantErr: false,
		},
		{
			name:           "service error",
			userID:         1,
			permissionCode: "admin.role.create",
			mockSetup: func(m *mockPermissionChecker) {
				m.On("GetUserPermissions", mock.Anything, 1).Return(nil, errors.New("db error")).Once()
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mockPermissionChecker{}
			tt.mockSetup(mockService)

			cache := auth.NewPermissionCache(mockService, 5*time.Minute)

			got, err := cache.HasPermission(context.Background(), tt.userID, tt.permissionCode)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPermissionCache_CacheHit(t *testing.T) {
	mockService := &mockPermissionChecker{}
	// Service should only be called once (first request caches the result)
	mockService.On("GetUserPermissions", mock.Anything, 1).Return([]string{"admin.role.create"}, nil).Once()

	cache := auth.NewPermissionCache(mockService, 5*time.Minute)

	// First call - cache miss
	got1, err1 := cache.HasPermission(context.Background(), 1, "admin.role.create")
	assert.NoError(t, err1)
	assert.True(t, got1)

	// Second call - cache hit (service should NOT be called again)
	got2, err2 := cache.HasPermission(context.Background(), 1, "admin.role.create")
	assert.NoError(t, err2)
	assert.True(t, got2)

	// Verify service was only called once
	mockService.AssertNumberOfCalls(t, "GetUserPermissions", 1)
}

func TestPermissionCache_InvalidateUser(t *testing.T) {
	mockService := &mockPermissionChecker{}
	// Service should be called twice (once before invalidation, once after)
	mockService.On("GetUserPermissions", mock.Anything, 1).Return([]string{"admin.role.create"}, nil).Twice()

	cache := auth.NewPermissionCache(mockService, 5*time.Minute)

	// First call - cache miss
	_, _ = cache.HasPermission(context.Background(), 1, "admin.role.create")

	// Invalidate user cache
	cache.InvalidateUser(1)

	// Second call - should hit service again
	_, _ = cache.HasPermission(context.Background(), 1, "admin.role.create")

	// Verify service was called twice
	mockService.AssertNumberOfCalls(t, "GetUserPermissions", 2)
}

func TestPermissionCache_InvalidateAll(t *testing.T) {
	mockService := &mockPermissionChecker{}
	// Service should be called 4 times (2 users x 2 calls each)
	mockService.On("GetUserPermissions", mock.Anything, 1).Return([]string{"admin.role.create"}, nil).Twice()
	mockService.On("GetUserPermissions", mock.Anything, 2).Return([]string{"user.role.read"}, nil).Twice()

	cache := auth.NewPermissionCache(mockService, 5*time.Minute)

	// First calls - cache miss for both users
	_, _ = cache.HasPermission(context.Background(), 1, "admin.role.create")
	_, _ = cache.HasPermission(context.Background(), 2, "user.role.read")

	// Invalidate all
	cache.InvalidateAll()

	// Second calls - should hit service again
	_, _ = cache.HasPermission(context.Background(), 1, "admin.role.create")
	_, _ = cache.HasPermission(context.Background(), 2, "user.role.read")

	// Verify service was called appropriately
	mockService.AssertExpectations(t)
}

func TestPermissionCache_Stats(t *testing.T) {
	mockService := &mockPermissionChecker{}
	mockService.On("GetUserPermissions", mock.Anything, mock.Anything).Return([]string{"admin.role.create"}, nil)

	cache := auth.NewPermissionCache(mockService, 5*time.Minute)

	// Initially empty
	stats := cache.Stats()
	assert.Equal(t, 0, stats.CachedUsers)

	// Cache some users
	_, _ = cache.HasPermission(context.Background(), 1, "admin.role.create")
	_, _ = cache.HasPermission(context.Background(), 2, "admin.role.create")
	_, _ = cache.HasPermission(context.Background(), 3, "admin.role.create")

	stats = cache.Stats()
	assert.Equal(t, 3, stats.CachedUsers)
	assert.Equal(t, 5*time.Minute, stats.TTL)
}
