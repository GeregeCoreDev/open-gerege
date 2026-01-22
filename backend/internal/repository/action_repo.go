// Package repository provides implementation for repository
//
// File: action_repo.go
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

	"time"

	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
)

type ActionRepository interface {
	List(ctx context.Context, q dto.ActionQuery) ([]domain.Action, int64, int, int, error)
	ByID(ctx context.Context, id int64) (domain.Action, error)
	Create(ctx context.Context, m domain.Action) error
	Update(ctx context.Context, id int64, m domain.Action) error
	Delete(ctx context.Context, id int64) error
}

type actionRepository struct {
	db *gorm.DB
}

func NewActionRepository(db *gorm.DB) ActionRepository {
	return &actionRepository{
		db: db,
	}
}

func (r *actionRepository) List(ctx context.Context, q dto.ActionQuery) ([]domain.Action, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	colMap := scopes.ColumnMap{
		"id":          "actions.id",
		"code":        "actions.code",
		"name":        "actions.name",
		"description": "actions.description",
	}

	tx := r.db.WithContext(ctx).Model(&domain.Action{}).Scopes(
		scopes.SearchScope(colMap, utils.ParseSearch(q.Search)),
		scopes.DateScope(q.CreatedFrom, q.CreatedTo),
	)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.Action
	if err := tx.Scopes(
		scopes.SortScope(colMap, utils.ParseSort(q.Sort), "id DESC"),
	).Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	return items, total, page, size, nil
}

func (r *actionRepository) ByID(ctx context.Context, id int64) (domain.Action, error) {
	var m domain.Action
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return domain.Action{}, err
	}
	return m, nil
}

func (r *actionRepository) Create(uctx context.Context, m domain.Action) error {
	// ctx-оос CreatedUser/Org онооно
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.CreatedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.CreatedOrgId = oid
	}

	// Action үүсгэх
	return r.db.WithContext(uctx).Create(&m).Error
}

func (r *actionRepository) Update(uctx context.Context, id int64, m domain.Action) error {
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.UpdatedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.UpdatedOrgId = oid
	}
	return r.db.WithContext(uctx).Model(&domain.Action{}).Where("id = ?", id).Updates(&m).Error
}

func (r *actionRepository) Delete(uctx context.Context, id int64) error {
	m := domain.Action{}
	if uid, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = uid
	}
	if oid, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = oid
	}
	m.DeletedDate = gorm.DeletedAt{Valid: true, Time: time.Now()}
	return r.db.WithContext(uctx).Where("id = ?", id).Updates(&m).Error
}
