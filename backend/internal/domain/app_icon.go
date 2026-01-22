// Package domain provides implementation for domain
//
// File: platform_icon.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type AppServiceIcon struct {
	Id            uint                 `json:"id" gorm:"primaryKey"`
	Name          string               `json:"name" gorm:"type:varchar(255)" validate:"required"`
	NameEn        string               `json:"name_en" gorm:"type:varchar(255)"`
	Icon          string               `json:"icon" gorm:"type:varchar(255)"`
	IconApp       string               `json:"icon_app" gorm:"type:varchar(255)"`
	IconTablet    string               `json:"icon_tablet" gorm:"type:varchar(255)"`
	IconKiosk     string               `json:"icon_kiosk" gorm:"type:varchar(255)"`
	Link          string               `json:"link" gorm:"type:varchar(255)"`
	GroupId       uint                 `json:"group_id" validate:"required"`
	Group         *AppServiceIconGroup `json:"group,omitempty" gorm:"foreignKey:Id;references:GroupId"`
	Seq           uint                 `json:"seq" gorm:"default:1"`
	IsNative      *bool                `json:"is_native" gorm:"default:false"`
	IsPublic      *bool                `json:"is_public" gorm:"default:true"`
	IsFeatured    *bool                `json:"is_featured" gorm:"default:false"`
	FeaturedIcon  string               `json:"featured_icon" gorm:"type:varchar(255)"`
	IsBestSelling *bool                `json:"is_best_selling" gorm:"default:false"`
	FeatureSeq    uint                 `json:"feature_seq" gorm:"default:1"`
	Description   string               `json:"description" gorm:"type:text"`
	SystemCode    string               `json:"system_code" gorm:"type:varchar(2)"`
	IsGroup       *bool                `json:"is_group"`
	ParentId      uint                 `json:"parent_id"`
	Parent        *AppServiceIcon      `json:"parent" gorm:"foreignKey:ParentId"`
	Childs        []AppServiceIcon     `json:"childs" gorm:"foreignKey:ParentId"`
	WebLink       string               `json:"web_link" gorm:"type:varchar(255)"`
	ExtraFields
}
type AppServiceIconGroup struct {
	Id          uint             `json:"id" gorm:"primaryKey"`
	Name        string           `json:"name" gorm:"type:varchar(80)" validate:"required"`
	NameEn      string           `json:"name_en" gorm:"type:varchar(80)"`
	Icon        string           `json:"icon" gorm:"type:varchar(255)"`
	AppServices []AppServiceIcon `json:"services" gorm:"foreignKey:GroupId;references:Id"`
	TypeName    string           `json:"type_name" gorm:"type:varchar(255);default:group"`
	Seq         uint             `json:"seq" gorm:"default:1"`
	ExtraFields
}
