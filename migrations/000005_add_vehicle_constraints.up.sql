-- Add vehicle status enum
CREATE TYPE vehicle_status AS ENUM ('AVAILABLE', 'ASSIGNED', 'MAINTENANCE', 'INACTIVE');

-- Add status and capacity columns to vehicles table
ALTER TABLE vehicles 
ADD COLUMN status vehicle_status NOT NULL DEFAULT 'AVAILABLE',
ADD COLUMN capacity INTEGER NULL CHECK (capacity > 0 AND capacity <= 60);

-- Add unique constraint to ensure a driver can only have one vehicle assigned
-- This will be enforced at the database level
ALTER TABLE vehicles
ADD CONSTRAINT unique_driver_vehicle UNIQUE (driver_id);

-- Add index for status queries
CREATE INDEX idx_vehicles_status ON vehicles(status);

-- Add index for driver_id queries
CREATE INDEX idx_vehicles_driver_id ON vehicles(driver_id);

-- Add index for type queries
CREATE INDEX idx_vehicles_type ON vehicles(type);

-- Update the status when a vehicle is assigned to a driver
CREATE OR REPLACE FUNCTION update_vehicle_status_on_assignment()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.driver_id IS NOT NULL AND (OLD.driver_id IS NULL OR OLD.driver_id != NEW.driver_id) THEN
        NEW.status = 'ASSIGNED';
    ELSIF NEW.driver_id IS NULL AND OLD.driver_id IS NOT NULL THEN
        NEW.status = 'AVAILABLE';
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_vehicle_status 
BEFORE UPDATE ON vehicles 
FOR EACH ROW 
EXECUTE FUNCTION update_vehicle_status_on_assignment();

