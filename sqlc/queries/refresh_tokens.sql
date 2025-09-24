-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (id, user_id, token, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT rt.*, u.id as user_id, u.name as user_name, u.email as user_email, u.role as user_role, u.status as user_status, u.org_id as user_org_id
FROM refresh_tokens rt
JOIN users u ON rt.user_id = u.id
WHERE rt.token = $1 AND rt.expires_at > NOW();

-- name: GetRefreshTokensByUserID :many
SELECT * FROM refresh_tokens 
WHERE user_id = $1 AND expires_at > NOW()
ORDER BY created_at DESC;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE token = $1;

-- name: DeleteRefreshTokensByUserID :exec
DELETE FROM refresh_tokens WHERE user_id = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens WHERE expires_at <= NOW();

-- name: CountRefreshTokensByUserID :one
SELECT COUNT(*) FROM refresh_tokens 
WHERE user_id = $1 AND expires_at > NOW();
