// Package handlers provides implementation for handlers
//
// File: tpay_account_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"templatev25/internal/http/dto"

	"context"
	"strconv"

	"git.gerege.mn/backend-packages/config"

	"templatev25/internal/service"

	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/resp"

	"github.com/gofiber/fiber/v2"
)

type tpayAccountHandler struct {
	cfg *config.Config
	svc *service.AccountService
}

func newTpayAccountHandler(cfg *config.Config, svc *service.AccountService) *tpayAccountHandler {
	return &tpayAccountHandler{
		cfg: cfg,
		svc: svc,
	}
}

// GetMyAccounts godoc
// @Summary      Get user's Tpay accounts
// @Tags         me
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /me/accounts [get]
func (h *tpayAccountHandler) GetMyAccounts(c *fiber.Ctx) error {
	response, err := h.svc.GetMyAccount(c.UserContext())
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, response.Data)
}

// SetDefaultAccount godoc
// @Summary      Set default Tpay account
// @Tags         me
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.SetDefaultAccountRequest true "Account data"
// @Success      200 {object} map[string]interface{}
// @Router       /me/accounts/default [put]
func (h *tpayAccountHandler) SetDefaultAccount(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.SetDefaultAccountRequest](c)
	if !ok {
		return nil
	}

	err := h.svc.SetDefaultAccount(c.UserContext(), &req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c)
}

// GetStatement godoc
// @Summary      Get Tpay account statement
// @Tags         me
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /me/accounts/statement [get]
func (h *tpayAccountHandler) GetStatement(c *fiber.Ctx) error {
	uctx := c.UserContext()

	// Query map авах
	queries := c.Queries()

	// Type-safe key ашиглах
	uctx = context.WithValue(uctx, ctx.KeyQueries, queries)

	res, err := h.svc.GetStatement(uctx)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, res.Data)
}

// GenerateQR godoc
// @Summary      Generate QR code for Tpay account
// @Tags         me
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        account_id path int true "Account ID"
// @Param        body body dto.AccountQRGenerateRequest true "QR data"
// @Success      200 {object} map[string]interface{}
// @Router       /me/accounts/{account_id}/qr [post]
func (h *tpayAccountHandler) GenerateQR(c *fiber.Ctx) error {
	accountIDStr := c.Params("account_id")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil || accountID <= 0 {
		return resp.BadRequest(c, "invalid account_id", nil)
	}

	req, ok := resp.BodyBindAndValidate[dto.AccountQRGenerateRequest](c)
	if !ok {
		return nil
	}

	res, err := h.svc.GenerateQR(c.UserContext(), accountID, &req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, res.Data)
}
