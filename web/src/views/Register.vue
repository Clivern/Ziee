<template>
  <div class="min-h-screen flex items-center justify-center bg-theme-bg px-4">
    <div class="max-w-md w-full">
      <div class="text-center mb-10">
        <div class="flex justify-center mb-8">
          <img src="/logo.png" :alt="$t('common.logo_alt')" class="h-24 w-auto">
        </div>
      </div>

      <div class="bg-white rounded-lg border border-theme-border p-8 shadow-sm">
        <h2 class="text-2xl font-semibold text-theme-text mb-2 text-center">{{ $t('register.title') }}</h2>
        <p class="text-sm text-theme-textLight text-center mb-6">
          {{ $t('register.complete_form') }}
        </p>

        <form class="space-y-5" @submit.prevent="handleRegister">
          <div>
            <label for="email" class="block text-sm font-medium text-theme-text mb-2">{{ $t('register.email') }}</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="input-field"
              :placeholder="$t('register.email_placeholder')"
              :disabled="loading"
            >
          </div>
          <div>
            <label for="name" class="block text-sm font-medium text-theme-text mb-2">{{ $t('register.name') }}</label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              required
              class="input-field"
              :placeholder="$t('register.name_placeholder')"
              :disabled="loading"
            >
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-theme-text mb-2">{{ $t('register.password') }}</label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              minlength="8"
              class="input-field"
              :placeholder="$t('register.password_placeholder')"
              :disabled="loading"
            >
          </div>
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-theme-text mb-2">{{ $t('register.confirm_password') }}</label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              required
              minlength="8"
              class="input-field"
              :placeholder="$t('register.confirm_password_placeholder')"
              :disabled="loading"
            >
          </div>

          <div v-if="error" class="rounded-md border border-red-200 bg-red-50 p-3">
            <p class="text-sm text-red-800">{{ error }}</p>
          </div>

          <button
            type="submit"
            class="w-full btn-primary py-2.5 disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="loading"
          >
            <span v-if="!loading">{{ $t('register.submit') }}</span>
            <span v-else class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24" aria-hidden="true">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
              </svg>
              {{ $t('register.submitting') }}
            </span>
          </button>
          <p class="text-center text-sm text-theme-textLight">
            {{ $t('register.already_have_account') }}
            <router-link :to="{ path: '/login', query: redirectQuery }" class="font-medium text-primary-800 hover:text-primary-900">{{ $t('register.sign_in') }}</router-link>
          </p>
        </form>
      </div>

      <p class="text-center text-xs text-theme-textLight mt-8">
        {{ $t('common.copyright') }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { authAPI } from '@/api'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const form = reactive({
  email: '',
  name: '',
  password: '',
  confirmPassword: ''
})
const loading = ref(false)
const error = ref(null)

const redirectQuery = computed(() => {
  return route.query.redirect ? { redirect: route.query.redirect } : {}
})

const handleRegister = async () => {
  if (form.password !== form.confirmPassword) {
    error.value = t('register.error_passwords_match')
    return
  }
  if (form.password.length < 8) {
    error.value = t('register.error_password_length')
    return
  }
  loading.value = true
  error.value = null
  try {
    await authAPI.register({
      email: form.email.trim(),
      name: form.name.trim(),
      password: form.password
    })
    router.push({
      path: '/login',
      query: {
        registered: '1',
        ...redirectQuery.value
      }
    })
  } catch (err) {
    error.value = err.response?.data?.errorMessage || t('register.error_default')
  } finally {
    loading.value = false
  }
}
</script>
