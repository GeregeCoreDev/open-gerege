// Package service provides implementation for service
//
// File: notification_service_test.go
// Description: Unit tests for notification service
package service_test

import (
	"context"
	"errors"
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/service"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockNotificationRepository implements repository.NotificationRepository
type mockNotificationRepository struct {
	mock.Mock
}

func (m *mockNotificationRepository) ListByUser(ctx context.Context, userID int, p common.PaginationQuery) ([]domain.Notification, int64, int, int, error) {
	args := m.Called(ctx, userID, p)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.Notification), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockNotificationRepository) ListGroups(ctx context.Context, p common.PaginationQuery) ([]domain.NotificationGroup, int64, int, int, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.NotificationGroup), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockNotificationRepository) MarkGroupRead(ctx context.Context, userID, groupID int) error {
	args := m.Called(ctx, userID, groupID)
	return args.Error(0)
}

func (m *mockNotificationRepository) MarkAllRead(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockNotificationRepository) CreateGroup(ctx context.Context, g domain.NotificationGroup) (domain.NotificationGroup, error) {
	args := m.Called(ctx, g)
	return args.Get(0).(domain.NotificationGroup), args.Error(1)
}

func (m *mockNotificationRepository) CreateNotification(ctx context.Context, n domain.Notification) (domain.Notification, error) {
	args := m.Called(ctx, n)
	return args.Get(0).(domain.Notification), args.Error(1)
}

func (m *mockNotificationRepository) CreateNotificationsBulk(ctx context.Context, ns []domain.Notification) error {
	args := m.Called(ctx, ns)
	return args.Error(0)
}

func (m *mockNotificationRepository) AllUserIDs(ctx context.Context) ([]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func TestNotificationService_List(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		query     common.PaginationQuery
		mockSetup func(*mockNotificationRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:   "success - returns notifications",
			userID: 1,
			query:  common.PaginationQuery{Page: 1, Size: 10},
			mockSetup: func(m *mockNotificationRepository) {
				notifications := []domain.Notification{
					{Id: 1, UserId: 1, Title: "Notification 1"},
					{Id: 2, UserId: 1, Title: "Notification 2"},
				}
				m.On("ListByUser", mock.Anything, 1, mock.AnythingOfType("common.PaginationQuery")).
					Return(notifications, int64(2), 1, 10, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:   "success - empty list",
			userID: 2,
			query:  common.PaginationQuery{Page: 1, Size: 10},
			mockSetup: func(m *mockNotificationRepository) {
				m.On("ListByUser", mock.Anything, 2, mock.AnythingOfType("common.PaginationQuery")).
					Return([]domain.Notification{}, int64(0), 1, 10, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:   "error - db error",
			userID: 3,
			query:  common.PaginationQuery{Page: 1, Size: 10},
			mockSetup: func(m *mockNotificationRepository) {
				m.On("ListByUser", mock.Anything, 3, mock.AnythingOfType("common.PaginationQuery")).
					Return(nil, int64(0), 0, 0, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockNotificationRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewNotificationService(mockRepo, &config.Config{})

			notifications, _, _, _, err := svc.List(context.Background(), tt.userID, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, notifications, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNotificationService_Groups(t *testing.T) {
	tests := []struct {
		name      string
		query     common.PaginationQuery
		mockSetup func(*mockNotificationRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "success - returns groups",
			query: common.PaginationQuery{Page: 1, Size: 10},
			mockSetup: func(m *mockNotificationRepository) {
				groups := []domain.NotificationGroup{
					{Id: 1, Title: "Group 1"},
					{Id: 2, Title: "Group 2"},
				}
				m.On("ListGroups", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return(groups, int64(2), 1, 10, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:  "error - db error",
			query: common.PaginationQuery{Page: 1, Size: 10},
			mockSetup: func(m *mockNotificationRepository) {
				m.On("ListGroups", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return(nil, int64(0), 0, 0, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockNotificationRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewNotificationService(mockRepo, &config.Config{})

			groups, _, _, _, err := svc.Groups(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, groups, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNotificationService_MarkGroupRead(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		groupID   int
		mockSetup func(*mockNotificationRepository)
		wantErr   bool
	}{
		{
			name:    "success - marked as read",
			userID:  1,
			groupID: 10,
			mockSetup: func(m *mockNotificationRepository) {
				m.On("MarkGroupRead", mock.Anything, 1, 10).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "error - db error",
			userID:  2,
			groupID: 20,
			mockSetup: func(m *mockNotificationRepository) {
				m.On("MarkGroupRead", mock.Anything, 2, 20).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockNotificationRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewNotificationService(mockRepo, &config.Config{})

			err := svc.MarkGroupRead(context.Background(), tt.userID, tt.groupID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNotificationService_MarkAllRead(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		mockSetup func(*mockNotificationRepository)
		wantErr   bool
	}{
		{
			name:   "success - all marked as read",
			userID: 1,
			mockSetup: func(m *mockNotificationRepository) {
				m.On("MarkAllRead", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error - db error",
			userID: 2,
			mockSetup: func(m *mockNotificationRepository) {
				m.On("MarkAllRead", mock.Anything, 2).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockNotificationRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewNotificationService(mockRepo, &config.Config{})

			err := svc.MarkAllRead(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
