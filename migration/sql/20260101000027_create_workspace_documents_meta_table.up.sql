CREATE TABLE workspace_documents_meta (
	id UUID PRIMARY KEY,
	workspace_document_id UUID NOT NULL REFERENCES workspace_documents(id) ON DELETE CASCADE,
	key VARCHAR(128) NOT NULL,
	value TEXT,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (workspace_document_id, key)
);
CREATE INDEX idx_workspace_documents_meta_document_id ON workspace_documents_meta(workspace_document_id);
CREATE INDEX idx_workspace_documents_meta_key ON workspace_documents_meta(key);
