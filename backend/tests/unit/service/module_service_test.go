// Package service provides implementation for service
//
// File: module_service_test.go
// Description: Unit tests for module service
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
)

// mockModuleRepository implements repository.ModuleRepository
type mockModuleRepository struct {
	mock.Mock
}

func (m *mockModuleRepository) List(ctx context.Context, q dto.ModuleListQuery) ([]domain.Module, int64, int, int, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.Module), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockModuleRepository) ByID(ctx context.Context, id int) (domain.Module, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Module), args.Error(1)
}

func (m *mockModuleRepository) Create(ctx context.Context, module domain.Module) error {
	args := m.Called(ctx, module)
	return args.Error(0)
}

func (m *mockModuleRepository) Update(ctx context.Context, id int, module domain.Module) error {
	args := m.Called(ctx, id, module)
	return args.Error(0)
}

func (m *mockModuleRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestModuleService_List(t *testing.T) {
	tests := []struct {
		name      string
		query     dto.ModuleListQuery
		mockSetup func(*mockModuleRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "success - returns modules",
			query: dto.ModuleListQuery{},
			mockSetup: func(m *mockModuleRepository) {
				modules := []domain.Module{
					{ID: 1, Code: "user", Name: "User Module"},
					{ID: 2, Code: "role", Name: "Role Module"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.ModuleListQuery")).
					Return(modules, int64(2), 1, 10, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:  "success - empty list",
			query: dto.ModuleListQuery{},
			mockSetup: func(m *mockModuleRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.ModuleListQuery")).
					Return([]domain.Module{}, int64(0), 1, 10, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:  "error - db error",
			query: dto.ModuleListQuery{},
			mockSetup: func(m *mockModuleRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.ModuleListQuery")).
					Return(nil, int64(0), 0, 0, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockModuleRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewModuleService(mockRepo)

			modules, _, _, _, err := svc.List(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, modules, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestModuleService_ByID(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		mockSetup func(*mockModuleRepository)
		wantErr   bool
	}{
		{
			name: "success - found",
			id:   1,
			mockSetup: func(m *mockModuleRepository) {
				m.On("ByID", mock.Anything, 1).
					Return(domain.Module{ID: 1, Code: "user", Name: "User Module"}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - not found",
			id:   999,
			mockSetup: func(m *mockModuleRepository) {
				m.On("ByID", mock.Anything, 999).
					Return(domain.Module{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockModuleRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewModuleService(mockRepo)

			module, err := svc.ByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, module.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestModuleService_Create(t *testing.T) {
	isActive := true
	tests := []struct {
		name      string
		input     dto.ModuleCreateDto
		mockSetup func(*mockModuleRepository)
		wantErr   bool
	}{
		{
			name: "success - created with lowercase code",
			input: dto.ModuleCreateDto{
				Code:        "USER",
				Name:        "User Module",
				Description: "Manages users",
				IsActive:    &isActive,
				SystemID:    1,
			},
			mockSetup: func(m *mockModuleRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(module domain.Module) bool {
					return module.Code == "user" && module.Name == "User Module"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - create fails",
			input: dto.ModuleCreateDto{
				Code:     "fail",
				Name:     "Fail Module",
				SystemID: 1,
			},
			mockSetup: func(m *mockModuleRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.Module")).
					Return(errors.New("create failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockModuleRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewModuleService(mockRepo)

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

func TestModuleService_Update(t *testing.T) {
	isActive := true
	tests := []struct {
		name      string
		id        int
		input     dto.ModuleUpdateDto
		mockSetup func(*mockModuleRepository)
		wantErr   bool
	}{
		{
			name: "success - updated with lowercase code",
			id:   1,
			input: dto.ModuleUpdateDto{
				Code:        "UPDATED",
				Name:        "Updated Module",
				Description: "Updated description",
				IsActive:    &isActive,
				SystemID:    1,
			},
			mockSetup: func(m *mockModuleRepository) {
				m.On("Update", mock.Anything, 1, mock.MatchedBy(func(module domain.Module) bool {
					return module.Code == "updated" && module.Name == "Updated Module"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - update fails",
			id:   999,
			input: dto.ModuleUpdateDto{
				Code: "fail",
				Name: "Fail",
			},
			mockSetup: func(m *mockModuleRepository) {
				m.On("Update", mock.Anything, 999, mock.AnythingOfType("domain.Module")).
					Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockModuleRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewModuleService(mockRepo)

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

func TestModuleService_Delete(t *testing.T) {
	isActive := true
	isInactive := false

	tests := []struct {
		name      string
		id        int
		mockSetup func(*mockModuleRepository)
		wantErr   bool
	}{
		{
			name: "success - deleted inactive module",
			id:   1,
			mockSetup: func(m *mockModuleRepository) {
				m.On("ByID", mock.Anything, 1).
					Return(domain.Module{ID: 1, Code: "inactive", IsActive: &isInactive}, nil)
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success - deleted module with nil IsActive",
			id:   2,
			mockSetup: func(m *mockModuleRepository) {
				m.On("ByID", mock.Anything, 2).
					Return(domain.Module{ID: 2, Code: "nilactive", IsActive: nil}, nil)
				m.On("Delete", mock.Anything, 2).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - cannot delete active module",
			id:   3,
			mockSetup: func(m *mockModuleRepository) {
				m.On("ByID", mock.Anything, 3).
					Return(domain.Module{ID: 3, Code: "active", IsActive: &isActive}, nil)
			},
			wantErr: true,
		},
		{
			name: "error - module not found",
			id:   999,
			mockSetup: func(m *mockModuleRepository) {
				m.On("ByID", mock.Anything, 999).
					Return(domain.Module{}, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "error - delete fails",
			id:   4,
			mockSetup: func(m *mockModuleRepository) {
				m.On("ByID", mock.Anything, 4).
					Return(domain.Module{ID: 4, IsActive: &isInactive}, nil)
				m.On("Delete", mock.Anything, 4).Return(errors.New("delete failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockModuleRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewModuleService(mockRepo)

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
