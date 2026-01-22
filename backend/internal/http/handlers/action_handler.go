// Package handlers provides implementation for handlers
//
// File: action_handler.go
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
	"git.gerege.mn/backend-packages/resp"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ActionHandler struct {
	*app.Dependencies
}

func NewActionHandler(d *app.Dependencies) *ActionHandler {
	return &ActionHandler{Dependencies: d}
}

// List godoc
// @Summary      List actions (paginated)
// @Tags         actions
// @Security     BearerAuth
// @Param        search    query   string false "Search (code/name/description)"
// @Param        module_id query   int    false "Filter by module_id"
// @Param        page      query   int    false "Page number"
// @Param        size      query   int    false "Page size"
// @Param        sort      query   string false "Sort (e.g. code:asc,name:desc)"
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /actions [get]
func (h *ActionHandler) List(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.ActionQuery](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	items, total, page, size, err := h.Service.Action.ListFilteredPaged(ctx, q)
	if err != nil {
		h.Log.Error("action_list_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Create godoc
// @Summary      Create action
// @Tags         actions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body   body dto.ActionCreateDto true "payload"
// @Success      201 {object} map[string]interface{}
// @Router       /actions [post]
func (h *ActionHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.ActionCreateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Action.Create(ctx, req); err != nil {
		h.Log.Error("action_create_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Created(c)
}

// Update godoc
// @Summary      Update action
// @Tags         actions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path int64                   true "ID"
// @Param        body   body dto.ActionUpdateDto     true "payload"
// @Success      200    {object} map[string]interface{}
// @Router       /actions/{id} [put]
func (h *ActionHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return resp.BadRequest(c, "invalid action id", err.Error())
	}

	req, ok := resp.BodyBindAndValidate[dto.ActionUpdateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Action.Update(ctx, id64, req); err != nil {
		h.Log.Warn("action_update_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Delete godoc
// @Summary      Delete action (soft)
// @Tags         actions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path int64 true "ID"
// @Success      200  {object} map[string]interface{}
// @Router       /actions/{id} [delete]
func (h *ActionHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return resp.BadRequest(c, "invalid action id", err.Error())
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Action.Delete(ctx, id64); err != nil {
		h.Log.Warn("action_delete_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

