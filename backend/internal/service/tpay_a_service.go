// Package service provides implementation for service
//
// File: tpay_a_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"git.gerege.mn/backend-packages/config"
)

type TpayService struct {
	cfg     *config.Config
	Account *AccountService
	Payment *PaymentService
	Card    *CardService
}

func NewTpayService(cfg *config.Config) *TpayService {
	return &TpayService{
		cfg:     cfg,
		Account: NewAccountService(cfg),
		Payment: NewPaymentService(cfg),
		Card:    NewCardService(cfg),
	}
}
