//go:build integration

// Package integration contains integration tests
//
// File: news_repo_test.go
// Description: News repository integration tests
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

func TestNewsRepository_Create(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNewsRepository(db)
	ctx := CreateTestContext()

	tests := []struct {
		name    string
		news    domain.News
		wantErr bool
	}{
		{
			name: "success - create new news",
			news: domain.News{
				Title:    "Breaking News",
				Text:     "This is breaking news content",
				ImageUrl: "https://example.com/breaking.jpg",
			},
			wantErr: false,
		},
		{
			name: "success - create news with minimal fields",
			news: domain.News{
				Title: "Minimal News",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.news)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewsRepository_GetByID(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNewsRepository(db)
	ctx := CreateTestContext()

	// Seed a test news
	seededNews := SeedTestNews(t, db)

	tests := []struct {
		name    string
		newsID  int
		wantErr bool
	}{
		{
			name:    "success - news found",
			newsID:  seededNews.Id,
			wantErr: false,
		},
		{
			name:    "error - news not found",
			newsID:  99999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			news, err := repo.GetByID(ctx, tt.newsID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.newsID, news.Id)
		})
	}
}

func TestNewsRepository_Update(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNewsRepository(db)
	ctx := CreateTestContext()

	// Seed a test news
	seededNews := SeedTestNews(t, db)

	tests := []struct {
		name      string
		newsID    int
		update    domain.News
		wantErr   bool
		checkFunc func(t *testing.T)
	}{
		{
			name:   "success - update news title",
			newsID: seededNews.Id,
			update: domain.News{
				Title: "Updated News Title",
			},
			wantErr: false,
			checkFunc: func(t *testing.T) {
				updated, err := repo.GetByID(ctx, seededNews.Id)
				require.NoError(t, err)
				assert.Equal(t, "Updated News Title", updated.Title)
			},
		},
		{
			name:   "success - update news content",
			newsID: seededNews.Id,
			update: domain.News{
				Text: "Updated news content",
			},
			wantErr: false,
			checkFunc: func(t *testing.T) {
				updated, err := repo.GetByID(ctx, seededNews.Id)
				require.NoError(t, err)
				assert.Equal(t, "Updated news content", updated.Text)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Update(ctx, tt.newsID, tt.update)

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

func TestNewsRepository_Delete(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNewsRepository(db)
	ctx := CreateTestContext()

	// Seed a test news
	seededNews := SeedTestNews(t, db)

	tests := []struct {
		name    string
		newsID  int
		wantErr bool
	}{
		{
			name:    "success - news deleted (soft delete)",
			newsID:  seededNews.Id,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.newsID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewsRepository_List(t *testing.T) {
	db := GetTestDBWithTx(t)
	repo := repository.NewNewsRepository(db)
	ctx := CreateTestContext()

	// Seed test news items
	SeedTestNewsItems(t, db, 15)

	tests := []struct {
		name         string
		query        dto.NewsListQuery
		wantMinItems int
		wantMaxItems int
		wantTotalGte int64
		wantErr      bool
	}{
		{
			name: "success - default pagination",
			query: dto.NewsListQuery{
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
			query: dto.NewsListQuery{
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
			query: dto.NewsListQuery{
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
			news, total, page, size, err := repo.List(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(news), tt.wantMinItems)
			assert.LessOrEqual(t, len(news), tt.wantMaxItems)
			assert.GreaterOrEqual(t, total, tt.wantTotalGte)
			assert.Equal(t, tt.query.Page, page)
			assert.Equal(t, tt.query.Size, size)
		})
	}
}
