-- +goose Up
-- +goose StatementBegin
SET search_path TO template_backend;

-- ============================================
-- STANDARDIZE CODE FIELDS
-- All codes: lowercase with dot separator
-- Systems: system
-- Actions: action
-- Modules: system.module
-- Roles: system.role
-- Permissions: system.module.action
-- ============================================

-- ============================================
-- UPDATE SYSTEMS: lowercase
-- ============================================
UPDATE systems SET code = 'admin' WHERE code = 'ADMIN';
UPDATE systems SET code = 'grgb' WHERE code = 'GRGB';
UPDATE systems SET code = 'grga' WHERE code = 'GRGA';
UPDATE systems SET code = 'tpay' WHERE code = 'TPAY';

-- ============================================
-- UPDATE ACTIONS: lowercase
-- ============================================
UPDATE actions SET code = 'create' WHERE code = 'CREATE';
UPDATE actions SET code = 'read' WHERE code = 'READ';
UPDATE actions SET code = 'update' WHERE code = 'UPDATE';
UPDATE actions SET code = 'delete' WHERE code = 'DELETE';

-- ============================================
-- UPDATE MODULES: system.module format
-- ============================================
-- Admin system modules
UPDATE modules SET code = 'admin.system' WHERE code = 'ADMIN-SYSTEM';
UPDATE modules SET code = 'admin.module' WHERE code = 'ADMIN-MODULE';
UPDATE modules SET code = 'admin.role' WHERE code = 'ADMIN-ROLE';
UPDATE modules SET code = 'admin.permission' WHERE code = 'ADMIN-MODULE-PERMISSION';
UPDATE modules SET code = 'admin.org' WHERE code = 'ADMIN-ORG';
UPDATE modules SET code = 'admin.user' WHERE code = 'ADMIN-USER';
UPDATE modules SET code = 'admin.dashboard' WHERE code = 'ADMIN-DASHBOARD';
UPDATE modules SET code = 'admin.orgtype' WHERE code = 'ADMIN-ORGTYPE';
UPDATE modules SET code = 'admin.aichat' WHERE code = 'ADMIN-AI-CHAT';

-- Gerege App modules
UPDATE modules SET code = 'grga.icon' WHERE code = 'GRGA-ICON';
UPDATE modules SET code = 'grga.dashboard' WHERE code = 'GRGA-DASHBOARD';
UPDATE modules SET code = 'grga.role' WHERE code = 'GRGA-ROLE';

-- Gerege Business modules
UPDATE modules SET code = 'grgb.icon' WHERE code = 'GRGB-ICON';
UPDATE modules SET code = 'grgb.dashboard' WHERE code = 'GRGB-DASHBOARD';
UPDATE modules SET code = 'grgb.agents' WHERE code = 'GRGB-AGENTS';
UPDATE modules SET code = 'grgb.role' WHERE code = 'GRGB-ROLE';

-- TPay modules
UPDATE modules SET code = 'tpay.dashboard' WHERE code = 'TPAY-DASHBOARD';
UPDATE modules SET code = 'tpay.wallet' WHERE code = 'TPAY-WALLET';
UPDATE modules SET code = 'tpay.trx' WHERE code = 'TPAY-TRX';
UPDATE modules SET code = 'tpay.balance' WHERE code = 'TPAY-BALANCE';
UPDATE modules SET code = 'tpay.role' WHERE code = 'TPAY-ROLE';
UPDATE modules SET code = 'tpay.users' WHERE code = 'TPAY-USERS';
UPDATE modules SET code = 'tpay.products' WHERE code = 'TPAY-PRODUCTS';

