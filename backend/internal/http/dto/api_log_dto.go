// Package dto provides implementation for dto
//
// File: api_log_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-01-09
// Last Updated: 2025-01-09
package dto

import "git.gerege.mn/backend-packages/common"

type APILogListQuery struct {
	Method     string `query:"method"`
	Path       string `query:"path"`
	StatusCode *int   `query:"status_code"`
	UserID     *int64 `query:"user_id"`
	OrgID      *int64 `query:"org_id"`
	IP         string `query:"ip"`
	common.PaginationQuery
}
