-- +goose Up
-- Migration: 008_performance_indexes.sql
-- Description: Add missing indexes for query optimization
-- Author: Performance Improvement
-- Date: 2026-01-10

SET search_path TO template_backend;

-- user_roles table - queries often filter by role_id alone or user_id alone
-- Composite unique already exists, but single-column indexes help with partial matches
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);

-- role_permissions table - permission lookups
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);

-- organization_users table - queries filter by user_id alone
CREATE INDEX IF NOT EXISTS idx_organization_users_user_id ON organization_users(user_id);
CREATE INDEX IF NOT EXISTS idx_organization_users_org_id ON organization_users(org_id);

-- menus table - frequently filtered by permission_id and is_active
CREATE INDEX IF NOT EXISTS idx_menus_permission_id ON menus(permission_id);
CREATE INDEX IF NOT EXISTS idx_menus_is_active ON menus(is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_menus_parent_id ON menus(parent_id);
CREATE INDEX IF NOT EXISTS idx_menus_sequence ON menus(sequence);

-- permissions table - code lookups
CREATE INDEX IF NOT EXISTS idx_permissions_code ON permissions(code);
CREATE INDEX IF NOT EXISTS idx_permissions_is_active ON permissions(is_active) WHERE is_active = true;

-- roles table - code lookups
CREATE INDEX IF NOT EXISTS idx_roles_code ON roles(code);
CREATE INDEX IF NOT EXISTS idx_roles_system_id ON roles(system_id);
CREATE INDEX IF NOT EXISTS idx_roles_is_active ON roles(is_active) WHERE is_active = true;

-- Composite index for menu queries with JOINs
CREATE INDEX IF NOT EXISTS idx_menus_active_deleted ON menus(is_active, deleted_date) WHERE is_active = true AND deleted_date IS NULL;

RESET search_path;

-- +goose Down
SET search_path TO template_backend;

DROP INDEX IF EXISTS idx_menus_active_deleted;
DROP INDEX IF EXISTS idx_roles_is_active;
DROP INDEX IF EXISTS idx_roles_system_id;
DROP INDEX IF EXISTS idx_roles_code;
DROP INDEX IF EXISTS idx_permissions_is_active;
DROP INDEX IF EXISTS idx_permissions_code;
DROP INDEX IF EXISTS idx_menus_sequence;
DROP INDEX IF EXISTS idx_menus_parent_id;
DROP INDEX IF EXISTS idx_menus_is_active;
DROP INDEX IF EXISTS idx_menus_permission_id;
DROP INDEX IF EXISTS idx_organization_users_org_id;
DROP INDEX IF EXISTS idx_organization_users_user_id;
DROP INDEX IF EXISTS idx_role_permissions_role_id;
DROP INDEX IF EXISTS idx_role_permissions_permission_id;
DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP INDEX IF EXISTS idx_user_roles_role_id;

RESET search_path;
