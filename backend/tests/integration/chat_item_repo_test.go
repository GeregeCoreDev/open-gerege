//go:build integration

// Package integration contains integration tests
//
// File: chat_item_repo_test.go
// Description: ChatItem repository integration tests
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

func TestChatItemRepository_Create(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewChatItemRepository(db)
	ctx := CreateTestContext()

	tests := []struct {
		name    string
		item    domain.ChatItem
		wantErr bool
	}{
		{
			name: "success - create chat item",
			item: domain.ChatItem{
				Key:    "hello",
				Answer: "Hello, how can I help you?",
			},
			wantErr: false,
		},
		{
			name: "success - create another chat item",
			item: domain.ChatItem{
				Key:    "goodbye",
				Answer: "Goodbye! Have a nice day.",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.item)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestChatItemRepository_ByID(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewChatItemRepository(db)
	ctx := CreateTestContext()

	// Seed a test chat item
	seededItem := SeedTestChatItem(t, db)

	tests := []struct {
		name    string
		itemID  int
		wantErr bool
	}{
		{
			name:    "success - item found",
			itemID:  seededItem.ID,
			wantErr: false,
		},
		{
			name:    "error - item not found",
			itemID:  99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := repo.ByID(ctx, tt.itemID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.itemID, item.ID)
		})
	}
}

func TestChatItemRepository_FindByKey(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewChatItemRepository(db)
	ctx := CreateTestContext()

	// Seed a test chat item
	seededItem := SeedTestChatItem(t, db)

	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "success - item found by key",
			key:     seededItem.Key,
			wantErr: false,
		},
		{
			name:    "error - item not found by key",
			key:     "nonexistent-key",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := repo.FindByKey(ctx, tt.key)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.key, item.Key)
		})
	}
}

func TestChatItemRepository_Update(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewChatItemRepository(db)
	ctx := CreateTestContext()

	// Seed a test chat item
	seededItem := SeedTestChatItem(t, db)

	tests := []struct {
		name      string
		itemID    int
		update    domain.ChatItem
		wantErr   bool
		checkFunc func(t *testing.T)
	}{
		{
			name:   "success - update answer",
			itemID: seededItem.ID,
			update: domain.ChatItem{
				Key:    seededItem.Key,
				Answer: "updated answer",
			},
			wantErr: false,
			checkFunc: func(t *testing.T) {
				updated, err := repo.ByID(ctx, seededItem.ID)
				require.NoError(t, err)
				assert.Equal(t, "updated answer", updated.Answer)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.itemID, tt.update)

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

func TestChatItemRepository_Delete(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewChatItemRepository(db)
	ctx := CreateTestContext()

	// Seed a test chat item
	seededItem := SeedTestChatItem(t, db)

	tests := []struct {
		name    string
		itemID  int
		wantErr bool
	}{
		{
			name:    "success - item deleted (soft delete)",
			itemID:  seededItem.ID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.itemID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestChatItemRepository_List(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewChatItemRepository(db)
	ctx := CreateTestContext()

	// Seed test chat items
	SeedTestChatItems(t, db, 15)

	tests := []struct {
		name         string
		query        dto.ChatItemQuery
		wantMinItems int
		wantMaxItems int
		wantTotalGte int64
		wantErr      bool
	}{
		{
			name: "success - default pagination",
			query: dto.ChatItemQuery{
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
			query: dto.ChatItemQuery{
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
		{
			name: "success - small page size",
			query: dto.ChatItemQuery{
				PaginationQuery: common.PaginationQuery{
					Page: 1,
					Size: 5,
				},
			},
			wantMinItems: 1,
			wantMaxItems: 5,
			wantTotalGte: 15,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, total, page, size, err := repo.List(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(items), tt.wantMinItems)
			assert.LessOrEqual(t, len(items), tt.wantMaxItems)
			assert.GreaterOrEqual(t, total, tt.wantTotalGte)
			assert.Equal(t, tt.query.Page, page)
			assert.Equal(t, tt.query.Size, size)
		})
	}
}