-- ============================================
-- UPDATE ROLES: system.role format
-- ============================================
UPDATE roles SET code = 'admin.superadmin' WHERE code = 'ADMIN-SUPERADMIN';
UPDATE roles SET code = 'admin.manager' WHERE code = 'ADMIN-MANAGER';
UPDATE roles SET code = 'tpay.sysadmin' WHERE code = 'TPAY-SYSADMIN';
UPDATE roles SET code = 'tpay.manager' WHERE code = 'TPAY-MANAGER';
UPDATE roles SET code = 'grga.sysadmin' WHERE code = 'GRGA-SYSADMIN';
UPDATE roles SET code = 'grga.manager' WHERE code = 'GRGA-MANAGER';
UPDATE roles SET code = 'grga.user' WHERE code = 'GRGA-USER';
UPDATE roles SET code = 'grgb.sysadmin' WHERE code = 'GRGB-SYSADMIN';
UPDATE roles SET code = 'grgb.user' WHERE code = 'GRGB-USER';

-- ============================================
-- UPDATE PERMISSIONS: system.module.action format
-- ============================================
-- Admin Dashboard
UPDATE permissions SET code = 'admin.dashboard.read' WHERE code = 'ADMIN-DASHBOARD-VIEW';

-- Admin OrgType
UPDATE permissions SET code = 'admin.orgtype.create' WHERE code = 'ADMIN-ORG-TYPE-CREATE';
UPDATE permissions SET code = 'admin.orgtype.read' WHERE code = 'ADMIN-ORG-TYPE-READ';
UPDATE permissions SET code = 'admin.orgtype.update' WHERE code = 'ADMIN-ORG-TYPE-UPDATE';
UPDATE permissions SET code = 'admin.orgtype.delete' WHERE code = 'ADMIN-ORG-TYPE-DELETE';

-- Admin Org
UPDATE permissions SET code = 'admin.org.create' WHERE code = 'ADMIN-ORG-CREATE';
UPDATE permissions SET code = 'admin.org.read' WHERE code = 'ADMIN-ORG-READ';
UPDATE permissions SET code = 'admin.org.update' WHERE code = 'ADMIN-ORG-UPDATE';
UPDATE permissions SET code = 'admin.org.delete' WHERE code = 'ADMIN-ORG-DELETE';

-- Admin System
UPDATE permissions SET code = 'admin.system.create' WHERE code = 'ADMIN-SYSTEM-CREATE';
UPDATE permissions SET code = 'admin.system.read' WHERE code = 'ADMIN-SYSTEM-READ';
UPDATE permissions SET code = 'admin.system.update' WHERE code = 'ADMIN-SYSTEM-UPDATE';
UPDATE permissions SET code = 'admin.system.delete' WHERE code = 'ADMIN-SYSTEM-DELETE';

-- Admin Module
UPDATE permissions SET code = 'admin.module.create' WHERE code = 'ADMIN-MODULE-CREATE';
UPDATE permissions SET code = 'admin.module.read' WHERE code = 'ADMIN-MODULE-READ';
UPDATE permissions SET code = 'admin.module.update' WHERE code = 'ADMIN-MODULE-UPDATE';
UPDATE permissions SET code = 'admin.module.delete' WHERE code = 'ADMIN-MODULE-DELETE';

-- Admin Permission
UPDATE permissions SET code = 'admin.permission.create' WHERE code = 'ADMIN-MODULE-PERMISSION-CREATE';
UPDATE permissions SET code = 'admin.permission.read' WHERE code = 'ADMIN-MODULE-PERMISSION-READ';
UPDATE permissions SET code = 'admin.permission.update' WHERE code = 'ADMIN-MODULE-PERMISSION-UPDATE';
UPDATE permissions SET code = 'admin.permission.delete' WHERE code = 'ADMIN-MODULE-PERMISSION-DELETE';

-- Admin Role
UPDATE permissions SET code = 'admin.role.create' WHERE code = 'ADMIN-ROLE-CREATE';
UPDATE permissions SET code = 'admin.role.read' WHERE code = 'ADMIN-ROLE-READ';
UPDATE permissions SET code = 'admin.role.update' WHERE code = 'ADMIN-ROLE-UPDATE';
UPDATE permissions SET code = 'admin.role.delete' WHERE code = 'ADMIN-ROLE-DELETE';

