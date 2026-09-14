<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main v-if="repo" class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-6">
        <router-link to="/repositories" class="text-sm font-medium text-primary-800 hover:underline">
          {{ $t('repo_page.back') }}
        </router-link>
      </div>

      <div class="mb-8 flex flex-wrap items-start justify-between gap-4">
        <div>
          <div class="flex flex-wrap items-center gap-2">
            <h1 class="text-3xl font-semibold text-theme-text">{{ repo.fullName }}</h1>
            <span class="rounded-full bg-theme-hover px-2.5 py-0.5 text-xs font-medium text-theme-textLight">
              {{ repo.private ? $t('repos_page.private') : $t('repos_page.public') }}
            </span>
            <span class="rounded-full border border-amber-200 bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
              {{ $t('dashboard.mock_badge') }}
            </span>
          </div>
          <p class="text-theme-textLight mt-2">{{ $t('repo_page.subtitle') }}</p>
        </div>
        <div class="flex flex-wrap gap-2">
          <button
            type="button"
            class="rounded-full border px-3 py-1.5 text-sm font-medium"
            :class="triageEnabled ? 'border-emerald-200 bg-emerald-50 text-emerald-800' : 'border-theme-border bg-white text-theme-textLight'"
            @click="triageEnabled = !triageEnabled"
          >
            {{ $t('repo_page.triage_toggle') }} · {{ triageEnabled ? $t('repos_page.on') : $t('repos_page.off') }}
          </button>
          <button
            type="button"
            class="rounded-full border px-3 py-1.5 text-sm font-medium"
            :class="queueEnabled ? 'border-emerald-200 bg-emerald-50 text-emerald-800' : 'border-theme-border bg-white text-theme-textLight'"
            @click="queueEnabled = !queueEnabled"
          >
            {{ $t('repo_page.queue_toggle') }} · {{ queueEnabled ? $t('repos_page.on') : $t('repos_page.off') }}
          </button>
        </div>
      </div>

      <section class="mb-8 grid grid-cols-2 gap-4 lg:grid-cols-4" aria-label="Repository stats">
        <div class="stat-card">
          <p class="stat-label">{{ $t('repo_page.stat_inbox') }}</p>
          <p class="stat-value mt-1">{{ inboxCount }}</p>
        </div>
        <div class="stat-card">
          <p class="stat-label">{{ $t('repo_page.stat_drafting') }}</p>
          <p class="stat-value mt-1">{{ draftingCount }}</p>
        </div>
        <div class="stat-card">
          <p class="stat-label">{{ $t('repo_page.stat_queue') }}</p>
          <p class="stat-value mt-1">{{ activeQueueCount }}</p>
        </div>
        <div class="stat-card">
          <p class="stat-label">{{ $t('repo_page.stat_blocked') }}</p>
          <p class="stat-value mt-1">{{ blockedCount }}</p>
        </div>
      </section>

      <div class="grid grid-cols-1 gap-6 xl:grid-cols-2">
        <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
          <div class="border-b border-theme-border px-6 py-4">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('repo_page.triage_title') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('repo_page.triage_desc') }}</p>
          </div>

          <div v-if="!triageEnabled" class="p-8 text-sm text-theme-textLight">
            {{ $t('repo_page.triage_paused') }}
          </div>
          <div v-else>
            <div class="flex flex-wrap gap-2 border-b border-theme-border px-6 py-3">
              <button
                v-for="filter in issueFilters"
                :key="filter"
                type="button"
                class="rounded-full px-3 py-1 text-xs font-medium"
                :class="issueFilter === filter ? 'bg-primary-800 text-white' : 'bg-theme-hover text-theme-text'"
                @click="issueFilter = filter"
              >
                {{ $t(`repo_page.filter_${filter}`) }}
              </button>
            </div>

            <ul class="divide-y divide-theme-border">
              <li v-for="issue in visibleIssues" :key="issue.id" class="px-6 py-4">
                <div class="flex flex-wrap items-start justify-between gap-2">
                  <div class="min-w-0">
                    <p class="text-sm font-medium text-theme-text">
                      <span class="text-theme-textLight">#{{ issue.number }}</span>
                      {{ issue.title }}
                    </p>
                    <p class="mt-1 text-xs text-theme-textLight">
                      {{ issue.author }} · {{ issue.openedAt }}
                      <span v-for="label in issue.labels" :key="label" class="ml-1 rounded bg-theme-hover px-1.5 py-0.5 font-mono">{{ label }}</span>
                    </p>
                  </div>
                  <span
                    class="shrink-0 rounded-full px-2.5 py-0.5 text-xs font-medium"
                    :class="issueStatusClass(issue.status)"
                  >
                    {{ $t(`repo_page.status_${issue.status}`) }}
                  </span>
                </div>
                <p class="mt-2 text-sm text-theme-textLight leading-relaxed">{{ issue.suggestionText }}</p>
                <div v-if="issue.status === 'inbox' || issue.status === 'waiting'" class="mt-3 flex flex-wrap gap-2">
                  <button type="button" class="btn-primary py-1.5 text-xs" @click="setIssueStatus(issue, 'drafting')">
                    {{ $t('repo_page.action_draft') }}
                  </button>
                  <button type="button" class="btn-secondary py-1.5 text-xs" @click="setIssueStatus(issue, 'waiting')">
                    {{ $t('repo_page.action_ask') }}
                  </button>
                  <button type="button" class="btn-secondary py-1.5 text-xs" @click="setIssueStatus(issue, 'dismissed')">
                    {{ $t('repo_page.action_dismiss') }}
                  </button>
                </div>
                <div v-else-if="issue.status === 'ready'" class="mt-3 flex flex-wrap gap-2">
                  <button type="button" class="btn-primary py-1.5 text-xs" @click="setIssueStatus(issue, 'drafting')">
                    {{ $t('repo_page.action_enqueue') }}
                  </button>
                </div>
              </li>
            </ul>
          </div>
        </section>

        <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
          <div class="border-b border-theme-border px-6 py-4">
            <h2 class="text-lg font-semibold text-theme-text">{{ $t('repo_page.queue_title') }}</h2>
            <p class="mt-1 text-sm text-theme-textLight">{{ $t('repo_page.queue_desc') }}</p>
          </div>

          <div v-if="!queueEnabled" class="p-8 text-sm text-theme-textLight">
            {{ $t('repo_page.queue_paused') }}
          </div>
          <div v-else-if="queue.length === 0" class="p-8 text-sm text-theme-textLight">
            {{ $t('repo_page.queue_empty') }}
          </div>
          <ol v-else class="divide-y divide-theme-border">
            <li v-for="(item, index) in queue" :key="item.id" class="flex gap-4 px-6 py-4">
              <span
                class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-sm font-semibold"
                :class="item.status === 'blocked' ? 'bg-amber-50 text-amber-800' : 'bg-theme-hover text-theme-text'"
              >
                {{ item.status === 'blocked' ? '—' : index + 1 }}
              </span>
              <div class="min-w-0 flex-1">
                <div class="flex flex-wrap items-center justify-between gap-2">
                  <p class="text-sm font-medium text-theme-text">
                    <span class="text-theme-textLight">#{{ item.number }}</span>
                    {{ item.title }}
                  </p>
                  <span
                    class="rounded-full px-2.5 py-0.5 text-xs font-medium"
                    :class="queueStatusClass(item.status)"
                  >
                    {{ $t(`repo_page.queue_${item.status}`) }}
                  </span>
                </div>
                <p class="mt-1 text-xs text-theme-textLight">
                  {{ item.author }} · {{ $t('repo_page.checks', { value: item.checks }) }}
                </p>
                <p class="mt-1.5 text-sm text-theme-textLight leading-relaxed">{{ item.note }}</p>
              </div>
            </li>
          </ol>
        </section>
      </div>
    </main>

    <main v-else class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <p class="text-sm text-theme-textLight">{{ $t('repo_page.not_found') }}</p>
      <router-link to="/repositories" class="mt-3 inline-block text-sm font-medium text-primary-800 hover:underline">
        {{ $t('repo_page.back') }}
      </router-link>
    </main>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppNav from '@/components/AppNav.vue'
