//go:build integration

// Package integration provides integration tests for HTTP handlers
//
// File: role_handler_test.go
// Description: Integration tests for RoleHandler
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockRoleService implements role service interface for testing
type mockRoleService struct {
	mock.Mock
}

func (m *mockRoleService) ListFilteredPaged(ctx context.Context, p dto.RoleListQuery) ([]domain.Role, int64, int, int, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.Role), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockRoleService) Create(ctx context.Context, req dto.RoleCreateDto) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockRoleService) Update(ctx context.Context, id int, req dto.RoleUpdateDto) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *mockRoleService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRoleService) GetPermissions(ctx context.Context, q dto.RolePermissionsQuery) ([]domain.Permission, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Permission), args.Error(1)
}

func (m *mockRoleService) SetPermissions(ctx context.Context, req dto.RolePermissionsUpdateDto) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// setupRoleTestApp creates a test Fiber app with role routes
func setupRoleTestApp(svc *mockRoleService) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"code":    "error",
				"msg":     err.Error(),
			})
		},
	})

	// Role routes
	role := app.Group("/api/v1/role")

	// List roles
	role.Get("/", func(c *fiber.Ctx) error {
		p := dto.RoleListQuery{}
		if err := c.QueryParser(&p); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid query parameters")
		}

		items, total, page, size, err := svc.ListFilteredPaged(c.UserContext(), p)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"code":    "ok",
			"data": fiber.Map{
				"items": items,
				"total": total,
				"page":  page,
				"size":  size,
			},
		})
	})

	// Create role
	role.Post("/", func(c *fiber.Ctx) error {
		var req dto.RoleCreateDto
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		// Basic validation
		if req.Code == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "code is required")
		}
		if req.Name == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "name is required")
		}
		if req.SystemID <= 0 {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "system_id is required")
		}

		if err := svc.Create(c.UserContext(), req); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"code":    "created",
		})
	})

	// Update role
	role.Put("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid role ID")
		}

		var req dto.RoleUpdateDto
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		if err := svc.Update(c.UserContext(), id, req); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"code":    "ok",
		})
	})

	// Delete role
	role.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid role ID")
		}

		if err := svc.Delete(c.UserContext(), id); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"code":    "ok",
		})
	})

	// Get role permissions
	role.Get("/permissions", func(c *fiber.Ctx) error {
		roleID := c.QueryInt("role_id")
		if roleID <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "role_id is required")
		}

		q := dto.RolePermissionsQuery{RoleID: roleID}
		perms, err := svc.GetPermissions(c.UserContext(), q)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"code":    "ok",
			"data":    perms,
		})
	})

	// Set role permissions
	role.Post("/permissions", func(c *fiber.Ctx) error {
		var req dto.RolePermissionsUpdateDto
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		if req.RoleID <= 0 {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "role_id is required")
		}

		if err := svc.SetPermissions(c.UserContext(), req); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"code":    "created",
		})
	})

	return app
}

