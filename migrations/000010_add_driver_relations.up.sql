-- Add new columns to drivers table for relationships
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS user_id UUID NULL REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS company_id UUID NULL REFERENCES companies(id) ON DELETE SET NULL;
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS vehicle_id VARCHAR(20) NULL REFERENCES vehicles(id) ON DELETE SET NULL;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_drivers_user_id ON drivers(user_id);
CREATE INDEX IF NOT EXISTS idx_drivers_company_id ON drivers(company_id);
CREATE INDEX IF NOT EXISTS idx_drivers_vehicle_id ON drivers(vehicle_id);

-- Add constraint to ensure a driver can only be assigned to one vehicle
CREATE UNIQUE INDEX IF NOT EXISTS idx_drivers_unique_vehicle ON drivers(vehicle_id) WHERE vehicle_id IS NOT NULL;