-- Admin User
UPDATE permissions SET code = 'admin.user.create' WHERE code = 'ADMIN-USER-CREATE';
UPDATE permissions SET code = 'admin.user.read' WHERE code = 'ADMIN-USER-READ';
UPDATE permissions SET code = 'admin.user.update' WHERE code = 'ADMIN-USER-UPDATE';
UPDATE permissions SET code = 'admin.user.delete' WHERE code = 'ADMIN-USER-DELETE';

-- Admin AI Chat
UPDATE permissions SET code = 'admin.aichat.create' WHERE code = 'ADMIN-AI-CHAT-CREATE';
UPDATE permissions SET code = 'admin.aichat.read' WHERE code = 'ADMIN-AI-CHAT-READ';
UPDATE permissions SET code = 'admin.aichat.update' WHERE code = 'ADMIN-AI-CHAT-UPDATE';
UPDATE permissions SET code = 'admin.aichat.delete' WHERE code = 'ADMIN-AI-CHAT-DELETE';

-- Gerege App Dashboard
UPDATE permissions SET code = 'grga.dashboard.read' WHERE code = 'GRGA-DASHBOARD-READ';

-- Gerege App Icon
UPDATE permissions SET code = 'grga.icon.create' WHERE code = 'GRGA-ICON-CREATE';
UPDATE permissions SET code = 'grga.icon.read' WHERE code = 'GRGA-ICON-READ';
UPDATE permissions SET code = 'grga.icon.update' WHERE code = 'GRGA-ICON-UPDATE';
UPDATE permissions SET code = 'grga.icon.delete' WHERE code = 'GRGA-ICON-DELETE';

-- Gerege App Role
UPDATE permissions SET code = 'grga.role.create' WHERE code = 'GRGA-ROLE-CREATE';
UPDATE permissions SET code = 'grga.role.read' WHERE code = 'GRGA-ROLE-READ';
UPDATE permissions SET code = 'grga.role.update' WHERE code = 'GRGA-ROLE-UPDATE';
UPDATE permissions SET code = 'grga.role.delete' WHERE code = 'GRGA-ROLE-DELETE';

-- Gerege Business Dashboard
UPDATE permissions SET code = 'grgb.dashboard.read' WHERE code = 'GRGB-DASHBOARD-READ';

-- Gerege Business Icon
UPDATE permissions SET code = 'grgb.icon.create' WHERE code = 'GRGB-ICON-CREATE';
UPDATE permissions SET code = 'grgb.icon.read' WHERE code = 'GRGB-ICON-READ';
UPDATE permissions SET code = 'grgb.icon.update' WHERE code = 'GRGB-ICON-UPDATE';
UPDATE permissions SET code = 'grgb.icon.delete' WHERE code = 'GRGB-ICON-DELETE';

-- Gerege Business Role
UPDATE permissions SET code = 'grgb.role.create' WHERE code = 'GRGB-ROLE-CREATE';
UPDATE permissions SET code = 'grgb.role.read' WHERE code = 'GRGB-ROLE-READ';
UPDATE permissions SET code = 'grgb.role.update' WHERE code = 'GRGB-ROLE-UPDATE';
UPDATE permissions SET code = 'grgb.role.delete' WHERE code = 'GRGB-ROLE-DELETE';

-- Gerege Business Agents
UPDATE permissions SET code = 'grgb.agents.create' WHERE code = 'GRGB-AGENTS-CREATE';
UPDATE permissions SET code = 'grgb.agents.read' WHERE code = 'GRGB-AGENTS-READ';
UPDATE permissions SET code = 'grgb.agents.update' WHERE code = 'GRGB-AGENTS-UPDATE';
UPDATE permissions SET code = 'grgb.agents.delete' WHERE code = 'GRGB-AGENTS-DELETE';

-- TPay Dashboard
UPDATE permissions SET code = 'tpay.dashboard.read' WHERE code = 'TPAY-DASHBOARD-READ';

-- TPay Balance
UPDATE permissions SET code = 'tpay.balance.read' WHERE code = 'TPAY-BALANCE-READ';

