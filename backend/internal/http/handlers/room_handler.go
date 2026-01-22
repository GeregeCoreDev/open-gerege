// Package handlers provides implementation for handlers
//
// File: room_handler.go
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
	"templatev25/internal/app"

	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/resp"

	"github.com/gofiber/fiber/v2"
)

type RoomHandler struct {
	*app.Dependencies
}

func NewRoomHandler(d *app.Dependencies) *RoomHandler {
	return &RoomHandler{Dependencies: d}
}

func (h *RoomHandler) List(c *fiber.Ctx) error {
	response, err := h.Service.Meet.List(c.UserContext())
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, response.Data)
}

func (h *RoomHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.CreateRoomRequest](c)
	if !ok {
		return nil
	}

	room, err := h.Service.Meet.Create(c.UserContext(), &req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.Created(c, room.Data)
}

func (h *RoomHandler) Join(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.JoinRoomRequest](c)
	if !ok {
		return nil
	}

	room, err := h.Service.Meet.Join(c.UserContext(), &req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.Created(c, room.Data)
}

func (h *RoomHandler) Delete(c *fiber.Ctx) error {
	roomID, err := strconv.Atoi(c.Params("id"))
	if err != nil || roomID <= 0 {
		return resp.BadRequest(c, "invalid room id", nil)
	}

	if err := h.Service.Meet.Delete(c.UserContext(), roomID); err != nil {
		return err
	}

	return resp.OK(c, fiber.Map{"message": "room deleted successfully"})
}

func (h *RoomHandler) AddUsers(c *fiber.Ctx) error {
	roomID, err := strconv.Atoi(c.Params("id"))
	if err != nil || roomID <= 0 {
		return resp.BadRequest(c, "invalid room id", nil)
	}

	req, ok := resp.BodyBindAndValidate[dto.AddUsersRequest](c)
	if !ok {
		return nil
	}

	if err := h.Service.Meet.AddUsers(c.UserContext(), roomID, &req); err != nil {
		return err
	}

	return resp.OK(c, fiber.Map{"message": "users added successfully"})
}

func (h *RoomHandler) RemoveUser(c *fiber.Ctx) error {
	roomID, err := strconv.Atoi(c.Params("id"))
	if err != nil || roomID <= 0 {
		return resp.BadRequest(c, "invalid room id", nil)
	}

	uid, err := strconv.Atoi(c.Params("user_id"))
	if err != nil || uid <= 0 {
		return resp.BadRequest(c, "invalid user id", nil)
	}

	if err := h.Service.Meet.RemoveUser(c.UserContext(), roomID, uid); err != nil {
		return err
	}

	return resp.OK(c, fiber.Map{"message": "user removed successfully"})
}

func (h *RoomHandler) GenerateToken(c *fiber.Ctx) error {
	uctx := c.UserContext()

	queries := c.Queries()
	uctx = context.WithValue(uctx, ctx.KeyQueries, queries)

	res, err := h.Service.Meet.GenerateToken(uctx)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, res.Data)
}
