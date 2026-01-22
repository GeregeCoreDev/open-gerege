// Package service provides implementation for service
//
// File: system_service_test.go
// Description: Unit tests for system service
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

// mockSystemRepository for testing
type mockSystemRepository struct {
	mock.Mock
}

func (m *mockSystemRepository) List(ctx context.Context, q dto.SystemListQuery) ([]domain.System, int64, int, int, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.System), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockSystemRepository) Create(ctx context.Context, sys domain.System) error {
	args := m.Called(ctx, sys)
	return args.Error(0)
}

func (m *mockSystemRepository) Update(ctx context.Context, id int, sys domain.System) error {
	args := m.Called(ctx, id, sys)
	return args.Error(0)
}

func (m *mockSystemRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockSystemRepository) ByID(ctx context.Context, id int) (domain.System, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.System), args.Error(1)
}

func (m *mockSystemRepository) GetActiveModuleCount(ctx context.Context, systemID int) int64 {
	args := m.Called(ctx, systemID)
	return int64(args.Int(0))
}

func (m *mockSystemRepository) GetActiveRoleCount(ctx context.Context, systemID int) int64 {
	args := m.Called(ctx, systemID)
	return int64(args.Int(0))
}

func TestSystemService_List(t *testing.T) {
	tests := []struct {
		name      string
		query     dto.SystemListQuery
		mockSetup func(*mockSystemRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "success - returns systems",
			query: dto.SystemListQuery{PaginationQuery: common.PaginationQuery{Page: 1, Size: 10}},
			mockSetup: func(m *mockSystemRepository) {
				systems := []domain.System{
					{ID: 1, Name: "System1", Code: "SYS1"},
					{ID: 2, Name: "System2", Code: "SYS2"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.SystemListQuery")).
					Return(systems, int64(2), 1, 10, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:  "error - db error",
			query: dto.SystemListQuery{PaginationQuery: common.PaginationQuery{Page: 1, Size: 10}},
			mockSetup: func(m *mockSystemRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.SystemListQuery")).
					Return(nil, int64(0), 0, 0, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockSystemRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewSystemService(mockRepo, zap.NewNop())

			systems, _, _, _, err := svc.List(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, systems, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSystemService_ByID(t *testing.T) {
	tests := []struct {
		name      string
		systemID  int
		mockSetup func(*mockSystemRepository)
		wantErr   bool
	}{
		{
			name:     "success - system found",
			systemID: 1,
			mockSetup: func(m *mockSystemRepository) {
				m.On("ByID", mock.Anything, 1).Return(domain.System{
					ID:   1,
					Name: "Test System",
					Code: "TEST",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:     "error - system not found",
			systemID: 999,
			mockSetup: func(m *mockSystemRepository) {
				m.On("ByID", mock.Anything, 999).Return(domain.System{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockSystemRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewSystemService(mockRepo, zap.NewNop())

			sys, err := svc.ByID(context.Background(), tt.systemID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.systemID, sys.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSystemService_Create(t *testing.T) {
	isActive := true

	tests := []struct {
		name      string
		input     dto.SystemCreateDto
		mockSetup func(*mockSystemRepository)
		wantErr   bool
	}{
		{
			name: "success - system created",
			input: dto.SystemCreateDto{
				Code:     "NEW_SYS",
				Name:     "New System",
				IsActive: &isActive,
			},
			mockSetup: func(m *mockSystemRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(sys domain.System) bool {
					return sys.Code == "new_sys" && sys.Key == "new_sys" // code is lowercased
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success - key defaults to code",
			input: dto.SystemCreateDto{
				Code: "TEST_CODE",
				Name: "Test System",
			},
			mockSetup: func(m *mockSystemRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(sys domain.System) bool {
					return sys.Key == "test_code"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - create fails",
			input: dto.SystemCreateDto{
				Code: "FAIL",
				Name: "Fail System",
			},
			mockSetup: func(m *mockSystemRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.System")).
					Return(errors.New("duplicate code"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockSystemRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewSystemService(mockRepo, zap.NewNop())

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

func TestSystemService_Update(t *testing.T) {
	isActiveFalse := false

	tests := []struct {
		name      string
		systemID  int
		input     dto.SystemUpdateDto
		mockSetup func(*mockSystemRepository)
		wantErr   bool
	}{
		{
			name:     "success - system updated",
			systemID: 1,
			input: dto.SystemUpdateDto{
				Code: "UPDATED",
				Name: "Updated System",
			},
			mockSetup: func(m *mockSystemRepository) {
				m.On("Update", mock.Anything, 1, mock.AnythingOfType("domain.System")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "error - cannot deactivate system with active modules",
			systemID: 2,
			input: dto.SystemUpdateDto{
				Code:     "SYS",
				Name:     "System with modules",
				IsActive: &isActiveFalse,
			},
			mockSetup: func(m *mockSystemRepository) {
				m.On("GetActiveModuleCount", mock.Anything, 2).Return(5)
			},
			wantErr: true,
		},
		{
			name:     "error - cannot deactivate system with active roles",
			systemID: 3,
			input: dto.SystemUpdateDto{
				Code:     "SYS",
				Name:     "System with roles",
				IsActive: &isActiveFalse,
			},
			mockSetup: func(m *mockSystemRepository) {
				m.On("GetActiveModuleCount", mock.Anything, 3).Return(0)
				m.On("GetActiveRoleCount", mock.Anything, 3).Return(3)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockSystemRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewSystemService(mockRepo, zap.NewNop())

			err := svc.Update(context.Background(), tt.systemID, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSystemService_Delete(t *testing.T) {
	isActiveFalse := false
	isActiveTrue := true

	tests := []struct {
		name      string
		systemID  int
		mockSetup func(*mockSystemRepository)
		wantErr   bool
	}{
		{
			name:     "success - system deleted",
			systemID: 1,
			mockSetup: func(m *mockSystemRepository) {
				m.On("ByID", mock.Anything, 1).Return(domain.System{ID: 1, IsActive: &isActiveFalse}, nil)
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "error - cannot delete active system",
			systemID: 2,
			mockSetup: func(m *mockSystemRepository) {
				m.On("ByID", mock.Anything, 2).Return(domain.System{ID: 2, IsActive: &isActiveTrue}, nil)
			},
			wantErr: true,
		},
		{
			name:     "error - system not found",
			systemID: 999,
			mockSetup: func(m *mockSystemRepository) {
				m.On("ByID", mock.Anything, 999).Return(domain.System{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockSystemRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewSystemService(mockRepo, zap.NewNop())

			err := svc.Delete(context.Background(), tt.systemID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
