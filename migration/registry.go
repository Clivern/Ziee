// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package migration

import (
	"database/sql"
)

// GetAll returns every registered migration.
func GetAll() []Migration {
	return []Migration{
		{
			Version:     "20260101000001",
			Description: "Create configs table",
			Up:          createConfigsTable,
			Down:        dropConfigsTable,
		},
		{
			Version:     "20260101000002",
			Description: "Create users table",
			Up:          createUsersTable,
			Down:        dropUsersTable,
		},
		{
			Version:     "20260101000003",
			Description: "Create users_meta table",
			Up:          createUsersMetaTable,
			Down:        dropUsersMetaTable,
		},
		{
			Version:     "20260101000004",
			Description: "Create user_sessions table",
			Up:          createUserSessionsTable,
			Down:        dropUserSessionsTable,
		}, {
			Version:     "20260101000005",
			Description: "Create user_api_keys table",
			Up:          createUserAPIKeysTable,
			Down:        dropUserAPIKeysTable,
		},
		{
			Version:     "20260101000006",
			Description: "Create workspaces table",
			Up:          createWorkspacesTable,
			Down:        dropWorkspacesTable,
		},
		{
			Version:     "20260101000007",
			Description: "Create user_invites table",
			Up:          createUserInvitesTable,
			Down:        dropUserInvitesTable,
		},
		{
			Version:     "20260101000008",
			Description: "Create workspaces_meta table",
			Up:          createWorkspacesMetaTable,
			Down:        dropWorkspacesMetaTable,
		},
		{
			Version:     "20260101000009",
			Description: "Create workspace_users table",
			Up:          createWorkspaceUsersTable,
			Down:        dropWorkspaceUsersTable,
		},
		{
			Version:     "20260101000010",
			Description: "Create password_reset_tokens table",
			Up:          createPasswordResetTokensTable,
			Down:        dropPasswordResetTokensTable,
		},
		{
			Version:     "20260101000011",
			Description: "Create integrations table",
			Up:          createIntegrationsTable,
			Down:        dropIntegrationsTable,
		},
		{
			Version:     "20260101000012",
			Description: "Create integrations_meta table",
			Up:          createIntegrationsMetaTable,
			Down:        dropIntegrationsMetaTable,
		},
		{
			Version:     "20260101000013",
			Description: "Create subscriptions table",
			Up:          createSubscriptionsTable,
			Down:        dropSubscriptionsTable,
		},
		{
			Version:     "20260101000014",
			Description: "Create subscriptions_meta table",
			Up:          createSubscriptionsMetaTable,
			Down:        dropSubscriptionsMetaTable,
		},
		{
			Version:     "20260101000015",
			Description: "Create usage table",
			Up:          createUsageTable,
			Down:        dropUsageTable,
		},
		{
			Version:     "20260101000016",
			Description: "Create audit table",
			Up:          createAuditTable,
			Down:        dropAuditTable,
		},
		{
			Version:     "20260101000017",
			Description: "Create prompts table",
			Up:          createPromptsTable,
			Down:        dropPromptsTable,
		},
		{
			Version:     "20260101000018",
			Description: "Create prompts_meta table",
			Up:          createPromptsMetaTable,
			Down:        dropPromptsMetaTable,
		},
		{
			Version:     "20260101000019",
			Description: "Create prompt_versions table",
			Up:          createPromptVersionsTable,
			Down:        dropPromptVersionsTable,
		},
		{
			Version:     "20260101000020",
			Description: "Create agents table",
			Up:          createAgentsTable,
			Down:        dropAgentsTable,
		},
		{
			Version:     "20260101000021",
			Description: "Create agents_meta table",
			Up:          createAgentsMetaTable,
			Down:        dropAgentsMetaTable,
		},
		{
			Version:     "20260101000022",
			Description: "Create sessions table",
			Up:          createSessionsTable,
			Down:        dropSessionsTable,
		},
		{
			Version:     "20260101000023",
			Description: "Create sessions_meta table",
			Up:          createSessionsMetaTable,
			Down:        dropSessionsMetaTable,
		},
		{
			Version:     "20260101000024",
			Description: "Create session_messages table",
			Up:          createSessionMessagesTable,
			Down:        dropSessionMessagesTable,
		},
		{
			Version:     "20260101000025",
			Description: "Create session_memories table",
			Up:          createSessionMemoriesTable,
			Down:        dropSessionMemoriesTable,
		},
		{
			Version:     "20260101000026",
			Description: "Create workspace_documents table",
			Up:          createWorkspaceDocumentsTable,
			Down:        dropWorkspaceDocumentsTable,
		},
		{
			Version:     "20260101000027",
			Description: "Create workspace_documents_meta table",
			Up:          createWorkspaceDocumentsMetaTable,
			Down:        dropWorkspaceDocumentsMetaTable,
		},
		{
			Version:     "20260101000028",
			Description: "Create async_tasks table",
			Up:          createAsyncTasksTable,
			Down:        dropAsyncTasksTable,
		},
		{
			Version:     "20260101000029",
			Description: "Create async_tasks_meta table",
			Up:          createAsyncTasksMetaTable,
			Down:        dropAsyncTasksMetaTable,
		},
		{
			Version:     "20260101000030",
			Description: "Create workspace_access_keys table",
			Up:          createWorkspaceAccessKeysTable,
			Down:        dropWorkspaceAccessKeysTable,
		},
	}
}

