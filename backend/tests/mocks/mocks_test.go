// Package mocks provides mock implementations for testing
//
// File: mocks_test.go
// Description: Tests for mock implementations
package mocks

import (
	"context"
	"testing"

	"templatev25/internal/domain"

	"git.gerege.mn/backend-packages/common"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_MockMethods(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		mockRepo := NewUserRepository(t)

		expectedUser := domain.User{
			Id:        1,
			FirstName: "Test",
			LastName:  "User",
		}

		mockRepo.On("GetByID", context.Background(), 1).Return(expectedUser, nil)

		user, err := mockRepo.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Id, user.Id)
		assert.Equal(t, expectedUser.FirstName, user.FirstName)
	})

	t.Run("Create", func(t *testing.T) {
		mockRepo := NewUserRepository(t)

		inputUser := domain.User{
			FirstName: "New",
			LastName:  "User",
		}
		expectedUser := domain.User{
			Id:        1,
			FirstName: "New",
			LastName:  "User",
		}

		mockRepo.On("Create", context.Background(), inputUser).Return(expectedUser, nil)

		user, err := mockRepo.Create(context.Background(), inputUser)

		assert.NoError(t, err)
		assert.Equal(t, 1, user.Id)
	})

	t.Run("List", func(t *testing.T) {
		mockRepo := NewUserRepository(t)

		query := common.PaginationQuery{Page: 1, Size: 10}
		expectedUsers := []domain.User{
			{Id: 1, FirstName: "User1"},
			{Id: 2, FirstName: "User2"},
		}

		mockRepo.On("List", context.Background(), query).Return(expectedUsers, int64(2), 1, 10, nil)

		users, total, page, size, err := mockRepo.List(context.Background(), query)

		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 1, page)
		assert.Equal(t, 10, size)
	})

	t.Run("Update", func(t *testing.T) {
		mockRepo := NewUserRepository(t)

		inputUser := domain.User{
			Id:        1,
			FirstName: "Updated",
		}

		mockRepo.On("Update", context.Background(), inputUser).Return(inputUser, nil)

		user, err := mockRepo.Update(context.Background(), inputUser)

		assert.NoError(t, err)
		assert.Equal(t, "Updated", user.FirstName)
	})

	t.Run("Delete", func(t *testing.T) {
		mockRepo := NewUserRepository(t)

		expectedUser := domain.User{Id: 1}

		mockRepo.On("Delete", context.Background(), 1).Return(expectedUser, nil)

		user, err := mockRepo.Delete(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, user.Id)
	})
}

func TestRoleRepository_MockMethods(t *testing.T) {
	t.Run("ByID", func(t *testing.T) {
		mockRepo := NewRoleRepository(t)

		expectedRole := domain.Role{
			ID:   1,
			Name: "Admin",
			Code: "ADMIN",
		}

		mockRepo.On("ByID", context.Background(), 1).Return(expectedRole, nil)

		role, err := mockRepo.ByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedRole.ID, role.ID)
		assert.Equal(t, expectedRole.Name, role.Name)
	})
}

func TestOrganizationRepository_MockMethods(t *testing.T) {
	t.Run("ByID", func(t *testing.T) {
		mockRepo := NewOrganizationRepository(t)

		expectedOrg := domain.Organization{
			Id:   1,
			Name: "Test Org",
		}

		mockRepo.On("ByID", context.Background(), 1).Return(expectedOrg, nil)

		org, err := mockRepo.ByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedOrg.Id, org.Id)
		assert.Equal(t, expectedOrg.Name, org.Name)
	})
}
