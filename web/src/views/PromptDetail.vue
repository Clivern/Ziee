<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-6">
        <router-link to="/prompts" class="text-sm text-primary-600 hover:text-primary-700 hover:underline">
          ← {{ $t('prompts_page.back_to_prompts') }}
        </router-link>
      </div>

      <div v-if="errorMessage" class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
        {{ errorMessage }}
      </div>

      <div class="mb-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-theme-text">{{ promptName }}</h1>
          <p v-if="promptHandle" class="mt-1 text-sm font-mono text-theme-textLight">{{ promptHandle }}</p>
          <p v-if="promptDescription" class="mt-2 text-sm text-theme-textLight">{{ promptDescription }}</p>
          <p v-else class="text-sm text-theme-textLight mt-1">{{ $t('prompts_page.detail_subtitle') }}</p>
        </div>
        <button
          v-if="canEdit"
          type="button"
          class="btn-primary self-start"
          @click="openNewVersionModal"
        >
          {{ $t('prompts_page.new_version') }}
        </button>
      </div>

      <div v-if="loading" class="p-12 text-center text-theme-textLight">{{ $t('prompts_page.loading') }}</div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-[280px_minmax(0,1fr)] gap-6">
        <!-- Version sidebar -->
        <aside class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
          <div class="border-b border-theme-border px-4 py-3 text-sm font-medium text-theme-text">
            {{ $t('prompts_page.versions') }}
          </div>
          <ul class="divide-y divide-theme-border max-h-[70vh] overflow-y-auto">
            <li v-for="version in versions" :key="version.id">
              <button
                type="button"
                class="w-full text-left px-4 py-3 hover:bg-theme-hover transition-colors"
                :class="selectedVersion?.id === version.id ? 'bg-sky-50 border-l-2 border-l-primary-600' : ''"
                @click="selectVersion(version)"
              >
                <div class="flex items-center justify-between gap-2">
                  <span class="text-sm font-medium text-theme-text"># {{ version.version }}</span>
                  <div class="flex flex-wrap gap-1 justify-end">
                  <span
                    v-if="version.production"
                    class="inline-flex rounded-full px-1.5 py-0.5 text-[10px] font-medium bg-emerald-100 text-emerald-800"
                  >
                    production
                  </span>
                  <span
                    v-for="label in parsePromptLabels(version.labels)"
                    :key="label"
                    class="inline-flex rounded-full px-1.5 py-0.5 text-[10px] font-medium"
                    :class="labelClass(label)"
                  >
                    {{ label }}
                  </span>
                  </div>
                </div>
                <p class="mt-1 text-xs text-theme-textLight">{{ formatPromptDateTime(version.createdAt) }}</p>
              </button>
            </li>
          </ul>
        </aside>

        <!-- Version detail -->
        <section v-if="selectedVersion" class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
          <div class="border-b border-theme-border px-6 py-4 flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 class="text-lg font-semibold text-theme-text">
                # {{ selectedVersion.version }}
              </h2>
              <p class="text-sm text-theme-textLight">{{ formatPromptDateTime(selectedVersion.createdAt) }}</p>
            </div>
            <div v-if="canEdit" class="flex flex-wrap items-center gap-3">
              <button
                v-if="!selectedVersion.production"
                type="button"
                class="btn-secondary"
                @click="openPromoteProductionModal(selectedVersion)"
              >
                {{ $t('prompts_page.promote_production') }}
              </button>
              <button
                type="button"
                class="text-sm text-red-600 hover:text-red-700 hover:underline"
                @click="openDeleteVersionModal(selectedVersion)"
              >
                {{ $t('prompts_page.delete_version') }}
              </button>
            </div>
          </div>

          <div class="border-b border-theme-border px-6">
            <nav class="flex gap-6">
              <button
                type="button"
                class="py-3 text-sm font-medium border-b-2 -mb-px"
                :class="activeTab === 'prompt' ? 'border-primary-600 text-primary-600' : 'border-transparent text-theme-textLight hover:text-theme-text'"
                @click="activeTab = 'prompt'"
              >
                {{ $t('prompts_page.tab_prompt') }}
              </button>
              <button
                type="button"
                class="py-3 text-sm font-medium border-b-2 -mb-px"
                :class="activeTab === 'config' ? 'border-primary-600 text-primary-600' : 'border-transparent text-theme-textLight hover:text-theme-text'"
                @click="activeTab = 'config'"
              >
                {{ $t('prompts_page.tab_config') }}
              </button>
            </nav>
          </div>

          <div class="p-6">
            <div v-if="formatVersionCommit(selectedVersion)" class="mb-4">
              <p class="font-mono text-xs text-theme-text">{{ formatVersionCommit(selectedVersion) }}</p>
            </div>

            <div v-show="activeTab === 'prompt'">
              <pre class="rounded-lg border border-theme-border bg-theme-hover p-4 text-sm text-theme-text whitespace-pre-wrap font-mono">{{ selectedVersion.content }}</pre>
              <p v-if="promptVariables.length" class="mt-3 text-xs text-theme-textLight">
                {{ $t('prompts_page.variables_available') }}
                <span v-for="v in promptVariables" :key="v" class="ml-1 font-mono text-theme-text">{{ v }}</span>
              </p>
            </div>

            <div v-show="activeTab === 'config'">
              <pre class="rounded-lg border border-theme-border bg-theme-hover p-4 text-sm text-theme-text whitespace-pre-wrap font-mono">{{ formattedConfig }}</pre>
            </div>
          </div>
        </section>
      </div>
    </main>

    <!-- New version modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showVersionModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closeVersionModal">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-2xl p-6 max-h-[90vh] overflow-y-auto">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('prompts_page.new_version_title', { name: promptName }) }}</h2>
            <div v-if="versionModalError" class="mt-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ versionModalError }}</div>
            <form @submit.prevent="handleSaveVersion" class="mt-4 space-y-4">
              <PromptFormFields v-model="form" :is-new-version="true" />
              <div class="flex gap-3 justify-end">
                <button type="button" class="btn-secondary" @click="closeVersionModal" :disabled="versionLoading">{{ $t('common.cancel') }}</button>
                <button type="submit" class="btn-primary" :disabled="versionLoading">{{ versionLoading ? $t('prompts_page.creating') : $t('prompts_page.save_version') }}</button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Promote to production modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-show="promoteVersion" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closePromoteProductionModal">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('prompts_page.promote_production_title') }}</h2>
            <div v-if="promoteProductionError" class="mt-3 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ promoteProductionError }}</div>
            <p class="mt-2 text-sm text-theme-textLight">
              {{ $t('prompts_page.promote_production_confirm', { name: promptName, version: promoteVersion?.version }) }}
            </p>
            <div class="mt-6 flex gap-3 justify-end">
              <button type="button" class="btn-secondary" :disabled="promotingProduction" @click="closePromoteProductionModal">{{ $t('common.cancel') }}</button>
              <button type="button" class="btn-primary" :disabled="promotingProduction" @click="confirmPromoteProduction">
                {{ promotingProduction ? $t('prompts_page.promoting_production') : $t('prompts_page.promote_production') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Delete version modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-show="deleteVersion" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closeDeleteVersionModal">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('prompts_page.delete_version_title') }}</h2>
            <div v-if="deleteVersionError" class="mt-3 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ deleteVersionError }}</div>
            <p class="mt-2 text-sm text-theme-textLight">{{ $t('prompts_page.delete_version_confirm', { name: promptName, version: deleteVersion?.version }) }}</p>
            <div class="mt-4">
              <label class="block text-sm font-medium text-theme-text mb-1.5">{{ $t('prompts_page.type_to_confirm', { name: promptName }) }}</label>
              <input v-model="deleteConfirmName" type="text" class="input-field" :placeholder="promptName" autocomplete="off" :disabled="deletingVersion" />
            </div>
            <div class="mt-6 flex gap-3 justify-end">
              <button type="button" @click="closeDeleteVersionModal" class="btn-secondary" :disabled="deletingVersion">{{ $t('common.cancel') }}</button>
              <button type="button" @click="confirmDeleteVersion" class="btn-primary bg-red-600 hover:bg-red-700 focus:ring-red-500" :disabled="deletingVersion || !canConfirmDeleteVersion">{{ deletingVersion ? $t('prompts_page.deleting') : $t('prompts_page.delete') }}</button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppNav from '@/components/AppNav.vue'
import { showFlash } from '@/lib/flash'
import { useWorkspaceContext } from '@/lib/permission'
import PromptFormFields from '@/components/PromptFormFields.vue'
import { promptAPI } from '@/api'
import {
  emptyPromptForm,
  formFromPrompt,
  buildCreatePayload,
  formatPromptDateTime,
  formatVersionCommit,
  formatPromptConfig,
  serializePromptConfigForApi,
  parsePromptLabels,
  extractPromptVariables,
  hasPromptLabel,
} from '@/lib/prompt'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const { currentWorkspace, canManage: canEdit } = useWorkspaceContext()

const promptId = computed(() => route.params.promptId || '')
const promptName = computed(() => selectedVersion.value?.name || versions.value[0]?.name || '')
const promptHandle = computed(() => selectedVersion.value?.handle || versions.value[0]?.handle || '')
const promptDescription = computed(() => selectedVersion.value?.description || versions.value[0]?.description || '')
const versions = ref([])
const selectedVersion = ref(null)
const loading = ref(false)
const errorMessage = ref(null)
const activeTab = ref('prompt')

const showVersionModal = ref(false)
const versionLoading = ref(false)
const versionModalError = ref(null)
const form = ref(emptyPromptForm())

const deleteVersion = ref(null)
const deleteConfirmName = ref('')
const deletingVersion = ref(false)
const deleteVersionError = ref(null)
const promoteVersion = ref(null)
const promotingProduction = ref(false)
const promoteProductionError = ref(null)

const promptVariables = computed(() => extractPromptVariables(selectedVersion.value?.content))
const formattedConfig = computed(() => formatPromptConfig(selectedVersion.value?.config))

const canConfirmDeleteVersion = computed(() => {
  return promptName.value && deleteConfirmName.value.trim() === promptName.value
})

function labelClass(label) {
  if (label === 'production') return 'bg-emerald-100 text-emerald-800'
  if (label === 'latest') return 'bg-sky-100 text-sky-800'
  return 'bg-theme-hover text-theme-text'
}

function selectVersion(version) {
  selectedVersion.value = version
  activeTab.value = 'prompt'
}

async function loadVersions(selectLatest = false) {
  if (!promptId.value) return
  loading.value = true
  errorMessage.value = null
  try {
    const res = await promptAPI.listVersions(currentWorkspace.id, promptId.value, { limit: 100, offset: 0 })
    versions.value = res.data.versions || []
    if (versions.value.length === 0) {
      router.replace('/prompts')
      return
    }
    const currentId = selectedVersion.value?.id
    const keep = versions.value.find((v) => v.id === currentId)
    selectedVersion.value = keep || (selectLatest ? versions.value[0] : versions.value[0])
  } catch (err) {
    if (err.response?.status === 404) {
      router.replace('/prompts')
      return
    }
    errorMessage.value = err.response?.data?.errorMessage || t('prompts_page.failed_load')
  } finally {
    loading.value = false
  }
}

function openNewVersionModal() {
  form.value = formFromPrompt(selectedVersion.value)
  versionModalError.value = null
  showVersionModal.value = true
}

function closeVersionModal() {
  if (versionLoading.value) return
  showVersionModal.value = false
}

async function handleSaveVersion() {
  const content = form.value.content.trim()
  if (!content) {
    versionModalError.value = t('prompts_page.content_required')
    return
  }

  versionLoading.value = true
  versionModalError.value = null
  try {
    const created = await promptAPI.createVersion(currentWorkspace.id, promptId.value, buildCreatePayload(form.value, { omitName: true }))
    showVersionModal.value = false
    await loadVersions(true)
    selectedVersion.value = versions.value.find((v) => v.id === created.data.id) || versions.value[0]
    showFlash(t('common.saved'))
  } catch (err) {
    versionModalError.value = err.response?.data?.errorMessage || t('prompts_page.failed_create')
  } finally {
    versionLoading.value = false
  }
}

function openDeleteVersionModal(version) {
  deleteVersion.value = version
  deleteConfirmName.value = ''
  deleteVersionError.value = null
}

function closeDeleteVersionModal() {
  if (!deletingVersion.value) {
    deleteVersion.value = null
    deleteConfirmName.value = ''
    deleteVersionError.value = null
  }
}

function openPromoteProductionModal(version) {
  promoteVersion.value = version
  promoteProductionError.value = null
}

function closePromoteProductionModal() {
  if (!promotingProduction.value) {
    promoteVersion.value = null
    promoteProductionError.value = null
  }
}

async function confirmPromoteProduction() {
  const version = promoteVersion.value
  if (!version || version.production) return

  promotingProduction.value = true
  promoteProductionError.value = null
  errorMessage.value = null
  try {
    await promptAPI.updateVersion(currentWorkspace.id, promptId.value, version.id, {
      type: version.type,
      content: version.content,
      config: serializePromptConfigForApi(version.config),
      production: true,
    })
    promoteVersion.value = null
    await loadVersions()
    selectedVersion.value = versions.value.find((v) => v.id === version.id) || selectedVersion.value
    showFlash(t('prompts_page.flash_promoted'))
  } catch (err) {
    promoteProductionError.value = err.response?.data?.errorMessage || t('prompts_page.failed_promote_production')
  } finally {
    promotingProduction.value = false
  }
}

async function confirmDeleteVersion() {
  const version = deleteVersion.value
  if (!version) return
  deletingVersion.value = true
  deleteVersionError.value = null
  try {
    await promptAPI.deleteVersion(currentWorkspace.id, promptId.value, version.id)
    deleteVersion.value = null
    deleteConfirmName.value = ''
    await loadVersions(true)
    if (versions.value.length === 0) {
      router.push('/prompts')
      return
    }
    showFlash(t('common.deleted'))
  } catch (err) {
    deleteVersionError.value = err.response?.data?.errorMessage || t('prompts_page.failed_delete')
  } finally {
    deletingVersion.value = false
  }
}

watch(() => route.params.promptId, () => {
  selectedVersion.value = null
  loadVersions(true)
})

onMounted(() => loadVersions(true))
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.15s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
