// Package handlers provides implementation for handlers
//
// File: module_handler.go
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

type ModuleHandler struct {
	*app.Dependencies
}

func NewModuleHandler(d *app.Dependencies) *ModuleHandler {
	return &ModuleHandler{Dependencies: d}
}

// List godoc
// @Summary      List modules
// @Description  Get paginated list of modules
// @Tags         module
// @Security     BearerAuth
// @Produce      json
// @Param        page      query int    false "Page number"
// @Param        size      query int    false "Page size"
// @Param        system_id query int    false "Filter by system_id"
// @Param        search    query string false "Search by name/code"
// @Success      200 {object} map[string]interface{}
// @Router       /module [get]
func (h *ModuleHandler) List(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.ModuleListQuery](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	items, total, page, size, err := h.Service.Module.List(ctx, q)
	if err != nil {
		h.Log.Error("module_list_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Create godoc
// @Summary      Create module
// @Tags         module
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.ModuleCreateDto true "payload"
// @Success      201 {object} map[string]interface{}
// @Router       /module [post]
func (h *ModuleHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.ModuleCreateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Module.Create(ctx, req); err != nil {
		h.Log.Error("module_create_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Created(c)
}

// Update godoc
// @Summary      Update module
// @Tags         module
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path int true "Module ID"
// @Param        body body dto.ModuleUpdateDto true "payload"
// @Success      200 {object} map[string]interface{}
// @Router       /module/{id} [put]
func (h *ModuleHandler) Update(c *fiber.Ctx) error {
	params, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	req, ok := resp.BodyBindAndValidate[dto.ModuleUpdateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Module.Update(ctx, params.ID, req); err != nil {
		h.Log.Warn("module_update_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Delete godoc
// @Summary      Delete module
// @Tags         module
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Module ID"
// @Success      200 {object} map[string]interface{}
// @Router       /module/{id} [delete]
func (h *ModuleHandler) Delete(c *fiber.Ctx) error {
	params, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Module.Delete(ctx, params.ID); err != nil {
		h.Log.Warn("module_delete_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
