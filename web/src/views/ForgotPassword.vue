<template>
  <div class="min-h-screen flex items-center justify-center bg-theme-bg px-4">
    <div class="max-w-md w-full">
      <div class="text-center mb-10">
        <div class="flex justify-center mb-8">
          <router-link to="/">
            <img src="/logo.png" :alt="$t('common.logo_alt')" class="h-24 w-auto">
          </router-link>
        </div>
      </div>

      <div class="bg-white rounded-lg border border-theme-border p-8 shadow-sm">
        <h2 class="text-2xl font-semibold text-theme-text mb-2 text-center">{{ $t('forgot_password.title') }}</h2>
        <p class="text-sm text-theme-textLight text-center mb-6">
          {{ $t('forgot_password.description') }}
        </p>
        <form class="space-y-5" @submit.prevent="handleSubmit">
          <div>
            <label for="email" class="block text-sm font-medium text-theme-text mb-2">{{ $t('forgot_password.email') }}</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="input-field"
              :placeholder="$t('forgot_password.email_placeholder')"
              :disabled="loading"
            >
          </div>

          <div v-if="error" class="rounded-md border border-red-200 bg-red-50 p-3">
            <p class="text-sm text-red-800">{{ error }}</p>
          </div>

          <div>
            <button
              type="submit"
              class="w-full btn-primary py-2.5 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="loading"
            >
              <span v-if="!loading">{{ $t('forgot_password.submit') }}</span>
              <span v-else class="flex items-center justify-center">
                <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
                {{ $t('forgot_password.submitting') }}
              </span>
            </button>
          </div>
        </form>

        <p class="text-center mt-6">
          <router-link to="/login" class="text-sm text-primary-800 hover:text-primary-900 font-medium">
            {{ $t('forgot_password.back_to_login') }}
          </router-link>
        </p>
      </div>

      <p class="text-center text-xs text-theme-textLight mt-8">
        {{ $t('common.copyright') }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { authAPI } from '@/api'
import { showFlash } from '@/lib/flash'

const { t } = useI18n()
const form = reactive({ email: '' })
const loading = ref(false)
const error = ref(null)

const handleSubmit = async () => {
  loading.value = true
  error.value = null

  try {
    await authAPI.forgotPassword({ email: form.email })
    showFlash(t('common.sent'))
  } catch (err) {
    error.value = err.response?.data?.errorMessage || t('forgot_password.error_default')
  } finally {
    loading.value = false
  }
}
</script>
