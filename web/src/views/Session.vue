<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-6">
        <router-link :to="agentDetailPath(agentId)" class="text-sm text-primary-600 hover:text-primary-700 hover:underline">
          ← {{ $t('agents_page.back_to_agent') }}
        </router-link>
      </div>

      <div v-if="errorMessage" class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
        {{ errorMessage }}
      </div>

      <div v-if="loading" class="rounded-xl border border-theme-border bg-white p-12 text-center text-theme-textLight">
        {{ $t('common.loading') }}
      </div>

      <div v-else-if="!session || !agent" class="rounded-xl border border-theme-border bg-white p-12 text-center text-theme-textLight">
        {{ $t('agents_page.session_not_found') }}
      </div>

      <template v-else>
        <div class="mb-6">
          <div class="flex flex-wrap items-center gap-3">
            <h1 class="text-2xl font-semibold text-theme-text">{{ session.title || $t('agents_page.untitled_session') }}</h1>
            <span class="inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium" :class="sessionStatusClass(session.status)">
              {{ sessionStatusLabel(session.status) }}
            </span>
          </div>
          <p class="mt-1 text-sm font-mono text-theme-textLight break-all">{{ session.id }}</p>
          <p class="mt-1 text-sm text-theme-textLight">
            {{ agent.name }} · {{ formatDateTime(session.createdAt) }}
          </p>
          <div v-if="formatLabels(session.labels).length" class="mt-3 inline-flex flex-wrap gap-1">
            <span
              v-for="label in formatLabels(session.labels)"
              :key="label"
              class="inline-flex items-center rounded-full bg-theme-hover px-2 py-0.5 text-xs font-mono text-theme-text"
            >
              {{ label }}
            </span>
          </div>
          <div v-if="hasMetadata(session.meta)" class="mt-3">
            <button
              type="button"
              class="inline-flex items-center gap-1.5 text-xs font-medium text-theme-textLight hover:text-theme-text"
              :aria-expanded="sessionMetaOpen"
              @click="sessionMetaOpen = !sessionMetaOpen"
            >
              <span class="transition-transform" :class="sessionMetaOpen ? 'rotate-90' : ''">▸</span>
              {{ sessionMetaOpen ? $t('agents_page.hide_metadata') : $t('agents_page.show_metadata') }}
            </button>
            <pre
              v-if="sessionMetaOpen"
              class="mt-2 max-w-2xl overflow-x-auto rounded-lg border border-theme-border bg-theme-hover px-3 py-2 text-xs font-mono text-theme-text"
            >{{ formatMetadata(session.meta) }}</pre>
          </div>
        </div>

        <div class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
          <div class="border-b border-theme-border px-6">
            <nav class="flex gap-6" aria-label="Session tabs">
              <button
                type="button"
                class="py-3 text-sm font-medium border-b-2 -mb-px transition-colors"
                :class="activeTab === 'messages' ? 'border-primary-600 text-primary-600' : 'border-transparent text-theme-textLight hover:text-theme-text'"
                @click="activeTab = 'messages'"
              >
                {{ $t('agents_page.tab_messages') }}
                <span class="ml-1.5 text-xs text-theme-textLight">({{ messagesTotal }})</span>
              </button>
              <button
                type="button"
                class="py-3 text-sm font-medium border-b-2 -mb-px transition-colors"
                :class="activeTab === 'memory' ? 'border-primary-600 text-primary-600' : 'border-transparent text-theme-textLight hover:text-theme-text'"
                @click="activeTab = 'memory'"
              >
                {{ $t('agents_page.tab_memory') }}
                <span class="ml-1.5 text-xs text-theme-textLight">({{ memoriesTotal }})</span>
              </button>
            </nav>
          </div>

          <!-- Messages -->
          <div v-if="activeTab === 'messages'" class="p-6">
            <div v-if="messagesLoading" class="py-12 text-center text-sm text-theme-textLight">
              {{ $t('common.loading') }}
            </div>
            <div v-else-if="messagesError" class="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
              {{ messagesError }}
            </div>
            <div v-else-if="messages.length === 0" class="py-12 text-center text-sm text-theme-textLight">
              {{ $t('agents_page.no_messages') }}
            </div>
            <div v-else class="space-y-4 max-h-[70vh] overflow-y-auto">
              <div
                v-for="message in messages"
                :key="message.id"
                class="flex"
                :class="message.role === 'user' ? 'justify-end' : 'justify-start'"
              >
                <div
                  class="max-w-[85%] rounded-xl px-4 py-3 text-sm"
                  :class="messageBubbleClass(message.role)"
                >
                  <p class="text-xs font-medium uppercase tracking-wide mb-1 opacity-70">{{ roleLabel(message.role) }}</p>
                  <p class="whitespace-pre-wrap">{{ message.content }}</p>
                  <div class="mt-2 flex flex-wrap items-center gap-3">
                    <p class="text-xs opacity-60">{{ formatDateTime(message.createdAt) }}</p>
                    <button
                      v-if="hasMetadata(message.meta)"
                      type="button"
                      class="inline-flex items-center gap-1 text-xs opacity-70 hover:opacity-100 underline-offset-2 hover:underline"
                      :aria-expanded="isMessageMetaOpen(message.id)"
                      @click="toggleMessageMeta(message.id)"
                    >
                      {{ isMessageMetaOpen(message.id) ? $t('agents_page.hide_metadata') : $t('agents_page.show_metadata') }}
                    </button>
                    <button
                      v-if="canManage"
                      type="button"
                      class="text-xs opacity-70 hover:opacity-100 hover:underline"
                      @click="openDeleteMessage(message)"
                    >
                      {{ $t('common.delete') }}
                    </button>
                  </div>
                  <pre
                    v-if="hasMetadata(message.meta) && isMessageMetaOpen(message.id)"
                    class="mt-2 overflow-x-auto rounded-lg border border-theme-border/60 bg-black/5 px-3 py-2 text-xs font-mono opacity-90"
                    :class="message.role === 'user' ? 'border-white/20 bg-black/10' : ''"
                  >{{ formatMetadata(message.meta) }}</pre>
                </div>
              </div>
            </div>
          </div>

          <!-- Memory -->
          <div v-else class="p-6">
            <div v-if="memoriesLoading" class="py-12 text-center text-sm text-theme-textLight">
              {{ $t('common.loading') }}
            </div>
            <div v-else-if="memoriesError" class="rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
              {{ memoriesError }}
            </div>
            <div v-else-if="memories.length === 0" class="py-12 text-center text-sm text-theme-textLight">
              {{ $t('agents_page.no_memory') }}
            </div>
            <ul v-else class="divide-y divide-theme-border rounded-lg border border-theme-border">
              <li v-for="memory in memories" :key="memory.id" class="px-4 py-4">
                <div class="flex items-start justify-between gap-4">
                  <div class="min-w-0 flex-1">
                    <div class="flex flex-wrap items-center gap-2">
                      <span class="inline-flex rounded-full bg-theme-hover px-2 py-0.5 text-xs font-medium text-theme-text capitalize">
                        {{ memoryKindLabel(memory.kind) }}
                      </span>
                      <span class="text-xs text-theme-textLight">{{ formatDateTime(memory.updatedAt) }}</span>
                    </div>
                    <p class="mt-2 text-sm text-theme-text whitespace-pre-wrap">{{ memory.content }}</p>
                    <div v-if="hasMetadata(memory.meta)" class="mt-2">
                      <button
                        type="button"
                        class="inline-flex items-center gap-1 text-xs text-theme-textLight hover:text-theme-text underline-offset-2 hover:underline"
                        :aria-expanded="isMemoryMetaOpen(memory.id)"
                        @click="toggleMemoryMeta(memory.id)"
                      >
                        {{ isMemoryMetaOpen(memory.id) ? $t('agents_page.hide_metadata') : $t('agents_page.show_metadata') }}
                      </button>
                      <pre
                        v-if="isMemoryMetaOpen(memory.id)"
                        class="mt-2 overflow-x-auto rounded-lg border border-theme-border bg-theme-hover px-3 py-2 text-xs font-mono text-theme-text"
                      >{{ formatMetadata(memory.meta) }}</pre>
                    </div>
                  </div>
                  <button
                    v-if="canManage"
                    type="button"
                    class="shrink-0 text-sm text-red-600 hover:text-red-700 hover:underline"
                    @click="openDeleteMemory(memory)"
                  >
                    {{ $t('common.delete') }}
                  </button>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </template>
    </main>

    <!-- Delete confirmation -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="pendingDelete" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closeDeleteModal">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 class="text-lg font-semibold text-theme-text">
              {{ pendingDelete.type === 'message' ? $t('agents_page.delete_message_title') : $t('agents_page.delete_memory_title') }}
            </h2>
            <p class="mt-2 text-sm text-theme-textLight">
              {{ pendingDelete.type === 'message' ? $t('agents_page.delete_message_confirm') : $t('agents_page.delete_memory_confirm') }}
            </p>
            <p class="mt-3 rounded-lg border border-theme-border bg-theme-bg px-3 py-2 text-sm text-theme-text line-clamp-3">
              {{ pendingDelete.preview }}
            </p>
            <p v-if="deleteError" class="mt-3 text-sm text-red-600">{{ deleteError }}</p>
            <div class="mt-6 flex gap-3 justify-end">
              <button type="button" class="btn-secondary" :disabled="deleting" @click="closeDeleteModal">{{ $t('common.cancel') }}</button>
              <button type="button" class="btn-primary bg-red-600 hover:bg-red-700 focus:ring-red-500" :disabled="deleting" @click="confirmDelete">
                {{ deleting ? $t('agents_page.deleting') : $t('common.delete') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppNav from '@/components/AppNav.vue'
import { agentAPI, sessionAPI, messageAPI, memoryAPI } from '@/api'
import { showFlash } from '@/lib/flash'
import { useWorkspaceContext } from '@/lib/permission'
import {
  agentDetailPath,
  formatLabels,
  formatDateTime,
  hasMetadata,
  formatMetadata,
} from '@/lib/agent'

const LIST_LIMIT = 100

const { t } = useI18n()
const route = useRoute()
const { currentWorkspace, canManage } = useWorkspaceContext()

const activeTab = ref('messages')
const sessionMetaOpen = ref(false)
const openMessageMetaIds = ref({})
const openMemoryMetaIds = ref({})
const pendingDelete = ref(null)
const deleting = ref(false)
const deleteError = ref(null)

const loading = ref(false)
const errorMessage = ref(null)
const agent = ref(null)
const session = ref(null)

const messages = ref([])
const messagesTotal = ref(0)
const messagesLoading = ref(false)
const messagesError = ref(null)

const memories = ref([])
const memoriesTotal = ref(0)
const memoriesLoading = ref(false)
const memoriesError = ref(null)

const agentId = computed(() => route.params.agentId)
const sessionId = computed(() => route.params.sessionId)

async function loadSession() {
  if (!currentWorkspace?.id || !agentId.value || !sessionId.value) return

  loading.value = true
  errorMessage.value = null
  agent.value = null
  session.value = null

  try {
    const [agentRes, sessionRes] = await Promise.all([
      agentAPI.get(currentWorkspace.id, agentId.value),
      sessionAPI.get(currentWorkspace.id, agentId.value, sessionId.value),
    ])
    agent.value = agentRes.data
    session.value = sessionRes.data
    await Promise.all([loadMessages(), loadMemories()])
  } catch (err) {
    if (err.response?.status === 404) {
      agent.value = null
      session.value = null
    } else {
      errorMessage.value = err.response?.data?.errorMessage || t('agents_page.failed_load_session')
    }
  } finally {
    loading.value = false
  }
}

async function loadMessages() {
  messagesLoading.value = true
  messagesError.value = null
  try {
    const res = await messageAPI.list(currentWorkspace.id, agentId.value, sessionId.value, {
      limit: LIST_LIMIT,
      offset: 0,
    })
    messages.value = res.data.messages ?? []
    messagesTotal.value = res.data._meta?.total ?? messages.value.length
  } catch (err) {
    messages.value = []
    messagesTotal.value = 0
    messagesError.value = err.response?.data?.errorMessage || t('agents_page.failed_load_messages')
  } finally {
    messagesLoading.value = false
  }
}

async function loadMemories() {
  memoriesLoading.value = true
  memoriesError.value = null
  try {
    const res = await memoryAPI.list(currentWorkspace.id, agentId.value, sessionId.value, {
      limit: LIST_LIMIT,
      offset: 0,
    })
    memories.value = res.data.memories ?? []
    memoriesTotal.value = res.data._meta?.total ?? memories.value.length
  } catch (err) {
    memories.value = []
    memoriesTotal.value = 0
    memoriesError.value = err.response?.data?.errorMessage || t('agents_page.failed_load_memories')
  } finally {
    memoriesLoading.value = false
  }
}

function sessionStatusLabel(status) {
  if (status === 'closed') return t('agents_page.session_closed')
  return t('agents_page.session_active')
}

function sessionStatusClass(status) {
  if (status === 'closed') return 'bg-theme-hover text-theme-textLight'
  return 'bg-sky-100 text-sky-800'
}

function roleLabel(role) {
  if (role === 'assistant') return t('agents_page.role_assistant')
  if (role === 'system') return t('agents_page.role_system')
  return t('agents_page.role_user')
}

function messageBubbleClass(role) {
  if (role === 'user') return 'bg-primary-600 text-white'
  if (role === 'system') return 'bg-theme-hover text-theme-text border border-theme-border'
  return 'bg-white text-theme-text border border-theme-border shadow-sm'
}

function isMessageMetaOpen(messageId) {
  return Boolean(openMessageMetaIds.value[messageId])
}

function toggleMessageMeta(messageId) {
  openMessageMetaIds.value = {
    ...openMessageMetaIds.value,
    [messageId]: !openMessageMetaIds.value[messageId],
  }
}

function isMemoryMetaOpen(memoryId) {
  return Boolean(openMemoryMetaIds.value[memoryId])
}

function toggleMemoryMeta(memoryId) {
  openMemoryMetaIds.value = {
    ...openMemoryMetaIds.value,
    [memoryId]: !openMemoryMetaIds.value[memoryId],
  }
}

function memoryKindLabel(kind) {
  if (kind === 'summary') return t('agents_page.memory_kind_summary')
  if (kind === 'fact') return t('agents_page.memory_kind_fact')
  if (kind === 'preference') return t('agents_page.memory_kind_preference')
  return kind
}

function openDeleteMessage(message) {
  deleteError.value = null
  pendingDelete.value = { type: 'message', id: message.id, preview: message.content }
}

function openDeleteMemory(memory) {
  deleteError.value = null
  pendingDelete.value = { type: 'memory', id: memory.id, preview: memory.content }
}

function closeDeleteModal() {
  if (deleting.value) return
  pendingDelete.value = null
  deleteError.value = null
}

async function confirmDelete() {
  if (!pendingDelete.value || deleting.value) return

  deleting.value = true
  deleteError.value = null

  try {
    if (pendingDelete.value.type === 'message') {
      await messageAPI.delete(currentWorkspace.id, agentId.value, sessionId.value, pendingDelete.value.id)
      await loadMessages()
    } else {
      await memoryAPI.delete(currentWorkspace.id, agentId.value, sessionId.value, pendingDelete.value.id)
      await loadMemories()
    }
    pendingDelete.value = null
    showFlash(t('common.deleted'))
  } catch (err) {
    deleteError.value = err.response?.data?.errorMessage || t('agents_page.failed_delete_item')
  } finally {
    deleting.value = false
  }
}

onMounted(loadSession)
watch([() => currentWorkspace?.id, agentId, sessionId], loadSession)
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.15s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
