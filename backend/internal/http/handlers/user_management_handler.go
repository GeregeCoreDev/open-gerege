// Package handlers provides implementation for handlers
//
// File: user_management_handler.go
// Description: Handler for user management (MFA, sessions, password, security)
package handlers

import (
	"errors"
	"strconv"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/service"

	"git.gerege.mn/backend-packages/resp"
	"github.com/gofiber/fiber/v2"
)

// UserManagementHandler handles user management endpoints
type UserManagementHandler struct {
	authService *service.AuthService
}

// NewUserManagementHandler creates a new user management handler
func NewUserManagementHandler(authService *service.AuthService) *UserManagementHandler {
	return &UserManagementHandler{
		authService: authService,
	}
}

// ============================================================
// MFA ENDPOINTS
// ============================================================

// GetMFAStatus godoc
// @Summary      Get MFA status
// @Tags         local-auth-user
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.MFAStatusResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/me/mfa [get]
func (h *UserManagementHandler) GetMFAStatus(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return resp.Unauthorized(c)
	}

	enabled, hasBackup, err := h.authService.GetMFAStatus(c.UserContext(), userID)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, dto.MFAStatusResponse{
		Enabled:        enabled,
		HasBackupCodes: hasBackup,
	})
}

// SetupTOTP godoc
// @Summary      Setup TOTP MFA
// @Tags         local-auth-user
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.TOTPSetupResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/me/mfa/totp/setup [post]
func (h *UserManagementHandler) SetupTOTP(c *fiber.Ctx) error {
	userID := getUserID(c)
	email := getEmail(c)
	if userID == 0 {
		return resp.Unauthorized(c)
	}

	result, err := h.authService.SetupTOTP(c.UserContext(), userID, email)
	if err != nil {
		if errors.Is(err, service.ErrMFAAlreadyEnabled) {
			return resp.BadRequest(c, "MFA is already enabled", nil)
		}
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, dto.TOTPSetupResponse{
		Secret:    result.Secret,
		QRCodeURL: result.QRCodeURL,
	})
}

// ConfirmTOTP godoc
// @Summary      Confirm TOTP setup
// @Tags         local-auth-user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.ConfirmTOTPRequest true "TOTP code"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/me/mfa/totp/confirm [post]
func (h *UserManagementHandler) ConfirmTOTP(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return resp.Unauthorized(c)
	}

	req, ok := resp.BodyBindAndValidate[dto.ConfirmTOTPRequest](c)
	if !ok {
		return nil
	}

	err := h.authService.ConfirmTOTP(
		c.UserContext(),
		userID,
		req.Code,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidMFACode) {
			return resp.BadRequest(c, "invalid code", nil)
		}
		if errors.Is(err, service.ErrMFAAlreadyEnabled) {
			return resp.BadRequest(c, "MFA is already enabled", nil)
		}
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, fiber.Map{"message": "MFA enabled successfully"})
}

// DisableTOTP godoc
// @Summary      Disable TOTP MFA
// @Tags         local-auth-user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.DisableTOTPRequest true "TOTP code"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/me/mfa/totp [delete]
func (h *UserManagementHandler) DisableTOTP(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return resp.Unauthorized(c)
	}

	req, ok := resp.BodyBindAndValidate[dto.DisableTOTPRequest](c)
	if !ok {
		return nil
	}

	err := h.authService.DisableTOTP(
		c.UserContext(),
		userID,
		req.Code,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidMFACode) {
			return resp.BadRequest(c, "invalid code", nil)
		}
		if errors.Is(err, service.ErrMFANotEnabled) {
			return resp.BadRequest(c, "MFA is not enabled", nil)
		}
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, fiber.Map{"message": "MFA disabled successfully"})
}

// GenerateBackupCodes godoc
// @Summary      Generate new backup codes
// @Tags         local-auth-user
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.BackupCodesResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/me/mfa/backup-codes [post]
func (h *UserManagementHandler) GenerateBackupCodes(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return resp.Unauthorized(c)
	}

	codes, err := h.authService.GenerateBackupCodes(
		c.UserContext(),
		userID,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		if errors.Is(err, service.ErrMFANotEnabled) {
			return resp.BadRequest(c, "MFA is not enabled", nil)
		}
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, dto.BackupCodesResponse{Codes: codes})
}

