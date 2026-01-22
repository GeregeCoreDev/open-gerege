// Package repository provides implementation for repository
//
// File: module_repo.go
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
)

type ModuleRepository interface {
	List(ctx context.Context, q dto.ModuleListQuery) ([]domain.Module, int64, int, int, error)
	ByID(ctx context.Context, id int) (domain.Module, error)
	Create(ctx context.Context, m domain.Module) error
	Update(ctx context.Context, id int, m domain.Module) error
	Delete(ctx context.Context, id int) error
}

type moduleRepository struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewModuleRepository(db *gorm.DB, cfg *config.Config) ModuleRepository {
	return &moduleRepository{db: db, cfg: cfg}
}

func (r *moduleRepository) List(ctx context.Context, q dto.ModuleListQuery) ([]domain.Module, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	colMap := scopes.ColumnMap{
		"id":          "modules.id",
		"code":        "modules.code",
		"name":        "modules.name",
		"description": "modules.description",
		"is_active":   "modules.is_active",
		"system_id":   "modules.system_id",
	}

	tx := r.db.WithContext(ctx).Model(&domain.Module{}).Scopes(
		scopes.SearchScope(colMap, utils.ParseSearch(q.Search)),
		scopes.DateScope(q.CreatedFrom, q.CreatedTo),
	)

	if q.Code != "" {
		tx = tx.Where("code ILIKE ?", "%"+q.Code+"%")
	}
	if q.Name != "" {
		tx = tx.Where("name ILIKE ?", "%"+q.Name+"%")
	}
	if q.IsActive != nil {
		tx = tx.Where("is_active = ?", *q.IsActive)
	}
	if q.SystemID != 0 {
		tx = tx.Where("system_id = ?", q.SystemID)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.Module
	if err := tx.Scopes(
		scopes.SortScope(colMap, utils.ParseSort(q.Sort), "id DESC"),
	).Preload("System").Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	return items, total, page, size, nil
}

func (r *moduleRepository) ByID(ctx context.Context, id int) (domain.Module, error) {
	var m domain.Module
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return domain.Module{}, err
	}
	return m, nil
}

func (r *moduleRepository) Create(uctx context.Context, m domain.Module) error {
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.CreatedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.CreatedOrgId = oid
	}
	return r.db.WithContext(uctx).Create(&m).Error
}

func (r *moduleRepository) Update(uctx context.Context, id int, m domain.Module) error {
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.UpdatedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.UpdatedOrgId = oid
	}
	return r.db.WithContext(uctx).
		Model(&domain.Module{}).
		Where("id = ?", id).
		Updates(&m).Error
}

func (r *moduleRepository) Delete(uctx context.Context, id int) error {
	var m domain.Module
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = oid
	}
	m.DeletedDate = gorm.DeletedAt{Valid: true, Time: time.Now()}
	return r.db.WithContext(uctx).Model(&domain.Module{}).Where("id = ?", id).Updates(&m).Error
}
