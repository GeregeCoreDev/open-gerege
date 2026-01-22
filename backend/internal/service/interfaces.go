// Package service provides business logic implementations
//
// File: interfaces.go
// Description: Service layer interfaces for dependency injection and testing
//
// These interfaces define the contracts for all service implementations.
// Handlers should depend on these interfaces, not concrete implementations.
// This enables:
//   - Easy mocking for unit tests
//   - Loose coupling between layers
//   - Flexibility to swap implementations
package service

import (
	"context"

	"templatev25/internal/auth"
	"templatev25/internal/domain"
	"templatev25/internal/http/dto"

	"git.gerege.mn/backend-packages/common"
)

// ============================================================
// USER SERVICE
// ============================================================

// UserServiceInterface defines user management operations
type UserServiceInterface interface {
	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id int) (domain.User, error)

	// List retrieves paginated users
	List(ctx context.Context, p common.PaginationQuery) ([]domain.User, int64, int, int, error)

	// Create creates a new user or returns existing if already exists
	Create(ctx context.Context, req dto.UserCreateDto) (domain.User, error)

	// Update updates an existing user
	Update(ctx context.Context, req dto.UserUpdateDto) (domain.User, error)

	// Delete soft-deletes a user
	Delete(ctx context.Context, id int) (domain.User, error)

	// Organizations retrieves user's organizations
	Organizations(ctx context.Context, userID, currentOrgID int, fields []string) (orgID int, org *domain.Organization, items []domain.Organization, err error)
}

// ============================================================
// ROLE SERVICE
// ============================================================

// RoleServiceInterface defines role management operations
type RoleServiceInterface interface {
	// SetCacheInvalidator sets the permission cache invalidator
	SetCacheInvalidator(cache auth.CacheInvalidator)

	// ListFilteredPaged retrieves paginated roles with filtering
	ListFilteredPaged(ctx context.Context, p dto.RoleListQuery) ([]domain.Role, int64, int, int, error)

	// Create creates a new role
	Create(ctx context.Context, req dto.RoleCreateDto) error

	// Update updates an existing role
	Update(ctx context.Context, id int, req dto.RoleUpdateDto) error

	// Delete soft-deletes a role (only if inactive)
	Delete(ctx context.Context, id int) error

	// GetPermissions retrieves permissions for a role
	GetPermissions(ctx context.Context, q dto.RolePermissionsQuery) ([]domain.Permission, error)

	// SetPermissions replaces all permissions for a role
	SetPermissions(ctx context.Context, req dto.RolePermissionsUpdateDto) error
}

// ============================================================
// PERMISSION SERVICE
// ============================================================

// PermissionServiceInterface defines permission management operations
type PermissionServiceInterface interface {
	// SetCacheInvalidator sets the permission cache invalidator
	SetCacheInvalidator(cache auth.CacheInvalidator)

	// ListFilteredPaged retrieves paginated permissions with filtering
	ListFilteredPaged(ctx context.Context, q dto.PermissionQuery) ([]domain.Permission, int64, int, int, error)

	// ByID retrieves a permission by ID
	ByID(ctx context.Context, id int) (domain.Permission, error)

	// ByCode retrieves a permission by code
	ByCode(ctx context.Context, code string) (domain.Permission, error)

	// Create creates new permissions in batch
	Create(ctx context.Context, req dto.PermissionCreateDto) error

	// Update updates an existing permission
	Update(ctx context.Context, id int, req dto.PermissionUpdateDto) error

	// Delete deletes a permission
	Delete(ctx context.Context, id int) error
}

// ============================================================
// ORGANIZATION SERVICE
// ============================================================

// OrganizationServiceInterface defines organization management operations
type OrganizationServiceInterface interface {
	// List retrieves paginated organizations
	List(ctx context.Context, p common.PaginationQuery) ([]domain.Organization, int64, int, int, error)

	// Create creates a new organization
	Create(ctx context.Context, req dto.OrganizationDto) (domain.Organization, error)

	// Update updates an existing organization
	Update(ctx context.Context, id int, req dto.OrganizationUpdateDto) (domain.Organization, error)

	// Delete deletes an organization
	Delete(ctx context.Context, id int) error

	// ByID retrieves an organization by ID
	ByID(ctx context.Context, id int) (domain.Organization, error)

	// Tree retrieves organization hierarchy tree starting from rootID
	Tree(ctx context.Context, rootID int) ([]domain.Organization, error)
}

// ============================================================
// NEWS SERVICE
// ============================================================

// NewsServiceInterface defines news management operations
type NewsServiceInterface interface {
	// List retrieves paginated news
	List(ctx context.Context, q dto.NewsListQuery) ([]domain.News, int64, int, int, error)

	// GetByID retrieves a news item by ID
	GetByID(ctx context.Context, id int) (domain.News, error)

	// Create creates a new news item
	Create(ctx context.Context, req dto.NewsDto) error

	// Update updates an existing news item
	Update(ctx context.Context, id int, req dto.NewsDto) error

	// Delete deletes a news item
	Delete(ctx context.Context, id int) error
}

// ============================================================
// SYSTEM SERVICE
// ============================================================

// SystemServiceInterface defines system management operations
type SystemServiceInterface interface {
	// List retrieves paginated systems
	List(ctx context.Context, q dto.SystemListQuery) ([]domain.System, int64, int, int, error)

	// ByID retrieves a system by ID
	ByID(ctx context.Context, id int) (domain.System, error)

	// Create creates a new system
	Create(ctx context.Context, req dto.SystemCreateDto) error

	// Update updates an existing system
	Update(ctx context.Context, id int, req dto.SystemUpdateDto) error

	// Delete deletes a system
	Delete(ctx context.Context, id int) error
}

// ============================================================
// MODULE SERVICE
// ============================================================

// ModuleServiceInterface defines module management operations
type ModuleServiceInterface interface {
	// List retrieves paginated modules
	List(ctx context.Context, q dto.ModuleListQuery) ([]domain.Module, int64, int, int, error)

	// ByID retrieves a module by ID
	ByID(ctx context.Context, id int) (domain.Module, error)

	// Create creates a new module
	Create(ctx context.Context, req dto.ModuleCreateDto) error

	// Update updates an existing module
	Update(ctx context.Context, id int, req dto.ModuleUpdateDto) error

	// Delete deletes a module
	Delete(ctx context.Context, id int) error
}

// NOTE: Additional service interfaces (NotificationService, UserRoleService, VerifyService, etc.)
// can be added here as the corresponding DTOs are defined.

// ============================================================
// Compile-time interface implementation checks
// ============================================================

var (
	_ UserServiceInterface         = (*UserService)(nil)
	_ RoleServiceInterface         = (*RoleService)(nil)
	_ PermissionServiceInterface   = (*PermissionService)(nil)
	_ OrganizationServiceInterface = (*OrganizationService)(nil)
	_ NewsServiceInterface         = (*NewsService)(nil)
)
