// Package handlers provides implementation for handlers
//
// File: terminal_handler.go
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

type TerminalHandler struct {
	*app.Dependencies
}

func NewTerminalHandler(d *app.Dependencies) *TerminalHandler {
	return &TerminalHandler{Dependencies: d}
}

// GET /terminal
// @Tags terminal
// @Produce json
// @Param page query int false "page>=1"
// @Param size query int false "pageSize"
// @Param search query string false "JSON search (name,serial,org_id)"
// @Param sort query string false "JSON sort"
// @Param createdFrom query string false "YYYY-MM-DD"
// @Param createdTo   query string false "YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Router /terminal [get]
func (h *TerminalHandler) List(c *fiber.Ctx) error {
	p, ok := resp.QueryBindAndValidate[common.PaginationQuery](c)
	if !ok {
		return nil
	}
	items, total, page, size, err := h.Service.Terminal.List(c.UserContext(), p)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// POST /terminal
// @Tags terminal
// @Accept json
// @Produce json
// @Param body body dto.TerminalCreateDto true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /terminal [post]
func (h *TerminalHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.TerminalCreateDto](c)
	if !ok {
		return nil
	}

	err := h.Service.Terminal.Create(c.UserContext(), req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// PUT /terminal/{id}
// @Tags terminal
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param body body dto.TerminalUpdateDto true "payload"
// @Success 200 {object} map[string]interface{}
// @Router /terminal/{id} [put]
func (h *TerminalHandler) Update(c *fiber.Ctx) error {
	idp, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}

	req, ok := resp.BodyBindAndValidate[dto.TerminalUpdateDto](c)
	if !ok {
		return nil
	}

	err := h.Service.Terminal.Update(c.UserContext(), idp.ID, req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// DELETE /terminal/{id}
// @Tags terminal
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} map[string]interface{}
// @Router /terminal/{id} [delete]
func (h *TerminalHandler) Delete(c *fiber.Ctx) error {
	idp, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}

	err := h.Service.Terminal.Delete(c.UserContext(), idp.ID)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
