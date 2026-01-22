// Package handlers provides implementation for handlers
//
// File: tpay_card_handler.go
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

type tpayCardHandler struct {
	cfg *config.Config
	svc *service.CardService
}

func newTpayCardHandler(cfg *config.Config, svc *service.CardService) *tpayCardHandler {
	return &tpayCardHandler{
		cfg: cfg,
		svc: svc,
	}
}

// CardList godoc
// @Summary      Get Tpay card list
// @Tags         me
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /me/card/list [get]
func (h *tpayCardHandler) CardList(c *fiber.Ctx) error {
	res, err := h.svc.List(c.UserContext())
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, res.Data)
}

// AddCard godoc
// @Summary      Add Tpay card
// @Tags         me
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.CreateCardDto true "Card data"
// @Success      200 {object} map[string]interface{}
// @Router       /me/card/create [post]
func (h *tpayCardHandler) AddCard(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.CreateCardDto](c)
	if !ok {
		return nil
	}

	res, err := h.svc.Create(c.UserContext(), &req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, res.Data)
}

// Confirm godoc
// @Summary      Confirm Tpay card
// @Tags         me
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.ConfirmCardReq true "Confirmation data"
// @Success      200 {object} map[string]interface{}
// @Router       /me/card/confirm [post]
func (h *tpayCardHandler) Confirm(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.ConfirmCardReq](c)
	if !ok {
		return nil
	}

	res, err := h.svc.Confirm(c.UserContext(), &req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, res.Data)
}

// SendOtp godoc
// @Summary      Send OTP for Tpay card
// @Tags         me
// @Security     BearerAuth
// @Produce      json
// @Param        id query int true "Card ID"
// @Success      200 {object} map[string]interface{}
// @Router       /me/card/otp [get]
func (h *tpayCardHandler) SendOtp(c *fiber.Ctx) error {
	id := c.QueryInt("id")
	err := h.svc.SendOtp(c.UserContext(), id)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// VerifyCard godoc
// @Summary      Verify Tpay card
// @Tags         me
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.ReqVerifyCard true "Verification data"
// @Success      200 {object} map[string]interface{}
// @Router       /me/card/verify [post]
func (h *tpayCardHandler) VerifyCard(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.ReqVerifyCard](c)
	if !ok {
		return nil
	}
	err := h.svc.VerifyCard(c.UserContext(), &req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
