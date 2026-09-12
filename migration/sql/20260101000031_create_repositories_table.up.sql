CREATE TABLE repositories (
	id UUID PRIMARY KEY,
	workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
	installation_id BIGINT NOT NULL,
	github_id BIGINT NOT NULL UNIQUE,
	node_id VARCHAR(64) NOT NULL,
	owner VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	full_name VARCHAR(255) NOT NULL,
	private BOOLEAN NOT NULL DEFAULT FALSE,
	meta JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);
CREATE INDEX idx_repositories_workspace_id ON repositories(workspace_id);
CREATE INDEX idx_repositories_installation_id ON repositories(installation_id);
