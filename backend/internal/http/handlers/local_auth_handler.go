// Package handlers provides implementation for handlers
//
// File: local_auth_handler.go
// Description: Handler for local authentication (login, MFA, logout)
package handlers

import (
	"errors"

	"templatev25/internal/http/dto"
	"templatev25/internal/service"

	"git.gerege.mn/backend-packages/resp"
	"github.com/gofiber/fiber/v2"
)

// LocalAuthHandler handles local authentication endpoints
type LocalAuthHandler struct {
	authService *service.AuthService
}

// NewLocalAuthHandler creates a new local auth handler
func NewLocalAuthHandler(authService *service.AuthService) *LocalAuthHandler {
	return &LocalAuthHandler{
		authService: authService,
	}
}

// Login godoc
// @Summary      Login with email and password
// @Description  Authenticate user with email and password. Returns MFA token if MFA is enabled.
// @Tags         local-auth
// @Accept       json
// @Produce      json
// @Param        body body dto.LoginRequest true "Login credentials"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse "Invalid credentials"
// @Failure      423 {object} dto.ErrorResponse "Account locked"
// @Router       /auth/local/login [post]
func (h *LocalAuthHandler) Login(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.LoginRequest](c)
	if !ok {
		return nil
	}

	loginReq := service.LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	}

	result, err := h.authService.Login(c.UserContext(), loginReq)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "invalid email or password",
			})
		case errors.Is(err, service.ErrAccountLocked):
			return c.Status(fiber.StatusLocked).JSON(fiber.Map{
				"success": false,
				"message": "account is locked due to too many failed attempts",
			})
		case errors.Is(err, service.ErrAccountNotActive):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "account is not active",
			})
		case errors.Is(err, service.ErrCredentialsNotFound):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "local authentication not set up for this account",
			})
		default:
			return resp.InternalServerError(c, err.Error())
		}
	}

	response := dto.LoginResponse{
		RequiresMFA: result.RequiresMFA,
		MFAToken:    result.MFAToken,
	}

	if !result.RequiresMFA && result.Session != nil {
		response.AccessToken = result.Session.SessionID
		response.ExpiresAt = result.Session.ExpiresAt.Unix()
		if result.User != nil {
			response.User = &dto.UserInfo{
				ID:        result.User.Id,
				Email:     result.User.Email,
				FirstName: result.User.FirstName,
				LastName:  result.User.LastName,
				Status:    result.User.Status,
			}
		}
	}

	return resp.OK(c, response)
}

// VerifyMFA godoc
// @Summary      Verify MFA code
// @Description  Verify TOTP code to complete login
// @Tags         local-auth
// @Accept       json
// @Produce      json
// @Param        body body dto.VerifyMFARequest true "MFA verification"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse "Invalid code"
// @Router       /auth/local/verify-mfa [post]
func (h *LocalAuthHandler) VerifyMFA(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.VerifyMFARequest](c)
	if !ok {
		return nil
	}

	mfaReq := service.VerifyMFARequest{
		MFAToken:  req.MFAToken,
		Code:      req.Code,
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	}

	result, err := h.authService.VerifyMFA(c.UserContext(), mfaReq)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidMFACode):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "invalid MFA code",
			})
		case errors.Is(err, service.ErrInvalidSession):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "invalid or expired MFA token",
			})
		case errors.Is(err, service.ErrMFANotEnabled):
			return resp.BadRequest(c, "MFA is not enabled", nil)
		default:
			return resp.InternalServerError(c, err.Error())
		}
	}

	response := dto.LoginResponse{
		AccessToken: result.Session.SessionID,
		ExpiresAt:   result.Session.ExpiresAt.Unix(),
	}

	if result.User != nil {
		response.User = &dto.UserInfo{
			ID:        result.User.Id,
			Email:     result.User.Email,
			FirstName: result.User.FirstName,
			LastName:  result.User.LastName,
			Status:    result.User.Status,
		}
	}

	return resp.OK(c, response)
}

// VerifyBackupCode godoc
// @Summary      Verify backup code
// @Description  Verify backup code to complete login when TOTP is unavailable
// @Tags         local-auth
// @Accept       json
// @Produce      json
// @Param        body body dto.VerifyBackupCodeRequest true "Backup code verification"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse "Invalid code"
// @Router       /auth/local/verify-backup-code [post]
func (h *LocalAuthHandler) VerifyBackupCode(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.VerifyBackupCodeRequest](c)
	if !ok {
		return nil
	}

	result, err := h.authService.VerifyBackupCode(
		c.UserContext(),
		req.MFAToken,
		req.Code,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidMFACode):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "invalid backup code",
			})
		case errors.Is(err, service.ErrInvalidSession):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "invalid or expired MFA token",
			})
		default:
			return resp.InternalServerError(c, err.Error())
		}
	}

	response := dto.LoginResponse{
		AccessToken: result.Session.SessionID,
		ExpiresAt:   result.Session.ExpiresAt.Unix(),
	}

	if result.User != nil {
		response.User = &dto.UserInfo{
			ID:        result.User.Id,
			Email:     result.User.Email,
			FirstName: result.User.FirstName,
			LastName:  result.User.LastName,
			Status:    result.User.Status,
		}
	}

	return resp.OK(c, response)
}

// Logout godoc
// @Summary      Logout current session
// @Description  Logout and invalidate current session
// @Tags         local-auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.Response
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/logout [post]
func (h *LocalAuthHandler) Logout(c *fiber.Ctx) error {
	sessionID := getSessionID(c)
	if sessionID == "" {
		return resp.Unauthorized(c)
	}

	err := h.authService.Logout(
		c.UserContext(),
		sessionID,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, fiber.Map{"message": "logged out successfully"})
}

// LogoutAll godoc
// @Summary      Logout all sessions
// @Description  Logout and invalidate all user sessions
// @Tags         local-auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.Response
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/logout-all [post]
func (h *LocalAuthHandler) LogoutAll(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return resp.Unauthorized(c)
	}

	err := h.authService.LogoutAll(
		c.UserContext(),
		userID,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, fiber.Map{"message": "all sessions revoked"})
}

// RefreshSession godoc
// @Summary      Refresh session
// @Description  Extend session expiration time
// @Tags         local-auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.LoginResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/refresh [post]
func (h *LocalAuthHandler) RefreshSession(c *fiber.Ctx) error {
	sessionID := getSessionID(c)
	if sessionID == "" {
		return resp.Unauthorized(c)
	}

	session, err := h.authService.RefreshSession(c.UserContext(), sessionID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidSession) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "session expired",
			})
		}
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, dto.LoginResponse{
		AccessToken: session.SessionID,
		ExpiresAt:   session.ExpiresAt.Unix(),
	})
}

// Helper functions
func getSessionID(c *fiber.Ctx) string {
	// Try to get from context (set by session auth middleware)
	if sid, ok := c.Locals("session_id").(string); ok {
		return sid
	}

	// Try Authorization header
	auth := c.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}

	return ""
}

func getUserID(c *fiber.Ctx) int {
	if userID, ok := c.Locals("user_id").(int); ok {
		return userID
	}
	return 0
}
