// Package handlers provides implementation for handlers
//
// File: platform_icon_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"templatev25/internal/http/dto"

	"templatev25/internal/app"
	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/resp"

	"github.com/gofiber/fiber/v2"
)

type AppServiceIconHandler struct {
	*app.Dependencies
}

type AppServiceIconGroupHandler struct {
	*app.Dependencies
}

func NewAppServiceIconHandler(d *app.Dependencies) *AppServiceIconHandler {
	return &AppServiceIconHandler{Dependencies: d}
}

// app-service-group
func NewAppServiceGroupHandler(d *app.Dependencies) *AppServiceIconGroupHandler {
	return &AppServiceIconGroupHandler{Dependencies: d}
}


// ----- App Service Icon -----

// GET /app-service-icon
func (h *AppServiceIconHandler) List(c *fiber.Ctx) error {
	items, err := h.Service.AppServiceIcon.List(c.UserContext())
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, items)
}

// POST /app-service-icon
func (h *AppServiceIconHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.AppServiceIconDto](c)
	if !ok {
		return nil
	}
	err := h.Service.AppServiceIcon.Create(c.UserContext(), req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// PUT /app-service-icon/{id}
func (h *AppServiceIconHandler) Update(c *fiber.Ctx) error {
	idp, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	req, ok := resp.BodyBindAndValidate[dto.AppServiceIconDto](c)
	if !ok {
		return nil
	}
	err := h.Service.AppServiceIcon.Update(c.UserContext(), idp.ID, req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// DELETE /app-service-icon/{id}
func (h *AppServiceIconHandler) Delete(c *fiber.Ctx) error {
	idp, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	if err := h.Service.AppServiceIcon.Delete(c.UserContext(), idp.ID); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// GET /app-service-group
func (h *AppServiceIconGroupHandler) List(c *fiber.Ctx) error {
	items, err := h.Service.AppServiceGroup.List(c.UserContext())
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, items)
}

// GET /app-service-group/with-icons
func (h *AppServiceIconGroupHandler) ListGroupsWithIcons(c *fiber.Ctx) error {
	items, err := h.Service.AppServiceGroup.ListGroupsWithIcons(c.UserContext())
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, items)
}

// POST /app-service-group
func (h *AppServiceIconGroupHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.AppServiceIconGroupDto](c)
	if !ok {
		return nil
	}
	err := h.Service.AppServiceGroup.Create(c.UserContext(), req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// PUT /app-service-group/{id}
func (h *AppServiceIconGroupHandler) Update(c *fiber.Ctx) error {
	idp, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	req, ok := resp.BodyBindAndValidate[dto.AppServiceIconGroupDto](c)
	if !ok {
		return nil
	}
	err := h.Service.AppServiceGroup.Update(c.UserContext(), idp.ID, req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// DELETE /app-service-group/{id}
func (h *AppServiceIconGroupHandler) Delete(c *fiber.Ctx) error {
	idp, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	if err := h.Service.AppServiceGroup.Delete(c.UserContext(), idp.ID); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
