-- name: CreateTimelineEvent :one
INSERT INTO reservation_timeline (id, reservation_id, title, description, at, variant)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTimelineEventByID :one
SELECT * FROM reservation_timeline WHERE id = $1;

-- name: GetTimelineByReservationID :many
SELECT * FROM reservation_timeline 
WHERE reservation_id = $1 
ORDER BY at ASC;

-- name: ListTimelineEvents :many
SELECT rt.*, r.pickup, r.destination
FROM reservation_timeline rt
LEFT JOIN reservations r ON rt.reservation_id = r.id
WHERE ($1::text IS NULL OR rt.reservation_id ILIKE '%' || $1 || '%' OR rt.title ILIKE '%' || $1 || '%')
ORDER BY rt.at DESC
LIMIT $2 OFFSET $3;

-- name: CountTimelineEvents :one
SELECT COUNT(*) FROM reservation_timeline rt
WHERE ($1::text IS NULL OR rt.reservation_id ILIKE '%' || $1 || '%' OR rt.title ILIKE '%' || $1 || '%');

-- name: UpdateTimelineEvent :one
UPDATE reservation_timeline
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    at = COALESCE($4, at),
    variant = COALESCE($5, variant)
WHERE id = $1
RETURNING *;

-- name: DeleteTimelineEvent :exec
DELETE FROM reservation_timeline WHERE id = $1;

-- name: DeleteTimelineByReservationID :exec
DELETE FROM reservation_timeline WHERE reservation_id = $1;
