// Package dto provides implementation for dto
//
// File: chat_item_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "git.gerege.mn/backend-packages/common"

type ChatItemCreateDto struct {
	Key    string `json:"key" validate:"required"`
	Answer string `json:"answer" validate:"required"`
}

type ChatItemUpdateDto ChatItemCreateDto

type ChatItemQuery struct {
	common.PaginationQuery
	Search string `query:"search"`
}
type ChatItemKeyDto struct {
	Key string `json:"key" validate:"required"`
}
