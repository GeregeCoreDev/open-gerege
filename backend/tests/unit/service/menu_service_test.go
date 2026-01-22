// Package service provides implementation for service
//
// File: menu_service_test.go
// Description: Unit tests for menu service
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

// mockMenuRepository implements repository.MenuRepository for testing
type mockMenuRepository struct {
	mock.Mock
}

func (m *mockMenuRepository) List(ctx context.Context, q dto.MenuListQuery) ([]domain.Menu, int64, int, int, error) {
	args := m.Called(ctx, q)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.Menu), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockMenuRepository) ListAll(ctx context.Context) ([]domain.Menu, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Menu), args.Error(1)
}

func (m *mockMenuRepository) ListByUserRoles(ctx context.Context, userID int) ([]domain.Menu, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Menu), args.Error(1)
}

func (m *mockMenuRepository) GetMenusByIDs(ctx context.Context, ids []int64) ([]domain.Menu, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Menu), args.Error(1)
}

func (m *mockMenuRepository) GetMenusByPermissionIDs(ctx context.Context, permissionIDs []int) ([]domain.Menu, error) {
	args := m.Called(ctx, permissionIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Menu), args.Error(1)
}

func (m *mockMenuRepository) ByID(ctx context.Context, id int64) (domain.Menu, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Menu), args.Error(1)
}

func (m *mockMenuRepository) Create(ctx context.Context, menu domain.Menu) error {
	args := m.Called(ctx, menu)
	return args.Error(0)
}

func (m *mockMenuRepository) Update(ctx context.Context, id int64, menu domain.Menu) error {
	args := m.Called(ctx, id, menu)
	return args.Error(0)
}

func (m *mockMenuRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestMenuService_List(t *testing.T) {
	tests := []struct {
		name      string
		query     dto.MenuListQuery
		mockSetup func(*mockMenuRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name: "success - returns menus",
			query: dto.MenuListQuery{
				PaginationQuery: common.PaginationQuery{Page: 1, Size: 10},
			},
			mockSetup: func(m *mockMenuRepository) {
				menus := []domain.Menu{
					{ID: 1, Name: "Menu 1"},
					{ID: 2, Name: "Menu 2"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("dto.MenuListQuery")).
					Return(menus, int64(2), 1, 10, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "error - db error",
			query: dto.MenuListQuery{
				PaginationQuery: common.PaginationQuery{Page: 1, Size: 10},
			},
			mockSetup: func(m *mockMenuRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("dto.MenuListQuery")).
					Return(nil, int64(0), 0, 0, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockMenuRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewMenuService(mockRepo)
			menus, _, _, _, err := svc.List(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, menus, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMenuService_ListAll(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*mockMenuRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name: "success - returns all menus",
			mockSetup: func(m *mockMenuRepository) {
				menus := []domain.Menu{
					{ID: 1, Name: "Menu 1"},
					{ID: 2, Name: "Menu 2"},
					{ID: 3, Name: "Menu 3"},
				}
				m.On("ListAll", mock.Anything).Return(menus, nil)
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name: "success - empty list",
			mockSetup: func(m *mockMenuRepository) {
				m.On("ListAll", mock.Anything).Return([]domain.Menu{}, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockMenuRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewMenuService(mockRepo)
			menus, err := svc.ListAll(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, menus, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMenuService_ListByUserRoles(t *testing.T) {
	parentID := int64(1)

	tests := []struct {
		name      string
		userID    int
		mockSetup func(*mockMenuRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:   "success - user with menus",
			userID: 1,
			mockSetup: func(m *mockMenuRepository) {
				menus := []domain.Menu{
					{ID: 2, Name: "Child Menu", ParentID: &parentID, Sequence: 1},
				}
				parentMenus := []domain.Menu{
					{ID: 1, Name: "Parent Menu", ParentID: nil, Sequence: 1},
				}
				m.On("ListByUserRoles", mock.Anything, 1).Return(menus, nil)
				m.On("GetMenusByIDs", mock.Anything, mock.AnythingOfType("[]int64")).Return(parentMenus, nil)
			},
			wantCount: 1, // Root menu count
			wantErr:   false,
		},
		{
			name:   "success - user without menus",
			userID: 2,
			mockSetup: func(m *mockMenuRepository) {
				m.On("ListByUserRoles", mock.Anything, 2).Return([]domain.Menu{}, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:   "error - db error",
			userID: 3,
			mockSetup: func(m *mockMenuRepository) {
				m.On("ListByUserRoles", mock.Anything, 3).Return(nil, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockMenuRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewMenuService(mockRepo)
			menus, err := svc.ListByUserRoles(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, menus, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMenuService_ByID(t *testing.T) {
	tests := []struct {
		name      string
		menuID    int64
		mockSetup func(*mockMenuRepository)
		wantErr   bool
	}{
		{
			name:   "success - menu found",
			menuID: 1,
			mockSetup: func(m *mockMenuRepository) {
				m.On("ByID", mock.Anything, int64(1)).Return(domain.Menu{
					ID:   1,
					Name: "Test Menu",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:   "error - menu not found",
			menuID: 999,
			mockSetup: func(m *mockMenuRepository) {
				m.On("ByID", mock.Anything, int64(999)).Return(domain.Menu{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockMenuRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewMenuService(mockRepo)
			menu, err := svc.ByID(context.Background(), tt.menuID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.menuID, menu.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMenuService_Create(t *testing.T) {
	zeroParent := int64(0)
	validParent := int64(1)

	tests := []struct {
		name      string
		input     dto.MenuCreateDto
		mockSetup func(*mockMenuRepository)
		wantErr   bool
	}{
		{
			name: "success - menu created",
			input: dto.MenuCreateDto{
				Code:     "NEW_MENU",
				Key:      "new-menu",
				Name:     "New Menu",
				Sequence: 1,
			},
			mockSetup: func(m *mockMenuRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.Menu")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success - menu with parent",
			input: dto.MenuCreateDto{
				Code:     "CHILD_MENU",
				Key:      "child-menu",
				Name:     "Child Menu",
				ParentID: &validParent,
				Sequence: 1,
			},
			mockSetup: func(m *mockMenuRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.Menu")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "success - zero parent converted to nil",
			input: dto.MenuCreateDto{
				Code:     "ROOT_MENU",
				Key:      "root-menu",
				Name:     "Root Menu",
				ParentID: &zeroParent,
				Sequence: 1,
			},
			mockSetup: func(m *mockMenuRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(menu domain.Menu) bool {
					return menu.ParentID == nil
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error - create fails",
			input: dto.MenuCreateDto{
				Code: "FAIL_MENU",
				Name: "Fail Menu",
			},
			mockSetup: func(m *mockMenuRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.Menu")).
					Return(errors.New("create failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockMenuRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewMenuService(mockRepo)
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

func TestMenuService_Update(t *testing.T) {
	tests := []struct {
		name      string
		menuID    int64
		input     dto.MenuUpdateDto
		mockSetup func(*mockMenuRepository)
		wantErr   bool
	}{
		{
			name:   "success - menu updated",
			menuID: 1,
			input: dto.MenuUpdateDto{
				Name: "Updated Menu",
			},
			mockSetup: func(m *mockMenuRepository) {
				m.On("Update", mock.Anything, int64(1), mock.AnythingOfType("domain.Menu")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error - update fails",
			menuID: 999,
			input: dto.MenuUpdateDto{
				Name: "Fail Update",
			},
			mockSetup: func(m *mockMenuRepository) {
				m.On("Update", mock.Anything, int64(999), mock.AnythingOfType("domain.Menu")).
					Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockMenuRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewMenuService(mockRepo)
			err := svc.Update(context.Background(), tt.menuID, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestMenuService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		menuID    int64
		mockSetup func(*mockMenuRepository)
		wantErr   bool
	}{
		{
			name:   "success - menu deleted",
			menuID: 1,
			mockSetup: func(m *mockMenuRepository) {
				m.On("Delete", mock.Anything, int64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "error - delete fails",
			menuID: 999,
			mockSetup: func(m *mockMenuRepository) {
				m.On("Delete", mock.Anything, int64(999)).Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockMenuRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewMenuService(mockRepo)
			err := svc.Delete(context.Background(), tt.menuID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
