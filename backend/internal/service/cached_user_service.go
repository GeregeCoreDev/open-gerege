// Package service provides implementation for service
//
// File: cached_user_service.go
// Description: Cached wrapper for UserService
//
// This wrapper adds caching to frequently accessed user data.
// Cache is automatically invalidated on write operations.
package service

import (
	"context"
	"fmt"
	"time"

	"templatev25/internal/cache"
	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/config"
	"go.uber.org/zap"
)

// CachedUserService wraps UserService with caching
type CachedUserService struct {
	*UserService
	cache *cache.Cache[domain.User]
}

// NewCachedUserService creates a new cached user service
func NewCachedUserService(svc *UserService, cfg *config.Config) *CachedUserService {
	// Configure cache based on environment
	cacheCfg := cache.Config{
		MaxSize:         5000,              // Cache up to 5000 users
		TTL:             10 * time.Minute,  // 10 minute TTL
		CleanupInterval: 5 * time.Minute,
	}

	return &CachedUserService{
		UserService: svc,
		cache:       cache.New[domain.User](cacheCfg),
	}
}

// cacheKey generates a cache key for a user ID
func (s *CachedUserService) cacheKey(id int) string {
	return fmt.Sprintf("user:%d", id)
}

// GetByID retrieves a user by ID with caching
func (s *CachedUserService) GetByID(ctx context.Context, id int) (domain.User, error) {
	key := s.cacheKey(id)

	// Try cache first
	if user, found := s.cache.Get(key); found {
		return user, nil
	}

	// Cache miss - fetch from database
	user, err := s.UserService.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	// Store in cache
	s.cache.Set(key, user)
	return user, nil
}

// List is not cached (pagination varies)
func (s *CachedUserService) List(ctx context.Context, p common.PaginationQuery) ([]domain.User, int64, int, int, error) {
	return s.UserService.List(ctx, p)
}

// Create invalidates cache after creating
func (s *CachedUserService) Create(ctx context.Context, req dto.UserCreateDto) (domain.User, error) {
	user, err := s.UserService.Create(ctx, req)
	if err != nil {
		return domain.User{}, err
	}

	// Cache the new user
	s.cache.Set(s.cacheKey(user.Id), user)
	return user, nil
}

// Update invalidates cache after updating
func (s *CachedUserService) Update(ctx context.Context, req dto.UserUpdateDto) (domain.User, error) {
	user, err := s.UserService.Update(ctx, req)
	if err != nil {
		return domain.User{}, err
	}

	// Update cache
	s.cache.Set(s.cacheKey(user.Id), user)
	return user, nil
}

// Delete invalidates cache after deleting
func (s *CachedUserService) Delete(ctx context.Context, id int) (domain.User, error) {
	user, err := s.UserService.Delete(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	// Remove from cache
	s.cache.Delete(s.cacheKey(id))
	return user, nil
}

// Stats returns cache statistics
func (s *CachedUserService) CacheStats() cache.Stats {
	return s.cache.Stats()
}

// InvalidateUser removes a specific user from cache
func (s *CachedUserService) InvalidateUser(id int) {
	s.cache.Delete(s.cacheKey(id))
}

// InvalidateAll clears the entire user cache
func (s *CachedUserService) InvalidateAll() {
	s.cache.Clear()
}

// Stop stops the cache cleanup goroutine
func (s *CachedUserService) Stop() {
	s.cache.Stop()
}

// LogCacheStats logs current cache statistics
func (s *CachedUserService) LogCacheStats(log *zap.Logger) {
	stats := s.cache.Stats()
	log.Info("user_cache_stats",
		zap.Int("size", stats.Size),
		zap.Int64("hits", stats.Hits),
		zap.Int64("misses", stats.Misses),
		zap.Float64("hit_ratio", stats.Ratio),
	)
}
