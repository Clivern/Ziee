import axios from 'axios'
import { removeUserFromStorage, removeWorkspaceFromStorage } from '@/utils/storage'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

api.interceptors.request.use((config) => {
  try {
    const locale = typeof localStorage !== 'undefined' && localStorage.getItem('_flocale')
    if (locale === 'en' || locale === 'fr') {
      config.headers['Accept-Language'] = locale
    }
  } catch (_) {}
  return config
})

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    const isLoginEndpoint = error.config?.url === 'public/action/login'
    const isProfileEndpoint = error.config?.url === '/action/profile'
    const isLogoutEndpoint = error.config?.url === 'public/action/logout'
    const isForgotPasswordEndpoint = error.config?.url === 'public/action/forgot-password'
    const isResetPasswordEndpoint = error.config?.url === 'public/action/reset-password'
    const isOnPublicPage = window.location.pathname === '/' ||
                           window.location.pathname === '/status' ||
                           window.location.pathname === '/docs' ||
                           window.location.pathname === '/login' ||
                           window.location.pathname === '/setup' ||
                           window.location.pathname.startsWith('/forgot-password') ||
                           window.location.pathname.startsWith('/reset-password') ||
                           window.location.pathname === '/register'

    if (error.response?.status === 401 &&
        !isLoginEndpoint &&
        !isProfileEndpoint &&
        !isLogoutEndpoint &&
        !isForgotPasswordEndpoint &&
        !isResetPasswordEndpoint &&
        !isOnPublicPage) {
      removeUserFromStorage()
      removeWorkspaceFromStorage()
      const redirect = `${window.location.pathname}${window.location.search}`
      window.location.href = `/login?redirect=${encodeURIComponent(redirect)}`
    }
    return Promise.reject(error)
  }
)

// API endpoints
export const healthAPI = {
  check: () => api.get('public/_health'),
  ready: () => api.get('public/_ready'),
}

export const authAPI = {
  login: (data) => api.post('public/action/login', data),
  logout: () => api.post('public/action/logout'),
  getProfile: () => api.get('/action/profile'),
  updateProfile: (data) => api.put('/action/profile', data),
  forgotPassword: (data) => api.post('public/action/forgot-password', data),
  resetPassword: (data) => api.post('public/action/reset-password', data),
  register: (data) => api.post('public/action/register', data),
}

export const setupAPI = {
  install: (data) => api.post('public/action/setup', data),
  checkInstalled: () => api.get('public/action/setup/status'),
}

// API endpoints for settings
export const settingsAPI = {
  getSettings: () => api.get('/action/settings'),
  updateSettings: (data) => api.put('/action/settings', data),
}

// API endpoints for workspaces (projects)
export const workspaceAPI = {
  list: (params) => api.get('/workspaces', { params }),
  get: (id) => api.get(`/workspaces/${id}`),
  create: (data) => api.post('/workspaces', data),
  update: (id, data) => api.put(`/workspaces/${id}`, data),
  delete: (id) => api.delete(`/workspaces/${id}`),
}

export const inviteAPI = {
  getByToken: (token) => api.get(`/action/invite-by-token/${encodeURIComponent(token)}`),
  acceptByToken: (token) => api.post(`/action/accept-invite/${encodeURIComponent(token)}`),
  rejectByToken: (token) => api.post(`/action/reject-invite/${encodeURIComponent(token)}`),
}

export const workspaceMemberAPI = {
  list: (workspaceId, params) => api.get(`/workspaces/${workspaceId}/members`, { params }),
  updateRole: (workspaceId, userId, data) => api.put(`/workspaces/${workspaceId}/members/${userId}`, data),
  delete: (workspaceId, userId) => api.delete(`/workspaces/${workspaceId}/members/${userId}`),
}

export const workspaceInviteAPI = {
  list: (workspaceId, params) => api.get(`/workspaces/${workspaceId}/invites`, { params }),
  create: (workspaceId, data) => api.post(`/workspaces/${workspaceId}/invites`, data),
  delete: (workspaceId, id) => api.delete(`/workspaces/${workspaceId}/invites/${id}`),
}

export const statsAPI = {
  get: (workspaceId) => api.get(`/workspaces/${workspaceId}/stats`),
}

