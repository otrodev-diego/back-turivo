-- Drop existing driver_feedback table if it exists
DROP TABLE IF EXISTS driver_feedback;

-- Create driver_feedback table for real driver ratings and feedback
CREATE TABLE driver_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    driver_id VARCHAR(20) NOT NULL REFERENCES drivers(id) ON DELETE CASCADE,
    reservation_id VARCHAR(20) NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
    rating DECIMAL(2,1) NOT NULL CHECK (rating >= 1.0 AND rating <= 5.0),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_driver_feedback_driver_id ON driver_feedback(driver_id);
CREATE INDEX idx_driver_feedback_reservation_id ON driver_feedback(reservation_id);
CREATE INDEX idx_driver_feedback_created_at ON driver_feedback(created_at);

-- Add distance_km and arrived_on_time columns to reservations for real KPIs
ALTER TABLE reservations ADD COLUMN IF NOT EXISTS distance_km DECIMAL(10,2);
ALTER TABLE reservations ADD COLUMN IF NOT EXISTS arrived_on_time BOOLEAN DEFAULT true;
