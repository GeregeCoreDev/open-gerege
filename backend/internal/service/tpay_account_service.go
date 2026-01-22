// Package service provides implementation for service
//
// File: tpay_account_service.go
// Description: implementation for service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service

import (
	"context"
	"fmt"
	"net/url"
	"git.gerege.mn/backend-packages/config"
	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/httpx"
	"templatev25/internal/http/dto"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AccountService struct {
	cfg *config.Config
}

func NewAccountService(cfg *config.Config) *AccountService {
	return &AccountService{cfg: cfg}
}

func (s *AccountService) GetMyAccount(uctx context.Context) (*common.APIResponse[[]dto.Account], error) {
	url := fmt.Sprintf("%s/accounts/me", s.cfg.URLS.Tpay)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	accounts, _, err := httpx.GetJSON[common.APIResponse[[]dto.Account]](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	})

	return &accounts, err
}

func (s *AccountService) SetDefaultAccount(uctx context.Context, req *dto.SetDefaultAccountRequest) error {
	url := fmt.Sprintf("%s/accounts/set-default", s.cfg.URLS.Tpay)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	_, _, err := httpx.PutJSON[*dto.SetDefaultAccountRequest, any](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, req)

	return err
}

func (s *AccountService) GetStatement(uctx context.Context) (*common.APIResponse[dto.AccountStatementResponse], error) {
	baseURL := fmt.Sprintf("%s/accounts/statement", s.cfg.URLS.Tpay)

	queries, _ := uctx.Value(ctx.KeyQueries).(map[string]string)

	params := url.Values{}
	for k, v := range queries {
		if v != "" {
			params.Add(k, v)
		}
	}

	fullURL := baseURL
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	resp, _, err := httpx.GetJSON[common.APIResponse[dto.AccountStatementResponse]](
		uctx,
		client,
		fullURL,
		map[string]string{
			fiber.HeaderCookie: "sid=" + sid,
		},
	)

	return &resp, err
}

func (s *AccountService) GenerateQR(uctx context.Context, accountID int64, req *dto.AccountQRGenerateRequest) (*common.APIResponse[dto.AccountQRResponse], error) {
	url := fmt.Sprintf("%s/accounts/%d/qr", s.cfg.URLS.Tpay, accountID)
	client := httpx.New(3 * time.Second)
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	resp, _, err := httpx.PostJSON[*dto.AccountQRGenerateRequest, common.APIResponse[dto.AccountQRResponse]](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, req)

	return &resp, err
}
