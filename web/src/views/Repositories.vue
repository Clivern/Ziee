<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8 flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 class="text-3xl font-semibold text-theme-text">{{ $t('repos_page.title') }}</h1>
          <p class="text-theme-textLight mt-2">{{ $t('repos_page.subtitle') }}</p>
        </div>
        <span class="rounded-full border border-amber-200 bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
          {{ $t('dashboard.mock_badge') }}
        </span>
      </div>

      <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
        <div class="border-b border-theme-border px-6 py-4">
          <h2 class="text-lg font-semibold text-theme-text">{{ $t('repos_page.connected') }}</h2>
          <p class="mt-1 text-sm text-theme-textLight">{{ $t('repos_page.connected_desc') }}</p>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-theme-border">
            <thead class="bg-theme-hover">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('repos_page.repository') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('repos_page.triage') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('repos_page.queue') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('repos_page.activity') }}</th>
                <th class="px-6 py-3 text-right text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('repos_page.actions') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-theme-border bg-white">
              <tr v-for="repo in repositories" :key="repo.id" class="hover:bg-theme-hover/60">
                <td class="px-6 py-4">
                  <router-link :to="`/repositories/${repo.id}`" class="text-sm font-medium text-theme-text hover:underline">
                    {{ repo.fullName }}
                  </router-link>
                  <p class="mt-0.5 text-xs text-theme-textLight">
                    {{ repo.private ? $t('repos_page.private') : $t('repos_page.public') }}
                    · {{ repo.defaultBranch }}
                  </p>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                    :class="repo.triageEnabled ? 'bg-emerald-50 text-emerald-800' : 'bg-theme-hover text-theme-textLight'"
                  >
                    {{ repo.triageEnabled ? $t('repos_page.on') : $t('repos_page.off') }}
                  </span>
                  <p class="mt-1 text-xs text-theme-textLight">{{ $t('repos_page.inbox_count', { count: inboxCount(repo.id) }) }}</p>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                    :class="repo.queueEnabled ? 'bg-emerald-50 text-emerald-800' : 'bg-theme-hover text-theme-textLight'"
                  >
                    {{ repo.queueEnabled ? $t('repos_page.on') : $t('repos_page.off') }}
                  </span>
                  <p class="mt-1 text-xs text-theme-textLight">{{ $t('repos_page.queue_count', { count: queueCount(repo.id) }) }}</p>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-textLight">{{ formatDate(repo.lastActivity) }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                  <router-link :to="`/repositories/${repo.id}`" class="font-medium text-primary-800 hover:underline">
                    {{ $t('repos_page.open') }}
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import AppNav from '@/components/AppNav.vue'
import { getMockIssues, getMockQueue, MOCK_REPOSITORIES } from '@/mocks/repositories'

const repositories = MOCK_REPOSITORIES

function inboxCount(id) {
  return getMockIssues(id).filter((issue) => issue.status === 'inbox').length
}

function queueCount(id) {
  return getMockQueue(id).filter((item) => item.status !== 'blocked').length
}

function formatDate(iso) {
  return new Date(iso).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  })
}
</script>
