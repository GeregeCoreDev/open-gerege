//go:build integration

// Package integration contains integration tests
//
// File: user_repo_test.go
// Description: User repository integration tests
package integration

import (
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewUserRepository(db)
	ctx := CreateTestContext()

	tests := []struct {
		name    string
		user    domain.User
		wantErr bool
	}{
		{
			name: "success - create new user",
			user: domain.User{
				RegNo:     "ZZ12345678",
				FirstName: "New",
				LastName:  "User",
				Email:     "newuser@example.com",
				PhoneNo:   "88001122",
				Gender:    1,
			},
			wantErr: false,
		},
		{
			name: "success - create user with minimal fields",
			user: domain.User{
				RegNo:     "YY98765432",
				FirstName: "Minimal",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := repo.Create(ctx, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Greater(t, created.Id, 0)
			assert.Equal(t, tt.user.FirstName, created.FirstName)
			assert.Equal(t, tt.user.RegNo, created.RegNo)
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewUserRepository(db)
	ctx := CreateTestContext()

	// Seed a test user
	seededUser := SeedTestUser(t, db)

	tests := []struct {
		name    string
		userID  int
		wantErr bool
	}{
		{
			name:    "success - user found",
			userID:  seededUser.Id,
			wantErr: false,
		},
		{
			name:    "error - user not found",
			userID:  99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByID(ctx, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.userID, user.Id)
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewUserRepository(db)
	ctx := CreateTestContext()

	// Seed a test user
	seededUser := SeedTestUser(t, db)

	tests := []struct {
		name      string
		update    domain.User
		wantErr   bool
		checkFunc func(t *testing.T, updated domain.User)
	}{
		{
			name: "success - update first name",
			update: domain.User{
				Id:        seededUser.Id,
				FirstName: "Updated",
				LastName:  seededUser.LastName,
				RegNo:     seededUser.RegNo,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, updated domain.User) {
				assert.Equal(t, "Updated", updated.FirstName)
			},
		},
		{
			name: "success - update email",
			update: domain.User{
				Id:    seededUser.Id,
				Email: "updated@example.com",
				RegNo: seededUser.RegNo,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, updated domain.User) {
				assert.Equal(t, "updated@example.com", updated.Email)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := repo.Update(ctx, tt.update)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.checkFunc != nil {
				tt.checkFunc(t, updated)
			}
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewUserRepository(db)
	ctx := CreateTestContext()

	// Seed a test user
	seededUser := SeedTestUser(t, db)

	tests := []struct {
		name    string
		userID  int
		wantErr bool
	}{
		{
			name:    "success - user deleted (soft delete)",
			userID:  seededUser.Id,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleted, err := repo.Delete(ctx, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.userID, deleted.Id)

			// Verify soft delete - should not be found by normal query
			_, err = repo.GetByID(ctx, tt.userID)
			assert.Error(t, err) // Should fail because of soft delete
		})
	}
}

func TestUserRepository_List(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewUserRepository(db)
	ctx := CreateTestContext()

	// Clean up and seed test users
	SeedTestUsers(t, db, 15)

	tests := []struct {
		name           string
		query          common.PaginationQuery
		wantMinItems   int
		wantMaxItems   int
		wantTotalGte   int64
		wantErr        bool
	}{
		{
			name: "success - default pagination",
			query: common.PaginationQuery{
				Page: 1,
				Size: 10,
			},
			wantMinItems: 1,
			wantMaxItems: 10,
			wantTotalGte: 15,
			wantErr:      false,
		},
		{
			name: "success - page 2",
			query: common.PaginationQuery{
				Page: 2,
				Size: 10,
			},
			wantMinItems: 1,
			wantMaxItems: 10,
			wantTotalGte: 15,
			wantErr:      false,
		},
		{
			name: "success - small page size",
			query: common.PaginationQuery{
				Page: 1,
				Size: 5,
			},
			wantMinItems: 1,
			wantMaxItems: 5,
			wantTotalGte: 15,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, total, page, size, err := repo.List(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(users), tt.wantMinItems)
			assert.LessOrEqual(t, len(users), tt.wantMaxItems)
			assert.GreaterOrEqual(t, total, tt.wantTotalGte)
			assert.Equal(t, tt.query.Page, page)
			assert.Equal(t, tt.query.Size, size)
		})
	}
}

func TestUserRepository_UserOrgIDs(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewUserRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	user := SeedTestUser(t, db)
	org1 := SeedTestOrganization(t, db)

	// Create org_users relation (manual insert since we don't have that seeder)
	db.Exec("INSERT INTO org_users (user_id, org_id, created_date) VALUES (?, ?, NOW())", user.Id, org1.Id)

	tests := []struct {
		name         string
		userID       int
		wantMinOrgs  int
		wantErr      bool
	}{
		{
			name:        "success - user with organizations",
			userID:      user.Id,
			wantMinOrgs: 1,
			wantErr:     false,
		},
		{
			name:        "success - user without organizations",
			userID:      99999, // Non-existent user
			wantMinOrgs: 0,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgIDs, err := repo.UserOrgIDs(ctx, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(orgIDs), tt.wantMinOrgs)
		})
	}
}
