-- +goose Up
-- +goose StatementBegin
SET search_path TO template_backend;

-- ============================================
-- SYSTEMS TABLE (domain/system.go)
-- ============================================
CREATE TABLE IF NOT EXISTS systems (
    id              SERIAL PRIMARY KEY,
    code            VARCHAR(255) UNIQUE,
    key             VARCHAR(255),
    name            VARCHAR(255),
    description     VARCHAR(255),
    is_active       BOOLEAN DEFAULT TRUE,
    icon            VARCHAR(255),
    sequence        INTEGER DEFAULT 0,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER DEFAULT 0,
    created_org_id  INTEGER DEFAULT 0,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER DEFAULT 0,
    updated_org_id  INTEGER DEFAULT 0,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER DEFAULT 0,
    deleted_org_id  INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_systems_deleted ON systems(deleted_date);

-- Triggers
DROP TRIGGER IF EXISTS trg_systems_ins ON systems;
CREATE TRIGGER trg_systems_ins BEFORE INSERT ON systems FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_systems_upd ON systems;
CREATE TRIGGER trg_systems_upd BEFORE UPDATE ON systems FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- ACTIONS TABLE (domain/action.go)
-- ============================================
CREATE TABLE IF NOT EXISTS actions (
    id              BIGSERIAL PRIMARY KEY,
    code            VARCHAR(255) UNIQUE NOT NULL,
    name            VARCHAR(255),
    description     VARCHAR(255),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER DEFAULT 0,
    created_org_id  INTEGER DEFAULT 0,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER DEFAULT 0,
    updated_org_id  INTEGER DEFAULT 0,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER DEFAULT 0,
    deleted_org_id  INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_actions_deleted ON actions(deleted_date);

DROP TRIGGER IF EXISTS trg_actions_ins ON actions;
CREATE TRIGGER trg_actions_ins BEFORE INSERT ON actions FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_actions_upd ON actions;
CREATE TRIGGER trg_actions_upd BEFORE UPDATE ON actions FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- MODULES TABLE (domain/module.go)
-- ============================================
CREATE TABLE IF NOT EXISTS modules (
    id              SERIAL PRIMARY KEY,
    code            VARCHAR(255) UNIQUE,
    name            VARCHAR(255),
    description     VARCHAR(255),
    is_active       BOOLEAN DEFAULT TRUE,
    system_id       INTEGER REFERENCES systems(id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER DEFAULT 0,
    created_org_id  INTEGER DEFAULT 0,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER DEFAULT 0,
    updated_org_id  INTEGER DEFAULT 0,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER DEFAULT 0,
    deleted_org_id  INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_modules_system ON modules(system_id);
CREATE INDEX IF NOT EXISTS idx_modules_deleted ON modules(deleted_date);

DROP TRIGGER IF EXISTS trg_modules_ins ON modules;
CREATE TRIGGER trg_modules_ins BEFORE INSERT ON modules FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_modules_upd ON modules;
CREATE TRIGGER trg_modules_upd BEFORE UPDATE ON modules FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- PERMISSIONS TABLE (domain/permission.go)
-- ============================================
CREATE TABLE IF NOT EXISTS permissions (
    id              SERIAL PRIMARY KEY,
    code            VARCHAR(255) UNIQUE NOT NULL,
    name            VARCHAR(255),
    description     VARCHAR(255),
    system_id       INTEGER REFERENCES systems(id) ON UPDATE CASCADE ON DELETE SET NULL,
    module_id       INTEGER REFERENCES modules(id) ON UPDATE CASCADE ON DELETE SET NULL,
    action_id       BIGINT REFERENCES actions(id) ON UPDATE CASCADE ON DELETE SET NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER DEFAULT 0,
    created_org_id  INTEGER DEFAULT 0,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER DEFAULT 0,
    updated_org_id  INTEGER DEFAULT 0,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER DEFAULT 0,
    deleted_org_id  INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_permissions_system ON permissions(system_id);
CREATE INDEX IF NOT EXISTS idx_permissions_module ON permissions(module_id);
CREATE INDEX IF NOT EXISTS idx_permissions_action ON permissions(action_id);
CREATE INDEX IF NOT EXISTS idx_permissions_deleted ON permissions(deleted_date);

DROP TRIGGER IF EXISTS trg_permissions_ins ON permissions;
CREATE TRIGGER trg_permissions_ins BEFORE INSERT ON permissions FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_permissions_upd ON permissions;
CREATE TRIGGER trg_permissions_upd BEFORE UPDATE ON permissions FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- MENUS TABLE (domain/menu.go)
-- ============================================
CREATE TABLE IF NOT EXISTS menus (
    id              BIGSERIAL PRIMARY KEY,
    code            VARCHAR(255),
    key             VARCHAR(255) UNIQUE,
    name            VARCHAR(255),
    description     VARCHAR(255),
    icon            VARCHAR(255),
    path            VARCHAR(255),
    sequence        BIGINT DEFAULT 0,
    parent_id       BIGINT REFERENCES menus(id) ON UPDATE CASCADE ON DELETE SET NULL,
    permission_id   BIGINT REFERENCES permissions(id) ON UPDATE CASCADE ON DELETE SET NULL,
    is_active       BOOLEAN DEFAULT TRUE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER DEFAULT 0,
    created_org_id  INTEGER DEFAULT 0,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER DEFAULT 0,
    updated_org_id  INTEGER DEFAULT 0,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER DEFAULT 0,
    deleted_org_id  INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_menus_parent ON menus(parent_id);
CREATE INDEX IF NOT EXISTS idx_menus_permission ON menus(permission_id);
CREATE INDEX IF NOT EXISTS idx_menus_deleted ON menus(deleted_date);

DROP TRIGGER IF EXISTS trg_menus_ins ON menus;
CREATE TRIGGER trg_menus_ins BEFORE INSERT ON menus FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_menus_upd ON menus;
CREATE TRIGGER trg_menus_upd BEFORE UPDATE ON menus FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- ROLES TABLE (domain/role.go)
-- ============================================
CREATE TABLE IF NOT EXISTS roles (
    id              SERIAL PRIMARY KEY,
    system_id       INTEGER REFERENCES systems(id) ON UPDATE CASCADE ON DELETE SET NULL,
    code            VARCHAR(255) UNIQUE NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     VARCHAR(255),
    is_active       BOOLEAN DEFAULT TRUE,
    is_system_role  BOOLEAN DEFAULT FALSE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER DEFAULT 0,
    created_org_id  INTEGER DEFAULT 0,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER DEFAULT 0,
    updated_org_id  INTEGER DEFAULT 0,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER DEFAULT 0,
    deleted_org_id  INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_roles_system ON roles(system_id);
CREATE INDEX IF NOT EXISTS idx_roles_deleted ON roles(deleted_date);

DROP TRIGGER IF EXISTS trg_roles_ins ON roles;
CREATE TRIGGER trg_roles_ins BEFORE INSERT ON roles FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_roles_upd ON roles;
CREATE TRIGGER trg_roles_upd BEFORE UPDATE ON roles FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- ROLE_PERMISSIONS TABLE (domain/role.go)
-- ============================================
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id         INTEGER NOT NULL REFERENCES roles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    permission_id   INTEGER NOT NULL REFERENCES permissions(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER DEFAULT 0,
    created_org_id  INTEGER DEFAULT 0,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER DEFAULT 0,
    updated_org_id  INTEGER DEFAULT 0,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER DEFAULT 0,
    deleted_org_id  INTEGER DEFAULT 0,
    PRIMARY KEY (role_id, permission_id)
);
CREATE INDEX IF NOT EXISTS idx_role_permissions_deleted ON role_permissions(deleted_date);

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO template_backend;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS menus CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS modules CASCADE;
DROP TABLE IF EXISTS actions CASCADE;
DROP TABLE IF EXISTS systems CASCADE;
-- +goose StatementEnd
