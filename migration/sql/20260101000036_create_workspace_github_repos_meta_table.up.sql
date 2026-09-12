CREATE TABLE workspace_github_repos_meta (
	id UUID PRIMARY KEY,
	workspace_github_repo_id UUID NOT NULL REFERENCES workspace_github_repos(id) ON DELETE CASCADE,
	key VARCHAR(60) NOT NULL,
	value JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (workspace_github_repo_id, key)
);
