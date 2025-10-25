-- Drop driver_feedback table
DROP TABLE IF EXISTS driver_feedback;

-- Remove added columns from reservations
ALTER TABLE reservations DROP COLUMN IF EXISTS distance_km;
ALTER TABLE reservations DROP COLUMN IF EXISTS arrived_on_time;
