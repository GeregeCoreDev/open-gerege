// Package auth provides implementation for auth
//
// File: permission_cache.go
// Description: Permission caching layer for performance optimization
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
/*
Package auth нь permission caching-ийг хариуцна.

Энэ файл нь permission-уудыг cache-лэх layer-ийг тодорхойлно.
Cache нь DB руу очих request-ийн тоог бууруулж, хурдыг нэмэгдүүлнэ.

Cache бүтэц:
  - In-memory cache (sync.Map ашиглана)
  - TTL-тэй (default 5 минут)
  - User ID-гаар key хадгална

Invalidation:
  - InvalidateUser: Хэрэглэгчийн cache-ийг цэвэрлэх
  - InvalidateRole: Role-д хамаарах бүх хэрэглэгчийн cache цэвэрлэх
  - InvalidateAll: Бүх cache цэвэрлэх

Ашиглалт:

	// Cache үүсгэх
	permCache := auth.NewPermissionCache(permService, 5*time.Minute)

	// Middleware-д ашиглах
	app.Post("/role",
	    auth.RequirePermission(permCache, "admin.role.create"),
	    handler.Create,
	)

	// Cache цэвэрлэх (role өөрчлөгдөхөд)
	permCache.InvalidateRole(roleID)
*/
package auth

import (
	"context"
	"slices"
	"sync"
	"time"
)

// ============================================================
// CACHED PERMISSIONS STRUCTURE
// ============================================================

// cachedPermissions нь хэрэглэгчийн permission-уудыг TTL-тэй хадгална.
type cachedPermissions struct {
	codes     []string  // Permission кодуудын жагсаалт
	expiresAt time.Time // Cache хүчинтэй хугацаа
}

// isExpired нь cache хүчингүй болсон эсэхийг шалгана.
func (cp *cachedPermissions) isExpired() bool {
	return time.Now().After(cp.expiresAt)
}

// ============================================================
// CACHE INVALIDATOR INTERFACE
// ============================================================

// CacheInvalidator нь cache цэвэрлэх интерфейс.
// Service layer энэ интерфейсийг ашиглаж cache invalidation хийнэ.
type CacheInvalidator interface {
	// InvalidateUser нь нэг хэрэглэгчийн cache цэвэрлэнэ
	InvalidateUser(userID int)
	// InvalidateUsers нь олон хэрэглэгчийн cache цэвэрлэнэ
	InvalidateUsers(userIDs []int)
	// InvalidateAll нь бүх cache цэвэрлэнэ
	InvalidateAll()
}

// ============================================================
// PERMISSION CACHE
// ============================================================

// PermissionCache нь permission-уудыг cache-лэх layer.
// PermissionChecker интерфейсийг implement хийнэ.
type PermissionCache struct {
	service PermissionChecker // Underlying service (DB руу хандах)
	cache   sync.Map          // userID -> *cachedPermissions
	ttl     time.Duration     // Cache TTL
	mu      sync.RWMutex      // Role invalidation-д ашиглах
}

// NewPermissionCache нь шинэ permission cache үүсгэнэ.
//
// Parameters:
//   - service: Underlying permission service
//   - ttl: Cache-ийн хүчинтэй хугацаа (жишээ: 5*time.Minute)
//
// Returns:
//   - *PermissionCache: Cache instance
func NewPermissionCache(service PermissionChecker, ttl time.Duration) *PermissionCache {
	return &PermissionCache{
		service: service,
		ttl:     ttl,
	}
}

// ============================================================
// PERMISSION CHECKER IMPLEMENTATION
// ============================================================

// HasPermission нь хэрэглэгч тодорхой permission-тэй эсэхийг шалгана.
// Cache-д байвал DB руу явахгүй.
//
// Parameters:
//   - ctx: Context
//   - userID: Хэрэглэгчийн ID
//   - permissionCode: Permission код
//
// Returns:
//   - bool: Permission байвал true
//   - error: Алдаа
func (pc *PermissionCache) HasPermission(ctx context.Context, userID int, permissionCode string) (bool, error) {
	// Бүх permission-уудыг авах (cache ашиглана)
	perms, err := pc.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	// Permission байгаа эсэхийг шалгах
	return slices.Contains(perms, permissionCode), nil
}

// GetUserPermissions нь хэрэглэгчийн бүх permission-уудыг буцаана.
// Cache-д байвал DB руу явахгүй.
//
// Parameters:
//   - ctx: Context
//   - userID: Хэрэглэгчийн ID
//
// Returns:
//   - []string: Permission кодуудын жагсаалт
//   - error: Алдаа
func (pc *PermissionCache) GetUserPermissions(ctx context.Context, userID int) ([]string, error) {
	// ============================================================
	// STEP 1: Cache-ээс хайх
	// ============================================================
	if cached, ok := pc.cache.Load(userID); ok {
		cp := cached.(*cachedPermissions)
		if !cp.isExpired() {
			return cp.codes, nil
		}
		// Хүчингүй болсон бол устгах
		pc.cache.Delete(userID)
	}

	// ============================================================
	// STEP 2: DB-ээс авах
	// ============================================================
	perms, err := pc.service.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// ============================================================
	// STEP 3: Cache-д хадгалах
	// ============================================================
	pc.cache.Store(userID, &cachedPermissions{
		codes:     perms,
		expiresAt: time.Now().Add(pc.ttl),
	})

	return perms, nil
}

// ============================================================
// CACHE INVALIDATION
// ============================================================

// InvalidateUser нь тодорхой хэрэглэгчийн cache-ийг цэвэрлэнэ.
// Хэрэглэгчийн role-ууд өөрчлөгдөхөд дуудна.
//
// Parameters:
//   - userID: Хэрэглэгчийн ID
func (pc *PermissionCache) InvalidateUser(userID int) {
	pc.cache.Delete(userID)
}

// InvalidateUsers нь олон хэрэглэгчийн cache-ийг цэвэрлэнэ.
//
// Parameters:
//   - userIDs: Хэрэглэгчдийн ID-ууд
func (pc *PermissionCache) InvalidateUsers(userIDs []int) {
	for _, id := range userIDs {
		pc.cache.Delete(id)
	}
}

// InvalidateAll нь бүх cache-ийг цэвэрлэнэ.
// Permission эсвэл role_permission өөрчлөгдөхөд дуудна.
func (pc *PermissionCache) InvalidateAll() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.cache = sync.Map{}
}

// ============================================================
// CACHE STATS (DEBUG/MONITORING)
// ============================================================

// CacheStats нь cache-ийн статистик мэдээлэл.
type CacheStats struct {
	CachedUsers int           // Cache-д байгаа хэрэглэгчийн тоо
	TTL         time.Duration // Cache TTL
}

// Stats нь cache-ийн статистикийг буцаана.
func (pc *PermissionCache) Stats() CacheStats {
	count := 0
	pc.cache.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return CacheStats{
		CachedUsers: count,
		TTL:         pc.ttl,
	}
}
