// Package service provides implementation for service
//
// File: role_service_test.go
// Description: implementation for service
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockRoleRepository for testing - implements repository.RoleRepository
type mockRoleRepository struct {
	mock.Mock
}

func (m *mockRoleRepository) List(ctx context.Context, p dto.RoleListQuery) ([]domain.Role, int64, int, int, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.Role), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockRoleRepository) Create(ctx context.Context, r domain.Role) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *mockRoleRepository) Update(ctx context.Context, id int, r domain.Role) error {
	args := m.Called(ctx, id, r)
	return args.Error(0)
}

func (m *mockRoleRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRoleRepository) ByID(ctx context.Context, id int) (domain.Role, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Role), args.Error(1)
}

func (m *mockRoleRepository) Permissions(ctx context.Context, q dto.RolePermissionsQuery) ([]domain.Permission, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Permission), args.Error(1)
}

func (m *mockRoleRepository) ReplacePermissions(ctx context.Context, roleID int, permIDs []int) error {
	args := m.Called(ctx, roleID, permIDs)
	return args.Error(0)
}

func (m *mockRoleRepository) GetUserCount(ctx context.Context, roleID int) int64 {
	args := m.Called(ctx, roleID)
	return int64(args.Int(0))
}

func TestRoleService_ListFilteredPaged(t *testing.T) {
	tests := []struct {
		name      string
		query     dto.RoleListQuery
		mockSetup func(*mockRoleRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "success - returns roles",
			query: dto.RoleListQuery{PaginationQuery: common.PaginationQuery{Page: 1, Size: 10}},
			mockSetup: func(m *mockRoleRepository) {
				roles := []domain.Role{
					{ID: 1, Name: "Admin", Code: "ADMIN"},
					{ID: 2, Name: "User", Code: "USER"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.RoleListQuery")).
					Return(roles, int64(2), 1, 10, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:  "error - db error",
			query: dto.RoleListQuery{PaginationQuery: common.PaginationQuery{Page: 1, Size: 10}},
			mockSetup: func(m *mockRoleRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.RoleListQuery")).
					Return(nil, int64(0), 0, 0, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRoleRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewRoleService(mockRepo, zap.NewNop())

			roles, _, _, _, err := svc.ListFilteredPaged(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, roles, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRoleService_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     dto.RoleCreateDto
		mockSetup func(*mockRoleRepository)
		wantErr   bool
	}{
		{
			name: "success - role created",
			input: dto.RoleCreateDto{
				Name:        "New Role",
				Code:        "NEW_ROLE",
				Description: "Test role",
			},
			mockSetup: func(m *mockRoleRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.Role")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - create fails",
			input: dto.RoleCreateDto{
				Name: "Fail Role",
				Code: "FAIL",
			},
			mockSetup: func(m *mockRoleRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.Role")).
					Return(errors.New("duplicate code"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRoleRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewRoleService(mockRepo, zap.NewNop())

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

func TestRoleService_Delete(t *testing.T) {
	isActiveFalse := false
	isActiveTrue := true

	tests := []struct {
		name      string
		roleID    int
		mockSetup func(*mockRoleRepository)
		wantErr   bool
	}{
		{
			name:   "success - role deleted",
			roleID: 1,
			mockSetup: func(m *mockRoleRepository) {
				m.On("ByID", mock.Anything, 1).Return(domain.Role{ID: 1, IsActive: &isActiveFalse}, nil)
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error - role is active",
			roleID: 2,
			mockSetup: func(m *mockRoleRepository) {
				m.On("ByID", mock.Anything, 2).Return(domain.Role{ID: 2, IsActive: &isActiveTrue}, nil)
			},
			wantErr: true,
		},
		{
			name:   "error - role not found",
			roleID: 999,
			mockSetup: func(m *mockRoleRepository) {
				m.On("ByID", mock.Anything, 999).Return(domain.Role{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRoleRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewRoleService(mockRepo, zap.NewNop())

			err := svc.Delete(context.Background(), tt.roleID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRoleService_GetPermissions(t *testing.T) {
	tests := []struct {
		name      string
		query     dto.RolePermissionsQuery
		mockSetup func(*mockRoleRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "success - returns permissions",
			query: dto.RolePermissionsQuery{RoleID: 1},
			mockSetup: func(m *mockRoleRepository) {
				perms := []domain.Permission{
					{ID: 1, Code: "READ", Name: "Read"},
					{ID: 2, Code: "WRITE", Name: "Write"},
				}
				m.On("Permissions", mock.Anything, mock.AnythingOfType("dto.RolePermissionsQuery")).Return(perms, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:  "success - no permissions",
			query: dto.RolePermissionsQuery{RoleID: 2},
			mockSetup: func(m *mockRoleRepository) {
				m.On("Permissions", mock.Anything, mock.AnythingOfType("dto.RolePermissionsQuery")).Return([]domain.Permission{}, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRoleRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewRoleService(mockRepo, zap.NewNop())

			perms, err := svc.GetPermissions(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, perms, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRoleService_SetPermissions(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.RolePermissionsUpdateDto
		mockSetup func(*mockRoleRepository)
		wantErr   bool
	}{
		{
			name: "success - permissions set",
			req:  dto.RolePermissionsUpdateDto{RoleID: 1, PermissionIDs: []int{1, 2, 3}},
			mockSetup: func(m *mockRoleRepository) {
				m.On("ReplacePermissions", mock.Anything, 1, []int{1, 2, 3}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - set fails",
			req:  dto.RolePermissionsUpdateDto{RoleID: 1, PermissionIDs: []int{999}},
			mockSetup: func(m *mockRoleRepository) {
				m.On("ReplacePermissions", mock.Anything, 1, []int{999}).Return(errors.New("invalid permission"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRoleRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewRoleService(mockRepo, zap.NewNop())

			err := svc.SetPermissions(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRoleService_Update(t *testing.T) {
	isActiveFalse := false

	tests := []struct {
		name      string
		roleID    int
		req       dto.RoleUpdateDto
		mockSetup func(*mockRoleRepository)
		wantErr   bool
	}{
		{
			name:   "success - role updated",
			roleID: 1,
			req: dto.RoleUpdateDto{
				Name: "Updated Role",
				Code: "UPDATED",
			},
			mockSetup: func(m *mockRoleRepository) {
				m.On("Update", mock.Anything, 1, mock.AnythingOfType("domain.Role")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error - cannot deactivate role with users",
			roleID: 2,
			req: dto.RoleUpdateDto{
				Name:     "Role with users",
				IsActive: &isActiveFalse,
			},
			mockSetup: func(m *mockRoleRepository) {
				m.On("GetUserCount", mock.Anything, 2).Return(5)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRoleRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewRoleService(mockRepo, zap.NewNop())

			err := svc.Update(context.Background(), tt.roleID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
