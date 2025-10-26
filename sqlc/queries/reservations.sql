-- name: CreateReservation :one
INSERT INTO reservations (id, user_id, org_id, pickup, destination, datetime, passengers, status, amount, notes, assigned_driver_id, distance_km)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetReservationByID :one
SELECT r.*, d.id as driver_id, d.first_name, d.last_name, d.phone, d.email, d.status as driver_status
FROM reservations r
LEFT JOIN drivers d ON r.assigned_driver_id = d.id
WHERE r.id = $1;

-- name: ListReservations :many
SELECT * FROM reservations
WHERE ($1::text IS NULL OR pickup ILIKE '%' || $1 || '%' OR destination ILIKE '%' || $1 || '%' OR id ILIKE '%' || $1 || '%')
  AND ($2::reservation_status IS NULL OR status = $2)
  AND ($3::uuid IS NULL OR user_id = $3)
  AND ($4::uuid IS NULL OR org_id = $4)
ORDER BY 
  CASE WHEN $5 = 'datetime' THEN datetime END DESC,
  CASE WHEN $5 = 'created_at' THEN created_at END DESC,
  CASE WHEN $5 = 'id' THEN id END
LIMIT $6 OFFSET $7;

-- name: CountReservations :one
SELECT COUNT(*) FROM reservations
WHERE ($1::text IS NULL OR pickup ILIKE '%' || $1 || '%' OR destination ILIKE '%' || $1 || '%' OR id ILIKE '%' || $1 || '%')
  AND ($2::reservation_status IS NULL OR status = $2)
  AND ($3::uuid IS NULL OR user_id = $3)
  AND ($4::uuid IS NULL OR org_id = $4);

-- name: UpdateReservation :one
UPDATE reservations
SET pickup = COALESCE($2, pickup),
    destination = COALESCE($3, destination),
    datetime = COALESCE($4, datetime),
    passengers = COALESCE($5, passengers),
    status = COALESCE($6, status),
    amount = COALESCE($7, amount),
    notes = COALESCE($8, notes),
    assigned_driver_id = COALESCE($9, assigned_driver_id),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateReservationStatus :one
UPDATE reservations
SET status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteReservation :exec
DELETE FROM reservations WHERE id = $1;

-- name: GetReservationsByStatus :many
SELECT * FROM reservations
WHERE status = $1
ORDER BY datetime ASC;

-- name: GetReservationsByDateRange :many
SELECT * FROM reservations
WHERE datetime BETWEEN $1 AND $2
ORDER BY datetime ASC;
