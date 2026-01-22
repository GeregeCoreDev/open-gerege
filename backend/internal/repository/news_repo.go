// Package repository provides implementation for repository
//
// File: news_repo.go
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

	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
)

type NewsRepository interface {
	List(ctx context.Context, q dto.NewsListQuery) ([]domain.News, int64, int, int, error)
	GetByID(ctx context.Context, id int) (domain.News, error)
	Create(ctx context.Context, m domain.News) error
	Update(ctx context.Context, id int, m domain.News) error
	Delete(uctx context.Context, id int) error
}

type newsRepository struct{ db *gorm.DB }

func NewNewsRepository(db *gorm.DB) NewsRepository { return &newsRepository{db: db} }

func (r *newsRepository) List(ctx context.Context, q dto.NewsListQuery) ([]domain.News, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	colMap := scopes.ColumnMap{
		"id":        "news.id",
		"title":     "news.title",
		"text":      "news.text",
		"image_url": "news.image_url",
	}

	tx := r.db.WithContext(ctx).Model(&domain.News{}).
		Scopes(
			scopes.SearchScope(colMap, utils.ParseSearch(q.Search)),
			scopes.DateScope(q.CreatedFrom, q.CreatedTo),
		)

	if q.CategoryID != 0 {
		tx = tx.Where("category_id = ?", q.CategoryID)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.News
	if err := tx.Scopes(scopes.SortScope(colMap, utils.ParseSort(q.Sort), "id DESC")).
		Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	return items, total, page, size, nil
}

func (r *newsRepository) GetByID(ctx context.Context, id int) (domain.News, error) {
	var m domain.News
	err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error
	return m, err
}

func (r *newsRepository) Create(uctx context.Context, m domain.News) error {
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.CreatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.CreatedOrgId = orgId
	}
	if err := r.db.WithContext(uctx).Create(&m).Error; err != nil {
		return err
	}
	return nil
}

func (r *newsRepository) Update(uctx context.Context, id int, m domain.News) error {
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.UpdatedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.UpdatedOrgId = orgId
	}
	if err := r.db.WithContext(uctx).Model(&domain.News{}).Where("id = ?", id).Updates(&m).Error; err != nil {
		return err
	}
	return nil
}

func (r *newsRepository) Delete(uctx context.Context, id int) error {
	m := domain.News{}
	if userId, ok := ctx.GetValue[int](uctx, ctx.KeyUserID); ok {
		m.DeletedUserId = userId
	}
	if orgId, ok := ctx.GetValue[int](uctx, ctx.KeyOrgID); ok {
		m.DeletedOrgId = orgId
	}
	m.Title += " (Deleted)"
	m.DeletedDate = gorm.DeletedAt{Valid: true, Time: time.Now()}
	if err := r.db.WithContext(uctx).Where("id = ?", id).Updates(&m).Error; err != nil {
		return err
	}
	return nil

}
