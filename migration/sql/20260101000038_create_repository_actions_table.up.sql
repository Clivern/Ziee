CREATE TABLE repository_actions (
	id UUID PRIMARY KEY,
	repository_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
	async_task_id UUID NULL REFERENCES async_tasks(id) ON DELETE SET NULL,
	delivery_id VARCHAR(255) NOT NULL,
	event_kind VARCHAR(64) NOT NULL,
	issue_number INT NOT NULL,
	verb VARCHAR(64) NULL,
	plan JSONB NOT NULL DEFAULT '[]',
	status VARCHAR(20) NOT NULL,
	error TEXT NULL,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);

CREATE INDEX idx_repository_actions_repo_issue
	ON repository_actions (repository_id, issue_number, created_at DESC);
