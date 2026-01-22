// Package domain provides implementation for domain
//
// File: citizen.go
// Description: implementation for domain
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package domain

type Citizen struct {
	Id                int    `json:"id" gorm:"primaryKey"`
	CivilId           int    `json:"civil_id"`
	RegNo             string `json:"reg_no" gorm:"type:varchar(10)"`
	FamilyName        string `json:"family_name" gorm:"type:varchar(80)"`
	LastName          string `json:"last_name" gorm:"type:varchar(150)"`
	FirstName         string `json:"first_name" gorm:"type:varchar(150)"`
	Gender            int    `json:"gender"`
	BirthDate         string `json:"birth_date" gorm:"type:varchar(10)"`
	PhoneNo           string `json:"phone_no" gorm:"type:varchar(8)"`
	Email             string `json:"email" gorm:"type:varchar(80)"`
	IsForeign         int    `json:"is_foreign"`
	CountryCode       string `json:"country_code" gorm:"type:varchar(3)"`
	Hash              string `json:"hash" gorm:"type:varchar(200)"`
	ParentAddressId   int    `json:"parent_address_id"`
	ParentAddressName string `json:"parent_address_name" gorm:"type:varchar(20)"`
	AimagId           int    `json:"aimag_id"`
	AimagCode         string `json:"aimag_code" gorm:"type:varchar(3)"`
	AimagName         string `json:"aimag_name" gorm:"type:varchar(255)"`
	SumId             int    `json:"sum_id"`
	SumCode           string `json:"sum_code" gorm:"type:varchar(3)"`
	SumName           string `json:"sum_name" gorm:"type:varchar(255)"`
	BagId             int    `json:"bag_id"`
	BagCode           string `json:"bag_code" gorm:"type:varchar(3)"`
	BagName           string `json:"bag_name" gorm:"type:varchar(255)"`
	AddressDetail     string `json:"address_detail" gorm:"type:varchar(255)"`
	AddressType       string `json:"address_type" gorm:"type:varchar(255)"`
	AddressTypeName   string `json:"address_type_name" gorm:"type:varchar(255)"`
	Nationality       string `json:"nationality" gorm:"type:varchar(255)"`
	CountryName       string `json:"country_name" gorm:"type:varchar(255)"`
	CountryNameEn     string `json:"country_name_en" gorm:"type:varchar(255)"`
	ProfileImgUrl     string `json:"profile_img_url" gorm:"type:varchar(255)"`
	ExtraFields
}
