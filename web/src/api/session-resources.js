/**
 * Session messages and memories share the same `meta` contract.
 *
 * ## Write (create, update, batch create)
 * - `meta` is optional.
 * - Must be a JSON **string** (not a nested object in the request body).
 * - Must contain valid JSON when present (object, array, string, number, etc.).
 * - Omit the field or pass `null` when there is no metadata.
 *
 * ## Read (list, get, batch create response)
 * - `meta` is omitted when empty.
 * - When present, it is a parsed JSON **value** (usually an object).
 *
 * @typedef {Record<string, unknown>} SessionResourceMeta
 */

/**
 * @typedef {Object} SessionMessageCreatePayload
 * @property {'system'|'user'|'assistant'} role
 * @property {string} content
 * @property {string} [meta] JSON string; use {@link stringifyMeta} when building from an object.
 */

/**
 * @typedef {Object} SessionMessageUpdatePayload
 * @property {string} content
 * @property {'system'|'user'|'assistant'} [role]
 * @property {string} [meta] JSON string; use {@link stringifyMeta} when building from an object.
 */

/**
 * @typedef {Object} SessionMessageBatchCreatePayload
 * @property {SessionMessageCreatePayload[]} messages Max 100 items.
 */

/**
 * @typedef {Object} SessionMessageBatchDeletePayload
 * @property {string[]} ids Max 100 message ids.
 */

/**
 * @typedef {Object} SessionMessage
 * @property {string} id
 * @property {string} sessionId
 * @property {'system'|'user'|'assistant'} role
 * @property {string} content
 * @property {SessionResourceMeta|null|undefined} [meta] Parsed JSON object when present.
 * @property {string} createdAt RFC3339 timestamp.
 */

/**
 * @typedef {Object} SessionMemoryCreatePayload
 * @property {'summary'|'fact'|'preference'} kind
 * @property {string} content
 * @property {string} [meta] JSON string; use {@link stringifyMeta} when building from an object.
 */

/**
 * @typedef {Object} SessionMemoryUpdatePayload
 * @property {string} content
 * @property {'summary'|'fact'|'preference'} [kind]
 * @property {string} [meta] JSON string; use {@link stringifyMeta} when building from an object.
 */

/**
 * @typedef {Object} SessionMemoryBatchCreatePayload
 * @property {SessionMemoryCreatePayload[]} memories Max 100 items.
 */

/**
 * @typedef {Object} SessionMemoryBatchDeletePayload
 * @property {string[]} ids Max 100 memory ids.
 */

/**
 * @typedef {Object} SessionMemory
 * @property {string} id
 * @property {string} sessionId
 * @property {'summary'|'fact'|'preference'} kind
 * @property {string} content
 * @property {SessionResourceMeta|null|undefined} [meta] Parsed JSON object when present.
 * @property {string} createdAt RFC3339 timestamp.
 * @property {string} updatedAt RFC3339 timestamp.
 */

/**
 * Converts client-side metadata to the JSON string the API expects on write.
 *
 * @param {SessionResourceMeta|string|null|undefined} meta
 * @returns {string|undefined}
 *
 * @example
 * stringifyMeta({ source: 'sdk', model: 'gpt-4' })
 * // => '{"source":"sdk","model":"gpt-4"}'
 *
 * @example
 * messageAPI.create(workspaceId, agentId, sessionId, {
 *   role: 'assistant',
 *   content: 'Done.',
 *   meta: stringifyMeta({ toolCalls: [{ name: 'search', id: 'call_1' }] }),
 * })
 */
export function stringifyMeta(meta) {
  if (meta == null) return undefined
  if (typeof meta === 'string') {
    const trimmed = meta.trim()
    return trimmed === '' ? undefined : trimmed
  }
  return JSON.stringify(meta)
}

/**
 * Builds a batch message create payload with optional meta on each item.
 *
 * @param {Array<{ role: SessionMessageCreatePayload['role'], content: string, meta?: SessionResourceMeta|string }>} items
 * @returns {SessionMessageBatchCreatePayload}
 */
export function buildMessageBatchPayload(items) {
  return {
    messages: items.map(({ role, content, meta }) => {
      const metaValue = stringifyMeta(meta)
      return metaValue ? { role, content, meta: metaValue } : { role, content }
    }),
  }
}

/**
 * Builds a batch memory create payload with optional meta on each item.
 *
 * @param {Array<{ kind: SessionMemoryCreatePayload['kind'], content: string, meta?: SessionResourceMeta|string }>} items
 * @returns {SessionMemoryBatchCreatePayload}
 */
export function buildMemoryBatchPayload(items) {
  return {
    memories: items.map(({ kind, content, meta }) => {
      const metaValue = stringifyMeta(meta)
      return metaValue ? { kind, content, meta: metaValue } : { kind, content }
    }),
  }
}