export const billingAPI = {
  status: (workspaceId) => api.get(`/workspaces/${workspaceId}/billing`),
  usage: (workspaceId) => api.get(`/workspaces/${workspaceId}/billing/usage`),
  checkout: (workspaceId, data) => api.post(`/workspaces/${workspaceId}/billing/checkout`, data),
  portal: (workspaceId) => api.post(`/workspaces/${workspaceId}/billing/portal`),
}

export const workspaceAccessKeyAPI = {
  list: (workspaceId, params) => api.get(`/workspaces/${workspaceId}/keys`, { params }),
  create: (workspaceId, data) => api.post(`/workspaces/${workspaceId}/keys`, data),
  delete: (workspaceId, id) => api.delete(`/workspaces/${workspaceId}/keys/${id}`),
}

// API endpoints for current user's API keys
export const apiKeysAPI = {
  list: (params) => api.get('/apiKeys', { params }),
  create: (data) => api.post('/apiKeys', data),
  delete: (id) => api.delete(`/apiKeys/${id}`),
}

export const promptAPI = {
  list: (workspaceId, params) => api.get(`/workspaces/${workspaceId}/prompts`, { params }),
  get: (workspaceId, promptId) => api.get(`/workspaces/${workspaceId}/prompts/${promptId}`),
  create: (workspaceId, data) => api.post(`/workspaces/${workspaceId}/prompts`, data),
  deletePrompt: (workspaceId, promptId) => api.delete(`/workspaces/${workspaceId}/prompts/${promptId}`),
  listVersions: (workspaceId, promptId, params) => api.get(`/workspaces/${workspaceId}/prompts/${promptId}/versions`, { params }),
  getVersion: (workspaceId, promptId, versionId) => api.get(`/workspaces/${workspaceId}/prompts/${promptId}/versions/${versionId}`),
  createVersion: (workspaceId, promptId, data) => api.post(`/workspaces/${workspaceId}/prompts/${promptId}/versions`, data),
  updateVersion: (workspaceId, promptId, versionId, data) => api.put(`/workspaces/${workspaceId}/prompts/${promptId}/versions/${versionId}`, data),
  deleteVersion: (workspaceId, promptId, versionId) => api.delete(`/workspaces/${workspaceId}/prompts/${promptId}/versions/${versionId}`),
}

export const documentAPI = {
  list: (workspaceId, params, config) => api.get(`/workspaces/${workspaceId}/documents`, { ...config, params }),
  upload: (workspaceId, formData) => api.post(`/workspaces/${workspaceId}/documents`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    timeout: 60000,
  }),
  delete: (workspaceId, documentId) => api.delete(`/workspaces/${workspaceId}/documents/${documentId}`),
}

export const agentAPI = {
  list: (workspaceId, params) => api.get(`/workspaces/${workspaceId}/agents`, { params }),
  get: (workspaceId, agentId) => api.get(`/workspaces/${workspaceId}/agents/${agentId}`),
  create: (workspaceId, data) => api.post(`/workspaces/${workspaceId}/agents`, data),
  update: (workspaceId, agentId, data) => api.put(`/workspaces/${workspaceId}/agents/${agentId}`, data),
  delete: (workspaceId, agentId) => api.delete(`/workspaces/${workspaceId}/agents/${agentId}`),
}

