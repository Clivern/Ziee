<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <h1 class="text-3xl font-semibold text-theme-text">{{ $t('prompts_page.title') }}</h1>
        <p class="text-theme-textLight mt-2">{{ $t('prompts_page.subtitle') }}</p>
      </div>

      <div v-if="errorMessage" class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
        {{ errorMessage }}
      </div>

      <div class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
        <div v-if="canEdit" class="flex justify-end border-b border-theme-border bg-white px-6 py-4">
          <button
            type="button"
            @click="openCreateModal"
            class="btn-primary inline-flex items-center justify-center gap-2"
          >
            <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            {{ $t('prompts_page.new_prompt') }}
          </button>
        </div>

        <div v-if="loading" class="p-12 text-center">
          <p class="text-theme-textLight">{{ $t('prompts_page.loading') }}</p>
        </div>

        <div v-else-if="prompts.length === 0" class="p-12 text-center">
          <p class="text-theme-textLight">{{ $t('prompts_page.no_prompts') }}</p>
          <p v-if="canEdit" class="text-sm text-theme-textLight mt-1">{{ $t('prompts_page.no_prompts_hint') }}</p>
        </div>

        <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-theme-border">
          <thead class="bg-theme-hover">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('prompts_page.name') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('prompts_page.handle') }}</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('prompts_page.versions') }}</th>
              <th class="px-6 py-3 text-right text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('prompts_page.actions') }}</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-theme-border">
            <tr v-for="prompt in prompts" :key="prompt.promptId" class="hover:bg-theme-hover">
              <td class="px-6 py-4">
                <router-link :to="promptDetailPath(prompt.promptId)" class="text-sm font-medium text-primary-600 hover:text-primary-700 hover:underline">
                  {{ prompt.name }}
                </router-link>
                <p v-if="prompt.description" class="mt-0.5 text-xs text-theme-textLight line-clamp-1">{{ prompt.description }}</p>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-theme-textLight">{{ prompt.handle }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-textLight">
                {{ prompt.versionCount }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                <span class="inline-flex items-center gap-3">
                  <router-link :to="promptDetailPath(prompt.promptId)" class="text-primary-600 hover:text-primary-700 hover:underline">
                    {{ $t('prompts_page.view') }}
                  </router-link>
                  <button
                    v-if="canEdit"
                    type="button"
                    @click="openDeleteModal(prompt)"
                    class="text-red-600 hover:text-red-700 hover:underline disabled:opacity-50"
                    :disabled="deletingName === prompt.name"
                  >
                    {{ deletingName === prompt.name ? $t('prompts_page.deleting') : $t('prompts_page.delete') }}
                  </button>
                </span>
              </td>
            </tr>
          </tbody>
        </table>
        </div>

        <div v-if="total > 0" class="bg-white px-6 py-4 border-t border-theme-border flex items-center justify-between">
          <div class="text-sm text-theme-textLight">
            {{ $t('prompts_page.showing', { from: offset + 1, to: Math.min(offset + limit, total), total }) }}
          </div>
          <div class="flex items-center gap-2">
            <button type="button" @click="goToPage(offset - limit)" :disabled="offset === 0" class="btn-secondary text-sm disabled:opacity-50">{{ $t('prompts_page.previous') }}</button>
            <button type="button" @click="goToPage(offset + limit)" :disabled="offset + limit >= total" class="btn-secondary text-sm disabled:opacity-50">{{ $t('prompts_page.next') }}</button>
          </div>
        </div>
      </div>
    </main>

    <!-- Create modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closeCreateModal">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-2xl p-6 max-h-[90vh] overflow-y-auto">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('prompts_page.create_title') }}</h2>
            <div v-if="createModalError" class="mt-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ createModalError }}</div>
            <form @submit.prevent="handleCreate" class="mt-4 space-y-4">
              <div>
                <label for="create-prompt-name" class="form-label">{{ $t('prompts_page.name') }}</label>
                <input
                  id="create-prompt-name"
                  v-model="form.name"
                  type="text"
                  required
                  class="input-field"
                  :placeholder="$t('prompts_page.name_placeholder')"
                />
                <p class="mt-1.5 text-xs text-theme-textLight">
                  <template v-if="createHandlePreview">{{ $t('prompts_page.handle_preview', { handle: createHandlePreview }) }}</template>
                  <template v-else>{{ $t('prompts_page.handle_preview_empty') }}</template>
                </p>
              </div>
              <div>
                <label for="create-prompt-description" class="form-label">{{ $t('prompts_page.description') }}</label>
                <textarea
                  id="create-prompt-description"
                  v-model="form.description"
                  rows="2"
                  class="input-field"
                  :placeholder="$t('prompts_page.description_placeholder')"
                />
              </div>
              <div>
                <label for="create-prompt-type" class="form-label">{{ $t('prompts_page.type') }}</label>
                <select id="create-prompt-type" v-model="form.type" class="input-field">
                  <option value="text">{{ $t('prompts_page.type_text') }}</option>
                  <option value="chat">{{ $t('prompts_page.type_chat') }}</option>
                </select>
              </div>
              <div>
                <label for="create-prompt-content" class="form-label">{{ $t('prompts_page.content') }}</label>
                <p class="text-xs text-theme-textLight mb-1.5">
                  {{ $t('prompts_page.variables_hint_prefix') }}
                  <span class="font-mono text-theme-text">{{ PROMPT_VARIABLE_EXAMPLE }}</span>.
                  {{ $t('prompts_page.variables_hint_suffix') }}
                </p>
                <textarea
                  id="create-prompt-content"
                  v-model="form.content"
                  required
                  rows="8"
                  class="input-field font-mono text-sm"
                  :placeholder="`${$t('prompts_page.content_placeholder')}\n${PROMPT_VARIABLE_EXAMPLE}`"
                />
                <p v-if="createFormVariables.length" class="mt-2 text-xs text-theme-textLight">
                  {{ $t('prompts_page.variables_available') }}
                  <span v-for="v in createFormVariables" :key="v" class="ml-1 font-mono text-theme-text">{{ v }}</span>
                </p>
              </div>
              <div>
                <label for="create-prompt-config" class="form-label">{{ $t('prompts_page.config') }}</label>
                <p class="text-xs text-theme-textLight mb-1.5">{{ $t('prompts_page.config_hint') }}</p>
                <textarea
                  id="create-prompt-config"
                  v-model="form.config"
                  rows="4"
                  class="input-field font-mono text-sm"
                  placeholder="{}"
                />
              </div>
              <label class="flex items-center gap-2 text-sm text-theme-text cursor-pointer">
                <input
                  v-model="form.production"
                  type="checkbox"
                  class="rounded border-theme-border text-primary-600 focus:ring-primary-800"
                />
                {{ $t('prompts_page.set_production') }}
              </label>
              <div>
                <label for="create-prompt-commit" class="form-label">{{ $t('prompts_page.commit_message') }}</label>
                <textarea
                  id="create-prompt-commit"
                  v-model="form.commitMessage"
                  rows="2"
                  class="input-field"
                  :placeholder="$t('prompts_page.commit_message_placeholder')"
                />
              </div>
              <div class="flex gap-3 justify-end">
                <button type="button" class="btn-secondary" @click="closeCreateModal" :disabled="createLoading">{{ $t('common.cancel') }}</button>
                <button type="submit" class="btn-primary" :disabled="createLoading">{{ createLoading ? $t('prompts_page.creating') : $t('prompts_page.create') }}</button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Delete modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-show="deleteModalPrompt" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closeDeleteModal">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('prompts_page.delete_title') }}</h2>
            <div v-if="deleteModalError" class="mt-3 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ deleteModalError }}</div>
            <p class="mt-2 text-sm text-theme-textLight">{{ $t('prompts_page.delete_all_confirm', { name: deleteModalPrompt?.name, count: deleteModalPrompt?.versionCount }) }}</p>
            <div class="mt-4">
              <label for="delete-confirm-prompt" class="block text-sm font-medium text-theme-text mb-1.5">{{ $t('prompts_page.type_to_confirm', { name: deleteModalPrompt?.name }) }}</label>
              <input id="delete-confirm-prompt" v-model="deleteConfirmName" type="text" class="input-field" :placeholder="deleteModalPrompt?.name" autocomplete="off" :disabled="deleting" />
            </div>
            <div class="mt-6 flex gap-3 justify-end">
              <button type="button" @click="closeDeleteModal" class="btn-secondary" :disabled="deleting">{{ $t('common.cancel') }}</button>
              <button type="button" @click="confirmDelete" class="btn-primary bg-red-600 hover:bg-red-700 focus:ring-red-500" :disabled="deleting || !canConfirmDelete">{{ deleting ? $t('prompts_page.deleting') : $t('prompts_page.delete') }}</button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppNav from '@/components/AppNav.vue'
