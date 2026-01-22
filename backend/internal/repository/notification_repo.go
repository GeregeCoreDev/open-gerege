// Package repository provides implementation for repository
//
// File: notification_repo.go
// Description: implementation for repository
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package repository

import (
	"context"

	"templatev25/internal/domain"
	"git.gerege.mn/backend-packages/common"

	"git.gerege.mn/backend-packages/scopes"
	"git.gerege.mn/backend-packages/utils"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	ListByUser(ctx context.Context, userID int, p common.PaginationQuery) ([]domain.Notification, int64, int, int, error)
	MarkGroupRead(ctx context.Context, userID, groupID int) error
	MarkAllRead(ctx context.Context, userID int) error

	ListGroups(ctx context.Context, p common.PaginationQuery) ([]domain.NotificationGroup, int64, int, int, error)
	CreateGroup(ctx context.Context, g domain.NotificationGroup) (domain.NotificationGroup, error)

	CreateNotification(ctx context.Context, n domain.Notification) (domain.Notification, error)
	CreateNotificationsBulk(ctx context.Context, ns []domain.Notification) error

	AllUserIDs(ctx context.Context) ([]int, error)
}

type notificationRepository struct{ db *gorm.DB }

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) ListByUser(ctx context.Context, userID int, p common.PaginationQuery) ([]domain.Notification, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(p)
	colMap := scopes.ColumnMap{
		"id":       "notifications.id",
		"user_id":  "notifications.user_id",
		"is_read":  "notifications.is_read",
		"type":     "notifications.type",
		"tenant":   "notifications.tenant",
		"title":    "notifications.title",
		"content":  "notifications.content",
		"group_id": "notifications.group_id",
	}
	tx := r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("user_id = ?", userID).
		Scopes(
			scopes.SearchScope(colMap, utils.ParseSearch(p.Search)),
			scopes.DateScope(p.CreatedFrom, p.CreatedTo),
		)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.Notification
	if err := tx.Scopes(scopes.SortScope(colMap, utils.ParseSort(p.Sort), "id DESC")).
		Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	return items, total, page, size, nil
}

func (r *notificationRepository) MarkGroupRead(ctx context.Context, userID, groupID int) error {
	return r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllRead(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}

func (r *notificationRepository) ListGroups(ctx context.Context, p common.PaginationQuery) ([]domain.NotificationGroup, int64, int, int, error) {
	page, size, offset := utils.OffsetLimit(p)
	colMap := scopes.ColumnMap{
		"id":      "notification_groups.id",
		"user_id": "notification_groups.user_id",
		"type":    "notification_groups.type",
		"tenant":  "notification_groups.tenant",
		"title":   "notification_groups.title",
	}
	tx := r.db.WithContext(ctx).
		Model(&domain.NotificationGroup{}).
		Scopes(
			scopes.SearchScope(colMap, utils.ParseSearch(p.Search)),
			scopes.DateScope(p.CreatedFrom, p.CreatedTo),
		)

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, 0, 0, err
	}

	var items []domain.NotificationGroup
	if err := tx.Scopes(scopes.SortScope(colMap, utils.ParseSort(p.Sort), "id DESC")).
		Offset(offset).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, 0, 0, err
	}
	return items, total, page, size, nil
}

func (r *notificationRepository) CreateGroup(ctx context.Context, g domain.NotificationGroup) (domain.NotificationGroup, error) {
	if err := r.db.WithContext(ctx).Create(&g).Error; err != nil {
		return domain.NotificationGroup{}, err
	}
	return g, nil
}

func (r *notificationRepository) CreateNotification(ctx context.Context, n domain.Notification) (domain.Notification, error) {
	if err := r.db.WithContext(ctx).Create(&n).Error; err != nil {
		return domain.Notification{}, err
	}
	return n, nil
}

func (r *notificationRepository) CreateNotificationsBulk(ctx context.Context, ns []domain.Notification) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, n := range ns {
			if err := tx.Create(&n).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *notificationRepository) AllUserIDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}
