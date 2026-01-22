// Package app provides dependency injection
//
// File: dependency_test.go
// Description: Unit tests for app package
package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDependencies_Structure tests that Dependencies struct has expected fields
func TestDependencies_Structure(t *testing.T) {
	// Document the expected structure of Dependencies
	deps := &Dependencies{}

	// Core dependencies
	assert.Nil(t, deps.DB)
	assert.Nil(t, deps.Log)
	assert.Nil(t, deps.Cfg)
	assert.Nil(t, deps.AuthCache)
	assert.Nil(t, deps.SSO)
	assert.Nil(t, deps.PermCache)

	// Layer containers
	assert.Nil(t, deps.Repo)
	assert.Nil(t, deps.Service)
}

// TestRepoContainer_Structure tests that RepoContainer has expected repositories
func TestRepoContainer_Structure(t *testing.T) {
	repo := &RepoContainer{}

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

// TestServiceContainer_Structure tests that ServiceContainer has expected services
func TestServiceContainer_Structure(t *testing.T) {
	svc := &ServiceContainer{}

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
