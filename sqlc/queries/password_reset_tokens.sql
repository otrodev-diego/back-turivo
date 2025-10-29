-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPasswordResetTokenByToken :one
SELECT * FROM password_reset_tokens
WHERE token = $1;

-- name: MarkPasswordResetTokenAsUsed :exec
UPDATE password_reset_tokens
SET used = TRUE
WHERE token = $1;

-- name: DeleteExpiredPasswordResetTokens :exec
DELETE FROM password_reset_tokens
WHERE expires_at < NOW();

