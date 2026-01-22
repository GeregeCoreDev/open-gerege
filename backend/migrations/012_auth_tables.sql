-- +goose Up
-- +goose StatementBegin
SET search_path TO template_backend;

-- ============================================
-- USER CREDENTIALS TABLE (Local Authentication)
-- ============================================
-- Stores password hashes and account lockout info for local auth
CREATE TABLE IF NOT EXISTS user_credentials (
    id                     SERIAL PRIMARY KEY,
    user_id                INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    password_hash          TEXT NOT NULL,
    password_changed_at    TIMESTAMPTZ DEFAULT NOW(),
    failed_login_attempts  INTEGER DEFAULT 0,
    locked_until           TIMESTAMPTZ,
    must_change_password   BOOLEAN DEFAULT false,
    created_date           TIMESTAMPTZ DEFAULT NOW(),
    created_user_id        INTEGER DEFAULT 0,
    created_org_id         INTEGER DEFAULT 0,
    updated_date           TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id        INTEGER DEFAULT 0,
    updated_org_id         INTEGER DEFAULT 0,
    deleted_date           TIMESTAMPTZ,
    deleted_user_id        INTEGER DEFAULT 0,
    deleted_org_id         INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_user_credentials_user_id ON user_credentials(user_id);
CREATE INDEX IF NOT EXISTS idx_user_credentials_deleted ON user_credentials(deleted_date);

DROP TRIGGER IF EXISTS trg_user_credentials_ins ON user_credentials;
CREATE TRIGGER trg_user_credentials_ins BEFORE INSERT ON user_credentials FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_user_credentials_upd ON user_credentials;
CREATE TRIGGER trg_user_credentials_upd BEFORE UPDATE ON user_credentials FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- USER MFA TOTP TABLE (Multi-Factor Authentication)
-- ============================================
-- Stores TOTP secrets for Google Authenticator style MFA
CREATE TABLE IF NOT EXISTS user_mfa_totp (
    id                     SERIAL PRIMARY KEY,
    user_id                INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    secret_encrypted       TEXT NOT NULL,
    is_enabled             BOOLEAN DEFAULT false,
    verified_at            TIMESTAMPTZ,
    created_date           TIMESTAMPTZ DEFAULT NOW(),
    created_user_id        INTEGER DEFAULT 0,
    created_org_id         INTEGER DEFAULT 0,
    updated_date           TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id        INTEGER DEFAULT 0,
    updated_org_id         INTEGER DEFAULT 0,
    deleted_date           TIMESTAMPTZ,
    deleted_user_id        INTEGER DEFAULT 0,
    deleted_org_id         INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_user_mfa_totp_user_id ON user_mfa_totp(user_id);
CREATE INDEX IF NOT EXISTS idx_user_mfa_totp_deleted ON user_mfa_totp(deleted_date);

DROP TRIGGER IF EXISTS trg_user_mfa_totp_ins ON user_mfa_totp;
CREATE TRIGGER trg_user_mfa_totp_ins BEFORE INSERT ON user_mfa_totp FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_user_mfa_totp_upd ON user_mfa_totp;
CREATE TRIGGER trg_user_mfa_totp_upd BEFORE UPDATE ON user_mfa_totp FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- USER MFA BACKUP CODES TABLE
-- ============================================
-- One-time use recovery codes for MFA
CREATE TABLE IF NOT EXISTS user_mfa_backup_codes (
    id                     SERIAL PRIMARY KEY,
    user_id                INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash              TEXT NOT NULL,
    used_at                TIMESTAMPTZ,
    created_date           TIMESTAMPTZ DEFAULT NOW(),
    created_user_id        INTEGER DEFAULT 0,
    created_org_id         INTEGER DEFAULT 0,
    updated_date           TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id        INTEGER DEFAULT 0,
    updated_org_id         INTEGER DEFAULT 0,
    deleted_date           TIMESTAMPTZ,
    deleted_user_id        INTEGER DEFAULT 0,
    deleted_org_id         INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_user_mfa_backup_codes_user_id ON user_mfa_backup_codes(user_id);
CREATE INDEX IF NOT EXISTS idx_user_mfa_backup_codes_deleted ON user_mfa_backup_codes(deleted_date);

DROP TRIGGER IF EXISTS trg_user_mfa_backup_codes_ins ON user_mfa_backup_codes;
CREATE TRIGGER trg_user_mfa_backup_codes_ins BEFORE INSERT ON user_mfa_backup_codes FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_user_mfa_backup_codes_upd ON user_mfa_backup_codes;
CREATE TRIGGER trg_user_mfa_backup_codes_upd BEFORE UPDATE ON user_mfa_backup_codes FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- SESSIONS TABLE (DB-backed session storage)
-- ============================================
-- Stores session metadata (Redis is primary, DB for audit/backup)
CREATE TABLE IF NOT EXISTS sessions (
    id                     TEXT PRIMARY KEY,
    user_id                INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_address             INET,
    user_agent             TEXT,
    expires_at             TIMESTAMPTZ NOT NULL,
    last_activity_at       TIMESTAMPTZ DEFAULT NOW(),
    revoked_at             TIMESTAMPTZ,
    revoked_reason         TEXT,
    created_date           TIMESTAMPTZ DEFAULT NOW(),
    created_user_id        INTEGER DEFAULT 0,
    created_org_id         INTEGER DEFAULT 0,
    updated_date           TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id        INTEGER DEFAULT 0,
    updated_org_id         INTEGER DEFAULT 0,
    deleted_date           TIMESTAMPTZ,
    deleted_user_id        INTEGER DEFAULT 0,
    deleted_org_id         INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_deleted ON sessions(deleted_date);

DROP TRIGGER IF EXISTS trg_sessions_ins ON sessions;
CREATE TRIGGER trg_sessions_ins BEFORE INSERT ON sessions FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_sessions_upd ON sessions;
CREATE TRIGGER trg_sessions_upd BEFORE UPDATE ON sessions FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- LOGIN HISTORY TABLE
-- ============================================
-- Tracks all login attempts (success and failure)
CREATE TABLE IF NOT EXISTS login_history (
    id                     SERIAL PRIMARY KEY,
    user_id                INTEGER REFERENCES users(id) ON DELETE SET NULL,
    email                  TEXT,
    ip_address             INET,
    user_agent             TEXT,
    login_method           TEXT NOT NULL, -- 'local', 'sso'
    success                BOOLEAN NOT NULL,
    failure_reason         TEXT,
    mfa_used               BOOLEAN DEFAULT false,
    created_date           TIMESTAMPTZ DEFAULT NOW(),
    created_user_id        INTEGER DEFAULT 0,
    created_org_id         INTEGER DEFAULT 0,
    updated_date           TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id        INTEGER DEFAULT 0,
    updated_org_id         INTEGER DEFAULT 0,
    deleted_date           TIMESTAMPTZ,
    deleted_user_id        INTEGER DEFAULT 0,
    deleted_org_id         INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_login_history_user_id ON login_history(user_id);
CREATE INDEX IF NOT EXISTS idx_login_history_created_date ON login_history(created_date);
CREATE INDEX IF NOT EXISTS idx_login_history_deleted ON login_history(deleted_date);

DROP TRIGGER IF EXISTS trg_login_history_ins ON login_history;
CREATE TRIGGER trg_login_history_ins BEFORE INSERT ON login_history FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_login_history_upd ON login_history;
CREATE TRIGGER trg_login_history_upd BEFORE UPDATE ON login_history FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- SECURITY AUDIT TRAIL TABLE
-- ============================================
-- Tracks all security-related events
CREATE TABLE IF NOT EXISTS security_audit_trail (
    id                     SERIAL PRIMARY KEY,
    user_id                INTEGER REFERENCES users(id) ON DELETE SET NULL,
    action                 TEXT NOT NULL, -- 'password_change', 'mfa_enable', 'mfa_disable', 'session_revoke', etc.
    target_type            TEXT, -- 'user', 'session', etc.
    target_id              TEXT, -- ID of the affected entity
    old_value              JSONB,
    new_value              JSONB,
    ip_address             INET,
    user_agent             TEXT,
    created_date           TIMESTAMPTZ DEFAULT NOW(),
    created_user_id        INTEGER DEFAULT 0,
    created_org_id         INTEGER DEFAULT 0,
    updated_date           TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id        INTEGER DEFAULT 0,
    updated_org_id         INTEGER DEFAULT 0,
    deleted_date           TIMESTAMPTZ,
    deleted_user_id        INTEGER DEFAULT 0,
    deleted_org_id         INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_security_audit_trail_user_id ON security_audit_trail(user_id);
CREATE INDEX IF NOT EXISTS idx_security_audit_trail_action ON security_audit_trail(action);
CREATE INDEX IF NOT EXISTS idx_security_audit_trail_created_date ON security_audit_trail(created_date);
CREATE INDEX IF NOT EXISTS idx_security_audit_trail_deleted ON security_audit_trail(deleted_date);

DROP TRIGGER IF EXISTS trg_security_audit_trail_ins ON security_audit_trail;
CREATE TRIGGER trg_security_audit_trail_ins BEFORE INSERT ON security_audit_trail FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_security_audit_trail_upd ON security_audit_trail;
CREATE TRIGGER trg_security_audit_trail_upd BEFORE UPDATE ON security_audit_trail FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- PASSWORD HISTORY TABLE
-- ============================================
-- Prevents password reuse
CREATE TABLE IF NOT EXISTS password_history (
    id                     SERIAL PRIMARY KEY,
    user_id                INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    password_hash          TEXT NOT NULL,
    created_date           TIMESTAMPTZ DEFAULT NOW(),
    created_user_id        INTEGER DEFAULT 0,
    created_org_id         INTEGER DEFAULT 0,
    updated_date           TIMESTAMPTZ DEFAULT NOW(),
    updated_user_id        INTEGER DEFAULT 0,
    updated_org_id         INTEGER DEFAULT 0,
    deleted_date           TIMESTAMPTZ,
    deleted_user_id        INTEGER DEFAULT 0,
    deleted_org_id         INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_password_history_user_id ON password_history(user_id);
CREATE INDEX IF NOT EXISTS idx_password_history_deleted ON password_history(deleted_date);

DROP TRIGGER IF EXISTS trg_password_history_ins ON password_history;
CREATE TRIGGER trg_password_history_ins BEFORE INSERT ON password_history FOR EACH ROW EXECUTE FUNCTION set_timestamps_on_insert();
DROP TRIGGER IF EXISTS trg_password_history_upd ON password_history;
CREATE TRIGGER trg_password_history_upd BEFORE UPDATE ON password_history FOR EACH ROW EXECUTE FUNCTION set_updated_date_timestamp();

-- ============================================
-- EXTEND USERS TABLE
-- ============================================
-- Add status and login tracking columns
ALTER TABLE users ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'active';
ALTER TABLE users ADD COLUMN IF NOT EXISTS status_reason TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS status_changed_at TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS status_changed_by INTEGER REFERENCES users(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS login_count INTEGER DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO template_backend;

-- Remove user status columns
ALTER TABLE users DROP COLUMN IF EXISTS status;
ALTER TABLE users DROP COLUMN IF EXISTS status_reason;
ALTER TABLE users DROP COLUMN IF EXISTS status_changed_at;
ALTER TABLE users DROP COLUMN IF EXISTS status_changed_by;
ALTER TABLE users DROP COLUMN IF EXISTS last_login_at;
ALTER TABLE users DROP COLUMN IF EXISTS login_count;

-- Drop tables in reverse order
DROP TABLE IF EXISTS password_history CASCADE;
DROP TABLE IF EXISTS security_audit_trail CASCADE;
DROP TABLE IF EXISTS login_history CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS user_mfa_backup_codes CASCADE;
DROP TABLE IF EXISTS user_mfa_totp CASCADE;
DROP TABLE IF EXISTS user_credentials CASCADE;

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd
