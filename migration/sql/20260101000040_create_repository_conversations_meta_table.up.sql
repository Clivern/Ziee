CREATE TABLE repository_conversations_meta (
	id UUID PRIMARY KEY,
	conversation_id UUID NOT NULL REFERENCES repository_conversations(id) ON DELETE CASCADE,
	key VARCHAR(60) NOT NULL,
	value JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (conversation_id, key)
);
