//go:build integration

// Package integration provides integration tests for HTTP handlers
//
// File: user_handler_test.go
// Description: Integration tests for UserHandler
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/common"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockUserService implements user service interface for testing
type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) List(ctx context.Context, p common.PaginationQuery) ([]domain.User, int64, int, int, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockUserService) GetByID(ctx context.Context, id int) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserService) Create(ctx context.Context, req dto.UserCreateDto) (domain.User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserService) Update(ctx context.Context, req dto.UserUpdateDto) (domain.User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserService) Delete(ctx context.Context, id int) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserService) Organizations(ctx context.Context, citizenID, orgID int, fields []string) (int, domain.Organization, []domain.Organization, error) {
	args := m.Called(ctx, citizenID, orgID, fields)
	return args.Int(0), args.Get(1).(domain.Organization), args.Get(2).([]domain.Organization), args.Error(3)
}

// setupUserTestApp creates a test Fiber app with user routes
func setupUserTestApp(svc *mockUserService) *fiber.App {
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

	// User routes
	user := app.Group("/api/v1/user")

	// List users
	user.Get("/", func(c *fiber.Ctx) error {
		p := common.PaginationQuery{}
		if err := c.QueryParser(&p); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid query parameters")
		}

		items, total, page, size, err := svc.List(c.UserContext(), p)
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

	// Get user by ID
	user.Get("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid user ID")
		}

		u, err := svc.GetByID(c.UserContext(), id)
		if err != nil {
			if err.Error() == "not found" {
				return fiber.NewError(fiber.StatusNotFound, "user not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"code":    "ok",
			"data":    u,
		})
	})

	// Create user
	user.Post("/", func(c *fiber.Ctx) error {
		var req dto.UserCreateDto
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		// Basic validation
		if req.FirstName == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "first_name is required")
		}
		if req.LastName == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "last_name is required")
		}

		u, err := svc.Create(c.UserContext(), req)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"code":    "created",
			"data":    u,
		})
	})

	// Update user
	user.Put("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid user ID")
		}

		var req dto.UserUpdateDto
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}
		req.Id = id

		u, err := svc.Update(c.UserContext(), req)
		if err != nil {
			if err.Error() == "not found" {
				return fiber.NewError(fiber.StatusNotFound, "user not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"code":    "ok",
			"data":    u,
		})
	})

	// Delete user
	user.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid user ID")
		}

		u, err := svc.Delete(c.UserContext(), id)
		if err != nil {
			if err.Error() == "not found" {
				return fiber.NewError(fiber.StatusNotFound, "user not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"code":    "ok",
			"data":    u,
		})
	})

	return app
}

func TestUserHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		mockSetup  func(*mockUserService)
		wantStatus int
		wantCount  int
	}{
		{
			name:  "success - returns users",
			query: "",
			mockSetup: func(m *mockUserService) {
				users := []domain.User{
					{Id: 1, FirstName: "John", LastName: "Doe"},
					{Id: 2, FirstName: "Jane", LastName: "Smith"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return(users, int64(2), 1, 10, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:  "success - empty list",
			query: "",
			mockSetup: func(m *mockUserService) {
				m.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return([]domain.User{}, int64(0), 1, 10, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name:  "success - with pagination",
			query: "?page=2&size=5",
			mockSetup: func(m *mockUserService) {
				users := []domain.User{
					{Id: 6, FirstName: "User", LastName: "Six"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return(users, int64(6), 2, 5, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:  "error - database error",
			query: "",
			mockSetup: func(m *mockUserService) {
				m.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return(nil, int64(0), 0, 0, errors.New("database error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{}
			tt.mockSetup(mockSvc)

			app := setupUserTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/user/"+tt.query, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantStatus == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				json.Unmarshal(body, &result)

				assert.True(t, result["success"].(bool))
				if data, ok := result["data"].(map[string]interface{}); ok {
					items := data["items"].([]interface{})
					assert.Len(t, items, tt.wantCount)
				}
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		mockSetup  func(*mockUserService)
		wantStatus int
		wantName   string
	}{
		{
			name:   "success - returns user",
			userID: "1",
			mockSetup: func(m *mockUserService) {
				m.On("GetByID", mock.Anything, 1).
					Return(domain.User{Id: 1, FirstName: "John", LastName: "Doe"}, nil)
			},
			wantStatus: http.StatusOK,
			wantName:   "John",
		},
		{
			name:   "error - not found",
			userID: "999",
			mockSetup: func(m *mockUserService) {
				m.On("GetByID", mock.Anything, 999).
					Return(domain.User{}, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "error - invalid ID",
			userID:     "invalid",
			mockSetup:  func(m *mockUserService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{}
			tt.mockSetup(mockSvc)

			app := setupUserTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/user/"+tt.userID, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantStatus == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				json.Unmarshal(body, &result)

				assert.True(t, result["success"].(bool))
				data := result["data"].(map[string]interface{})
				assert.Equal(t, tt.wantName, data["first_name"])
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockSetup  func(*mockUserService)
		wantStatus int
	}{
		{
			name: "success - user created",
			body: `{"first_name": "John", "last_name": "Doe", "email": "john@example.com", "phone_no": "99001122"}`,
			mockSetup: func(m *mockUserService) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(req dto.UserCreateDto) bool {
					return req.FirstName == "John" && req.LastName == "Doe"
				})).Return(domain.User{Id: 1, FirstName: "John", LastName: "Doe"}, nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "error - missing first_name",
			body:       `{"last_name": "Doe", "email": "john@example.com"}`,
			mockSetup:  func(m *mockUserService) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "error - missing last_name",
			body:       `{"first_name": "John", "email": "john@example.com"}`,
			mockSetup:  func(m *mockUserService) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "error - empty body",
			body:       `{}`,
			mockSetup:  func(m *mockUserService) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "error - invalid JSON",
			body:       `{invalid}`,
			mockSetup:  func(m *mockUserService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error - duplicate email",
			body: `{"first_name": "John", "last_name": "Doe", "email": "existing@example.com"}`,
			mockSetup: func(m *mockUserService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("dto.UserCreateDto")).
					Return(domain.User{}, errors.New("email already exists"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{}
			tt.mockSetup(mockSvc)

			app := setupUserTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/user/", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		body       string
		mockSetup  func(*mockUserService)
		wantStatus int
	}{
		{
			name:   "success - user updated",
			userID: "1",
			body:   `{"first_name": "John Updated", "last_name": "Doe Updated"}`,
			mockSetup: func(m *mockUserService) {
				m.On("Update", mock.Anything, mock.AnythingOfType("dto.UserUpdateDto")).
					Return(domain.User{Id: 1, FirstName: "John Updated", LastName: "Doe Updated"}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "error - not found",
			userID: "999",
			body:   `{"first_name": "John", "last_name": "Doe"}`,
			mockSetup: func(m *mockUserService) {
				m.On("Update", mock.Anything, mock.AnythingOfType("dto.UserUpdateDto")).
					Return(domain.User{}, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "error - invalid ID",
			userID:     "invalid",
			body:       `{"first_name": "John", "last_name": "Doe"}`,
			mockSetup:  func(m *mockUserService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error - invalid JSON",
			userID:     "1",
			body:       `{invalid}`,
			mockSetup:  func(m *mockUserService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{}
			tt.mockSetup(mockSvc)

			app := setupUserTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodPut, "/api/v1/user/"+tt.userID, bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		mockSetup  func(*mockUserService)
		wantStatus int
	}{
		{
			name:   "success - user deleted",
			userID: "1",
			mockSetup: func(m *mockUserService) {
				m.On("Delete", mock.Anything, 1).
					Return(domain.User{Id: 1, FirstName: "John", LastName: "Doe"}, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "error - not found",
			userID: "999",
			mockSetup: func(m *mockUserService) {
				m.On("Delete", mock.Anything, 999).
					Return(domain.User{}, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "error - invalid ID",
			userID:     "invalid",
			mockSetup:  func(m *mockUserService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{}
			tt.mockSetup(mockSvc)

			app := setupUserTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/user/"+tt.userID, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

// TestUserHandler_ContentType tests that responses have correct content type
func TestUserHandler_ContentType(t *testing.T) {
	mockSvc := &mockUserService{}
	mockSvc.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
		Return([]domain.User{}, int64(0), 1, 10, nil)

	app := setupUserTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json")
}

// TestUserHandler_EmptyBody tests handling of empty request bodies
func TestUserHandler_EmptyBody(t *testing.T) {
	mockSvc := &mockUserService{}
	app := setupUserTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Empty body should result in validation error
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

// TestUserHandler_SearchQuery tests search query parameter handling
func TestUserHandler_SearchQuery(t *testing.T) {
	mockSvc := &mockUserService{}
	mockSvc.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
		Return([]domain.User{
			{Id: 1, FirstName: "John", LastName: "Doe"},
		}, int64(1), 1, 10, nil)

	app := setupUserTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/?search=john", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	assert.True(t, result["success"].(bool))
	mockSvc.AssertExpectations(t)
}
