// Package service provides implementation for service
//
// File: action_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"
	"strings"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"templatev25/internal/repository"

	"go.uber.org/zap"
)

type ActionService struct {
	repo repository.ActionRepository
	log  *zap.Logger
}

func NewActionService(repo repository.ActionRepository, log *zap.Logger) *ActionService {
	return &ActionService{
		repo: repo,
		log:  log,
	}
}

func (s *ActionService) ListFilteredPaged(ctx context.Context, q dto.ActionQuery) ([]domain.Action, int64, int, int, error) {
	return s.repo.List(ctx, q)
}

func (s *ActionService) ByID(ctx context.Context, id int64) (domain.Action, error) {
	return s.repo.ByID(ctx, id)
}

func (s *ActionService) Create(ctx context.Context, req dto.ActionCreateDto) error {
	// Code-г lower case болгох
	code := strings.ToLower(req.Code)
	
	// Action үүсгэх
	m := domain.Action{
		Code:        code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}
	return s.repo.Create(ctx, m)
}

func (s *ActionService) Update(ctx context.Context, id int64, req dto.ActionUpdateDto) error {
	// Code-г lower case болгох
	code := strings.ToLower(req.Code)
	
	// Action засах
	m := domain.Action{
		Code:        code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}
	return s.repo.Update(ctx, id, m)
}

func (s *ActionService) Delete(ctx context.Context, id int64) error {
	// Action устгах
	return s.repo.Delete(ctx, id)
}
