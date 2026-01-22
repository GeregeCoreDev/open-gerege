//go:build integration

// Package integration contains integration tests
package integration

import (
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/http/dto"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/common"
	"git.gerege.mn/backend-packages/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModuleRepository_Create(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewModuleRepository(db, &config.Config{})
	ctx := CreateTestContext()

	// Seed system
	system := SeedTestSystem(t, db)

	tests := []struct {
		name    string
		module  domain.Module
		wantErr bool
	}{
		{
			name: "success - create module",
			module: domain.Module{
				SystemID:    system.ID,
				Code:        "NEW_MODULE",
				Name:        "New Module",
				Description: "A new test module",
				IsActive:    boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "success - create with minimal fields",
			module: domain.Module{
				SystemID: system.ID,
				Code:     "MIN_MODULE",
				Name:     "Minimal Module",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.module)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify
			var created domain.Module
			err = db.Where("code = ?", tt.module.Code).First(&created).Error
			require.NoError(t, err)
			assert.Equal(t, tt.module.Name, created.Name)
		})
	}
}

func TestModuleRepository_ByID(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewModuleRepository(db, &config.Config{})
	ctx := CreateTestContext()

	// Seed
	system := SeedTestSystem(t, db)
	module := seedTestModule(t, db, system.ID)

	tests := []struct {
		name     string
		moduleID int
		wantErr  bool
	}{
		{
			name:     "success - found",
			moduleID: module.ID,
			wantErr:  false,
		},
		{
			name:     "error - not found",
			moduleID: 99999,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.ByID(ctx, tt.moduleID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.moduleID, result.ID)
		})
	}
}

func TestModuleRepository_Update(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewModuleRepository(db, &config.Config{})
	ctx := CreateTestContext()

	// Seed
	system := SeedTestSystem(t, db)
	module := seedTestModule(t, db, system.ID)

	tests := []struct {
		name     string
		moduleID int
		update   domain.Module
		wantErr  bool
	}{
		{
			name:     "success - update name",
			moduleID: module.ID,
			update: domain.Module{
				Name:        "Updated Module Name",
				Description: "Updated description",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.moduleID, tt.update)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify
			updated, err := repo.ByID(ctx, tt.moduleID)
			require.NoError(t, err)
			assert.Equal(t, tt.update.Name, updated.Name)
		})
	}
}

func TestModuleRepository_Delete(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewModuleRepository(db, &config.Config{})
	ctx := CreateTestContext()

	// Seed
	system := SeedTestSystem(t, db)
	module := seedTestModule(t, db, system.ID)

	tests := []struct {
		name     string
		moduleID int
		wantErr  bool
	}{
		{
			name:     "success - soft delete",
			moduleID: module.ID,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.moduleID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify soft delete
			_, err = repo.ByID(ctx, tt.moduleID)
			assert.Error(t, err)
		})
	}
}

func TestModuleRepository_List(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewModuleRepository(db, &config.Config{})
	ctx := CreateTestContext()

	// Seed
	system := SeedTestSystem(t, db)

	// Create multiple modules
	for i := 0; i < 5; i++ {
		mod := domain.Module{
			SystemID: system.ID,
			Code:     "MOD_LIST_" + string(rune('A'+i)),
			Name:     "Module " + string(rune('A'+i)),
			IsActive: boolPtr(true),
		}
		db.Create(&mod)
	}

	tests := []struct {
		name         string
		query        dto.ModuleListQuery
		wantMinItems int
		wantErr      bool
	}{
		{
			name: "success - list all",
			query: dto.ModuleListQuery{
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
			query: dto.ModuleListQuery{
				PaginationQuery: common.PaginationQuery{
					Page: 1,
					Size: 10,
				},
				SystemID: system.ID,
			},
			wantMinItems: 5,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modules, total, _, _, err := repo.List(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(modules), tt.wantMinItems)
			assert.GreaterOrEqual(t, total, int64(tt.wantMinItems))
		})
	}
}
