<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <h1 class="text-3xl font-semibold text-theme-text">{{ $t('knowledge_page.title') }}</h1>
        <p class="text-theme-textLight mt-2">{{ $t('knowledge_page.subtitle') }}</p>
      </div>

      <div v-if="errorMessage" class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-800">
        {{ errorMessage }}
      </div>

      <div class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
        <div class="border-b border-theme-border px-6 py-4">
          <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <h2 class="text-lg font-semibold text-theme-text">{{ $t('knowledge_page.documents') }}</h2>
              <p class="mt-1 text-sm text-theme-textLight">{{ $t('knowledge_page.documents_desc') }}</p>
            </div>
            <button type="button" class="btn-primary inline-flex items-center justify-center gap-2 shrink-0" @click="fileInputRef?.click()">
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
              </svg>
              {{ $t('knowledge_page.upload') }}
            </button>
          </div>
        </div>

        <div
          class="mx-6 my-4 rounded-lg border-2 border-dashed border-theme-border bg-theme-bg px-6 py-8 text-center transition-colors"
          :class="dragActive ? 'border-primary-500 bg-primary-50/40' : ''"
          @dragenter.prevent="dragActive = true"
          @dragover.prevent="dragActive = true"
          @dragleave.prevent="dragActive = false"
          @drop.prevent="onDropFiles"
        >
          <svg class="mx-auto h-10 w-10 text-theme-textLight" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <p class="mt-3 text-sm font-medium text-theme-text">{{ $t('knowledge_page.dropzone_title') }}</p>
          <p class="mt-1 text-xs text-theme-textLight">{{ $t('knowledge_page.dropzone_hint') }}</p>
          <label class="mt-4 inline-flex cursor-pointer items-center justify-center gap-2 rounded-md border border-theme-border bg-white px-4 py-2 text-sm font-medium text-theme-text hover:bg-theme-hover">
            {{ $t('knowledge_page.browse_files') }}
            <input
              ref="fileInputRef"
              type="file"
              accept=".txt,.md,text/plain,text/markdown"
              multiple
              class="sr-only"
              @change="onFileChange"
            />
          </label>
        </div>

        <div v-if="loadingDocuments" class="p-12 text-center">
          <svg class="animate-spin h-8 w-8 mx-auto text-theme-text" fill="none" viewBox="0 0 24 24" aria-hidden="true">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <p class="text-theme-textLight mt-3">{{ $t('knowledge_page.loading_documents') }}</p>
        </div>

        <div v-else-if="documents.length === 0" class="p-12 text-center">
          <p class="text-theme-textLight">{{ $t('knowledge_page.no_documents') }}</p>
          <p class="text-sm text-theme-textLight mt-1">{{ $t('knowledge_page.no_documents_hint') }}</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-theme-border">
            <thead class="bg-theme-hover">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('knowledge_page.title_col') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('knowledge_page.file') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('knowledge_page.labels') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('knowledge_page.status') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('knowledge_page.created_at') }}</th>
                <th class="px-6 py-3 text-right text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('knowledge_page.actions') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-theme-border bg-white">
              <tr v-for="doc in documents" :key="doc.id" class="hover:bg-theme-hover/60">
                <td class="px-6 py-4 text-sm font-medium text-theme-text">{{ doc.title }}</td>
                <td class="px-6 py-4">
                  <p class="text-sm text-theme-text">{{ doc.filename }}</p>
                  <p class="mt-0.5 text-xs text-theme-textLight">
                    {{ fmtSize(doc.size) }} · {{ $t('knowledge_page.chars', { count: doc.charCount.toLocaleString() }) }}
                  </p>
                </td>
                <td class="px-6 py-4">
                  <span v-if="doc.labels?.length" class="inline-flex flex-wrap gap-1">
                    <span
                      v-for="label in doc.labels"
                      :key="label"
                      class="inline-flex items-center rounded-full bg-theme-hover px-2 py-0.5 text-xs font-mono text-theme-text"
                    >
                      <span class="text-theme-textLight">{{ splitLabel(label).key }}</span><span>=</span><span>{{ splitLabel(label).value }}</span>
                    </span>
                  </span>
                  <span v-else class="text-sm text-theme-textLight">—</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium"
                    :class="statusMeta(doc.status).class"
                  >
                    <svg
                      v-if="doc.status === KNOWLEDGE_STATUS.PROCESSING"
                      class="h-3.5 w-3.5 animate-spin"
                      fill="none"
                      viewBox="0 0 24 24"
                      aria-hidden="true"
                    >
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                    </svg>
                    {{ $t(statusMeta(doc.status).labelKey) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-textLight">{{ fmtDate(doc.uploadedAt) }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                  <button type="button" class="text-red-600 hover:text-red-700 hover:underline" @click="deleteDoc = doc">
                    {{ $t('knowledge_page.delete') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div
          v-if="total > 0"
          class="flex flex-wrap items-center justify-between gap-3 border-t border-theme-border px-6 py-4 text-sm text-theme-textLight"
        >
          <p>{{ $t('knowledge_page.showing', { from: pageStart, to: pageEnd, total }) }}</p>
          <div class="flex items-center gap-2">
            <button type="button" class="btn-secondary py-1.5 text-sm disabled:opacity-50" :disabled="offset === 0 || loadingDocuments" @click="goToPage(offset - limit)">
              {{ $t('knowledge_page.previous') }}
            </button>
            <button type="button" class="btn-secondary py-1.5 text-sm disabled:opacity-50" :disabled="offset + limit >= total || loadingDocuments" @click="goToPage(offset + limit)">
              {{ $t('knowledge_page.next') }}
            </button>
          </div>
        </div>
      </div>
    </main>

    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showUploadModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closeUpload">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-lg p-6 max-h-[90vh] overflow-y-auto">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('knowledge_page.upload_title') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('knowledge_page.upload_desc') }}</p>

            <div v-if="uploadError" class="mt-4 rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-800">{{ uploadError }}</div>

            <div v-if="pendingUploads.length" class="mt-4 space-y-3">
              <p class="text-xs font-medium uppercase tracking-wide text-theme-textLight">{{ $t('knowledge_page.selected_files') }}</p>
              <div
                v-for="(entry, index) in pendingUploads"
                :key="`${entry.file.name}-${index}`"
                class="rounded-lg border border-theme-border p-3 space-y-2"
              >
                <div class="flex items-center justify-between gap-3 text-sm">
                  <span class="truncate text-theme-text">{{ entry.file.name }}</span>
                  <span class="shrink-0 text-theme-textLight">{{ fmtSize(entry.file.size) }}</span>
                </div>
                <div>
                  <label :for="`upload-title-${index}`" class="form-label">{{ $t('knowledge_page.document_title') }}</label>
                  <input
                    :id="`upload-title-${index}`"
                    v-model="entry.title"
                    type="text"
                    required
                    class="input-field"
                    :placeholder="$t('knowledge_page.title_placeholder')"
                  />
                </div>
              </div>
            </div>

            <div class="mt-4">
              <label class="form-label">{{ $t('knowledge_page.labels') }}</label>
              <p class="text-xs text-theme-textLight mb-2">{{ $t('knowledge_page.labels_hint') }}</p>
              <div class="rounded-lg border border-theme-border bg-theme-bg px-3 py-2">
                <div v-if="uploadLabelsList.length" class="mb-2 flex flex-wrap gap-1.5">
                  <span
                    v-for="label in uploadLabelsList"
                    :key="label"
                    class="inline-flex items-center gap-1 rounded-full bg-white px-2 py-0.5 text-xs font-mono text-theme-text border border-theme-border"
                  >
                    <span class="text-theme-textLight">{{ splitLabel(label).key }}</span><span>=</span><span>{{ splitLabel(label).value }}</span>
                    <button type="button" class="text-theme-textLight hover:text-theme-text" :aria-label="$t('knowledge_page.remove_label', { label })" @click="rmLabel(label)">
                      ×
                    </button>
                  </span>
                </div>
                <input
                  v-model="uploadLabelInput"
                  type="text"
                  class="w-full bg-transparent text-sm font-mono text-theme-text placeholder:text-theme-textLight focus:outline-none"
                  :placeholder="$t('knowledge_page.labels_placeholder')"
                  @keydown.enter.prevent="addLabel"
                  @keydown="onLabelKey"
                  @blur="addLabel"
                />
              </div>
              <p v-if="uploadLabelError" class="mt-1.5 text-xs text-red-600">{{ uploadLabelError }}</p>
            </div>

            <div class="mt-6 flex gap-3 justify-end">
              <button type="button" class="btn-secondary" :disabled="uploading" @click="closeUpload">{{ $t('common.cancel') }}</button>
              <button type="button" class="btn-primary" :disabled="uploading || !canConfirmUpload" @click="doUpload">
                {{ uploading ? $t('knowledge_page.uploading') : $t('knowledge_page.upload_confirm') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <Teleport to="body">
      <Transition name="modal">
        <div v-if="deleteDoc" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40" @click.self="closeDelete">
          <div class="bg-white rounded-lg shadow-xl border border-theme-border w-full max-w-sm p-6">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('knowledge_page.delete_title') }}</h2>
            <p class="mt-2 text-sm text-theme-textLight">{{ $t('knowledge_page.delete_confirm', { name: deleteDoc.title }) }}</p>
            <div class="mt-6 flex gap-3 justify-end">
              <button type="button" class="btn-secondary" :disabled="deleting" @click="closeDelete">{{ $t('common.cancel') }}</button>
              <button type="button" class="btn-primary bg-red-600 hover:bg-red-700 focus:ring-red-500" :disabled="deleting" @click="doDelete">
                {{ $t('knowledge_page.delete') }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import AppNav from '@/components/AppNav.vue'
import { showFlash } from '@/lib/flash'
import { useUrlPagination } from '@/lib/pagination'
import { documentAPI } from '@/api'
import { useWorkspaceContext } from '@/lib/permission'

const KNOWLEDGE_STATUS = {
  PROCESSING: 'processing',
  INDEXED: 'indexed',
  FAILED: 'failed',
}

const KNOWLEDGE_STATUS_META = {
  [KNOWLEDGE_STATUS.PROCESSING]: {
    class: 'bg-amber-100 text-amber-800',
    labelKey: 'knowledge_page.status_indexing',
  },
  [KNOWLEDGE_STATUS.FAILED]: {
    class: 'bg-red-100 text-red-800',
    labelKey: 'knowledge_page.status_failed',
  },
  [KNOWLEDGE_STATUS.INDEXED]: {
    class: 'bg-emerald-100 text-emerald-800',
    labelKey: 'knowledge_page.status_indexed',
  },
}

const LABEL_PART_PATTERN = /^[a-z0-9][a-z0-9_.-]*$/
const POLL_INTERVAL_MS = 3000

const { t } = useI18n()
const route = useRoute()
const { currentWorkspace } = useWorkspaceContext()
const { offset, limit, goToOffset, goToPage: setPage } = useUrlPagination({ limit: 10 })

function fmtSize(bytes) {
  if (!bytes || bytes < 1024) return `${bytes || 0} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function parseLabel(raw) {
  const text = String(raw).trim()
  const eq = text.indexOf('=')
  if (eq <= 0 || eq === text.length - 1) return null

  const key = text.slice(0, eq).trim().toLowerCase()
  const value = text.slice(eq + 1).trim().toLowerCase()
  if (!LABEL_PART_PATTERN.test(key) || !LABEL_PART_PATTERN.test(value)) return null

  return { key, value, label: `${key}=${value}` }
}

function titleFrom(filename) {
  const base = String(filename).replace(/\.[^.]+$/, '').trim()
  if (!base) return ''
  return base
    .replace(/[-_]+/g, ' ')
    .replace(/\b\w/g, (char) => char.toUpperCase())
}

function isTextFile(file) {
  const name = file.name.toLowerCase()
  return (
    name.endsWith('.txt')
    || name.endsWith('.md')
    || file.type === 'text/plain'
    || file.type === 'text/markdown'
  )
}

function splitLabel(label) {
  return parseLabel(label) || { key: label, value: '' }
}

function fmtDate(iso) {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })
  } catch {
    return iso
  }
}

function statusMeta(status) {
  return KNOWLEDGE_STATUS_META[status] || KNOWLEDGE_STATUS_META[KNOWLEDGE_STATUS.INDEXED]
}

function normDoc(doc) {
  return {
    id: doc.id,
    title: doc.title,
    filename: doc.filename,
    contentType: doc.contentType,
    checksum: doc.checksum,
    size: doc.size?.value ?? doc.size ?? 0,
    charCount: doc.charCount ?? 0,
    labels: doc.labels || [],
    status: doc.status,
    uploadedAt: doc.createdAt,
  }
}

const uploadLabelsList = ref([])
const uploadLabelInput = ref('')
const uploadLabelError = ref(null)

function resetLabels() {
  uploadLabelsList.value = []
  uploadLabelInput.value = ''
  uploadLabelError.value = null
}

function addLabel() {
  const text = uploadLabelInput.value.trim()
  if (!text) return

  const parsed = parseLabel(text)
  if (!parsed) {
    uploadLabelError.value = t('knowledge_page.invalid_label_format')
    return
  }

  if (!uploadLabelsList.value.includes(parsed.label)) {
    uploadLabelsList.value = [...uploadLabelsList.value, parsed.label]
  }
  uploadLabelInput.value = ''
  uploadLabelError.value = null
}

function rmLabel(label) {
  uploadLabelsList.value = uploadLabelsList.value.filter((item) => item !== label)
}

function onLabelKey(event) {
  if (event.key === ',') {
    event.preventDefault()
    addLabel()
    return
  }
  if (event.key === 'Backspace' && !uploadLabelInput.value && uploadLabelsList.value.length) {
    uploadLabelsList.value = uploadLabelsList.value.slice(0, -1)
    uploadLabelError.value = null
  }
}

const documents = ref([])
const total = ref(0)
const loadingDocuments = ref(false)
const errorMessage = ref(null)
const dragActive = ref(false)
const fileInputRef = ref(null)
let pollTimer = null

const showUploadModal = ref(false)
const pendingUploads = ref([])
const uploading = ref(false)
const uploadError = ref(null)

const deleteDoc = ref(null)
const deleting = ref(false)

const pageStart = computed(() => (total.value ? offset.value + 1 : 0))
const pageEnd = computed(() => Math.min(offset.value + limit, total.value))
const hasProcessingDocuments = computed(() =>
  documents.value.some((doc) => doc.status === KNOWLEDGE_STATUS.PROCESSING),
)
const canConfirmUpload = computed(() =>
  pendingUploads.value.length > 0 && pendingUploads.value.every((entry) => entry.title.trim()),
)

async function loadDocs({ silent = false } = {}) {
  if (!currentWorkspace?.id) return

  if (!silent) loadingDocuments.value = true
  errorMessage.value = null

  try {
    const res = await documentAPI.list(currentWorkspace.id, {
      limit,
      offset: offset.value,
    })
    documents.value = (res.data.documents || []).map(normDoc)
    total.value = res.data._meta?.total ?? documents.value.length
  } catch (err) {
    errorMessage.value = err.response?.data?.errorMessage || t('knowledge_page.failed_load_documents')
    documents.value = []
    total.value = 0
  } finally {
    if (!silent) loadingDocuments.value = false
  }
}

watch(hasProcessingDocuments, (processing) => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
  if (processing) {
    pollTimer = setInterval(() => loadDocs({ silent: true }), POLL_INTERVAL_MS)
  }
})

async function goToPage(newOffset) {
  if (newOffset < 0 || newOffset >= total.value) return
  goToOffset(newOffset)
}

watch(() => route.query.page, () => loadDocs())

function openUpload(files) {
  const textFiles = [...files].filter(isTextFile)
  if (!textFiles.length) {
    uploadError.value = t('knowledge_page.invalid_file_type')
    showUploadModal.value = true
    return
  }

  pendingUploads.value = textFiles.map((file) => ({
    file,
    title: titleFrom(file.name),
  }))
  resetLabels()
  uploadError.value = null
  showUploadModal.value = true
}

function onDropFiles(event) {
  dragActive.value = false
  openUpload(event.dataTransfer?.files || [])
}

function onFileChange(event) {
  openUpload(event.target.files || [])
  event.target.value = ''
}

function resetUpload() {
  showUploadModal.value = false
  pendingUploads.value = []
  resetLabels()
  uploadError.value = null
}

function closeUpload() {
  if (uploading.value) return
  resetUpload()
}

async function doUpload() {
  if (!canConfirmUpload.value || uploading.value || !currentWorkspace?.id) return

  uploading.value = true
  uploadError.value = null

  try {
    const labels = [...uploadLabelsList.value]
    for (const entry of pendingUploads.value) {
      const formData = new FormData()
      formData.append('file', entry.file)
      formData.append('title', entry.title.trim())
      if (labels.length) {
        formData.append('labels', JSON.stringify(labels))
      }
      await documentAPI.upload(currentWorkspace.id, formData)
    }

    resetUpload()
    setPage(1)
    await loadDocs()
    showFlash(t('common.created'))
  } catch (err) {
    uploadError.value = err.response?.data?.errorMessage || t('knowledge_page.upload_failed')
  } finally {
    uploading.value = false
  }
}

function closeDelete() {
  if (deleting.value) return
  deleteDoc.value = null
}

async function doDelete() {
  if (!deleteDoc.value || !currentWorkspace?.id || deleting.value) return

  deleting.value = true
  errorMessage.value = null

  try {
    await documentAPI.delete(currentWorkspace.id, deleteDoc.value.id)
    deleteDoc.value = null

    if (offset.value >= total.value - 1 && offset.value > 0) {
      goToOffset(Math.max(0, offset.value - limit))
    }

    await loadDocs()
    showFlash(t('common.deleted'))
  } catch (err) {
    errorMessage.value = err.response?.data?.errorMessage || t('knowledge_page.upload_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(loadDocs)
onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.15s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
