// Package service provides implementation for service
//
// File: news_service_test.go
// Description: Unit tests for news service
package service_test

import (
	"context"
	"errors"
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/service"

	"git.gerege.mn/backend-packages/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockNewsRepository implements repository.NewsRepository for testing
type mockNewsRepository struct {
	mock.Mock
}

func (m *mockNewsRepository) List(ctx context.Context, q dto.NewsListQuery) ([]domain.News, int64, int, int, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.News), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockNewsRepository) GetByID(ctx context.Context, id int) (domain.News, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.News), args.Error(1)
}

func (m *mockNewsRepository) Create(ctx context.Context, news domain.News) error {
	args := m.Called(ctx, news)
	return args.Error(0)
}

func (m *mockNewsRepository) Update(ctx context.Context, id int, news domain.News) error {
	args := m.Called(ctx, id, news)
	return args.Error(0)
}

func (m *mockNewsRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestNewsService_List(t *testing.T) {
	tests := []struct {
		name      string
		query     dto.NewsListQuery
		mockSetup func(*mockNewsRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name: "success - returns news",
			query: dto.NewsListQuery{
				PaginationQuery: common.PaginationQuery{Page: 1, Size: 10},
			},
			mockSetup: func(m *mockNewsRepository) {
				news := []domain.News{
					{Id: 1, Title: "News 1"},
					{Id: 2, Title: "News 2"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.NewsListQuery")).
					Return(news, int64(2), 1, 10, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "success - empty list",
			query: dto.NewsListQuery{
				PaginationQuery: common.PaginationQuery{Page: 1, Size: 10},
			},
			mockSetup: func(m *mockNewsRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.NewsListQuery")).
					Return([]domain.News{}, int64(0), 1, 10, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "error - db error",
			query: dto.NewsListQuery{
				PaginationQuery: common.PaginationQuery{Page: 1, Size: 10},
			},
			mockSetup: func(m *mockNewsRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.NewsListQuery")).
					Return(nil, int64(0), 0, 0, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockNewsRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewNewsService(mockRepo)
			news, _, _, _, err := svc.List(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, news, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewsService_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		newsID    int
		mockSetup func(*mockNewsRepository)
		wantErr   bool
	}{
		{
			name:   "success - news found",
			newsID: 1,
			mockSetup: func(m *mockNewsRepository) {
				m.On("GetByID", mock.Anything, 1).Return(domain.News{
					Id:    1,
					Title: "Test News",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:   "error - news not found",
			newsID: 999,
			mockSetup: func(m *mockNewsRepository) {
				m.On("GetByID", mock.Anything, 999).Return(domain.News{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockNewsRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewNewsService(mockRepo)
			news, err := svc.GetByID(context.Background(), tt.newsID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newsID, news.Id)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewsService_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     dto.NewsDto
		mockSetup func(*mockNewsRepository)
		wantErr   bool
	}{
		{
			name: "success - news created",
			input: dto.NewsDto{
				Title:    "New News",
				Text:     "News content",
				ImageUrl: "https://example.com/image.jpg",
			},
			mockSetup: func(m *mockNewsRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.News")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - create fails",
			input: dto.NewsDto{
				Title: "Fail News",
			},
			mockSetup: func(m *mockNewsRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.News")).
					Return(errors.New("create failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockNewsRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewNewsService(mockRepo)
			err := svc.Create(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewsService_Update(t *testing.T) {
	tests := []struct {
		name      string
		newsID    int
		input     dto.NewsDto
		mockSetup func(*mockNewsRepository)
		wantErr   bool
	}{
		{
			name:   "success - news updated",
			newsID: 1,
			input: dto.NewsDto{
				Title: "Updated News",
				Text:  "Updated content",
			},
			mockSetup: func(m *mockNewsRepository) {
				m.On("Update", mock.Anything, 1, mock.AnythingOfType("domain.News")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error - update fails",
			newsID: 999,
			input: dto.NewsDto{
				Title: "Fail Update",
			},
			mockSetup: func(m *mockNewsRepository) {
				m.On("Update", mock.Anything, 999, mock.AnythingOfType("domain.News")).
					Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockNewsRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewNewsService(mockRepo)
			err := svc.Update(context.Background(), tt.newsID, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewsService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		newsID    int
		mockSetup func(*mockNewsRepository)
		wantErr   bool
	}{
		{
			name:   "success - news deleted",
			newsID: 1,
			mockSetup: func(m *mockNewsRepository) {
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error - delete fails",
			newsID: 999,
			mockSetup: func(m *mockNewsRepository) {
				m.On("Delete", mock.Anything, 999).Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockNewsRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewNewsService(mockRepo)
			err := svc.Delete(context.Background(), tt.newsID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
