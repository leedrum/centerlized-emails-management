-- name: UpsertToken :exec
INSERT INTO gmail_credentials (
    tenant_name, email, client_id, client_secret,
    access_token, refresh_token, token_expiry
) VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (tenant_name, email)
DO UPDATE SET
    client_id = EXCLUDED.client_id,
    client_secret = EXCLUDED.client_secret,
    access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token,
    token_expiry = EXCLUDED.token_expiry,
    updated_at = CURRENT_TIMESTAMP;

-- name: GetGmailCredentialsByTenantName :many
SELECT * FROM gmail_credentials
WHERE tenant_name = $1;

-- name: GetOneGmailCredentialByEmail :one
SELECT * FROM gmail_credentials
WHERE email = $1
LIMIT 1;

-- name: ListGmailCredentials :many
SELECT * FROM gmail_credentials
WHERE tenant_name = $1;
