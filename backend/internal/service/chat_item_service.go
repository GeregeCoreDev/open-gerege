// Package service provides implementation for service
//
// File: chat_item_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"
	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"templatev25/internal/repository"

	"go.uber.org/zap"
)

type ChatItemService struct {
	repo repository.ChatItemRepository
	log  *zap.Logger
}

func NewChatItemService(r repository.ChatItemRepository, log *zap.Logger) *ChatItemService {
	return &ChatItemService{repo: r, log: log}
}

func (s *ChatItemService) GetByKey(ctx context.Context, key string) (domain.ChatItem, error) {
	return s.repo.FindByKey(ctx, key)
}

func (s *ChatItemService) List(ctx context.Context, q dto.ChatItemQuery) ([]domain.ChatItem, int64, int, int, error) {
	return s.repo.List(ctx, q)
}

func (s *ChatItemService) Create(ctx context.Context, d dto.ChatItemCreateDto) error {
	m := domain.ChatItem{
		Key:    d.Key,
		Answer: d.Answer,
	}
	return s.repo.Create(ctx, m)
}

func (s *ChatItemService) Update(ctx context.Context, id int, d dto.ChatItemUpdateDto) error {
	m := domain.ChatItem{
		Key:    d.Key,
		Answer: d.Answer,
	}
	return s.repo.Update(ctx, id, m)
}

func (s *ChatItemService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
