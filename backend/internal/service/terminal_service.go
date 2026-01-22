// Package service provides implementation for service
//
// File: terminal_service.go
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
	"git.gerege.mn/backend-packages/common"
	"templatev25/internal/repository"
)

type TerminalService struct{ repo repository.TerminalRepository }

func NewTerminalService(repo repository.TerminalRepository) *TerminalService {
	return &TerminalService{repo: repo}
}

func (s *TerminalService) List(ctx context.Context, p common.PaginationQuery) ([]domain.Terminal, int64, int, int, error) {
	return s.repo.List(ctx, p)
}

func (s *TerminalService) Create(ctx context.Context, req dto.TerminalCreateDto) error {
	m := domain.Terminal{
		Serial: req.Serial,
		Name:   req.Name,
	}
	if req.OrgId != nil {
		m.OrgId = req.OrgId
	}
	return s.repo.Create(ctx, m)
}

func (s *TerminalService) Update(ctx context.Context, id int, req dto.TerminalUpdateDto) error {
	m := domain.Terminal{
		Serial: req.Serial,
		Name:   req.Name,
	}
	if req.OrgId != nil {
		m.OrgId = req.OrgId
	}
	return s.repo.Update(ctx, id, m)
}

func (s *TerminalService) Delete(ctx context.Context, id int) error {
	return s.repo.DeleteSoft(ctx, id)
}
