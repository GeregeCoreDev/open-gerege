// Package handlers provides implementation for handlers
//
// File: client_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
// internal/http/handlers/client.go
package handlers

import (
	"net/http"
	"templatev25/internal/app"
	"templatev25/internal/auth"
	"git.gerege.mn/backend-packages/httpx"
	"git.gerege.mn/backend-packages/resp"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ClientHandler struct{ *app.Dependencies }

func NewClientHandler(d *app.Dependencies) *ClientHandler {
	return &ClientHandler{Dependencies: d}
}

// --- helpers ---

func (h *ClientHandler) headers(c *fiber.Ctx) (map[string]string, error) {
	sid := auth.ExtractSID(c, h.Cfg)
	if sid == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "missing session")

	}
	return map[string]string{
		"Cookie":        "sid=" + sid,
		"Authorization": httpx.BasicAuth(h.Cfg.Auth.ClientID, h.Cfg.Auth.ClientSecret),
	}, nil
}

func (h *ClientHandler) doSSO(c *fiber.Ctx, method, path string, body any) error {
	headers, err := h.headers(c)
	if err != nil {
		return err
	}

	client := httpx.New(3 * time.Second)
	var res any
	res, _, err = httpx.DoJSON[any](c.UserContext(), client, httpx.Request{URL: h.Cfg.URLS.SSO + path, Method: method, Headers: headers, Body: body})
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, res)
}

// --- endpoints ---

// List godoc
// @Summary      List OAuth clients
// @Tags         client
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /client [get]
func (h *ClientHandler) List(c *fiber.Ctx) error {

	return h.doSSO(c, http.MethodGet, "/client/mini", nil)
}

// ScopeList godoc
// @Summary      List OAuth scopes
// @Tags         client
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /client/scope [get]
func (h *ClientHandler) ScopeList(c *fiber.Ctx) error {
	return h.doSSO(c, http.MethodGet, "/authz/scope/v2", nil)
}

// ScopeCreate godoc
// @Summary      Create OAuth scope
// @Tags         client
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body map[string]interface{} true "Scope data"
// @Success      200 {object} map[string]interface{}
// @Router       /client/scope [post]
func (h *ClientHandler) ScopeCreate(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[map[string]any](c)
	if !ok {
		return nil
	}
	return h.doSSO(c, http.MethodPost, "/authz/scope", req)
}

// ScopeDelete godoc
// @Summary      Delete OAuth scope
// @Tags         client
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body map[string]interface{} true "Scope data"
// @Success      200 {object} map[string]interface{}
// @Router       /client/scope [delete]
func (h *ClientHandler) ScopeDelete(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[map[string]any](c)
	if !ok {
		return nil
	}
	return h.doSSO(c, http.MethodDelete, "/authz/scope", req)
}
