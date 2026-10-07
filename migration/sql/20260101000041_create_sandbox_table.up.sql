CREATE TABLE sandbox (
	id UUID PRIMARY KEY,
	repository_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
	owner VARCHAR(255) NOT NULL,
	config JSONB NOT NULL DEFAULT '{}',
	usage JSONB NOT NULL DEFAULT '{}',
	port INT NOT NULL,
	status VARCHAR(20) NOT NULL DEFAULT 'active',
	token VARCHAR(100) NOT NULL UNIQUE,
	remote_id VARCHAR(200) NOT NULL,
	run_id UUID NOT NULL UNIQUE,
	expires_at TIMESTAMP NOT NULL,
	last_activity_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);
CREATE INDEX idx_sandbox_repository_id ON sandbox(repository_id);
CREATE INDEX idx_sandbox_owner ON sandbox(owner);
CREATE INDEX idx_sandbox_expires_at ON sandbox(expires_at);
CREATE INDEX idx_sandbox_status ON sandbox(status);
CREATE INDEX idx_sandbox_token ON sandbox(token);
CREATE INDEX idx_sandbox_remote_id ON sandbox(remote_id);
CREATE INDEX idx_sandbox_run_id ON sandbox(run_id);
CREATE INDEX idx_sandbox_last_activity_at ON sandbox(last_activity_at);
