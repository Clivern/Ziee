<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <h1 class="text-3xl font-semibold text-theme-text">{{ $t('agents_page.title') }}</h1>
        <p class="text-theme-textLight mt-2">{{ $t('agents_page.subtitle') }}</p>
      </div>

      <div v-if="errorMessage" class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
        {{ errorMessage }}
      </div>

      <div class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
        <div class="border-b border-theme-border px-6 py-4 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.list_title') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('agents_page.list_desc') }}</p>
          </div>
          <button
            v-if="canEdit"
            type="button"
            class="btn-primary inline-flex items-center justify-center gap-2 shrink-0"
            @click="openCreateModal"
          >
            <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            {{ $t('agents_page.create') }}
          </button>
        </div>

        <div v-if="loading" class="p-12 text-center">
          <p class="text-theme-textLight">{{ $t('agents_page.loading') }}</p>
        </div>

        <div v-else-if="agents.length === 0" class="p-12 text-center">
          <p class="text-theme-textLight">{{ $t('agents_page.no_agents') }}</p>
          <p v-if="canEdit" class="text-sm text-theme-textLight mt-1">{{ $t('agents_page.no_agents_hint') }}</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-theme-border">
            <thead class="bg-theme-hover">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.name') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.created_at') }}</th>
                <th class="px-6 py-3 text-right text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.actions') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-theme-border bg-white">
              <tr v-for="agent in agents" :key="agent.id" class="hover:bg-theme-hover/60">
                <td class="px-6 py-4">
                  <router-link :to="agentDetailPath(agent.id)" class="text-sm font-medium text-primary-600 hover:text-primary-700 hover:underline">
                    {{ agent.name }}
                  </router-link>
                  <p v-if="agent.description" class="mt-0.5 text-xs text-theme-textLight line-clamp-1">{{ agent.description }}</p>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-textLight">{{ formatAgentDate(agent.createdAt) }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                  <span class="inline-flex items-center gap-3">
                    <router-link :to="agentDetailPath(agent.id)" class="text-primary-600 hover:text-primary-700 hover:underline">
                      {{ $t('agents_page.view') }}
                    </router-link>
                    <button
                      v-if="canEdit"
                      type="button"
                      class="text-red-600 hover:text-red-700 hover:underline disabled:opacity-50"
                      :disabled="deletingId === agent.id"
                      @click="openDeleteModal(agent)"
                    >
                      {{ deletingId === agent.id ? $t('agents_page.deleting') : $t('agents_page.delete') }}
                    </button>
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-if="total > 0" class="bg-white px-6 py-4 border-t border-theme-border flex items-center justify-between">
          <div class="text-sm text-theme-textLight">
            {{ $t('agents_page.showing', { from: offset + 1, to: Math.min(offset + limit, total), total }) }}
          </div>
          <div class="flex items-center gap-2">
            <button type="button" class="btn-secondary text-sm disabled:opacity-50" :disabled="offset === 0" @click="goToPage(offset - limit)">
              {{ $t('agents_page.previous') }}
            </button>
            <button type="button" class="btn-secondary text-sm disabled:opacity-50" :disabled="offset + limit >= total" @click="goToPage(offset + limit)">
              {{ $t('agents_page.next') }}
            </button>
          </div>
        </div>
      </div>
    </main>

    <!-- Create modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closeCreateModal">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-lg p-6 max-h-[90vh] overflow-y-auto">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.create_title') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('agents_page.create_desc') }}</p>
            <div v-if="createError" class="mt-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ createError }}</div>
            <form class="mt-4 space-y-4" @submit.prevent="handleCreate">
              <div>
                <label for="agent-name" class="form-label">{{ $t('agents_page.name') }}</label>
                <input id="agent-name" v-model="form.name" type="text" required class="input-field" :placeholder="$t('agents_page.name_placeholder')" />
              </div>
              <div>
                <label for="agent-description" class="form-label">{{ $t('agents_page.description') }}</label>
                <textarea id="agent-description" v-model="form.description" rows="2" required class="input-field" :placeholder="$t('agents_page.description_placeholder')" />
              </div>
              <div class="flex gap-3 justify-end">
                <button type="button" class="btn-secondary" :disabled="createLoading" @click="closeCreateModal">{{ $t('common.cancel') }}</button>
                <button type="submit" class="btn-primary" :disabled="createLoading">
                  {{ createLoading ? $t('agents_page.creating') : $t('agents_page.create') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Delete modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="deleteTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closeDeleteModal">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.delete_title') }}</h2>
            <div v-if="deleteError" class="mt-3 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ deleteError }}</div>
            <p class="mt-2 text-sm text-theme-textLight">{{ $t('agents_page.delete_confirm', { name: deleteTarget.name }) }}</p>
            <div class="mt-6 flex gap-3 justify-end">
              <button type="button" class="btn-secondary" :disabled="deleting" @click="closeDeleteModal">{{ $t('common.cancel') }}</button>
              <button type="button" class="btn-primary bg-red-600 hover:bg-red-700 focus:ring-red-500" :disabled="deleting" @click="confirmDelete">
                {{ deleting ? $t('agents_page.deleting') : $t('agents_page.delete') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppNav from '@/components/AppNav.vue'
import { showFlash } from '@/lib/flash'
import { useUrlPagination } from '@/lib/pagination'
import { useWorkspaceContext } from '@/lib/permission'
import { agentAPI } from '@/api'
import { agentDetailPath, formatAgentDate, buildCreatePayload } from '@/lib/agent'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { currentWorkspace, canManage: canEdit } = useWorkspaceContext()
const { offset, limit, goToOffset } = useUrlPagination({ limit: 50 })

const agents = ref([])
const loading = ref(false)
const total = ref(0)
const errorMessage = ref(null)

const showCreateModal = ref(false)
const createLoading = ref(false)
const createError = ref(null)
const form = ref({ name: '', description: '' })

const deleteTarget = ref(null)
const deleteError = ref(null)
const deleting = ref(false)
const deletingId = ref(null)

async function loadPage() {
  loading.value = true
  errorMessage.value = null
  try {
    const res = await agentAPI.list(currentWorkspace.id, { limit, offset: offset.value })
    agents.value = res.data.agents || []
    total.value = res.data._meta?.total ?? 0
  } catch (err) {
    errorMessage.value = err.response?.data?.errorMessage || t('agents_page.failed_load')
    agents.value = []
  } finally {
    loading.value = false
  }
}

function goToPage(newOffset) {
  if (newOffset < 0 || newOffset >= total.value) return
  goToOffset(newOffset)
}

watch(() => route.query.page, loadPage)

function openCreateModal() {
  form.value = { name: '', description: '' }
  createError.value = null
  showCreateModal.value = true
}

function closeCreateModal() {
  if (createLoading.value) return
  showCreateModal.value = false
}

async function handleCreate() {
  const name = form.value.name.trim()
  const description = form.value.description.trim()
  if (!name) {
    createError.value = t('agents_page.name_required')
    return
  }
  if (!description) {
    createError.value = t('agents_page.description_required')
    return
  }

  createLoading.value = true
  createError.value = null
  try {
    const created = await agentAPI.create(currentWorkspace.id, buildCreatePayload(form.value))
    showCreateModal.value = false
    showFlash(t('common.created'))
    router.push(agentDetailPath(created.data.id))
  } catch (err) {
    createError.value = err.response?.data?.errorMessage || t('agents_page.failed_create')
  } finally {
    createLoading.value = false
  }
}

function openDeleteModal(agent) {
  deleteTarget.value = agent
  deleteError.value = null
}

function closeDeleteModal() {
  if (!deleting.value) {
    deleteTarget.value = null
    deleteError.value = null
  }
}

async function confirmDelete() {
  const agent = deleteTarget.value
  if (!agent) return

  deleting.value = true
  deletingId.value = agent.id
  deleteError.value = null
  try {
    await agentAPI.delete(currentWorkspace.id, agent.id)
    deleteTarget.value = null
    await loadPage()
    showFlash(t('common.deleted'))
  } catch (err) {
    deleteError.value = err.response?.data?.errorMessage || t('agents_page.failed_delete')
  } finally {
    deleting.value = false
    deletingId.value = null
  }
}

onMounted(loadPage)
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.15s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
