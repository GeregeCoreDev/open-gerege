// Package handlers provides implementation for handlers
//
// File: menu_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"context"
	"strconv"
	"time"

	"templatev25/internal/app"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/resp"
	ssoclient "git.gerege.mn/backend-packages/sso-client"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type MenuHandler struct {
	*app.Dependencies
}

func NewMenuHandler(d *app.Dependencies) *MenuHandler {
	return &MenuHandler{Dependencies: d}
}

// List godoc
// @Summary      List menus
// @Description  Get all menus with children preloaded (tree structure, no pagination)
// @Tags         menu
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /menu [get]
func (h *MenuHandler) List(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	items, err := h.Service.Menu.ListAll(ctx)
	if err != nil {
		h.Log.Error("menu_list_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, items)
}

// ListByRole godoc
// @Summary      List menus by user roles
// @Description  Get all menus accessible by current user's roles (tree structure)
// @Tags         menu
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /menu/my [get]
func (h *MenuHandler) ListByRole(c *fiber.Ctx) error {
	// Session-оос CitizenID авах
	claims, ok := ssoclient.GetClaims(c)
	if !ok {
		return resp.Unauthorized(c)
	}

	if claims.CitizenID == 0 {
		return resp.BadRequest(c, "citizen id required", nil)
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	// User-ийн role-уудын permission-уудтай холбоотой menu-уудыг олох
	items, err := h.Service.Menu.ListByUserRoles(ctx, claims.CitizenID)
	if err != nil {
		h.Log.Error("menu_list_by_user_roles_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, items)
}

// Get godoc
// @Summary      Get menu by ID
// @Description  Get menu by ID
// @Tags         menu
// @Security     BearerAuth
// @Produce      json
// @Param        id path int64 true "Menu ID"
// @Success      200 {object} map[string]interface{}
// @Router       /menu/{id} [get]
func (h *MenuHandler) Get(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return resp.BadRequest(c, "invalid menu id", err.Error())
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	item, err := h.Service.Menu.ByID(ctx, id64)
	if err != nil {
		h.Log.Error("menu_get_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, item)
}

// Create godoc
// @Summary      Create menu
// @Description  Create a new menu
// @Tags         menu
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.MenuCreateDto true "payload"
// @Success      201 {object} map[string]interface{}
// @Router       /menu [post]
func (h *MenuHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.MenuCreateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Menu.Create(ctx, req); err != nil {
		h.Log.Error("menu_create_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Created(c)
}

// Update godoc
// @Summary      Update menu
// @Description  Update an existing menu
// @Tags         menu
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path int64 true "Menu ID"
// @Param        body body dto.MenuUpdateDto true "payload"
// @Success      200 {object} map[string]interface{}
// @Router       /menu/{id} [put]
func (h *MenuHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return resp.BadRequest(c, "invalid menu id", err.Error())
	}

	req, ok := resp.BodyBindAndValidate[dto.MenuUpdateDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Menu.Update(ctx, id64, req); err != nil {
		h.Log.Warn("menu_update_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Delete godoc
// @Summary      Delete menu
// @Description  Delete a menu
// @Tags         menu
// @Security     BearerAuth
// @Produce      json
// @Param        id path int64 true "Menu ID"
// @Success      200 {object} map[string]interface{}
// @Router       /menu/{id} [delete]
func (h *MenuHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return resp.BadRequest(c, "invalid menu id", err.Error())
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Menu.Delete(ctx, id64); err != nil {
		h.Log.Warn("menu_delete_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
