//go:build integration

// Package integration contains integration tests
//
// File: notification_repo_test.go
// Description: Notification repository integration tests
package integration

import (
	"testing"

	"templatev25/internal/domain"
	"templatev25/internal/repository"

	"git.gerege.mn/backend-packages/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationRepository_CreateGroup(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNotificationRepository(db)
	ctx := CreateTestContext()

	// Seed a test user
	user := SeedTestUser(t, db)

	tests := []struct {
		name    string
		group   domain.NotificationGroup
		wantErr bool
	}{
		{
			name: "success - create notification group",
			group: domain.NotificationGroup{
				UserId:  user.Id,
				Title:   "New Group",
				Content: "Group content",
				Type:    "info",
				Tenant:  "test",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := repo.CreateGroup(ctx, tt.group)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Greater(t, created.Id, 0)
			assert.Equal(t, tt.group.Title, created.Title)
		})
	}
}

func TestNotificationRepository_CreateNotification(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNotificationRepository(db)
	ctx := CreateTestContext()

	// Seed a test user and group
	user := SeedTestUser(t, db)
	group := SeedTestNotificationGroup(t, db, user.Id)

	tests := []struct {
		name         string
		notification domain.Notification
		wantErr      bool
	}{
		{
			name: "success - create notification",
			notification: domain.Notification{
				UserId:  user.Id,
				Title:   "New Notification",
				Content: "Notification content",
				Type:    "info",
				Tenant:  "test",
				GroupId: group.Id,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := repo.CreateNotification(ctx, tt.notification)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Greater(t, created.Id, 0)
			assert.Equal(t, tt.notification.Title, created.Title)
		})
	}
}

func TestNotificationRepository_ListByUser(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNotificationRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	user := SeedTestUser(t, db)
	group := SeedTestNotificationGroup(t, db, user.Id)
	SeedTestNotifications(t, db, user.Id, group.Id, 15)

	tests := []struct {
		name         string
		userID       int
		query        common.PaginationQuery
		wantMinItems int
		wantMaxItems int
		wantTotalGte int64
		wantErr      bool
	}{
		{
			name:   "success - default pagination",
			userID: user.Id,
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
			name:   "success - page 2",
			userID: user.Id,
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
			name:   "success - no notifications for other user",
			userID: 99999,
			query: common.PaginationQuery{
				Page: 1,
				Size: 10,
			},
			wantMinItems: 0,
			wantMaxItems: 0,
			wantTotalGte: 0,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifications, total, page, size, err := repo.ListByUser(ctx, tt.userID, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(notifications), tt.wantMinItems)
			assert.LessOrEqual(t, len(notifications), tt.wantMaxItems)
			assert.GreaterOrEqual(t, total, tt.wantTotalGte)
			assert.Equal(t, tt.query.Page, page)
			assert.Equal(t, tt.query.Size, size)
		})
	}
}

func TestNotificationRepository_MarkGroupRead(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNotificationRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	user := SeedTestUser(t, db)
	group := SeedTestNotificationGroup(t, db, user.Id)
	SeedTestNotifications(t, db, user.Id, group.Id, 5)

	tests := []struct {
		name    string
		userID  int
		groupID int
		wantErr bool
	}{
		{
			name:    "success - mark group as read",
			userID:  user.Id,
			groupID: group.Id,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.MarkGroupRead(ctx, tt.userID, tt.groupID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify all notifications in group are marked as read
			notifications, _, _, _, err := repo.ListByUser(ctx, tt.userID, common.PaginationQuery{Page: 1, Size: 100})
			require.NoError(t, err)
			for _, n := range notifications {
				if n.GroupId == tt.groupID {
					assert.True(t, n.IsRead)
				}
			}
		})
	}
}

func TestNotificationRepository_MarkAllRead(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNotificationRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	user := SeedTestUser(t, db)
	group := SeedTestNotificationGroup(t, db, user.Id)
	SeedTestNotifications(t, db, user.Id, group.Id, 5)

	tests := []struct {
		name    string
		userID  int
		wantErr bool
	}{
		{
			name:    "success - mark all as read",
			userID:  user.Id,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.MarkAllRead(ctx, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify all notifications are marked as read
			notifications, _, _, _, err := repo.ListByUser(ctx, tt.userID, common.PaginationQuery{Page: 1, Size: 100})
			require.NoError(t, err)
			for _, n := range notifications {
				assert.True(t, n.IsRead)
			}
		})
	}
}

func TestNotificationRepository_ListGroups(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNotificationRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	user := SeedTestUser(t, db)
	for i := 0; i < 5; i++ {
		SeedTestNotificationGroup(t, db, user.Id)
	}

	tests := []struct {
		name         string
		query        common.PaginationQuery
		wantMinItems int
		wantErr      bool
	}{
		{
			name: "success - list groups",
			query: common.PaginationQuery{
				Page: 1,
				Size: 10,
			},
			wantMinItems: 5,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups, total, _, _, err := repo.ListGroups(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(groups), tt.wantMinItems)
			assert.GreaterOrEqual(t, total, int64(tt.wantMinItems))
		})
	}
}

func TestNotificationRepository_CreateNotificationsBulk(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNotificationRepository(db)
	ctx := CreateTestContext()

	// Seed test data
	user := SeedTestUser(t, db)
	group := SeedTestNotificationGroup(t, db, user.Id)

	notifications := []domain.Notification{
		{UserId: user.Id, Title: "Bulk 1", GroupId: group.Id, Type: "info", Tenant: "test"},
		{UserId: user.Id, Title: "Bulk 2", GroupId: group.Id, Type: "info", Tenant: "test"},
		{UserId: user.Id, Title: "Bulk 3", GroupId: group.Id, Type: "info", Tenant: "test"},
	}

	err := repo.CreateNotificationsBulk(ctx, notifications)
	require.NoError(t, err)

	// Verify notifications were created
	list, total, _, _, err := repo.ListByUser(ctx, user.Id, common.PaginationQuery{Page: 1, Size: 10})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 3)
	assert.GreaterOrEqual(t, total, int64(3))
}

func TestNotificationRepository_AllUserIDs(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNotificationRepository(db)
	ctx := CreateTestContext()

	// Seed test users
	SeedTestUsers(t, db, 5)

	ids, err := repo.AllUserIDs(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(ids), 5)
}
