-- Add company_profile column to users table
-- This field is only used when role = 'COMPANY'
-- Valid values: 'COMPANY_ADMIN', 'COMPANY_USER'

ALTER TABLE users 
ADD COLUMN company_profile VARCHAR(50);

-- Add comment to explain usage
COMMENT ON COLUMN users.company_profile IS 'Profile type for COMPANY role users. Valid values: COMPANY_ADMIN (full access), COMPANY_USER (limited access)';

-- Set default value for existing COMPANY users (backwards compatibility)
UPDATE users 
SET company_profile = 'COMPANY_ADMIN' 
WHERE role = 'COMPANY' AND company_profile IS NULL;

-- Create index for filtering by company_profile
CREATE INDEX idx_users_company_profile ON users(company_profile) WHERE company_profile IS NOT NULL;

