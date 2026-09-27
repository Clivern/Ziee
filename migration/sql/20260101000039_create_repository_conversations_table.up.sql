CREATE TABLE repository_conversations (
	id UUID PRIMARY KEY,
	repository_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
	target_type VARCHAR(20) NOT NULL,
	target_number INT NOT NULL,
	role VARCHAR(20) NOT NULL,
	author VARCHAR(255) NULL,
	comment_id BIGINT NULL,
	body TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);

CREATE INDEX idx_repository_conversations_target
	ON repository_conversations (repository_id, target_type, target_number, created_at);
