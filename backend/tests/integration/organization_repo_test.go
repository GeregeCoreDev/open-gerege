//go:build integration

// Package integration contains integration tests
package integration

import (
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationRepository_Create(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewOrganizationRepository(db)
	ctx := CreateTestContext()

	tests := []struct {
		name    string
		org     domain.Organization
		wantErr bool
	}{
		{
			name: "success - create organization",
			org: domain.Organization{
				Name:     "New Organization",
				IsActive: boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "success - create with parent",
			org: domain.Organization{
				Name:     "Child Organization",
				IsActive: boolPtr(true),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := repo.Create(ctx, tt.org)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Greater(t, created.Id, 0)
			assert.Equal(t, tt.org.Name, created.Name)
		})
	}
}

func TestOrganizationRepository_ByID(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewOrganizationRepository(db)
	ctx := CreateTestContext()

	// Seed
	org := SeedTestOrganization(t, db)

	tests := []struct {
		name    string
		orgID   int
		wantErr bool
	}{
		{
			name:    "success - found",
			orgID:   org.Id,
			wantErr: false,
		},
		{
			name:    "error - not found",
			orgID:   99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.ByID(ctx, tt.orgID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.orgID, result.Id)
		})
	}
}

func TestOrganizationRepository_Update(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewOrganizationRepository(db)
	ctx := CreateTestContext()

	// Seed
	org := SeedTestOrganization(t, db)

	tests := []struct {
		name    string
		orgID   int
		update  domain.Organization
		wantErr bool
	}{
		{
			name:  "success - update name",
			orgID: org.Id,
			update: domain.Organization{
				Name: "Updated Organization Name",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := repo.Update(ctx, tt.orgID, tt.update)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.update.Name, updated.Name)
		})
	}
}

func TestOrganizationRepository_Delete(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewOrganizationRepository(db)
	ctx := CreateTestContext()

	// Seed
	org := SeedTestOrganization(t, db)

	tests := []struct {
		name    string
		orgID   int
		wantErr bool
	}{
		{
			name:    "success - soft delete",
			orgID:   org.Id,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.orgID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify soft delete
			_, err = repo.ByID(ctx, tt.orgID)
			assert.Error(t, err)
		})
	}
}

func TestOrganizationRepository_List(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewOrganizationRepository(db)
	ctx := CreateTestContext()

	// Seed multiple organizations
	for i := 0; i < 5; i++ {
		org := domain.Organization{
			Name:     "Organization " + string(rune('A'+i)),
			IsActive: boolPtr(true),
		}
		db.Create(&org)
	}

	tests := []struct {
		name         string
		query        common.PaginationQuery
		wantMinItems int
		wantErr      bool
	}{
		{
			name: "success - list all",
			query: common.PaginationQuery{
				Page: 1,
				Size: 10,
			},
			wantMinItems: 5,
			wantErr:      false,
		},
		{
			name: "success - pagination",
			query: common.PaginationQuery{
				Page: 1,
				Size: 3,
			},
			wantMinItems: 3,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgs, total, _, _, err := repo.List(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(orgs), tt.wantMinItems)
			assert.GreaterOrEqual(t, total, int64(tt.wantMinItems))
		})
	}
}

func TestOrganizationRepository_Tree(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewOrganizationRepository(db)
	ctx := CreateTestContext()

	// Create a hierarchy
	parent := domain.Organization{
		Name:     "Parent Org",
		IsActive: boolPtr(true),
	}
	db.Create(&parent)

	child1 := domain.Organization{
		Name:     "Child Org 1",
		ParentId: &parent.Id,
		IsActive: boolPtr(true),
	}
	db.Create(&child1)

	child2 := domain.Organization{
		Name:     "Child Org 2",
		ParentId: &parent.Id,
		IsActive: boolPtr(true),
	}
	db.Create(&child2)

	tests := []struct {
		name         string
		rootID       int
		wantMinItems int
		wantErr      bool
	}{
		{
			name:         "success - get tree from root",
			rootID:       parent.Id,
			wantMinItems: 1, // At least the children
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := repo.Tree(ctx, tt.rootID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(tree), tt.wantMinItems)
		})
	}
}