export const sessionAPI = {
  list: (workspaceId, agentId, params) => api.get(`/workspaces/${workspaceId}/agents/${agentId}/sessions`, { params }),
  get: (workspaceId, agentId, sessionId) => api.get(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}`),
  getByExternalId: (workspaceId, agentId, externalId) => api.get(`/workspaces/${workspaceId}/agents/${agentId}/sessions`, {
    params: { id: externalId },
  }),
  delete: (workspaceId, agentId, params) =>
    api.delete(`/workspaces/${workspaceId}/agents/${agentId}/sessions/by-labels`, { params }),
}

export const messageAPI = {
  /**
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {{ limit?: number, offset?: number }} [params]
   * @returns {Promise<import('axios').AxiosResponse<{ messages: import('./session-resources').SessionMessage[], _meta: { limit: number, offset: number, total: number } }>>}
   */
  list: (workspaceId, agentId, sessionId, params) =>
    api.get(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/messages`, { params }),
  /**
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {string} messageId
   * @returns {Promise<import('axios').AxiosResponse<import('./session-resources').SessionMessage>>}
   */
  get: (workspaceId, agentId, sessionId, messageId) =>
    api.get(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/messages/${messageId}`),
  /**
   * Create one message. `meta` must be a JSON string on write; see {@link module:session-resources.stringifyMeta}.
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {import('./session-resources').SessionMessageCreatePayload} data
   * @returns {Promise<import('axios').AxiosResponse<import('./session-resources').SessionMessage>>}
   */
  create: (workspaceId, agentId, sessionId, data) =>
    api.post(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/messages`, data),
  /**
   * Batch create up to 100 messages. Each item supports optional `meta` (JSON string).
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {import('./session-resources').SessionMessageBatchCreatePayload} data
   * @returns {Promise<import('axios').AxiosResponse<{ messages: import('./session-resources').SessionMessage[] }>>}
   */
  createBatch: (workspaceId, agentId, sessionId, data) =>
    api.post(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/messages/batch`, data),
  /**
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {string} messageId
   * @param {import('./session-resources').SessionMessageUpdatePayload} data
   * @returns {Promise<import('axios').AxiosResponse<import('./session-resources').SessionMessage>>}
   */
  update: (workspaceId, agentId, sessionId, messageId, data) =>
    api.put(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/messages/${messageId}`, data),
  delete: (workspaceId, agentId, sessionId, messageId) =>
    api.delete(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/messages/${messageId}`),
  /**
   * Batch delete up to 100 messages by id.
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {import('./session-resources').SessionMessageBatchDeletePayload} data
   */
  deleteBatch: (workspaceId, agentId, sessionId, data) =>
    api.delete(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/messages/batch`, { data }),
}

export const memoryAPI = {
  /**
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {{ limit?: number, offset?: number }} [params]
   * @returns {Promise<import('axios').AxiosResponse<{ memories: import('./session-resources').SessionMemory[], _meta: { limit: number, offset: number, total: number } }>>}
   */
  list: (workspaceId, agentId, sessionId, params) =>
    api.get(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/memories`, { params }),
  /**
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {string} memoryId
   * @returns {Promise<import('axios').AxiosResponse<import('./session-resources').SessionMemory>>}
   */
  get: (workspaceId, agentId, sessionId, memoryId) =>
    api.get(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/memories/${memoryId}`),
  /**
   * Create one memory. `meta` must be a JSON string on write; see {@link module:session-resources.stringifyMeta}.
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {import('./session-resources').SessionMemoryCreatePayload} data
   * @returns {Promise<import('axios').AxiosResponse<import('./session-resources').SessionMemory>>}
   */
  create: (workspaceId, agentId, sessionId, data) =>
    api.post(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/memories`, data),
  /**
   * Batch create up to 100 memories. Each item supports optional `meta` (JSON string).
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {import('./session-resources').SessionMemoryBatchCreatePayload} data
   * @returns {Promise<import('axios').AxiosResponse<{ memories: import('./session-resources').SessionMemory[] }>>}
   */
  createBatch: (workspaceId, agentId, sessionId, data) =>
    api.post(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/memories/batch`, data),
  /**
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {string} memoryId
   * @param {import('./session-resources').SessionMemoryUpdatePayload} data
   * @returns {Promise<import('axios').AxiosResponse<import('./session-resources').SessionMemory>>}
   */
  update: (workspaceId, agentId, sessionId, memoryId, data) =>
    api.put(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/memories/${memoryId}`, data),
  delete: (workspaceId, agentId, sessionId, memoryId) =>
    api.delete(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/memories/${memoryId}`),
  /**
   * Batch delete up to 100 memories by id.
   * @param {string} workspaceId
   * @param {string} agentId
   * @param {string} sessionId
   * @param {import('./session-resources').SessionMemoryBatchDeletePayload} data
   */
  deleteBatch: (workspaceId, agentId, sessionId, data) =>
    api.delete(`/workspaces/${workspaceId}/agents/${agentId}/sessions/${sessionId}/memories/batch`, { data }),
}

export { stringifyMeta, buildMessageBatchPayload, buildMemoryBatchPayload } from './session-resources'
