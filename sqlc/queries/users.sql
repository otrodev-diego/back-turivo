-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, role, org_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users
WHERE ($1::text IS NULL OR name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')
  AND ($2::user_role IS NULL OR role = $2)
  AND ($3::user_status IS NULL OR status = $3)
ORDER BY 
  CASE WHEN $4 = 'name' THEN name END,
  CASE WHEN $4 = 'email' THEN email END,
  CASE WHEN $4 = 'created_at' THEN created_at END DESC
LIMIT $5 OFFSET $6;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE ($1::text IS NULL OR name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')
  AND ($2::user_role IS NULL OR role = $2)
  AND ($3::user_status IS NULL OR status = $3);

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE($2, name),
    email = COALESCE($3, email),
    password_hash = COALESCE($4, password_hash),
    role = COALESCE($5, role),
    status = COALESCE($6, status),
    org_id = COALESCE($7, org_id),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

