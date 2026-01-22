// Package dto provides implementation for dto
//
// File: system_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "git.gerege.mn/backend-packages/common"

type SystemListQuery struct {
	Code     string `query:"code"`
	Key      string `query:"key"`
	Name     string `query:"name"`
	IsActive *bool  `query:"is_active"`
	common.PaginationQuery
}

type SystemCreateDto struct {
	Code        string `json:"code"        validate:"required,max=255"`
	Key         string `json:"key"         validate:"omitempty,max=255"`
	Name        string `json:"name"        validate:"required,max=255"`
	Description string `json:"description" validate:"omitempty,max=255"`
	IsActive    *bool  `json:"is_active,omitempty"`
	Icon        string `json:"icon"        validate:"omitempty,max=255"`
	Path        string `json:"path"        validate:"omitempty,max=255"`
	Sequence    int    `json:"sequence"`
}

type SystemUpdateDto SystemCreateDto
