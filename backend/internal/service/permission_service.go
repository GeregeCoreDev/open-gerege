// Package service provides implementation for service
//
// File: permission_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"

	"templatev25/internal/auth"
	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"templatev25/internal/repository"

	"go.uber.org/zap"
)

type PermissionService struct {
	repo  repository.PermissionRepository
	log   *zap.Logger
	cache auth.CacheInvalidator // Permission cache invalidation (optional)
}

func NewPermissionService(repo repository.PermissionRepository, log *zap.Logger) *PermissionService {
	return &PermissionService{repo: repo, log: log}
}

// SetCacheInvalidator нь permission cache invalidator-ийг тохируулна.
// Permission-ууд өөрчлөгдөхөд cache цэвэрлэхэд ашиглагдана.
func (s *PermissionService) SetCacheInvalidator(cache auth.CacheInvalidator) {
	s.cache = cache
}

func (s *PermissionService) ListFilteredPaged(ctx context.Context, q dto.PermissionQuery) ([]domain.Permission, int64, int, int, error) {
	return s.repo.List(ctx, q)
}

func (s *PermissionService) ByID(ctx context.Context, id int) (domain.Permission, error) {
	return s.repo.ByID(ctx, id)
}

func (s *PermissionService) ByCode(ctx context.Context, code string) (domain.Permission, error) {
	return s.repo.ByCode(ctx, code)
}

func (s *PermissionService) Create(ctx context.Context, req dto.PermissionCreateDto) error {
	// ActionIDs-ээс Permission үүсгэх (нэг system-ийн нэг module-д олон action-д зориулсан permission үүсгэх)
	// Transaction ашиглаж бүх Permission-г нэгэн зэрэг үүсгэх
	return s.repo.CreateBatch(ctx, req.SystemID, req.ModuleID, req.ActionIDs)
}

func (s *PermissionService) Update(ctx context.Context, id int, req dto.PermissionUpdateDto) error {
	m := domain.Permission{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		ModuleID:    req.ModuleID,
		SystemID:    req.SystemID,
		ActionID:    req.ActionID,
	}
	if err := s.repo.Update(ctx, id, m); err != nil {
		return err
	}
	// Permission өөрчлөгдөхөд бүх cache цэвэрлэх
	if s.cache != nil {
		s.cache.InvalidateAll()
	}
	return nil
}

func (s *PermissionService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	// Permission устгахад бүх cache цэвэрлэх
	if s.cache != nil {
		s.cache.InvalidateAll()
	}
	return nil
}

// HasPermission нь хэрэглэгч тодорхой permission-тэй эсэхийг шалгана.
//
// Parameters:
//   - ctx: Context
//   - userID: Хэрэглэгчийн ID
//   - permissionCode: Permission код (жишээ: "admin.role.create")
//
// Returns:
//   - bool: Permission байвал true
//   - error: Алдаа
func (s *PermissionService) HasPermission(ctx context.Context, userID int, permissionCode string) (bool, error) {
	return s.repo.UserHasPermission(ctx, userID, permissionCode)
}

// GetUserPermissions нь хэрэглэгчийн бүх permission код-уудыг буцаана.
//
// Parameters:
//   - ctx: Context
//   - userID: Хэрэглэгчийн ID
//
// Returns:
//   - []string: Permission кодуудын жагсаалт
//   - error: Алдаа
func (s *PermissionService) GetUserPermissions(ctx context.Context, userID int) ([]string, error) {
	return s.repo.GetUserPermissionCodes(ctx, userID)
}
