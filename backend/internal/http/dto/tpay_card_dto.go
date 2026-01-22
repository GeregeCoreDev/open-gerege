// Package dto provides implementation for dto
//
// File: tpay_card_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import (
	"templatev25/internal/domain"
	"time"
)

type Bank struct {
	Id                 int    `json:"id" gorm:"primaryKey"`
	Code               string `json:"code" gorm:"varchar(40)"`
	CodeTxt            string `json:"code_txt"`
	Name               string `json:"name" gorm:"varchar(255)"`
	Img                string `json:"img"`
	GolomtBankCode     string `json:"golomt_bank_code"`
	NegdiBankCode      string `json:"negdi_bank_code" gorm:"uniqueIndex"`
	BackgroundImgBig   string `json:"background_img_big"`
	BackgroundImgSmall string `json:"background_img_small"`
}

type Card struct {
	Id          int        `json:"id"`
	BankCode    string     `json:"bank_code"`
	Bank        Bank       `json:"bank"`
	CardNumber  string     `json:"card_number"`
	Token       int        `json:"token"`
	UserId      int        `json:"user_id"`
	CompanyName string     `json:"company_name"`
	IsDefault   bool       `json:"is_default"`
	Status      int8       `json:"status"`
	IsVerified  bool       `json:"is_verified"`
	CreatedAt   *time.Time `json:"created_date,omitempty"`
	UpdatedAt   *time.Time `json:"updated_date,omitempty"`

	// Relations
	User *domain.User `json:"user,omitempty"`
}

type CreateCardDto struct {
	CardNumber string `json:"card_number" validate:"required"`
	CardExp    string `json:"card_exp" validate:"required"`
}

type Order struct {
	Amount       float64  `json:"amount"`
	ApprovalCode string   `json:"approvalCode"`
	CheckId      string   `json:"checkid"`
	Currency     string   `json:"currency"`
	Customer     Customer `json:"customer"`
	Description  string   `json:"description"`
	Ordertype    string   `json:"ordertype"`
	Reason       string   `json:"reason"`
	Status       string   `json:"status"`
	Token        []Token  `json:"token"`
	Tranid       string   `json:"tranid"`
	Detail       string   `json:"detail"`
	MaskedPan    string   `json:"maskedpan"`
	ExpDate      string   `json:"expdate"`
	Brand        string   `json:"brand"`
}

type Customer struct {
	Customerid         string `json:"customerid"`
	Customername       string `json:"customername"`
	Customerregisterid uint   `json:"customerregisterid"`
}

type Token struct {
	Expdate   string `json:"expdate"`
	Maskedpan string `json:"maskedpan"`
	Tokenid   uint   `json:"tokenid"`
	Brand     string `json:"brand"`
	Bankname  string `json:"bankname"`
}

type ConfirmCardReq struct {
	TranId  string `json:"tranid" validate:"required"`
	CheckId string `json:"checkid" validate:"required"`
	CardCvv uint   `json:"cardcvv" validate:"required"`
}

type ReqVerifyCard struct {
	Id  uint   `json:"id"`
	Otp string `json:"otp"`
}
