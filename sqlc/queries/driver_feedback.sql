-- name: CreateDriverFeedback :one
INSERT INTO driver_feedback (driver_id, reservation_id, rating, comment)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetDriverFeedback :many
SELECT * FROM driver_feedback 
WHERE driver_id = $1 
ORDER BY created_at DESC;

-- name: GetDriverRealKPIs :one
SELECT 
    -- Total trips completed
    (SELECT COUNT(*) FROM reservations r
     WHERE r.assigned_driver_id = $1 AND r.status = 'COMPLETADA') as total_trips,
    
    -- Total kilometers
    (SELECT COALESCE(SUM(r.distance_km), 0) FROM reservations r
     WHERE r.assigned_driver_id = $1 AND r.status = 'COMPLETADA') as total_km,
    
    -- On-time rate (percentage)
    (SELECT CASE 
        WHEN COUNT(*) = 0 THEN 0
        ELSE ROUND(
            (COUNT(CASE WHEN r.arrived_on_time = true THEN 1 END) * 100.0 / COUNT(*))::DECIMAL, 1
        )
     END FROM reservations r
     WHERE r.assigned_driver_id = $1 AND r.status = 'COMPLETADA') as on_time_rate,
    
    -- Cancel rate (percentage)
    (SELECT CASE 
        WHEN COUNT(*) = 0 THEN 0
        ELSE ROUND(
            (COUNT(CASE WHEN r.status = 'CANCELADA' THEN 1 END) * 100.0 / COUNT(*))::DECIMAL, 1
        )
     END FROM reservations r
     WHERE r.assigned_driver_id = $1) as cancel_rate,
    
    -- Average rating
    (SELECT COALESCE(ROUND(AVG(df.rating)::DECIMAL, 1), 0) FROM driver_feedback df
     WHERE df.driver_id = $1) as average_rating;
