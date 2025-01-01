export const INTEGRATION_TYPES = {
  WEBHOOK: 'webhook',
}

export const WEBHOOK_SIGNATURE_HEADER = 'X-Ziee-Signature-256'

export const WEBHOOK_EVENTS = [
  { value: 'memory.add', labelKey: 'integrations_page.event_memory_add' },
  { value: 'memory.search', labelKey: 'integrations_page.event_memory_search' },
  { value: 'memory.delete', labelKey: 'integrations_page.event_memory_delete' },
  { value: 'access_key.created', labelKey: 'integrations_page.event_access_key_created' },
  { value: 'member.invited', labelKey: 'integrations_page.event_member_invited' },
]


export const MOCK_INTEGRATIONS = [
  {
    id: 'int_01k9m4a1b2c3',
    workspaceId: 'ws_01k9m2x8f1a0',
    type: INTEGRATION_TYPES.WEBHOOK,
    name: 'Slack alerts',
    config: {
      url: 'https://hooks.slack.com/services/T04/B08/q7xK9mNp2vL',
      events: ['memory.add', 'memory.search'],
      enabled: true,
    },
    createdAt: '2026-06-18T09:12:00Z',
    updatedAt: '2026-06-22T14:30:00Z',
    meta: [
      {
        id: 'meta_01k9m4a1b2c4',
        integrationId: 'int_01k9m4a1b2c3',
        key: 'last_delivery',
        value: { at: '2026-06-22T14:28:41Z', status: 'success', statusCode: 200, latencyMs: 142 },
        createdAt: '2026-06-18T09:12:00Z',
        updatedAt: '2026-06-22T14:28:41Z',
      },
      {
        id: 'meta_01k9m4a1b2c5',
        integrationId: 'int_01k9m4a1b2c3',
        key: 'delivery_stats',
        value: { total: 284, failed: 3 },
        createdAt: '2026-06-18T09:12:00Z',
        updatedAt: '2026-06-22T14:28:41Z',
      },
      {
        id: 'meta_01k9m4a1b2c6',
        integrationId: 'int_01k9m4a1b2c3',
        key: 'signing_secret',
        value: { prefix: 'whsec_a8f3' },
        createdAt: '2026-06-18T09:12:00Z',
        updatedAt: '2026-06-18T09:12:00Z',
      },
    ],
  },
  {
    id: 'int_01k9m4a2d4e5',
    workspaceId: 'ws_01k9m2x8f1a0',
    type: INTEGRATION_TYPES.WEBHOOK,
    name: 'CI pipeline',
    config: {
      url: 'https://api.github.com/repos/acme/app/dispatches',
      events: ['access_key.created'],
      enabled: true,
    },
    createdAt: '2026-06-20T16:44:00Z',
    updatedAt: '2026-06-21T08:15:00Z',
    meta: [
      {
        id: 'meta_01k9m4a2d4e6',
        integrationId: 'int_01k9m4a2d4e5',
        key: 'last_delivery',
        value: { at: '2026-06-21T08:14:22Z', status: 'success', statusCode: 204, latencyMs: 89 },
        createdAt: '2026-06-20T16:44:00Z',
        updatedAt: '2026-06-21T08:14:22Z',
      },
      {
        id: 'meta_01k9m4a2d4e7',
        integrationId: 'int_01k9m4a2d4e5',
        key: 'delivery_stats',
        value: { total: 12, failed: 0 },
        createdAt: '2026-06-20T16:44:00Z',
        updatedAt: '2026-06-21T08:14:22Z',
      },
      {
        id: 'meta_01k9m4a2d4e8',
        integrationId: 'int_01k9m4a2d4e5',
        key: 'signing_secret',
        value: { prefix: 'whsec_b2c9' },
        createdAt: '2026-06-20T16:44:00Z',
        updatedAt: '2026-06-20T16:44:00Z',
      },
    ],
  },
  {
    id: 'int_01k9m4a3f6g7',
    workspaceId: 'ws_01k9m2x8f1a0',
    type: INTEGRATION_TYPES.WEBHOOK,
    name: 'Legacy CRM sync',
    config: {
      url: 'https://crm.example.com/webhooks/ziee',
      events: ['memory.add'],
      enabled: false,
    },
    createdAt: '2026-06-10T11:00:00Z',
    updatedAt: '2026-06-19T10:22:00Z',
    meta: [
      {
        id: 'meta_01k9m4a3f6g8',
        integrationId: 'int_01k9m4a3f6g7',
        key: 'last_delivery',
        value: { at: '2026-06-19T10:20:05Z', status: 'failed', statusCode: 503, latencyMs: 30012, error: 'Connection timed out' },
        createdAt: '2026-06-10T11:00:00Z',
        updatedAt: '2026-06-19T10:20:05Z',
      },
      {
        id: 'meta_01k9m4a3f6g9',
        integrationId: 'int_01k9m4a3f6g7',
        key: 'delivery_stats',
        value: { total: 47, failed: 8 },
        createdAt: '2026-06-10T11:00:00Z',
        updatedAt: '2026-06-19T10:20:05Z',
      },
      {
        id: 'meta_01k9m4a3f6ga',
        integrationId: 'int_01k9m4a3f6g7',
        key: 'signing_secret',
        value: { prefix: 'whsec_c1d4' },
        createdAt: '2026-06-10T11:00:00Z',
        updatedAt: '2026-06-10T11:00:00Z',
      },
    ],
  },
]

export function createIntegrationId() {
  return `int_${crypto.randomUUID().replace(/-/g, '').slice(0, 12)}`
}

export function createMetaId() {
  return `meta_${crypto.randomUUID().replace(/-/g, '').slice(0, 12)}`
}

export function createSigningSecret() {
  const chars = 'abcdefghijklmnopqrstuvwxyz0123456789'
  let suffix = ''
  for (let i = 0; i < 28; i++) {
    suffix += chars[Math.floor(Math.random() * chars.length)]
  }
  return `whsec_${suffix}`
}

export function getMetaValue(integration, key) {
  const entry = integration.meta?.find((m) => m.key === key)
  return entry?.value ?? null
}

export function getSigningSecretMeta(integration) {
  return getMetaValue(integration, 'signing_secret')
}

export function setSigningSecretMeta(integration, prefix) {
  const now = new Date().toISOString()
  const existing = integration.meta?.find((m) => m.key === 'signing_secret')
  if (existing) {
    existing.value = { prefix }
    existing.updatedAt = now
    return
  }
  integration.meta = integration.meta ?? []
  integration.meta.push({
    id: createMetaId(),
    integrationId: integration.id,
    key: 'signing_secret',
    value: { prefix },
    createdAt: now,
    updatedAt: now,
  })
}

export function maskSigningSecret(prefix) {
  return `${prefix ?? 'whsec_'}••••••••`
}
