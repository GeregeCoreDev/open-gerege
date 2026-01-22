// Package service provides implementation for service
//
// File: organization_service_test.go
// Description: Unit tests for organization service
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

// mockOrganizationRepository for testing
type mockOrganizationRepository struct {
	mock.Mock
}

func (m *mockOrganizationRepository) List(ctx context.Context, p common.PaginationQuery) ([]domain.Organization, int64, int, int, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, 0, 0, 0, args.Error(4)
	}
	return args.Get(0).([]domain.Organization), args.Get(1).(int64), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

func (m *mockOrganizationRepository) Create(ctx context.Context, org domain.Organization) (domain.Organization, error) {
	args := m.Called(ctx, org)
	return args.Get(0).(domain.Organization), args.Error(1)
}

func (m *mockOrganizationRepository) Update(ctx context.Context, id int, org domain.Organization) (domain.Organization, error) {
	args := m.Called(ctx, id, org)
	return args.Get(0).(domain.Organization), args.Error(1)
}

func (m *mockOrganizationRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockOrganizationRepository) ByID(ctx context.Context, id int) (domain.Organization, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Organization), args.Error(1)
}

func (m *mockOrganizationRepository) Tree(ctx context.Context, rootID int) ([]domain.Organization, error) {
	args := m.Called(ctx, rootID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Organization), args.Error(1)
}

func TestOrganizationService_List(t *testing.T) {
	tests := []struct {
		name      string
		query     common.PaginationQuery
		mockSetup func(*mockOrganizationRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "success - returns organizations",
			query: common.PaginationQuery{Page: 1, Size: 10},
			mockSetup: func(m *mockOrganizationRepository) {
				orgs := []domain.Organization{
					{Id: 1, Name: "Org1"},
					{Id: 2, Name: "Org2"},
				}
				m.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return(orgs, int64(2), 1, 10, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:  "success - empty list",
			query: common.PaginationQuery{Page: 1, Size: 10},
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return([]domain.Organization{}, int64(0), 1, 10, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:  "error - db error",
			query: common.PaginationQuery{Page: 1, Size: 10},
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("List", mock.Anything, mock.AnythingOfType("common.PaginationQuery")).
					Return(nil, int64(0), 0, 0, errors.New("db error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockOrganizationRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewOrganizationService(mockRepo, zap.NewNop())

			orgs, _, _, _, err := svc.List(context.Background(), tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, orgs, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestOrganizationService_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     dto.OrganizationDto
		mockSetup func(*mockOrganizationRepository)
		wantErr   bool
	}{
		{
			name: "success - organization created",
			input: dto.OrganizationDto{
				Name:      "New Org",
				RegNo:     "1234567",
				ShortName: "NO",
			},
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.Organization")).
					Return(domain.Organization{Id: 1, Name: "New Org"}, nil)
			},
			wantErr: false,
		},
		{
			name: "success - short name defaults to name",
			input: dto.OrganizationDto{
				Name:  "New Org Without Short",
				RegNo: "1234567",
			},
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(org domain.Organization) bool {
					return org.ShortName == "New Org Without Short"
				})).Return(domain.Organization{Id: 1, Name: "New Org Without Short"}, nil)
			},
			wantErr: false,
		},
		{
			name: "error - create fails",
			input: dto.OrganizationDto{
				Name: "Fail Org",
			},
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("domain.Organization")).
					Return(domain.Organization{}, errors.New("create failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockOrganizationRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewOrganizationService(mockRepo, zap.NewNop())

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

func TestOrganizationService_Update(t *testing.T) {
	tests := []struct {
		name      string
		orgID     int
		input     dto.OrganizationUpdateDto
		mockSetup func(*mockOrganizationRepository)
		wantErr   bool
	}{
		{
			name:  "success - organization updated",
			orgID: 1,
			input: dto.OrganizationUpdateDto{
				Name:      "Updated Org",
				ShortName: "UO",
			},
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("Update", mock.Anything, 1, mock.AnythingOfType("domain.Organization")).
					Return(domain.Organization{Id: 1, Name: "Updated Org"}, nil)
			},
			wantErr: false,
		},
		{
			name:  "error - update fails",
			orgID: 999,
			input: dto.OrganizationUpdateDto{
				Name: "Fail Update",
			},
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("Update", mock.Anything, 999, mock.AnythingOfType("domain.Organization")).
					Return(domain.Organization{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockOrganizationRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewOrganizationService(mockRepo, zap.NewNop())

			_, err := svc.Update(context.Background(), tt.orgID, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestOrganizationService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		orgID     int
		mockSetup func(*mockOrganizationRepository)
		wantErr   bool
	}{
		{
			name:  "success - organization deleted",
			orgID: 1,
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("Delete", mock.Anything, 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "error - delete fails",
			orgID: 999,
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("Delete", mock.Anything, 999).Return(errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockOrganizationRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewOrganizationService(mockRepo, zap.NewNop())

			err := svc.Delete(context.Background(), tt.orgID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestOrganizationService_ByID(t *testing.T) {
	tests := []struct {
		name      string
		orgID     int
		mockSetup func(*mockOrganizationRepository)
		wantErr   bool
	}{
		{
			name:  "success - organization found",
			orgID: 1,
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("ByID", mock.Anything, 1).Return(domain.Organization{
					Id:   1,
					Name: "Test Org",
				}, nil)
			},
			wantErr: false,
		},
		{
			name:  "error - organization not found",
			orgID: 999,
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("ByID", mock.Anything, 999).Return(domain.Organization{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockOrganizationRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewOrganizationService(mockRepo, zap.NewNop())

			org, err := svc.ByID(context.Background(), tt.orgID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.orgID, org.Id)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestOrganizationService_Tree(t *testing.T) {
	tests := []struct {
		name      string
		rootID    int
		mockSetup func(*mockOrganizationRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:   "success - returns tree",
			rootID: 1,
			mockSetup: func(m *mockOrganizationRepository) {
				orgs := []domain.Organization{
					{Id: 1, Name: "Parent"},
					{Id: 2, Name: "Child1"},
					{Id: 3, Name: "Child2"},
				}
				m.On("Tree", mock.Anything, 1).Return(orgs, nil)
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:   "error - tree fetch fails",
			rootID: 999,
			mockSetup: func(m *mockOrganizationRepository) {
				m.On("Tree", mock.Anything, 999).Return(nil, errors.New("not found"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockOrganizationRepository{}
			tt.mockSetup(mockRepo)

			svc := service.NewOrganizationService(mockRepo, zap.NewNop())

			orgs, err := svc.Tree(context.Background(), tt.rootID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, orgs, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
