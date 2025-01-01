CREATE TABLE subscriptions (
	id UUID PRIMARY KEY,
	workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
	plan VARCHAR(20) NOT NULL,
	status VARCHAR(20) NOT NULL DEFAULT 'active',
	provider VARCHAR(20),
	provider_subscription_id VARCHAR(200),
	provider_customer_id VARCHAR(200),
	period_start TIMESTAMP,
	period_end TIMESTAMP,
	created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
	UNIQUE (workspace_id)
);
