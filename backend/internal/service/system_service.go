// Package service provides implementation for service
//
// File: system_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"
	"errors"
	"strings"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/repository"

	"go.uber.org/zap"
)

type SystemService interface {
	List(ctx context.Context, q dto.SystemListQuery) ([]domain.System, int64, int, int, error)
	ByID(ctx context.Context, id int) (domain.System, error)
	Create(ctx context.Context, req dto.SystemCreateDto) error
	Update(ctx context.Context, id int, req dto.SystemUpdateDto) error
	Delete(ctx context.Context, id int) error
}

type systemService struct {
	repo repository.SystemRepository
	log  *zap.Logger
}

func NewSystemService(repo repository.SystemRepository, log *zap.Logger) SystemService {
	return &systemService{repo: repo, log: log}
}

// List
func (s *systemService) List(ctx context.Context, q dto.SystemListQuery) ([]domain.System, int64, int, int, error) {
	items, total, page, size, err := s.repo.List(ctx, q)
	if err != nil {
		s.log.Error("system_list_failed", zap.Error(err))
		return nil, 0, 0, 0, err
	}
	s.log.Debug("system_list_success", zap.Int64("total", total), zap.Int("page", page))
	return items, total, page, size, nil
}

// ByID
func (s *systemService) ByID(ctx context.Context, id int) (domain.System, error) {
	sys, err := s.repo.ByID(ctx, id)
	if err != nil {
		s.log.Error("system_get_by_id_failed", zap.Int("system_id", id), zap.Error(err))
		return domain.System{}, err
	}
	return sys, nil
}

// Create
func (s *systemService) Create(ctx context.Context, req dto.SystemCreateDto) error {
	// Code-г lower case болгох
	code := strings.ToLower(req.Code)
	
	// Key хоосон бол code-ийн утгыг key-д оноох
	key := req.Key
	if key == "" {
		key = code
	} else {
		key = strings.ToLower(key)
	}

	m := domain.System{
		Code:        code,
		Key:         key,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
		Icon:        req.Icon,
		Sequence:    req.Sequence,
	}
	// CreatedUserId/CreatedOrgId нь repo талд ctx-ээс ононо
	if err := s.repo.Create(ctx, m); err != nil {
		s.log.Error("system_create_failed", zap.String("code", code), zap.Error(err))
		return err
	}
	s.log.Info("system_created", zap.String("code", code), zap.String("name", req.Name))
	return nil
}

// Update
func (s *systemService) Update(ctx context.Context, id int, req dto.SystemUpdateDto) error {
	if req.IsActive != nil && !*req.IsActive {
		if s.repo.GetActiveModuleCount(ctx, id) > 0 {
			s.log.Warn("system_update_blocked_has_modules", zap.Int("system_id", id))
			return errors.New("системд идэвхитэй модуль бүртгэлтэй тул идэвхигүй болгох боломжгүй")
		}

		if s.repo.GetActiveRoleCount(ctx, id) > 0 {
			s.log.Warn("system_update_blocked_has_roles", zap.Int("system_id", id))
			return errors.New("системд идэвхитэй эрх бүртгэлтэй тул идэвхигүй болгох боломжгүй")
		}
	}

	// Code-г lower case болгох
	code := strings.ToLower(req.Code)

	// Key хоосон бол code-ийн утгыг key-д оноох
	key := req.Key
	if key == "" {
		key = code
	} else {
		key = strings.ToLower(key)
	}

	m := domain.System{
		Code:        code,
		Key:         key,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
		Icon:        req.Icon,
		Sequence:    req.Sequence,
	}

	// UpdatedUserId/UpdatedOrgId нь repo талд ctx-ээс ононо
	if err := s.repo.Update(ctx, id, m); err != nil {
		s.log.Error("system_update_failed", zap.Int("system_id", id), zap.Error(err))
		return err
	}
	s.log.Info("system_updated", zap.Int("system_id", id))
	return nil
}

// Delete (soft delete)
func (s *systemService) Delete(ctx context.Context, id int) error {
	existing, err := s.repo.ByID(ctx, id)
	if err != nil {
		s.log.Error("system_delete_not_found", zap.Int("system_id", id), zap.Error(err))
		return err
	}
	if existing.IsActive != nil && *existing.IsActive {
		s.log.Warn("system_delete_blocked_active", zap.Int("system_id", id))
		return errors.New("систем идэвхитэй тул устгах боломжгүй")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("system_delete_failed", zap.Int("system_id", id), zap.Error(err))
		return err
	}
	s.log.Info("system_deleted", zap.Int("system_id", id))
	return nil
}
