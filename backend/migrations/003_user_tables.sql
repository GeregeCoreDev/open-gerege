-- +goose Up
-- +goose StatementBegin
SET search_path TO template_backend;

-- ============================================
-- CITIZENS TABLE (domain/citizen.go)
-- ============================================
CREATE TABLE IF NOT EXISTS citizens (
    id                  SERIAL PRIMARY KEY,
    civil_id            BIGINT DEFAULT 0,
    reg_no              VARCHAR(10),
    family_name         VARCHAR(80),
    last_name           VARCHAR(150),
    first_name          VARCHAR(150),
    gender              INTEGER DEFAULT 0,
    birth_date          VARCHAR(10),
    phone_no            VARCHAR(8),
    email               VARCHAR(80),
    is_foreign          INTEGER DEFAULT 0,
    country_code        VARCHAR(3),
    hash                VARCHAR(200),
    parent_address_id   INTEGER DEFAULT 0,
    parent_address_name VARCHAR(20),
    aimag_id            INTEGER DEFAULT 0,
    aimag_code          VARCHAR(3),
    aimag_name          VARCHAR(255),
    sum_id              INTEGER DEFAULT 0,
    sum_code            VARCHAR(3),
    sum_name            VARCHAR(255),
    bag_id              INTEGER DEFAULT 0,
    bag_code            VARCHAR(3),
    bag_name            VARCHAR(255),
    address_detail      VARCHAR(255),
    address_type        VARCHAR(255),
    address_type_name   VARCHAR(255),
    nationality         VARCHAR(255),
    country_name        VARCHAR(255),
    country_name_en     VARCHAR(255),
    profile_img_url     VARCHAR(255),
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
CREATE INDEX IF NOT EXISTS idx_citizens_reg_no ON citizens(reg_no);
CREATE INDEX IF NOT EXISTS idx_citizens_deleted ON citizens(deleted_date);

DROP TRIGGER IF EXISTS trg_citizens_ins ON citizens;
CREATE TRIGGER trg_citizens_ins BEFORE INSERT ON citizens FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_citizens_upd ON citizens;
CREATE TRIGGER trg_citizens_upd BEFORE UPDATE ON citizens FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- USERS TABLE (domain/user.go)
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    civil_id        BIGINT DEFAULT 0,
    reg_no          VARCHAR(10),
    family_name     VARCHAR(80),
    last_name       VARCHAR(150),
    first_name      VARCHAR(150),
    gender          INTEGER DEFAULT 0,
    birth_date      VARCHAR(10),
    phone_no        VARCHAR(8),
    email           VARCHAR(80),
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
CREATE INDEX IF NOT EXISTS idx_users_reg_no ON users(reg_no);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted ON users(deleted_date);

DROP TRIGGER IF EXISTS trg_users_ins ON users;
CREATE TRIGGER trg_users_ins BEFORE INSERT ON users FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_users_upd ON users;
CREATE TRIGGER trg_users_upd BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- USER_ROLES TABLE (domain/user.go)
-- ============================================
CREATE TABLE IF NOT EXISTS user_roles (
    user_id         INTEGER NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
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
    PRIMARY KEY (user_id, role_id)
);
CREATE INDEX IF NOT EXISTS idx_user_roles_deleted ON user_roles(deleted_date);

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO template_backend;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS citizens CASCADE;
-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd
