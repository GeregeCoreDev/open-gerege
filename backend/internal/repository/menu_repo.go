// Package repository provides implementation for repository
//
// File: menu_repo.go
// Description: implementation for repository
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package repository

import (
	"context"
	"time"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/config"
	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MenuRepository interface {
	List(ctx context.Context, q dto.MenuListQuery) ([]domain.Menu, int64, int, int, error)
	ListAll(ctx context.Context) ([]domain.Menu, error)
	ListByUserRoles(ctx context.Context, userID int) ([]domain.Menu, error)
	GetMenusByPermissionIDs(ctx context.Context, permissionIDs []int) ([]domain.Menu, error)
	GetMenusByIDs(ctx context.Context, ids []int64) ([]domain.Menu, error)
	ByID(ctx context.Context, id int64) (domain.Menu, error)
	Create(ctx context.Context, m domain.Menu) error
	Update(ctx context.Context, id int64, m domain.Menu) error
	Delete(ctx context.Context, id int64) error
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB, cfg *config.Config) MenuRepository {
	// cfg parameter kept for backward compatibility but no longer used
	// search_path is now set in DSN, so schema name is not needed
	return &menuRepository{db: db}
}

func (r *menuRepository) List(ctx context.Context, q dto.MenuListQuery) ([]domain.Menu, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	colMap := scopes.ColumnMap{
		"id":          "menus.id",
		"code":        "menus.code",
		"key":         "menus.key",
		"name":        "menus.name",
		"description": "menus.description",
		"icon":        "menus.icon",
		"path":        "menus.path",
		"sequence":    "menus.sequence",
		"parent_id":   "menus.parent_id",
	}

	tx := r.db.WithContext(ctx).Model(&domain.Menu{}).Scopes(
		scopes.SearchScope(colMap, utils.ParseSearch(q.Search)),
		scopes.DateScope(q.CreatedFrom, q.CreatedTo),
	)

	if q.Code != "" {
		tx = tx.Where("code ILIKE ?", "%"+q.Code+"%")
	}
	if q.Key != "" {
		tx = tx.Where("key ILIKE ?", "%"+q.Key+"%")
	}
	if q.Name != "" {
		tx = tx.Where("name ILIKE ?", "%"+q.Name+"%")
	}
	if q.ParentID != nil {
		tx = tx.Where("parent_id = ?", *q.ParentID)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.Menu
	if err := tx.Scopes(
		scopes.SortScope(colMap, utils.ParseSort(q.Sort), "sequence ASC, id DESC"),
	).Preload("Parent").Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	return items, total, page, size, nil
}

func preload(d *gorm.DB) *gorm.DB {
	return d.Preload("Children", func(db *gorm.DB) *gorm.DB {
		return db.Order("sequence ASC, id ASC").Preload("Children", preload)
	})
}

func (r *menuRepository) ListAll(ctx context.Context) ([]domain.Menu, error) {
	var allMenus []domain.Menu
	tx := r.db.WithContext(ctx).
		Order("sequence ASC, id ASC")

	if err := tx.Preload(clause.Associations, preload).Find(&allMenus, "parent_id IS NULL").Error; err != nil {
		return nil, err
	}

	return allMenus, nil
}

// ListByUserRoles returns menus accessible to a user based on their roles.
// Optimized: Single JOIN query instead of 3 sequential queries.
func (r *menuRepository) ListByUserRoles(ctx context.Context, userID int) ([]domain.Menu, error) {
	var menus []domain.Menu

	// Single query with JOINs: menus <- role_permissions <- user_roles
	// This replaces 3 sequential queries with 1 efficient JOIN
	err := r.db.WithContext(ctx).
		Model(&domain.Menu{}).
		Distinct().
		Joins("JOIN role_permissions rp ON rp.permission_id = menus.permission_id").
		Joins("JOIN user_roles ur ON ur.role_id = rp.role_id").
		Where("ur.user_id = ? AND menus.is_active = true AND menus.deleted_date IS NULL", userID).
		Order("menus.sequence ASC, menus.id ASC").
		Find(&menus).Error

	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (r *menuRepository) GetMenusByPermissionIDs(ctx context.Context, permissionIDs []int) ([]domain.Menu, error) {
	if len(permissionIDs) == 0 {
		return []domain.Menu{}, nil
	}

	var menus []domain.Menu
	if err := r.db.WithContext(ctx).Model(&domain.Menu{}).
		Where("permission_id IN ? AND is_active = true AND deleted_date IS NULL", permissionIDs).
		Order("sequence ASC, id ASC").
		Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

func (r *menuRepository) GetMenusByIDs(ctx context.Context, ids []int64) ([]domain.Menu, error) {
	if len(ids) == 0 {
		return []domain.Menu{}, nil
	}

	var menus []domain.Menu
	if err := r.db.WithContext(ctx).Model(&domain.Menu{}).
		Where("id IN ? AND deleted_date IS NULL AND is_active = true", ids).
		Order("sequence ASC, id ASC").
		Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

func (r *menuRepository) ByID(ctx context.Context, id int64) (domain.Menu, error) {
	var m domain.Menu
	if err := r.db.WithContext(ctx).Where("id = ?", id).Preload("Parent").Preload("System").First(&m).Error; err != nil {
		return domain.Menu{}, err
	}
	return m, nil
}

func (r *menuRepository) Create(uctx context.Context, m domain.Menu) error {
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.CreatedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.CreatedOrgId = oid
	}
	return r.db.WithContext(uctx).Create(&m).Error
}

func (r *menuRepository) Update(uctx context.Context, id int64, m domain.Menu) error {
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.UpdatedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.UpdatedOrgId = oid
	}
	return r.db.WithContext(uctx).
		Model(&domain.Menu{}).
		Where("id = ?", id).
		Updates(&m).Error
}

func (r *menuRepository) Delete(uctx context.Context, id int64) error {
	var m domain.Menu
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = oid
	}
	m.DeletedDate = gorm.DeletedAt{Valid: true, Time: time.Now()}
	return r.db.WithContext(uctx).Model(&domain.Menu{}).Where("id = ?", id).Updates(&m).Error
}
