CREATE TABLE pqueue (
	id UUID PRIMARY KEY,
	repo_id UUID NOT NULL REFERENCES workspace_github_repos(id) ON DELETE CASCADE,
	github_pr_id BIGINT NOT NULL,
	priority VARCHAR(20) NOT NULL DEFAULT 'medium',
	rank INTEGER NOT NULL,
	status VARCHAR(20) NOT NULL,
	pr_checksum TEXT NOT NULL,
	stack_checksum TEXT NOT NULL,
	meta JSONB NOT NULL DEFAULT '{}',
	opened_at TIMESTAMP NULL,
	merged_at TIMESTAMP NULL,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (repo_id, github_pr_id)
);
CREATE INDEX idx_pqueue_repo_id_status_priority_rank
	ON pqueue (repo_id, status, priority, rank);