-- TPay Transaction
UPDATE permissions SET code = 'tpay.trx.read' WHERE code = 'TPAY-TRX-READ';

-- TPay Wallet
UPDATE permissions SET code = 'tpay.wallet.read' WHERE code = 'TPAY-WALLET-READ';

-- TPay Role
UPDATE permissions SET code = 'tpay.role.create' WHERE code = 'TPAY-ROLE-CREATE';
UPDATE permissions SET code = 'tpay.role.read' WHERE code = 'TPAY-ROLE-READ';
UPDATE permissions SET code = 'tpay.role.update' WHERE code = 'TPAY-ROLE-UPDATE';
UPDATE permissions SET code = 'tpay.role.delete' WHERE code = 'TPAY-ROLE-DELETE';

-- TPay Users
UPDATE permissions SET code = 'tpay.users.create' WHERE code = 'TPAY-USER-READ-CREATE';
UPDATE permissions SET code = 'tpay.users.read' WHERE code = 'TPAY-USER-READ-READ';
UPDATE permissions SET code = 'tpay.users.update' WHERE code = 'TPAY-USER-READ-UPDATE';
UPDATE permissions SET code = 'tpay.users.delete' WHERE code = 'TPAY-USER-READ-DELETE';

-- TPay Products
UPDATE permissions SET code = 'tpay.products.create' WHERE code = 'TPAY-PRODUCTS-CREATE';
UPDATE permissions SET code = 'tpay.products.read' WHERE code = 'TPAY-PRODUCTS-READ';
UPDATE permissions SET code = 'tpay.products.update' WHERE code = 'TPAY-PRODUCTS-UPDATE';
UPDATE permissions SET code = 'tpay.products.delete' WHERE code = 'TPAY-PRODUCTS-DELETE';

-- ============================================
-- UPDATE MENUS: lowercase codes
-- ============================================
UPDATE menus SET code = 'admin' WHERE code = 'ADMIN';
UPDATE menus SET code = 'admin.dashboard' WHERE code = 'ADMIN-DASHBOARD';
UPDATE menus SET code = 'admin.system' WHERE code = 'ADMIN-SYSTEM';
UPDATE menus SET code = 'admin.module' WHERE code = 'ADMIN-MODULE';
UPDATE menus SET code = 'admin.role' WHERE code = 'ADMIN-ROLE';
UPDATE menus SET code = 'admin.permission' WHERE code = 'ADMIN-MODULE-PERMISSION';
UPDATE menus SET code = 'admin.org' WHERE code = 'ADMIN-ORG';
UPDATE menus SET code = 'admin.orgtype' WHERE code = 'ADMIN-ORGTYPE';
UPDATE menus SET code = 'admin.user' WHERE code = 'ADMIN-USER';
UPDATE menus SET code = 'admin.aichat' WHERE code = 'ADMIN-AI-CHAT';

UPDATE menus SET code = 'tpay' WHERE code = 'TPAY';
UPDATE menus SET code = 'tpay.dashboard' WHERE code = 'TPAY-DASHBOARD';
UPDATE menus SET code = 'tpay.wallet' WHERE code = 'TPAY-WALLET';
UPDATE menus SET code = 'tpay.trx' WHERE code = 'TPAY-TRX';
UPDATE menus SET code = 'tpay.balance' WHERE code = 'TPAY-BALANCE';
UPDATE menus SET code = 'tpay.users' WHERE code = 'TPAY-USERS';
UPDATE menus SET code = 'tpay.products' WHERE code = 'TPAY-PRODUCTS';
UPDATE menus SET code = 'tpay.role' WHERE code = 'TPAY-ROLE';

UPDATE menus SET code = 'grga' WHERE code = 'GRGA';
UPDATE menus SET code = 'grga.dashboard' WHERE code = 'GRGA-DASHBOARD';
UPDATE menus SET code = 'grga.icon' WHERE code = 'GRGA-ICON';
UPDATE menus SET code = 'grga.role' WHERE code = 'GRGA-ROLE';

