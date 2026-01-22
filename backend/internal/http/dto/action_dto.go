// Package dto provides implementation for dto
//
// File: action_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "git.gerege.mn/backend-packages/common"

// Query: /actions?search=...&page=1&size=20&sort=code:asc,name:desc
type ActionQuery struct {
	common.PaginationQuery
	Search string `query:"search"`
	Sort   string `query:"sort"`
}

type ActionCreateDto struct {
	Code        string `json:"code"        validate:"required"`
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type ActionUpdateDto ActionCreateDto

