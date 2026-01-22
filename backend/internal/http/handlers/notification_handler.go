// Package handlers provides implementation for handlers
//
// File: notification_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"templatev25/internal/http/dto"

	"templatev25/internal/app"
	"git.gerege.mn/backend-packages/sso-client"
	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/resp"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	*app.Dependencies
}

func NewNotificationHandler(d *app.Dependencies) *NotificationHandler {
	return &NotificationHandler{Dependencies: d}
}

// List godoc
// @Summary      List notifications
// @Description  Get paginated list of user notifications
// @Tags         notification
// @Security     BearerAuth
// @Produce      json
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Router       /notification [get]
func (h *NotificationHandler) List(c *fiber.Ctx) error {
	p, ok := resp.ParamsBindAndValidate[common.PaginationQuery](c)
	if !ok {
		return nil
	}

	claims, ok := ssoclient.GetClaims(c)
	if !ok {
		return resp.Unauthorized(c)
	}

	items, total, page, size, err := h.Service.Notification.List(c.UserContext(), claims.UserID, p)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Groups godoc
// @Summary      List notification groups
// @Tags         notification
// @Security     BearerAuth
// @Produce      json
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Router       /notification/groups [get]
func (h *NotificationHandler) Groups(c *fiber.Ctx) error {
	p, ok := resp.ParamsBindAndValidate[common.PaginationQuery](c)
	if !ok {
		return nil
	}
	items, total, page, size, err := h.Service.Notification.Groups(c.UserContext(), p)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Read godoc
// @Summary      Mark notification group as read
// @Tags         notification
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.NotificationReadDto true "Group ID"
// @Success      200 {object} map[string]interface{}
// @Router       /notification/read [post]
func (h *NotificationHandler) Read(c *fiber.Ctx) error {
	req, ok := resp.ParamsBindAndValidate[dto.NotificationReadDto](c)
	if !ok {
		return nil
	}

	claims, ok := ssoclient.GetClaims(c)
	if !ok {
		return resp.Unauthorized(c)
	}

	if err := h.Service.Notification.MarkGroupRead(c.UserContext(), claims.UserID, req.GroupId); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// ReadAll godoc
// @Summary      Mark all notifications as read
// @Tags         notification
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /notification/read-all [post]
func (h *NotificationHandler) ReadAll(c *fiber.Ctx) error {
	claims, ok := ssoclient.GetClaims(c)
	if !ok {
		return resp.Unauthorized(c)
	}
	if err := h.Service.Notification.MarkAllRead(c.UserContext(), claims.UserID); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Send godoc
// @Summary      Send notification
// @Tags         notification
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.NotificationSendDto true "Notification data"
// @Success      200 {object} map[string]interface{}
// @Router       /notification [post]
func (h *NotificationHandler) Send(c *fiber.Ctx) error {
	req, ok := resp.ParamsBindAndValidate[dto.NotificationSendDto](c)
	if !ok {
		return nil
	}

	claims, ok := ssoclient.GetClaims(c)
	if !ok {
		return resp.Unauthorized(c)
	}

	if err := h.Service.Notification.Send(
		c.UserContext(),
		req,
		claims.Username,
	); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
