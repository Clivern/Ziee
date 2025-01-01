import { reactive } from 'vue'

export const AGENT_STATUS = {
  ACTIVE: 'active',
  INACTIVE: 'inactive',
}

export const SESSION_STATUS = {
  ACTIVE: 'active',
  CLOSED: 'closed',
}

export const MESSAGE_ROLES = ['user', 'assistant', 'system']

export const MEMORY_KINDS = ['summary', 'fact', 'preference']

/** Known agent config fields — extend this list as new settings are added. */
export const AGENT_CONFIG_FIELDS = [
  {
    key: 'session_retention',
    type: 'select',
    labelKey: 'agents_page.config_session_retention',
    hintKey: 'agents_page.config_session_retention_hint',
    options: [
      { value: 'disabled', labelKey: 'agents_page.config_session_retention_disabled', stored: false },
      { value: '7d_inactivity', labelKey: 'agents_page.config_session_retention_7d', stored: '7d_inactivity' },
    ],
  },
]

function id(prefix) {
  return `${prefix}_${Math.random().toString(36).slice(2, 10)}`
}

const INITIAL_AGENTS = [
  {
    id: 'agt_support01',
    name: 'Support Bot',
    description: 'Handles customer support inquiries with knowledge-base retrieval.',
    status: AGENT_STATUS.ACTIVE,
    config: {
      session_retention: false,
    },
    meta: {
      model: 'gpt-4o',
      system_prompt: 'You are a helpful support agent. Use the knowledge base to answer accurately.',
      temperature: 0.7,
    },
    createdAt: '2026-06-10T09:00:00Z',
    updatedAt: '2026-06-20T14:30:00Z',
  },
  {
    id: 'agt_sales01',
    name: 'Sales Assistant',
    description: 'Qualifies leads and answers product pricing questions.',
    status: AGENT_STATUS.ACTIVE,
    config: {
      session_retention: '7d_inactivity',
    },
    meta: {
      model: 'gpt-4o-mini',
      system_prompt: 'You are a sales assistant. Be concise and guide users toward the right plan.',
      temperature: 0.5,
    },
    createdAt: '2026-06-12T11:15:00Z',
    updatedAt: '2026-06-18T16:45:00Z',
  },
]

const INITIAL_SESSIONS = [
  {
    id: '8f2a91bc-3d4e-5f60-a1b2-c3d4e5f67890',
    agentId: 'agt_support01',
    title: 'Refund policy question',
    status: SESSION_STATUS.CLOSED,
    labels: { channel: 'web', topic: 'billing' },
    meta: {
      ip_address: '52.14.88.201',
      user_agent: 'Mozilla/5.0',
      source: 'web_widget',
    },
    createdAt: '2026-06-19T10:05:00Z',
    updatedAt: '2026-06-19T10:12:00Z',
  },
  {
    id: '3c7d44ef-1a2b-3c4d-5e6f-7890abcdef12',
    agentId: 'agt_support01',
    title: 'API integration help',
    status: SESSION_STATUS.ACTIVE,
    labels: { channel: 'api', topic: 'developer' },
    meta: {
      ip_address: '52.14.88.201',
      user_agent: 'actx0-python/0.1.0',
      access_key: 'ak_support_prod',
    },
    createdAt: '2026-06-21T08:30:00Z',
    updatedAt: '2026-06-21T09:01:00Z',
  },
  {
    id: '9a1b22cd-4e5f-6789-abcd-ef0123456789',
    agentId: 'agt_sales01',
    title: 'Enterprise plan inquiry',
    status: SESSION_STATUS.ACTIVE,
    labels: { channel: 'web', topic: 'pricing' },
    meta: {
      ip_address: '203.0.113.42',
      user_agent: 'Mozilla/5.0',
      referrer: '/pricing',
    },
    createdAt: '2026-06-20T15:20:00Z',
    updatedAt: '2026-06-20T15:35:00Z',
  },
]

