CREATE TABLE repository_spam_users (
	id UUID PRIMARY KEY,
	repository_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
	github_id BIGINT,
	username VARCHAR(60),
	email VARCHAR(60),
	meta JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (repository_id, github_id)
);
CREATE INDEX idx_repository_spam_users_repository_id ON repository_spam_users(repository_id);
