export const AGENT_STATUS_ACTIVE = 'active'

export function agentDetailPath(agentId) {
  return `/agents/${agentId}`
}

export function sessionDetailPath(agentId, sessionId) {
  return `/agents/${agentId}/sessions/${sessionId}`
}

export function formatLabels(labels) {
  if (!labels || typeof labels !== 'object') return []
  return Object.entries(labels).map(([key, value]) => `${key}=${value}`)
}

export function formatAgentDate(iso) {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
  } catch {
    return iso
  }
}

export function buildCreatePayload(form) {
  return {
    name: form.name.trim(),
    description: form.description.trim(),
  }
}

export function buildUpdatePayload(form) {
  return buildCreatePayload(form)
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

export function hasMetadata(meta) {
  if (meta == null) return false
  if (typeof meta === 'string') {
    const trimmed = meta.trim()
    return trimmed !== '' && trimmed !== '{}' && trimmed !== 'null'
  }
  if (typeof meta === 'object') return Object.keys(meta).length > 0
  return false
}

export function formatMetadata(meta) {
  if (!hasMetadata(meta)) return ''
  try {
    if (typeof meta === 'string') {
      return JSON.stringify(JSON.parse(meta), null, 2)
    }
    return JSON.stringify(meta, null, 2)
  } catch {
    return String(meta)
  }
}