const INITIAL_MESSAGES = [
  {
    id: 'msg_001',
    sessionId: '8f2a91bc-3d4e-5f60-a1b2-c3d4e5f67890',
    role: 'user',
    content: 'What is your refund policy for annual subscriptions?',
    meta: { locale: 'en-US' },
    createdAt: '2026-06-19T10:05:12Z',
  },
  {
    id: 'msg_002',
    sessionId: '8f2a91bc-3d4e-5f60-a1b2-c3d4e5f67890',
    role: 'assistant',
    content: 'Annual subscriptions can be refunded within 14 days of purchase. After that window, we offer prorated credits toward future billing cycles.',
    meta: {
      model: 'gpt-4o',
      tokens_in: 48,
      tokens_out: 62,
      latency_ms: 842,
    },
    createdAt: '2026-06-19T10:05:28Z',
  },
  {
    id: 'msg_003',
    sessionId: '8f2a91bc-3d4e-5f60-a1b2-c3d4e5f67890',
    role: 'user',
    content: 'Does that apply if I used the API heavily?',
    createdAt: '2026-06-19T10:06:05Z',
  },
  {
    id: 'msg_004',
    sessionId: '8f2a91bc-3d4e-5f60-a1b2-c3d4e5f67890',
    role: 'assistant',
    content: 'Yes, the 14-day policy applies regardless of API usage. Heavy usage may affect the prorated credit amount after the refund window.',
    createdAt: '2026-06-19T10:06:22Z',
  },
  {
    id: 'msg_005',
    sessionId: '3c7d44ef-1a2b-3c4d-5e6f-7890abcdef12',
    role: 'user',
    content: 'How do I authenticate API requests for agent sessions?',
    meta: { locale: 'en-US' },
    createdAt: '2026-06-21T08:30:45Z',
  },
  {
    id: 'msg_006',
    sessionId: '3c7d44ef-1a2b-3c4d-5e6f-7890abcdef12',
    role: 'assistant',
    content: 'Use a workspace access key in the Authorization header: Bearer <token>. Create keys under Settings → Access keys.',
    meta: {
      model: 'gpt-4o',
      tokens_in: 36,
      tokens_out: 41,
      latency_ms: 612,
      docs_retrieved: 2,
    },
    createdAt: '2026-06-21T08:31:10Z',
  },
  {
    id: 'msg_007',
    sessionId: '9a1b22cd-4e5f-6789-abcd-ef0123456789',
    role: 'user',
    content: 'We need SSO and dedicated support. What does Enterprise include?',
    createdAt: '2026-06-20T15:20:18Z',
  },
  {
    id: 'msg_008',
    sessionId: '9a1b22cd-4e5f-6789-abcd-ef0123456789',
    role: 'assistant',
    content: 'Enterprise includes SAML SSO, dedicated support, custom SLAs, and higher usage limits. I can connect you with sales for a tailored quote.',
    createdAt: '2026-06-20T15:20:44Z',
  },
]

const INITIAL_MEMORIES = [
  {
    id: 'mem_001',
    sessionId: '8f2a91bc-3d4e-5f60-a1b2-c3d4e5f67890',
    kind: 'summary',
    content: 'User asked about annual subscription refunds and whether heavy API usage affects eligibility.',
    createdAt: '2026-06-19T10:12:00Z',
    updatedAt: '2026-06-19T10:12:00Z',
  },
  {
    id: 'mem_002',
    sessionId: '8f2a91bc-3d4e-5f60-a1b2-c3d4e5f67890',
    kind: 'fact',
    content: 'Customer is on an annual plan and concerned about API usage impact on refunds.',
    createdAt: '2026-06-19T10:12:00Z',
    updatedAt: '2026-06-19T10:12:00Z',
  },
  {
    id: 'mem_003',
    sessionId: '3c7d44ef-1a2b-3c4d-5e6f-7890abcdef12',
    kind: 'summary',
    content: 'Developer needs API authentication guidance for agent sessions.',
    createdAt: '2026-06-21T09:01:00Z',
    updatedAt: '2026-06-21T09:01:00Z',
  },
  {
    id: 'mem_004',
    sessionId: '9a1b22cd-4e5f-6789-abcd-ef0123456789',
    kind: 'preference',
    content: 'Prospect requires SSO and dedicated support — enterprise tier fit.',
    createdAt: '2026-06-20T15:35:00Z',
    updatedAt: '2026-06-20T15:35:00Z',
  },
]

export const agentStore = reactive({
  agents: [...INITIAL_AGENTS],
  sessions: [...INITIAL_SESSIONS],
  messages: [...INITIAL_MESSAGES],
  memories: [...INITIAL_MEMORIES],
})

