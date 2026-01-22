// Package service provides implementation for service
//
// File: verify_service.go
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
	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/httpx"
	"time"

	"github.com/gofiber/fiber/v2"
)

type VerifyService struct {
	cfg *config.Config
}

func NewVerifyService(cfg *config.Config) *VerifyService {
	return &VerifyService{cfg: cfg}
}

func (s *VerifyService) EmailVerify(uctx context.Context, email string) (err error) {
	url := fmt.Sprintf("%s/citizen/email/verify", s.cfg.URLS.SSO)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	type request struct {
		Email string `json:"email"`
	}

	_, _, err = httpx.PostJSON[request, any](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, request{Email: email})

	return
}

func (s *VerifyService) EmailVerifyConfirm(uctx context.Context, email, code string) (err error) {
	url := fmt.Sprintf("%s/citizen/email/verify/confirm", s.cfg.URLS.SSO)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	type request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	_, _, err = httpx.PostJSON[request, any](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, request{Email: email, Code: code})

	return
}

func (s *VerifyService) PhoneVerify(uctx context.Context, phone string) (err error) {
	url := fmt.Sprintf("%s/citizen/phone/verify", s.cfg.URLS.SSO)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	type request struct {
		PhoneNo string `json:"phone_no"`
	}

	_, _, err = httpx.PostJSON[request, any](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, request{PhoneNo: phone})

	return
}

func (s *VerifyService) PhoneVerifyConfirm(uctx context.Context, phone, code string) (err error) {
	url := fmt.Sprintf("%s/citizen/phone/verify/confirm", s.cfg.URLS.SSO)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	type request struct {
		PhoneNo string `json:"phone_no"`
		Code    string `json:"code"`
	}

	_, _, err = httpx.PostJSON[request, any](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, request{PhoneNo: phone, Code: code})

	return
}
