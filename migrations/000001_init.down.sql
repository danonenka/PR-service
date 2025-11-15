DROP INDEX IF EXISTS idx_reviewer_assignments_reviewer_id;
DROP INDEX IF EXISTS idx_reviewer_assignments_pr_id;
DROP INDEX IF EXISTS idx_pull_requests_status;
DROP INDEX IF EXISTS idx_pull_requests_author_id;
DROP INDEX IF EXISTS idx_users_is_active;
DROP INDEX IF EXISTS idx_users_team_id;

DROP TABLE IF EXISTS reviewer_assignments;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;