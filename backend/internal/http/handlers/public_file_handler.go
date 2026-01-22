// Package handlers provides implementation for handlers
//
// File: public_file_handler.go
// Description: implementation for handlers
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package handlers

import (
	"templatev25/internal/http/dto"

	"templatev25/internal/app"

	"templatev25/internal/service"
	"git.gerege.mn/backend-packages/resp"

	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	*app.Dependencies
}

func NewFileHandler(d *app.Dependencies) *FileHandler {
	return &FileHandler{Dependencies: d}
}

// GET /file/:uuid  -> файл serve хийх
func (h *FileHandler) GetFile(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		return resp.InternalServerError(c, "uuid is required")
	}
	// service талд тусгайлсан хэрэгсэл байхгүй тул шууд дискенээс уншина (хуучин логиктой адил)
	return c.SendFile(service.PublicImageDir + "/" + uuid)
}

// POST /file/upload (multipart/form-data)
func (h *FileHandler) Upload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	// description
	desc := ""
	if v, ok := form.Value["description"]; ok && len(v) > 0 {
		desc = v[0]
	}
	// old name (хуучныг солих)
	oldName := ""
	if v, ok := form.Value["name"]; ok && len(v) > 0 {
		oldName = v[0]
	}

	// file header
	files, ok := form.File["file"]
	if !ok || len(files) == 0 {
		return resp.InternalServerError(c, "file is required")
	}
	if len(files) != 1 {
		return resp.InternalServerError(c, "only one file is allowed")
	}

	created, err := h.Service.PublicFile.Upload(c.UserContext(), files[0], desc, oldName)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c, created.FileUrl)
}

// GET /file/list
func (h *FileHandler) GetPublicFileList(c *fiber.Ctx) error {
	q, ok := resp.ParamsBindAndValidate[dto.PublicFileListQuery](c)
	if !ok {
		return nil
	}
	items, total, page, size, err := h.Service.PublicFile.List(c.UserContext(), q)
	if err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.Paginated(c, items, total, page, size)
}

// DELETE /file  (body: { "id": number })
func (h *FileHandler) DeletePublicFile(c *fiber.Ctx) error {
	req, ok := resp.BodyBindAndValidate[dto.PublicFileDeleteDto](c)
	if !ok {
		return nil
	}
	if err := h.Service.PublicFile.Delete(c.UserContext(), req.Name); err != nil {
		return resp.InternalServerError(c, err.Error())
	}
	return resp.OK(c)
}
