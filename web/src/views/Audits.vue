<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <h1 class="text-3xl font-semibold text-theme-text">{{ $t('audits_page.title') }}</h1>
        <p class="text-theme-textLight mt-2">{{ $t('audits_page.subtitle') }}</p>
      </div>

      <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden opacity-70 max-w-2xl">
        <div class="p-6">
          <div class="flex items-start gap-4">
            <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg bg-theme-hover text-theme-textLight">
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
              </svg>
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <h2 class="text-base font-semibold text-theme-text">{{ $t('audits_page.event_log') }}</h2>
                <span class="rounded-full bg-theme-hover px-2 py-0.5 text-xs font-medium text-theme-textLight">
                  {{ $t('audits_page.coming_soon') }}
                </span>
              </div>
              <p class="mt-1.5 text-sm text-theme-textLight">{{ $t('audits_page.coming_soon_desc') }}</p>
            </div>
            <button
              type="button"
              class="btn-secondary shrink-0 text-sm py-1.5 px-3 opacity-60 cursor-not-allowed"
              disabled
            >
              {{ $t('audits_page.coming_soon') }}
            </button>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import AppNav from '@/components/AppNav.vue'
</script>
