-- name: CreateDriver :one
INSERT INTO drivers (id, first_name, last_name, rut_or_dni, birth_date, phone, email, photo_url, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetDriverByID :one
SELECT d.*, 
       dl.number as license_number, dl.class as license_class, dl.issued_at as license_issued_at, 
       dl.expires_at as license_expires_at, dl.file_url as license_file_url,
       dbc.status as background_status, dbc.file_url as background_file_url, dbc.checked_at as background_checked_at,
       da.regions, da.days, da.time_ranges
FROM drivers d
LEFT JOIN driver_licenses dl ON d.id = dl.driver_id
LEFT JOIN driver_background_checks dbc ON d.id = dbc.driver_id
LEFT JOIN driver_availability da ON d.id = da.driver_id
WHERE d.id = $1;

-- name: ListDrivers :many
SELECT d.*, 
       dl.number as license_number, dl.class as license_class,
       dbc.status as background_status
FROM drivers d
LEFT JOIN driver_licenses dl ON d.id = dl.driver_id
LEFT JOIN driver_background_checks dbc ON d.id = dbc.driver_id
WHERE ($1 = '' OR d.first_name ILIKE '%' || $1 || '%' OR d.last_name ILIKE '%' || $1 || '%' OR d.id ILIKE '%' || $1 || '%')
  AND ($2 = '' OR d.status = $2::driver_status)
ORDER BY 
  CASE WHEN $3 = 'name' THEN d.first_name || ' ' || d.last_name END,
  CASE WHEN $3 = 'id' THEN d.id END,
  CASE WHEN $3 = 'created_at' THEN d.created_at END DESC
LIMIT $4 OFFSET $5;

-- name: CountDrivers :one
SELECT COUNT(*) FROM drivers d
WHERE ($1 = '' OR d.first_name ILIKE '%' || $1 || '%' OR d.last_name ILIKE '%' || $1 || '%' OR d.id ILIKE '%' || $1 || '%')
  AND ($2 = '' OR d.status = $2::driver_status);

-- name: UpdateDriver :one
UPDATE drivers
SET first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    rut_or_dni = COALESCE($4, rut_or_dni),
    birth_date = COALESCE($5, birth_date),
    phone = COALESCE($6, phone),
    email = COALESCE($7, email),
    photo_url = COALESCE($8, photo_url),
    status = COALESCE($9, status),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDriver :exec
DELETE FROM drivers WHERE id = $1;

-- name: CreateDriverLicense :one
INSERT INTO driver_licenses (driver_id, number, class, issued_at, expires_at, file_url)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (driver_id) DO UPDATE SET
    number = EXCLUDED.number,
    class = EXCLUDED.class,
    issued_at = EXCLUDED.issued_at,
    expires_at = EXCLUDED.expires_at,
    file_url = EXCLUDED.file_url
RETURNING *;

-- name: CreateDriverBackgroundCheck :one
INSERT INTO driver_background_checks (driver_id, status, file_url, checked_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (driver_id) DO UPDATE SET
    status = EXCLUDED.status,
    file_url = EXCLUDED.file_url,
    checked_at = EXCLUDED.checked_at
RETURNING *;

-- name: CreateDriverAvailability :one
INSERT INTO driver_availability (driver_id, regions, days, time_ranges)
VALUES ($1, $2, $3, $4)
ON CONFLICT (driver_id) DO UPDATE SET
    regions = EXCLUDED.regions,
    days = EXCLUDED.days,
    time_ranges = EXCLUDED.time_ranges,
    updated_at = NOW()
RETURNING *;

