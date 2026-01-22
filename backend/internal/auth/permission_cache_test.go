// Package auth provides authentication and authorization utilities
//
// File: permission_cache_test.go
// Description: Unit tests for permission cache
package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPermissionChecker is a mock implementation of PermissionChecker
type mockPermissionChecker struct {
	permissions map[int][]string
	callCount   int
}

func (m *mockPermissionChecker) HasPermission(ctx context.Context, userID int, permissionCode string) (bool, error) {
	if perms, ok := m.permissions[userID]; ok {
		for _, p := range perms {
			if p == permissionCode {
				return true, nil
			}
		}
	}
	return false, nil
}

func (m *mockPermissionChecker) GetUserPermissions(ctx context.Context, userID int) ([]string, error) {
	m.callCount++
	if perms, ok := m.permissions[userID]; ok {
		return perms, nil
	}
	return []string{}, nil
}

func newMockChecker(permissions map[int][]string) *mockPermissionChecker {
	return &mockPermissionChecker{permissions: permissions}
}

func TestNewPermissionCache(t *testing.T) {
	mock := newMockChecker(nil)
	ttl := 5 * time.Minute

	cache := NewPermissionCache(mock, ttl)

	assert.NotNil(t, cache)
	assert.Equal(t, ttl, cache.ttl)
}

func TestPermissionCache_HasPermission(t *testing.T) {
	mock := newMockChecker(map[int][]string{
		1: {"admin.user.read", "admin.user.write"},
		2: {"admin.role.read"},
	})
	cache := NewPermissionCache(mock, 5*time.Minute)
	ctx := context.Background()

	tests := []struct {
		name           string
		userID         int
		permissionCode string
		want           bool
	}{
		{
			name:           "user has permission",
			userID:         1,
			permissionCode: "admin.user.read",
			want:           true,
		},
		{
			name:           "user has another permission",
			userID:         1,
			permissionCode: "admin.user.write",
			want:           true,
		},
		{
			name:           "user does not have permission",
			userID:         1,
			permissionCode: "admin.role.delete",
			want:           false,
		},
		{
			name:           "different user with permission",
			userID:         2,
			permissionCode: "admin.role.read",
			want:           true,
		},
		{
			name:           "user without any permissions",
			userID:         999,
			permissionCode: "admin.user.read",
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cache.HasPermission(ctx, tt.userID, tt.permissionCode)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPermissionCache_GetUserPermissions(t *testing.T) {
	mock := newMockChecker(map[int][]string{
		1: {"perm1", "perm2", "perm3"},
	})
	cache := NewPermissionCache(mock, 5*time.Minute)
	ctx := context.Background()

	// First call - should hit the mock
	perms, err := cache.GetUserPermissions(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, []string{"perm1", "perm2", "perm3"}, perms)
	assert.Equal(t, 1, mock.callCount)

	// Second call - should use cache
	perms, err = cache.GetUserPermissions(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, []string{"perm1", "perm2", "perm3"}, perms)
	assert.Equal(t, 1, mock.callCount) // Still 1, cache was used
}

func TestPermissionCache_InvalidateUser(t *testing.T) {
	mock := newMockChecker(map[int][]string{
		1: {"perm1"},
	})
	cache := NewPermissionCache(mock, 5*time.Minute)
	ctx := context.Background()

	// First call - populate cache
	_, _ = cache.GetUserPermissions(ctx, 1)
	assert.Equal(t, 1, mock.callCount)

	// Invalidate user
	cache.InvalidateUser(1)

	// Next call should hit mock again
	_, _ = cache.GetUserPermissions(ctx, 1)
	assert.Equal(t, 2, mock.callCount)
}

func TestPermissionCache_InvalidateUsers(t *testing.T) {
	mock := newMockChecker(map[int][]string{
		1: {"perm1"},
		2: {"perm2"},
		3: {"perm3"},
	})
	cache := NewPermissionCache(mock, 5*time.Minute)
	ctx := context.Background()

	// Populate cache for all users
	_, _ = cache.GetUserPermissions(ctx, 1)
	_, _ = cache.GetUserPermissions(ctx, 2)
	_, _ = cache.GetUserPermissions(ctx, 3)
	assert.Equal(t, 3, mock.callCount)

	// Invalidate users 1 and 2
	cache.InvalidateUsers([]int{1, 2})

	// User 3 should still use cache
	_, _ = cache.GetUserPermissions(ctx, 3)
	assert.Equal(t, 3, mock.callCount) // No change

	// Users 1 and 2 should hit mock again
	_, _ = cache.GetUserPermissions(ctx, 1)
	_, _ = cache.GetUserPermissions(ctx, 2)
	assert.Equal(t, 5, mock.callCount) // +2
}

func TestPermissionCache_InvalidateAll(t *testing.T) {
	mock := newMockChecker(map[int][]string{
		1: {"perm1"},
		2: {"perm2"},
	})
	cache := NewPermissionCache(mock, 5*time.Minute)
	ctx := context.Background()

	// Populate cache
	_, _ = cache.GetUserPermissions(ctx, 1)
	_, _ = cache.GetUserPermissions(ctx, 2)
	assert.Equal(t, 2, mock.callCount)

	// Invalidate all
	cache.InvalidateAll()

	// All calls should hit mock
	_, _ = cache.GetUserPermissions(ctx, 1)
	_, _ = cache.GetUserPermissions(ctx, 2)
	assert.Equal(t, 4, mock.callCount) // +2
}

func TestPermissionCache_Stats(t *testing.T) {
	mock := newMockChecker(map[int][]string{
		1: {"perm1"},
		2: {"perm2"},
	})
	ttl := 5 * time.Minute
	cache := NewPermissionCache(mock, ttl)
	ctx := context.Background()

	// Initially empty
	stats := cache.Stats()
	assert.Equal(t, 0, stats.CachedUsers)
	assert.Equal(t, ttl, stats.TTL)

	// Add some entries
	_, _ = cache.GetUserPermissions(ctx, 1)
	_, _ = cache.GetUserPermissions(ctx, 2)

	stats = cache.Stats()
	assert.Equal(t, 2, stats.CachedUsers)
}

func TestCachedPermissions_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "not expired",
			expiresAt: time.Now().Add(5 * time.Minute),
			want:      false,
		},
		{
			name:      "expired",
			expiresAt: time.Now().Add(-1 * time.Minute),
			want:      true,
		},
		{
			name:      "just expired",
			expiresAt: time.Now().Add(-1 * time.Millisecond),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := &cachedPermissions{
				codes:     []string{"test"},
				expiresAt: tt.expiresAt,
			}
			assert.Equal(t, tt.want, cp.isExpired())
		})
	}
}

func TestPermissionCache_TTLExpiration(t *testing.T) {
	mock := newMockChecker(map[int][]string{
		1: {"perm1"},
	})
	// Very short TTL for testing
	cache := NewPermissionCache(mock, 50*time.Millisecond)
	ctx := context.Background()

	// First call
	_, _ = cache.GetUserPermissions(ctx, 1)
	assert.Equal(t, 1, mock.callCount)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should hit mock again due to expiration
	_, _ = cache.GetUserPermissions(ctx, 1)
	assert.Equal(t, 2, mock.callCount)
}
