-- +goose Up
-- +goose StatementBegin
SET search_path TO template_backend;

-- ============================================
-- PUBLIC_FILES TABLE (domain/file.go)
-- ============================================
CREATE TABLE IF NOT EXISTS public_files (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255),
    extension       VARCHAR(10),
    description     VARCHAR(255),
    file_url        VARCHAR(255),
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
CREATE INDEX IF NOT EXISTS idx_public_files_deleted ON public_files(deleted_date);

DROP TRIGGER IF EXISTS trg_public_files_ins ON public_files;
CREATE TRIGGER trg_public_files_ins BEFORE INSERT ON public_files FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_public_files_upd ON public_files;
CREATE TRIGGER trg_public_files_upd BEFORE UPDATE ON public_files FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- NEWS TABLE (domain/news.go)
-- ============================================
CREATE TABLE IF NOT EXISTS news (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255),
    text            TEXT,
    image_url       VARCHAR(255),
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
CREATE INDEX IF NOT EXISTS idx_news_deleted ON news(deleted_date);

DROP TRIGGER IF EXISTS trg_news_ins ON news;
CREATE TRIGGER trg_news_ins BEFORE INSERT ON news FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_news_upd ON news;
CREATE TRIGGER trg_news_upd BEFORE UPDATE ON news FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- NOTIFICATION_GROUPS TABLE (domain/notification.go)
-- ============================================
CREATE TABLE IF NOT EXISTS notification_groups (
    id               SERIAL PRIMARY KEY,
    user_id          INTEGER DEFAULT 0,
    title            VARCHAR(255),
    content          TEXT,
    type             VARCHAR(20),
    tenant           VARCHAR(50),
    created_username VARCHAR(100),
    created_date     TIMESTAMPTZ DEFAULT NOW(),
    created_user_id  INTEGER DEFAULT 0,
    created_org_id   INTEGER DEFAULT 0,
    updated_date     TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id  INTEGER DEFAULT 0,
    updated_org_id   INTEGER DEFAULT 0,
    deleted_date     TIMESTAMPTZ,
    deleted_user_id  INTEGER DEFAULT 0,
    deleted_org_id   INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_notification_groups_user ON notification_groups(user_id);
CREATE INDEX IF NOT EXISTS idx_notification_groups_deleted ON notification_groups(deleted_date);

DROP TRIGGER IF EXISTS trg_notification_groups_ins ON notification_groups;
CREATE TRIGGER trg_notification_groups_ins BEFORE INSERT ON notification_groups FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_notification_groups_upd ON notification_groups;
CREATE TRIGGER trg_notification_groups_upd BEFORE UPDATE ON notification_groups FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- NOTIFICATIONS TABLE (domain/notification.go)
-- ============================================
CREATE TABLE IF NOT EXISTS notifications (
    id               SERIAL PRIMARY KEY,
    user_id          INTEGER DEFAULT 0,
    title            VARCHAR(255),
    content          TEXT,
    is_read          BOOLEAN DEFAULT FALSE,
    type             VARCHAR(20),
    tenant           VARCHAR(50),
    group_id         INTEGER DEFAULT 0,
    created_username VARCHAR(100),
    created_date     TIMESTAMPTZ DEFAULT NOW(),
    created_user_id  INTEGER DEFAULT 0,
    created_org_id   INTEGER DEFAULT 0,
    updated_date     TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id  INTEGER DEFAULT 0,
    updated_org_id   INTEGER DEFAULT 0,
    deleted_date     TIMESTAMPTZ,
    deleted_user_id  INTEGER DEFAULT 0,
    deleted_org_id   INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_group ON notifications(group_id);
CREATE INDEX IF NOT EXISTS idx_notifications_deleted ON notifications(deleted_date);

DROP TRIGGER IF EXISTS trg_notifications_ins ON notifications;
CREATE TRIGGER trg_notifications_ins BEFORE INSERT ON notifications FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_notifications_upd ON notifications;
CREATE TRIGGER trg_notifications_upd BEFORE UPDATE ON notifications FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- CHAT_ITEMS TABLE (domain/chatbot.go)
-- ============================================
CREATE TABLE IF NOT EXISTS chat_items (
    id              SERIAL PRIMARY KEY,
    key             VARCHAR(255),
    answer          TEXT,
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
CREATE INDEX IF NOT EXISTS idx_chat_items_key ON chat_items(key);
CREATE INDEX IF NOT EXISTS idx_chat_items_deleted ON chat_items(deleted_date);

DROP TRIGGER IF EXISTS trg_chat_items_ins ON chat_items;
CREATE TRIGGER trg_chat_items_ins BEFORE INSERT ON chat_items FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_chat_items_upd ON chat_items;
CREATE TRIGGER trg_chat_items_upd BEFORE UPDATE ON chat_items FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- LOGS TABLE - API Request Logging (domain/log.go - APILog)
-- ============================================
CREATE TABLE IF NOT EXISTS logs (
    id              BIGSERIAL PRIMARY KEY,
    org_id          BIGINT,
    user_id         BIGINT,
    username        VARCHAR(50),
    path            VARCHAR(255),
    method          VARCHAR(10),
    params          JSONB,
    queries         JSONB,
    body            JSONB,
    status_code     INTEGER,
    response        JSONB,
    latency_ms      BIGINT,
    req_size        BIGINT,
    res_size        BIGINT,
    ip              VARCHAR(45),
    created_date    TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_logs_user ON logs(user_id);
CREATE INDEX IF NOT EXISTS idx_logs_org ON logs(org_id);
CREATE INDEX IF NOT EXISTS idx_logs_path ON logs(path);
CREATE INDEX IF NOT EXISTS idx_logs_created ON logs(created_date);

-- ============================================
-- AUDIT_LOGS TABLE (domain/log.go - AuditLog)
-- ============================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER,
    action_id       BIGINT REFERENCES actions(id) ON UPDATE CASCADE ON DELETE SET NULL,
    module_id       BIGINT,
    metadata        JSONB,
    created_date    TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_module ON audit_logs(module_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_date);

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO template_backend;
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS logs CASCADE;
DROP TABLE IF EXISTS chat_items CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS notification_groups CASCADE;
DROP TABLE IF EXISTS news CASCADE;
DROP TABLE IF EXISTS public_files CASCADE;
-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd
