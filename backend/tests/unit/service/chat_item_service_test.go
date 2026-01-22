// Package service provides implementation for service
//
// File: chat_item_service_test.go
// Description: Unit tests for chat item service
package service_test

import (
	"context"
	"errors"
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// mockChatItemRepository implements repository.ChatItemRepository
type mockChatItemRepository struct {
	mock.Mock
}

func (m *mockChatItemRepository) List(ctx context.Context, q dto.ChatItemQuery) ([]domain.ChatItem, int64, int, int, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.ChatItem), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockChatItemRepository) ByID(ctx context.Context, id int) (domain.ChatItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.ChatItem), args.Error(1)
}

func (m *mockChatItemRepository) FindByKey(ctx context.Context, key string) (domain.ChatItem, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(domain.ChatItem), args.Error(1)
}

func (m *mockChatItemRepository) Create(ctx context.Context, item domain.ChatItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *mockChatItemRepository) Update(ctx context.Context, id int, item domain.ChatItem) error {
	args := m.Called(ctx, id, item)
	return args.Error(0)
}

func (m *mockChatItemRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestChatItemService_GetByKey(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		mockSetup func(*mockChatItemRepository)
		wantErr   bool
	}{
		{
			name: "success - found by key",
			key:  "greeting",
			mockSetup: func(m *mockChatItemRepository) {
				m.On("FindByKey", mock.Anything, "greeting").
					Return(domain.ChatItem{ID: 1, Key: "greeting", Answer: "Hello!"}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - not found",
			key:  "unknown",
			mockSetup: func(m *mockChatItemRepository) {
				m.On("FindByKey", mock.Anything, "unknown").
					Return(domain.ChatItem{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockChatItemRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewChatItemService(mockRepo, zap.NewNop())

			item, err := svc.GetByKey(context.Background(), tt.key)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.key, item.Key)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestChatItemService_List(t *testing.T) {
	tests := []struct {
		name      string
		query     dto.ChatItemQuery
		mockSetup func(*mockChatItemRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "success - returns items",
			query: dto.ChatItemQuery{},
			mockSetup: func(m *mockChatItemRepository) {
				items := []domain.ChatItem{
					{ID: 1, Key: "hello", Answer: "Hi there!"},
					{ID: 2, Key: "bye", Answer: "Goodbye!"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.ChatItemQuery")).
					Return(items, int64(2), 1, 10, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:  "success - empty list",
			query: dto.ChatItemQuery{},
			mockSetup: func(m *mockChatItemRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.ChatItemQuery")).
					Return([]domain.ChatItem{}, int64(0), 1, 10, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:  "error - db error",
			query: dto.ChatItemQuery{},
			mockSetup: func(m *mockChatItemRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.ChatItemQuery")).
					Return(nil, int64(0), 0, 0, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockChatItemRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewChatItemService(mockRepo, zap.NewNop())

			items, _, _, _, err := svc.List(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, items, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestChatItemService_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     dto.ChatItemCreateDto
		mockSetup func(*mockChatItemRepository)
		wantErr   bool
	}{
		{
			name: "success - created",
			input: dto.ChatItemCreateDto{
				Key:    "new_key",
				Answer: "New answer",
			},
			mockSetup: func(m *mockChatItemRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(item domain.ChatItem) bool {
					return item.Key == "new_key" && item.Answer == "New answer"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - create fails",
			input: dto.ChatItemCreateDto{
				Key:    "fail_key",
				Answer: "Fail answer",
			},
			mockSetup: func(m *mockChatItemRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.ChatItem")).
					Return(errors.New("create failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockChatItemRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewChatItemService(mockRepo, zap.NewNop())

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

func TestChatItemService_Update(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		input     dto.ChatItemUpdateDto
		mockSetup func(*mockChatItemRepository)
		wantErr   bool
	}{
		{
			name: "success - updated",
			id:   1,
			input: dto.ChatItemUpdateDto{
				Key:    "updated_key",
				Answer: "Updated answer",
			},
			mockSetup: func(m *mockChatItemRepository) {
				m.On("Update", mock.Anything, 1, mock.MatchedBy(func(item domain.ChatItem) bool {
					return item.Key == "updated_key" && item.Answer == "Updated answer"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - update fails",
			id:   999,
			input: dto.ChatItemUpdateDto{
				Key:    "fail",
				Answer: "fail",
			},
			mockSetup: func(m *mockChatItemRepository) {
				m.On("Update", mock.Anything, 999, mock.AnythingOfType("domain.ChatItem")).
					Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockChatItemRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewChatItemService(mockRepo, zap.NewNop())

			err := svc.Update(context.Background(), tt.id, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestChatItemService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		mockSetup func(*mockChatItemRepository)
		wantErr   bool
	}{
		{
			name: "success - deleted",
			id:   1,
			mockSetup: func(m *mockChatItemRepository) {
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - delete fails",
			id:   999,
			mockSetup: func(m *mockChatItemRepository) {
				m.On("Delete", mock.Anything, 999).Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockChatItemRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewChatItemService(mockRepo, zap.NewNop())

			err := svc.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
