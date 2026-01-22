// Package service provides implementation for service
//
// File: notification_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"
	"fmt"
	"time"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/config"
	"git.gerege.mn/backend-packages/httpx"
)

// Default Socket API base URL (fallback if config not provided)
const defaultSocketAPIBase = "https://socket.gerege.mn/api"

type NotificationService struct {
	repo repository.NotificationRepository
	http *httpx.Client
	cfg  *config.Config
}

func NewNotificationService(repo repository.NotificationRepository, cfg *config.Config) *NotificationService {
	return &NotificationService{
		repo: repo,
		http: httpx.New(3 * time.Second),
		cfg:  cfg,
	}
}

// getSocketAPIBase returns the socket API base URL
// TODO: Add Socket field to config.URLConfig when available
func (s *NotificationService) getSocketAPIBase() string {
	return defaultSocketAPIBase
}

// List for current user
func (s *NotificationService) List(ctx context.Context, userID int, p common.PaginationQuery) ([]domain.Notification, int64, int, int, error) {
	return s.repo.ListByUser(ctx, userID, p)
}

func (s *NotificationService) Groups(ctx context.Context, p common.PaginationQuery) ([]domain.NotificationGroup, int64, int, int, error) {
	return s.repo.ListGroups(ctx, p)
}

func (s *NotificationService) MarkGroupRead(ctx context.Context, userID, groupID int) error {
	return s.repo.MarkGroupRead(ctx, userID, groupID)
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID int) error {
	return s.repo.MarkAllRead(ctx, userID)
}

// Send: if UserID==0 => broadcast_all, else direct (dm)
func (s *NotificationService) Send(ctx context.Context, req dto.NotificationSendDto, createdUsername string) error {
	// 1) Create group
	group := domain.NotificationGroup{
		UserId:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
		Type:    typeOf(req.UserID), // "dm" | "broadcast_all"
		Tenant:  req.Tenant,
		// CreatedUserId:   createdBy,
		CreatedUsername: createdUsername,
	}
	g, err := s.repo.CreateGroup(ctx, group)
	if err != nil {
		return err
	}

	if req.UserID != 0 {
		// 2a) Direct notification
		n := domain.Notification{
			UserId:  req.UserID,
			Title:   req.Title,
			Content: req.Content,
			IsRead:  false,
			Type:    "dm",
			Tenant:  req.Tenant,
			GroupId: g.Id,
			// CreatedUserId:   createdBy,
			CreatedUsername: createdUsername,
		}
		if _, err := s.repo.CreateNotification(ctx, n); err != nil {
			return err
		}
		// 3a) Call socket /send
		body := map[string]any{
			"to":              fmt.Sprintf("%d", req.UserID),
			"idempotency_key": req.IdempotentKey,
			"body":            n,
		}
		_, _, err = httpx.PostJSON[map[string]any, any](ctx, s.http, s.getSocketAPIBase()+"/send", nil, body)
		return err
	}

	// 2b) Broadcast: call socket first
	bcast := domain.Notification{
		UserId:  0,
		Title:   req.Title,
		Content: req.Content,
		IsRead:  false,
		Tenant:  req.Tenant,
		Type:    "broadcast_all",
		GroupId: g.Id,
	}
	body := map[string]any{
		"tenant":          req.Tenant,
		"idempotency_key": req.IdempotentKey,
		"body":            bcast,
	}
	if _, _, err := httpx.PostJSON[map[string]any, any](ctx, s.http, s.getSocketAPIBase()+"/broadcast", nil, body); err != nil {
		return err
	}

	// 3b) Create notifications for all users
	ids, err := s.repo.AllUserIDs(ctx)
	if err != nil {
		return err
	}
	bulk := make([]domain.Notification, 0, len(ids))
	for _, uid := range ids {
		bulk = append(bulk, domain.Notification{
			UserId:  uid,
			Title:   req.Title,
			Content: req.Content,
			IsRead:  false,
			Tenant:  req.Tenant,
			Type:    "broadcast_all",
			GroupId: g.Id,
			// CreatedUserId:   createdBy,
			CreatedUsername: createdUsername,
		})
	}
	return s.repo.CreateNotificationsBulk(ctx, bulk)
}

func typeOf(userID int) string {
	if userID == 0 {
		return "broadcast_all"
	}
	return "dm"
}
