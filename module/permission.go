// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/clivern/actx0/db"

	"github.com/samber/lo"
)

// Workspace permission constants. Use with Perm.Can for workspace-scoped checks.
const (
	CanGetWorkspace    = "CAN_GET_WORKSPACE"
	CanUpdateWorkspace = "CAN_UPDATE_WORKSPACE"
	CanDeleteWorkspace = "CAN_DELETE_WORKSPACE"

	CanListWorkspaceMembers = "CAN_LIST_WORKSPACE_MEMBERS"
	CanUpdateMemberRole     = "CAN_UPDATE_MEMBER_ROLE"
	CanRemoveMember         = "CAN_REMOVE_MEMBER"

	CanListWorkspaceInvites = "CAN_LIST_WORKSPACE_INVITES"
	CanGetWorkspaceInvite   = "CAN_GET_WORKSPACE_INVITE"
	CanInviteMember         = "CAN_INVITE_MEMBER"
	CanRemoveInvite         = "CAN_REMOVE_INVITE"

	CanGetWorkspaceBilling    = "CAN_GET_WORKSPACE_BILLING"
	CanUpdateWorkspaceBilling = "CAN_UPDATE_WORKSPACE_BILLING"

	CanListWorkspaceAccessKeys  = "CAN_LIST_WORKSPACE_ACCESS_KEYS"
	CanCreateWorkspaceAccessKey = "CAN_CREATE_WORKSPACE_ACCESS_KEY"
	CanGetWorkspaceAccessKey    = "CAN_GET_WORKSPACE_ACCESS_KEY"
	CanDeleteWorkspaceAccessKey = "CAN_DELETE_WORKSPACE_ACCESS_KEY"

	CanListPrompts  = "CAN_LIST_PROMPTS"
	CanGetPrompt    = "CAN_GET_PROMPT"
	CanCreatePrompt = "CAN_CREATE_PROMPT"
	CanUpdatePrompt = "CAN_UPDATE_PROMPT"
	CanDeletePrompt = "CAN_DELETE_PROMPT"

	CanCreateWorkspaceDocument = "CAN_CREATE_WORKSPACE_DOCUMENT"
	CanListWorkspaceDocuments  = "CAN_LIST_WORKSPACE_DOCUMENTS"
	CanQueryWorkspaceDocuments = "CAN_QUERY_WORKSPACE_DOCUMENTS"
	CanDeleteWorkspaceDocument = "CAN_DELETE_WORKSPACE_DOCUMENT"

	CanListAgents  = "CAN_LIST_AGENTS"
	CanGetAgent    = "CAN_GET_AGENT"
	CanCreateAgent = "CAN_CREATE_AGENT"
	CanUpdateAgent = "CAN_UPDATE_AGENT"
	CanDeleteAgent = "CAN_DELETE_AGENT"

	CanListAgentSessions  = "CAN_LIST_AGENT_SESSIONS"
	CanGetAgentSession    = "CAN_GET_AGENT_SESSION"
	CanCreateAgentSession = "CAN_CREATE_AGENT_SESSION"
	CanUpdateAgentSession = "CAN_UPDATE_AGENT_SESSION"
	CanDeleteAgentSession = "CAN_DELETE_AGENT_SESSION"

	CanListSessionMessages  = "CAN_LIST_SESSION_MESSAGES"
	CanGetSessionMessage    = "CAN_GET_SESSION_MESSAGE"
	CanCreateSessionMessage = "CAN_CREATE_SESSION_MESSAGE"
	CanUpdateSessionMessage = "CAN_UPDATE_SESSION_MESSAGE"
	CanDeleteSessionMessage = "CAN_DELETE_SESSION_MESSAGE"

	CanListSessionMemories = "CAN_LIST_SESSION_MEMORIES"
	CanGetSessionMemory    = "CAN_GET_SESSION_MEMORY"
	CanCreateSessionMemory = "CAN_CREATE_SESSION_MEMORY"
	CanUpdateSessionMemory = "CAN_UPDATE_SESSION_MEMORY"
	CanDeleteSessionMemory = "CAN_DELETE_SESSION_MEMORY"

	CanListWorkspaceAudits = "CAN_LIST_WORKSPACE_AUDITS"
	CanGetWorkspaceAudit   = "CAN_GET_WORKSPACE_AUDIT"
)

var (
	ErrPermissionPrincipalMissing = errors.New("permission principal missing")
	ErrPermissionWorkspaceMissing = errors.New("permission workspace missing")
)

// AccessKeyPermissions lists permissions that may be assigned to workspace access keys.
var AccessKeyPermissions = map[string]bool{
	CanListPrompts:             true,
	CanGetPrompt:               true,
	CanListWorkspaceDocuments:  true,
	CanQueryWorkspaceDocuments: true,
	CanListAgents:              true,
	CanGetAgent:                true,
	CanCreateAgent:             true,
	CanDeleteAgent:             true,
	CanListAgentSessions:       true,
	CanGetAgentSession:         true,
	CanCreateAgentSession:      true,
	CanUpdateAgentSession:      true,
	CanDeleteAgentSession:      true,
	CanListSessionMessages:     true,
	CanGetSessionMessage:       true,
	CanCreateSessionMessage:    true,
	CanUpdateSessionMessage:    true,
	CanDeleteSessionMessage:    true,
	CanListSessionMemories:     true,
	CanGetSessionMemory:        true,
	CanCreateSessionMemory:     true,
	CanUpdateSessionMemory:     true,
	CanDeleteSessionMemory:     true,
}

