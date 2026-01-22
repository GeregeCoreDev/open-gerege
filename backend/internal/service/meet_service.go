// Package service provides implementation for service
//
// File: meet_service.go
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

type MeetService struct {
	cfg    *config.Config
	client *httpx.Client // Reusable HTTP client
}

func NewMeetService(cfg *config.Config) *MeetService {
	return &MeetService{
		cfg:    cfg,
		client: httpx.New(3 * time.Second), // Create once, reuse everywhere
	}
}

type tokenRes struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func (s *MeetService) GenerateToken(uctx context.Context) (*common.APIResponse[tokenRes], error) {
	baseURL := fmt.Sprintf("%s/livekit/token", s.cfg.URLS.Meet)

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

	client := s.client
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	resp, _, err := httpx.GetJSON[*common.APIResponse[tokenRes]](
		uctx,
		client,
		fullURL,
		map[string]string{
			fiber.HeaderCookie: "sid=" + sid,
		},
	)

	return resp, err
}

func (s *MeetService) List(uctx context.Context) (common.APIResponse[[]dto.RoomResponse], error) {
	url := fmt.Sprintf("%s/room", s.cfg.URLS.Meet)
	client := s.client
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	rooms, _, err := httpx.GetJSON[common.APIResponse[[]dto.RoomResponse]](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	})

	return rooms, err
}

func (s *MeetService) Create(uctx context.Context, req *dto.CreateRoomRequest) (*common.APIResponse[dto.RoomResponse], error) {
	url := fmt.Sprintf("%s/room", s.cfg.URLS.Meet)
	client := s.client
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	rooms, _, err := httpx.PostJSON[*dto.CreateRoomRequest, *common.APIResponse[dto.RoomResponse]](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, req)

	return rooms, err
}

func (s *MeetService) Join(uctx context.Context, req *dto.JoinRoomRequest) (*common.APIResponse[dto.RoomResponse], error) {
	url := fmt.Sprintf("%s/room/join", s.cfg.URLS.Meet)
	client := s.client
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	rooms, _, err := httpx.PostJSON[*dto.JoinRoomRequest, *common.APIResponse[dto.RoomResponse]](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, req)

	return rooms, err
}

func (s *MeetService) Delete(uctx context.Context, roomID int) error {
	url := fmt.Sprintf("%s/room/%d", s.cfg.URLS.Meet, roomID)
	client := s.client
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	_, _, err := httpx.DeleteJSON[any, any](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, nil)

	return err
}

func (s *MeetService) AddUsers(uctx context.Context, roomID int, req *dto.AddUsersRequest) error {
	url := fmt.Sprintf("%s/room/%d/users", s.cfg.URLS.Meet, roomID)
	client := s.client
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	_, _, err := httpx.PostJSON[*dto.AddUsersRequest, any](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, req)

	return err
}

func (s *MeetService) RemoveUser(uctx context.Context, roomID, userID int) error {
	url := fmt.Sprintf("%s/room/%d/users/%d", s.cfg.URLS.Meet, roomID, userID)
	client := s.client
	sid, _ := ctx.GetValue[string](uctx, ctx.KeySID)

	_, _, err := httpx.DeleteJSON[any, any](uctx, client, url, map[string]string{
		fiber.HeaderCookie: "sid=" + sid,
	}, nil)

	return err
}
