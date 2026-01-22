// Package dto provides implementation for dto
//
// File: module_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "git.gerege.mn/backend-packages/common"

type ModuleListQuery struct {
	Code     string `query:"code"`
	Name     string `query:"name"`
	IsActive *bool  `query:"is_active"`
	SystemID int    `query:"system_id"`
	common.PaginationQuery
}

type ModuleCreateDto struct {
	Code        string `json:"code"        validate:"required,max=255"`
	Name        string `json:"name"        validate:"required,max=255"`
	Description string `json:"description" validate:"omitempty,max=255"`
	IsActive    *bool  `json:"is_active"`
	SystemID    int    `json:"system_id"   validate:"required"`
}

type ModuleUpdateDto ModuleCreateDto

type ModuleByRoleQuery struct {
	RoleID int `query:"role_id" validate:"required,gt=0"`
}

// Nested response: System -> Module -> Permission
type PermissionNode struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModuleNode struct {
	ID          int              `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Permissions []PermissionNode `json:"permissions"`
}

type SystemNode struct {
	ID          int          `json:"id"`
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Icon        string       `json:"icon"`
	Path        string       `json:"path"`
	Sequence    int          `json:"sequence"`
	Modules     []ModuleNode `json:"modules"`
}