// createConfigsTable creates the configs table.
func createConfigsTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE configs (
			id UUID PRIMARY KEY,
			key VARCHAR(60) NOT NULL UNIQUE,
			value JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_configs_key ON configs(key)`)
}

// dropConfigsTable drops the configs table.
func dropConfigsTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS configs")
}

// createUsersTable creates the users table.
func createUsersTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE users (
			id UUID PRIMARY KEY,
			name VARCHAR(60) NOT NULL,
			email VARCHAR(60) NOT NULL UNIQUE,
			pwd_hash VARCHAR(200) NOT NULL,
			provider VARCHAR(20) NOT NULL DEFAULT 'local',
			provider_user_id VARCHAR(255),
			role VARCHAR(20) NOT NULL DEFAULT 'regular',
			is_active BOOLEAN DEFAULT true,
			is_email_verified BOOLEAN DEFAULT false,
			email_verify_token VARCHAR(100) NULL UNIQUE,
			last_login_at TIMESTAMP NULL,
			language VARCHAR(20) NOT NULL DEFAULT 'en',
			theme VARCHAR(20) NOT NULL DEFAULT 'default',
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_users_email ON users(email)`)
}

// dropUsersTable drops the users table.
func dropUsersTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS users")
}

// createUsersMetaTable creates the users meta table.
func createUsersMetaTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE users_meta (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			key VARCHAR(60) NOT NULL,
			value JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (user_id, key)
		);
		CREATE INDEX idx_users_meta_user_id_key ON users_meta(user_id, key)`)
}

// dropUsersMetaTable drops the users meta table.
func dropUsersMetaTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS users_meta")
}

// createUserSessionsTable creates the user sessions table.
func createUserSessionsTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE user_sessions (
			id UUID PRIMARY KEY,
			token VARCHAR(100) NOT NULL UNIQUE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			ip_address VARCHAR(45),
			user_agent VARCHAR(200),
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_user_sessions_token ON user_sessions(token)`)
}

// dropUserSessionsTable drops the user sessions table.
func dropUserSessionsTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS user_sessions")
}

// createUserAPIKeysTable creates the user apikeys table.
func createUserAPIKeysTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE user_api_keys (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(60) NOT NULL,
			token VARCHAR(100) NOT NULL UNIQUE,
			expires_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_user_api_keys_token ON user_api_keys(token)`)
}

// dropUserAPIKeysTable drops the user apikeys table.
func dropUserAPIKeysTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS user_api_keys")
}

// createUserInvitesTable creates the user invites table.
func createUserInvitesTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE user_invites (
			id UUID PRIMARY KEY,
			email VARCHAR(60) NOT NULL,
			role VARCHAR(20) NOT NULL DEFAULT 'regular',
			token VARCHAR(100) NOT NULL UNIQUE,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			inviter_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			expires_at TIMESTAMP NOT NULL,
			accepted_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_user_invites_token ON user_invites(token)`)
}

// dropUserInvitesTable drops the user invites table.
func dropUserInvitesTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS user_invites")
}

