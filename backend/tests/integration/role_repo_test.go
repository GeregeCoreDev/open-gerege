//go:build integration

// Package integration contains integration tests
//
// File: role_repo_test.go
// Description: Role repository integration tests
package integration

import (
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleRepository_Create(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewRoleRepository(db)
	ctx := CreateTestContext()

	// Seed a system first (role depends on system)
	system := SeedTestSystem(t, db)

	tests := []struct {
		name    string
		role    domain.Role
		wantErr bool
	}{
		{
			name: "success - create new role",
			role: domain.Role{
				SystemID:    system.ID,
				Code:        "NEW_ROLE",
				Name:        "New Role",
				Description: "A new test role",
				IsActive:    boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "success - create role with minimal fields",
			role: domain.Role{
				SystemID: system.ID,
				Code:     "MINIMAL_ROLE",
				Name:     "Minimal Role",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify the role was created by fetching it
			var createdRole domain.Role
			err = db.Where("code = ?", tt.role.Code).First(&createdRole).Error
			require.NoError(t, err)
			assert.Equal(t, tt.role.Name, createdRole.Name)
		})
	}
}

func TestRoleRepository_ByID(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewRoleRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	system := SeedTestSystem(t, db)
	seededRole := SeedTestRole(t, db, system.ID)

	tests := []struct {
		name    string
		roleID  int
		wantErr bool
	}{
		{
			name:    "success - role found",
			roleID:  seededRole.ID,
			wantErr: false,
		},
		{
			name:    "error - role not found",
			roleID:  99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := repo.ByID(ctx, tt.roleID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.roleID, role.ID)
		})
	}
}

func TestRoleRepository_Update(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewRoleRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	system := SeedTestSystem(t, db)
	seededRole := SeedTestRole(t, db, system.ID)

	tests := []struct {
		name    string
		roleID  int
		update  domain.Role
		wantErr bool
	}{
		{
			name:   "success - update role name",
			roleID: seededRole.ID,
			update: domain.Role{
				Name:        "Updated Role Name",
				Description: "Updated description",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.roleID, tt.update)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify the update
			updated, err := repo.ByID(ctx, tt.roleID)
			require.NoError(t, err)
			assert.Equal(t, tt.update.Name, updated.Name)
		})
	}
}

func TestRoleRepository_Delete(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewRoleRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	system := SeedTestSystem(t, db)
	seededRole := SeedTestRole(t, db, system.ID)

	tests := []struct {
		name    string
		roleID  int
		wantErr bool
	}{
		{
			name:    "success - role deleted (soft delete)",
			roleID:  seededRole.ID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.roleID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify soft delete - should not be found by normal query
			_, err = repo.ByID(ctx, tt.roleID)
			assert.Error(t, err)
		})
	}
}

func TestRoleRepository_List(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewRoleRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	system := SeedTestSystem(t, db)
	for i := 0; i < 5; i++ {
		isActive := true
		role := domain.Role{
			SystemID:    system.ID,
			Code:        "ROLE_" + string(rune('A'+i)),
			Name:        "Role " + string(rune('A'+i)),
			Description: "Test role",
			IsActive:    &isActive,
		}
		db.Create(&role)
	}

	tests := []struct {
		name         string
		query        dto.RoleListQuery
		wantMinItems int
		wantErr      bool
	}{
		{
			name: "success - list all roles",
			query: dto.RoleListQuery{
				PaginationQuery: common.PaginationQuery{
					Page: 1,
					Size: 10,
				},
			},
			wantMinItems: 5,
			wantErr:      false,
		},
		{
			name: "success - filter by system_id",
			query: dto.RoleListQuery{
				PaginationQuery: common.PaginationQuery{
					Page: 1,
					Size: 10,
				},
				SystemId: system.ID,
			},
			wantMinItems: 5,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roles, total, _, _, err := repo.List(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(roles), tt.wantMinItems)
			assert.GreaterOrEqual(t, total, int64(tt.wantMinItems))
		})
	}
}

func TestRoleRepository_ReplacePermissions(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewRoleRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	system := SeedTestSystem(t, db)
	role := SeedTestRole(t, db, system.ID)

	// Create module first
	module := domain.Module{
		SystemID:    system.ID,
		Code:        "TEST_MOD",
		Name:        "Test Module",
		IsActive:    boolPtr(true),
	}
	db.Create(&module)

	// Create permissions
	perm1 := domain.Permission{
		ModuleID:    module.ID,
		Code:        "PERM_1",
		Name:        "Permission 1",
		IsActive:    boolPtr(true),
	}
	perm2 := domain.Permission{
		ModuleID:    module.ID,
		Code:        "PERM_2",
		Name:        "Permission 2",
		IsActive:    boolPtr(true),
	}
	db.Create(&perm1)
	db.Create(&perm2)

	tests := []struct {
		name    string
		roleID  int
		permIDs []int
		wantErr bool
	}{
		{
			name:    "success - assign permissions to role",
			roleID:  role.ID,
			permIDs: []int{perm1.ID, perm2.ID},
			wantErr: false,
		},
		{
			name:    "success - replace with single permission",
			roleID:  role.ID,
			permIDs: []int{perm1.ID},
			wantErr: false,
		},
		{
			name:    "success - clear all permissions",
			roleID:  role.ID,
			permIDs: []int{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.ReplacePermissions(ctx, tt.roleID, tt.permIDs)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify the permissions
			permissions, err := repo.Permissions(ctx, dto.RolePermissionsQuery{
				RoleID: tt.roleID,
			})
			require.NoError(t, err)
			assert.Equal(t, len(tt.permIDs), len(permissions))
		})
	}
}

func TestRoleRepository_GetUserCount(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewRoleRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	system := SeedTestSystem(t, db)
	role := SeedTestRole(t, db, system.ID)
	user := SeedTestUser(t, db)

	// Assign user to role
	db.Create(&domain.UserRole{
		UserId: user.Id,
		RoleID: role.ID,
	})

	tests := []struct {
		name      string
		roleID    int
		wantCount int64
	}{
		{
			name:      "success - role with one user",
			roleID:    role.ID,
			wantCount: 1,
		},
		{
			name:      "success - role with no users",
			roleID:    99999,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := repo.GetUserCount(ctx, tt.roleID)
			assert.Equal(t, tt.wantCount, count)
		})
	}
}

// Helper function for bool pointer
func boolPtr(b bool) *bool {
	return &b
}
