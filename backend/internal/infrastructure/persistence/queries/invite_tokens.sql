-- name: CreateInviteToken :one
INSERT INTO invite_tokens (token, created_by, expires_at)
VALUES ($1, $2, $3)
RETURNING id, token, created_by, used, expires_at, created_at;

-- name: GetInviteToken :one
SELECT id, token, created_by, used, expires_at, created_at FROM invite_tokens WHERE token = $1;

-- name: MarkInviteTokenAsUsed :one
UPDATE invite_tokens SET used = true WHERE id = $1 RETURNING id, token, created_by, used, expires_at, created_at;

-- name: ListInviteTokensByCreator :many
SELECT id, token, created_by, used, expires_at, created_at FROM invite_tokens WHERE created_by = $1 ORDER BY created_at DESC;

-- name: GetInviteTokenByID :one
SELECT id, token, created_by, used, expires_at, created_at FROM invite_tokens WHERE id = $1;

-- name: ListAllInviteTokens :many
SELECT it.id, it.token, it.created_by, it.used, it.expires_at, it.created_at, u.full_name AS creator_name
FROM invite_tokens it
LEFT JOIN users u ON u.id = it.created_by
ORDER BY it.created_at DESC;

-- name: DeleteInviteToken :exec
DELETE FROM invite_tokens WHERE id = $1;

-- name: CleanupExpiredInviteTokens :exec
DELETE FROM invite_tokens WHERE expires_at < NOW() OR used = true;
