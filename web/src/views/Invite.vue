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
        {{ $t('invite.title') }}
      </h1>
      <p class="mt-2 text-sm text-theme-textLight text-center">
        {{ $t('invite.description') }}
      </p>

      <div class="mt-6 rounded-lg border border-theme-border bg-white shadow-sm overflow-hidden">
        <div v-if="loading" class="p-8 text-center">
          <svg class="animate-spin h-8 w-8 mx-auto text-theme-text" fill="none" viewBox="0 0 24 24" aria-hidden="true">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
          <p class="text-theme-textLight mt-3">{{ $t('invite.loading') }}</p>
        </div>

        <div v-else-if="error" class="p-6">
          <div class="rounded-md border border-red-200 bg-red-50 p-3" role="alert">
            <p class="text-sm text-red-800">{{ error }}</p>
          </div>
          <button
            type="button"
            @click="router.push('/select-workspace')"
            class="mt-4 w-full btn-secondary"
          >
            {{ $t('invite.back_to_workspaces') }}
          </button>
        </div>

        <div v-else-if="invite" class="p-6">
          <div class="space-y-4">
            <div>
              <p class="text-xs font-medium uppercase tracking-wide text-theme-textLight">{{ $t('invite.email') }}</p>
              <p class="mt-1 text-sm font-semibold text-theme-text break-words">{{ invite.email }}</p>
            </div>
            <div>
              <p class="text-xs font-medium uppercase tracking-wide text-theme-textLight">{{ $t('invite.workspace') }}</p>
              <p class="mt-1 text-sm font-semibold text-theme-text break-words">{{ invite.workspaceName || invite.workspaceId }}</p>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-xs font-medium uppercase tracking-wide text-theme-textLight">{{ $t('invite.role') }}</p>
                <p class="mt-1 text-sm font-semibold text-theme-text capitalize">{{ invite.role }}</p>
              </div>
              <div>
                <p class="text-xs font-medium uppercase tracking-wide text-theme-textLight">{{ $t('invite.expires') }}</p>
                <p class="mt-1 text-sm font-semibold text-theme-text">{{ formatDate(invite.expiresAt) }}</p>
              </div>
            </div>
          </div>

          <div v-if="actionError" class="mt-5 rounded-md border border-red-200 bg-red-50 p-3" role="alert">
            <p class="text-sm text-red-800">{{ actionError }}</p>
          </div>

          <div class="mt-6 grid grid-cols-1 gap-3 sm:grid-cols-2">
            <button
              type="button"
              @click="rejectInvite"
              class="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="!!actionLoading"
            >
              {{ actionLoading === 'reject' ? $t('invite.rejecting') : $t('invite.reject') }}
            </button>
            <button
              type="button"
              @click="acceptInvite"
              class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="!!actionLoading"
            >
              {{ actionLoading === 'accept' ? $t('invite.accepting') : $t('invite.accept') }}
            </button>
          </div>
        </div>

        <div v-else class="p-6">
          <div class="rounded-md border border-red-200 bg-red-50 p-3" role="alert">
            <p class="text-sm text-red-800">{{ $t('invite.not_found') }}</p>
          </div>
          <button
            type="button"
            @click="router.push('/select-workspace')"
            class="mt-4 w-full btn-secondary"
          >
            {{ $t('invite.back_to_workspaces') }}
          </button>
        </div>
      </div>

      <p class="text-center text-xs text-theme-textLight mt-8">
        {{ $t('common.copyright') }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { invite_api } from '@/api'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const invite = ref(null)
const loading = ref(true)
const actionLoading = ref(null)
const error = ref(null)
const actionError = ref(null)

function isValidInvite(data) {
  return data && typeof data === 'object' && !Array.isArray(data) && typeof data.email === 'string'
}

function formatDate(iso) {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
  } catch {
    return iso
  }
}

async function loadInvite() {
  if (!route.params.token) {
    error.value = t('invite.not_found')
    loading.value = false
    return
  }

  loading.value = true
  error.value = null
  invite.value = null

  try {
    const res = await invite_api.getByToken(route.params.token)
    if (!isValidInvite(res.data)) {
      error.value = t('invite.not_found')
      return
    }
    invite.value = res.data
  } catch (err) {
    const status = err.response?.status
    const message = err.response?.data?.errorMessage
    error.value = (status === 401 || status === 403 || status === 404)
      ? (message || t('invite.not_found'))
      : (message || t('invite.failed_load'))
  } finally {
    loading.value = false
  }
}

async function acceptInvite() {
  actionLoading.value = 'accept'
  actionError.value = null

  try {
    await invite_api.acceptByToken(route.params.token)
    router.push('/select-workspace')
  } catch (err) {
    actionError.value = err.response?.data?.errorMessage || t('invite.failed_accept')
  } finally {
    actionLoading.value = null
  }
}

async function rejectInvite() {
  actionLoading.value = 'reject'
  actionError.value = null

  try {
    await invite_api.rejectByToken(route.params.token)
    router.push('/select-workspace')
  } catch (err) {
    actionError.value = err.response?.data?.errorMessage || t('invite.failed_reject')
  } finally {
    actionLoading.value = null
  }
}

onMounted(loadInvite)
</script>
