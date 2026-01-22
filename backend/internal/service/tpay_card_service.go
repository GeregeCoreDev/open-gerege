// Package service provides implementation for service
//
// File: tpay_card_service.go
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

type CardService struct {
	cfg *config.Config
}

func NewCardService(cfg *config.Config) *CardService {
	return &CardService{
		cfg: cfg,
	}
}

func (s *CardService) List(uctx context.Context) (*common.APIResponse[[]dto.Card], error) {
	url := fmt.Sprintf("%s/card/list", s.cfg.URLS.Tpay)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	resp, _, err := httpx.GetJSON[common.APIResponse[[]dto.Card]](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	})

	return &resp, err
}

func (s *CardService) Create(uctx context.Context, req *dto.CreateCardDto) (*common.APIResponse[dto.Order], error) {
	url := fmt.Sprintf("%s/card/create", s.cfg.URLS.Tpay)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	resp, _, err := httpx.PostJSON[*dto.CreateCardDto, *common.APIResponse[dto.Order]](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, req)

	return resp, err
}

func (s *CardService) Confirm(uctx context.Context, req *dto.ConfirmCardReq) (*common.APIResponse[dto.Card], error) {
	url := fmt.Sprintf("%s/card/create", s.cfg.URLS.Tpay)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	resp, _, err := httpx.PostJSON[*dto.ConfirmCardReq, *common.APIResponse[dto.Card]](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, req)

	return resp, err
}

func (s *CardService) SendOtp(uctx context.Context, id int) error {
	url := fmt.Sprintf("%s/card/send_otp?id=%d", s.cfg.URLS.Tpay, id)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	_, _, err := httpx.GetJSON[any](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	})

	return err
}

func (s *CardService) VerifyCard(uctx context.Context, req *dto.ReqVerifyCard) error {
	url := fmt.Sprintf("%s/card/verify_card", s.cfg.URLS.Tpay)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	_, _, err := httpx.PostJSON[*dto.ReqVerifyCard, any](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, req)

	return err
}
