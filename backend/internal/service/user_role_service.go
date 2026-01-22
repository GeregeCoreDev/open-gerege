// Package service provides implementation for service
//
// File: user_role_service.go
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
)

type UserRoleService interface {
	UsersByRole(ctx context.Context, q dto.UserRoleUsersQuery) ([]domain.UserRole, int64, int, int, error)
	RolesByUser(ctx context.Context, q dto.UserRoleRolesQuery) ([]domain.UserRole, int64, int, int, error)
	AssignByRole(ctx context.Context, req dto.UserRoleAssignByRole) error
	AssignByUser(ctx context.Context, req dto.UserRoleAssignByUser) error
	Remove(ctx context.Context, req dto.UserRoleRemoveDto) error
	SetCacheInvalidator(cache auth.CacheInvalidator)
}

type userRoleService struct {
	repo  repository.UserRoleRepository
	cache auth.CacheInvalidator // Permission cache invalidation (optional)
}

func NewUserRoleService(repo repository.UserRoleRepository) UserRoleService {
	return &userRoleService{repo: repo}
}

// SetCacheInvalidator нь permission cache invalidator-ийг тохируулна.
// User role-ууд өөрчлөгдөхөд cache цэвэрлэхэд ашиглагдана.
func (s *userRoleService) SetCacheInvalidator(cache auth.CacheInvalidator) {
	s.cache = cache
}

func (s *userRoleService) UsersByRole(ctx context.Context, q dto.UserRoleUsersQuery) ([]domain.UserRole, int64, int, int, error) {
	return s.repo.UsersByRole(ctx, q)
}
func (s *userRoleService) RolesByUser(ctx context.Context, q dto.UserRoleRolesQuery) ([]domain.UserRole, int64, int, int, error) {
	return s.repo.RolesByUser(ctx, q)
}
func (s *userRoleService) AssignByRole(ctx context.Context, req dto.UserRoleAssignByRole) error {
	if err := s.repo.AddUsersToRole(ctx, req.RoleID, req.UserIDs); err != nil {
		return err
	}
	// Cache цэвэрлэх (role-д нэмэгдсэн хэрэглэгчид)
	if s.cache != nil {
		s.cache.InvalidateUsers(req.UserIDs)
	}
	return nil
}
func (s *userRoleService) AssignByUser(ctx context.Context, req dto.UserRoleAssignByUser) error {
	if err := s.repo.AddRolesToUser(ctx, req.UserID, req.RoleIDs); err != nil {
		return err
	}
	// Cache цэвэрлэх (хэрэглэгчийн role өөрчлөгдсөн)
	if s.cache != nil {
		s.cache.InvalidateUser(req.UserID)
	}
	return nil
}
func (s *userRoleService) Remove(ctx context.Context, req dto.UserRoleRemoveDto) error {
	if err := s.repo.Remove(ctx, req.UserID, req.RoleID); err != nil {
		return err
	}
	// Cache цэвэрлэх (хэрэглэгчийн role устсан)
	if s.cache != nil {
		s.cache.InvalidateUser(req.UserID)
	}
	return nil
}
