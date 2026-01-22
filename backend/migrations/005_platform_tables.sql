-- +goose Up
-- +goose StatementBegin
SET search_path TO template_backend;

-- ============================================
-- APP_SERVICE_ICON_GROUPS TABLE (domain/app_icon.go)
-- ============================================
CREATE TABLE IF NOT EXISTS app_service_icon_groups (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(80) NOT NULL,
    name_en         VARCHAR(80),
    icon            VARCHAR(255),
    type_name       VARCHAR(255) DEFAULT 'group',
    seq             INTEGER DEFAULT 1,
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
CREATE INDEX IF NOT EXISTS idx_app_service_icon_groups_deleted ON app_service_icon_groups(deleted_date);

DROP TRIGGER IF EXISTS trg_app_service_icon_groups_ins ON app_service_icon_groups;
CREATE TRIGGER trg_app_service_icon_groups_ins BEFORE INSERT ON app_service_icon_groups FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_app_service_icon_groups_upd ON app_service_icon_groups;
CREATE TRIGGER trg_app_service_icon_groups_upd BEFORE UPDATE ON app_service_icon_groups FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- APP_SERVICE_ICONS TABLE (domain/app_icon.go)
-- ============================================
CREATE TABLE IF NOT EXISTS app_service_icons (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    name_en         VARCHAR(255),
    icon            VARCHAR(255),
    icon_app        VARCHAR(255),
    icon_tablet     VARCHAR(255),
    icon_kiosk      VARCHAR(255),
    link            VARCHAR(255),
    group_id        INTEGER REFERENCES app_service_icon_groups(id) ON UPDATE CASCADE ON DELETE SET NULL,
    seq             INTEGER DEFAULT 1,
    is_native       BOOLEAN DEFAULT FALSE,
    is_public       BOOLEAN DEFAULT TRUE,
    is_featured     BOOLEAN DEFAULT FALSE,
    featured_icon   VARCHAR(255),
    is_best_selling BOOLEAN DEFAULT FALSE,
    feature_seq     INTEGER DEFAULT 1,
    description     TEXT,
    system_code     VARCHAR(2),
    is_group        BOOLEAN DEFAULT FALSE,
    parent_id       INTEGER REFERENCES app_service_icons(id) ON UPDATE CASCADE ON DELETE SET NULL,
    web_link        VARCHAR(255),
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
CREATE INDEX IF NOT EXISTS idx_app_service_icons_group ON app_service_icons(group_id);
CREATE INDEX IF NOT EXISTS idx_app_service_icons_parent ON app_service_icons(parent_id);
CREATE INDEX IF NOT EXISTS idx_app_service_icons_deleted ON app_service_icons(deleted_date);

DROP TRIGGER IF EXISTS trg_app_service_icons_ins ON app_service_icons;
CREATE TRIGGER trg_app_service_icons_ins BEFORE INSERT ON app_service_icons FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_app_service_icons_upd ON app_service_icons;
CREATE TRIGGER trg_app_service_icons_upd BEFORE UPDATE ON app_service_icons FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- VEHICLES TABLE (domain/vehicle.go)
-- ============================================
CREATE TABLE IF NOT EXISTS vehicles (
    id              SERIAL PRIMARY KEY,
    plate_no        VARCHAR(7) UNIQUE NOT NULL
);

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO template_backend;
DROP TABLE IF EXISTS vehicles CASCADE;
DROP TABLE IF EXISTS app_service_icons CASCADE;
DROP TABLE IF EXISTS app_service_icon_groups CASCADE;
-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd
