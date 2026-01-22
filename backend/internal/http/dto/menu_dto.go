// Package dto provides implementation for dto
//
// File: menu_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "git.gerege.mn/backend-packages/common"

type MenuListQuery struct {
	Code     string `query:"code"`
	Key      string `query:"key"`
	Name     string `query:"name"`
	ParentID *int64 `query:"parent_id"`
	common.PaginationQuery
}

type MenuCreateDto struct {
	Code         string `json:"code"         validate:"required,max=255"`
	Key          string `json:"key"         validate:"required,max=255"`
	Name         string `json:"name"        validate:"required,max=255"`
	Description  string `json:"description" validate:"omitempty,max=255"`
	Icon         string `json:"icon"        validate:"omitempty,max=255"`
	Path         string `json:"path"        validate:"omitempty,max=255"`
	Sequence     int64  `json:"sequence"`
	ParentID     *int64 `json:"parent_id"`
	PermissionID *int64 `json:"permission_id"`
	IsActive     *bool  `json:"is_active"`
}

type MenuUpdateDto MenuCreateDto
