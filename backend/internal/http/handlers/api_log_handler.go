// Package handlers provides implementation for handlers
//
// File: api_log_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-01-09
// Last Updated: 2025-01-09
package handlers

import (
	"context"
	"time"

	"templatev25/internal/app"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/resp"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type APILogHandler struct {
	*app.Dependencies
}

func NewAPILogHandler(d *app.Dependencies) *APILogHandler {
	return &APILogHandler{Dependencies: d}
}

// List godoc
// @Summary      List API logs (paginated)
// @Description  Get paginated list of API logs with filtering
// @Tags         api-logs
// @Security     BearerAuth
// @Produce      json
// @Param        page        query int    false "Page number"
// @Param        size        query int    false "Page size"
// @Param        method      query string false "Filter by HTTP method (GET, POST, etc.)"
// @Param        path        query string false "Filter by path (ILIKE)"
// @Param        status_code query int    false "Filter by status code"
// @Param        user_id     query int64  false "Filter by user ID"
// @Param        org_id      query int64  false "Filter by organization ID"
// @Param        ip          query string false "Filter by IP address (ILIKE)"
// @Param        search      query string false "Search (method/path/ip/username)"
// @Param        sort        query string false "Sort (e.g. created_date:desc,id:desc)"
// @Param        created_from query string false "Filter from date (YYYY-MM-DD)"
// @Param        created_to   query string false "Filter to date (YYYY-MM-DD)"
// @Success      200 {object} map[string]interface{}
// @Router       /api-logs [get]
func (h *APILogHandler) List(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.APILogListQuery](c)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	items, total, page, size, err := h.Service.APILog.List(ctx, q)
	if err != nil {
		h.Log.Error("api_log_list_failed", zap.Error(err))
		return resp.InternalServerError(c, err.Error())
	}

	return resp.Paginated(c, items, total, page, size)
}