// createWorkspacesTable creates the workspaces table.
func createWorkspacesTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE workspaces (
			id UUID PRIMARY KEY,
			name VARCHAR(60) NOT NULL,
			handle VARCHAR(100) NOT NULL UNIQUE,
			meta JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_workspaces_handle ON workspaces(handle)`)
}

// dropWorkspacesTable drops the workspaces table.
func dropWorkspacesTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS workspaces")
}

// createWorkspacesMetaTable creates the workspaces meta table.
func createWorkspacesMetaTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE workspaces_meta (
			id UUID PRIMARY KEY,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			key VARCHAR(60) NOT NULL,
			value JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (workspace_id, key)
		)`)
}

// dropWorkspacesMetaTable drops the workspaces meta table.
func dropWorkspacesMetaTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS workspaces_meta")
}

// createWorkspaceUsersTable creates the workspace users table.
func createWorkspaceUsersTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE workspace_users (
			id UUID PRIMARY KEY,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			role VARCHAR(20) NOT NULL DEFAULT 'regular',
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (workspace_id, user_id)
		)`)
}

// dropWorkspaceUsersTable drops the workspace users table.
func dropWorkspaceUsersTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS workspace_users")
}

// createPasswordResetTokensTable creates the password reset tokens table.
func createPasswordResetTokensTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE password_reset_tokens (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token VARCHAR(100) NOT NULL UNIQUE,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_password_reset_tokens_token ON password_reset_tokens(token)`)
}

// dropPasswordResetTokensTable drops the password reset tokens table.
func dropPasswordResetTokensTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS password_reset_tokens")
}

// createIntegrationsTable creates the integrations table.
func createIntegrationsTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE integrations (
			id UUID PRIMARY KEY,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			type VARCHAR(20) NOT NULL,
			name VARCHAR(60) NOT NULL,
			config JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		)`)
}

// dropIntegrationsTable drops the integrations table.
func dropIntegrationsTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS integrations")
}

// createIntegrationsMetaTable creates the integrations meta table.
func createIntegrationsMetaTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE integrations_meta (
			id UUID PRIMARY KEY,
			integration_id UUID NOT NULL REFERENCES integrations(id) ON DELETE CASCADE,
			key VARCHAR(60) NOT NULL,
			value JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (integration_id, key)
		)`)
}

// dropIntegrationsMetaTable drops the integrations meta table.
func dropIntegrationsMetaTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS integrations_meta")
}

// createSubscriptionsTable creates the subscriptions table.
func createSubscriptionsTable(db *sql.DB) error {
	return exec(db, `
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
		)`)
}

// dropSubscriptionsTable drops the subscriptions table.
func dropSubscriptionsTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS subscriptions")
}

// createSubscriptionsMetaTable creates the subscriptions meta table.
func createSubscriptionsMetaTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE subscriptions_meta (
			id UUID PRIMARY KEY,
			subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
			key VARCHAR(60) NOT NULL,
			value JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (subscription_id, key)
		)`)
}

// dropSubscriptionsMetaTable drops the subscriptions meta table.
func dropSubscriptionsMetaTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS subscriptions_meta")
}

// createUsageTable creates the usage table.
func createUsageTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE usage (
			id UUID PRIMARY KEY,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			type VARCHAR(60) NOT NULL,
			quantity BIGINT NOT NULL DEFAULT 1,
			unit VARCHAR(20),
			period_start TIMESTAMP NOT NULL,
			period_end TIMESTAMP NOT NULL,
			meta JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (workspace_id, type, period_start)
		)`)
}

// dropUsageTable drops the usage table.
func dropUsageTable(db *sql.DB) error {
	return exec(db, `DROP TABLE IF EXISTS usage`)
}

// createAuditTable creates the audit table.
func createAuditTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE audit (
			id UUID PRIMARY KEY,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			user_id UUID NULL REFERENCES users(id),
			action VARCHAR(100) NOT NULL,
			resource_type VARCHAR(100),
			resource_id UUID,
			ip_address VARCHAR(45),
			user_agent VARCHAR(200),
			meta JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		)`)
}

// dropAuditTable drops the audit table.
func dropAuditTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS audit")
}

// createAgentsTable creates the agents table.
func createAgentsTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE agents (
			id UUID PRIMARY KEY,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			kind VARCHAR(50) NOT NULL DEFAULT 'unmanaged',
			prompt_id UUID NULL REFERENCES prompts(id) ON DELETE SET NULL,
			kb_labels JSONB NULL,
			handle VARCHAR(100) NOT NULL UNIQUE,
			description TEXT NOT NULL DEFAULT '',
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			config JSONB,
			meta JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_agents_workspace_id ON agents(workspace_id);
		CREATE INDEX idx_agents_status ON agents(status)`)
}

// dropAgentsTable drops the agents table.
func dropAgentsTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS agents")
}

// createAgentsMetaTable creates the agents meta table.
func createAgentsMetaTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE agents_meta (
			id UUID PRIMARY KEY,
			agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
			key VARCHAR(128) NOT NULL,
			value TEXT,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (agent_id, key)
		);
		CREATE INDEX idx_agents_meta_agent_id ON agents_meta(agent_id);
		CREATE INDEX idx_agents_meta_key ON agents_meta(key)`)
}

// dropAgentsMetaTable drops the agents meta table.
func dropAgentsMetaTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS agents_meta")
}

// createSessionsTable creates the sessions table.
func createSessionsTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE sessions (
			id UUID PRIMARY KEY,
			external_id VARCHAR(60) NOT NULL,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
			title VARCHAR(255),
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			labels JSONB,
			meta JSONB,
			last_activity_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (agent_id, external_id)
		);
		CREATE INDEX idx_sessions_workspace_id ON sessions(workspace_id);
		CREATE INDEX idx_sessions_agent_id ON sessions(agent_id);
		CREATE INDEX idx_sessions_status ON sessions(status);
		CREATE INDEX idx_sessions_labels ON sessions USING GIN (labels jsonb_path_ops)`)
}

// dropSessionsTable drops the sessions table.
func dropSessionsTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS sessions")
}

// createSessionsMetaTable creates the sessions meta table.
func createSessionsMetaTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE sessions_meta (
			id UUID PRIMARY KEY,
			session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
			key VARCHAR(128) NOT NULL,
			value TEXT,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (session_id, key)
		);
		CREATE INDEX idx_sessions_meta_session_id ON sessions_meta(session_id);
		CREATE INDEX idx_sessions_meta_key ON sessions_meta(key)`)
}

// dropSessionsMetaTable drops the sessions meta table.
func dropSessionsMetaTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS sessions_meta")
}

// createSessionMessagesTable creates the session messages table.
func createSessionMessagesTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE session_messages (
			id UUID PRIMARY KEY,
			internal_id UUID NOT NULL UNIQUE,
			session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
			role VARCHAR(50) NOT NULL,
			content TEXT NOT NULL,
			meta JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_session_messages_session_id ON session_messages(session_id);
		CREATE INDEX idx_session_messages_role ON session_messages(role);
		CREATE INDEX idx_session_messages_created_at ON session_messages(created_at)`)
}

// dropSessionMessagesTable drops the session messages table.
func dropSessionMessagesTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS session_messages")
}

// createSessionMemoriesTable creates the session memories table.
func createSessionMemoriesTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE session_memories (
			id UUID PRIMARY KEY,
			internal_id UUID NOT NULL UNIQUE,
			session_id UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
			kind VARCHAR(100) NOT NULL DEFAULT 'summary',
			content TEXT NOT NULL,
			meta JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_session_memories_session_id ON session_memories(session_id);
		CREATE INDEX idx_session_memories_kind ON session_memories(kind)`)
}

// dropSessionMemoriesTable drops the session memories table.
func dropSessionMemoriesTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS session_memories")
}

// createWorkspaceDocumentsTable creates the workspace documents table.
func createWorkspaceDocumentsTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE workspace_documents (
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
		CREATE INDEX idx_workspace_documents_workspace_id ON workspace_documents(workspace_id);
		CREATE INDEX idx_workspace_documents_status ON workspace_documents(status);
		CREATE INDEX idx_workspace_documents_labels ON workspace_documents USING GIN (labels jsonb_path_ops)`)
}

// dropWorkspaceDocumentsTable drops the workspace documents table.
func dropWorkspaceDocumentsTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS workspace_documents")
}

// createWorkspaceDocumentsMetaTable creates the workspace documents meta table.
func createWorkspaceDocumentsMetaTable(db *sql.DB) error {
	return exec(db, `
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
		CREATE INDEX idx_workspace_documents_meta_key ON workspace_documents_meta(key)`)
}

// dropWorkspaceDocumentsMetaTable drops the workspace documents meta table.
func dropWorkspaceDocumentsMetaTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS workspace_documents_meta")
}

// createAsyncTasksTable creates the async tasks table.
func createAsyncTasksTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE async_tasks (
			id UUID PRIMARY KEY,
			workspace_id UUID NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			type VARCHAR(60) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			payload JSONB,
			result JSONB,
			error JSONB,
			attempts INTEGER NOT NULL DEFAULT 0,
			priority SMALLINT NOT NULL DEFAULT 50,
			run_at TIMESTAMP NULL,
			locked_at TIMESTAMP NULL,
			completed_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_async_tasks_status_priority_created_at
			ON async_tasks (status, priority DESC, created_at)`)
}

// dropAsyncTasksTable drops the async tasks table.
func dropAsyncTasksTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS async_tasks")
}

// createAsyncTasksMetaTable creates the async tasks meta table.
func createAsyncTasksMetaTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE async_tasks_meta (
			id UUID PRIMARY KEY,
			async_task_id UUID NOT NULL REFERENCES async_tasks(id) ON DELETE CASCADE,
			key VARCHAR(60) NOT NULL,
			value JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (async_task_id, key)
		)`)
}

// dropAsyncTasksMetaTable drops the async tasks meta table.
func dropAsyncTasksMetaTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS async_tasks_meta")
}

// createPromptsTable creates the prompts table.
func createPromptsTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE prompts (
			id UUID PRIMARY KEY,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			handle VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (workspace_id, handle)
		);
		CREATE INDEX idx_prompts_workspace_id ON prompts(workspace_id)`)
}

// dropPromptsTable drops the prompts table.
func dropPromptsTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS prompts")
}

// createPromptsMetaTable creates the prompts meta table.
func createPromptsMetaTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE prompts_meta (
			id UUID PRIMARY KEY,
			prompt_id UUID NOT NULL REFERENCES prompts(id) ON DELETE CASCADE,
			key VARCHAR(128) NOT NULL,
			value TEXT,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (prompt_id, key)
		);
		CREATE INDEX idx_prompts_meta_prompt_id ON prompts_meta(prompt_id);
		CREATE INDEX idx_prompts_meta_key ON prompts_meta(key)`)
}

// dropPromptsMetaTable drops the prompts meta table.
func dropPromptsMetaTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS prompts_meta")
}

// createPromptVersionsTable creates the prompt versions table.
func createPromptVersionsTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE prompt_versions (
			id UUID PRIMARY KEY,
			prompt_id UUID NOT NULL REFERENCES prompts(id) ON DELETE CASCADE,
			version INTEGER NOT NULL,
			type VARCHAR(20) NOT NULL DEFAULT 'text',
			content TEXT NOT NULL,
			config JSONB,
			labels JSONB,
			commit_message VARCHAR(255),
			commit_hash VARCHAR(40) NOT NULL,
			meta JSONB,
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			production BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			UNIQUE (prompt_id, version)
		);
		CREATE INDEX idx_prompt_versions_prompt_id ON prompt_versions(prompt_id);
		CREATE INDEX idx_prompt_versions_status ON prompt_versions(status);
		CREATE INDEX idx_prompt_versions_labels ON prompt_versions USING GIN (labels);
		CREATE UNIQUE INDEX idx_prompt_versions_one_production ON prompt_versions(prompt_id) WHERE production = TRUE`)
}

// dropPromptVersionsTable drops the prompt versions table.
func dropPromptVersionsTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS prompt_versions")
}

// createWorkspaceAccessKeysTable creates the workspace access keys table.
func createWorkspaceAccessKeysTable(db *sql.DB) error {
	return exec(db, `
		CREATE TABLE workspace_access_keys (
			id UUID PRIMARY KEY,
			workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
			name VARCHAR(60) NOT NULL,
			token VARCHAR(100) NOT NULL UNIQUE,
			expires_at TIMESTAMP NULL,
			meta JSONB,
			created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
			updated_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
		);
		CREATE INDEX idx_workspace_access_keys_workspace_id ON workspace_access_keys(workspace_id);
		CREATE INDEX idx_workspace_access_keys_token ON workspace_access_keys(token)`)
}

// dropWorkspaceAccessKeysTable drops the workspace access keys table.
func dropWorkspaceAccessKeysTable(db *sql.DB) error {
	return exec(db, "DROP TABLE IF EXISTS workspace_access_keys")
}
