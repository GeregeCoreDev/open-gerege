// Package service provides implementation for service
//
// File: role_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"
	"errors"

	"templatev25/internal/auth"
	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/middleware"
	"templatev25/internal/repository"

	"go.uber.org/zap"
)

type RoleService struct {
	repo  repository.RoleRepository
	log   *zap.Logger
	cache auth.CacheInvalidator // Permission cache invalidation (optional)
}

func NewRoleService(repo repository.RoleRepository, log *zap.Logger) *RoleService {
	return &RoleService{repo: repo, log: log}
}

// SetCacheInvalidator нь permission cache invalidator-ийг тохируулна.
// Role permission-ууд өөрчлөгдөхөд cache цэвэрлэхэд ашиглагдана.
func (s *RoleService) SetCacheInvalidator(cache auth.CacheInvalidator) {
	s.cache = cache
}

// ListFilteredPaged — model_service & menu_service загвартай адил
func (s *RoleService) ListFilteredPaged(ctx context.Context, p dto.RoleListQuery) ([]domain.Role, int64, int, int, error) {
	log := middleware.LoggerOrDefault(ctx, s.log)
	items, total, page, size, err := s.repo.List(ctx, p)
	if err != nil {
		log.Error("role_list_failed", zap.Error(err))
		return nil, 0, 0, 0, err
	}
	log.Debug("role_list_success", zap.Int64("total", total), zap.Int("page", page))
	return items, total, page, size, nil
}

// Create — handler аль хэдийн validate хийсэн гэж үзэж repo руу шууд дамжуулна
func (s *RoleService) Create(ctx context.Context, req dto.RoleCreateDto) error {
	log := middleware.LoggerOrDefault(ctx, s.log)
	m := domain.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		SystemID:    req.SystemID,
		IsActive:    req.IsActive,
	}
	if err := s.repo.Create(ctx, m); err != nil {
		log.Error("role_create_failed", zap.String("code", req.Code), zap.Error(err))
		return err
	}
	log.Info("role_created", zap.String("code", req.Code), zap.String("name", req.Name))
	return nil
}

// Update — model_service-тай ижил signature
func (s *RoleService) Update(ctx context.Context, id int, req dto.RoleUpdateDto) error {
	log := middleware.LoggerOrDefault(ctx, s.log)
	if req.IsActive != nil && !*req.IsActive && s.repo.GetUserCount(ctx, id) > 0 {
		log.Warn("role_update_blocked_has_users", zap.Int("role_id", id))
		return errors.New("эрх нь хэрэглэгчтэй бүртгэлтэй тул идэвхигүй болгох боломжгүй")
	}
	m := domain.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		SystemID:    req.SystemID,
		IsActive:    req.IsActive,
	}
	if err := s.repo.Update(ctx, id, m); err != nil {
		log.Error("role_update_failed", zap.Int("role_id", id), zap.Error(err))
		return err
	}
	log.Info("role_updated", zap.Int("role_id", id))
	return nil
}

// Delete — model_repo шиг soft-delete, repo буцаасан объектод Deleted* талбарууд populate-лагдана
func (s *RoleService) Delete(ctx context.Context, id int) error {
	log := middleware.LoggerOrDefault(ctx, s.log)
	existing, err := s.repo.ByID(ctx, id)
	if err != nil {
		log.Error("role_delete_not_found", zap.Int("role_id", id), zap.Error(err))
		return err
	}
	if existing.IsActive != nil && *existing.IsActive {
		log.Warn("role_delete_blocked_active", zap.Int("role_id", id))
		return errors.New("эрх идэвхитэй тул устгах боломжгүй")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		log.Error("role_delete_failed", zap.Int("role_id", id), zap.Error(err))
		return err
	}
	log.Info("role_deleted", zap.Int("role_id", id))
	return nil
}

func (s *RoleService) GetPermissions(ctx context.Context, q dto.RolePermissionsQuery) ([]domain.Permission, error) {
	log := middleware.LoggerOrDefault(ctx, s.log)
	perms, err := s.repo.Permissions(ctx, q)
	if err != nil {
		log.Error("role_permissions_get_failed", zap.Int("role_id", q.RoleID), zap.Error(err))
		return nil, err
	}
	log.Debug("role_permissions_fetched", zap.Int("role_id", q.RoleID), zap.Int("count", len(perms)))
	return perms, nil
}

func (s *RoleService) SetPermissions(ctx context.Context, req dto.RolePermissionsUpdateDto) error {
	log := middleware.LoggerOrDefault(ctx, s.log)
	if err := s.repo.ReplacePermissions(ctx, req.RoleID, req.PermissionIDs); err != nil {
		log.Error("role_permissions_set_failed", zap.Int("role_id", req.RoleID), zap.Error(err))
		return err
	}

	// Permission cache цэвэрлэх (role-д хамаарах бүх хэрэглэгчид)
	if s.cache != nil {
		s.cache.InvalidateAll() // Role permission өөрчлөгдөхөд бүх cache цэвэрлэх
		log.Debug("permission_cache_invalidated", zap.Int("role_id", req.RoleID))
	}

	log.Info("role_permissions_updated", zap.Int("role_id", req.RoleID), zap.Int("permission_count", len(req.PermissionIDs)))
	return nil
}