UPDATE menus SET code = 'grgb' WHERE code = 'GRGB';
UPDATE menus SET code = 'grgb.dashboard' WHERE code = 'GRGB-DASHBOARD';
UPDATE menus SET code = 'grgb.icon' WHERE code = 'GRGB-ICON';
UPDATE menus SET code = 'grgb.agents' WHERE code = 'GRGB-AGENTS';
UPDATE menus SET code = 'grgb.role' WHERE code = 'GRGB-ROLE';

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO template_backend;

-- ============================================
-- REVERT SYSTEMS
-- ============================================
UPDATE systems SET code = 'ADMIN' WHERE code = 'admin';
UPDATE systems SET code = 'GRGB' WHERE code = 'grgb';
UPDATE systems SET code = 'GRGA' WHERE code = 'grga';
UPDATE systems SET code = 'TPAY' WHERE code = 'tpay';

-- ============================================
-- REVERT ACTIONS
-- ============================================
UPDATE actions SET code = 'CREATE' WHERE code = 'create';
UPDATE actions SET code = 'READ' WHERE code = 'read';
UPDATE actions SET code = 'UPDATE' WHERE code = 'update';
UPDATE actions SET code = 'DELETE' WHERE code = 'delete';

-- ============================================
-- REVERT MODULES
-- ============================================
UPDATE modules SET code = 'ADMIN-SYSTEM' WHERE code = 'admin.system';
UPDATE modules SET code = 'ADMIN-MODULE' WHERE code = 'admin.module';
UPDATE modules SET code = 'ADMIN-ROLE' WHERE code = 'admin.role';
UPDATE modules SET code = 'ADMIN-MODULE-PERMISSION' WHERE code = 'admin.permission';
UPDATE modules SET code = 'ADMIN-ORG' WHERE code = 'admin.org';
UPDATE modules SET code = 'ADMIN-USER' WHERE code = 'admin.user';
UPDATE modules SET code = 'ADMIN-DASHBOARD' WHERE code = 'admin.dashboard';
UPDATE modules SET code = 'ADMIN-ORGTYPE' WHERE code = 'admin.orgtype';
UPDATE modules SET code = 'ADMIN-AI-CHAT' WHERE code = 'admin.aichat';
UPDATE modules SET code = 'GRGA-ICON' WHERE code = 'grga.icon';
UPDATE modules SET code = 'GRGA-DASHBOARD' WHERE code = 'grga.dashboard';
UPDATE modules SET code = 'GRGA-ROLE' WHERE code = 'grga.role';
UPDATE modules SET code = 'GRGB-ICON' WHERE code = 'grgb.icon';
UPDATE modules SET code = 'GRGB-DASHBOARD' WHERE code = 'grgb.dashboard';
UPDATE modules SET code = 'GRGB-AGENTS' WHERE code = 'grgb.agents';
UPDATE modules SET code = 'GRGB-ROLE' WHERE code = 'grgb.role';
UPDATE modules SET code = 'TPAY-DASHBOARD' WHERE code = 'tpay.dashboard';
UPDATE modules SET code = 'TPAY-WALLET' WHERE code = 'tpay.wallet';
UPDATE modules SET code = 'TPAY-TRX' WHERE code = 'tpay.trx';
UPDATE modules SET code = 'TPAY-BALANCE' WHERE code = 'tpay.balance';
UPDATE modules SET code = 'TPAY-ROLE' WHERE code = 'tpay.role';
UPDATE modules SET code = 'TPAY-USERS' WHERE code = 'tpay.users';
UPDATE modules SET code = 'TPAY-PRODUCTS' WHERE code = 'tpay.products';

-- ============================================
-- REVERT ROLES
-- ============================================
UPDATE roles SET code = 'ADMIN-SUPERADMIN' WHERE code = 'admin.superadmin';
UPDATE roles SET code = 'ADMIN-MANAGER' WHERE code = 'admin.manager';
UPDATE roles SET code = 'TPAY-SYSADMIN' WHERE code = 'tpay.sysadmin';
UPDATE roles SET code = 'TPAY-MANAGER' WHERE code = 'tpay.manager';
UPDATE roles SET code = 'GRGA-SYSADMIN' WHERE code = 'grga.sysadmin';
UPDATE roles SET code = 'GRGA-MANAGER' WHERE code = 'grga.manager';
UPDATE roles SET code = 'GRGA-USER' WHERE code = 'grga.user';
UPDATE roles SET code = 'GRGB-SYSADMIN' WHERE code = 'grgb.sysadmin';
UPDATE roles SET code = 'GRGB-USER' WHERE code = 'grgb.user';

-- ============================================
-- REVERT PERMISSIONS (abbreviated - full list in Up migration)
-- ============================================
UPDATE permissions SET code = 'ADMIN-DASHBOARD-VIEW' WHERE code = 'admin.dashboard.read';
UPDATE permissions SET code = 'ADMIN-ORG-TYPE-CREATE' WHERE code = 'admin.orgtype.create';
UPDATE permissions SET code = 'ADMIN-ORG-TYPE-READ' WHERE code = 'admin.orgtype.read';
UPDATE permissions SET code = 'ADMIN-ORG-TYPE-UPDATE' WHERE code = 'admin.orgtype.update';
UPDATE permissions SET code = 'ADMIN-ORG-TYPE-DELETE' WHERE code = 'admin.orgtype.delete';
UPDATE permissions SET code = 'ADMIN-ORG-CREATE' WHERE code = 'admin.org.create';
UPDATE permissions SET code = 'ADMIN-ORG-READ' WHERE code = 'admin.org.read';
UPDATE permissions SET code = 'ADMIN-ORG-UPDATE' WHERE code = 'admin.org.update';
UPDATE permissions SET code = 'ADMIN-ORG-DELETE' WHERE code = 'admin.org.delete';
UPDATE permissions SET code = 'ADMIN-SYSTEM-CREATE' WHERE code = 'admin.system.create';
UPDATE permissions SET code = 'ADMIN-SYSTEM-READ' WHERE code = 'admin.system.read';
UPDATE permissions SET code = 'ADMIN-SYSTEM-UPDATE' WHERE code = 'admin.system.update';
UPDATE permissions SET code = 'ADMIN-SYSTEM-DELETE' WHERE code = 'admin.system.delete';
UPDATE permissions SET code = 'ADMIN-MODULE-CREATE' WHERE code = 'admin.module.create';
UPDATE permissions SET code = 'ADMIN-MODULE-READ' WHERE code = 'admin.module.read';
UPDATE permissions SET code = 'ADMIN-MODULE-UPDATE' WHERE code = 'admin.module.update';
UPDATE permissions SET code = 'ADMIN-MODULE-DELETE' WHERE code = 'admin.module.delete';
UPDATE permissions SET code = 'ADMIN-MODULE-PERMISSION-CREATE' WHERE code = 'admin.permission.create';
UPDATE permissions SET code = 'ADMIN-MODULE-PERMISSION-READ' WHERE code = 'admin.permission.read';
UPDATE permissions SET code = 'ADMIN-MODULE-PERMISSION-UPDATE' WHERE code = 'admin.permission.update';
UPDATE permissions SET code = 'ADMIN-MODULE-PERMISSION-DELETE' WHERE code = 'admin.permission.delete';
UPDATE permissions SET code = 'ADMIN-ROLE-CREATE' WHERE code = 'admin.role.create';
UPDATE permissions SET code = 'ADMIN-ROLE-READ' WHERE code = 'admin.role.read';
UPDATE permissions SET code = 'ADMIN-ROLE-UPDATE' WHERE code = 'admin.role.update';
UPDATE permissions SET code = 'ADMIN-ROLE-DELETE' WHERE code = 'admin.role.delete';
UPDATE permissions SET code = 'ADMIN-USER-CREATE' WHERE code = 'admin.user.create';
UPDATE permissions SET code = 'ADMIN-USER-READ' WHERE code = 'admin.user.read';
UPDATE permissions SET code = 'ADMIN-USER-UPDATE' WHERE code = 'admin.user.update';
UPDATE permissions SET code = 'ADMIN-USER-DELETE' WHERE code = 'admin.user.delete';
UPDATE permissions SET code = 'ADMIN-AI-CHAT-CREATE' WHERE code = 'admin.aichat.create';
UPDATE permissions SET code = 'ADMIN-AI-CHAT-READ' WHERE code = 'admin.aichat.read';
UPDATE permissions SET code = 'ADMIN-AI-CHAT-UPDATE' WHERE code = 'admin.aichat.update';
UPDATE permissions SET code = 'ADMIN-AI-CHAT-DELETE' WHERE code = 'admin.aichat.delete';
UPDATE permissions SET code = 'GRGA-DASHBOARD-READ' WHERE code = 'grga.dashboard.read';
UPDATE permissions SET code = 'GRGA-ICON-CREATE' WHERE code = 'grga.icon.create';
UPDATE permissions SET code = 'GRGA-ICON-READ' WHERE code = 'grga.icon.read';
UPDATE permissions SET code = 'GRGA-ICON-UPDATE' WHERE code = 'grga.icon.update';
UPDATE permissions SET code = 'GRGA-ICON-DELETE' WHERE code = 'grga.icon.delete';
UPDATE permissions SET code = 'GRGA-ROLE-CREATE' WHERE code = 'grga.role.create';
UPDATE permissions SET code = 'GRGA-ROLE-READ' WHERE code = 'grga.role.read';
UPDATE permissions SET code = 'GRGA-ROLE-UPDATE' WHERE code = 'grga.role.update';
UPDATE permissions SET code = 'GRGA-ROLE-DELETE' WHERE code = 'grga.role.delete';
UPDATE permissions SET code = 'GRGB-DASHBOARD-READ' WHERE code = 'grgb.dashboard.read';
UPDATE permissions SET code = 'GRGB-ICON-CREATE' WHERE code = 'grgb.icon.create';
UPDATE permissions SET code = 'GRGB-ICON-READ' WHERE code = 'grgb.icon.read';
UPDATE permissions SET code = 'GRGB-ICON-UPDATE' WHERE code = 'grgb.icon.update';
UPDATE permissions SET code = 'GRGB-ICON-DELETE' WHERE code = 'grgb.icon.delete';
UPDATE permissions SET code = 'GRGB-ROLE-CREATE' WHERE code = 'grgb.role.create';
UPDATE permissions SET code = 'GRGB-ROLE-READ' WHERE code = 'grgb.role.read';
UPDATE permissions SET code = 'GRGB-ROLE-UPDATE' WHERE code = 'grgb.role.update';
UPDATE permissions SET code = 'GRGB-ROLE-DELETE' WHERE code = 'grgb.role.delete';
UPDATE permissions SET code = 'GRGB-AGENTS-CREATE' WHERE code = 'grgb.agents.create';
UPDATE permissions SET code = 'GRGB-AGENTS-READ' WHERE code = 'grgb.agents.read';
UPDATE permissions SET code = 'GRGB-AGENTS-UPDATE' WHERE code = 'grgb.agents.update';
UPDATE permissions SET code = 'GRGB-AGENTS-DELETE' WHERE code = 'grgb.agents.delete';
UPDATE permissions SET code = 'TPAY-DASHBOARD-READ' WHERE code = 'tpay.dashboard.read';
UPDATE permissions SET code = 'TPAY-BALANCE-READ' WHERE code = 'tpay.balance.read';
UPDATE permissions SET code = 'TPAY-TRX-READ' WHERE code = 'tpay.trx.read';
UPDATE permissions SET code = 'TPAY-WALLET-READ' WHERE code = 'tpay.wallet.read';
UPDATE permissions SET code = 'TPAY-ROLE-CREATE' WHERE code = 'tpay.role.create';
UPDATE permissions SET code = 'TPAY-ROLE-READ' WHERE code = 'tpay.role.read';
UPDATE permissions SET code = 'TPAY-ROLE-UPDATE' WHERE code = 'tpay.role.update';
UPDATE permissions SET code = 'TPAY-ROLE-DELETE' WHERE code = 'tpay.role.delete';
UPDATE permissions SET code = 'TPAY-USER-READ-CREATE' WHERE code = 'tpay.users.create';
UPDATE permissions SET code = 'TPAY-USER-READ-READ' WHERE code = 'tpay.users.read';
UPDATE permissions SET code = 'TPAY-USER-READ-UPDATE' WHERE code = 'tpay.users.update';
UPDATE permissions SET code = 'TPAY-USER-READ-DELETE' WHERE code = 'tpay.users.delete';
UPDATE permissions SET code = 'TPAY-PRODUCTS-CREATE' WHERE code = 'tpay.products.create';
UPDATE permissions SET code = 'TPAY-PRODUCTS-READ' WHERE code = 'tpay.products.read';
UPDATE permissions SET code = 'TPAY-PRODUCTS-UPDATE' WHERE code = 'tpay.products.update';
UPDATE permissions SET code = 'TPAY-PRODUCTS-DELETE' WHERE code = 'tpay.products.delete';

