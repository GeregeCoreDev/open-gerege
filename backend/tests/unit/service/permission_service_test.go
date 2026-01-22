// Package service provides implementation for service
//
// File: permission_service_test.go
// Description: Unit tests for permission service
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

// mockPermissionRepository implements repository.PermissionRepository
type mockPermissionRepository struct {
	mock.Mock
}

func (m *mockPermissionRepository) List(ctx context.Context, q dto.PermissionQuery) ([]domain.Permission, int64, int, int, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.Permission), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockPermissionRepository) ByID(ctx context.Context, id int) (domain.Permission, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Permission), args.Error(1)
}

func (m *mockPermissionRepository) ByCode(ctx context.Context, code string) (domain.Permission, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(domain.Permission), args.Error(1)
}

func (m *mockPermissionRepository) Create(ctx context.Context, p domain.Permission) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockPermissionRepository) CreateBatch(ctx context.Context, systemID, moduleID int, actionIDs []int64) error {
	args := m.Called(ctx, systemID, moduleID, actionIDs)
	return args.Error(0)
}

func (m *mockPermissionRepository) Update(ctx context.Context, id int, p domain.Permission) error {
	args := m.Called(ctx, id, p)
	return args.Error(0)
}

func (m *mockPermissionRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockPermissionRepository) UserHasPermission(ctx context.Context, userID int, permissionCode string) (bool, error) {
	args := m.Called(ctx, userID, permissionCode)
	return args.Bool(0), args.Error(1)
}

func (m *mockPermissionRepository) GetUserPermissionCodes(ctx context.Context, userID int) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// mockCacheInvalidator implements auth.CacheInvalidator
type mockCacheInvalidator struct {
	mock.Mock
}

func (m *mockCacheInvalidator) InvalidateAll() {
	m.Called()
}

func (m *mockCacheInvalidator) InvalidateUser(userID int) {
	m.Called(userID)
}

func (m *mockCacheInvalidator) InvalidateUsers(userIDs []int) {
	m.Called(userIDs)
}

