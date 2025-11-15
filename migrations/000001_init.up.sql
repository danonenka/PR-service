CREATE TABLE IF NOT EXISTS teams (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    team_id VARCHAR(255) NOT NULL,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pull_requests (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN',
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reviewer_assignments (
    pr_id VARCHAR(255) NOT NULL,
    reviewer_id VARCHAR(255) NOT NULL,
    PRIMARY KEY (pr_id, reviewer_id),
    FOREIGN KEY (pr_id) REFERENCES pull_requests(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewer_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_users_team_id ON users(team_id);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_pull_requests_status ON pull_requests(status);
CREATE INDEX IF NOT EXISTS idx_reviewer_assignments_pr_id ON reviewer_assignments(pr_id);
CREATE INDEX IF NOT EXISTS idx_reviewer_assignments_reviewer_id ON reviewer_assignments(reviewer_id);

