-- Add company_profile column to registration_tokens table
ALTER TABLE registration_tokens 
ADD COLUMN IF NOT EXISTS company_profile VARCHAR(50);

COMMENT ON COLUMN registration_tokens.company_profile IS 'Profile type for COMPANY role registrations. Valid values: COMPANY_ADMIN, COMPANY_USER';

