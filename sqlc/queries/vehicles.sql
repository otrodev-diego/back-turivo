-- name: CreateVehicle :one
INSERT INTO vehicles (
    driver_id,
    type,
    brand,
    model,
    year,
    plate,
    vin,
    color,
    capacity,
    insurance_policy,
    insurance_expires_at,
    inspection_expires_at,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: GetVehicleByID :one
SELECT * FROM vehicles
WHERE id = $1;

-- name: GetVehicleByDriverID :one
SELECT * FROM vehicles
WHERE driver_id = $1
LIMIT 1;

-- name: ListVehicles :many
SELECT * FROM vehicles
WHERE 
    ($1::text IS NULL OR 
        brand ILIKE '%' || $1 || '%' OR 
        model ILIKE '%' || $1 || '%' OR 
        plate ILIKE '%' || $1 || '%')
    AND ($2::vehicle_type IS NULL OR type = $2)
    AND ($3::vehicle_status IS NULL OR status = $3)
    AND ($4::text IS NULL OR driver_id = $4)
ORDER BY
    CASE WHEN $5 = 'brand' THEN brand END,
    CASE WHEN $5 = 'model' THEN model END,
    CASE WHEN $5 = 'year' THEN year::text END,
    CASE WHEN $5 = 'created_at' THEN created_at::text END DESC,
    created_at DESC
LIMIT $6 OFFSET $7;

-- name: CountVehicles :one
SELECT COUNT(*) FROM vehicles
WHERE 
    ($1::text IS NULL OR 
        brand ILIKE '%' || $1 || '%' OR 
        model ILIKE '%' || $1 || '%' OR 
        plate ILIKE '%' || $1 || '%')
    AND ($2::vehicle_type IS NULL OR type = $2)
    AND ($3::vehicle_status IS NULL OR status = $3)
    AND ($4::text IS NULL OR driver_id = $4);

-- name: UpdateVehicle :one
UPDATE vehicles
SET
    type = COALESCE($2, type),
    brand = COALESCE($3, brand),
    model = COALESCE($4, model),
    year = COALESCE($5, year),
    plate = COALESCE($6, plate),
    vin = COALESCE($7, vin),
    color = COALESCE($8, color),
    capacity = COALESCE($9, capacity),
    insurance_policy = COALESCE($10, insurance_policy),
    insurance_expires_at = COALESCE($11, insurance_expires_at),
    inspection_expires_at = COALESCE($12, inspection_expires_at),
    status = COALESCE($13, status),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteVehicle :exec
DELETE FROM vehicles
WHERE id = $1;

-- name: AssignVehicleToDriver :one
UPDATE vehicles
SET 
    driver_id = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UnassignVehicleFromDriver :exec
UPDATE vehicles
SET 
    driver_id = NULL,
    status = 'AVAILABLE',
    updated_at = NOW()
WHERE driver_id = $1;

-- name: CreateVehiclePhoto :one
INSERT INTO vehicle_photos (
    vehicle_id,
    url
) VALUES (
    $1, $2
) RETURNING *;

-- name: DeleteVehiclePhoto :exec
DELETE FROM vehicle_photos
WHERE id = $1;

-- name: GetVehiclePhotos :many
SELECT url FROM vehicle_photos
WHERE vehicle_id = $1
ORDER BY created_at;

