// Package handlers provides implementation for handlers
//
// File: chat_item_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"templatev25/internal/http/dto"

	"context"
	"strings"
	"templatev25/internal/app"
	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/resp"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ChatItemHandler struct {
	*app.Dependencies
}

func NewChatItemHandler(d *app.Dependencies) *ChatItemHandler {
	return &ChatItemHandler{Dependencies: d}
}

// GetByKey godoc
// @Summary      Get chat item by key
// @Tags         chat
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.ChatItemKeyDto true "Key"
// @Success      200 {object} map[string]interface{}
// @Router       /chat/key [post]
func (h *ChatItemHandler) GetByKey(c *fiber.Ctx) error {
	dto, ok := resp.BodyBindAndValidate[dto.ChatItemKeyDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	item, err := h.Service.ChatItem.GetByKey(ctx, strings.ToLower(dto.Key))
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, item)
}

// List godoc
// @Summary      List chat items
// @Tags         chat
// @Security     BearerAuth
// @Produce      json
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Router       /chat [get]
func (h *ChatItemHandler) List(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.ChatItemQuery](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	items, total, page, size, err := h.Service.ChatItem.List(ctx, q)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.Paginated(c, items, total, page, size)
}

// Create godoc
// @Summary      Create chat item
// @Tags         chat
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.ChatItemCreateDto true "Chat item data"
// @Success      201 {object} map[string]interface{}
// @Router       /chat [post]
func (h *ChatItemHandler) Create(c *fiber.Ctx) error {
	body, ok := resp.BodyBindAndValidate[dto.ChatItemCreateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.ChatItem.Create(ctx, body); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Created(c)
}

// Update godoc
// @Summary      Update chat item
// @Tags         chat
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path int true "Chat item ID"
// @Param        body body dto.ChatItemUpdateDto true "Chat item data"
// @Success      200 {object} map[string]interface{}
// @Router       /chat/{id} [put]
func (h *ChatItemHandler) Update(c *fiber.Ctx) error {
	param, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}

	body, ok := resp.BodyBindAndValidate[dto.ChatItemUpdateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.ChatItem.Update(ctx, param.ID, body); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Delete godoc
// @Summary      Delete chat item
// @Tags         chat
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Chat item ID"
// @Success      200 {object} map[string]interface{}
// @Router       /chat/{id} [delete]
func (h *ChatItemHandler) Delete(c *fiber.Ctx) error {
	param, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.ChatItem.Delete(ctx, param.ID); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
