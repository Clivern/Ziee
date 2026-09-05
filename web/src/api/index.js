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
    const isProfileEndpoint = error.config?.url === '/action/profile'
    const isLogoutEndpoint = error.config?.url === 'public/action/logout'
    const isOnPublicPage = window.location.pathname === '/' ||
                           window.location.pathname === '/login' ||
                           window.location.pathname === '/setup'

    if (error.response?.status === 401 &&
        !isProfileEndpoint &&
        !isLogoutEndpoint &&
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
export const health_api = {
  check: () => api.get('public/_health'),
  ready: () => api.get('public/_ready'),
}

export const auth_api = {
  logout: () => api.post('public/action/logout'),
  getProfile: () => api.get('/action/profile'),
  updateProfile: (data) => api.put('/action/profile', data),
}

export const setup_api = {
  install: (data) => api.post('public/action/setup', data),
  checkInstalled: () => api.get('public/action/setup/status'),
}

// API endpoints for settings
export const settings_api = {
  getSettings: () => api.get('/action/settings'),
  updateSettings: (data) => api.put('/action/settings', data),
}

// API endpoints for workspaces (projects)
export const workspace_api = {
  list: (params) => api.get('/workspaces', { params }),
  get: (id) => api.get(`/workspaces/${id}`),
  create: (data) => api.post('/workspaces', data),
  update: (id, data) => api.put(`/workspaces/${id}`, data),
  delete: (id) => api.delete(`/workspaces/${id}`),
}

export const invite_api = {
  getByToken: (token) => api.get(`/action/invite-by-token/${encodeURIComponent(token)}`),
  acceptByToken: (token) => api.post(`/action/accept-invite/${encodeURIComponent(token)}`),
  rejectByToken: (token) => api.post(`/action/reject-invite/${encodeURIComponent(token)}`),
}

export const workspace_member_api = {
  list: (workspaceId, params) => api.get(`/workspaces/${workspaceId}/members`, { params }),
  updateRole: (workspaceId, userId, data) => api.put(`/workspaces/${workspaceId}/members/${userId}`, data),
  delete: (workspaceId, userId) => api.delete(`/workspaces/${workspaceId}/members/${userId}`),
}

export const workspace_invite_api = {
  list: (workspaceId, params) => api.get(`/workspaces/${workspaceId}/invites`, { params }),
  create: (workspaceId, data) => api.post(`/workspaces/${workspaceId}/invites`, data),
  delete: (workspaceId, id) => api.delete(`/workspaces/${workspaceId}/invites/${id}`),
}

export const stats_api = {
  get: (workspaceId) => api.get(`/workspaces/${workspaceId}/stats`),
}

export const billing_api = {
  status: (workspaceId) => api.get(`/workspaces/${workspaceId}/billing`),
  usage: (workspaceId) => api.get(`/workspaces/${workspaceId}/billing/usage`),
  checkout: (workspaceId, data) => api.post(`/workspaces/${workspaceId}/billing/checkout`, data),
  portal: (workspaceId) => api.post(`/workspaces/${workspaceId}/billing/portal`),
}

export const workspace_access_key_api = {
  list: (workspaceId, params) => api.get(`/workspaces/${workspaceId}/keys`, { params }),
  create: (workspaceId, data) => api.post(`/workspaces/${workspaceId}/keys`, data),
  delete: (workspaceId, id) => api.delete(`/workspaces/${workspaceId}/keys/${id}`),
}

// API endpoints for current user's API keys
export const api_keys_api = {
  list: (params) => api.get('/apiKeys', { params }),
  create: (data) => api.post('/apiKeys', data),
  delete: (id) => api.delete(`/apiKeys/${id}`),
}

export const document_api = {
  list: (workspaceId, params, config) => api.get(`/workspaces/${workspaceId}/documents`, { ...config, params }),
  upload: (workspaceId, formData) => api.post(`/workspaces/${workspaceId}/documents`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    timeout: 60000,
  }),
  delete: (workspaceId, documentId) => api.delete(`/workspaces/${workspaceId}/documents/${documentId}`),
}

export const github_api = {
  listInstallations: () => api.get('/action/github/installations'),
  attachInstallation: (id, data) => api.post(`/action/github/installations/${id}/attach`, data),
}