func TestRoleHandler_List(t *testing.T) {
	isActive := true
	tests := []struct {
		name       string
		query      string
		mockSetup  func(*mockRoleService)
		wantStatus int
		wantCount  int
	}{
		{
			name:  "success - returns roles",
			query: "",
			mockSetup: func(m *mockRoleService) {
				roles := []domain.Role{
					{ID: 1, Code: "admin", Name: "Administrator", IsActive: &isActive},
					{ID: 2, Code: "user", Name: "User", IsActive: &isActive},
				}
				m.On("ListFilteredPaged", mock.Anything, mock.AnythingOfType("dto.RoleListQuery")).
					Return(roles, int64(2), 1, 10, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:  "success - empty list",
			query: "",
			mockSetup: func(m *mockRoleService) {
				m.On("ListFilteredPaged", mock.Anything, mock.AnythingOfType("dto.RoleListQuery")).
					Return([]domain.Role{}, int64(0), 1, 10, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name:  "success - with pagination",
			query: "?page=2&size=5",
			mockSetup: func(m *mockRoleService) {
				roles := []domain.Role{
					{ID: 6, Code: "role6", Name: "Role 6"},
				}
				m.On("ListFilteredPaged", mock.Anything, mock.AnythingOfType("dto.RoleListQuery")).
					Return(roles, int64(6), 2, 5, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:  "success - filter by system_id",
			query: "?system_id=1",
			mockSetup: func(m *mockRoleService) {
				roles := []domain.Role{
					{ID: 1, Code: "admin", Name: "Administrator", SystemID: 1},
				}
				m.On("ListFilteredPaged", mock.Anything, mock.AnythingOfType("dto.RoleListQuery")).
					Return(roles, int64(1), 1, 10, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockRoleService{}
			tt.mockSetup(mockSvc)

			app := setupRoleTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/role/"+tt.query, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			var result map[string]interface{}
			json.Unmarshal(body, &result)

			assert.True(t, result["success"].(bool))
			if data, ok := result["data"].(map[string]interface{}); ok {
				items := data["items"].([]interface{})
				assert.Len(t, items, tt.wantCount)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestRoleHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockSetup  func(*mockRoleService)
		wantStatus int
	}{
		{
			name: "success - role created",
			body: `{"code": "manager", "name": "Manager", "system_id": 1, "description": "Manager role"}`,
			mockSetup: func(m *mockRoleService) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(req dto.RoleCreateDto) bool {
					return req.Code == "manager" && req.Name == "Manager" && req.SystemID == 1
				})).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "error - missing code",
			body:       `{"name": "Manager", "system_id": 1}`,
			mockSetup:  func(m *mockRoleService) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "error - missing name",
			body:       `{"code": "manager", "system_id": 1}`,
			mockSetup:  func(m *mockRoleService) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "error - missing system_id",
			body:       `{"code": "manager", "name": "Manager"}`,
			mockSetup:  func(m *mockRoleService) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "error - invalid JSON",
			body:       `{invalid}`,
			mockSetup:  func(m *mockRoleService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockRoleService{}
			tt.mockSetup(mockSvc)

			app := setupRoleTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/role/", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestRoleHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		roleID     string
		body       string
		mockSetup  func(*mockRoleService)
		wantStatus int
	}{
		{
			name:   "success - role updated",
			roleID: "1",
			body:   `{"code": "admin_updated", "name": "Admin Updated", "system_id": 1}`,
			mockSetup: func(m *mockRoleService) {
				m.On("Update", mock.Anything, 1, mock.AnythingOfType("dto.RoleUpdateDto")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "error - invalid role ID",
			roleID:     "invalid",
			body:       `{"code": "admin", "name": "Admin", "system_id": 1}`,
			mockSetup:  func(m *mockRoleService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error - invalid JSON",
			roleID:     "1",
			body:       `{invalid}`,
			mockSetup:  func(m *mockRoleService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockRoleService{}
			tt.mockSetup(mockSvc)

			app := setupRoleTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodPut, "/api/v1/role/"+tt.roleID, bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestRoleHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		roleID     string
		mockSetup  func(*mockRoleService)
		wantStatus int
	}{
		{
			name:   "success - role deleted",
			roleID: "1",
			mockSetup: func(m *mockRoleService) {
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "error - invalid role ID",
			roleID:     "invalid",
			mockSetup:  func(m *mockRoleService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockRoleService{}
			tt.mockSetup(mockSvc)

			app := setupRoleTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/role/"+tt.roleID, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestRoleHandler_GetPermissions(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		mockSetup  func(*mockRoleService)
		wantStatus int
		wantCount  int
	}{
		{
			name:  "success - returns permissions",
			query: "?role_id=1",
			mockSetup: func(m *mockRoleService) {
				perms := []domain.Permission{
					{ID: 1, Code: "admin.user.read", Name: "Read Users"},
					{ID: 2, Code: "admin.user.write", Name: "Write Users"},
				}
				m.On("GetPermissions", mock.Anything, dto.RolePermissionsQuery{RoleID: 1}).
					Return(perms, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:       "error - missing role_id",
			query:      "",
			mockSetup:  func(m *mockRoleService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error - invalid role_id",
			query:      "?role_id=0",
			mockSetup:  func(m *mockRoleService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockRoleService{}
			tt.mockSetup(mockSvc)

			app := setupRoleTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/role/permissions"+tt.query, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantStatus == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				json.Unmarshal(body, &result)

				assert.True(t, result["success"].(bool))
				data := result["data"].([]interface{})
				assert.Len(t, data, tt.wantCount)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestRoleHandler_SetPermissions(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockSetup  func(*mockRoleService)
		wantStatus int
	}{
		{
			name: "success - permissions set",
			body: `{"role_id": 1, "permission_ids": [1, 2, 3]}`,
			mockSetup: func(m *mockRoleService) {
				m.On("SetPermissions", mock.Anything, mock.MatchedBy(func(req dto.RolePermissionsUpdateDto) bool {
					return req.RoleID == 1 && len(req.PermissionIDs) == 3
				})).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "success - empty permissions (remove all)",
			body: `{"role_id": 1, "permission_ids": []}`,
			mockSetup: func(m *mockRoleService) {
				m.On("SetPermissions", mock.Anything, mock.AnythingOfType("dto.RolePermissionsUpdateDto")).
					Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "error - missing role_id",
			body:       `{"permission_ids": [1, 2, 3]}`,
			mockSetup:  func(m *mockRoleService) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "error - invalid JSON",
			body:       `{invalid}`,
			mockSetup:  func(m *mockRoleService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockRoleService{}
			tt.mockSetup(mockSvc)

			app := setupRoleTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/role/permissions", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

// TestRoleHandler_ContentType tests that responses have correct content type
func TestRoleHandler_ContentType(t *testing.T) {
	mockSvc := &mockRoleService{}
	mockSvc.On("ListFilteredPaged", mock.Anything, mock.AnythingOfType("dto.RoleListQuery")).
		Return([]domain.Role{}, int64(0), 1, 10, nil)

	app := setupRoleTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/role/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json")
}

// TestRoleHandler_NotFound tests 404 for non-existent routes
func TestRoleHandler_NotFound(t *testing.T) {
	mockSvc := &mockRoleService{}
	app := setupRoleTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/role/non-existent/path", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
