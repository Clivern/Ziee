<template>
  <div class="min-h-screen flex items-center justify-center bg-theme-bg px-4 py-12">
    <div class="max-w-md w-full">
      <div class="text-center mb-10">
        <router-link to="/" class="inline-flex justify-center">
          <img src="/logo.png" :alt="$t('common.logo_alt')" class="h-20 w-auto">
        </router-link>
      </div>

      <div class="text-center mb-8">
        <h1 class="text-2xl font-semibold text-theme-text tracking-tight">
          {{ title }}
        </h1>
        <p class="mt-2 text-sm text-theme-textLight leading-relaxed">
          {{ subtitle }}
        </p>
        <p v-if="installationId" class="mt-2 text-xs font-mono text-theme-textLight">
          {{ $t('getting_started.installation_id', { id: installationId }) }}
        </p>
      </div>

      <div class="bg-white rounded-lg border border-theme-border p-8 shadow-sm">
        <ol class="space-y-5">
          <li
            v-for="(step, index) in steps"
            :key="step.title"
            class="flex gap-4"
          >
            <span
              class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary-200 text-sm font-semibold text-theme-text"
            >
              {{ index + 1 }}
            </span>
            <div class="min-w-0 pt-0.5">
              <p class="text-sm font-semibold text-theme-text">{{ step.title }}</p>
              <p class="mt-1 text-sm text-theme-textLight leading-relaxed">{{ step.description }}</p>
            </div>
          </li>
        </ol>

        <button
          type="button"
          class="mt-8 w-full btn-primary py-2.5"
          @click="continueSetup"
        >
          {{ ctaLabel }}
        </button>
      </div>

      <p class="text-center text-xs text-theme-textLight mt-8">
        {{ $t('common.copyright') }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { user } from '@/lib/auth'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const installationId = computed(() => {
  const id = route.query.installation_id
  return typeof id === 'string' ? id : ''
})

const setupAction = computed(() => {
  const action = route.query.setup_action
  return typeof action === 'string' ? action : ''
})

const isUpdate = computed(() => setupAction.value === 'update')

const title = computed(() => {
  if (isUpdate.value) return t('getting_started.title_update')
  if (setupAction.value === 'install' || installationId.value) return t('getting_started.title_install')
  return t('getting_started.title')
})

const subtitle = computed(() => {
  if (isUpdate.value) return t('getting_started.subtitle_update')
  if (setupAction.value === 'install' || installationId.value) return t('getting_started.subtitle_install')
  return t('getting_started.subtitle')
})

const steps = computed(() => [
  {
    title: t('getting_started.step_sign_in'),
    description: t('getting_started.step_sign_in_desc'),
  },
  {
    title: t('getting_started.step_workspace'),
    description: t('getting_started.step_workspace_desc'),
  },
  {
    title: t('getting_started.step_merge'),
    description: t('getting_started.step_merge_desc'),
  },
])

const ctaLabel = computed(() => {
  return user.value ? t('getting_started.cta_continue') : t('getting_started.cta_sign_in')
})

function continueSetup() {
  if (user.value) {
    router.push({
      path: '/select-workspace',
      query: installationId.value ? { installation_id: installationId.value } : {},
    })
    return
  }

  router.push({
    path: '/login',
    query: { redirect: route.fullPath },
  })
}
</script>
