// Package dto provides implementation for dto
//
// File: permission_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "git.gerege.mn/backend-packages/common"

// Query: /permissions?search=...&module_id=...&page=1&size=20&sort=code:asc,name:desc
type PermissionQuery struct {
	common.PaginationQuery
	SystemID int    `query:"system_id"`
	ModuleID int    `query:"module_id"`
	Search   string `query:"search"`
	Sort     string `query:"sort"`
}

type PermissionCreateDto struct {
	SystemID  int     `json:"system_id"   validate:"required,gt=0"`
	ModuleID  int     `json:"module_id"   validate:"required,gt=0"`
	ActionIDs []int64 `json:"action_ids" validate:"required,min=1,dive,gt=0"`
}

type PermissionUpdateDto struct {
	Code        string `json:"code"        validate:"required"`
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description"`
	ModuleID    int    `json:"module_id"   validate:"required,gt=0"`
	SystemID    int    `json:"system_id"   validate:"required,gt=0"`
	ActionID    *int64 `json:"action_id"`
	IsActive    *bool  `json:"is_active"`
}
