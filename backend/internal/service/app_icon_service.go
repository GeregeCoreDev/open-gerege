// Package service provides implementation for service
//
// File: platform_icon_service.go
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

type AppServiceIconService struct {
	repo repository.AppServiceIconRepository
}

type AppServiceIconGroup struct {
	repo repository.AppServiceIconGroupRepository
}

func NewAppServiceIconService(repo repository.AppServiceIconRepository) *AppServiceIconService {
	return &AppServiceIconService{repo: repo}
}

func NewAppServiceIconGroup(repo repository.AppServiceIconGroupRepository) *AppServiceIconGroup {
	return &AppServiceIconGroup{repo: repo}
}


// ---- App Service Icon ----

func (s *AppServiceIconService) List(ctx context.Context) ([]domain.AppServiceIcon, error) {
	return s.repo.List(ctx)
}

func (s *AppServiceIconService) Create(ctx context.Context, req dto.AppServiceIconDto) error {
	m := domain.AppServiceIcon{
		Name:          req.Name,
		NameEn:        req.NameEn,
		Icon:          req.Icon,
		IconApp:       req.IconApp,
		IconTablet:    req.IconTablet,
		IconKiosk:     req.IconKiosk,
		Link:          req.Link,
		GroupId:       req.GroupId,
		Seq:           req.Seq,
		IsNative:      req.IsNative,
		IsPublic:      req.IsPublic,
		IsFeatured:    req.IsFeatured,
		FeaturedIcon:  req.FeaturedIcon,
		IsBestSelling: req.IsBestSelling,
		FeatureSeq:    req.FeatureSeq,
		Description:   req.Description,
		SystemCode:    req.SystemCode,
		IsGroup:       req.IsGroup,
		ParentId:      req.ParentId,
		WebLink:       req.WebLink,
	}
	return s.repo.Create(ctx, m)
}

func (s *AppServiceIconService) Update(ctx context.Context, id int, req dto.AppServiceIconDto) error {
	m := domain.AppServiceIcon{
		Name:          req.Name,
		NameEn:        req.NameEn,
		Icon:          req.Icon,
		IconApp:       req.IconApp,
		IconTablet:    req.IconTablet,
		IconKiosk:     req.IconKiosk,
		Link:          req.Link,
		GroupId:       req.GroupId,
		Seq:           req.Seq,
		IsNative:      req.IsNative,
		IsPublic:      req.IsPublic,
		IsFeatured:    req.IsFeatured,
		FeaturedIcon:  req.FeaturedIcon,
		IsBestSelling: req.IsBestSelling,
		FeatureSeq:    req.FeatureSeq,
		Description:   req.Description,
		SystemCode:    req.SystemCode,
		IsGroup:       req.IsGroup,
		ParentId:      req.ParentId,
		WebLink:       req.WebLink,
	}
	return s.repo.Update(ctx, id, m)
}

func (s *AppServiceIconService) Delete(ctx context.Context, id int) error {
	return s.repo.DeleteSoft(ctx, id)
}

// ---- App Service Icon Group ----

func (s *AppServiceIconGroup) List(ctx context.Context) ([]domain.AppServiceIconGroup, error) {
	return s.repo.List(ctx)
}

func (s *AppServiceIconGroup) ListGroupsWithIcons(ctx context.Context) ([]domain.AppServiceIconGroup, error) {
	return s.repo.ListGroupsWithIcons(ctx)
}

func (s *AppServiceIconGroup) Create(ctx context.Context, req dto.AppServiceIconGroupDto) error {
	m := domain.AppServiceIconGroup{
		Name:     req.Name,
		NameEn:   req.NameEn,
		Icon:     req.Icon,
		TypeName: req.TypeName,
		Seq:      req.Seq,
	}
	return s.repo.Create(ctx, m)
}

func (s *AppServiceIconGroup) Update(ctx context.Context, id int, req dto.AppServiceIconGroupDto) error {
	m := domain.AppServiceIconGroup{
		Name:     req.Name,
		NameEn:   req.NameEn,
		Icon:     req.Icon,
		TypeName: req.TypeName,
		Seq:      req.Seq,
	}
	return s.repo.Update(ctx, id, m)
}

func (s *AppServiceIconGroup) Delete(ctx context.Context, id int) error {
	return s.repo.DeleteSoft(ctx, id)
}
