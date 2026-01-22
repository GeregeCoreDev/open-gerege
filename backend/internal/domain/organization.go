// Package domain provides implementation for domain
//
// File: organization.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type OrganizationType struct {
	Id          int    `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"type:varchar(255)"`
	Name        string `json:"name" gorm:"type:varchar(255)"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	ExtraFields
}

type OrgTypeSystem struct {
	TypeId   int     `json:"type_id" gorm:"primaryKey"`
	SystemID int     `json:"system_id" gorm:"primaryKey"`
	System   *System `json:"system,omitempty" gorm:"foreignKey:SystemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ExtraFields
}

type OrgTypeRole struct {
	TypeId int   `json:"type_id" gorm:"primaryKey"`
	RoleID int   `json:"role_id" gorm:"primaryKey"`
	Role   *Role `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ExtraFields
}

type Organization struct {
	Id                int               `json:"id,omitempty" gorm:"primaryKey"`
	RegNo             string            `json:"reg_no,omitempty" gorm:"type:varchar(7)"`
	Name              string            `json:"name,omitempty" gorm:"type:varchar(255)"`
	ShortName         string            `json:"short_name,omitempty" gorm:"type:varchar(255)"`
	TypeId            int               `json:"type_id"`
	Type              *OrganizationType `json:"type,omitempty" gorm:"foreignKey:TypeId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PhoneNo           string            `json:"phone_no,omitempty" gorm:"type:varchar(8)"`
	Email             string            `json:"email,omitempty" gorm:"type:varchar(50)" `
	Longitude         float64           `json:"longitude,omitempty" gorm:"type:float8;default:106.91758628931501"`
	Latitude          float64           `json:"latitude,omitempty" gorm:"type:float8;default:47.918825014251915"`
	IsActive          *bool             `json:"is_active,omitempty"`
	AimagId           int               `json:"aimag_id,omitempty"`
	SumId             int               `json:"sum_id,omitempty"`
	BagId             int               `json:"bag_id,omitempty"`
	AddressDetail     string            `json:"address_detail,omitempty" gorm:"type:varchar(255)"`
	AimagName         string            `json:"aimag_name,omitempty" gorm:"type:varchar(255)"`
	SumName           string            `json:"sum_name,omitempty" gorm:"type:varchar(255)"`
	BagName           string            `json:"bag_name,omitempty" gorm:"type:varchar(255)"`
	CountryCode       string            `json:"country_code,omitempty"`
	CountryName       string            `json:"country_name,omitempty"`
	Sequence          int               `json:"sequence,omitempty"`
	ParentAddressId   int               `json:"parent_address_id,omitempty"`
	ParentAddressName string            `json:"parent_address_name,omitempty" gorm:"type:varchar(25)"`
	CountryNameEn     string            `json:"country_name_en,omitempty"`
	ParentId          *int              `json:"parent_id"`
	Children          *[]Organization   `json:"children,omitempty" gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ExtraFields
}

type OrganizationUser struct {
	OrgId        int           `json:"org_id"`
	UserId       int           `json:"user_id"`
	Organization *Organization `json:"organization,omitempty" gorm:"foreignKey:OrgId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User         *User         `json:"user,omitempty" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ExtraFields
}
