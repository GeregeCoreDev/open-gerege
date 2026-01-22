// Package handlers provides implementation for handlers
//
// File: user_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"templatev25/internal/http/dto"

	"fmt"
	"templatev25/internal/app"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/resp"
	ssoclient "git.gerege.mn/backend-packages/sso-client"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	*app.Dependencies
}

func NewUserHandler(d *app.Dependencies) *UserHandler {
	return &UserHandler{Dependencies: d}
}

// Me godoc
// @Summary      Get current user (claims)
// @Tags         me
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.Response
// @Failure      401 {object} dto.ErrorResponse
// @Router       /me [get]
func (h *UserHandler) Me(c *fiber.Ctx) error {
	claims, ok := ssoclient.GetClaims(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "no claims")
	}

	return resp.OK(c, claims)
}

// FindFromCore godoc
// @Summary      Find user from Core system
// @Tags         user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body ssoclient.ReqFind true "Search payload"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /user/find-from-core [post]
func (h *UserHandler) FindFromCore(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[ssoclient.ReqFind](c)
	if !ok {
		return nil
	}
	// SSO client-ээр Core руу шууд дуудна
	out, err := ssoclient.FindUserFromCore(c.UserContext(), req, h.Cfg, h.Log)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	_, uerr := h.Service.User.Create(c.UserContext(), dto.UserCreateDto{
		Id:         out.Id,
		CivilId:    out.CivilId,
		RegNo:      out.RegNo,
		FamilyName: out.FamilyName,
		LastName:   out.LastName,
		FirstName:  out.FirstName,
		Gender:     out.Gender,
		BirthDate:  out.BirthDate,
		PhoneNo:    out.PhoneNo,
		Email:      out.Email,
	})
	if uerr != nil {
		return resp.InternalServerError(c, uerr.Error())
	}

	return resp.OK(c, out)
}

// List godoc
// @Summary      List users
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Param        page query int false "Page number (>=1)"
// @Param        size query int false "Page size"
// @Param        search query string false "JSON search (first_name,last_name,reg_no,phone_no,...)"
// @Param        sort query string false "JSON sort"
// @Param        createdFrom query string false "Created from (YYYY-MM-DD)"
// @Param        createdTo query string false "Created to (YYYY-MM-DD)"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /user [get]
func (h *UserHandler) List(c *fiber.Ctx) error {
	p, ok := resp.QueryBindAndValidate[common.PaginationQuery](c)
	if !ok {
		return nil
	}
	items, total, page, size, err := h.Service.User.List(c.UserContext(), p)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size) // <- Paginated-г хэрэглэж байна
}

// Create godoc
// @Summary      Create user
// @Tags         user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.UserCreateDto true "User data"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /user [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.UserCreateDto](c)
	if !ok {
		return nil
	}
	out, err := h.Service.User.Create(c.UserContext(), req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, out)
}

// Update godoc
// @Summary      Update user
// @Tags         user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Param        body body dto.UserUpdateDto true "User data"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /user/{id} [put]
func (h *UserHandler) Update(c *fiber.Ctx) error {

	req, ok := resp.BodyBindAndValidate[dto.UserUpdateDto](c)
	if !ok {
		return nil
	}
	out, err := h.Service.User.Update(c.UserContext(), req)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, out)
}

// Delete godoc
// @Summary      Delete user
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /user/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	params, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	out, err := h.Service.User.Delete(c.UserContext(), params.ID)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, out)
}

// Profile godoc
// @Summary      Get user profile with organizations
// @Tags         me
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.Response
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /me/profile [get]
func (h *UserHandler) Profile(c *fiber.Ctx) error {
	claims, ok := ssoclient.GetClaims(c)
	if !ok {
		return resp.Unauthorized(c)
	}

	// if claims.IsOrg {
	// 	var (
	// 		org domain.Organization
	// 		err error
	// 	)
	// 	// Org profile — Core-оос HTTPX-ээр татна
	// 	// org, err := h.fetchOrgFromCore(c, claims.OrgID)
	// 	if err != nil {
	// 		return resp.InternalServerError(c, err.Error())
	// 	}
	// 	return resp.OK(c, fiber.Map{
	// 		"is_org": true,
	// 		"org":    org,
	// 	})
	// }

	// Citizen profile: эхлээд DB; байхгүй бол Core
	u, err := h.Service.User.GetByID(c.UserContext(), claims.CitizenID)
	if err != nil {
		out, ferr := ssoclient.FindUserFromCore(c.UserContext(),
			ssoclient.ReqFind{SearchText: fmt.Sprintf("%d", claims.CitizenID)},
			h.Cfg, h.Log,
		)
		if ferr != nil {
			return resp.InternalServerError(c, ferr.Error())
		}
		return resp.OK(c, fiber.Map{
			"is_org": false,
			"user":   out,
		})
	}
	return resp.OK(c, fiber.Map{
		"is_org": false,
		"user":   u,
	})
}

// ProfileSSO godoc
// @Summary      Get profile from SSO
// @Tags         me
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.Response
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /me/profile/sso [get]
func (h *UserHandler) ProfileSSO(c *fiber.Ctx) error {
	res, err := ssoclient.GetUserProfileSSO(c.UserContext(), h.Cfg, h.Log)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, res)
}

// Organizations godoc
// @Summary      Get user's organizations list
// @Tags         me
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.Response
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /me/organizations [get]
func (h *UserHandler) Organizations(c *fiber.Ctx) error {
	claims, ok := ssoclient.GetClaims(c)
	if !ok {
		return resp.Unauthorized(c)
	}
	fields := []string{"id", "name", "reg_no", "parent_id", "type_id"}
	orgID, org, items, err := h.Service.User.Organizations(c.UserContext(), claims.CitizenID, claims.OrgID, fields)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, fiber.Map{
		"org_id": orgID,
		"org":    org,
		"items":  items,
	})
}
