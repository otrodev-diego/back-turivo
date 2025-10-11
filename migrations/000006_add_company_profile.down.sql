-- Remove company_profile column from users table

DROP INDEX IF EXISTS idx_users_company_profile;

ALTER TABLE users 
DROP COLUMN IF EXISTS company_profile;


