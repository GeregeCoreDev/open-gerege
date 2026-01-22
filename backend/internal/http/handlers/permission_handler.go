// Package handlers provides implementation for handlers
//
// File: permission_handler.go
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
	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/resp"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type PermissionHandler struct {
	*app.Dependencies
}

func NewPermissionHandler(d *app.Dependencies) *PermissionHandler {
	return &PermissionHandler{Dependencies: d}
}

// List godoc
// @Summary      List permissions (paginated)
// @Tags         permissions
// @Security     BearerAuth
// @Param        search    query   string false "Search (code/name/description)"
// @Param        module_id query   int    false "Filter by module_id"
// @Param        page      query   int    false "Page number"
// @Param        size      query   int    false "Page size"
// @Param        sort      query   string false "Sort (e.g. code:asc,name:desc)"
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /permissions [get]
func (h *PermissionHandler) List(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.PermissionQuery](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	items, total, page, size, err := h.Service.Permission.ListFilteredPaged(ctx, q)
	if err != nil {
		h.Log.Error("permission_list_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Create godoc
// @Summary      Create permission
// @Tags         permissions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body   body dto.PermissionCreateDto true "payload"
// @Success      201 {object} map[string]interface{}
// @Router       /permissions [post]
func (h *PermissionHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.PermissionCreateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Permission.Create(ctx, req); err != nil {
		h.Log.Error("permission_create_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Created(c)
}

// Update godoc
// @Summary      Update permission
// @Tags         permissions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path int                        true "ID"
// @Param        body   body dto.PermissionCreateDto    true "payload"
// @Success      200    {object} map[string]interface{}
// @Router       /permissions/{id} [put]
func (h *PermissionHandler) Update(c *fiber.Ctx) error {
	params, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}

	req, ok := resp.BodyBindAndValidate[dto.PermissionUpdateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Permission.Update(ctx, params.ID, req); err != nil {
		h.Log.Warn("permission_update_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Delete godoc
// @Summary      Delete permission (soft)
// @Tags         permissions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path int true "ID"
// @Success      200  {object} map[string]interface{}
// @Router       /permissions/{id} [delete]
func (h *PermissionHandler) Delete(c *fiber.Ctx) error {
	params, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Permission.Delete(ctx, params.ID); err != nil {
		h.Log.Warn("permission_delete_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