import { getMockIssues, getMockQueue, getMockRepository } from '@/mocks/repositories'

const route = useRoute()
const issueFilters = ['all', 'inbox', 'drafting', 'waiting', 'ready', 'dismissed']

const repo = computed(() => getMockRepository(route.params.repoId))
const triageEnabled = ref(true)
const queueEnabled = ref(true)
const issueFilter = ref('all')
const issues = ref([])
const queue = ref([])

watch(
  repo,
  (value) => {
    if (!value) {
      issues.value = []
      queue.value = []
      return
    }
    triageEnabled.value = value.triageEnabled
    queueEnabled.value = value.queueEnabled
    issueFilter.value = 'all'
    issues.value = getMockIssues(value.id)
    queue.value = getMockQueue(value.id)
  },
  { immediate: true },
)

const visibleIssues = computed(() => {
  if (issueFilter.value === 'all') return issues.value
  return issues.value.filter((issue) => issue.status === issueFilter.value)
})

const inboxCount = computed(() => issues.value.filter((issue) => issue.status === 'inbox').length)
const draftingCount = computed(() => issues.value.filter((issue) => issue.status === 'drafting').length)
const activeQueueCount = computed(() => queue.value.filter((item) => item.status !== 'blocked').length)
const blockedCount = computed(() => queue.value.filter((item) => item.status === 'blocked').length)

function setIssueStatus(issue, status) {
  issue.status = status
}

function issueStatusClass(status) {
  if (status === 'inbox') return 'bg-sky-50 text-sky-800'
  if (status === 'drafting') return 'bg-violet-50 text-violet-800'
  if (status === 'waiting') return 'bg-amber-50 text-amber-800'
  if (status === 'ready') return 'bg-emerald-50 text-emerald-800'
  return 'bg-theme-hover text-theme-textLight'
}

function queueStatusClass(status) {
  if (status === 'merging') return 'bg-emerald-50 text-emerald-800'
  if (status === 'queued') return 'bg-sky-50 text-sky-800'
  if (status === 'blocked') return 'bg-amber-50 text-amber-800'
  return 'bg-theme-hover text-theme-textLight'
}
</script>
