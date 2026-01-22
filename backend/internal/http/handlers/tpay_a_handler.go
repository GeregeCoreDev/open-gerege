// Package handlers provides implementation for handlers
//
// File: tpay_a_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import "templatev25/internal/app"

type TpayHandler struct {
	*app.Dependencies
	Account *tpayAccountHandler
	Payment *tpayPaymentHandler
	Card    *tpayCardHandler
}

func NewTpayHandler(d *app.Dependencies) *TpayHandler {
	return &TpayHandler{
		Dependencies: d,
		Account:      newTpayAccountHandler(d.Cfg, d.Service.Tpay.Account),
		Payment:      newTpayPaymentHandler(d.Cfg, d.Service.Tpay.Payment),
		Card:         newTpayCardHandler(d.Cfg, d.Service.Tpay.Card),
	}
}
