// Package service provides implementation for service
//
// File: cached_organization_service.go
// Description: Cached wrapper for OrganizationService
//
// Caches organization data for frequent lookups.
// Automatically invalidates on organization modifications.
package service

import (
	"context"
	"fmt"
	"time"

	"templatev25/internal/cache"
	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/common"
	"go.uber.org/zap"
)

// CachedOrganizationService wraps OrganizationService with caching
type CachedOrganizationService struct {
	*OrganizationService
	orgCache  *cache.Cache[domain.Organization]
	treeCache *cache.Cache[[]domain.Organization]
}

// NewCachedOrganizationService creates a new cached organization service
func NewCachedOrganizationService(svc *OrganizationService) *CachedOrganizationService {
	orgCfg := cache.Config{
		MaxSize:         2000,             // Cache up to 2000 organizations
		TTL:             15 * time.Minute, // 15 minute TTL
		CleanupInterval: 5 * time.Minute,
	}

	treeCfg := cache.Config{
		MaxSize:         100,              // Cache up to 100 tree queries
		TTL:             10 * time.Minute, // 10 minute TTL
		CleanupInterval: 5 * time.Minute,
	}

	return &CachedOrganizationService{
		OrganizationService: svc,
		orgCache:            cache.New[domain.Organization](orgCfg),
		treeCache:           cache.New[[]domain.Organization](treeCfg),
	}
}

// orgKey generates a cache key for an organization ID
func (s *CachedOrganizationService) orgKey(id int) string {
	return fmt.Sprintf("org:%d", id)
}

// treeKey generates a cache key for organization tree
func (s *CachedOrganizationService) treeKey(rootID int) string {
	return fmt.Sprintf("org:tree:%d", rootID)
}

// ByID retrieves an organization by ID with caching
func (s *CachedOrganizationService) ByID(ctx context.Context, id int) (domain.Organization, error) {
	key := s.orgKey(id)

	// Try cache first
	if org, found := s.orgCache.Get(key); found {
		return org, nil
	}

	// Cache miss - fetch from database
	org, err := s.OrganizationService.ByID(ctx, id)
	if err != nil {
		return domain.Organization{}, err
	}

	// Store in cache
	s.orgCache.Set(key, org)
	return org, nil
}

// Tree retrieves organization tree with caching
func (s *CachedOrganizationService) Tree(ctx context.Context, rootID int) ([]domain.Organization, error) {
	key := s.treeKey(rootID)

	// Try cache first
	if tree, found := s.treeCache.Get(key); found {
		return tree, nil
	}

	// Cache miss - fetch from database
	tree, err := s.OrganizationService.Tree(ctx, rootID)
	if err != nil {
		return nil, err
	}

	// Store in cache
	s.treeCache.Set(key, tree)
	return tree, nil
}

// List is not cached (pagination varies)
func (s *CachedOrganizationService) List(ctx context.Context, p common.PaginationQuery) ([]domain.Organization, int64, int, int, error) {
	return s.OrganizationService.List(ctx, p)
}

// Create invalidates cache after creating
func (s *CachedOrganizationService) Create(ctx context.Context, req dto.OrganizationDto) (domain.Organization, error) {
	org, err := s.OrganizationService.Create(ctx, req)
	if err != nil {
		return domain.Organization{}, err
	}

	// Cache the new organization
	s.orgCache.Set(s.orgKey(org.Id), org)

	// Invalidate tree caches (structure changed)
	s.treeCache.Clear()

	return org, nil
}

// Update invalidates cache after updating
func (s *CachedOrganizationService) Update(ctx context.Context, id int, req dto.OrganizationUpdateDto) (domain.Organization, error) {
	org, err := s.OrganizationService.Update(ctx, id, req)
	if err != nil {
		return domain.Organization{}, err
	}

	// Update cache
	s.orgCache.Set(s.orgKey(id), org)

	// Invalidate tree caches if parent changed
	if req.ParentID != nil {
		s.treeCache.Clear()
	}

	return org, nil
}

// Delete invalidates cache after deleting
func (s *CachedOrganizationService) Delete(ctx context.Context, id int) error {
	err := s.OrganizationService.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Remove from cache
	s.orgCache.Delete(s.orgKey(id))

	// Invalidate tree caches (structure changed)
	s.treeCache.Clear()

	return nil
}

// CacheStats returns cache statistics
func (s *CachedOrganizationService) CacheStats() (orgStats, treeStats cache.Stats) {
	return s.orgCache.Stats(), s.treeCache.Stats()
}

// InvalidateOrg removes a specific organization from cache
func (s *CachedOrganizationService) InvalidateOrg(id int) {
	s.orgCache.Delete(s.orgKey(id))
}

// InvalidateAll clears all organization caches
func (s *CachedOrganizationService) InvalidateAll() {
	s.orgCache.Clear()
	s.treeCache.Clear()
}

// Stop stops the cache cleanup goroutines
func (s *CachedOrganizationService) Stop() {
	s.orgCache.Stop()
	s.treeCache.Stop()
}

// LogCacheStats logs current cache statistics
func (s *CachedOrganizationService) LogCacheStats(log *zap.Logger) {
	orgStats, treeStats := s.CacheStats()
	log.Info("organization_cache_stats",
		zap.Int("org_size", orgStats.Size),
		zap.Int64("org_hits", orgStats.Hits),
		zap.Int64("org_misses", orgStats.Misses),
		zap.Float64("org_hit_ratio", orgStats.Ratio),
		zap.Int("tree_size", treeStats.Size),
		zap.Int64("tree_hits", treeStats.Hits),
		zap.Int64("tree_misses", treeStats.Misses),
		zap.Float64("tree_hit_ratio", treeStats.Ratio),
	)
}