export function generateAgentId() {
  return id('agt')
}

export function generateSessionId() {
  return crypto.randomUUID()
}

export function createAgent({ name, description }) {
  const now = new Date().toISOString()
  const agent = {
    id: generateAgentId(),
    name: name.trim(),
    description: description?.trim() || null,
    status: AGENT_STATUS.ACTIVE,
    config: null,
    meta: null,
    createdAt: now,
    updatedAt: now,
  }
  agentStore.agents.unshift(agent)
  return agent
}

export function getAgent(agentId) {
  return agentStore.agents.find((a) => a.id === agentId) || null
}

export function agentConfigToForm(config) {
  const out = {}
  for (const field of AGENT_CONFIG_FIELDS) {
    if (field.type !== 'select') continue
    const stored = config?.[field.key]
    const option = field.options.find((o) => o.stored === stored) ?? field.options[0]
    out[field.key] = option.value
  }
  return out
}

export function buildAgentConfig(existingConfig, formConfig) {
  const next = { ...(existingConfig || {}) }
  for (const field of AGENT_CONFIG_FIELDS) {
    if (field.type !== 'select') continue
    const option = field.options.find((o) => o.value === formConfig[field.key])
    if (option) next[field.key] = option.stored
  }
  return Object.keys(next).length ? next : null
}

export function updateAgent(agentId, { name, description, config }) {
  const agent = getAgent(agentId)
  if (!agent) return null
  const now = new Date().toISOString()
  agent.name = name.trim()
  agent.description = description?.trim() || null
  agent.config = buildAgentConfig(agent.config, config)
  agent.updatedAt = now
  return agent
}

export function deleteAgent(agentId) {
  const sessionIds = agentStore.sessions.filter((s) => s.agentId === agentId).map((s) => s.id)
  agentStore.agents = agentStore.agents.filter((a) => a.id !== agentId)
  agentStore.sessions = agentStore.sessions.filter((s) => s.agentId !== agentId)
  agentStore.messages = agentStore.messages.filter((m) => !sessionIds.includes(m.sessionId))
  agentStore.memories = agentStore.memories.filter((m) => !sessionIds.includes(m.sessionId))
}

export function listSessionsByAgent(agentId) {
  return agentStore.sessions
    .filter((s) => s.agentId === agentId)
    .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
}

export function getSession(sessionId) {
  return agentStore.sessions.find((s) => s.id === sessionId) || null
}

export function countSessionsByAgent(agentId) {
  return agentStore.sessions.filter((s) => s.agentId === agentId).length
}

export function listMessagesBySession(sessionId) {
  return agentStore.messages
    .filter((m) => m.sessionId === sessionId)
    .sort((a, b) => new Date(a.createdAt) - new Date(b.createdAt))
}

export function listMemoriesBySession(sessionId) {
  return agentStore.memories
    .filter((m) => m.sessionId === sessionId)
    .sort((a, b) => new Date(b.updatedAt) - new Date(a.updatedAt))
}

export function sessionMessageCount(sessionId) {
  return agentStore.messages.filter((m) => m.sessionId === sessionId).length
}

export function sessionMemoryCount(sessionId) {
  return agentStore.memories.filter((m) => m.sessionId === sessionId).length
}

export function deleteMessage(messageId) {
  agentStore.messages = agentStore.messages.filter((m) => m.id !== messageId)
}

export function deleteMemory(memoryId) {
  agentStore.memories = agentStore.memories.filter((m) => m.id !== memoryId)
}

export function formatLabels(labels) {
  if (!labels || typeof labels !== 'object') return []
  return Object.entries(labels).map(([key, value]) => `${key}=${value}`)
}

export function hasMetadata(meta) {
  return meta != null && typeof meta === 'object' && Object.keys(meta).length > 0
}

export function formatMetadata(meta) {
  if (!hasMetadata(meta)) return ''
  try {
    return JSON.stringify(meta, null, 2)
  } catch {
    return String(meta)
  }
}

export function agentDetailPath(agentId) {
  return `/agents/${agentId}`
}

export function sessionDetailPath(agentId, sessionId) {
  return `/agents/${agentId}/sessions/${sessionId}`
}

export function formatDate(iso) {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
  } catch {
    return iso
  }
}

export function formatDateTime(iso) {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleString(undefined, {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return iso
  }
}
