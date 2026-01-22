// Package repository provides implementation for repository
//
// File: system_repo.go
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

	xctx "git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
)

type SystemRepository interface {
	List(ctx context.Context, q dto.SystemListQuery) ([]domain.System, int64, int, int, error)
	ByID(ctx context.Context, id int) (domain.System, error)
	Create(ctx context.Context, m domain.System) error
	Update(ctx context.Context, id int, m domain.System) error
	Delete(ctx context.Context, id int) error // soft delete
	GetActiveModuleCount(uctx context.Context, id int) int64
	GetActiveRoleCount(uctx context.Context, id int) int64
}

type systemRepository struct {
	db *gorm.DB
}

func NewSystemRepository(db *gorm.DB) SystemRepository {
	return &systemRepository{db: db}
}

func (r *systemRepository) List(ctx context.Context, q dto.SystemListQuery) ([]domain.System, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	colMap := scopes.ColumnMap{
		"id":          "systems.id",
		"code":        "systems.code",
		"key":         "systems.key",
		"name":        "systems.name",
		"description": "systems.description",
		"is_active":   "systems.is_active",
		"icon":        "systems.icon",
		"path":        "systems.path",
		"sequence":    "systems.sequence",
	}

	tx := r.db.WithContext(ctx).Model(&domain.System{}).Scopes(
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
	if q.IsActive != nil {
		tx = tx.Where("is_active = ?", *q.IsActive)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.System
	if err := tx.Scopes(
		scopes.SortScope(colMap, utils.ParseSort(q.Sort), "sequence ASC, id DESC"),
	).Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	return items, total, page, size, nil
}

func (r *systemRepository) ByID(ctx context.Context, id int) (domain.System, error) {
	var m domain.System
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return domain.System{}, err
	}
	return m, nil
}

func (r *systemRepository) Create(uctx context.Context, m domain.System) error {
	// ctx-оос CreatedUser/Org онооно
	if userId, ok := xctx.GetValue[int](uctx, xctx.KeyUserID); ok {
		m.CreatedUserId = userId
	}
	if orgId, ok := xctx.GetValue[int](uctx, xctx.KeyOrgID); ok {
		m.CreatedOrgId = orgId
	}

	// Transaction ашиглаж system болон menu-г нэгэн зэрэг үүсгэх
	return WithTx(uctx, r.db, func(tx *gorm.DB) error {
		// System үүсгэх
		if err := tx.Create(&m).Error; err != nil {
			return err
		}

		// System үүсгэсний дараа автоматаар parent menu үүсгэх
		// Key хоосон бол code ашиглах
		menuKey := m.Key
		if menuKey == "" {
			menuKey = m.Code
		}

		isActive := true
		menu := domain.Menu{
			Code:         m.Code,
			Key:          menuKey,
			Name:         m.Name,
			Description:  m.Description,
			Icon:         m.Icon,
			Path:         "",
			Sequence:     int64(m.Sequence),
			PermissionID: nil,       // Permission байхгүй бол nil
			IsActive:     &isActive, // System үүсгэхэд menu идэвхтэй байх
			// ParentID-г оруулахгүй (NULL байх) - root menu
		}

		// Menu-ийн CreatedUser/Org-ийг system-ийнхтэй ижил болгох
		menu.CreatedUserId = m.CreatedUserId
		menu.CreatedOrgId = m.CreatedOrgId

		// Menu үүсгэх - parent_id болон permission_id-г NULL болгохын тулд зөвхөн шаардлагатай талбаруудыг сонгох
		return tx.Select("code", "key", "name", "description", "icon", "path", "sequence", "permission_id", "is_active", "created_user_id", "created_org_id").
			Create(&menu).Error
	})
}

func (r *systemRepository) Update(uctx context.Context, id int, m domain.System) error {
	// ctx-оос UpdatedUser/Org онооно
	if userId, ok := xctx.GetValue[int](uctx, xctx.KeyUserID); ok {
		m.UpdatedUserId = userId
	}
	if orgId, ok := xctx.GetValue[int](uctx, xctx.KeyOrgID); ok {
		m.UpdatedOrgId = orgId
	}

	// System update хийх
	if err := r.db.WithContext(uctx).
		Model(&domain.System{}).
		Where("id = ?", id).
		Updates(&m).Error; err != nil {
		return err
	}

	// System-ийн code-тай ижил code-тай menu-ийг олох
	var existingMenu domain.Menu
	err := r.db.WithContext(uctx).
		Where("code = ? AND deleted_date IS NULL", m.Code).
		First(&existingMenu).Error

	if err == nil && existingMenu.ID > 0 {
		// Menu олдвол update хийх
		menuUpdate := domain.Menu{
			Name:        m.Name,
			Description: m.Description,
			Icon:        m.Icon,
			Sequence:    int64(m.Sequence),
		}

		// UpdatedUser/Org-ийг system-ийнхтэй ижил болгох
		menuUpdate.UpdatedUserId = m.UpdatedUserId
		menuUpdate.UpdatedOrgId = m.UpdatedOrgId

		// Menu update хийх
		return r.db.WithContext(uctx).
			Model(&domain.Menu{}).
			Where("id = ?", existingMenu.ID).
			Updates(&menuUpdate).Error
	}

	// Menu олдохгүй бол зүгээр буцаах (system update хийгдсэн)
	return nil
}

func (r *systemRepository) Delete(uctx context.Context, id int) error {
	// Soft delete — таны загвартай адил
	var m domain.System
	if userId, ok := xctx.GetValue[int](uctx, xctx.KeyUserID); ok {
		m.DeletedUserId = userId
	}
	if orgId, ok := xctx.GetValue[int](uctx, xctx.KeyOrgID); ok {
		m.DeletedOrgId = orgId
	}
	m.DeletedDate = gorm.DeletedAt{Valid: true, Time: time.Now()}
	return r.db.WithContext(uctx).Model(&domain.System{}).Where("id = ?", id).Updates(&m).Error
}

func (r *systemRepository) GetActiveModuleCount(uctx context.Context, id int) int64 {
	cnt := int64(0)
	r.db.WithContext(uctx).Model(&domain.Module{}).Where("system_id = ? AND is_active = true", id).Count(&cnt)
	return cnt
}

func (r *systemRepository) GetActiveRoleCount(uctx context.Context, id int) int64 {
	cnt := int64(0)
	r.db.WithContext(uctx).Model(&domain.Role{}).Where("system_id = ? AND is_active = true", id).Count(&cnt)
	return cnt
}