import { showFlash } from '@/lib/flash'
import { useUrlPagination } from '@/lib/pagination'
import { useWorkspaceContext } from '@/lib/permission'
import { promptAPI } from '@/api'
import {
  emptyPromptForm,
  buildCreatePayload,
  promptDetailPath,
  handleFromName,
  extractPromptVariables,
  PROMPT_VARIABLE_EXAMPLE,
} from '@/lib/prompt'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { currentWorkspace, canManage: canEdit } = useWorkspaceContext()
const { offset, limit, goToOffset } = useUrlPagination({ limit: 50 })

const prompts = ref([])
const loading = ref(false)
const total = ref(0)
const errorMessage = ref(null)

const showCreateModal = ref(false)
const createLoading = ref(false)
const createModalError = ref(null)
const form = ref(emptyPromptForm())

const deleteModalPrompt = ref(null)
const deleteConfirmName = ref('')
const deleting = ref(false)
const deletingName = ref(null)
const deleteModalError = ref(null)

const canConfirmDelete = computed(() => {
  const name = deleteModalPrompt.value?.name
  return name && deleteConfirmName.value.trim() === name
})

const createFormVariables = computed(() => extractPromptVariables(form.value.content))
const createHandlePreview = computed(() => handleFromName(form.value.name))

async function loadPage() {
  loading.value = true
  errorMessage.value = null
  try {
    const res = await promptAPI.list(currentWorkspace.id, { limit, offset: offset.value })
    prompts.value = res.data.prompts || []
    total.value = res.data._meta?.total ?? 0
  } catch (err) {
    errorMessage.value = err.response?.data?.errorMessage || t('prompts_page.failed_load')
    prompts.value = []
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
  form.value = emptyPromptForm()
  createModalError.value = null
  showCreateModal.value = true
}

function closeCreateModal() {
  if (createLoading.value) return
  showCreateModal.value = false
}

async function handleCreate() {
  const name = form.value.name.trim()
  const content = form.value.content.trim()
  if (!name) {
    createModalError.value = t('prompts_page.name_required')
    return
  }
  if (!content) {
    createModalError.value = t('prompts_page.content_required')
    return
  }

  createLoading.value = true
  createModalError.value = null
  try {
    const created = await promptAPI.create(currentWorkspace.id, buildCreatePayload(form.value))
    showCreateModal.value = false
    await loadPage()
    showFlash(t('common.saved'))
    router.push(promptDetailPath(created.data.promptId))
  } catch (err) {
    createModalError.value = err.response?.data?.errorMessage || t('prompts_page.failed_create')
  } finally {
    createLoading.value = false
  }
}

function openDeleteModal(prompt) {
  deleteModalPrompt.value = prompt
  deleteConfirmName.value = ''
  deleteModalError.value = null
}

function closeDeleteModal() {
  if (!deleting.value) {
    deleteModalPrompt.value = null
    deleteConfirmName.value = ''
    deleteModalError.value = null
  }
}

async function confirmDelete() {
  const prompt = deleteModalPrompt.value
  if (!prompt) return
  deleting.value = true
  deletingName.value = prompt.name
  deleteModalError.value = null
  try {
    await promptAPI.deletePrompt(currentWorkspace.id, prompt.promptId)
    deleteModalPrompt.value = null
    deleteConfirmName.value = ''
    await loadPage()
    showFlash(t('common.deleted'))
  } catch (err) {
    deleteModalError.value = err.response?.data?.errorMessage || t('prompts_page.failed_delete')
  } finally {
    deleting.value = false
    deletingName.value = null
  }
}

onMounted(loadPage)
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.15s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
