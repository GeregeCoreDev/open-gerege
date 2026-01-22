// Package service provides implementation for service
//
// File: api_log_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-01-09
// Last Updated: 2025-01-09
package service

import (
	"context"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/repository"
)

type APILogService interface {
	List(ctx context.Context, q dto.APILogListQuery) ([]domain.APILog, int64, int, int, error)
}

type apiLogService struct {
	repo repository.APILogRepository
}

func NewAPILogService(repo repository.APILogRepository) APILogService {
	return &apiLogService{repo: repo}
}

func (s *apiLogService) List(ctx context.Context, q dto.APILogListQuery) ([]domain.APILog, int64, int, int, error) {
	return s.repo.List(ctx, q)
}
