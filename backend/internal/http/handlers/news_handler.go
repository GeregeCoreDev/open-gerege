// Package handlers provides implementation for handlers
//
// File: news_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"templatev25/internal/http/dto"

	"strconv"

	"templatev25/internal/app"
	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/resp"

	"github.com/gofiber/fiber/v2"
)

type NewsHandler struct{ *app.Dependencies }

func NewNewsHandler(d *app.Dependencies) *NewsHandler {
	return &NewsHandler{Dependencies: d}
}

// List godoc
// @Summary      List news
// @Description  Get paginated list of news articles
// @Tags         news
// @Produce      json
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} dto.PaginatedResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /news [get]
func (h *NewsHandler) List(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.NewsListQuery](c)
	if !ok {
		return nil
	}
	items, total, page, size, err := h.Service.News.List(c.UserContext(), q)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Get godoc
// @Summary      Get news by ID
// @Tags         news
// @Produce      json
// @Param        id path int true "News ID"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /news/{id} [get]
func (h *NewsHandler) Get(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	out, err := h.Service.News.GetByID(c.UserContext(), int(id64))
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, out)
}

// Create godoc
// @Summary      Create news
// @Tags         news
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.NewsDto true "News data"
// @Success      201 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /news [post]
func (h *NewsHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.NewsDto](c)
	if !ok {
		return nil
	}

	err := h.Service.News.Create(c.UserContext(), req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Created(c)
}

// Update godoc
// @Summary      Update news
// @Tags         news
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path int true "News ID"
// @Param        body body dto.NewsDto true "News data"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /news/{id} [put]
func (h *NewsHandler) Update(c *fiber.Ctx) error {
	idp, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	req, ok := resp.BodyBindAndValidate[dto.NewsDto](c)
	if !ok {
		return nil
	}

	err := h.Service.News.Update(c.UserContext(), idp.ID, req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Delete godoc
// @Summary      Delete news
// @Tags         news
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "News ID"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /news/{id} [delete]
func (h *NewsHandler) Delete(c *fiber.Ctx) error {
	idp, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	if err := h.Service.News.Delete(c.UserContext(), idp.ID); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
