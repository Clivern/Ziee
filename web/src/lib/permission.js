import { computed } from 'vue'
import { loadUserFromStorage, loadWorkspaceFromStorage } from '@/utils/storage'

export const WORKSPACE_KEY_PERMISSIONS = [
  { value: 'CAN_LIST_WORKSPACE_DOCUMENTS', labelKey: 'workspace_settings_page.perm_list_documents' },
  { value: 'CAN_QUERY_WORKSPACE_DOCUMENTS', labelKey: 'workspace_settings_page.perm_query_documents' },
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