// ============================================================
// SESSION ENDPOINTS
// ============================================================

// ListSessions godoc
// @Summary      List active sessions
// @Tags         local-auth-user
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.SessionListResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/me/sessions [get]
func (h *UserManagementHandler) ListSessions(c *fiber.Ctx) error {
	userID := getUserID(c)
	currentSessionID := getSessionID(c)
	if userID == 0 {
		return resp.Unauthorized(c)
	}

	sessions, err := h.authService.GetActiveSessions(c.UserContext(), userID)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	var sessionInfos []dto.SessionInfoResponse
	for _, s := range sessions {
		sessionInfos = append(sessionInfos, dto.SessionInfoResponse{
			SessionID:  s.SessionID,
			IPAddress:  s.IPAddress,
			UserAgent:  s.UserAgent,
			CreatedAt:  s.CreatedAt,
			LastActive: s.LastActivityAt,
			IsCurrent:  s.SessionID == currentSessionID,
		})
	}

	return resp.OK(c, dto.SessionListResponse{
		Sessions: sessionInfos,
		Total:    len(sessionInfos),
	})
}

// RevokeSession godoc
// @Summary      Revoke a specific session
// @Tags         local-auth-user
// @Security     BearerAuth
// @Param        id path string true "Session ID"
// @Produce      json
// @Success      200 {object} dto.Response
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/me/sessions/{id} [delete]
func (h *UserManagementHandler) RevokeSession(c *fiber.Ctx) error {
	sessionID := c.Params("id")
	if sessionID == "" {
		return resp.BadRequest(c, "session id required", nil)
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

	return resp.OK(c, fiber.Map{"message": "session revoked"})
}

// ============================================================
// PASSWORD ENDPOINTS
// ============================================================

// ChangePassword godoc
// @Summary      Change password
// @Tags         local-auth-user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.ChangePasswordRequest true "Password change"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/me/password [post]
func (h *UserManagementHandler) ChangePassword(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return resp.Unauthorized(c)
	}

	req, ok := resp.BodyBindAndValidate[dto.ChangePasswordRequest](c)
	if !ok {
		return nil
	}

	err := h.authService.ChangePassword(
		c.UserContext(),
		userID,
		req.CurrentPassword,
		req.NewPassword,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			return resp.BadRequest(c, "current password is incorrect", nil)
		case errors.Is(err, service.ErrPasswordTooWeak):
			return resp.BadRequest(c, "password does not meet requirements", nil)
		case errors.Is(err, service.ErrPasswordReused):
			return resp.BadRequest(c, "password was recently used", nil)
		case errors.Is(err, service.ErrCredentialsNotFound):
			return resp.BadRequest(c, "local authentication not set up", nil)
		default:
			return resp.InternalServerError(c, err.Error())
		}
	}

	return resp.OK(c, fiber.Map{"message": "password changed successfully"})
}

// ============================================================
// HISTORY ENDPOINTS
// ============================================================

// GetLoginHistory godoc
// @Summary      Get login history
// @Tags         local-auth-user
// @Security     BearerAuth
// @Param        limit query int false "Limit (default 50)"
// @Produce      json
// @Success      200 {object} dto.LoginHistoryResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/me/login-history [get]
func (h *UserManagementHandler) GetLoginHistory(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return resp.Unauthorized(c)
	}

	limit := c.QueryInt("limit", 50)
	if limit > 100 {
		limit = 100
	}

	history, err := h.authService.GetLoginHistory(c.UserContext(), userID, limit)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	var entries []dto.LoginHistoryEntry
	for _, h := range history {
		entries = append(entries, dto.LoginHistoryEntry{
			ID:            h.ID,
			Email:         h.Email,
			IPAddress:     h.IPAddress,
			UserAgent:     h.UserAgent,
			LoginMethod:   h.LoginMethod,
			Success:       h.Success,
			FailureReason: h.FailureReason,
			MFAUsed:       h.MFAUsed,
		})
	}

	return resp.OK(c, dto.LoginHistoryResponse{
		Entries: entries,
		Total:   len(entries),
	})
}

