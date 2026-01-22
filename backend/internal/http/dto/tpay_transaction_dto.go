// Package dto provides implementation for dto
//
// File: tpay_transaction_dto.go
// Description: implementation for dto
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package dto

// QRPayRequest represents QR code payment request
type QRPayRequest struct {
	QRString        string  `json:"qr_string" validate:"required"`
	PaymentMethodID int32   `json:"payment_method_id" validate:"required"` // 1=WALLET, 2=CARD, etc.
	CardID          *int64  `json:"card_id,omitempty"`                     // Required if payment_method_id is CARD
	Description     *string `json:"description,omitempty"`
}

// TransactionResponse represents transaction response
type TransactionResponse struct {
	ID                int64   `json:"id"`
	ReferenceNumber   string  `json:"reference_number"`
	Amount            float64 `json:"amount"`
	FeeAmount         float64 `json:"fee_amount"`
	SenderAccountID   int64   `json:"sender_account_id"`
	ReceiverAccountID int64   `json:"receiver_account_id"`
	TransactionType   string  `json:"transaction_type"`
	TransactionStatus string  `json:"transaction_status"`
	Description       *string `json:"description,omitempty"`
	CreatedAt         string  `json:"created_at"`
}

// P2PTransferRequest represents a peer-to-peer transfer request
type P2PTransferRequest struct {
	ReceiverPhone     string  `json:"receiver_phone,omitempty" validate:"omitempty,min=8,max=15"`
	ReceiverAccountID *int64  `json:"receiver_account_id,omitempty" validate:"omitempty,gt=0"`
	Amount            float64 `json:"amount" validate:"required,gt=0"`
	Description       *string `json:"description,omitempty" validate:"omitempty,max=500"`
	PIN               string  `json:"pin" validate:"required,len=4"`
}

// P2PTransferResponse represents the response after a successful P2P transfer
type P2PTransferResponse struct {
	TransactionID       int64   `json:"transaction_id"`
	ReferenceNumber     string  `json:"reference_number"`
	SenderAccountID     int64   `json:"sender_account_id"`
	ReceiverAccountID   int64   `json:"receiver_account_id"`
	ReceiverName        string  `json:"receiver_name"`
	ReceiverPhone       string  `json:"receiver_phone"`
	Amount              float64 `json:"amount"`
	FeeAmount           float64 `json:"fee_amount"`
	Description         *string `json:"description,omitempty"`
	SenderBalanceBefore float64 `json:"sender_balance_before"`
	SenderBalanceAfter  float64 `json:"sender_balance_after"`
	Status              string  `json:"status"`
	CreatedAt           string  `json:"created_at"`
}
