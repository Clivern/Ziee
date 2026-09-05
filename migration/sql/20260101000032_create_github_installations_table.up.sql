CREATE TABLE github_installations (
	id UUID PRIMARY KEY,
	github_id BIGINT NOT NULL UNIQUE,
	github_user_id VARCHAR(255) NOT NULL,
	account_id BIGINT NOT NULL,
	account_login VARCHAR(255) NOT NULL,
	account_type VARCHAR(20) NOT NULL,
	workspace_id UUID NULL REFERENCES workspaces(id) ON DELETE SET NULL,
	status VARCHAR(20) NOT NULL DEFAULT 'pending',
	repository_selection VARCHAR(20) NOT NULL,
	html_url VARCHAR(500) NOT NULL,
	meta JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);
CREATE INDEX idx_github_installations_github_user_id ON github_installations(github_user_id);
CREATE INDEX idx_github_installations_workspace_id ON github_installations(workspace_id);
CREATE INDEX idx_github_installations_status ON github_installations(status);
