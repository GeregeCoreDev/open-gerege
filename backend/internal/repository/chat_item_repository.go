// Package repository provides implementation for repository
//
// File: chat_item_repository.go
// Description: implementation for repository
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package repository

import (
	"context"
	"strings"
	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"
	"time"

	"gorm.io/gorm"
)

type ChatItemRepository interface {
	List(ctx context.Context, q dto.ChatItemQuery) ([]domain.ChatItem, int64, int, int, error)
	ByID(ctx context.Context, id int) (domain.ChatItem, error)
	Create(ctx context.Context, m domain.ChatItem) error
	Update(ctx context.Context, id int, m domain.ChatItem) error
	Delete(ctx context.Context, id int) error
	FindByKey(ctx context.Context, key string) (domain.ChatItem, error)
}

type chatItemRepository struct {
	db *gorm.DB
}

func NewChatItemRepository(db *gorm.DB) ChatItemRepository {
	return &chatItemRepository{db: db}
}

func (r *chatItemRepository) FindByKey(ctx context.Context, key string) (domain.ChatItem, error) {
	var item domain.ChatItem
	err := r.db.WithContext(ctx).
		Where("key = ? AND deleted_date IS NULL", key).
		Last(&item).Error
	if err != nil {
		return domain.ChatItem{}, err
	}
	return item, nil
}

func (r *chatItemRepository) List(ctx context.Context, q dto.ChatItemQuery) ([]domain.ChatItem, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	colMap := scopes.ColumnMap{
		"key":    "chat_items.key",
		"answer": "chat_items.answer",
	}

	tx := r.db.WithContext(ctx).Model(&domain.ChatItem{}).Scopes(
		scopes.SearchScope(colMap, utils.ParseSearch(q.Search)),
	)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.ChatItem
	if err := tx.Offset(offset).Limit(size).
		Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	return items, total, page, size, nil
}

func (r *chatItemRepository) ByID(ctx context.Context, id int) (domain.ChatItem, error) {
	var m domain.ChatItem
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return domain.ChatItem{}, err
	}
	return m, nil
}

func (r *chatItemRepository) Create(uctx context.Context, m domain.ChatItem) error {
	if user, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.CreatedUserId = user
	}
	if org, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.CreatedOrgId = org
	}

	m.Key = strings.ToLower(m.Key)
	m.Answer = strings.ToLower(m.Answer)
	return r.db.WithContext(uctx).Create(&m).Error
}

func (r *chatItemRepository) Update(uctx context.Context, id int, m domain.ChatItem) error {
	if user, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.UpdatedUserId = user
	}
	if org, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.UpdatedOrgId = org
	}

	m.Key = strings.ToLower(m.Key)
	m.Answer = strings.ToLower(m.Answer)
	return r.db.WithContext(uctx).Model(&domain.ChatItem{}).
		Where("id = ?", id).
		Updates(&m).Error
}

func (r *chatItemRepository) Delete(uctx context.Context, id int) error {
	m := domain.ChatItem{}
	m.DeletedDate = gorm.DeletedAt{Valid: true, Time: time.Now()}

	if user, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = user
	}
	if org, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = org
	}

	return r.db.WithContext(uctx).Where("id = ?", id).Updates(&m).Error
}
