// Package service provides implementation for service
//
// File: menu_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"
	"sort"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"templatev25/internal/repository"
)

type MenuService interface {
	List(ctx context.Context, q dto.MenuListQuery) ([]domain.Menu, int64, int, int, error)
	ListAll(ctx context.Context) ([]domain.Menu, error)
	ListByUserRoles(ctx context.Context, userID int) ([]domain.Menu, error)
	ByID(ctx context.Context, id int64) (domain.Menu, error)
	Create(ctx context.Context, req dto.MenuCreateDto) error
	Update(ctx context.Context, id int64, req dto.MenuUpdateDto) error
	Delete(ctx context.Context, id int64) error
}

type menuService struct {
	repo repository.MenuRepository
}

func NewMenuService(repo repository.MenuRepository) MenuService {
	return &menuService{repo: repo}
}

func (s *menuService) List(ctx context.Context, q dto.MenuListQuery) ([]domain.Menu, int64, int, int, error) {
	return s.repo.List(ctx, q)
}

func (s *menuService) ListAll(ctx context.Context) ([]domain.Menu, error) {
	return s.repo.ListAll(ctx)
}

func (s *menuService) ListByUserRoles(ctx context.Context, userID int) ([]domain.Menu, error) {
	// Get menus by user roles
	allMenus, err := s.repo.ListByUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(allMenus) == 0 {
		return []domain.Menu{}, nil
	}

	// Collect all parent IDs recursively until we reach root
	allParentIDs := make(map[int64]bool)
	processedIDs := make(map[int64]bool)

	// Initialize with permission menus
	for _, menu := range allMenus {
		processedIDs[menu.ID] = true
		if menu.ParentID != nil && !processedIDs[*menu.ParentID] {
			allParentIDs[*menu.ParentID] = true
		}
	}

	// Recursively fetch all parent menus until we reach root
	var parentMenus []domain.Menu
	currentLevelParentIDs := make(map[int64]bool)
	for id := range allParentIDs {
		currentLevelParentIDs[id] = true
	}

	for len(currentLevelParentIDs) > 0 {
		ids := make([]int64, 0, len(currentLevelParentIDs))
		for id := range currentLevelParentIDs {
			ids = append(ids, id)
		}

		currentParents, err := s.repo.GetMenusByIDs(ctx, ids)
		if err != nil {
			return nil, err
		}

		parentMenus = append(parentMenus, currentParents...)

		// Collect next level parent IDs
		nextLevelParentIDs := make(map[int64]bool)
		for _, menu := range currentParents {
			processedIDs[menu.ID] = true
			if menu.ParentID != nil && !allParentIDs[*menu.ParentID] && !processedIDs[*menu.ParentID] {
				allParentIDs[*menu.ParentID] = true
				nextLevelParentIDs[*menu.ParentID] = true
			}
		}

		currentLevelParentIDs = nextLevelParentIDs
	}

	// Combine all menus using map to avoid duplicates
	allMenusMap := make(map[int64]domain.Menu, len(allMenus)+len(parentMenus))
	for _, menu := range allMenus {
		menu.Children = nil
		allMenusMap[menu.ID] = menu
	}
	for _, menu := range parentMenus {
		if _, exists := allMenusMap[menu.ID]; !exists {
			menu.Children = nil
			allMenusMap[menu.ID] = menu
		}
	}

	// Convert to slice and sort - O(n log n) instead of O(n²)
	combinedMenus := make([]domain.Menu, 0, len(allMenusMap))
	for _, menu := range allMenusMap {
		combinedMenus = append(combinedMenus, menu)
	}
	sort.Slice(combinedMenus, func(i, j int) bool {
		if combinedMenus[i].Sequence != combinedMenus[j].Sequence {
			return combinedMenus[i].Sequence < combinedMenus[j].Sequence
		}
		return combinedMenus[i].ID < combinedMenus[j].ID
	})

	// Build tree in single pass using map of pointers
	menuMap := make(map[int64]*domain.Menu, len(combinedMenus))
	for i := range combinedMenus {
		combinedMenus[i].Children = []domain.Menu{}
		menuMap[combinedMenus[i].ID] = &combinedMenus[i]
	}

	// Attach children to parents (single pass, O(n))
	var rootMenus []*domain.Menu
	for i := range combinedMenus {
		menu := &combinedMenus[i]
		if menu.ParentID == nil {
			rootMenus = append(rootMenus, menu)
		} else if parent, exists := menuMap[*menu.ParentID]; exists {
			parent.Children = append(parent.Children, *menu)
		}
	}

	// Sort children recursively using sort.Slice - O(n log n)
	var sortChildren func(children []domain.Menu)
	sortChildren = func(children []domain.Menu) {
		sort.Slice(children, func(i, j int) bool {
			if children[i].Sequence != children[j].Sequence {
				return children[i].Sequence < children[j].Sequence
			}
			return children[i].ID < children[j].ID
		})
		for i := range children {
			if len(children[i].Children) > 0 {
				sortChildren(children[i].Children)
			}
		}
	}

	// Sort root menus using sort.Slice
	sort.Slice(rootMenus, func(i, j int) bool {
		if rootMenus[i].Sequence != rootMenus[j].Sequence {
			return rootMenus[i].Sequence < rootMenus[j].Sequence
		}
		return rootMenus[i].ID < rootMenus[j].ID
	})

	// Sort all children recursively
	for _, menu := range rootMenus {
		sortChildren(menu.Children)
	}

	// Convert to result slice
	result := make([]domain.Menu, len(rootMenus))
	for i, menu := range rootMenus {
		result[i] = *menu
	}

	return result, nil
}

func (s *menuService) ByID(ctx context.Context, id int64) (domain.Menu, error) {
	return s.repo.ByID(ctx, id)
}

func (s *menuService) Create(ctx context.Context, req dto.MenuCreateDto) error {
	// ParentID 0 байвал nil болгох
	parentID := req.ParentID
	if parentID != nil && *parentID == 0 {
		parentID = nil
	}

	// PermissionID 0 байвал nil болгох
	permissionID := req.PermissionID
	if permissionID != nil && *permissionID == 0 {
		permissionID = nil
	}

	m := domain.Menu{
		Code:         req.Code,
		Key:          req.Key,
		Name:         req.Name,
		Description:  req.Description,
		Icon:         req.Icon,
		Path:         req.Path,
		Sequence:     req.Sequence,
		ParentID:     parentID,
		PermissionID: permissionID,
		IsActive:     req.IsActive,
	}
	return s.repo.Create(ctx, m)
}

func (s *menuService) Update(ctx context.Context, id int64, req dto.MenuUpdateDto) error {
	// ParentID 0 байвал nil болгох
	parentID := req.ParentID
	if parentID != nil && *parentID == 0 {
		parentID = nil
	}

	// PermissionID 0 байвал nil болгох
	permissionID := req.PermissionID
	if permissionID != nil && *permissionID == 0 {
		permissionID = nil
	}

	m := domain.Menu{
		Code:         req.Code,
		Key:          req.Key,
		Name:         req.Name,
		Description:  req.Description,
		Icon:         req.Icon,
		Path:         req.Path,
		Sequence:     req.Sequence,
		ParentID:     parentID,
		PermissionID: permissionID,
		IsActive:     req.IsActive,
	}
	return s.repo.Update(ctx, id, m)
}

func (s *menuService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
