-- +goose Up
-- +goose StatementBegin
SET search_path TO template_backend;

-- ============================================
-- ORGANIZATION_TYPES TABLE (domain/organization.go)
-- ============================================
CREATE TABLE IF NOT EXISTS organization_types (
    id              SERIAL PRIMARY KEY,
    code            VARCHAR(255),
    name            VARCHAR(255),
    description     VARCHAR(255),
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
CREATE INDEX IF NOT EXISTS idx_organization_types_deleted ON organization_types(deleted_date);

DROP TRIGGER IF EXISTS trg_organization_types_ins ON organization_types;
CREATE TRIGGER trg_organization_types_ins BEFORE INSERT ON organization_types FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_organization_types_upd ON organization_types;
CREATE TRIGGER trg_organization_types_upd BEFORE UPDATE ON organization_types FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- ORG_TYPE_SYSTEMS TABLE (domain/organization.go)
-- ============================================
CREATE TABLE IF NOT EXISTS org_type_systems (
    type_id         INTEGER NOT NULL REFERENCES organization_types(id) ON UPDATE CASCADE ON DELETE CASCADE,
    system_id       INTEGER NOT NULL REFERENCES systems(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER DEFAULT 0,
    created_org_id  INTEGER DEFAULT 0,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER DEFAULT 0,
    updated_org_id  INTEGER DEFAULT 0,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER DEFAULT 0,
    deleted_org_id  INTEGER DEFAULT 0,
    PRIMARY KEY (type_id, system_id)
);
CREATE INDEX IF NOT EXISTS idx_org_type_systems_deleted ON org_type_systems(deleted_date);

-- ============================================
-- ORG_TYPE_ROLES TABLE (domain/organization.go)
-- ============================================
CREATE TABLE IF NOT EXISTS org_type_roles (
    type_id         INTEGER NOT NULL REFERENCES organization_types(id) ON UPDATE CASCADE ON DELETE CASCADE,
    role_id         INTEGER NOT NULL REFERENCES roles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER DEFAULT 0,
    created_org_id  INTEGER DEFAULT 0,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER DEFAULT 0,
    updated_org_id  INTEGER DEFAULT 0,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER DEFAULT 0,
    deleted_org_id  INTEGER DEFAULT 0,
    PRIMARY KEY (type_id, role_id)
);
CREATE INDEX IF NOT EXISTS idx_org_type_roles_deleted ON org_type_roles(deleted_date);

-- ============================================
-- ORGANIZATIONS TABLE (domain/organization.go)
-- ============================================
CREATE TABLE IF NOT EXISTS organizations (
    id                  SERIAL PRIMARY KEY,
    reg_no              VARCHAR(7),
    name                VARCHAR(255),
    short_name          VARCHAR(255),
    type_id             INTEGER REFERENCES organization_types(id) ON UPDATE CASCADE ON DELETE SET NULL,
    phone_no            VARCHAR(8),
    email               VARCHAR(50),
    longitude           FLOAT8 DEFAULT 106.91758628931501,
    latitude            FLOAT8 DEFAULT 47.918825014251915,
    is_active           BOOLEAN DEFAULT TRUE,
    aimag_id            INTEGER DEFAULT 0,
    sum_id              INTEGER DEFAULT 0,
    bag_id              INTEGER DEFAULT 0,
    address_detail      VARCHAR(255),
    aimag_name          VARCHAR(255),
    sum_name            VARCHAR(255),
    bag_name            VARCHAR(255),
    country_code        VARCHAR(10),
    country_name        VARCHAR(255),
    sequence            INTEGER DEFAULT 0,
    parent_address_id   INTEGER DEFAULT 0,
    parent_address_name VARCHAR(25),
    country_name_en     VARCHAR(255),
    parent_id           INTEGER REFERENCES organizations(id) ON UPDATE CASCADE ON DELETE SET NULL,
    created_date        TIMESTAMPTZ DEFAULT NOW(),
    created_user_id     INTEGER DEFAULT 0,
    created_org_id      INTEGER DEFAULT 0,
    updated_date        TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id     INTEGER DEFAULT 0,
    updated_org_id      INTEGER DEFAULT 0,
    deleted_date        TIMESTAMPTZ,
    deleted_user_id     INTEGER DEFAULT 0,
    deleted_org_id      INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_organizations_type ON organizations(type_id);
CREATE INDEX IF NOT EXISTS idx_organizations_parent ON organizations(parent_id);
CREATE INDEX IF NOT EXISTS idx_organizations_deleted ON organizations(deleted_date);

DROP TRIGGER IF EXISTS trg_organizations_ins ON organizations;
CREATE TRIGGER trg_organizations_ins BEFORE INSERT ON organizations FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_organizations_upd ON organizations;
CREATE TRIGGER trg_organizations_upd BEFORE UPDATE ON organizations FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- ORGANIZATION_USERS TABLE (domain/organization.go)
-- ============================================
CREATE TABLE IF NOT EXISTS organization_users (
    org_id          INTEGER NOT NULL REFERENCES organizations(id) ON UPDATE CASCADE ON DELETE CASCADE,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_date    TIMESTAMPTZ DEFAULT NOW(),
    created_user_id INTEGER DEFAULT 0,
    created_org_id  INTEGER DEFAULT 0,
    updated_date    TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id INTEGER DEFAULT 0,
    updated_org_id  INTEGER DEFAULT 0,
    deleted_date    TIMESTAMPTZ,
    deleted_user_id INTEGER DEFAULT 0,
    deleted_org_id  INTEGER DEFAULT 0,
    PRIMARY KEY (org_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_organization_users_deleted ON organization_users(deleted_date);

-- ============================================
-- TERMINALS TABLE (domain/terminal.go)
-- ============================================
CREATE TABLE IF NOT EXISTS terminals (
    id              SERIAL PRIMARY KEY,
    serial          VARCHAR(80) UNIQUE NOT NULL,
    name            VARCHAR(255),
    org_id          INTEGER REFERENCES organizations(id) ON UPDATE CASCADE ON DELETE SET NULL,
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
CREATE INDEX IF NOT EXISTS idx_terminals_org ON terminals(org_id);
CREATE INDEX IF NOT EXISTS idx_terminals_deleted ON terminals(deleted_date);

DROP TRIGGER IF EXISTS trg_terminals_ins ON terminals;
CREATE TRIGGER trg_terminals_ins BEFORE INSERT ON terminals FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_terminals_upd ON terminals;
CREATE TRIGGER trg_terminals_upd BEFORE UPDATE ON terminals FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO template_backend;
DROP TABLE IF EXISTS terminals CASCADE;
DROP TABLE IF EXISTS organization_users CASCADE;
DROP TABLE IF EXISTS organizations CASCADE;
DROP TABLE IF EXISTS org_type_roles CASCADE;
DROP TABLE IF EXISTS org_type_systems CASCADE;
DROP TABLE IF EXISTS organization_types CASCADE;
-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd
