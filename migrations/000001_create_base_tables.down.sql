-- Drop triggers
DROP TRIGGER IF EXISTS update_driver_availability_updated_at ON driver_availability;
DROP TRIGGER IF EXISTS update_reservations_updated_at ON reservations;
DROP TRIGGER IF EXISTS update_requests_updated_at ON requests;
DROP TRIGGER IF EXISTS update_vehicles_updated_at ON vehicles;
DROP TRIGGER IF EXISTS update_drivers_updated_at ON drivers;
DROP TRIGGER IF EXISTS update_hotels_updated_at ON hotels;
DROP TRIGGER IF EXISTS update_companies_updated_at ON companies;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS driver_feedback;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS reservation_timeline;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS driver_availability;
DROP TABLE IF EXISTS vehicle_photos;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS driver_background_checks;
DROP TABLE IF EXISTS driver_licenses;
DROP TABLE IF EXISTS drivers;
DROP TABLE IF EXISTS hotels;
DROP TABLE IF EXISTS companies;
DROP TABLE IF EXISTS users;

-- Drop ENUM types
DROP TYPE IF EXISTS language;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_gateway;
DROP TYPE IF EXISTS reservation_status;
DROP TYPE IF EXISTS request_status;
DROP TYPE IF EXISTS vehicle_type;
DROP TYPE IF EXISTS background_check_status;
DROP TYPE IF EXISTS license_class;
DROP TYPE IF EXISTS driver_status;
DROP TYPE IF EXISTS company_sector;
DROP TYPE IF EXISTS company_status;
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role;

