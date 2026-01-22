//go:build integration

// Package integration contains integration tests
package integration

import (
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestPermissionRepository_Create(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewPermissionRepository(db)
	ctx := CreateTestContext()

	// Seed dependencies
	system := SeedTestSystem(t, db)
	module := seedTestModule(t, db, system.ID)

	tests := []struct {
		name    string
		perm    domain.Permission
		wantErr bool
	}{
		{
			name: "success - create permission",
			perm: domain.Permission{
				ModuleID:    module.ID,
				Code:        "TEST_PERM_CREATE",
				Name:        "Test Permission Create",
				Description: "A test permission",
				IsActive:    boolPtr(true),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.perm)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify
			var created domain.Permission
			err = db.Where("code = ?", tt.perm.Code).First(&created).Error
			require.NoError(t, err)
			assert.Equal(t, tt.perm.Name, created.Name)
		})
	}
}

func TestPermissionRepository_ByID(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewPermissionRepository(db)
	ctx := CreateTestContext()

	// Seed
	system := SeedTestSystem(t, db)
	module := seedTestModule(t, db, system.ID)
	perm := seedTestPermission(t, db, module.ID)

	tests := []struct {
		name    string
		permID  int
		wantErr bool
	}{
		{
			name:    "success - found",
			permID:  perm.ID,
			wantErr: false,
		},
		{
			name:    "error - not found",
			permID:  99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.ByID(ctx, tt.permID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.permID, result.ID)
		})
	}
}

func TestPermissionRepository_ByCode(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewPermissionRepository(db)
	ctx := CreateTestContext()

	// Seed
	system := SeedTestSystem(t, db)
	module := seedTestModule(t, db, system.ID)
	perm := seedTestPermission(t, db, module.ID)

	tests := []struct {
		name     string
		permCode string
		wantErr  bool
	}{
		{
			name:     "success - found",
			permCode: perm.Code,
			wantErr:  false,
		},
		{
			name:     "error - not found",
			permCode: "NON_EXISTENT_CODE",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.ByCode(ctx, tt.permCode)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.permCode, result.Code)
		})
	}
}

func TestPermissionRepository_List(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewPermissionRepository(db)
	ctx := CreateTestContext()

	// Seed
	system := SeedTestSystem(t, db)
	module := seedTestModule(t, db, system.ID)

	// Create multiple permissions
	for i := 0; i < 5; i++ {
		perm := domain.Permission{
			ModuleID: module.ID,
			Code:     "PERM_LIST_" + string(rune('A'+i)),
			Name:     "Permission " + string(rune('A'+i)),
			IsActive: boolPtr(true),
		}
		db.Create(&perm)
	}

	tests := []struct {
		name         string
		query        dto.PermissionQuery
		wantMinItems int
		wantErr      bool
	}{
		{
			name: "success - list all",
			query: dto.PermissionQuery{
				PaginationQuery: common.PaginationQuery{
					Page: 1,
					Size: 10,
				},
			},
			wantMinItems: 5,
			wantErr:      false,
		},
		{
			name: "success - filter by module_id",
			query: dto.PermissionQuery{
				PaginationQuery: common.PaginationQuery{
					Page: 1,
					Size: 10,
				},
				ModuleID: module.ID,
			},
			wantMinItems: 5,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms, total, _, _, err := repo.List(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(perms), tt.wantMinItems)
			assert.GreaterOrEqual(t, total, int64(tt.wantMinItems))
		})
	}
}

func TestPermissionRepository_UserHasPermission(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewPermissionRepository(db)
	ctx := CreateTestContext()

	// Seed full permission chain
	user := SeedTestUser(t, db)
	system := SeedTestSystem(t, db)
	module := seedTestModule(t, db, system.ID)
	perm := seedTestPermission(t, db, module.ID)
	role := SeedTestRole(t, db, system.ID)

	// Assign permission to role
	db.Exec("INSERT INTO role_permissions (role_id, permission_id, created_date) VALUES (?, ?, NOW())", role.ID, perm.ID)

	// Assign role to user
	db.Create(&domain.UserRole{
		UserId: user.Id,
		RoleID: role.ID,
	})

	tests := []struct {
		name     string
		userID   int
		permCode string
		want     bool
		wantErr  bool
	}{
		{
			name:     "success - user has permission",
			userID:   user.Id,
			permCode: perm.Code,
			want:     true,
			wantErr:  false,
		},
		{
			name:     "success - user does not have permission",
			userID:   user.Id,
			permCode: "NON_EXISTENT_PERM",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "success - non-existent user",
			userID:   99999,
			permCode: perm.Code,
			want:     false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			has, err := repo.UserHasPermission(ctx, tt.userID, tt.permCode)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, has)
		})
	}
}

func TestPermissionRepository_GetUserPermissionCodes(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewPermissionRepository(db)
	ctx := CreateTestContext()

	// Seed
	user := SeedTestUser(t, db)
	system := SeedTestSystem(t, db)
	module := seedTestModule(t, db, system.ID)
	role := SeedTestRole(t, db, system.ID)

	// Create and assign permissions
	var permCodes []string
	for i := 0; i < 3; i++ {
		perm := domain.Permission{
			ModuleID: module.ID,
			Code:     "USER_PERM_" + string(rune('A'+i)),
			Name:     "User Permission " + string(rune('A'+i)),
			IsActive: boolPtr(true),
		}
		db.Create(&perm)
		permCodes = append(permCodes, perm.Code)
		db.Exec("INSERT INTO role_permissions (role_id, permission_id, created_date) VALUES (?, ?, NOW())", role.ID, perm.ID)
	}

	// Assign role to user
	db.Create(&domain.UserRole{
		UserId: user.Id,
		RoleID: role.ID,
	})

	tests := []struct {
		name          string
		userID        int
		wantMinCodes  int
		wantErr       bool
	}{
		{
			name:         "success - user with permissions",
			userID:       user.Id,
			wantMinCodes: 3,
			wantErr:      false,
		},
		{
			name:         "success - user without permissions",
			userID:       99999,
			wantMinCodes: 0,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codes, err := repo.GetUserPermissionCodes(ctx, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(codes), tt.wantMinCodes)
		})
	}
}

// Helper functions
func seedTestModule(t *testing.T, db *gorm.DB, systemID int) domain.Module {
	t.Helper()
	module := domain.Module{
		SystemID:    systemID,
		Code:        "TEST_MODULE",
		Name:        "Test Module",
		Description: "A test module",
		IsActive:    boolPtr(true),
	}
	if err := db.Create(&module).Error; err != nil {
		t.Fatalf("failed to seed test module: %v", err)
	}
	return module
}

func seedTestPermission(t *testing.T, db *gorm.DB, moduleID int) domain.Permission {
	t.Helper()
	perm := domain.Permission{
		ModuleID:    moduleID,
		Code:        "TEST_PERMISSION",
		Name:        "Test Permission",
		Description: "A test permission",
		IsActive:    boolPtr(true),
	}
	if err := db.Create(&perm).Error; err != nil {
		t.Fatalf("failed to seed test permission: %v", err)
	}
	return perm
}
