const VARIABLE_PATTERN = /\{\{([a-zA-Z_][a-zA-Z0-9_]*)\}\}/g

export const PROMPT_VARIABLE_EXAMPLE = '{{variable}}'

export function extractPromptVariables(content) {
  if (!content) return []
  const found = new Set()
  let match
  while ((match = VARIABLE_PATTERN.exec(content)) !== null) {
    found.add(match[1])
  }
  return [...found].sort()
}

export function parsePromptLabels(labels) {
  if (!labels) return []
  if (Array.isArray(labels)) return labels
  try {
    const parsed = JSON.parse(labels)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

export function formatPromptConfig(config) {
  if (config == null || config === '') return '{}'
  if (typeof config === 'string') {
    try {
      return JSON.stringify(JSON.parse(config), null, 2)
    } catch {
      return config
    }
  }
  if (typeof config === 'object') {
    return JSON.stringify(config, null, 2)
  }
  return '{}'
}

export function serializePromptConfigForApi(config) {
  if (config == null || config === '') return undefined
  if (typeof config === 'string') {
    const trimmed = config.trim()
    if (!trimmed || trimmed === '{}') return undefined
    return trimmed
  }
  if (typeof config === 'object') {
    const serialized = JSON.stringify(config)
    if (serialized === '{}') return undefined
    return serialized
  }
  return undefined
}

export function hasPromptLabel(labels, label) {
  return parsePromptLabels(labels).includes(label)
}

export function formatPromptDateTime(iso) {
  if (!iso) return ''
  try {
    const d = new Date(iso)
    return d.toLocaleString(undefined, {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    })
  } catch {
    return iso
  }
}

export function formatVersionCommit(version) {
  const hash = version?.commitHash?.trim()
  const message = version?.commitMessage?.trim()
  if (hash && message) return `${hash}- ${message}`
  if (hash) return hash
  return message || ''
}

export function promptDetailPath(promptId) {
  return `/prompts/${encodeURIComponent(promptId)}`
}

export function handleFromName(name) {
  const normalized = (name || '').trim().toLowerCase()
  const parts = normalized.split(/[^a-z0-9]+/).filter(Boolean)
  return parts.join('-').slice(0, 100).replace(/-+$/,'')
}

export function emptyPromptForm() {
  return {
    name: '',
    description: '',
    type: 'text',
    content: '',
    config: '{}',
    production: false,
    commitMessage: '',
  }
}

export function formFromPrompt(prompt) {
  return {
    name: prompt?.name || '',
    description: prompt?.description || '',
    type: prompt?.type || 'text',
    content: prompt?.content || '',
    config: formatPromptConfig(prompt?.config),
    production: !!prompt?.production,
    commitMessage: '',
  }
}

export function buildCreatePayload(form, { nameOverride, omitName } = {}) {
  const payload = {
    type: form.type,
    content: form.content,
    production: !!form.production,
  }
  if (!omitName) {
    payload.name = (nameOverride || form.name).trim()
  }
  payload.description = (form.description || '').trim()
  const config = (form.config || '').trim()
  if (config && config !== '{}') {
    payload.config = config
  }
  if (form.commitMessage.trim()) {
    payload.commitMessage = form.commitMessage.trim()
  }
  return payload
}
