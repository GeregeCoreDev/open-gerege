// Package repository provides implementation for repository
//
// File: public_file_repo.go
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

	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
)

type PublicFileRepository interface {
	List(ctx context.Context, q dto.PublicFileListQuery) ([]domain.PublicFile, int64, int, int, error)
	Create(ctx context.Context, m domain.PublicFile) (domain.PublicFile, error)
	DeleteByID(ctx context.Context, id int) (domain.PublicFile, error)
	GetByName(ctx context.Context, name string) (domain.PublicFile, error)
	GetByID(ctx context.Context, id int) (domain.PublicFile, error)
}

type publicFileRepository struct{ db *gorm.DB }

func NewPublicFileRepository(db *gorm.DB) PublicFileRepository {
	return &publicFileRepository{db: db}
}

func (r *publicFileRepository) List(ctx context.Context, q dto.PublicFileListQuery) ([]domain.PublicFile, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	colMap := scopes.ColumnMap{
		"id":          "public_files.id",
		"name":        "public_files.name",
		"extension":   "public_files.extension",
		"description": "public_files.description",
		"file_url":    "public_files.file_url",
	}

	tx := r.db.WithContext(ctx).
		Model(&domain.PublicFile{}).
		Scopes(
			scopes.SearchScope(colMap, utils.ParseSearch(q.Search)),
			scopes.DateScope(q.CreatedFrom, q.CreatedTo),
		)

	if d := q.Description; d != "" {
		tx = tx.Where("description ILIKE ?", "%"+d+"%")
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.PublicFile
	if err := tx.Scopes(
		scopes.SortScope(colMap, utils.ParseSort(q.Sort), "id DESC"),
	).Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	return items, total, page, size, nil
}

func (r *publicFileRepository) Create(ctx context.Context, m domain.PublicFile) (domain.PublicFile, error) {
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return domain.PublicFile{}, err
	}
	return m, nil
}

func (r *publicFileRepository) DeleteByID(ctx context.Context, id int) (domain.PublicFile, error) {
	var pf domain.PublicFile
	if err := r.db.WithContext(ctx).Take(&pf, "id = ?", id).Error; err != nil {
		return domain.PublicFile{}, err
	}
	if err := r.db.WithContext(ctx).Delete(&pf).Error; err != nil {
		return domain.PublicFile{}, err
	}
	return pf, nil
}

func (r *publicFileRepository) GetByName(ctx context.Context, name string) (domain.PublicFile, error) {
	var pf domain.PublicFile
	err := r.db.WithContext(ctx).Take(&pf, "name = ?", name).Error
	return pf, err
}

func (r *publicFileRepository) GetByID(ctx context.Context, id int) (domain.PublicFile, error) {
	var pf domain.PublicFile
	err := r.db.WithContext(ctx).Take(&pf, "id = ?", id).Error
	return pf, err
}
