// Package dto provides implementation for dto
//
// File: user_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

type UserCreateDto struct {
	Id         int    `json:"id"         validate:"required,gt=0"` // хуучин логикоор Id-тайгаар орж ирдэг
	CivilId    int    `json:"civil_id"`
	RegNo      string `json:"reg_no"      validate:"omitempty,max=10"`
	FamilyName string `json:"family_name" validate:"omitempty,max=80"`
	LastName   string `json:"last_name"   validate:"omitempty,max=150"`
	FirstName  string `json:"first_name"  validate:"omitempty,max=150"`
	Gender     int    `json:"gender"`
	BirthDate  string `json:"birth_date"  validate:"omitempty,max=10"`
	PhoneNo    string `json:"phone_no"    validate:"omitempty,max=8"`
	Email      string `json:"email"       validate:"omitempty,max=80,email"`
}

type UserUpdateDto struct {
	Id         int    `json:"id"         validate:"required,gt=0"`
	CivilId    int    `json:"civil_id"`
	RegNo      string `json:"reg_no"      validate:"omitempty,max=10"`
	FamilyName string `json:"family_name" validate:"omitempty,max=80"`
	LastName   string `json:"last_name"   validate:"omitempty,max=150"`
	FirstName  string `json:"first_name"  validate:"omitempty,max=150"`
	Gender     int    `json:"gender"`
	BirthDate  string `json:"birth_date"  validate:"omitempty,max=10"`
	PhoneNo    string `json:"phone_no"    validate:"omitempty,max=8"`
	Email      string `json:"email"       validate:"omitempty,max=80,email"`
}

// Core-оос хайх хүсэлт (хуучин models.ReqFind-тэй адилхан талбар)
type ReqFind struct {
	SearchText string `json:"search_text" validate:"required"`
}

type AccountInfo struct {
	UserID   uint   `json:"user_id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageURL string `json:"image_url"`
}

type CitizenVerification struct {
	CitizenID        uint `json:"citizen_id"`
	IsRegNoVerified  bool `json:"is_reg_no_verified"`
	IsDanVerified    bool `json:"is_dan_verified"`
	IsGsignMVerified bool `json:"is_gsign_m_verified"`
	IsGsignCVerified bool `json:"is_gsign_c_verified"`
	IsPhoneVerified  bool `json:"is_phone_verified"`
	IsEmailVerified  bool `json:"is_email_verified"`
}

type UserProfileInfo struct {
	LoginAccountInfo AccountInfo          `json:"login_account_info"`
	CitizenInfo      *UserCreateDto       `json:"citizen_info,omitempty"`
	Verifications    *CitizenVerification `json:"verifications,omitempty"`
}

// Core response struct (шаардлагатай талбараа нэмээрэй)
type CoreUser struct {
	Id         int    `json:"id"`
	CivilId    int    `json:"civil_id"`
	RegNo      string `json:"reg_no"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	FamilyName string `json:"family_name"`
	PhoneNo    string `json:"phone_no"`
	Email      string `json:"email"`
	BirthDate  string `json:"birth_date"`
	Gender     int    `json:"gender"`
}
