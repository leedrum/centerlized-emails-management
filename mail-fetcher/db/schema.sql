-- db/schema.sql
CREATE TABLE IF NOT EXISTS gmail_credentials (
    id SERIAL PRIMARY KEY,
    tenant_name TEXT NOT NULL,
    email TEXT NOT NULL,
    client_id TEXT NOT NULL,
    client_secret TEXT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    token_expiry TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_gmail_credentials_tenant_user
ON gmail_credentials (tenant_name, email);
