import { computed } from 'vue'
import { loadUserFromStorage, loadWorkspaceFromStorage } from '@/utils/storage'

export const WORKSPACE_KEY_PERMISSIONS = [
  { value: 'CAN_LIST_PROMPTS', labelKey: 'workspace_settings_page.perm_list_prompts' },
  { value: 'CAN_GET_PROMPT', labelKey: 'workspace_settings_page.perm_get_prompt' },
  { value: 'CAN_LIST_WORKSPACE_DOCUMENTS', labelKey: 'workspace_settings_page.perm_list_documents' },
  { value: 'CAN_QUERY_WORKSPACE_DOCUMENTS', labelKey: 'workspace_settings_page.perm_query_documents' },
  { value: 'CAN_LIST_AGENTS', labelKey: 'workspace_settings_page.perm_list_agents' },
  { value: 'CAN_GET_AGENT', labelKey: 'workspace_settings_page.perm_get_agent' },
  { value: 'CAN_CREATE_AGENT', labelKey: 'workspace_settings_page.perm_create_agent' },
  { value: 'CAN_DELETE_AGENT', labelKey: 'workspace_settings_page.perm_delete_agent' },
  { value: 'CAN_LIST_AGENT_SESSIONS', labelKey: 'workspace_settings_page.perm_list_agent_sessions' },
  { value: 'CAN_GET_AGENT_SESSION', labelKey: 'workspace_settings_page.perm_get_agent_session' },
  { value: 'CAN_CREATE_AGENT_SESSION', labelKey: 'workspace_settings_page.perm_create_agent_session' },
  { value: 'CAN_UPDATE_AGENT_SESSION', labelKey: 'workspace_settings_page.perm_update_agent_session' },
  { value: 'CAN_DELETE_AGENT_SESSION', labelKey: 'workspace_settings_page.perm_delete_agent_session' },
  { value: 'CAN_LIST_SESSION_MESSAGES', labelKey: 'workspace_settings_page.perm_list_session_messages' },
  { value: 'CAN_GET_SESSION_MESSAGE', labelKey: 'workspace_settings_page.perm_get_session_message' },
  { value: 'CAN_CREATE_SESSION_MESSAGE', labelKey: 'workspace_settings_page.perm_create_session_message' },
  { value: 'CAN_UPDATE_SESSION_MESSAGE', labelKey: 'workspace_settings_page.perm_update_session_message' },
  { value: 'CAN_DELETE_SESSION_MESSAGE', labelKey: 'workspace_settings_page.perm_delete_session_message' },
  { value: 'CAN_LIST_SESSION_MEMORIES', labelKey: 'workspace_settings_page.perm_list_session_memories' },
  { value: 'CAN_GET_SESSION_MEMORY', labelKey: 'workspace_settings_page.perm_get_session_memory' },
  { value: 'CAN_CREATE_SESSION_MEMORY', labelKey: 'workspace_settings_page.perm_create_session_memory' },
  { value: 'CAN_UPDATE_SESSION_MEMORY', labelKey: 'workspace_settings_page.perm_update_session_memory' },
  { value: 'CAN_DELETE_SESSION_MEMORY', labelKey: 'workspace_settings_page.perm_delete_session_memory' },
]

export const PLATFORM_ROLE_ADMIN = 'admin'
export const PLATFORM_ROLE_REGULAR = 'regular'
export const WORKSPACE_ROLE_READONLY = 'readonly'
export const WORKSPACE_ROLE_ADMIN = 'admin'
export const WORKSPACE_ROLE_OWNER = 'owner'
export const WORKSPACE_ROLE_REGULAR = 'regular'

const PLATFORM_ROLES = new Set([
  PLATFORM_ROLE_ADMIN,
  PLATFORM_ROLE_REGULAR,
])

const WORKSPACE_MEMBERSHIP_ROLES = new Set([
  WORKSPACE_ROLE_READONLY,
  WORKSPACE_ROLE_REGULAR,
  WORKSPACE_ROLE_ADMIN,
  WORKSPACE_ROLE_OWNER,
])

export function getPlatformRole(user) {
  const r = user?.role
  if (r == null || r === '') return null
  return PLATFORM_ROLES.has(r) ? r : null
}

export function getWorkspaceMembershipRole(workspace) {
  const r = workspace?.role
  if (r == null || r === '') return null
  return WORKSPACE_MEMBERSHIP_ROLES.has(r) ? r : null
}

export function canManageWorkspace(user, workspace) {
  if (!workspace?.id) return false
  return workspace.role === WORKSPACE_ROLE_ADMIN || workspace.role === WORKSPACE_ROLE_OWNER
}

export function useWorkspaceContext() {
  const currentUser = loadUserFromStorage()
  const currentWorkspace = loadWorkspaceFromStorage()
  const canManage = computed(() => canManageWorkspace(currentUser, currentWorkspace))

  return { currentUser, currentWorkspace, canManage }
}
