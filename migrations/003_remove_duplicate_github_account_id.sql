-- migrations/003_remove_duplicate_github_account_id.sql
-- Remove duplicate github_account_id column from licenses table
-- The correct column name is git_hub_account_id (with underscores)

-- Drop foreign key constraints on github_account_id
ALTER TABLE licenses DROP CONSTRAINT IF EXISTS licenses_github_account_id_fkey;

-- Drop the duplicate column
ALTER TABLE licenses DROP COLUMN IF EXISTS github_account_id;

-- Recreate the unique constraint with the correct column name
ALTER TABLE licenses DROP CONSTRAINT IF EXISTS unique_user_plugin_license;
ALTER TABLE licenses ADD CONSTRAINT unique_user_plugin_license UNIQUE (user_id, plugin_id, git_hub_account_id);
