<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-6">
        <router-link to="/agents" class="text-sm text-primary-600 hover:text-primary-700 hover:underline">
          ← {{ $t('agents_page.back_to_agents') }}
        </router-link>
      </div>

      <div v-if="errorMessage" class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
        {{ errorMessage }}
      </div>

      <div v-if="loading" class="rounded-xl border border-theme-border bg-white p-12 text-center text-theme-textLight">
        {{ $t('common.loading') }}
      </div>

      <div v-else-if="!agent" class="rounded-xl border border-theme-border bg-white p-12 text-center text-theme-textLight">
        {{ $t('agents_page.not_found') }}
      </div>

      <template v-else>
        <div class="mb-8 rounded-xl border border-theme-border bg-white shadow-sm p-6">
          <div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div class="min-w-0 flex-1">
              <h1 class="text-2xl font-semibold text-theme-text">{{ agent.name }}</h1>
              <p class="mt-2 text-sm text-theme-textLight">{{ agent.description }}</p>
              <p class="mt-2 text-xs font-mono text-theme-textLight break-all">{{ agent.id }}</p>
              <dl class="mt-4 grid grid-cols-1 gap-3 text-sm sm:grid-cols-2">
                <div>
                  <dt class="text-theme-textLight">{{ $t('agents_page.created_at') }}</dt>
                  <dd class="mt-0.5 text-theme-text">{{ formatAgentDate(agent.createdAt) }}</dd>
                </div>
              </dl>
            </div>
            <button v-if="canEdit" type="button" class="btn-secondary shrink-0" @click="openEditModal">
              {{ $t('common.edit') }}
            </button>
          </div>
        </div>

        <section class="mb-8 rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
          <div class="border-b border-theme-border bg-primary-50/40 px-6 py-5">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div>
                <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.getting_started') }}</h2>
                <p class="mt-1 max-w-2xl text-sm text-theme-textLight leading-relaxed">
                  {{ $t('agents_page.getting_started_desc') }}
                </p>
              </div>
              <span class="rounded-full border border-theme-border bg-white px-3 py-1 text-xs font-medium text-theme-textLight">
                {{ gettingStartedSteps.length }} {{ $t('agents_page.steps') }}
              </span>
            </div>
          </div>

          <ol class="divide-y divide-theme-border">
            <li
              v-for="(step, index) in gettingStartedSteps"
              :key="step.title"
              class="group flex gap-4 px-6 py-5 transition-colors hover:bg-theme-hover/40"
            >
              <div class="relative flex flex-col items-center">
                <span
                  class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-theme-border bg-white text-theme-text shadow-sm transition-colors group-hover:border-primary-300 group-hover:bg-primary-50"
                  :class="stepIconClass(step.icon)"
                >
                  <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" :d="stepIcons[step.icon]" />
                  </svg>
                </span>
                <span
                  v-if="index < gettingStartedSteps.length - 1"
                  class="mt-2 w-px flex-1 min-h-[1.5rem] bg-theme-border"
                  aria-hidden="true"
                />
              </div>
              <div class="min-w-0 flex-1 pb-1">
                <span class="text-xs font-medium uppercase tracking-wide text-theme-textLight">
                  {{ $t('agents_page.step_label', { number: step.number }) }}
                </span>
                <p class="mt-1 text-base font-medium text-theme-text">{{ step.title }}</p>
                <p class="mt-1 text-sm leading-relaxed text-theme-textLight">{{ step.description }}</p>
                <div v-if="step.code" class="mt-3 flex flex-wrap items-center gap-2">
                  <code class="rounded-lg border border-theme-border bg-theme-bg px-3 py-2 text-xs font-mono text-theme-text break-all">{{ step.code }}</code>
                  <button type="button" class="btn-secondary px-3 py-1.5 text-xs" @click="copyText(step.code)">
                    {{ copiedField === step.code ? $t('agents_page.copied') : $t('agents_page.copy') }}
                  </button>
                </div>
                <button
                  v-if="step.action === 'edit' && canEdit"
                  type="button"
                  class="btn-secondary mt-3 inline-flex items-center px-3 py-1.5 text-sm"
                  @click="openEditModal"
                >
                  {{ $t('agents_page.edit_agent') }}
                </button>
                <router-link
                  v-else-if="step.link"
                  :to="step.link.to"
                  class="btn-secondary mt-3 inline-flex items-center px-3 py-1.5 text-sm"
                >
                  {{ step.link.label }}
                </router-link>
              </div>
            </li>
          </ol>
        </section>

        <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
          <div class="border-b border-theme-border px-6 py-4">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.sessions_title') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('agents_page.sessions_desc') }}</p>
          </div>

          <div v-if="sessionsLoading" class="p-12 text-center text-sm text-theme-textLight">
            {{ $t('common.loading') }}
          </div>

          <div v-else-if="sessionsError" class="p-6 text-center text-sm text-red-800">
            {{ sessionsError }}
          </div>

          <div v-else-if="sessions.length === 0" class="p-12 text-center text-sm text-theme-textLight">
            {{ $t('agents_page.no_sessions') }}
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-theme-border">
              <thead class="bg-theme-hover">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.session') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.status') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.labels') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.created_at') }}</th>
                  <th class="px-6 py-3 text-right text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.actions') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-theme-border bg-white">
                <tr v-for="session in sessions" :key="session.id" class="hover:bg-theme-hover/60">
                  <td class="px-6 py-4">
                    <router-link :to="sessionDetailPath(agentId, session.id)" class="text-sm font-medium text-primary-600 hover:text-primary-700 hover:underline">
                      {{ session.title || $t('agents_page.untitled_session') }}
                    </router-link>
                    <p class="mt-0.5 text-xs font-mono text-theme-textLight break-all">{{ session.externalId }}</p>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span class="inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium" :class="sessionStatusClass(session.status)">
                      {{ sessionStatusLabel(session.status) }}
                    </span>
                  </td>
                  <td class="px-6 py-4">
                    <div v-if="formatLabels(session.labels).length" class="flex flex-wrap gap-1">
                      <span
                        v-for="label in formatLabels(session.labels)"
                        :key="label"
                        class="inline-flex items-center rounded-full bg-theme-hover px-2 py-0.5 text-xs font-mono text-theme-text"
                      >
                        {{ label }}
                      </span>
                    </div>
                    <span v-else class="text-sm text-theme-textLight">—</span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-textLight">{{ formatDateTime(session.createdAt) }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <span class="inline-flex items-center gap-3">
                      <router-link :to="sessionDetailPath(agentId, session.id)" class="text-primary-600 hover:text-primary-700 hover:underline">
                        {{ $t('agents_page.view') }}
                      </router-link>
                      <button
                        v-if="canEdit"
                        type="button"
                        class="text-red-600 hover:text-red-700 hover:underline disabled:opacity-50"
                        :disabled="deletingSessionId === session.id"
                        @click="openDeleteSessionModal(session)"
                      >
                        {{ deletingSessionId === session.id ? $t('agents_page.deleting') : $t('agents_page.delete') }}
                      </button>
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div v-if="sessionsTotal > 0" class="bg-white px-6 py-4 border-t border-theme-border flex items-center justify-between">
            <div class="text-sm text-theme-textLight">
              {{ $t('agents_page.showing', { from: sessionsOffset + 1, to: Math.min(sessionsOffset + sessionsLimit, sessionsTotal), total: sessionsTotal }) }}
            </div>
            <div class="flex items-center gap-2">
              <button type="button" class="btn-secondary text-sm disabled:opacity-50" :disabled="sessionsOffset === 0" @click="goToSessionsPage(sessionsOffset - sessionsLimit)">
                {{ $t('agents_page.previous') }}
              </button>
              <button type="button" class="btn-secondary text-sm disabled:opacity-50" :disabled="sessionsOffset + sessionsLimit >= sessionsTotal" @click="goToSessionsPage(sessionsOffset + sessionsLimit)">
                {{ $t('agents_page.next') }}
              </button>
            </div>
          </div>
        </section>
      </template>
    </main>

    <!-- Delete session modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="deleteSessionTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closeDeleteSessionModal">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.delete_session_title') }}</h2>
            <div v-if="deleteSessionError" class="mt-3 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ deleteSessionError }}</div>
            <p class="mt-2 text-sm text-theme-textLight">
              {{ $t('agents_page.delete_session_confirm', { name: deleteSessionTarget.title || $t('agents_page.untitled_session') }) }}
            </p>
            <div class="mt-6 flex gap-3 justify-end">
              <button type="button" class="btn-secondary" :disabled="deletingSession" @click="closeDeleteSessionModal">{{ $t('common.cancel') }}</button>
              <button type="button" class="btn-primary bg-red-600 hover:bg-red-700 focus:ring-red-500" :disabled="deletingSession" @click="confirmDeleteSession">
                {{ deletingSession ? $t('agents_page.deleting') : $t('agents_page.delete') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Edit modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closeEditModal">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-lg p-6 max-h-[90vh] overflow-y-auto">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.edit_title') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('agents_page.edit_desc') }}</p>
            <div v-if="editError" class="mt-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ editError }}</div>
            <form class="mt-4 space-y-4" @submit.prevent="handleSave">
              <div>
                <label for="edit-agent-name" class="form-label">{{ $t('agents_page.name') }}</label>
                <input id="edit-agent-name" v-model="editForm.name" type="text" required class="input-field" :placeholder="$t('agents_page.name_placeholder')" />
              </div>
              <div>
                <label for="edit-agent-description" class="form-label">{{ $t('agents_page.description') }}</label>
                <textarea id="edit-agent-description" v-model="editForm.description" rows="2" required class="input-field" :placeholder="$t('agents_page.description_placeholder')" />
              </div>

              <div class="flex gap-3 justify-end">
                <button type="button" class="btn-secondary" :disabled="saveLoading" @click="closeEditModal">{{ $t('common.cancel') }}</button>
                <button type="submit" class="btn-primary" :disabled="saveLoading">
                  {{ saveLoading ? $t('agents_page.saving') : $t('common.save') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppNav from '@/components/AppNav.vue'
import { showFlash } from '@/lib/flash'
import { useUrlPagination } from '@/lib/pagination'
import { useWorkspaceContext } from '@/lib/permission'
import { agentAPI, sessionAPI } from '@/api'
import { formatAgentDate, buildUpdatePayload, sessionDetailPath, formatLabels, formatDateTime } from '@/lib/agent'

const { t } = useI18n()
const route = useRoute()
const { currentWorkspace, canManage: canEdit } = useWorkspaceContext()
const { offset: sessionsOffset, limit: sessionsLimit, goToOffset: goToSessionsOffset, goToPage: setSessionsPage } = useUrlPagination({ limit: 50 })

const agent = ref(null)
const loading = ref(false)
const errorMessage = ref(null)

const sessions = ref([])
const sessionsLoading = ref(false)
const sessionsError = ref(null)
const sessionsTotal = ref(0)

const showEditModal = ref(false)
const editError = ref(null)
const saveLoading = ref(false)
const copiedField = ref(null)
const editForm = ref({ name: '', description: '' })

const deleteSessionTarget = ref(null)
const deleteSessionError = ref(null)
const deletingSession = ref(false)
const deletingSessionId = ref(null)

const agentId = computed(() => route.params.agentId)

const stepIcons = {
  settings: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z',
  book: 'M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253',
  chat: 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z',
  chart: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z',
}

const gettingStartedSteps = computed(() => [
  {
    number: 1,
    icon: 'settings',
    title: t('agents_page.step_configure'),
    description: t('agents_page.step_configure_desc'),
    action: 'edit',
  },
  {
    number: 2,
    icon: 'book',
    title: t('agents_page.step_knowledge'),
    description: t('agents_page.step_knowledge_desc'),
    link: { to: '/knowledge', label: t('agents_page.step_knowledge_link') },
  },
  {
    number: 3,
    icon: 'chat',
    title: t('agents_page.step_session'),
    description: t('agents_page.step_session_desc'),
    code: agent.value?.handle ?? '',
  },
  {
    number: 4,
    icon: 'chart',
    title: t('agents_page.step_monitor'),
    description: t('agents_page.step_monitor_desc'),
  },
])

function stepIconClass(icon) {
  const map = {
    settings: 'text-sky-700',
    book: 'text-violet-700',
    chat: 'text-emerald-700',
    chart: 'text-primary-800',
  }
  return map[icon] ?? 'text-theme-text'
}

async function loadAgent() {
  loading.value = true
  errorMessage.value = null
  agent.value = null
  sessions.value = []
  sessionsTotal.value = 0
  try {
    const res = await agentAPI.get(currentWorkspace.id, agentId.value)
    agent.value = res.data
    await loadSessions()
  } catch (err) {
    if (err.response?.status === 404) {
      agent.value = null
    } else {
      errorMessage.value = err.response?.data?.errorMessage || t('agents_page.failed_load')
    }
  } finally {
    loading.value = false
  }
}

async function loadSessions() {
  sessionsLoading.value = true
  sessionsError.value = null
  try {
    const res = await sessionAPI.list(currentWorkspace.id, agentId.value, {
      limit: sessionsLimit,
      offset: sessionsOffset.value,
    })
    sessions.value = res.data.sessions ?? []
    sessionsTotal.value = res.data._meta?.total ?? 0
  } catch (err) {
    sessions.value = []
    sessionsTotal.value = 0
    sessionsError.value = err.response?.data?.errorMessage || t('agents_page.failed_load_sessions')
  } finally {
    sessionsLoading.value = false
  }
}

function goToSessionsPage(offset) {
  if (offset < 0 || offset >= sessionsTotal.value) return
  goToSessionsOffset(offset)
}

watch(() => route.query.page, () => {
  if (agent.value) loadSessions()
})

function sessionStatusLabel(status) {
  if (status === 'closed') return t('agents_page.session_closed')
  return t('agents_page.session_active')
}

function sessionStatusClass(status) {
  if (status === 'closed') return 'bg-theme-hover text-theme-textLight'
  return 'bg-sky-100 text-sky-800'
}

function openEditModal() {
  if (!agent.value) return
  editForm.value = {
    name: agent.value.name,
    description: agent.value.description || '',
  }
  editError.value = null
  showEditModal.value = true
}

function closeEditModal() {
  if (saveLoading.value) return
  showEditModal.value = false
}

async function handleSave() {
  const name = editForm.value.name.trim()
  const description = editForm.value.description.trim()
  if (!name) {
    editError.value = t('agents_page.name_required')
    return
  }
  if (!description) {
    editError.value = t('agents_page.description_required')
    return
  }

  saveLoading.value = true
  editError.value = null
  try {
    const res = await agentAPI.update(currentWorkspace.id, agentId.value, buildUpdatePayload(editForm.value))
    agent.value = res.data
    showEditModal.value = false
    showFlash(t('common.saved'))
  } catch (err) {
    editError.value = err.response?.data?.errorMessage || t('agents_page.failed_save')
  } finally {
    saveLoading.value = false
  }
}

async function copyText(text) {
  try {
    await navigator.clipboard.writeText(text)
    copiedField.value = text
    setTimeout(() => {
      if (copiedField.value === text) copiedField.value = null
    }, 2000)
  } catch {
    // clipboard unavailable
  }
}

function openDeleteSessionModal(session) {
  deleteSessionTarget.value = session
  deleteSessionError.value = null
}

function closeDeleteSessionModal() {
  if (!deletingSession.value) {
    deleteSessionTarget.value = null
    deleteSessionError.value = null
  }
}

async function confirmDeleteSession() {
  const session = deleteSessionTarget.value
  if (!session) return

  deletingSession.value = true
  deletingSessionId.value = session.id
  deleteSessionError.value = null

  try {
    await sessionAPI.delete(currentWorkspace.id, agentId.value, { id: session.externalId })
    deleteSessionTarget.value = null

    if (sessionsOffset.value >= sessionsTotal.value - 1 && sessionsOffset.value > 0) {
      goToSessionsOffset(Math.max(0, sessionsOffset.value - sessionsLimit))
    }

    await loadSessions()
    showFlash(t('common.deleted'))
  } catch (err) {
    deleteSessionError.value = err.response?.data?.errorMessage || t('agents_page.failed_delete_session')
  } finally {
    deletingSession.value = false
    deletingSessionId.value = null
  }
}

onMounted(loadAgent)
watch(agentId, async (newId, oldId) => {
  if (oldId && newId !== oldId) {
    await setSessionsPage(1)
  }
  await loadAgent()
})
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.15s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
