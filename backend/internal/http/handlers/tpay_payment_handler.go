// Package handlers provides implementation for handlers
//
// File: tpay_payment_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/config"

	"templatev25/internal/service"

	"git.gerege.mn/backend-packages/resp"

	"github.com/gofiber/fiber/v2"
)

type tpayPaymentHandler struct {
	cfg *config.Config
	svc *service.PaymentService
}

func newTpayPaymentHandler(cfg *config.Config, svc *service.PaymentService) *tpayPaymentHandler {
	return &tpayPaymentHandler{
		cfg: cfg,
		svc: svc,
	}
}

// QrPay godoc
// @Summary      QR payment
// @Tags         me
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.QRPayRequest true "Payment data"
// @Success      200 {object} map[string]interface{}
// @Router       /me/tpay/transaction/qr-pay [post]
func (h *tpayPaymentHandler) QrPay(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.QRPayRequest](c)
	if !ok {
		return nil
	}

	res, err := h.svc.QrPay(c.UserContext(), &req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, res.Data)
}

// P2PTransfer godoc
// @Summary      P2P transfer
// @Tags         me
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.P2PTransferRequest true "Transfer data"
// @Success      200 {object} map[string]interface{}
// @Router       /me/tpay/transaction/p2p [post]
func (h *tpayPaymentHandler) P2PTransfer(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.P2PTransferRequest](c)
	if !ok {
		return nil
	}

	res, err := h.svc.P2PTransfer(c.UserContext(), &req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, res.Data)
}
