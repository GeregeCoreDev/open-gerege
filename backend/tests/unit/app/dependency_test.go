package app_test

import (
	"testing"

	"templatev25/internal/app"

	"github.com/stretchr/testify/assert"
)

func TestDependencies_StructFields(t *testing.T) {
	// Test that Dependencies struct has all required fields
	deps := &app.Dependencies{}

	// These should all be nil initially
	assert.Nil(t, deps.DB)
	assert.Nil(t, deps.Log)
	assert.Nil(t, deps.Cfg)
	assert.Nil(t, deps.AuthCache)
	assert.Nil(t, deps.SSO)
	assert.Nil(t, deps.PermCache)
	assert.Nil(t, deps.Repo)
	assert.Nil(t, deps.Service)
}

func TestRepoContainer_StructFields(t *testing.T) {
	// Test that RepoContainer struct has all required fields
	repo := &app.RepoContainer{}

	// User & Auth
	assert.Nil(t, repo.User)
	assert.Nil(t, repo.UserRole)

	// System & Module
	assert.Nil(t, repo.System)
	assert.Nil(t, repo.Module)
	assert.Nil(t, repo.Menu)

	// Permission & Role
	assert.Nil(t, repo.Permission)
	assert.Nil(t, repo.Action)
	assert.Nil(t, repo.Role)

	// Organization
	assert.Nil(t, repo.Organization)
	assert.Nil(t, repo.OrganizationType)
	assert.Nil(t, repo.OrgUser)

	// Terminal & Platform
	assert.Nil(t, repo.Terminal)
	assert.Nil(t, repo.AppServiceIcon)
	assert.Nil(t, repo.AppServiceIconGroup)

	// Content
	assert.Nil(t, repo.PublicFile)
	assert.Nil(t, repo.Notification)
	assert.Nil(t, repo.News)
	assert.Nil(t, repo.ChatItem)
	assert.Nil(t, repo.APILog)
}

func TestServiceContainer_StructFields(t *testing.T) {
	// Test that ServiceContainer struct has all required fields
	svc := &app.ServiceContainer{}

	// User & Auth
	assert.Nil(t, svc.User)
	assert.Nil(t, svc.UserRole)

	// System & Module
	assert.Nil(t, svc.System)
	assert.Nil(t, svc.Module)
	assert.Nil(t, svc.Menu)

	// Permission & Role
	assert.Nil(t, svc.Permission)
	assert.Nil(t, svc.Action)
	assert.Nil(t, svc.Role)

	// Organization
	assert.Nil(t, svc.Organization)
	assert.Nil(t, svc.OrganizationType)
	assert.Nil(t, svc.OrgUser)

	// Terminal & Platform
	assert.Nil(t, svc.Terminal)
	assert.Nil(t, svc.AppServiceIcon)
	assert.Nil(t, svc.AppServiceGroup)

	// Content
	assert.Nil(t, svc.PublicFile)
	assert.Nil(t, svc.Notification)
	assert.Nil(t, svc.News)
	assert.Nil(t, svc.ChatItem)
	assert.Nil(t, svc.APILog)

	// External Integrations
	assert.Nil(t, svc.Verify)
	assert.Nil(t, svc.Meet)
	assert.Nil(t, svc.Tpay)
}
