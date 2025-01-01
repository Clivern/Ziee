<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <h1 class="text-3xl font-semibold text-theme-text">{{ $t('integrations_page.title') }}</h1>
        <p class="text-theme-textLight mt-2">{{ $t('integrations_page.subtitle') }}</p>
      </div>

      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <div
          v-for="item in integrationTypes"
          :key="item.type"
          class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden opacity-70"
        >
          <div class="p-6">
            <div class="flex items-start gap-4">
              <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg bg-theme-hover text-theme-textLight">
                <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="item.icon" />
                </svg>
              </div>
              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-2">
                  <h2 class="text-base font-semibold text-theme-text">{{ $t(item.labelKey) }}</h2>
                  <span class="rounded-full bg-theme-hover px-2 py-0.5 text-xs font-medium text-theme-textLight">
                    {{ $t('integrations_page.coming_soon') }}
                  </span>
                </div>
                <p class="mt-1.5 text-sm text-theme-textLight">{{ $t(item.descKey) }}</p>
              </div>
              <button
                type="button"
                class="btn-secondary shrink-0 text-sm py-1.5 px-3 opacity-60 cursor-not-allowed"
                disabled
              >
                {{ $t('integrations_page.coming_soon') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import AppNav from '@/components/AppNav.vue'
import { INTEGRATION_TYPES } from '@/mocks/integrations'

const integrationTypes = [
  {
    type: INTEGRATION_TYPES.WEBHOOK,
    labelKey: 'integrations_page.type_webhook',
    descKey: 'integrations_page.type_webhook_desc',
    icon: 'M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1',
  },
  {
    type: 'slack',
    labelKey: 'integrations_page.type_slack',
    descKey: 'integrations_page.type_slack_desc',
    icon: 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z',
  },
  {
    type: 'zapier',
    labelKey: 'integrations_page.type_zapier',
    descKey: 'integrations_page.type_zapier_desc',
    icon: 'M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z',
  },
]
</script>
