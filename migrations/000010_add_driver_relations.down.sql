-- Remove constraints and indexes
DROP INDEX IF EXISTS idx_drivers_unique_vehicle;
DROP INDEX IF EXISTS idx_drivers_vehicle_id;
DROP INDEX IF EXISTS idx_drivers_company_id;
DROP INDEX IF EXISTS idx_drivers_user_id;

-- Remove columns
ALTER TABLE drivers DROP COLUMN IF EXISTS vehicle_id;
ALTER TABLE drivers DROP COLUMN IF EXISTS company_id;
ALTER TABLE drivers DROP COLUMN IF EXISTS user_id;
