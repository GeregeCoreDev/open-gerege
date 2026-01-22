// Package service provides implementation for service
//
// File: news_service.go
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
)

type NewsService struct{ repo repository.NewsRepository }

func NewNewsService(repo repository.NewsRepository) *NewsService { return &NewsService{repo: repo} }

func (s *NewsService) List(ctx context.Context, q dto.NewsListQuery) ([]domain.News, int64, int, int, error) {
	return s.repo.List(ctx, q)
}

func (s *NewsService) GetByID(ctx context.Context, id int) (domain.News, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *NewsService) Create(ctx context.Context, req dto.NewsDto) error {
	m := domain.News{
		Title:    req.Title,
		Text:     req.Text,
		ImageUrl: req.ImageUrl,
	}
	return s.repo.Create(ctx, m)
}

func (s *NewsService) Update(ctx context.Context, id int, req dto.NewsDto) error {
	m := domain.News{
		Title:    req.Title,
		Text:     req.Text,
		ImageUrl: req.ImageUrl,
	}
	return s.repo.Update(ctx, id, m)
}

func (s *NewsService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