func TestPermissionService_ListFilteredPaged(t *testing.T) {
	tests := []struct {
		name      string
		query     dto.PermissionQuery
		mockSetup func(*mockPermissionRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "success - returns permissions",
			query: dto.PermissionQuery{},
			mockSetup: func(m *mockPermissionRepository) {
				permissions := []domain.Permission{
					{ID: 1, Code: "admin.user.read", Name: "Read users"},
					{ID: 2, Code: "admin.user.write", Name: "Write users"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.PermissionQuery")).
					Return(permissions, int64(2), 1, 10, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:  "success - empty list",
			query: dto.PermissionQuery{},
			mockSetup: func(m *mockPermissionRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.PermissionQuery")).
					Return([]domain.Permission{}, int64(0), 1, 10, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:  "error - db error",
			query: dto.PermissionQuery{},
			mockSetup: func(m *mockPermissionRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.PermissionQuery")).
					Return(nil, int64(0), 0, 0, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockPermissionRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewPermissionService(mockRepo, zap.NewNop())

			permissions, _, _, _, err := svc.ListFilteredPaged(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, permissions, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPermissionService_ByID(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		mockSetup func(*mockPermissionRepository)
		wantErr   bool
	}{
		{
			name: "success - found",
			id:   1,
			mockSetup: func(m *mockPermissionRepository) {
				m.On("ByID", mock.Anything, 1).
					Return(domain.Permission{ID: 1, Code: "admin.user.read"}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - not found",
			id:   999,
			mockSetup: func(m *mockPermissionRepository) {
				m.On("ByID", mock.Anything, 999).
					Return(domain.Permission{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockPermissionRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewPermissionService(mockRepo, zap.NewNop())

			permission, err := svc.ByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, permission.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPermissionService_ByCode(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		mockSetup func(*mockPermissionRepository)
		wantErr   bool
	}{
		{
			name: "success - found",
			code: "admin.user.read",
			mockSetup: func(m *mockPermissionRepository) {
				m.On("ByCode", mock.Anything, "admin.user.read").
					Return(domain.Permission{ID: 1, Code: "admin.user.read"}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - not found",
			code: "unknown.permission",
			mockSetup: func(m *mockPermissionRepository) {
				m.On("ByCode", mock.Anything, "unknown.permission").
					Return(domain.Permission{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockPermissionRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewPermissionService(mockRepo, zap.NewNop())

			permission, err := svc.ByCode(context.Background(), tt.code)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.code, permission.Code)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPermissionService_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     dto.PermissionCreateDto
		mockSetup func(*mockPermissionRepository)
		wantErr   bool
	}{
		{
			name: "success - created batch",
			input: dto.PermissionCreateDto{
				SystemID:  1,
				ModuleID:  2,
				ActionIDs: []int64{1, 2, 3},
			},
			mockSetup: func(m *mockPermissionRepository) {
				m.On("CreateBatch", mock.Anything, 1, 2, []int64{1, 2, 3}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - create fails",
			input: dto.PermissionCreateDto{
				SystemID:  1,
				ModuleID:  2,
				ActionIDs: []int64{1},
			},
			mockSetup: func(m *mockPermissionRepository) {
				m.On("CreateBatch", mock.Anything, 1, 2, []int64{1}).Return(errors.New("create failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockPermissionRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewPermissionService(mockRepo, zap.NewNop())

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

func TestPermissionService_Update(t *testing.T) {
	tests := []struct {
		name           string
		id             int
		input          dto.PermissionUpdateDto
		useCache       bool
		mockSetup      func(*mockPermissionRepository, *mockCacheInvalidator)
		wantErr        bool
		wantCacheClean bool
	}{
		{
			name: "success - updated with cache invalidation",
			id:   1,
			input: dto.PermissionUpdateDto{
				Code:        "admin.user.updated",
				Name:        "Updated permission",
				Description: "Updated description",
			},
			useCache: true,
			mockSetup: func(m *mockPermissionRepository, c *mockCacheInvalidator) {
				m.On("Update", mock.Anything, 1, mock.AnythingOfType("domain.Permission")).Return(nil)
				c.On("InvalidateAll").Return()
			},
			wantErr:        false,
			wantCacheClean: true,
		},
		{
			name: "success - updated without cache",
			id:   2,
			input: dto.PermissionUpdateDto{
				Code: "admin.user.updated2",
				Name: "Updated permission 2",
			},
			useCache: false,
			mockSetup: func(m *mockPermissionRepository, c *mockCacheInvalidator) {
				m.On("Update", mock.Anything, 2, mock.AnythingOfType("domain.Permission")).Return(nil)
			},
			wantErr:        false,
			wantCacheClean: false,
		},
		{
			name: "error - update fails",
			id:   999,
			input: dto.PermissionUpdateDto{
				Code: "fail",
			},
			useCache: false,
			mockSetup: func(m *mockPermissionRepository, c *mockCacheInvalidator) {
				m.On("Update", mock.Anything, 999, mock.AnythingOfType("domain.Permission")).
					Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockPermissionRepository{}
			mockCache := &mockCacheInvalidator{}
			tt.mockSetup(mockRepo, mockCache)

			svc := service.NewPermissionService(mockRepo, zap.NewNop())
			if tt.useCache {
				svc.SetCacheInvalidator(mockCache)
			}

			err := svc.Update(context.Background(), tt.id, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			if tt.wantCacheClean {
				mockCache.AssertExpectations(t)
			}
		})
	}
}

func TestPermissionService_Delete(t *testing.T) {
	tests := []struct {
		name           string
		id             int
		useCache       bool
		mockSetup      func(*mockPermissionRepository, *mockCacheInvalidator)
		wantErr        bool
		wantCacheClean bool
	}{
		{
			name:     "success - deleted with cache invalidation",
			id:       1,
			useCache: true,
			mockSetup: func(m *mockPermissionRepository, c *mockCacheInvalidator) {
				m.On("Delete", mock.Anything, 1).Return(nil)
				c.On("InvalidateAll").Return()
			},
			wantErr:        false,
			wantCacheClean: true,
		},
		{
			name:     "success - deleted without cache",
			id:       2,
			useCache: false,
			mockSetup: func(m *mockPermissionRepository, c *mockCacheInvalidator) {
				m.On("Delete", mock.Anything, 2).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "error - delete fails",
			id:       999,
			useCache: false,
			mockSetup: func(m *mockPermissionRepository, c *mockCacheInvalidator) {
				m.On("Delete", mock.Anything, 999).Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockPermissionRepository{}
			mockCache := &mockCacheInvalidator{}
			tt.mockSetup(mockRepo, mockCache)

			svc := service.NewPermissionService(mockRepo, zap.NewNop())
			if tt.useCache {
				svc.SetCacheInvalidator(mockCache)
			}

			err := svc.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			if tt.wantCacheClean {
				mockCache.AssertExpectations(t)
			}
		})
	}
}

func TestPermissionService_HasPermission(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		permissionCode string
		mockSetup      func(*mockPermissionRepository)
		wantResult     bool
		wantErr        bool
	}{
		{
			name:           "success - has permission",
			userID:         1,
			permissionCode: "admin.user.read",
			mockSetup: func(m *mockPermissionRepository) {
				m.On("UserHasPermission", mock.Anything, 1, "admin.user.read").Return(true, nil)
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name:           "success - no permission",
			userID:         2,
			permissionCode: "admin.user.delete",
			mockSetup: func(m *mockPermissionRepository) {
				m.On("UserHasPermission", mock.Anything, 2, "admin.user.delete").Return(false, nil)
			},
			wantResult: false,
			wantErr:    false,
		},
		{
			name:           "error - db error",
			userID:         3,
			permissionCode: "admin.user.read",
			mockSetup: func(m *mockPermissionRepository) {
				m.On("UserHasPermission", mock.Anything, 3, "admin.user.read").Return(false, errors.New("db error"))
			},
			wantResult: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockPermissionRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewPermissionService(mockRepo, zap.NewNop())

			result, err := svc.HasPermission(context.Background(), tt.userID, tt.permissionCode)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResult, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPermissionService_GetUserPermissions(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		mockSetup func(*mockPermissionRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:   "success - returns permissions",
			userID: 1,
			mockSetup: func(m *mockPermissionRepository) {
				codes := []string{"admin.user.read", "admin.user.write", "admin.role.read"}
				m.On("GetUserPermissionCodes", mock.Anything, 1).Return(codes, nil)
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:   "success - empty permissions",
			userID: 2,
			mockSetup: func(m *mockPermissionRepository) {
				m.On("GetUserPermissionCodes", mock.Anything, 2).Return([]string{}, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:   "error - db error",
			userID: 3,
			mockSetup: func(m *mockPermissionRepository) {
				m.On("GetUserPermissionCodes", mock.Anything, 3).Return(nil, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockPermissionRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewPermissionService(mockRepo, zap.NewNop())

			codes, err := svc.GetUserPermissions(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, codes, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
