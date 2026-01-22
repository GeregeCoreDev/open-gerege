// Package repository provides implementation for repository
//
// File: platform_icon_repo.go
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
	"git.gerege.mn/backend-packages/ctx"

	"gorm.io/gorm"
)

type AppServiceIconRepository interface {
	List(ctx context.Context) ([]domain.AppServiceIcon, error)
	Create(ctx context.Context, m domain.AppServiceIcon) error
	Update(ctx context.Context, id int, m domain.AppServiceIcon) error
	DeleteSoft(ctx context.Context, id int) error
}

type AppServiceIconGroupRepository interface {
	List(ctx context.Context) ([]domain.AppServiceIconGroup, error)
	ListGroupsWithIcons(ctx context.Context) ([]domain.AppServiceIconGroup, error)
	Create(ctx context.Context, m domain.AppServiceIconGroup) error
	Update(ctx context.Context, id int, m domain.AppServiceIconGroup) error
	DeleteSoft(ctx context.Context, id int) error
}

type appServiceIconRepo struct{ db *gorm.DB }
type appServiceIconGroupRepo struct{ db *gorm.DB }

func NewAppServiceIconRepository(db *gorm.DB) AppServiceIconRepository {
	return &appServiceIconRepo{db: db}
}

func NewAppServiceIconGroupRepository(db *gorm.DB) AppServiceIconGroupRepository {
	return &appServiceIconGroupRepo{db: db}
}

func (r *appServiceIconGroupRepo) List(ctx context.Context) ([]domain.AppServiceIconGroup, error) {
	var items []domain.AppServiceIconGroup
	err := r.db.WithContext(ctx).
		Order("seq ASC").
		Find(&items).Error
	return items, err
}

func (r *appServiceIconGroupRepo) ListGroupsWithIcons(ctx context.Context) ([]domain.AppServiceIconGroup, error) {
	var items []domain.AppServiceIconGroup
	err := r.db.WithContext(ctx).
		Preload("AppServices", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Childs", func(tx *gorm.DB) *gorm.DB {
				isPublic := true
				return tx.Where("is_public = ?", &isPublic).Order("seq ASC")
			}).Where("is_public = ? AND parent_id = 0", true).Order("seq ASC")
		}).
		Order("seq ASC").
		Find(&items).Error
	return items, err
}

func (r *appServiceIconGroupRepo) Create(uctx context.Context, m domain.AppServiceIconGroup) error {
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.CreatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.CreatedOrgId = orgId
	}
	return r.db.WithContext(uctx).Create(&m).Error
}

func (r *appServiceIconGroupRepo) Update(uctx context.Context, id int, m domain.AppServiceIconGroup) error {
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.UpdatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.UpdatedOrgId = orgId
	}

	return r.db.WithContext(uctx).Model(&domain.AppServiceIconGroup{}).
		Where("id = ?", id).Updates(&m).Error
}

func (r *appServiceIconGroupRepo) DeleteSoft(uctx context.Context, id int) error {
	var m domain.AppServiceIconGroup

	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = orgId
	}

	m.DeletedDate = gorm.DeletedAt{Time: time.Now(), Valid: true}

	return r.db.WithContext(uctx).
		Model(&domain.AppServiceIconGroup{}).
		Where("id = ?", id).
		Updates(&m).Error
}

// ---------- App Service Icon ----------

func (r *appServiceIconRepo) List(ctx context.Context) ([]domain.AppServiceIcon, error) {
	var items []domain.AppServiceIcon
	err := r.db.WithContext(ctx).
		Order("seq ASC").
		Preload("Group").
		Preload("Parent").
		Preload("Childs").
		Find(&items).Error
	return items, err
}

func (r *appServiceIconRepo) Create(uctx context.Context, m domain.AppServiceIcon) error {
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.CreatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.CreatedOrgId = orgId
	}
	return r.db.WithContext(uctx).Create(&m).Error
}

func (r *appServiceIconRepo) Update(uctx context.Context, id int, m domain.AppServiceIcon) error {
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.UpdatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.UpdatedOrgId = orgId
	}

	return r.db.WithContext(uctx).Model(&domain.AppServiceIcon{}).
		Where("id = ?", id).Updates(&m).Error
}

func (r *appServiceIconRepo) DeleteSoft(uctx context.Context, id int) error {
	var m domain.AppServiceIcon

	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = orgId
	}

	m.DeletedDate = gorm.DeletedAt{Time: time.Now(), Valid: true}

	return r.db.WithContext(uctx).
		Model(&domain.AppServiceIcon{}).
		Where("id = ?", id).
		Updates(&m).Error
}

