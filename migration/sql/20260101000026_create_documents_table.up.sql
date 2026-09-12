CREATE TABLE documents (
	id UUID PRIMARY KEY,
	internal_id UUID NOT NULL UNIQUE,
	workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
	title VARCHAR(255) NOT NULL,
	filename VARCHAR(255) NOT NULL,
	content_type VARCHAR(100) NOT NULL,
	checksum VARCHAR(255) NOT NULL,
	size BIGINT NOT NULL,
	char_count BIGINT NOT NULL,
	labels JSONB,
	processed_at TIMESTAMP NULL,
	chunking_config JSONB,
	meta JSONB,
	status VARCHAR(50) NOT NULL DEFAULT 'pending',
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);
CREATE INDEX idx_documents_workspace_id ON documents(workspace_id);
CREATE INDEX idx_documents_status ON documents(status);
CREATE INDEX idx_documents_labels ON documents USING GIN (labels jsonb_path_ops);
