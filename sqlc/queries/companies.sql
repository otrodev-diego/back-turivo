-- name: CreateCompany :one
INSERT INTO companies (name, rut, contact_email, status, sector)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCompanyByID :one
SELECT * FROM companies WHERE id = $1;

-- name: GetCompanyByRUT :one
SELECT * FROM companies WHERE rut = $1;

-- name: ListCompanies :many
SELECT * FROM companies
WHERE ($1::text IS NULL OR name ILIKE '%' || $1 || '%' OR rut ILIKE '%' || $1 || '%')
  AND ($2::company_status IS NULL OR status = $2)
  AND ($3::company_sector IS NULL OR sector = $3)
ORDER BY 
  CASE WHEN $4 = 'name' THEN name END,
  CASE WHEN $4 = 'rut' THEN rut END,
  CASE WHEN $4 = 'created_at' THEN created_at END DESC
LIMIT $5 OFFSET $6;

-- name: CountCompanies :one
SELECT COUNT(*) FROM companies
WHERE ($1::text IS NULL OR name ILIKE '%' || $1 || '%' OR rut ILIKE '%' || $1 || '%')
  AND ($2::company_status IS NULL OR status = $2)
  AND ($3::company_sector IS NULL OR sector = $3);

-- name: UpdateCompany :one
UPDATE companies
SET name = COALESCE($2, name),
    rut = COALESCE($3, rut),
    contact_email = COALESCE($4, contact_email),
    status = COALESCE($5, status),
    sector = COALESCE($6, sector),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCompany :exec
DELETE FROM companies WHERE id = $1;

