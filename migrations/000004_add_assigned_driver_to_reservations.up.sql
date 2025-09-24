-- Add assigned_driver_id field to reservations table
ALTER TABLE reservations 
ADD COLUMN assigned_driver_id VARCHAR(20) NULL REFERENCES drivers(id) ON DELETE SET NULL;

-- Add index for performance
CREATE INDEX idx_reservations_assigned_driver ON reservations(assigned_driver_id);
