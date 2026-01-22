-- +goose Up
-- +goose StatementBegin

-- Schema and Extensions
CREATE SCHEMA IF NOT EXISTS template_backend;
SET search_path TO template_backend;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA template_backend;

-- Timestamp Functions
CREATE OR REPLACE FUNCTION set_timestamps_on_insert()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_date := COALESCE(NEW.created_date, NOW());
    NEW.updated_date := COALESCE(NEW.updated_date, NOW());
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION set_updated_date_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_date := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_updated_date_timestamp() CASCADE;
DROP FUNCTION IF EXISTS set_timestamps_on_insert() CASCADE;
DROP SCHEMA IF EXISTS template_backend CASCADE;
-- +goose StatementEnd