// Perm checks workspace permissions for users and workspace access keys.
type Perm struct {
	WorkspaceRepository db.WorkspaceRepository
	user                *db.User
	accessKey           *db.WorkspaceAccessKey
	workspace           *db.Workspace
}

// NewPerm creates a permission checker backed by the given database connection.
func NewPerm(dbConn *sql.DB) *Perm {
	return &Perm{
		WorkspaceRepository: db.NewWorkspaceRepository(dbConn),
	}
}

// WithUser sets the user principal. Mutually exclusive with WithAccessKey.
func (p *Perm) WithUser(user *db.User) *Perm {
	p.user = user
	p.accessKey = nil
	return p
}

// WithAccessKey sets the workspace access key principal. Mutually exclusive with WithUser.
func (p *Perm) WithAccessKey(accessKey *db.WorkspaceAccessKey) *Perm {
	p.accessKey = accessKey
	p.user = nil
	return p
}

// WithWorkspace sets the workspace being accessed.
func (p *Perm) WithWorkspace(workspace *db.Workspace) *Perm {
	p.workspace = workspace
	return p
}

// WithWorkspaceId loads the workspace by id from the repository.
func (p *Perm) WithWorkspaceId(workspaceId db.Id) *Perm {
	workspace, _ := p.WorkspaceRepository.GetById(workspaceId)
	p.workspace = workspace
	return p
}

// Can reports whether the configured principal may perform the given permission on the workspace.
func (p *Perm) Can(permission string) (bool, error) {
	if p.workspace == nil {
		return false, ErrPermissionWorkspaceMissing
	}
	switch {
	case p.user != nil:
		return p.CanAsUser(permission)
	case p.accessKey != nil:
		return p.CanAsAccessKey(permission)
	default:
		return false, ErrPermissionPrincipalMissing
	}
}

// CanAsUser checks if the user has the given permission.
func (p *Perm) CanAsUser(permission string) (bool, error) {
	membership, err := p.WorkspaceRepository.GetWorkspaceMembership(p.workspace.Id, p.user.Id)
	if err != nil {
		return false, err
	}
	if membership == nil {
		return false, nil
	}

	switch permission {
	case CanGetWorkspace,
		CanGetWorkspaceBilling,
		CanListPrompts,
		CanGetPrompt,
		CanListWorkspaceDocuments,
		CanQueryWorkspaceDocuments,
		CanListAgents,
		CanGetAgent,
		CanListAgentSessions,
		CanGetAgentSession,
		CanListSessionMessages,
		CanGetSessionMessage,
		CanListSessionMemories,
		CanGetSessionMemory:
		return true, nil
	case CanUpdateWorkspace,
		CanDeleteWorkspace,
		CanUpdateWorkspaceBilling,
		CanListWorkspaceMembers,
		CanUpdateMemberRole,
		CanRemoveMember,
		CanListWorkspaceInvites,
		CanGetWorkspaceInvite,
		CanInviteMember,
		CanRemoveInvite,
		CanListWorkspaceAccessKeys,
		CanCreateWorkspaceAccessKey,
		CanGetWorkspaceAccessKey,
		CanDeleteWorkspaceAccessKey,
		CanListWorkspaceAudits,
		CanGetWorkspaceAudit,
		CanCreatePrompt,
		CanUpdatePrompt,
		CanDeletePrompt,
		CanCreateWorkspaceDocument,
		CanDeleteWorkspaceDocument,
		CanCreateAgent,
		CanUpdateAgent,
		CanDeleteAgent,
		CanCreateAgentSession,
		CanUpdateAgentSession,
		CanDeleteAgentSession,
		CanCreateSessionMessage,
		CanUpdateSessionMessage,
		CanDeleteSessionMessage,
		CanCreateSessionMemory,
		CanUpdateSessionMemory,
		CanDeleteSessionMemory:
		return membership.Role == db.UserRoleAdmin || membership.Role == db.UserRoleOwner, nil
	default:
		return false, nil
	}
}

// CanAsAccessKey checks if the access key has the given permission.
func (p *Perm) CanAsAccessKey(permission string) (bool, error) {
	if p.accessKey.WorkspaceId != p.workspace.Id {
		return false, nil
	}
	if !AccessKeyPermissions[permission] {
		return false, nil
	}

	var meta AccessKeyMeta
	err := json.Unmarshal([]byte(*p.accessKey.Meta), &meta)
	if err != nil {
		return false, err
	}

	return lo.Contains(meta.Permissions, permission), nil
}

// ValidAccessKeyPermissions reports whether every permission may be assigned to an access key.
func ValidAccessKeyPermissions(permissions []string) bool {
	if len(permissions) == 0 {
		return false
	}
	return lo.EveryBy(permissions, func(permission string) bool {
		return AccessKeyPermissions[permission]
	})
}