// GetSecurityAudit godoc
// @Summary      Get security audit trail
// @Tags         local-auth-user
// @Security     BearerAuth
// @Param        limit query int false "Limit (default 50)"
// @Produce      json
// @Success      200 {object} dto.SecurityAuditResponse
// @Failure      401 {object} dto.ErrorResponse
// @Router       /auth/local/me/security-audit [get]
func (h *UserManagementHandler) GetSecurityAudit(c *fiber.Ctx) error {
	userID := getUserID(c)
	if userID == 0 {
		return resp.Unauthorized(c)
	}

	limit := c.QueryInt("limit", 50)
	if limit > 100 {
		limit = 100
	}

	audit, err := h.authService.GetSecurityAudit(c.UserContext(), userID, limit)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	var entries []dto.SecurityAuditEntry
	for _, a := range audit {
		entries = append(entries, dto.SecurityAuditEntry{
			ID:         a.ID,
			Action:     a.Action,
			TargetType: a.TargetType,
			TargetID:   a.TargetID,
			OldValue:   a.OldValue,
			NewValue:   a.NewValue,
			IPAddress:  a.IPAddress,
			UserAgent:  a.UserAgent,
		})
	}

	return resp.OK(c, dto.SecurityAuditResponse{
		Entries: entries,
		Total:   len(entries),
	})
}

// ============================================================
// ADMIN ENDPOINTS
// ============================================================

// UpdateUserStatus godoc
// @Summary      Update user status (admin)
// @Tags         user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Param        body body dto.UpdateUserStatusRequest true "Status update"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Router       /user/{id}/status [put]
func (h *UserManagementHandler) UpdateUserStatus(c *fiber.Ctx) error {
	adminID := getUserID(c)
	if adminID == 0 {
		return resp.Unauthorized(c)
	}

	targetID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return resp.BadRequest(c, "invalid user id", nil)
	}

	req, ok := resp.BodyBindAndValidate[dto.UpdateUserStatusRequest](c)
	if !ok {
		return nil
	}

	// Validate status
	status := domain.UserStatus(req.Status)
	if !status.IsValid() {
		return resp.BadRequest(c, "invalid status", nil)
	}

	err = h.authService.UpdateUserStatus(
		c.UserContext(),
		targetID,
		req.Status,
		req.Reason,
		adminID,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, fiber.Map{"message": "status updated"})
}

// UnlockUser godoc
// @Summary      Unlock user account (admin)
// @Tags         user
// @Security     BearerAuth
// @Param        id path int true "User ID"
// @Produce      json
// @Success      200 {object} dto.Response
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Router       /user/{id}/unlock [post]
func (h *UserManagementHandler) UnlockUser(c *fiber.Ctx) error {
	adminID := getUserID(c)
	if adminID == 0 {
		return resp.Unauthorized(c)
	}

	targetID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return resp.BadRequest(c, "invalid user id", nil)
	}

	err = h.authService.UnlockAccount(
		c.UserContext(),
		targetID,
		adminID,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, fiber.Map{"message": "account unlocked"})
}

// SetUserPassword godoc
// @Summary      Set user password (admin)
// @Tags         user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Param        body body dto.SetPasswordRequest true "Password"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Router       /user/{id}/password [post]
func (h *UserManagementHandler) SetUserPassword(c *fiber.Ctx) error {
	adminID := getUserID(c)
	if adminID == 0 {
		return resp.Unauthorized(c)
	}

	targetID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return resp.BadRequest(c, "invalid user id", nil)
	}

	req, ok := resp.BodyBindAndValidate[dto.SetPasswordRequest](c)
	if !ok {
		return nil
	}

	err = h.authService.SetPassword(c.UserContext(), targetID, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrPasswordTooWeak) {
			return resp.BadRequest(c, "password does not meet requirements", nil)
		}
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, fiber.Map{"message": "password set"})
}

// Helper function
func getEmail(c *fiber.Ctx) string {
	if email, ok := c.Locals("email").(string); ok {
		return email
	}
	return ""
}
