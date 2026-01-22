// Package dto provides implementation for dto
//
// File: tpay_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

import "time"

type Account struct {
	Id            int64      `json:"id"`
	AccountTypeId int32      `json:"account_type_id"`
	OwnerId       int64      `json:"owner_id"`
	Balance       float64    `json:"balance"`
	FreezeAmount  float64    `json:"freeze_amount"`
	IsActive      bool       `json:"is_active"`
	IsDefault     bool       `json:"is_default"`
	CreatedAt     *time.Time `json:"created_date,omitempty"`
	UpdatedAt     *time.Time `json:"updated_date,omitempty"`

	AccountType *AccountType `json:"account_type,omitempty"`
}

type AccountType struct {
	Id          int32      `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_date,omitempty"`
	UpdatedAt   *time.Time `json:"updated_date,omitempty"`
}

type SetDefaultAccountRequest struct {
	AccountID int64 `json:"account_id" validate:"required,gt=0"`
}

type AccountStatementRequest struct {
	AccountID   int64      `json:"account_id" validate:"required,gt=0"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	EntrySide   *string    `json:"entry_side" validate:"omitempty,oneof=DEBIT CREDIT"`
	EntryTypeID *int32     `json:"entry_type_id"`
	Page        int        `json:"page" validate:"omitempty,min=1"`
	PageSize    int        `json:"page_size" validate:"omitempty,min=1,max=100"`
}

type AccountStatementEntry struct {
	ID              int64   `json:"id"`
	TransactionID   int64   `json:"transaction_id"`
	TransactionDate string  `json:"transaction_date"`
	EntryType       string  `json:"entry_type"`
	EntrySide       string  `json:"entry_side"`
	Amount          float64 `json:"amount"`
	BalanceBefore   float64 `json:"balance_before"`
	BalanceAfter    float64 `json:"balance_after"`
	PaymentMethod   *string `json:"payment_method,omitempty"`
	Description     *string `json:"description,omitempty"`
	ReferenceNumber string  `json:"reference_number"`
}

type AccountStatementResponse struct {
	AccountID      int64                   `json:"account_id"`
	CurrentBalance float64                 `json:"current_balance"`
	TotalEntries   int64                   `json:"total_entries"`
	Page           int                     `json:"page"`
	PageSize       int                     `json:"page_size"`
	TotalPages     int                     `json:"total_pages"`
	Items          []AccountStatementEntry `json:"items"`
}

type AccountQR struct {
	Id        int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AccountId int64      `gorm:"column:account_id;not null" json:"account_id"`
	QRData    string     `gorm:"column:qr_data;type:text;not null;uniqueIndex" json:"qr_data"`
	Amount    float64    `gorm:"column:amount;type:numeric(19,4);default:0.0000" json:"amount"`
	IsActive  bool       `gorm:"column:is_active;default:true" json:"is_active"`
	ExpiresAt *time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	CreatedAt *time.Time `json:"created_date,omitempty"`
	UpdatedAt *time.Time `json:"updated_date,omitempty"`

	// Relations
	Account *Account `gorm:"foreignKey:AccountId;references:Id" json:"account,omitempty"`
}

type AccountQRGenerateRequest struct {
	Amount float64 `json:"amount" validate:"gte=0"`
}

type AccountQRResponse struct {
	Id           int64     `json:"id"`
	AccountId    int64     `json:"account_id"`
	QRData       string    `json:"qr_data"`
	Amount       float64   `json:"amount"`
	IsActive     bool      `json:"is_active"`
	ExpiresAt    time.Time `json:"expires_at"`
	RemainingSec int       `json:"remaining_second"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	IsExpired    bool      `json:"is_expired"`
	IsValid      bool      `json:"is_valid"`
}
