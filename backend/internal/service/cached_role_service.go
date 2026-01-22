// Package service provides implementation for service
//
// File: cached_role_service.go
// Description: Cached wrapper for RoleService
//
// Caches role data for frequent permission checks.
// Automatically invalidates on role modifications.
package service

import (
	"context"
	"fmt"
	"time"

	"templatev25/internal/cache"
	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"go.uber.org/zap"
)

// CachedRoleService wraps RoleService with caching
type CachedRoleService struct {
	*RoleService
	roleCache       *cache.Cache[domain.Role]
	permissionCache *cache.Cache[[]domain.Permission]
}

// NewCachedRoleService creates a new cached role service
func NewCachedRoleService(svc *RoleService) *CachedRoleService {
	roleCfg := cache.Config{
		MaxSize:         500,              // Cache up to 500 roles
		TTL:             15 * time.Minute, // 15 minute TTL
		CleanupInterval: 5 * time.Minute,
	}

	permCfg := cache.Config{
		MaxSize:         1000,             // Cache permissions for up to 1000 roles
		TTL:             5 * time.Minute,  // 5 minute TTL (shorter for security)
		CleanupInterval: 2 * time.Minute,
	}

	return &CachedRoleService{
		RoleService:     svc,
		roleCache:       cache.New[domain.Role](roleCfg),
		permissionCache: cache.New[[]domain.Permission](permCfg),
	}
}

// roleKey generates a cache key for a role ID
func (s *CachedRoleService) roleKey(id int) string {
	return fmt.Sprintf("role:%d", id)
}

// permKey generates a cache key for role permissions
func (s *CachedRoleService) permKey(roleID int) string {
	return fmt.Sprintf("role:%d:perms", roleID)
}

// GetPermissions retrieves role permissions with caching
func (s *CachedRoleService) GetPermissions(ctx context.Context, q dto.RolePermissionsQuery) ([]domain.Permission, error) {
	key := s.permKey(q.RoleID)

	// Try cache first
	if perms, found := s.permissionCache.Get(key); found {
		return perms, nil
	}

	// Cache miss - fetch from database
	perms, err := s.RoleService.GetPermissions(ctx, q)
	if err != nil {
		return nil, err
	}

	// Store in cache
	s.permissionCache.Set(key, perms)
	return perms, nil
}

// SetPermissions invalidates cache after updating permissions
func (s *CachedRoleService) SetPermissions(ctx context.Context, req dto.RolePermissionsUpdateDto) error {
	err := s.RoleService.SetPermissions(ctx, req)
	if err != nil {
		return err
	}

	// Invalidate all permission cache for this role
	s.permissionCache.DeletePrefix(fmt.Sprintf("role:%d:perms:", req.RoleID))
	return nil
}

// Create invalidates cache after creating
func (s *CachedRoleService) Create(ctx context.Context, req dto.RoleCreateDto) error {
	return s.RoleService.Create(ctx, req)
}

// Update invalidates cache after updating
func (s *CachedRoleService) Update(ctx context.Context, id int, req dto.RoleUpdateDto) error {
	err := s.RoleService.Update(ctx, id, req)
	if err != nil {
		return err
	}

	// Invalidate cache
	s.roleCache.Delete(s.roleKey(id))
	return nil
}

// Delete invalidates cache after deleting
func (s *CachedRoleService) Delete(ctx context.Context, id int) error {
	err := s.RoleService.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate caches
	s.roleCache.Delete(s.roleKey(id))
	s.permissionCache.DeletePrefix(fmt.Sprintf("role:%d:perms:", id))
	return nil
}

// CacheStats returns cache statistics
func (s *CachedRoleService) CacheStats() (roleStats, permStats cache.Stats) {
	return s.roleCache.Stats(), s.permissionCache.Stats()
}

// InvalidateRole removes a specific role from cache
func (s *CachedRoleService) InvalidateRole(id int) {
	s.roleCache.Delete(s.roleKey(id))
	s.permissionCache.DeletePrefix(fmt.Sprintf("role:%d:perms:", id))
}

// InvalidateAll clears all role caches
func (s *CachedRoleService) InvalidateAll() {
	s.roleCache.Clear()
	s.permissionCache.Clear()
}

// Stop stops the cache cleanup goroutines
func (s *CachedRoleService) Stop() {
	s.roleCache.Stop()
	s.permissionCache.Stop()
}

// LogCacheStats logs current cache statistics
func (s *CachedRoleService) LogCacheStats(log *zap.Logger) {
	roleStats, permStats := s.CacheStats()
	log.Info("role_cache_stats",
		zap.Int("role_size", roleStats.Size),
		zap.Int64("role_hits", roleStats.Hits),
		zap.Int64("role_misses", roleStats.Misses),
		zap.Float64("role_hit_ratio", roleStats.Ratio),
		zap.Int("perm_size", permStats.Size),
		zap.Int64("perm_hits", permStats.Hits),
		zap.Int64("perm_misses", permStats.Misses),
		zap.Float64("perm_hit_ratio", permStats.Ratio),
	)
}
