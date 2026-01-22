// Package dto provides implementation for dto
//
// File: platform_icon_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

// ---- App Service Icon Group ----
type AppServiceIconGroupDto struct {
	Name     string `json:"name" validate:"required,max=80"`
	NameEn   string `json:"name_en" validate:"omitempty,max=80"`
	Icon     string `json:"icon,omitempty" validate:"omitempty,max=255"`
	TypeName string `json:"type_name" validate:"omitempty,max=255"`
	Seq      uint   `json:"seq"`
}

// ---- App Service Icon ----
type AppServiceIconDto struct {
	Name          string `json:"name" validate:"required,max=255"`
	NameEn        string `json:"name_en" validate:"omitempty,max=255"`
	Icon          string `json:"icon,omitempty" validate:"omitempty,max=255"`
	IconApp       string `json:"icon_app,omitempty" validate:"omitempty,max=255"`
	IconTablet    string `json:"icon_tablet,omitempty" validate:"omitempty,max=255"`
	IconKiosk     string `json:"icon_kiosk,omitempty" validate:"omitempty,max=255"`
	Link          string `json:"link,omitempty" validate:"omitempty,max=255"`
	GroupId       uint   `json:"group_id" validate:"required"`
	Seq           uint   `json:"seq"`
	IsNative      *bool  `json:"is_native"`
	IsPublic      *bool  `json:"is_public"`
	IsFeatured    *bool  `json:"is_featured"`
	FeaturedIcon  string `json:"featured_icon,omitempty" validate:"omitempty,max=255"`
	IsBestSelling *bool  `json:"is_best_selling"`
	FeatureSeq    uint   `json:"feature_seq"`
	Description   string `json:"description,omitempty"`
	SystemCode    string `json:"system_code,omitempty" validate:"omitempty,max=2"`
	IsGroup       *bool  `json:"is_group"`
	ParentId      uint   `json:"parent_id,omitempty"`
	WebLink       string `json:"web_link,omitempty" validate:"omitempty,max=255"`
}
