-- 000001_init.sql
-- Initial baseline schema for TikTokUp PostgreSQL 18 database.

CREATE TABLE IF NOT EXISTS app_settings (
    key TEXT PRIMARY KEY,
    value JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO app_settings (key, value)
VALUES (
    'platform',
    '{"name": "TikTokUp", "schema_initialized": true}'::jsonb
)
ON CONFLICT (key) DO NOTHING;
