// Package repository provides implementation for repository
//
// File: user_role_repo.go
// Description: implementation for repository
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package repository

import (
	"context"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRoleRepository interface {
	UsersByRole(ctx context.Context, q dto.UserRoleUsersQuery) ([]domain.UserRole, int64, int, int, error)
	RolesByUser(ctx context.Context, q dto.UserRoleRolesQuery) ([]domain.UserRole, int64, int, int, error)
	AddUsersToRole(ctx context.Context, roleID int, userIDs []int) error
	AddRolesToUser(ctx context.Context, userID int, roleIDs []int) error
	Remove(ctx context.Context, userID, roleID int) error
}

type userRoleRepository struct{ db *gorm.DB }

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository { return &userRoleRepository{db: db} }

// GET /user-role/users
func (r *userRoleRepository) UsersByRole(ctx context.Context, q dto.UserRoleUsersQuery) ([]domain.UserRole, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	tx := r.db.WithContext(ctx).Model(&domain.UserRole{}).Where("role_id = ?", q.RoleID)
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.UserRole
	if err := tx.Preload("User", func(px *gorm.DB) *gorm.DB {
		return px.Select("id", "civil_id", "reg_no", "family_name", "last_name", "first_name", "gender", "birth_date", "phone_no", "email")
	}).Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	return items, total, page, size, nil
}

// GET /user-role/roles
func (r *userRoleRepository) RolesByUser(ctx context.Context, q dto.UserRoleRolesQuery) ([]domain.UserRole, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	tx := r.db.WithContext(ctx).Model(&domain.UserRole{}).Where("user_id = ?", q.UserID)
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.UserRole
	if err := tx.Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = TRUE")
	}).
		Preload("Role.System").Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	return items, total, page, size, nil
}

// POST assign by role
// Batch insert with ON CONFLICT - N queries -> 1 query
func (r *userRoleRepository) AddUsersToRole(ctx context.Context, roleID int, userIDs []int) error {
	if len(userIDs) == 0 {
		return nil
	}

	// Build batch of UserRole records
	links := make([]domain.UserRole, 0, len(userIDs))
	for _, uid := range userIDs {
		links = append(links, domain.UserRole{RoleID: roleID, UserId: uid})
	}

	// Single batch insert with ON CONFLICT DO NOTHING (idempotent)
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "role_id"}},
		DoNothing: true,
	}).Create(&links).Error
}

// POST assign by user
// Batch insert with ON CONFLICT - N queries -> 1 query
func (r *userRoleRepository) AddRolesToUser(ctx context.Context, userID int, roleIDs []int) error {
	if len(roleIDs) == 0 {
		return nil
	}

	// Build batch of UserRole records
	links := make([]domain.UserRole, 0, len(roleIDs))
	for _, rid := range roleIDs {
		links = append(links, domain.UserRole{RoleID: rid, UserId: userID})
	}

	// Single batch insert with ON CONFLICT DO NOTHING (idempotent)
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "role_id"}},
		DoNothing: true,
	}).Create(&links).Error
}

func (r *userRoleRepository) Remove(ctx context.Context, userID, roleID int) error {
	return r.db.WithContext(ctx).Where("role_id = ? AND user_id = ?", roleID, userID).Delete(&domain.UserRole{}).Error
}
