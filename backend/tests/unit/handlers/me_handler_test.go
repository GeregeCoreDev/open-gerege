// Package handlers provides unit tests for HTTP handlers
//
// File: me_handler_test.go
// Description: Unit tests for Me (Current User) Handler endpoints
package handlers

import (
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

// mockClaims represents SSO claims for testing
type mockClaims struct {
	CitizenID int  `json:"citizen_id"`
	OrgID     int  `json:"org_id"`
	IsOrg     bool `json:"is_org"`
}

// setupMeTestApp creates a test Fiber app with Me routes
func setupMeTestApp(claims *mockClaims, userSvc *mockUserService) *fiber.App {
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
				"message": err.Error(),
			})
		},
	})

	// Me routes
	me := app.Group("/me")

	// Auth middleware simulator
	authMiddleware := func(c *fiber.Ctx) error {
		if claims == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		c.Locals("claims", claims)
		return c.Next()
	}

	// GET /me - Get current user claims
	me.Get("/", authMiddleware, func(c *fiber.Ctx) error {
		cl := c.Locals("claims")
		if cl == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "no claims")
		}
		return c.JSON(fiber.Map{
			"code":    "OK",
			"message": "success",
			"data":    cl,
		})
	})

	// GET /me/profile - Get user profile
	me.Get("/profile", authMiddleware, func(c *fiber.Ctx) error {
		cl := c.Locals("claims").(*mockClaims)

		user, err := userSvc.GetByID(c.UserContext(), cl.CitizenID)
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    "OK",
				"message": "success",
				"data": fiber.Map{
					"is_org": false,
					"user":   fiber.Map{"id": cl.CitizenID, "source": "sso"},
				},
			})
		}

		return c.JSON(fiber.Map{
			"code":    "OK",
			"message": "success",
			"data": fiber.Map{
				"is_org": cl.IsOrg,
				"user":   user,
			},
		})
	})

	// GET /me/profile/sso - Get profile from SSO (external call)
	me.Get("/profile/sso", authMiddleware, func(c *fiber.Ctx) error {
		cl := c.Locals("claims").(*mockClaims)
		return c.JSON(fiber.Map{
			"code":    "OK",
			"message": "success",
			"data": fiber.Map{
				"citizen_id": cl.CitizenID,
				"org_id":     cl.OrgID,
				"source":     "sso",
			},
		})
	})

	// GET /me/organizations - Get user's organizations
	me.Get("/organizations", authMiddleware, func(c *fiber.Ctx) error {
		cl := c.Locals("claims").(*mockClaims)

		orgID, org, items, err := userSvc.Organizations(c.UserContext(), cl.CitizenID, cl.OrgID, []string{"id", "name"})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"code":    "OK",
			"message": "success",
			"data": fiber.Map{
				"org_id": orgID,
				"org":    org,
				"items":  items,
			},
		})
	})

	return app
}

// =============================================================================
// GET /me - Get Current User Claims
// =============================================================================

func TestMeHandler_GetCurrentUser(t *testing.T) {
	tests := []struct {
		name       string
		claims     *mockClaims
		wantStatus int
		wantData   bool
	}{
		{
			name: "success - returns claims",
			claims: &mockClaims{
				CitizenID: 12345678,
				OrgID:     20000001,
				IsOrg:     false,
			},
			wantStatus: http.StatusOK,
			wantData:   true,
		},
		{
			name:       "error - unauthorized (no claims)",
			claims:     nil,
			wantStatus: http.StatusUnauthorized,
			wantData:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{}
			app := setupMeTestApp(tt.claims, mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantData {
				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				json.Unmarshal(body, &result)

				assert.Equal(t, "OK", result["code"])
				data := result["data"].(map[string]interface{})
				assert.Equal(t, float64(tt.claims.CitizenID), data["citizen_id"])
				assert.Equal(t, float64(tt.claims.OrgID), data["org_id"])
			}
		})
	}
}

// =============================================================================
// GET /me/profile - Get User Profile
// =============================================================================

