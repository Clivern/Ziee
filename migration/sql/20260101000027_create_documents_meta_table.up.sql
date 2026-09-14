CREATE TABLE documents_meta (
	id UUID PRIMARY KEY,
	document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
	key VARCHAR(128) NOT NULL,
	value TEXT,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (document_id, key)
);
CREATE INDEX idx_documents_meta_document_id ON documents_meta(document_id);
CREATE INDEX idx_documents_meta_key ON documents_meta(key);