-- ============================================
-- REVERT MENUS
-- ============================================
UPDATE menus SET code = 'ADMIN' WHERE code = 'admin';
UPDATE menus SET code = 'ADMIN-DASHBOARD' WHERE code = 'admin.dashboard';
UPDATE menus SET code = 'ADMIN-SYSTEM' WHERE code = 'admin.system';
UPDATE menus SET code = 'ADMIN-MODULE' WHERE code = 'admin.module';
UPDATE menus SET code = 'ADMIN-ROLE' WHERE code = 'admin.role';
UPDATE menus SET code = 'ADMIN-MODULE-PERMISSION' WHERE code = 'admin.permission';
UPDATE menus SET code = 'ADMIN-ORG' WHERE code = 'admin.org';
UPDATE menus SET code = 'ADMIN-ORGTYPE' WHERE code = 'admin.orgtype';
UPDATE menus SET code = 'ADMIN-USER' WHERE code = 'admin.user';
UPDATE menus SET code = 'ADMIN-AI-CHAT' WHERE code = 'admin.aichat';
UPDATE menus SET code = 'TPAY' WHERE code = 'tpay';
UPDATE menus SET code = 'TPAY-DASHBOARD' WHERE code = 'tpay.dashboard';
UPDATE menus SET code = 'TPAY-WALLET' WHERE code = 'tpay.wallet';
UPDATE menus SET code = 'TPAY-TRX' WHERE code = 'tpay.trx';
UPDATE menus SET code = 'TPAY-BALANCE' WHERE code = 'tpay.balance';
UPDATE menus SET code = 'TPAY-USERS' WHERE code = 'tpay.users';
UPDATE menus SET code = 'TPAY-PRODUCTS' WHERE code = 'tpay.products';
UPDATE menus SET code = 'TPAY-ROLE' WHERE code = 'tpay.role';
UPDATE menus SET code = 'GRGA' WHERE code = 'grga';
UPDATE menus SET code = 'GRGA-DASHBOARD' WHERE code = 'grga.dashboard';
UPDATE menus SET code = 'GRGA-ICON' WHERE code = 'grga.icon';
UPDATE menus SET code = 'GRGA-ROLE' WHERE code = 'grga.role';
UPDATE menus SET code = 'GRGB' WHERE code = 'grgb';
UPDATE menus SET code = 'GRGB-DASHBOARD' WHERE code = 'grgb.dashboard';
UPDATE menus SET code = 'GRGB-ICON' WHERE code = 'grgb.icon';
UPDATE menus SET code = 'GRGB-AGENTS' WHERE code = 'grgb.agents';
UPDATE menus SET code = 'GRGB-ROLE' WHERE code = 'grgb.role';

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd
