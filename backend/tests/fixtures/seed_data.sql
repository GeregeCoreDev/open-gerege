-- Test Seed Data
-- This file contains test data for integration tests

-- Truncate all tables first (in correct order due to FK constraints)
TRUNCATE TABLE user_roles CASCADE;
TRUNCATE TABLE role_permissions CASCADE;
TRUNCATE TABLE permissions CASCADE;
TRUNCATE TABLE roles CASCADE;
TRUNCATE TABLE modules CASCADE;
TRUNCATE TABLE systems CASCADE;
TRUNCATE TABLE users CASCADE;
TRUNCATE TABLE organizations CASCADE;

-- ============================================================
-- ORGANIZATIONS
-- ============================================================
INSERT INTO organizations (id, name, code, is_active, created_date) VALUES
(1, 'Test Organization', 'TEST_ORG', true, NOW()),
(2, 'Demo Organization', 'DEMO_ORG', true, NOW());

-- ============================================================
-- USERS
-- ============================================================
INSERT INTO users (id, reg_no, first_name, last_name, email, phone_no, gender, created_date) VALUES
(1, 'AA00112233', 'Test', 'User', 'test@example.com', '99112233', 1, NOW()),
(2, 'BB00112233', 'Admin', 'User', 'admin@example.com', '99112234', 1, NOW()),
(3, 'CC00112233', 'Normal', 'User', 'normal@example.com', '99112235', 2, NOW());

-- ============================================================
-- SYSTEMS
-- ============================================================
INSERT INTO systems (id, code, key, name, description, is_active, sequence, created_date) VALUES
(1, 'ADMIN', 'admin', 'Admin System', 'Administration system', true, 1, NOW()),
(2, 'USER', 'user', 'User System', 'User management system', true, 2, NOW());

-- ============================================================
-- MODULES
-- ============================================================
INSERT INTO modules (id, system_id, code, key, name, description, is_active, sequence, created_date) VALUES
(1, 1, 'USER_MGMT', 'user-management', 'User Management', 'Manage users', true, 1, NOW()),
(2, 1, 'ROLE_MGMT', 'role-management', 'Role Management', 'Manage roles', true, 2, NOW()),
(3, 2, 'PROFILE', 'profile', 'User Profile', 'User profile module', true, 1, NOW());

-- ============================================================
-- ROLES
-- ============================================================
INSERT INTO roles (id, system_id, code, key, name, description, is_active, created_date) VALUES
(1, 1, 'SUPER_ADMIN', 'super-admin', 'Super Admin', 'Full system access', true, NOW()),
(2, 1, 'ADMIN', 'admin', 'Admin', 'Administrative access', true, NOW()),
(3, 2, 'USER', 'user', 'User', 'Normal user access', true, NOW());

-- ============================================================
-- PERMISSIONS
-- ============================================================
INSERT INTO permissions (id, module_id, code, key, name, description, is_active, created_date) VALUES
(1, 1, 'USER_READ', 'user.read', 'Read Users', 'Can view users', true, NOW()),
(2, 1, 'USER_WRITE', 'user.write', 'Write Users', 'Can create/edit users', true, NOW()),
(3, 1, 'USER_DELETE', 'user.delete', 'Delete Users', 'Can delete users', true, NOW()),
(4, 2, 'ROLE_READ', 'role.read', 'Read Roles', 'Can view roles', true, NOW()),
(5, 2, 'ROLE_WRITE', 'role.write', 'Write Roles', 'Can create/edit roles', true, NOW()),
(6, 3, 'PROFILE_READ', 'profile.read', 'Read Profile', 'Can view profile', true, NOW()),
(7, 3, 'PROFILE_WRITE', 'profile.write', 'Write Profile', 'Can edit profile', true, NOW());

-- ============================================================
-- ROLE_PERMISSIONS
-- ============================================================
INSERT INTO role_permissions (role_id, permission_id, created_date) VALUES
-- Super Admin has all permissions
(1, 1, NOW()), (1, 2, NOW()), (1, 3, NOW()),
(1, 4, NOW()), (1, 5, NOW()),
(1, 6, NOW()), (1, 7, NOW()),
-- Admin has read/write but not delete
(2, 1, NOW()), (2, 2, NOW()),
(2, 4, NOW()), (2, 5, NOW()),
(2, 6, NOW()), (2, 7, NOW()),
-- User has only profile permissions
(3, 6, NOW()), (3, 7, NOW());

-- ============================================================
-- USER_ROLES
-- ============================================================
INSERT INTO user_roles (user_id, role_id, created_date) VALUES
(1, 1, NOW()),  -- Test User is Super Admin
(2, 2, NOW()),  -- Admin User is Admin
(3, 3, NOW());  -- Normal User is User
