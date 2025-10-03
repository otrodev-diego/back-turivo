-- Remove company_profile column from registration_tokens table
ALTER TABLE registration_tokens 
DROP COLUMN IF EXISTS company_profile;

