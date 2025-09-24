-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create ENUM types
CREATE TYPE user_role AS ENUM ('ADMIN', 'USER', 'DRIVER', 'COMPANY');
CREATE TYPE user_status AS ENUM ('ACTIVE', 'BLOCKED');
CREATE TYPE company_status AS ENUM ('ACTIVE', 'SUSPENDED');
CREATE TYPE company_sector AS ENUM ('HOTEL', 'MINERIA', 'TURISMO');
CREATE TYPE driver_status AS ENUM ('ACTIVE', 'INACTIVE');
CREATE TYPE license_class AS ENUM ('A1', 'A2', 'A3', 'A4', 'A5', 'B', 'C', 'D', 'E');
CREATE TYPE background_check_status AS ENUM ('APPROVED', 'PENDING', 'REJECTED');
CREATE TYPE vehicle_type AS ENUM ('BUS', 'VAN', 'SEDAN', 'SUV');
CREATE TYPE request_status AS ENUM ('PENDIENTE', 'ASIGNADA', 'EN_RUTA', 'COMPLETADA', 'CANCELADA');
CREATE TYPE reservation_status AS ENUM ('ACTIVA', 'PROGRAMADA', 'COMPLETADA', 'CANCELADA');
CREATE TYPE payment_gateway AS ENUM ('WEBPAY_PLUS');
CREATE TYPE payment_status AS ENUM ('APPROVED', 'REJECTED', 'PENDING');
CREATE TYPE language AS ENUM ('es', 'en', 'pt', 'fr');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'USER',
    status user_status NOT NULL DEFAULT 'ACTIVE',
    org_id UUID NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Companies table
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    rut VARCHAR(50) UNIQUE NOT NULL,
    contact_email VARCHAR(255) NOT NULL,
    status company_status NOT NULL DEFAULT 'ACTIVE',
    sector company_sector NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Hotels table
CREATE TABLE hotels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Drivers table
CREATE TABLE drivers (
    id VARCHAR(20) PRIMARY KEY, -- e.g., CON-001
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    rut_or_dni VARCHAR(50) UNIQUE NOT NULL,
    birth_date DATE NULL,
    phone VARCHAR(50) NULL,
    email VARCHAR(255) NULL,
    photo_url VARCHAR(500) NULL,
    status driver_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Driver licenses table
CREATE TABLE driver_licenses (
    driver_id VARCHAR(20) NOT NULL REFERENCES drivers(id) ON DELETE CASCADE,
    number VARCHAR(100) NOT NULL,
    class license_class NOT NULL,
    issued_at DATE NULL,
    expires_at DATE NULL,
    file_url VARCHAR(500) NULL,
    PRIMARY KEY (driver_id)
);

-- Driver background checks table
CREATE TABLE driver_background_checks (
    driver_id VARCHAR(20) NOT NULL REFERENCES drivers(id) ON DELETE CASCADE,
    status background_check_status NOT NULL DEFAULT 'PENDING',
    file_url VARCHAR(500) NULL,
    checked_at TIMESTAMPTZ NULL,
    PRIMARY KEY (driver_id)
);

-- Vehicles table
CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    driver_id VARCHAR(20) NULL REFERENCES drivers(id) ON DELETE SET NULL,
    type vehicle_type NOT NULL,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER NULL,
    plate VARCHAR(20) NULL,
    vin VARCHAR(100) NULL,
    color VARCHAR(50) NULL,
    insurance_policy VARCHAR(100) NULL,
    insurance_expires_at DATE NULL,
    inspection_expires_at DATE NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Vehicle photos table
CREATE TABLE vehicle_photos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    url VARCHAR(500) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Driver availability table
CREATE TABLE driver_availability (
    driver_id VARCHAR(20) NOT NULL REFERENCES drivers(id) ON DELETE CASCADE,
    regions JSONB NOT NULL DEFAULT '[]',
    days JSONB NOT NULL DEFAULT '[]',
    time_ranges JSONB NOT NULL DEFAULT '[]',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (driver_id)
);

-- Requests table
CREATE TABLE requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hotel_id UUID NULL REFERENCES hotels(id) ON DELETE SET NULL,
    company_id UUID NULL REFERENCES companies(id) ON DELETE SET NULL,
    fecha TIMESTAMPTZ NOT NULL,
    origin JSONB NOT NULL,
    destination JSONB NOT NULL,
    pax INTEGER NOT NULL CHECK (pax > 0),
    vehicle_type vehicle_type NOT NULL,
    language language NULL,
    status request_status NOT NULL DEFAULT 'PENDIENTE',
    assigned_driver_id VARCHAR(20) NULL REFERENCES drivers(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT check_org_id CHECK ((hotel_id IS NOT NULL) OR (company_id IS NOT NULL))
);

-- Reservations table
CREATE TABLE reservations (
    id VARCHAR(20) PRIMARY KEY, -- e.g., RSV-1042
    user_id UUID NULL REFERENCES users(id) ON DELETE SET NULL,
    org_id UUID NULL, -- Can reference hotels or companies
    pickup VARCHAR(500) NOT NULL,
    destination VARCHAR(500) NOT NULL,
    datetime TIMESTAMPTZ NOT NULL,
    passengers INTEGER NOT NULL CHECK (passengers > 0),
    status reservation_status NOT NULL DEFAULT 'ACTIVA',
    amount NUMERIC(12,2) NULL,
    notes TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Reservation timeline table
CREATE TABLE reservation_timeline (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reservation_id VARCHAR(20) NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    at TIMESTAMPTZ NOT NULL,
    variant VARCHAR(50) NOT NULL DEFAULT 'default',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reservation_id VARCHAR(20) NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
    gateway payment_gateway NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'CLP',
    status payment_status NOT NULL DEFAULT 'PENDING',
    transaction_ref VARCHAR(255) NULL,
    payload JSONB NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Driver feedback table
CREATE TABLE driver_feedback (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    trip_id VARCHAR(20) NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
    passenger_name VARCHAR(255) NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add foreign key constraints
ALTER TABLE users ADD CONSTRAINT fk_users_org_id 
    FOREIGN KEY (org_id) REFERENCES companies(id) ON DELETE SET NULL;

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_org_id ON users(org_id);
CREATE INDEX idx_companies_rut ON companies(rut);
CREATE INDEX idx_drivers_status ON drivers(status);
CREATE INDEX idx_requests_status ON requests(status);
CREATE INDEX idx_requests_fecha ON requests(fecha);
CREATE INDEX idx_requests_assigned_driver ON requests(assigned_driver_id);
CREATE INDEX idx_reservations_status ON reservations(status);
CREATE INDEX idx_reservations_datetime ON reservations(datetime);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_reservation_id ON payments(reservation_id);

-- Create updated_at triggers
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_hotels_updated_at BEFORE UPDATE ON hotels 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_drivers_updated_at BEFORE UPDATE ON drivers 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vehicles_updated_at BEFORE UPDATE ON vehicles 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_requests_updated_at BEFORE UPDATE ON requests 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reservations_updated_at BEFORE UPDATE ON reservations 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_driver_availability_updated_at BEFORE UPDATE ON driver_availability 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

