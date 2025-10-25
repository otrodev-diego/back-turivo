-- Add user_id to drivers table to link drivers with users
ALTER TABLE drivers ADD COLUMN user_id UUID NULL REFERENCES users(id) ON DELETE SET NULL;

-- Create index for better performance
CREATE INDEX idx_drivers_user_id ON drivers(user_id);

-- Add company_id to drivers table to assign drivers to companies
ALTER TABLE drivers ADD COLUMN company_id UUID NULL REFERENCES companies(id) ON DELETE SET NULL;

-- Create index for better performance
CREATE INDEX idx_drivers_company_id ON drivers(company_id);

-- Add vehicle_id to drivers table for vehicle assignment
ALTER TABLE drivers ADD COLUMN vehicle_id UUID NULL REFERENCES vehicles(id) ON DELETE SET NULL;

-- Create index for better performance
CREATE INDEX idx_drivers_vehicle_id ON drivers(vehicle_id);

-- Add constraint to ensure a driver can only be assigned to one vehicle
CREATE UNIQUE INDEX idx_drivers_unique_vehicle ON drivers(vehicle_id) WHERE vehicle_id IS NOT NULL;
