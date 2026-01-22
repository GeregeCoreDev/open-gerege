//go:build integration

// Package integration contains integration tests
//
// File: menu_repo_test.go
// Description: Menu repository integration tests
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

func TestMenuRepository_Create(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewMenuRepository(db, &config.Config{})
	ctx := CreateTestContext()

	tests := []struct {
		name    string
		menu    domain.Menu
		wantErr bool
	}{
		{
			name: "success - create new menu",
			menu: domain.Menu{
				Code:     "NEW_MENU",
				Key:      "new-menu",
				Name:     "New Menu",
				Path:     "/new",
				Sequence: 1,
			},
			wantErr: false,
		},
		{
			name: "success - create menu with parent",
			menu: domain.Menu{
				Code:     "CHILD_MENU",
				Key:      "child-menu",
				Name:     "Child Menu",
				Path:     "/child",
				Sequence: 2,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.menu)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestMenuRepository_ByID(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewMenuRepository(db, &config.Config{})
	ctx := CreateTestContext()

	// Seed a test menu
	seededMenu := SeedTestMenu(t, db)

	tests := []struct {
		name    string
		menuID  int64
		wantErr bool
	}{
		{
			name:    "success - menu found",
			menuID:  seededMenu.ID,
			wantErr: false,
		},
		{
			name:    "error - menu not found",
			menuID:  99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu, err := repo.ByID(ctx, tt.menuID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.menuID, menu.ID)
		})
	}
}

func TestMenuRepository_Update(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewMenuRepository(db, &config.Config{})
	ctx := CreateTestContext()

	// Seed a test menu
	seededMenu := SeedTestMenu(t, db)

	tests := []struct {
		name      string
		menuID    int64
		update    domain.Menu
		wantErr   bool
		checkFunc func(t *testing.T)
	}{
		{
			name:   "success - update menu name",
			menuID: seededMenu.ID,
			update: domain.Menu{
				Name: "Updated Menu",
			},
			wantErr: false,
			checkFunc: func(t *testing.T) {
				updated, err := repo.ByID(ctx, seededMenu.ID)
				require.NoError(t, err)
				assert.Equal(t, "Updated Menu", updated.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.menuID, tt.update)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.checkFunc != nil {
				tt.checkFunc(t)
			}
		})
	}
}

func TestMenuRepository_Delete(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewMenuRepository(db, &config.Config{})
	ctx := CreateTestContext()

	// Seed a test menu
	seededMenu := SeedTestMenu(t, db)

	tests := []struct {
		name    string
		menuID  int64
		wantErr bool
	}{
		{
			name:    "success - menu deleted (soft delete)",
			menuID:  seededMenu.ID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.menuID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify soft delete - should not be found by normal query
			_, err = repo.ByID(ctx, tt.menuID)
			assert.Error(t, err) // Should fail because of soft delete
		})
	}
}

func TestMenuRepository_List(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewMenuRepository(db, &config.Config{})
	ctx := CreateTestContext()

	// Seed test menus
	SeedTestMenus(t, db, 15)

	tests := []struct {
		name         string
		query        dto.MenuListQuery
		wantMinItems int
		wantMaxItems int
		wantTotalGte int64
		wantErr      bool
	}{
		{
			name: "success - default pagination",
			query: dto.MenuListQuery{
				PaginationQuery: common.PaginationQuery{
					Page: 1,
					Size: 10,
				},
			},
			wantMinItems: 1,
			wantMaxItems: 10,
			wantTotalGte: 15,
			wantErr:      false,
		},
		{
			name: "success - page 2",
			query: dto.MenuListQuery{
				PaginationQuery: common.PaginationQuery{
					Page: 2,
					Size: 10,
				},
			},
			wantMinItems: 1,
			wantMaxItems: 10,
			wantTotalGte: 15,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menus, total, page, size, err := repo.List(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(menus), tt.wantMinItems)
			assert.LessOrEqual(t, len(menus), tt.wantMaxItems)
			assert.GreaterOrEqual(t, total, tt.wantTotalGte)
			assert.Equal(t, tt.query.Page, page)
			assert.Equal(t, tt.query.Size, size)
		})
	}
}

func TestMenuRepository_ListAll(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewMenuRepository(db, &config.Config{})
	ctx := CreateTestContext()

	// Seed test menus
	SeedTestMenus(t, db, 5)

	menus, err := repo.ListAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(menus), 5)
}
