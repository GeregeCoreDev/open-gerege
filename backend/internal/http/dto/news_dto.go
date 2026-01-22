// Package dto provides implementation for dto
//
// File: news_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "git.gerege.mn/backend-packages/common"

type NewsListQuery struct {
	CategoryID int `query:"category_id"`
	common.PaginationQuery
}

type NewsDto struct {
	Title    string `json:"title"     validate:"required,min=3,max=255"`
	Text     string `json:"text"      validate:"required,min=3"`
	ImageUrl string `json:"image_url" validate:"omitempty,min=3,max=255"`
}
