// Package service provides implementation for service
//
// File: tpay_payment_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"
	"fmt"
	"git.gerege.mn/backend-packages/config"
	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/httpx"
	"templatev25/internal/http/dto"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PaymentService struct {
	cfg *config.Config
}

func NewPaymentService(cfg *config.Config) *PaymentService {
	return &PaymentService{
		cfg: cfg,
	}
}

func (s *PaymentService) QrPay(uctx context.Context, req *dto.QRPayRequest) (*common.APIResponse[dto.TransactionResponse], error) {
	url := fmt.Sprintf("%s/transaction/qr_pay", s.cfg.URLS.Tpay)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	resp, _, err := httpx.PostJSON[*dto.QRPayRequest, *common.APIResponse[dto.TransactionResponse]](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, req)

	return resp, err
}

func (s *PaymentService) P2PTransfer(uctx context.Context, req *dto.P2PTransferRequest) (*common.APIResponse[dto.P2PTransferResponse], error) {
	url := fmt.Sprintf("%s/transaction/p2p", s.cfg.URLS.Tpay)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	resp, _, err := httpx.PostJSON[*dto.P2PTransferRequest, *common.APIResponse[dto.P2PTransferResponse]](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, req)

	return resp, err
}
