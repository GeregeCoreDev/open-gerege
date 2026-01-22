// Package service provides implementation for service
//
// File: user_service_test.go
// Description: Unit tests for user service
// Author: Bayarsaikhan Otgonbayar, CTO
// Company: Gerege Core Team
// Created: 2025-02-20
// Last Updated: 2025-02-20
package service_test

import (
	"context"
	"errors"
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/service"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockUserRepository for testing - implements repository.UserRepository
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) List(ctx context.Context, p common.PaginationQuery) ([]domain.User, int64, int, int, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockUserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserRepository) Update(ctx context.Context, u domain.User) (domain.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserRepository) Delete(ctx context.Context, id int) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserRepository) UserOrgIDs(ctx context.Context, userID int) ([]int, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]int), args.Error(1)
}

func (m *mockUserRepository) GetOrganizationsByIDs(ctx context.Context, ids []int, fields []string) ([]domain.Organization, error) {
	args := m.Called(ctx, ids, fields)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func (m *mockUserRepository) GetOrganization(ctx context.Context, id int, fields []string) (*domain.Organization, error) {
	args := m.Called(ctx, id, fields)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Organization), args.Error(1)
}

func TestUserService_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		mockSetup func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name:   "success - user found",
			userID: 1,
			mockSetup: func(m *mockUserRepository) {
				m.On("GetByID", mock.Anything, 1).Return(domain.User{
					Id:        1,
					FirstName: "Test",
					LastName:  "User",
					RegNo:     "AA12345678",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:   "error - user not found",
			userID: 999,
			mockSetup: func(m *mockUserRepository) {
				m.On("GetByID", mock.Anything, 999).Return(domain.User{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewUserService(mockRepo, &config.Config{}, zap.NewNop())

			user, err := svc.GetByID(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, user.Id)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_List(t *testing.T) {
	tests := []struct {
		name      string
		query     common.PaginationQuery
		mockSetup func(*mockUserRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "success - returns users",
			query: common.PaginationQuery{Page: 1, Size: 10},
			mockSetup: func(m *mockUserRepository) {
				users := []domain.User{
					{Id: 1, FirstName: "User1"},
					{Id: 2, FirstName: "User2"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return(users, int64(2), 1, 10, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:  "success - empty list",
			query: common.PaginationQuery{Page: 1, Size: 10},
			mockSetup: func(m *mockUserRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return([]domain.User{}, int64(0), 1, 10, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:  "error - db error",
			query: common.PaginationQuery{Page: 1, Size: 10},
			mockSetup: func(m *mockUserRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return(nil, int64(0), 0, 0, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewUserService(mockRepo, &config.Config{}, zap.NewNop())

			users, _, _, _, err := svc.List(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, users, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     dto.UserCreateDto
		mockSetup func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name: "success - new user created",
			input: dto.UserCreateDto{
				Id:        1,
				FirstName: "New",
				LastName:  "User",
				RegNo:     "AA12345678",
			},
			mockSetup: func(m *mockUserRepository) {
				m.On("GetByID", mock.Anything, 1).Return(domain.User{}, errors.New("not found"))
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.User")).
					Return(domain.User{Id: 1, FirstName: "New", LastName: "User"}, nil)
			},
			wantErr: false,
		},
		{
			name: "success - user already exists",
			input: dto.UserCreateDto{
				Id:        2,
				FirstName: "Existing",
				LastName:  "User",
			},
			mockSetup: func(m *mockUserRepository) {
				m.On("GetByID", mock.Anything, 2).Return(domain.User{Id: 2, FirstName: "Existing"}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - create fails",
			input: dto.UserCreateDto{
				Id:        3,
				FirstName: "Fail",
			},
			mockSetup: func(m *mockUserRepository) {
				m.On("GetByID", mock.Anything, 3).Return(domain.User{}, errors.New("not found"))
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.User")).
					Return(domain.User{}, errors.New("create failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewUserService(mockRepo, &config.Config{}, zap.NewNop())

			_, err := svc.Create(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	tests := []struct {
		name      string
		input     dto.UserUpdateDto
		mockSetup func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name: "success - user updated",
			input: dto.UserUpdateDto{
				Id:        1,
				FirstName: "Updated",
				LastName:  "User",
			},
			mockSetup: func(m *mockUserRepository) {
				m.On("GetByID", mock.Anything, 1).Return(domain.User{Id: 1}, nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("domain.User")).
					Return(domain.User{Id: 1, FirstName: "Updated"}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - user not found",
			input: dto.UserUpdateDto{
				Id:        999,
				FirstName: "NotFound",
			},
			mockSetup: func(m *mockUserRepository) {
				m.On("GetByID", mock.Anything, 999).Return(domain.User{}, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "error - update fails",
			input: dto.UserUpdateDto{
				Id:        2,
				FirstName: "Fail",
			},
			mockSetup: func(m *mockUserRepository) {
				m.On("GetByID", mock.Anything, 2).Return(domain.User{Id: 2}, nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("domain.User")).
					Return(domain.User{}, errors.New("update failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewUserService(mockRepo, &config.Config{}, zap.NewNop())

			_, err := svc.Update(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		mockSetup func(*mockUserRepository)
		wantErr   bool
	}{
		{
			name:   "success - user deleted",
			userID: 1,
			mockSetup: func(m *mockUserRepository) {
				m.On("Delete", mock.Anything, 1).Return(domain.User{Id: 1}, nil)
			},
			wantErr: false,
		},
		{
			name:   "error - delete fails",
			userID: 999,
			mockSetup: func(m *mockUserRepository) {
				m.On("Delete", mock.Anything, 999).Return(domain.User{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewUserService(mockRepo, &config.Config{}, zap.NewNop())

			_, err := svc.Delete(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Organizations(t *testing.T) {
	tests := []struct {
		name         string
		userID       int
		currentOrgID int
		mockSetup    func(*mockUserRepository)
		wantOrgCount int
		wantErr      bool
	}{
		{
			name:         "success - returns organizations",
			userID:       1,
			currentOrgID: 10,
			mockSetup: func(m *mockUserRepository) {
				m.On("UserOrgIDs", mock.Anything, 1).Return([]int{10, 20, 30}, nil)
				m.On("GetOrganizationsByIDs", mock.Anything, []int{10, 20, 30}, mock.Anything).
					Return([]domain.Organization{
						{Id: 10, Name: "Org1"},
						{Id: 20, Name: "Org2"},
						{Id: 30, Name: "Org3"},
					}, nil)
				m.On("GetOrganization", mock.Anything, 10, mock.Anything).
					Return(&domain.Organization{Id: 10, Name: "Org1"}, nil)
			},
			wantOrgCount: 3,
			wantErr:      false,
		},
		{
			name:         "success - no current org",
			userID:       2,
			currentOrgID: 0,
			mockSetup: func(m *mockUserRepository) {
				m.On("UserOrgIDs", mock.Anything, 2).Return([]int{10}, nil)
				m.On("GetOrganizationsByIDs", mock.Anything, []int{10}, mock.Anything).
					Return([]domain.Organization{{Id: 10, Name: "Org1"}}, nil)
			},
			wantOrgCount: 1,
			wantErr:      false,
		},
		{
			name:         "error - get org ids fails",
			userID:       3,
			currentOrgID: 0,
			mockSetup: func(m *mockUserRepository) {
				m.On("UserOrgIDs", mock.Anything, 3).Return(nil, errors.New("db error"))
			},
			wantOrgCount: 0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewUserService(mockRepo, &config.Config{}, zap.NewNop())

			_, _, orgs, err := svc.Organizations(context.Background(), tt.userID, tt.currentOrgID, []string{"id", "name"})

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, orgs, tt.wantOrgCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
