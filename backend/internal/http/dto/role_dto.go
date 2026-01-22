// Package dto provides implementation for dto
//
// File: role_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "git.gerege.mn/backend-packages/common"

type RoleListQuery struct {
	SystemId int   `query:"system_id" validate:"omitempty,gt=0"`
	IsActive *bool `query:"is_active"`
	common.PaginationQuery
}

type RoleCreateDto struct {
	SystemID    int    `json:"system_id" validate:"required,gt=0"`
	Code        string `json:"code"        validate:"required,min=2,max=255"`
	Name        string `json:"name"        validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"max=255"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type RoleUpdateDto RoleCreateDto

type RolePermissionsQuery struct {
	RoleID int `query:"role_id" validate:"required,gt=0"`
}

type RolePermissionsUpdateDto struct {
	RoleID        int   `json:"role_id"        validate:"required,gt=0"`
	PermissionIDs []int `json:"permission_ids" validate:"required,min=0,dive,gt=0"`
}
