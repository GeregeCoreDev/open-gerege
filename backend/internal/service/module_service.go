// Package service provides implementation for service
//
// File: module_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"
	"errors"
	"strings"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"templatev25/internal/repository"
)

type ModuleService interface {
	List(ctx context.Context, q dto.ModuleListQuery) ([]domain.Module, int64, int, int, error)
	ByID(ctx context.Context, id int) (domain.Module, error)
	Create(ctx context.Context, req dto.ModuleCreateDto) error
	Update(ctx context.Context, id int, req dto.ModuleUpdateDto) error
	Delete(ctx context.Context, id int) error
}

type moduleService struct{ repo repository.ModuleRepository }

func NewModuleService(repo repository.ModuleRepository) ModuleService {
	return &moduleService{repo: repo}
}

func (s *moduleService) List(ctx context.Context, q dto.ModuleListQuery) ([]domain.Module, int64, int, int, error) {
	return s.repo.List(ctx, q)
}

func (s *moduleService) ByID(ctx context.Context, id int) (domain.Module, error) {
	return s.repo.ByID(ctx, id)
}

func (s *moduleService) Create(ctx context.Context, req dto.ModuleCreateDto) error {
	// Code-г lower case болгох
	code := strings.ToLower(req.Code)
	
	m := domain.Module{
		Code:        code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
		SystemID:    req.SystemID,
	}
	return s.repo.Create(ctx, m)
}

func (s *moduleService) Update(ctx context.Context, id int, req dto.ModuleUpdateDto) error {
	// Code-г lower case болгох
	code := strings.ToLower(req.Code)
	
	m := domain.Module{
		Code:        code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
		SystemID:    req.SystemID,
	}
	return s.repo.Update(ctx, id, m)
}

func (s *moduleService) Delete(ctx context.Context, id int) error {
	existing, err := s.repo.ByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.IsActive != nil && *existing.IsActive {
		return errors.New("модуль идэвхитэй тул устгах боломжгүй")
	}
	return s.repo.Delete(ctx, id)
}
