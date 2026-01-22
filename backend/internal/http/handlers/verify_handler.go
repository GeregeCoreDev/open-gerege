// Package handlers provides implementation for handlers
//
// File: verify_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"fmt"
	"templatev25/internal/app"
	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/sso-client"
	"git.gerege.mn/backend-packages/resp"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

type VerifyHandler struct {
	*app.Dependencies
}

func NewVerifyHandler(d *app.Dependencies) *VerifyHandler {
	return &VerifyHandler{Dependencies: d}
}

func (h *VerifyHandler) Dan(c *fiber.Ctx) error {
	sid := c.Locals(ssoclient.LocalsSID).(string)

	claims, err := h.SSO.GetClaims(c.Context(), sid, ctx.RequestID(c))
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	cfg := &oauth2.Config{
		ClientID:     h.Cfg.Auth.ClientID,
		ClientSecret: h.Cfg.Auth.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  h.Cfg.Auth.DanURL,
			TokenURL: h.Cfg.Auth.TokenURL,
		},
	}

	authURL := cfg.AuthCodeURL(
		"xyz",
		oauth2.SetAuthURLParam("sid", sid),
		oauth2.SetAuthURLParam("uid", fmt.Sprintf("%d", claims.UserID)),
		oauth2.SetAuthURLParam("client", "android"),
		oauth2.SetAuthURLParam("redirect_url", h.Cfg.URLS.SSO+"/auth/dan/callback/v3"),
		oauth2.SetAuthURLParam("callback_url", h.Cfg.URLS.SSO+"/"),
	)

	return resp.OK(c, fiber.Map{"url": authURL})
}

func (h *VerifyHandler) Email(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[struct {
		Email string `json:"email" validate:"required"`
	}](c)
	if !ok {
		return nil
	}

	err := h.Service.Verify.EmailVerify(c.UserContext(), req.Email)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c)
}

func (h *VerifyHandler) EmailConfirm(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[struct {
		Email string `json:"email" validate:"required"`
		Code  string `json:"code"  validate:"required,len=6"`
	}](c)
	if !ok {
		return nil
	}

	err := h.Service.Verify.EmailVerifyConfirm(c.UserContext(), req.Email, req.Code)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c)
}

func (h *VerifyHandler) Phone(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[struct {
		PhoneNo string `json:"phone_no" validate:"required"`
	}](c)
	if !ok {
		return nil
	}

	err := h.Service.Verify.PhoneVerify(c.UserContext(), req.PhoneNo)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c)
}

func (h *VerifyHandler) PhoneConfirm(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[struct {
		Phone string `json:"phone_no" validate:"required"`
		Code  string `json:"code"  validate:"required,len=6"`
	}](c)
	if !ok {
		return nil
	}

	err := h.Service.Verify.PhoneVerifyConfirm(c.UserContext(), req.Phone, req.Code)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c)
}
