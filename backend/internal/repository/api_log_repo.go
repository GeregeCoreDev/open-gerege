// Package repository provides implementation for repository
//
// File: api_log_repo.go
// Description: implementation for repository
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-01-09
// Last Updated: 2025-01-09
package repository

import (
	"context"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/config"
	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"
	"gorm.io/gorm"
)

type APILogRepository interface {
	Create(ctx context.Context, log domain.APILog) error
	List(ctx context.Context, q dto.APILogListQuery) ([]domain.APILog, int64, int, int, error)
}

type apiLogRepository struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAPILogRepository(db *gorm.DB) APILogRepository {
	return &apiLogRepository{db: db, cfg: nil}
}

func NewAPILogRepositoryWithConfig(db *gorm.DB, cfg *config.Config) APILogRepository {
	return &apiLogRepository{db: db, cfg: cfg}
}

func (r *apiLogRepository) Create(ctx context.Context, log domain.APILog) error {
	// GORM automatically adds schema prefix via TablePrefix in db.NewPostgres
	// So we can use Model() directly, which will use TableName() from domain.APILog
	return r.db.WithContext(ctx).Model(&domain.APILog{}).Create(&log).Error
}

func (r *apiLogRepository) List(ctx context.Context, q dto.APILogListQuery) ([]domain.APILog, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(q.PaginationQuery)

	colMap := scopes.ColumnMap{
		"id":          "logs.id",
		"method":      "logs.method",
		"path":        "logs.path",
		"status_code": "logs.status_code",
		"user_id":     "logs.user_id",
		"org_id":      "logs.org_id",
		"ip":          "logs.ip",
		"username":    "logs.username",
	}

	tx := r.db.WithContext(ctx).Model(&domain.APILog{}).Scopes(
		scopes.SearchScope(colMap, utils.ParseSearch(q.Search)),
		scopes.DateScope(q.CreatedFrom, q.CreatedTo),
	)

	if q.Method != "" {
		tx = tx.Where("method = ?", q.Method)
	}
	if q.Path != "" {
		tx = tx.Where("path ILIKE ?", "%"+q.Path+"%")
	}
	if q.StatusCode != nil {
		tx = tx.Where("status_code = ?", *q.StatusCode)
	}
	if q.UserID != nil {
		tx = tx.Where("user_id = ?", *q.UserID)
	}
	if q.OrgID != nil {
		tx = tx.Where("org_id = ?", *q.OrgID)
	}
	if q.IP != "" {
		tx = tx.Where("ip ILIKE ?", "%"+q.IP+"%")
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.APILog
	if err := tx.Scopes(
		scopes.SortScope(colMap, utils.ParseSort(q.Sort), "created_date DESC, id DESC"),
	).Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	return items, total, page, size, nil
}
