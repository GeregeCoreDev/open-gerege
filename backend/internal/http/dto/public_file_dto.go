// Package dto provides implementation for dto
//
// File: public_file_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "git.gerege.mn/backend-packages/common"

type PublicFileListQuery struct {
	Description string `query:"description"`
	// Стандарт pagination/sort/search-аа зэрэг ашиглаж болно гэсэн санаагаар
	common.PaginationQuery
}

type PublicFileDeleteDto struct {
	Name string `json:"name" validate:"required"`
}
