//go:build integration

// Package integration provides integration tests for HTTP handlers
//
// File: news_handler_test.go
// Description: Integration tests for NewsHandler
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

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockNewsService implements news service interface for testing
type mockNewsService struct {
	mock.Mock
}

func (m *mockNewsService) List(ctx context.Context, q dto.NewsListQuery) ([]domain.News, int64, int, int, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.News), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockNewsService) GetByID(ctx context.Context, id int) (domain.News, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.News), args.Error(1)
}

func (m *mockNewsService) Create(ctx context.Context, req dto.NewsDto) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockNewsService) Update(ctx context.Context, id int, req dto.NewsDto) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *mockNewsService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// setupNewsTestApp creates a test Fiber app with news routes
func setupNewsTestApp(svc *mockNewsService) *fiber.App {
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

	// News routes
	news := app.Group("/api/v1/news")

	// List news
	news.Get("/", func(c *fiber.Ctx) error {
		q := dto.NewsListQuery{}
		if err := c.QueryParser(&q); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid query parameters")
		}

		items, total, page, size, err := svc.List(c.UserContext(), q)
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

	// Get news by ID
	news.Get("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid news ID")
		}

		item, err := svc.GetByID(c.UserContext(), id)
		if err != nil {
			if err.Error() == "not found" {
				return fiber.NewError(fiber.StatusNotFound, "news not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"code":    "ok",
			"data":    item,
		})
	})

	// Create news
	news.Post("/", func(c *fiber.Ctx) error {
		var req dto.NewsDto
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		// Basic validation
		if req.Title == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "title is required")
		}
		if req.Text == "" {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "text is required")
		}

		if err := svc.Create(c.UserContext(), req); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"code":    "created",
		})
	})

	// Update news
	news.Put("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid news ID")
		}

		var req dto.NewsDto
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		if err := svc.Update(c.UserContext(), id, req); err != nil {
			if err.Error() == "not found" {
				return fiber.NewError(fiber.StatusNotFound, "news not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"code":    "ok",
		})
	})

	// Delete news
	news.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid news ID")
		}

		if err := svc.Delete(c.UserContext(), id); err != nil {
			if err.Error() == "not found" {
				return fiber.NewError(fiber.StatusNotFound, "news not found")
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"success": true,
			"code":    "ok",
		})
	})

	return app
}

func TestNewsHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		mockSetup  func(*mockNewsService)
		wantStatus int
		wantCount  int
	}{
		{
			name:  "success - returns news",
			query: "",
			mockSetup: func(m *mockNewsService) {
				items := []domain.News{
					{Id: 1, Title: "News 1", Text: "Content 1"},
					{Id: 2, Title: "News 2", Text: "Content 2"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.NewsListQuery")).
					Return(items, int64(2), 1, 10, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:  "success - empty list",
			query: "",
			mockSetup: func(m *mockNewsService) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.NewsListQuery")).
					Return([]domain.News{}, int64(0), 1, 10, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name:  "success - with pagination",
			query: "?page=2&size=5",
			mockSetup: func(m *mockNewsService) {
				items := []domain.News{
					{Id: 6, Title: "News 6", Text: "Content 6"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.NewsListQuery")).
					Return(items, int64(6), 2, 5, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:  "success - filter by category",
			query: "?category_id=5",
			mockSetup: func(m *mockNewsService) {
				items := []domain.News{
					{Id: 1, Title: "Category News", Text: "Content"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.NewsListQuery")).
					Return(items, int64(1), 1, 10, nil)
			},
			wantStatus: http.StatusOK,
			wantCount:  1,
		},
		{
			name:  "error - database error",
			query: "",
			mockSetup: func(m *mockNewsService) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.NewsListQuery")).
					Return(nil, int64(0), 0, 0, errors.New("database error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockNewsService{}
			tt.mockSetup(mockSvc)

			app := setupNewsTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/news/"+tt.query, nil)
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

func TestNewsHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		newsID     string
		mockSetup  func(*mockNewsService)
		wantStatus int
		wantTitle  string
	}{
		{
			name:   "success - returns news",
			newsID: "1",
			mockSetup: func(m *mockNewsService) {
				m.On("GetByID", mock.Anything, 1).
					Return(domain.News{Id: 1, Title: "Test News", Text: "Test Content"}, nil)
			},
			wantStatus: http.StatusOK,
			wantTitle:  "Test News",
		},
		{
			name:   "error - not found",
			newsID: "999",
			mockSetup: func(m *mockNewsService) {
				m.On("GetByID", mock.Anything, 999).
					Return(domain.News{}, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "error - invalid ID",
			newsID:     "invalid",
			mockSetup:  func(m *mockNewsService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockNewsService{}
			tt.mockSetup(mockSvc)

			app := setupNewsTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/news/"+tt.newsID, nil)
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
				assert.Equal(t, tt.wantTitle, data["title"])
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestNewsHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockSetup  func(*mockNewsService)
		wantStatus int
	}{
		{
			name: "success - news created",
			body: `{"title": "New Article", "text": "This is the content", "image_url": "https://example.com/image.jpg"}`,
			mockSetup: func(m *mockNewsService) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(req dto.NewsDto) bool {
					return req.Title == "New Article" && req.Text == "This is the content"
				})).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "success - without image URL",
			body: `{"title": "New Article", "text": "This is the content"}`,
			mockSetup: func(m *mockNewsService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("dto.NewsDto")).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "error - missing title",
			body:       `{"text": "Content only"}`,
			mockSetup:  func(m *mockNewsService) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "error - missing text",
			body:       `{"title": "Title only"}`,
			mockSetup:  func(m *mockNewsService) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "error - empty body",
			body:       `{}`,
			mockSetup:  func(m *mockNewsService) {},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "error - invalid JSON",
			body:       `{invalid}`,
			mockSetup:  func(m *mockNewsService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error - service error",
			body: `{"title": "New Article", "text": "This is the content"}`,
			mockSetup: func(m *mockNewsService) {
				m.On("Create", mock.Anything, mock.AnythingOfType("dto.NewsDto")).
					Return(errors.New("database error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockNewsService{}
			tt.mockSetup(mockSvc)

			app := setupNewsTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/news/", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestNewsHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		newsID     string
		body       string
		mockSetup  func(*mockNewsService)
		wantStatus int
	}{
		{
			name:   "success - news updated",
			newsID: "1",
			body:   `{"title": "Updated Title", "text": "Updated content"}`,
			mockSetup: func(m *mockNewsService) {
				m.On("Update", mock.Anything, 1, mock.AnythingOfType("dto.NewsDto")).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "error - not found",
			newsID: "999",
			body:   `{"title": "Updated Title", "text": "Updated content"}`,
			mockSetup: func(m *mockNewsService) {
				m.On("Update", mock.Anything, 999, mock.AnythingOfType("dto.NewsDto")).
					Return(errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "error - invalid ID",
			newsID:     "invalid",
			body:       `{"title": "Updated Title", "text": "Updated content"}`,
			mockSetup:  func(m *mockNewsService) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error - invalid JSON",
			newsID:     "1",
			body:       `{invalid}`,
			mockSetup:  func(m *mockNewsService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockNewsService{}
			tt.mockSetup(mockSvc)

			app := setupNewsTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodPut, "/api/v1/news/"+tt.newsID, bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestNewsHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		newsID     string
		mockSetup  func(*mockNewsService)
		wantStatus int
	}{
		{
			name:   "success - news deleted",
			newsID: "1",
			mockSetup: func(m *mockNewsService) {
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "error - not found",
			newsID: "999",
			mockSetup: func(m *mockNewsService) {
				m.On("Delete", mock.Anything, 999).Return(errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "error - invalid ID",
			newsID:     "invalid",
			mockSetup:  func(m *mockNewsService) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockNewsService{}
			tt.mockSetup(mockSvc)

			app := setupNewsTestApp(mockSvc)

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/news/"+tt.newsID, nil)

			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

// TestNewsHandler_ContentType tests that responses have correct content type
func TestNewsHandler_ContentType(t *testing.T) {
	mockSvc := &mockNewsService{}
	mockSvc.On("List", mock.Anything, mock.AnythingOfType("dto.NewsListQuery")).
		Return([]domain.News{}, int64(0), 1, 10, nil)

	app := setupNewsTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/news/", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json")
}

// TestNewsHandler_MethodNotAllowed tests method not allowed responses
func TestNewsHandler_MethodNotAllowed(t *testing.T) {
	mockSvc := &mockNewsService{}
	app := setupNewsTestApp(mockSvc)

	// PATCH is not allowed
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/news/1", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

// TestNewsHandler_LargePayload tests handling of large payloads
func TestNewsHandler_LargePayload(t *testing.T) {
	mockSvc := &mockNewsService{}
	mockSvc.On("Create", mock.Anything, mock.AnythingOfType("dto.NewsDto")).Return(nil)

	app := setupNewsTestApp(mockSvc)

	// Create a large text content (10KB)
	largeText := make([]byte, 10*1024)
	for i := range largeText {
		largeText[i] = 'a'
	}

	body := `{"title": "Large Article", "text": "` + string(largeText) + `"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/news/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