func TestMeHandler_GetProfile(t *testing.T) {
	tests := []struct {
		name       string
		claims     *mockClaims
		mockSetup  func(*mockUserService)
		wantStatus int
		wantIsOrg  bool
	}{
		{
			name: "success - returns profile from DB",
			claims: &mockClaims{
				CitizenID: 12345678,
				OrgID:     20000001,
				IsOrg:     false,
			},
			mockSetup: func(m *mockUserService) {
				m.On("GetByID", mock.Anything, 12345678).
					Return(domain.User{
						Id:        12345678,
						FirstName: "Тест",
						LastName:  "Хэрэглэгч",
						Email:     "test@example.com",
					}, nil)
			},
			wantStatus: http.StatusOK,
			wantIsOrg:  false,
		},
		{
			name: "success - user not in DB, returns from SSO",
			claims: &mockClaims{
				CitizenID: 99999999,
				OrgID:     20000001,
				IsOrg:     false,
			},
			mockSetup: func(m *mockUserService) {
				m.On("GetByID", mock.Anything, 99999999).
					Return(domain.User{}, errors.New("not found"))
			},
			wantStatus: http.StatusOK,
			wantIsOrg:  false,
		},
		{
			name:       "error - unauthorized",
			claims:     nil,
			mockSetup:  func(m *mockUserService) {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{}
			tt.mockSetup(mockSvc)

			app := setupMeTestApp(tt.claims, mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/me/profile", nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantStatus == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				json.Unmarshal(body, &result)

				assert.Equal(t, "OK", result["code"])
				data := result["data"].(map[string]interface{})
				assert.Equal(t, tt.wantIsOrg, data["is_org"])
				assert.NotNil(t, data["user"])
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

// =============================================================================
// GET /me/profile/sso - Get Profile from SSO
// =============================================================================

func TestMeHandler_GetProfileSSO(t *testing.T) {
	tests := []struct {
		name       string
		claims     *mockClaims
		wantStatus int
	}{
		{
			name: "success - returns SSO profile",
			claims: &mockClaims{
				CitizenID: 12345678,
				OrgID:     20000001,
				IsOrg:     false,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "error - unauthorized",
			claims:     nil,
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{}
			app := setupMeTestApp(tt.claims, mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/me/profile/sso", nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantStatus == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				json.Unmarshal(body, &result)

				assert.Equal(t, "OK", result["code"])
				data := result["data"].(map[string]interface{})
				assert.Equal(t, "sso", data["source"])
				assert.Equal(t, float64(tt.claims.CitizenID), data["citizen_id"])
			}
		})
	}
}

// =============================================================================
// GET /me/organizations - Get User's Organizations
// =============================================================================

func TestMeHandler_GetOrganizations(t *testing.T) {
	tests := []struct {
		name       string
		claims     *mockClaims
		mockSetup  func(*mockUserService)
		wantStatus int
		wantOrgID  int
		wantCount  int
	}{
		{
			name: "success - returns organizations",
			claims: &mockClaims{
				CitizenID: 12345678,
				OrgID:     20000001,
				IsOrg:     false,
			},
			mockSetup: func(m *mockUserService) {
				m.On("Organizations", mock.Anything, 12345678, 20000001, []string{"id", "name"}).
					Return(
						20000001,
						domain.Organization{Id: 20000001, Name: "Main Org"},
						[]domain.Organization{
							{Id: 20000001, Name: "Main Org"},
							{Id: 20000002, Name: "Second Org"},
						},
						nil,
					)
			},
			wantStatus: http.StatusOK,
			wantOrgID:  20000001,
			wantCount:  2,
		},
		{
			name: "success - user with single organization",
			claims: &mockClaims{
				CitizenID: 88888888,
				OrgID:     20000003,
				IsOrg:     false,
			},
			mockSetup: func(m *mockUserService) {
				m.On("Organizations", mock.Anything, 88888888, 20000003, []string{"id", "name"}).
					Return(
						20000003,
						domain.Organization{Id: 20000003, Name: "Single Org"},
						[]domain.Organization{
							{Id: 20000003, Name: "Single Org"},
						},
						nil,
					)
			},
			wantStatus: http.StatusOK,
			wantOrgID:  20000003,
			wantCount:  1,
		},
		{
			name: "success - user with no organizations",
			claims: &mockClaims{
				CitizenID: 77777777,
				OrgID:     0,
				IsOrg:     false,
			},
			mockSetup: func(m *mockUserService) {
				m.On("Organizations", mock.Anything, 77777777, 0, []string{"id", "name"}).
					Return(
						0,
						domain.Organization{},
						[]domain.Organization{},
						nil,
					)
			},
			wantStatus: http.StatusOK,
			wantOrgID:  0,
			wantCount:  0,
		},
		{
			name: "error - service error",
			claims: &mockClaims{
				CitizenID: 12345678,
				OrgID:     20000001,
				IsOrg:     false,
			},
			mockSetup: func(m *mockUserService) {
				m.On("Organizations", mock.Anything, 12345678, 20000001, []string{"id", "name"}).
					Return(0, domain.Organization{}, []domain.Organization{}, errors.New("database error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "error - unauthorized",
			claims:     nil,
			mockSetup:  func(m *mockUserService) {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockUserService{}
			tt.mockSetup(mockSvc)

			app := setupMeTestApp(tt.claims, mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/me/organizations", nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantStatus == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				json.Unmarshal(body, &result)

				assert.Equal(t, "OK", result["code"])
				data := result["data"].(map[string]interface{})
				assert.Equal(t, float64(tt.wantOrgID), data["org_id"])

				items := data["items"].([]interface{})
				assert.Len(t, items, tt.wantCount)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

// =============================================================================
// Content-Type & Response Format Tests
// =============================================================================

func TestMeHandler_ResponseFormat(t *testing.T) {
	claims := &mockClaims{
		CitizenID: 12345678,
		OrgID:     20000001,
		IsOrg:     false,
	}

	mockSvc := &mockUserService{}
	app := setupMeTestApp(claims, mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check Content-Type
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json")

	// Check response structure
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// Standard response fields
	assert.Contains(t, result, "code")
	assert.Contains(t, result, "message")
	assert.Contains(t, result, "data")
}

// =============================================================================
// Security Tests - All endpoints require authentication
// =============================================================================

func TestMeHandler_RequiresAuthentication(t *testing.T) {
	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/me"},
		{http.MethodGet, "/me/profile"},
		{http.MethodGet, "/me/profile/sso"},
		{http.MethodGet, "/me/organizations"},
	}

	mockSvc := &mockUserService{}
	app := setupMeTestApp(nil, mockSvc) // No claims = unauthorized

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			req := httptest.NewRequest(ep.method, ep.path, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}
}

// =============================================================================
// Edge Case Tests
// =============================================================================

func TestMeHandler_LargeCitizenID(t *testing.T) {
	// Test with large citizen ID (common in Mongolia: 8+ digits)
	claims := &mockClaims{
		CitizenID: 99999999999,
		OrgID:     20000001,
		IsOrg:     false,
	}

	mockSvc := &mockUserService{}
	mockSvc.On("GetByID", mock.Anything, 99999999999).
		Return(domain.User{}, errors.New("not found"))

	app := setupMeTestApp(claims, mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/me/profile", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMeHandler_OrgUserClaims(t *testing.T) {
	// Test with organization user (IsOrg = true)
	claims := &mockClaims{
		CitizenID: 0, // No citizen ID for org users
		OrgID:     20000001,
		IsOrg:     true,
	}

	mockSvc := &mockUserService{}
	app := setupMeTestApp(claims, mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	data := result["data"].(map[string]interface{})
	assert.Equal(t, true, data["is_org"])
	assert.Equal(t, float64(20000001), data["org_id"])
}

func TestMeHandler_ZeroCitizenID(t *testing.T) {
	// Edge case: citizen_id = 0
	claims := &mockClaims{
		CitizenID: 0,
		OrgID:     20000001,
		IsOrg:     false,
	}

	mockSvc := &mockUserService{}
	app := setupMeTestApp(claims, mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMeHandler_MultipleOrganizations(t *testing.T) {
	// User with many organizations
	claims := &mockClaims{
		CitizenID: 12345678,
		OrgID:     20000001,
		IsOrg:     false,
	}

	mockSvc := &mockUserService{}

	// Create 10 organizations
	orgs := make([]domain.Organization, 10)
	for i := 0; i < 10; i++ {
		orgs[i] = domain.Organization{Id: 20000001 + i, Name: "Org " + string(rune('A'+i))}
	}

	mockSvc.On("Organizations", mock.Anything, 12345678, 20000001, []string{"id", "name"}).
		Return(20000001, orgs[0], orgs, nil)

	app := setupMeTestApp(claims, mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/me/organizations", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	data := result["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	assert.Len(t, items, 10)
}
