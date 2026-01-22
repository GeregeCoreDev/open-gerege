//go:build integration

// Package integration contains integration tests
//
// File: system_repo_test.go
// Description: System repository integration tests
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

func TestSystemRepository_Create(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewSystemRepository(db)
	ctx := CreateTestContext()

	tests := []struct {
		name    string
		system  domain.System
		wantErr bool
	}{
		{
			name: "success - create new system",
			system: domain.System{
				Code:        "NEW_SYS",
				Name:        "New System",
				Description: "A new test system",
				IsActive:    boolPtr(true),
				Sequence:    1,
			},
			wantErr: false,
		},
		{
			name: "success - create system with minimal fields",
			system: domain.System{
				Code:     "MIN_SYS",
				Name:     "Minimal System",
				Sequence: 2,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.system)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify the system was created
			var created domain.System
			err = db.Where("code = ?", tt.system.Code).First(&created).Error
			require.NoError(t, err)
			assert.Equal(t, tt.system.Name, created.Name)
			assert.Equal(t, tt.system.Key, created.Key)
		})
	}
}

func TestSystemRepository_ByID(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewSystemRepository(db)
	ctx := CreateTestContext()

	// Seed a test system
	seededSystem := SeedTestSystem(t, db)

	tests := []struct {
		name     string
		systemID int
		wantErr  bool
	}{
		{
			name:     "success - system found",
			systemID: seededSystem.ID,
			wantErr:  false,
		},
		{
			name:     "error - system not found",
			systemID: 99999,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			system, err := repo.ByID(ctx, tt.systemID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.systemID, system.ID)
		})
	}
}

func TestSystemRepository_Update(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewSystemRepository(db)
	ctx := CreateTestContext()

	// Seed a test system
	seededSystem := SeedTestSystem(t, db)

	tests := []struct {
		name     string
		systemID int
		update   domain.System
		wantErr  bool
	}{
		{
			name:     "success - update system name",
			systemID: seededSystem.ID,
			update: domain.System{
				Code:        seededSystem.Code,
				Name:        "Updated System Name",
				Description: "Updated description",
			},
			wantErr: false,
		},
		{
			name:     "success - update sequence",
			systemID: seededSystem.ID,
			update: domain.System{
				Code:     seededSystem.Code,
				Sequence: 99,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.systemID, tt.update)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify the update
			updated, err := repo.ByID(ctx, tt.systemID)
			require.NoError(t, err)
			if tt.update.Name != "" {
				assert.Equal(t, tt.update.Name, updated.Name)
			}
		})
	}
}

func TestSystemRepository_Delete(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewSystemRepository(db)
	ctx := CreateTestContext()

	// Seed a test system
	seededSystem := SeedTestSystem(t, db)

	tests := []struct {
		name     string
		systemID int
		wantErr  bool
	}{
		{
			name:     "success - system deleted (soft delete)",
			systemID: seededSystem.ID,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.systemID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify soft delete - should not be found by normal query
			_, err = repo.ByID(ctx, tt.systemID)
			assert.Error(t, err)
		})
	}
}

func TestSystemRepository_List(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewSystemRepository(db)
	ctx := CreateTestContext()

	// Seed multiple systems
	for i := 0; i < 5; i++ {
		system := domain.System{
			Code:        "SYS_" + string(rune('A'+i)),
			Name:        "System " + string(rune('A'+i)),
			Description: "Test system",
			IsActive:    boolPtr(true),
			Sequence:    i + 1,
		}
		db.Create(&system)
	}

	tests := []struct {
		name         string
		query        dto.SystemListQuery
		wantMinItems int
		wantErr      bool
	}{
		{
			name: "success - list all systems",
			query: dto.SystemListQuery{
				PaginationQuery: common.PaginationQuery{
					Page: 1,
					Size: 10,
				},
			},
			wantMinItems: 5,
			wantErr:      false,
		},
		{
			name: "success - filter by code",
			query: dto.SystemListQuery{
				PaginationQuery: common.PaginationQuery{
					Page: 1,
					Size: 10,
				},
				Code: "SYS_A",
			},
			wantMinItems: 1,
			wantErr:      false,
		},
		{
			name: "success - filter by active status",
			query: dto.SystemListQuery{
				PaginationQuery: common.PaginationQuery{
					Page: 1,
					Size: 10,
				},
				IsActive: boolPtr(true),
			},
			wantMinItems: 5,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systems, total, _, _, err := repo.List(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(systems), tt.wantMinItems)
			assert.GreaterOrEqual(t, total, int64(tt.wantMinItems))
		})
	}
}

func TestSystemRepository_GetActiveModuleCount(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewSystemRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	system := SeedTestSystem(t, db)

	// Create modules for the system
	for i := 0; i < 3; i++ {
		module := domain.Module{
			SystemID:    system.ID,
			Code:        "MOD_" + string(rune('A'+i)),
			Name:        "Module " + string(rune('A'+i)),
			IsActive:    boolPtr(true),
		}
		db.Create(&module)
	}

	tests := []struct {
		name      string
		systemID  int
		wantCount int64
	}{
		{
			name:      "success - system with modules",
			systemID:  system.ID,
			wantCount: 3,
		},
		{
			name:      "success - system with no modules",
			systemID:  99999,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := repo.GetActiveModuleCount(ctx, tt.systemID)
			assert.Equal(t, tt.wantCount, count)
		})
	}
}

func TestSystemRepository_GetActiveRoleCount(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewSystemRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	system := SeedTestSystem(t, db)

	// Create roles for the system
	for i := 0; i < 2; i++ {
		role := domain.Role{
			SystemID:    system.ID,
			Code:        "ROLE_" + string(rune('A'+i)),
			Name:        "Role " + string(rune('A'+i)),
			IsActive:    boolPtr(true),
		}
		db.Create(&role)
	}

	tests := []struct {
		name      string
		systemID  int
		wantCount int64
	}{
		{
			name:      "success - system with roles",
			systemID:  system.ID,
			wantCount: 2,
		},
		{
			name:      "success - system with no roles",
			systemID:  99999,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := repo.GetActiveRoleCount(ctx, tt.systemID)
			assert.Equal(t, tt.wantCount, count)
		})
	}
}
