// Package dto provides implementation for dto
//
// File: organization_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import (
	"templatev25/internal/domain"
	"git.gerege.mn/backend-packages/common"
)

type OrganizationDto struct {
	Id                int     `json:"id" validate:"omitempty,gt=0"`
	RegNo             string  `json:"reg_no" validate:"required,max=7"`
	Name              string  `json:"name" validate:"required,max=255"`
	ShortName         string  `json:"short_name" validate:"omitempty,max=255"`
	TypeId            int     `json:"type_id" validate:"required,gt=0"`
	PhoneNo           string  `json:"phone_no" validate:"omitempty,max=8"`
	Email             string  `json:"email" validate:"omitempty,max=50,email"`
	Longitude         float64 `json:"longitude"`
	Latitude          float64 `json:"latitude"`
	IsActive          *bool   `json:"is_active"`
	AimagId           int     `json:"aimag_id"`
	SumId             int     `json:"sum_id"`
	BagId             int     `json:"bag_id"`
	AddressDetail     string  `json:"address_detail" validate:"omitempty,max=255"`
	AimagName         string  `json:"aimag_name" validate:"omitempty,max=255"`
	SumName           string  `json:"sum_name" validate:"omitempty,max=255"`
	BagName           string  `json:"bag_name" validate:"omitempty,max=255"`
	CountryCode       string  `json:"country_code"`
	CountryName       string  `json:"country_name"`
	Sequence          int     `json:"sequence"`
	ParentAddressId   int     `json:"parent_address_id"`
	ParentAddressName string  `json:"parent_address_name" validate:"omitempty,max=25"`
	CountryNameEn     string  `json:"country_name_en"`
	ParentID          *int    `json:"parent_id"`
}

type OrganizationUpdateDto = OrganizationDto

type OrganizationTreeQuery struct {
	OrgId int `query:"org_id" validate:"required"`
}

type OrganizationTypeDto struct {
	Code        string `json:"code" validate:"required,max=255"`
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description" validate:"omitempty,max=255"`
}

type OrgTypeRolesQuery struct {
	TypeID int `query:"type_id" validate:"required,gt=0"` // org type id
}

type OrgTypeRolesAddDto struct {
	TypeID  int   `json:"type_id" validate:"required,gt=0"` // org type id
	RoleIDs []int `json:"role_ids" validate:"required,dive,gt=0"`
}

type OrgUserListQuery struct {
	OrgId  int    `query:"org_id"`
	UserId int    `query:"user_id"`
	Name   string `query:"name"`
	common.PaginationQuery
}

type OrgUserCreateDto struct {
	OrgId  int `json:"org_id" validate:"required,gt=0"`
	UserId int `json:"user_id" validate:"required,gt=0"`
}

type OrgUserDeleteDto OrgUserCreateDto
type ResOrguserUserItem struct {
	OrgId       int                  `json:"org_id"`
	UserId      int                  `json:"user_id"`
	LastName    string               `json:"last_name"`
	FirstName   string               `json:"first_name"`
	RegNo       string               `json:"reg_no"`
	BirthDate   string               `json:"birth_date"`
	Gender      int                  `json:"gender"`
	PhoneNo     string               `json:"phone_no"`
	Email       string               `json:"email"`
	CreatedDate domain.LocalDateTime `json:"created_date"`
}

type ResOrguserOrgItem struct {
	OrgId       int                  `json:"org_id"`
	Id          int                  `json:"id"`
	Name        string               `json:"name"`
	ShortName   string               `json:"short_name"`
	RegNo       string               `json:"reg_no"`
	CreatedDate domain.LocalDateTime `json:"created_date"`
}

type OrgTypeSystemsQuery struct {
	TypeID int `query:"type_id" validate:"required,gt=0"`
}

type OrgTypeAddSystemsDto struct {
	TypeID    int   `json:"type_id"    validate:"required,gt=0"`
	SystemIDs []int `json:"system_ids" validate:"required,min=1,dive,gt=0"`
}
