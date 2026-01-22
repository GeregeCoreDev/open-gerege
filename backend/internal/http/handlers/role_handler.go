// Package handlers provides implementation for handlers
//
// File: role_handler.go
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

type RoleHandler struct {
	*app.Dependencies
}

func NewRoleHandler(d *app.Dependencies) *RoleHandler {
	return &RoleHandler{Dependencies: d}
}

// -----------------------------------------------------------------------------
// CRUD
// -----------------------------------------------------------------------------

// List godoc
// @Summary      List roles
// @Tags         role
// @Security     BearerAuth
// @Produce      json
// @Param        page query int false "Page number (>=1)"
// @Param        pageSize query int false "Page size"
// @Param        search query string false "Search JSON (e.g. [{\"field\":\"name\",\"value\":\"adm\"}])"
// @Param        sort query string false "Sort JSON (e.g. [{\"field\":\"id\",\"desc\":true}])"
// @Param        createdFrom query string false "Created from (YYYY-MM-DD)"
// @Param        createdTo query string false "Created to (YYYY-MM-DD)"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /role [get]
func (h *RoleHandler) List(c *fiber.Ctx) error {
	p, ok := resp.QueryBindAndValidate[dto.RoleListQuery](c)
	if !ok {
		return nil
	}
	items, total, page, size, err := h.Service.Role.ListFilteredPaged(c.UserContext(), p)
	if err != nil {
		h.Log.Error("access_group_list_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Create godoc
// @Summary      Create role
// @Tags         role
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.RoleCreateDto true "Role data"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /role [post]
func (h *RoleHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.RoleCreateDto](c)
	if !ok {
		return nil
	}
	// actor context-д auth middleware аль хэдийн тавьсан
	err := h.Service.Role.Create(c.UserContext(), req)
	if err != nil {
		h.Log.Error("access_group_create_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Update godoc
// @Summary      Update role
// @Tags         role
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Role ID"
// @Param        body body dto.RoleUpdateDto true "Role data"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /role/{id} [put]
func (h *RoleHandler) Update(c *fiber.Ctx) error {
	params, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	req, ok := resp.BodyBindAndValidate[dto.RoleUpdateDto](c)
	if !ok {
		return nil
	}
	err := h.Service.Role.Update(c.UserContext(), params.ID, req)
	if err != nil {
		h.Log.Error("access_group_update_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Delete godoc
// @Summary      Delete (soft) role
// @Tags         role
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Role ID"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /role/{id} [delete]
func (h *RoleHandler) Delete(c *fiber.Ctx) error {
	params, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	// ctx доторх User/Org-г repo тал уншина
	err := h.Service.Role.Delete(c.UserContext(), params.ID)
	if err != nil {
		h.Log.Error("access_group_delete_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// --- ШИНЭ: GET /role/permissions?role_id=...

// GetRolePermissions godoc
// @Summary      List permissions of a role
// @Tags         role
// @Security     BearerAuth
// @Produce      json
// @Param        role_id query int true "Role ID"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /role/permissions [get]
func (h *RoleHandler) GetRolePermissions(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.RolePermissionsQuery](c)
	if !ok {
		return nil
	}
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	items, err := h.Service.Role.GetPermissions(ctx, q)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, items)
}

// --- ШИНЭ: POST /role/permissions (replace semantics)

// SetRolePermissions godoc
// @Summary      Replace permissions of a role
// @Description  Replaces all permissions of a role with the provided list
// @Tags         role
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.RolePermissionsUpdateDto true "Permission IDs"
// @Success      201 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /role/permissions [post]
func (h *RoleHandler) SetRolePermissions(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.RolePermissionsUpdateDto](c)
	if !ok {
		return nil
	}
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	if err := h.Service.Role.SetPermissions(ctx, req); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Created(c)
}
