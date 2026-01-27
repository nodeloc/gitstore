-- Make GitHub-related fields nullable in plugins table
ALTER TABLE plugins ALTER COLUMN github_repo_id DROP NOT NULL;
ALTER TABLE plugins ALTER COLUMN github_repo_url DROP NOT NULL;
ALTER TABLE plugins ALTER COLUMN github_repo_name DROP NOT NULL;

-- Drop unique constraint on github_repo_id since it can now be NULL
-- and we don't want multiple NULL values to be considered violations
ALTER TABLE plugins DROP CONSTRAINT IF EXISTS plugins_github_repo_id_key;

-- Add a conditional unique constraint that only applies to non-NULL values
CREATE UNIQUE INDEX plugins_github_repo_id_unique 
ON plugins (github_repo_id) 
WHERE github_repo_id IS NOT NULL AND github_repo_id != 0;
