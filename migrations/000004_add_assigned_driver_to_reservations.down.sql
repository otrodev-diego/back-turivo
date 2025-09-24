-- Remove assigned_driver_id field from reservations table
DROP INDEX IF EXISTS idx_reservations_assigned_driver;
ALTER TABLE reservations DROP COLUMN IF EXISTS assigned_driver_id;
