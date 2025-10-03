-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_update_vehicle_status ON vehicles;
DROP FUNCTION IF EXISTS update_vehicle_status_on_assignment();

-- Drop indexes
DROP INDEX IF EXISTS idx_vehicles_type;
DROP INDEX IF EXISTS idx_vehicles_driver_id;
DROP INDEX IF EXISTS idx_vehicles_status;

-- Remove unique constraint
ALTER TABLE vehicles DROP CONSTRAINT IF EXISTS unique_driver_vehicle;

-- Remove columns
ALTER TABLE vehicles 
DROP COLUMN IF EXISTS capacity,
DROP COLUMN IF EXISTS status;

-- Drop enum type
DROP TYPE IF EXISTS vehicle_status;

