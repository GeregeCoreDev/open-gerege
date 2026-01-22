// Package handlers provides implementation for handlers
//
// File: system_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"templatev25/internal/http/dto"

	"context"
	"templatev25/internal/app"
	"time"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/resp"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type SystemHandler struct {
	*app.Dependencies
}

func NewSystemHandler(d *app.Dependencies) *SystemHandler {
	return &SystemHandler{Dependencies: d}
}

// GET /system
// @Summary      List systems (paginated)
// @Tags         systems
// @Security     BearerAuth
// @Param        page query int    false "Page number"
// @Param        size query int    false "Page size"
// @Param        code query string false "Filter by code (ILIKE)"
// @Param        name query string false "Filter by name (ILIKE)"
// @Param        is_active query bool false "Filter by active"
// @Produce      json
// @Success      200 {object} map[string]interface{}
func (h *SystemHandler) List(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.SystemListQuery](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	items, total, page, size, err := h.Service.System.List(ctx, q)
	if err != nil {
		h.Log.Error("system_list_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// GET /system/:id
// @Summary      Get system by id
// @Tags         systems
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
func (h *SystemHandler) Get(c *fiber.Ctx) error {
	params, ok := resp.ParamsBindAndValidate[common.ID](c) // dto.IDInt{ Id int `params:"id" validate:"required"` }
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	item, err := h.Service.System.ByID(ctx, params.ID)
	if err != nil {
		h.Log.Warn("system_get_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, item)
}

// POST /system
// @Summary      Create system
// @Tags         systems
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.SystemCreateDto true "payload"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
func (h *SystemHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.SystemCreateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.System.Create(ctx, req); err != nil {
		h.Log.Error("system_create_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Created(c)
}

// PUT /system/:id
// @Summary      Update system
// @Tags         systems
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path int                    true "System ID"
// @Param        body body dto.SystemUpdateDto    true "payload"
// @Success      200 {object} map[string]interface{}
func (h *SystemHandler) Update(c *fiber.Ctx) error {
	params, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	req, ok := resp.BodyBindAndValidate[dto.SystemUpdateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.System.Update(ctx, params.ID, req); err != nil {
		h.Log.Warn("system_update_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// DELETE /system/:id
// @Summary      Delete system (soft)
// @Tags         systems
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
func (h *SystemHandler) Delete(c *fiber.Ctx) error {
	params, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.System.Delete(ctx, params.ID); err != nil {
		h.Log.Warn("system_delete_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
