// Package repository provides implementation for repository
//
// File: terminal_repo.go
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
	"git.gerege.mn/backend-packages/common"

	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
)

type TerminalRepository interface {
	List(ctx context.Context, p common.PaginationQuery) ([]domain.Terminal, int64, int, int, error)
	Create(ctx context.Context, m domain.Terminal) error
	Update(ctx context.Context, id int, m domain.Terminal) error
	DeleteSoft(uctx context.Context, id int) error
}

type terminalRepository struct{ db *gorm.DB }

func NewTerminalRepository(db *gorm.DB) TerminalRepository {
	return &terminalRepository{db: db}
}

func (r *terminalRepository) List(ctx context.Context, p common.PaginationQuery) ([]domain.Terminal, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(p)

	colMap := scopes.ColumnMap{
		"id":     "terminals.id",
		"name":   "terminals.name",
		"serial": "terminals.serial",
		"org_id": "terminals.org_id",
	}

	tx := r.db.WithContext(ctx).Model(&domain.Terminal{}).Preload("Organization").Scopes(
		scopes.SearchScope(colMap, utils.ParseSearch(p.Search)),
		scopes.DateScope(p.CreatedFrom, p.CreatedTo),
	)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.Terminal
	if err := tx.Scopes(
		scopes.SortScope(colMap, utils.ParseSort(p.Sort), "id DESC"),
	).Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	return items, total, page, size, nil
}

func (r *terminalRepository) Create(uctx context.Context, m domain.Terminal) error {
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.CreatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.CreatedOrgId = orgId
	}
	return r.db.WithContext(uctx).Create(&m).Error
}

func (r *terminalRepository) Update(uctx context.Context, id int, m domain.Terminal) error {
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.UpdatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.UpdatedOrgId = orgId
	}

	return r.db.WithContext(uctx).
		Model(&domain.Terminal{}).
		Where("id = ?", id).
		Updates(&m).Error
}

func (r *terminalRepository) DeleteSoft(uctx context.Context, id int) error {

	m := domain.Terminal{}
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = orgId
	}

	m.DeletedDate = gorm.DeletedAt{Valid: true, Time: time.Now()}
	if err := r.db.WithContext(uctx).Where("id = ?", id).Updates(&m).Error; err != nil {
		return err
	}

	return nil
}
