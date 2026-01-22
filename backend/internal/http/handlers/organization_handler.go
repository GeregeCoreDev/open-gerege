// Package handlers provides implementation for handlers
//
// File: organization_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"strings"

	"templatev25/internal/app"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/resp"
	ssoclient "git.gerege.mn/backend-packages/sso-client"

	"github.com/gofiber/fiber/v2"
)

type OrganizationHandler struct {
	*app.Dependencies
}

func NewOrganizationHandler(d *app.Dependencies) *OrganizationHandler {
	return &OrganizationHandler{Dependencies: d}
}

// FindFromCore godoc
// @Summary      Find organization from Core
// @Description  Search organization from SSO Core service
// @Tags         organization
// @Security     BearerAuth
// @Produce      json
// @Param        search_text query string true "Search text (reg_no or name)"
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /organization/find [get]
func (h *OrganizationHandler) FindFromCore(c *fiber.Ctx) error {
	req, ok := resp.QueryBindAndValidate[ssoclient.ReqFind](c)
	if !ok {
		return nil
	}

	// SSO client-ээр Core руу шууд дуудна
	out, err := ssoclient.FindOrganizationFromCore(c.UserContext(), req, h.Cfg, h.Log)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}

	return resp.OK(c, out)
}

// List godoc
// @Summary      List organizations
// @Description  Get paginated list of organizations
// @Tags         organization
// @Security     BearerAuth
// @Produce      json
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Router       /organization [get]
func (h *OrganizationHandler) List(c *fiber.Ctx) error {
	p, ok := resp.ParamsBindAndValidate[common.PaginationQuery](c)
	if !ok {
		return nil
	}
	items, total, page, size, err := h.Service.Organization.List(c.UserContext(), p)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Create godoc
// @Summary      Create organization
// @Description  Create a new organization
// @Tags         organization
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.OrganizationDto true "Organization data"
// @Success      200 {object} map[string]interface{}
// @Router       /organization [post]
func (h *OrganizationHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.OrganizationDto](c)
	if !ok {
		return nil
	}
	out, err := h.Service.Organization.Create(c.UserContext(), req)
	if err != nil {

		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, out)
}

// Update godoc
// @Summary      Update organization
// @Description  Update an existing organization
// @Tags         organization
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path int true "Organization ID"
// @Param        body body dto.OrganizationUpdateDto true "Organization data"
// @Success      200 {object} map[string]interface{}
// @Router       /organization/{id} [put]
func (h *OrganizationHandler) Update(c *fiber.Ctx) error {
	idParam, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	req, ok := resp.BodyBindAndValidate[dto.OrganizationUpdateDto](c)
	if !ok {
		return nil
	}
	out, err := h.Service.Organization.Update(c.UserContext(), idParam.ID, req)
	if err != nil {

		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, out)
}

