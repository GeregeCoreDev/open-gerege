// Package handlers provides implementation for handlers
//
// File: user_role_handler.go
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

	"git.gerege.mn/backend-packages/resp"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserRoleHandler struct {
	*app.Dependencies
}

func NewUserRoleHandler(d *app.Dependencies) *UserRoleHandler {
	return &UserRoleHandler{Dependencies: d}
}

// UsersByRole godoc
// @Summary      Get users by role
// @Description  Get paginated list of users with specified role
// @Tags         role-matrix
// @Security     BearerAuth
// @Produce      json
// @Param        role_id query int true "Role ID"
// @Param        page    query int false "Page number"
// @Param        size    query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Router       /role-matrix/users [get]
func (h *UserRoleHandler) UsersByRole(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.UserRoleUsersQuery](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	items, total, page, size, err := h.Service.UserRole.UsersByRole(ctx, q)
	if err != nil {
		h.Log.Error("userrole_users_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// RolesByUser godoc
// @Summary      Get roles by user
// @Description  Get paginated list of roles for specified user
// @Tags         role-matrix
// @Security     BearerAuth
// @Produce      json
// @Param        user_id query int true "User ID"
// @Param        page    query int false "Page number"
// @Param        size    query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Router       /role-matrix/roles [get]
func (h *UserRoleHandler) RolesByUser(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.UserRoleRolesQuery](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	items, total, page, size, err := h.Service.UserRole.RolesByUser(ctx, q)
	if err != nil {
		h.Log.Error("userrole_roles_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Create godoc
// @Summary      Assign user role
// @Description  Assign role to user (by role or by user)
// @Tags         role-matrix
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.UserRoleAssignByRole true "Assignment data"
// @Success      201 {object} map[string]interface{}
// @Router       /role-matrix [post]
func (h *UserRoleHandler) Create(c *fiber.Ctx) error {
	// эхэлж "assign by role" bind оролдоно
	if req, ok := resp.BodyBindAndValidate[dto.UserRoleAssignByRole](c); ok {
		ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		defer cancel()
		if err := h.Service.UserRole.AssignByRole(ctx, req); err != nil {
			h.Log.Error("userrole_assign_by_role_failed", zap.Error(err))
			return resp.InternalServerError(c, err.Error())
		}
		return resp.Created(c)
	}

	// эсрэг тохиолдолд "assign by user" гэж үзнэ
	req2, ok := resp.BodyBindAndValidate[dto.UserRoleAssignByUser](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	if err := h.Service.UserRole.AssignByUser(ctx, req2); err != nil {
		h.Log.Error("userrole_assign_by_user_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Created(c)
}

// Delete godoc
// @Summary      Remove user role
// @Description  Remove role assignment from user
// @Tags         role-matrix
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.UserRoleRemoveDto true "Remove data"
// @Success      200 {object} map[string]interface{}
// @Router       /role-matrix [delete]
func (h *UserRoleHandler) Delete(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.UserRoleRemoveDto](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()
	if err := h.Service.UserRole.Remove(ctx, req); err != nil {
		h.Log.Warn("userrole_remove_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
