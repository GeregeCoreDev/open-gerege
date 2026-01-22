// Package handlers provides implementation for handlers
//
// File: auth_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"fmt"
	"templatev25/internal/app"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/ctx"
	"git.gerege.mn/backend-packages/resp"
	ssoclient "git.gerege.mn/backend-packages/sso-client"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	*app.Dependencies
}

func NewAuthHandler(d *app.Dependencies) *AuthHandler {
	return &AuthHandler{Dependencies: d}
}

// InitDirection godoc
// @Summary      Redirect to OAuth login
// @Description  Redirects user to SSO login page with PKCE
// @Tags         auth
// @Produce      json
// @Success      302 "Redirect to SSO"
// @Router       /auth/login [get]
func (h *AuthHandler) InitDirection(c *fiber.Ctx) error {
	result, err := ssoclient.InitOAuthDirection(c, h.Cfg.Auth, h.Cfg.Cookie)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Redirect(result.AuthURL, fiber.StatusFound)
}

// OAuthCallback godoc
// @Summary      OAuth callback handler
// @Description  Handles OAuth callback, sets session cookie
// @Tags         auth
// @Param        state query string true "OAuth state"
// @Param        sid   query string true "Session ID"
// @Success      302 "Redirect to callback URL"
// @Failure      400 {object} map[string]interface{} "Invalid state"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Router       /auth/callback [get]
func (h *AuthHandler) OAuthCallback(c *fiber.Ctx) error {
	_, err := ssoclient.HandleOAuthCallbackAndSetCookie(
		c, h.SSO, c.Query("state"), c.Query("sid"), ctx.RequestID(c), h.Cfg.Cookie,
	)
	if err != nil {
		return err
	}

	return c.Redirect(h.Cfg.Auth.CallbackURL, fiber.StatusFound)
}

// Logout godoc
// @Summary      Logout user
// @Description  Invalidates session and clears cookie
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	_ = ssoclient.LogoutAndDeleteCookie(
		c, h.SSO, ssoclient.GetSID(c), ctx.RequestID(c), h.Cfg.Cookie,
	)

	return resp.OK(c, fiber.Map{"message": "logged out"})
}

// GoogleLogin godoc
// @Summary      Google OAuth login
// @Description  Login with Google ID token
// @Tags         auth
// @Param        id_token query string true "Google ID token"
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{} "id_token required"
// @Router       /auth/google/login [post]
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	idToken := c.Query("id_token")
	if idToken == "" {
		return resp.BadRequestValidation(c, fmt.Errorf("id_token is required"))
	}

	result, err := ssoclient.HandleGoogleLoginAndSetCookie(
		c, h.SSO, idToken, ctx.RequestID(c), h.Cfg.Cookie,
	)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, common.SID{Sid: result.SID})
}

// AuthVerify godoc
// @Summary      Verify authentication
// @Description  Verifies session and sets cookie
// @Tags         auth
// @Param        sid query string false "Session ID"
// @Success      302 "Redirect to callback URL"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Router       /auth/verify [get]
func (h *AuthHandler) AuthVerify(c *fiber.Ctx) error {
	result, err := ssoclient.VerifyAndSetCookie(
		c, h.SSO, c.Query("sid"), ctx.RequestID(c), h.Cfg.Cookie,
	)
	if err != nil {
		return err
	}

	if result == nil {
		return c.Redirect(h.Cfg.Auth.CallbackURL, fiber.StatusTemporaryRedirect)
	}

	return c.Redirect(h.Cfg.Auth.CallbackURL, fiber.StatusFound)
}

// ChangeOrganization godoc
// @Summary      Change current organization
// @Description  Switches user's active organization
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body common.ID true "Organization ID"
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /auth/org/change [post]
func (h *AuthHandler) ChangeOrganization(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}

	sid := ssoclient.GetSID(c)
	rid := ctx.RequestID(c)

	_, err := ssoclient.ChangeOrganizationAndSetCookie(
		c, h.SSO, sid, rid, req.ID, h.Cfg.Cookie,
	)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c)
}