// Delete godoc
// @Summary      Delete organization
// @Description  Delete an organization (soft delete)
// @Tags         organization
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Organization ID"
// @Success      200 {object} map[string]interface{}
// @Router       /organization/{id} [delete]
func (h *OrganizationHandler) Delete(c *fiber.Ctx) error {
	idParam, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	if err := h.Service.Organization.Delete(c.UserContext(), idParam.ID); err != nil {

		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Tree godoc
// @Summary      Get organization tree
// @Description  Get hierarchical organization tree
// @Tags         organization
// @Security     BearerAuth
// @Produce      json
// @Param        org_id query int false "Organization ID"
// @Success      200 {object} map[string]interface{}
// @Router       /organization/tree [get]
func (h *OrganizationHandler) Tree(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.OrganizationTreeQuery](c)
	if !ok {
		return nil
	}

	items, err := h.Service.Organization.Tree(c.UserContext(), q.OrgId)
	if err != nil {

		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, items)
}

type OrganizationTypeHandler struct {
	*app.Dependencies
}

func NewOrganizationTypeHandler(d *app.Dependencies) *OrganizationTypeHandler {
	return &OrganizationTypeHandler{Dependencies: d}
}

// List godoc
// @Summary      List organization types
// @Description  Get paginated list of organization types
// @Tags         orgtype
// @Security     BearerAuth
// @Produce      json
// @Param        page query int false "Page number"
// @Param        size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Router       /orgtype [get]
func (h *OrganizationTypeHandler) List(c *fiber.Ctx) error {
	p, ok := resp.QueryBindAndValidate[common.PaginationQuery](c)
	if !ok {
		return nil
	}
	items, total, page, size, err := h.Service.OrganizationType.List(c.UserContext(), p)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Create godoc
// @Summary      Create organization type
// @Tags         orgtype
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.OrganizationTypeDto true "payload"
// @Success      200 {object} map[string]interface{}
// @Router       /orgtype [post]
func (h *OrganizationTypeHandler) Create(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.OrganizationTypeDto](c)
	if !ok {
		return nil
	}
	err := h.Service.OrganizationType.Create(c.UserContext(), req)
	if err != nil {

		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Update godoc
// @Summary      Update organization type
// @Tags         orgtype
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path int true "ID"
// @Param        body body dto.OrganizationTypeDto true "payload"
// @Success      200 {object} map[string]interface{}
// @Router       /orgtype/{id} [put]
func (h *OrganizationTypeHandler) Update(c *fiber.Ctx) error {
	idp, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	req, ok := resp.BodyBindAndValidate[dto.OrganizationTypeDto](c)
	if !ok {
		return nil
	}
	err := h.Service.OrganizationType.Update(c.UserContext(), idp.ID, req)
	if err != nil {

		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Delete godoc
// @Summary      Delete organization type
// @Tags         orgtype
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} map[string]interface{}
// @Router       /orgtype/{id} [delete]
func (h *OrganizationTypeHandler) Delete(c *fiber.Ctx) error {
	idp, ok := resp.ParamsBindAndValidate[common.ID](c)
	if !ok {
		return nil
	}
	if err := h.Service.OrganizationType.Delete(c.UserContext(), idp.ID); err != nil {

		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Systems godoc
// @Summary      Get systems by organization type
// @Tags         orgtype
// @Security     BearerAuth
// @Produce      json
// @Param        id query int true "Organization Type ID"
// @Success      200 {object} map[string]interface{}
// @Router       /orgtype/system [get]
func (h *OrganizationTypeHandler) Systems(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.OrgTypeSystemsQuery](c)
	if !ok {
		return nil
	}
	items, err := h.Service.OrganizationType.Systems(c.UserContext(), q.TypeID)
	if err != nil {

		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, items)
}

// AddSystems godoc
// @Summary      Add systems to organization type
// @Tags         orgtype
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.OrgTypeAddSystemsDto true "payload"
// @Success      200 {object} map[string]interface{}
// @Router       /orgtype/system [post]
func (h *OrganizationTypeHandler) AddSystems(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.OrgTypeAddSystemsDto](c)
	if !ok {
		return nil
	}
	if err := h.Service.OrganizationType.AddSystems(c.UserContext(), req.TypeID, req.SystemIDs); err != nil {

		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Roles godoc
// @Summary      Get roles by organization type
// @Tags         orgtype
// @Security     BearerAuth
// @Produce      json
// @Param        type_id query int true "Organization Type ID"
// @Success      200 {object} map[string]interface{}
// @Router       /orgtype/role [get]
func (h *OrganizationTypeHandler) Roles(c *fiber.Ctx) error {
	q, ok := resp.QueryBindAndValidate[dto.OrgTypeRolesQuery](c)
	if !ok {
		return nil
	}
	items, err := h.Service.OrganizationType.Roles(c.UserContext(), q.TypeID)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, items)
}

// AddRoles godoc
// @Summary      Add roles to organization type
// @Tags         orgtype
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.OrgTypeRolesAddDto true "payload"
// @Success      200 {object} map[string]interface{}
// @Router       /orgtype/role [post]
func (h *OrganizationTypeHandler) AddRoles(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.OrgTypeRolesAddDto](c)
	if !ok {
		return nil
	}
	if err := h.Service.OrganizationType.AddRoles(c.UserContext(), req.TypeID, req.RoleIDs); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

type OrgUserHandler struct{ *app.Dependencies }

func NewOrgUserHandler(d *app.Dependencies) *OrgUserHandler {
	return &OrgUserHandler{Dependencies: d}
}

// List godoc
// @Summary      List organization users
// @Tags         orguser
// @Security     BearerAuth
// @Produce      json
// @Param        org_id  query int false "Organization ID"
// @Param        user_id query int false "User ID"
// @Param        page    query int false "Page number"
// @Param        size    query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Router       /orguser [get]
func (h *OrgUserHandler) List(c *fiber.Ctx) error {
	q, ok := resp.ParamsBindAndValidate[dto.OrgUserListQuery](c)
	if !ok {
		return nil
	}

	// default org_id токеноос
	if q.OrgId == 0 && q.UserId == 0 {
		if claims, ok := ssoclient.GetClaims(c); ok {
			q.OrgId = claims.OrgID
		}
	}
	items, total, page, size, err := h.Service.OrgUser.List(c.UserContext(), q)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Add godoc
// @Summary      Add user to organization
// @Tags         orguser
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.OrgUserCreateDto true "payload"
// @Success      200 {object} map[string]interface{}
// @Router       /orguser [post]
func (h *OrgUserHandler) Add(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.OrgUserCreateDto](c)
	if !ok {
		return nil
	}
	authHeader := c.Get(fiber.HeaderAuthorization)
	if err := h.Service.OrgUser.Add(c.UserContext(), req, authHeader); err != nil {
		msg := err.Error()
		if strings.Contains(msg, "duplicate") {
			return resp.InternalServerError(c, "Хэрэглэгч аль хэдийн бүртгэгдсэн байна")
		}
		return resp.InternalServerError(c, msg)
	}
	return resp.OK(c)
}

// Remove godoc
// @Summary      Remove user from organization
// @Tags         orguser
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body dto.OrgUserDeleteDto true "payload"
// @Success      200 {object} map[string]interface{}
// @Router       /orguser [delete]
func (h *OrgUserHandler) Remove(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.OrgUserDeleteDto](c)
	if !ok {
		return nil
	}
	if err := h.Service.OrgUser.Remove(c.UserContext(), req); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}

// Users godoc
// @Summary      Get users by organization
// @Tags         orguser
// @Security     BearerAuth
// @Produce      json
// @Param        org_id query int  false "Organization ID"
// @Param        name   query string false "Filter by name"
// @Param        page   query int  false "Page number"
// @Param        size   query int  false "Page size"
// @Success      200 {object} map[string]interface{}
// @Router       /orguser/users [get]
func (h *OrgUserHandler) Users(c *fiber.Ctx) error {
	// org_id default токеноос
	orgId := c.QueryInt("org_id")
	if orgId == 0 {
		if claims, ok := ssoclient.GetClaims(c); ok {
			orgId = claims.OrgID
		}
	}
	p, ok := resp.ParamsBindAndValidate[common.PaginationQuery](c)
	if !ok {
		return nil
	}

	name := strings.TrimSpace(c.Query("name"))
	items, total, page, size, err := h.Service.OrgUser.UsersByOrg(c.UserContext(), orgId, name, p)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// Orgs godoc
// @Summary      Get organizations by user
// @Tags         orguser
// @Security     BearerAuth
// @Produce      json
// @Param        user_id query int  false "User ID"
// @Param        name    query string false "Filter by name"
// @Param        page    query int  false "Page number"
// @Param        size    query int  false "Page size"
// @Success      200 {object} map[string]interface{}
// @Router       /orguser/organizations [get]
func (h *OrgUserHandler) Orgs(c *fiber.Ctx) error {
	// user_id заагдаагүй бол токеноос (анхны кодонд энд orgID-г токеноос авч байсан нь алдаа байсан байж магад)
	userId := c.QueryInt("user_id")
	if userId == 0 {
		if claims, ok := ssoclient.GetClaims(c); ok {
			userId = claims.UserID
		}
	}
	p, ok := resp.ParamsBindAndValidate[common.PaginationQuery](c)
	if !ok {
		return nil
	}

	name := strings.TrimSpace(c.Query("name"))
	items, total, page, size, err := h.Service.OrgUser.OrgsByUser(c.UserContext(), userId, name, p)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}
