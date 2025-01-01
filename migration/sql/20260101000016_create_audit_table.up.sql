CREATE TABLE audit (
	id UUID PRIMARY KEY,
	workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
	actor_id UUID NULL,
	actor_name VARCHAR(100) NULL,
	actor_type VARCHAR(20) NULL,
	action VARCHAR(100) NOT NULL,
	resource_type VARCHAR(100),
	resource_id UUID,
	ip_address VARCHAR(45),
	user_agent VARCHAR(200),
	meta JSONB,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
);
