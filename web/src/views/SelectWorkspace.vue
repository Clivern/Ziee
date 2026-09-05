<template>
  <div class="min-h-screen bg-theme-bg flex flex-col items-center justify-center px-4 py-12">
    <div class="w-full max-w-md">
      <div class="flex justify-center mb-6">
        <div class="flex h-14 w-14 items-center justify-center rounded-full bg-primary-200">
          <svg class="h-8 w-8 text-theme-text" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
          </svg>
        </div>
      </div>

      <h1 class="text-xl font-semibold text-theme-text text-center">
        {{ workspaces.length ? $t('select_workspace.select_title') : $t('select_workspace.create_title') }}
      </h1>
      <p class="mt-2 text-sm text-theme-textLight text-center">
        {{ workspaces.length ? $t('select_workspace.select_description') : $t('select_workspace.no_workspaces') }}
      </p>

      <div v-if="installations.length" class="mt-6">
        <h2 class="text-sm font-semibold text-theme-text">
          {{ $t('select_workspace.installations_title') }}
        </h2>
        <p class="mt-1 text-sm text-theme-textLight">
          {{ $t('select_workspace.installations_description') }}
        </p>
        <div class="mt-3 space-y-3">
          <div
            v-for="installation in installations"
            :key="installation.id"
            class="flex items-center gap-3 rounded-lg border border-theme-border bg-white shadow-sm px-4 py-3"
          >
            <div class="min-w-0 flex-1">
              <p class="text-sm font-semibold text-theme-text truncate">{{ installation.accountLogin }}</p>
              <p class="text-xs text-theme-textLight mt-0.5">
                {{ installation.accountType }} · {{ installation.repositorySelection }}
              </p>
            </div>
            <button
              type="button"
              class="btn-primary py-1.5 px-3 text-sm shrink-0"
              :disabled="attachingId === installation.id"
              @click="openAttachModal(installation)"
            >
              {{ attachingId === installation.id ? $t('select_workspace.attaching') : $t('select_workspace.attach') }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="workspaces.length" class="mt-6 space-y-3">
        <div
          v-for="workspace in workspaces"
          :key="workspace.id"
          class="flex items-center gap-2 rounded-lg border border-theme-border bg-white shadow-sm overflow-hidden"
        >
          <button
            type="button"
            @click="selectWorkspace(workspace)"
            class="flex-1 flex items-center justify-between px-4 py-3 text-left hover:bg-theme-hover transition-colors focus:outline-none focus:ring-2 focus:ring-primary-800 focus:ring-inset min-w-0"
          >
            <div class="flex min-w-0 flex-1 items-center gap-3">
              <div class="min-w-0 flex-1">
                <p class="text-sm font-semibold text-theme-text truncate">{{ workspace.name }}</p>
                <p class="text-xs text-theme-textLight mt-0.5">{{ formatDate(workspace.createdAt) }}</p>
              </div>
              <span
                v-if="workspace.status === 'active' || storedWorkspace?.id === workspace.id"
                class="inline-flex flex-shrink-0 rounded-full bg-primary-200 px-2.5 py-0.5 text-xs font-medium text-theme-text"
              >
                {{ $t('select_workspace.active') }}
              </span>
            </div>
            <svg class="ml-3 h-5 w-5 flex-shrink-0 text-theme-textLight" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
      </div>

      <div class="mt-6">
        <button
          type="button"
          @click="openCreateModal"
          :class="[
            'w-full flex items-center justify-center gap-2 rounded-lg px-4 py-3 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-primary-800 focus:ring-offset-1 transition-colors',
            workspaces.length ? 'border border-theme-border bg-white text-theme-text hover:bg-theme-hover' : 'border-2 border-dashed border-primary-400 bg-primary-50/50 text-primary-800 hover:bg-primary-100'
          ]"
        >
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          {{ workspaces.length ? $t('select_workspace.create_new') : $t('select_workspace.create_first') }}
        </button>
      </div>

      <p v-if="error" class="mt-4 text-sm text-red-600 text-center">{{ error }}</p>
    </div>

    <Teleport to="body">
      <Transition name="modal">
        <div
          v-show="showCreateModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
          role="dialog"
          aria-modal="true"
          aria-labelledby="create-workspace-title"
          @click.self="closeCreateModal"
        >
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 id="create-workspace-title" class="text-lg font-semibold text-theme-text">{{ $t('select_workspace.create_modal_title') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('select_workspace.create_modal_description') }}</p>
            <div v-if="createModalError" class="mt-4 rounded-md border border-red-200 bg-red-50 p-3" role="alert">
              <p class="text-sm text-red-800">{{ createModalError }}</p>
            </div>
            <form @submit.prevent="handleCreate" class="mt-4 space-y-4">
              <div>
                <label for="workspace-name" class="form-label">{{ $t('select_workspace.workspace_name') }}</label>
                <input
                  id="workspace-name"
                  v-model="newWorkspaceName"
                  type="text"
                  required
                  class="input-field"
                  :placeholder="$t('select_workspace.workspace_name_placeholder')"
                  autofocus
                />
              </div>
              <div class="flex gap-3 justify-end">
                <button type="button" @click="closeCreateModal" class="btn-secondary">{{ $t('common.cancel') }}</button>
                <button type="submit" :disabled="createLoading || !newWorkspaceName.trim()" class="btn-primary">
                  {{ createLoading ? $t('common.loading') : $t('common.create') }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <Teleport to="body">
      <Transition name="modal">
        <div
          v-show="showAttachModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40"
          role="dialog"
          aria-modal="true"
          aria-labelledby="attach-workspace-title"
          @click.self="closeAttachModal"
        >
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 id="attach-workspace-title" class="text-lg font-semibold text-theme-text">
              {{ $t('select_workspace.attach_modal_title') }}
            </h2>
            <p class="mt-1 text-sm text-theme-textLight">
              {{ $t('select_workspace.attach_modal_description', { account: pendingInstallation?.accountLogin }) }}
            </p>
            <div v-if="attachModalError" class="mt-4 rounded-md border border-red-200 bg-red-50 p-3" role="alert">
              <p class="text-sm text-red-800">{{ attachModalError }}</p>
            </div>
            <div class="mt-4 space-y-2">
              <button
                v-for="workspace in workspaces"
                :key="workspace.id"
                type="button"
                class="w-full flex items-center justify-between rounded-lg border px-4 py-3 text-left text-sm transition-colors"
                :class="attachWorkspaceId === workspace.id
                  ? 'border-primary-800 bg-primary-50 text-theme-text'
                  : 'border-theme-border bg-white text-theme-text hover:bg-theme-hover'"
                @click="attachWorkspaceId = workspace.id"
              >
                <span class="font-medium truncate">{{ workspace.name }}</span>
                <span
                  v-if="attachWorkspaceId === workspace.id"
                  class="ml-3 h-2.5 w-2.5 shrink-0 rounded-full bg-primary-800"
                />
              </button>
            </div>
            <div class="mt-5 flex gap-3 justify-end">
              <button type="button" class="btn-secondary" @click="closeAttachModal">
                {{ $t('common.cancel') }}
              </button>
              <button
                type="button"
                class="btn-primary"
                :disabled="!attachWorkspaceId || !!attachingId"
                @click="confirmAttach"
              >
                {{ attachingId ? $t('select_workspace.attaching') : $t('select_workspace.attach') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { github_api, workspace_api } from '@/api'
import { loadWorkspaceFromStorage, saveWorkspaceToStorage } from '@/utils/storage'

const { t } = useI18n()
const router = useRouter()
const storedWorkspace = loadWorkspaceFromStorage()
const workspaces = ref([])
const installations = ref([])
const pendingInstallation = ref(null)
const attachWorkspaceId = ref('')
const attachingId = ref('')
const error = ref(null)
const showCreateModal = ref(false)
const showAttachModal = ref(false)
const attachModalError = ref(null)
const createModalError = ref(null)
const newWorkspaceName = ref('')
const createLoading = ref(false)

function formatDate(iso) {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
  } catch {
    return iso
  }
}

function openCreateModal() {
  createModalError.value = null
  showCreateModal.value = true
}

function closeCreateModal() {
  showCreateModal.value = false
  createModalError.value = null
  newWorkspaceName.value = ''
}

function openAttachModal(installation) {
  if (!workspaces.value.length) {
    openCreateModal()
    return
  }

  pendingInstallation.value = installation
  attachModalError.value = null
  attachWorkspaceId.value = storedWorkspace?.id || workspaces.value[0]?.id || ''
  showAttachModal.value = true
}

function closeAttachModal() {
  showAttachModal.value = false
  attachModalError.value = null
  pendingInstallation.value = null
  attachingId.value = ''
}

function selectWorkspace(workspace) {
  saveWorkspaceToStorage(workspace)
  router.push('/dashboard')
}

async function loadWorkspaces() {
  error.value = null
  try {
    const res = await workspace_api.list({})
    workspaces.value = res.data?.workspaces ?? []
  } catch (err) {
    error.value = err.response?.data?.errorMessage || t('select_workspace.failed_load')
    workspaces.value = []
  }
}

async function loadInstallations() {
  try {
    const res = await github_api.listInstallations()
    installations.value = res.data?.installations ?? []
  } catch (err) {
    error.value = err.response?.data?.errorMessage || t('select_workspace.failed_load_installations')
    installations.value = []
  }
}

async function confirmAttach() {
  const installation = pendingInstallation.value
  attachingId.value = installation.id
  attachModalError.value = null

  try {
    await github_api.attachInstallation(installation.id, { workspaceId: attachWorkspaceId.value })
    installations.value = installations.value.filter((item) => item.id !== installation.id)
    closeAttachModal()
  } catch (err) {
    attachModalError.value = err.response?.data?.errorMessage || t('select_workspace.failed_attach')
    attachingId.value = ''
  }
}

async function handleCreate() {
  createLoading.value = true
  createModalError.value = null

  try {
    const res = await workspace_api.create({ name: newWorkspaceName.value })
    if (res.data?.id) {
      const workspace = { ...res.data, role: 'owner', status: 'active' }
      saveWorkspaceToStorage(workspace)
      closeCreateModal()
      await loadWorkspaces()
      if (!installations.value.length) {
        router.push('/dashboard')
      }
    } else {
      createModalError.value = 'Invalid response'
    }
  } catch (err) {
    createModalError.value = err.response?.data?.errorMessage || 'Failed to create workspace'
  } finally {
    createLoading.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadWorkspaces(), loadInstallations()])
  if (workspaces.value.length === 0 && !error.value) openCreateModal()
})
</script>
