<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full px-4 sm:px-6 lg:px-8 py-8">
      <header class="mb-8">
        <h1 class="text-2xl sm:text-3xl font-semibold text-theme-text tracking-tight">
          {{ $t('dashboard.title') }}
        </h1>
        <p class="mt-1 text-sm text-theme-textLight">
          {{ $t('dashboard.subtitle') }}
        </p>
        <p v-if="user?.email" class="mt-2 text-sm text-theme-text">
          {{ $t('dashboard.welcome_back') }}
          <span class="font-medium">{{ user.name || user.email }}</span>
        </p>
      </header>

      <section class="mb-8" aria-label="Key metrics">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div class="stat-card">
            <p class="stat-label">{{ $t('dashboard.api_calls_month') }}</p>
            <p class="stat-value mt-1">{{ loading ? '…' : stats.apiCallsMonth.toLocaleString() }}</p>
          </div>
          <div class="stat-card">
            <p class="stat-label">{{ $t('dashboard.documents_stored') }}</p>
            <p class="stat-value mt-1">{{ loading ? '…' : stats.documentsStored.toLocaleString() }}</p>
          </div>
        </div>
      </section>

      <section class="mb-8 grid grid-cols-1 gap-6 lg:grid-cols-2">
        <div class="section overflow-hidden !p-0">
          <div class="border-b border-theme-border px-6 py-4">
            <div class="flex flex-wrap items-center justify-between gap-2">
              <div>
                <h2 class="text-lg font-semibold text-theme-text">{{ $t('dashboard.attention_title') }}</h2>
                <p class="mt-0.5 text-sm text-theme-textLight">{{ $t('dashboard.attention_subtitle') }}</p>
              </div>
              <span class="rounded-full border border-amber-200 bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
                {{ $t('dashboard.mock_badge') }}
              </span>
            </div>
          </div>

          <ul class="divide-y divide-theme-border max-h-[28rem] overflow-y-auto">
            <li
              v-for="item in attentionItems"
              :key="item.id"
              class="flex gap-4 px-6 py-4"
            >
              <span
                class="mt-1.5 h-2 w-2 shrink-0 rounded-full"
                :class="attentionDotClass(item.level)"
                aria-hidden="true"
              />
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-theme-text">{{ item.title }}</p>
                <p class="mt-0.5 text-sm text-theme-textLight leading-relaxed">{{ item.description }}</p>
                <router-link
                  v-if="item.to"
                  :to="item.to"
                  class="mt-2 inline-flex text-sm font-medium text-primary-800 hover:underline"
                >
                  {{ item.action }}
                </router-link>
              </div>
            </li>
          </ul>
        </div>

        <div class="section overflow-hidden !p-0">
          <div class="border-b border-theme-border px-6 py-4">
            <div class="flex flex-wrap items-center justify-between gap-2">
              <div>
                <h2 class="text-lg font-semibold text-theme-text">{{ $t('dashboard.activity_title') }}</h2>
                <p class="mt-0.5 text-sm text-theme-textLight">{{ $t('dashboard.activity_subtitle') }}</p>
              </div>
              <span class="rounded-full border border-amber-200 bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
                {{ $t('dashboard.mock_badge') }}
              </span>
            </div>
          </div>

          <ul class="divide-y divide-theme-border">
            <li
              v-for="event in recentActivity"
              :key="event.id"
              class="flex gap-3 px-6 py-3.5"
            >
              <span
                class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-theme-hover text-theme-textLight"
                aria-hidden="true"
              >
                <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.75" :d="activityIcons[event.icon]" />
                </svg>
              </span>
              <div class="min-w-0 flex-1">
                <p class="text-sm text-theme-text">
                  <span class="font-medium">{{ event.actor }}</span>
                  {{ event.action }}
                  <span class="font-medium">{{ event.target }}</span>
                </p>
                <p class="mt-0.5 text-xs text-theme-textLight">{{ event.time }}</p>
              </div>
            </li>
          </ul>

          <div class="border-t border-theme-border px-6 py-3">
            <router-link to="/audits" class="text-sm font-medium text-primary-800 hover:underline">
              {{ $t('dashboard.activity_view_all') }}
            </router-link>
          </div>
        </div>
      </section>

      <section class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <div class="space-y-6 lg:col-span-2">
          <div
            v-if="showWhatsNew"
            class="section overflow-hidden !p-0"
          >
            <div class="border-b border-theme-border px-6 py-4">
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div class="min-w-0">
                  <div class="flex flex-wrap items-center gap-2">
                    <h2 class="text-lg font-semibold text-theme-text">{{ $t('dashboard.whats_new_title') }}</h2>
                    <span class="rounded-full border border-amber-200 bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
                      {{ $t('dashboard.mock_badge') }}
                    </span>
                  </div>
                  <p class="mt-0.5 text-sm text-theme-textLight">{{ $t('dashboard.whats_new_subtitle') }}</p>
                </div>
                <button
                  type="button"
                  class="rounded-md px-2 py-1 text-sm text-theme-textLight transition hover:bg-theme-hover hover:text-theme-text"
                  @click="dismissWhatsNew"
                >
                  {{ $t('dashboard.whats_new_dismiss') }}
                </button>
              </div>
            </div>
            <div class="px-6 py-5">
              <p class="text-xs font-medium uppercase tracking-wide text-theme-textLight">
                {{ $t('dashboard.whats_new_tip_label') }}
              </p>
              <p class="mt-2 text-base font-medium text-theme-text">{{ $t('dashboard.whats_new_tip_title') }}</p>
              <p class="mt-1 text-sm leading-relaxed text-theme-textLight">
                {{ $t('dashboard.whats_new_tip_body') }}
              </p>
            </div>
          </div>
        </div>

        <div class="space-y-6">
          <div class="section h-fit overflow-hidden">
            <div class="mx-auto mb-4 flex h-20 w-full items-center justify-center" aria-hidden="true">
              <svg class="h-20 w-auto" viewBox="0 0 120 96" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="8" y="16" width="72" height="64" rx="10" fill="#F5F0EB" stroke="#D4C4B0" stroke-width="1.5"/>
                <rect x="8" y="16" width="72" height="18" rx="10" fill="#37352F"/>
                <circle cx="24" cy="25" r="3" fill="#E8DDD0"/>
                <circle cx="36" cy="25" r="3" fill="#E8DDD0"/>
                <rect x="20" y="44" width="20" height="16" rx="4" fill="#FFFFFF" stroke="#D4C4B0" stroke-width="1.2"/>
                <rect x="48" y="44" width="20" height="16" rx="4" fill="#FFFFFF" stroke="#D4C4B0" stroke-width="1.2"/>
                <rect x="20" y="66" width="20" height="8" rx="3" fill="#E8DDD0"/>
                <rect x="48" y="66" width="20" height="8" rx="3" fill="#37352F" opacity="0.85"/>
                <path d="M88 34c8 0 14 6 14 14v10c0 8-6 14-14 14h-2l-6 8v-8h-6c-8 0-14-6-14-14V48c0-8 6-14 14-14h14z" fill="#FFFFFF" stroke="#37352F" stroke-width="1.5"/>
                <circle cx="82" cy="52" r="2" fill="#37352F"/>
                <circle cx="92" cy="52" r="2" fill="#37352F"/>
                <circle cx="102" cy="52" r="2" fill="#37352F"/>
              </svg>
            </div>
            <h2 class="section-title">{{ $t('dashboard.book_call_title') }}</h2>
            <p class="text-sm text-theme-textLight leading-relaxed">
              {{ $t('dashboard.book_call_desc') }}
            </p>
            <a
              :href="calendlyUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="btn-primary mt-4 inline-flex w-full items-center justify-center"
            >
              {{ $t('dashboard.book_call_cta') }}
            </a>
          </div>

          <div class="section h-fit">
            <h2 class="section-title">{{ $t('dashboard.workspace_info') }}</h2>
            <dl class="space-y-0">
              <div class="flex justify-between gap-4 py-3 border-b border-theme-border">
                <dt class="text-sm text-theme-textLight">{{ $t('dashboard.workspace_name') }}</dt>
                <dd class="text-sm font-medium text-theme-text text-right">{{ workspace?.name || '—' }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-3 border-b border-theme-border">
                <dt class="text-sm text-theme-textLight">{{ $t('dashboard.workspace_handle') }}</dt>
                <dd class="text-sm font-medium text-theme-text text-right font-mono">{{ workspace?.handle || '—' }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-3 border-b border-theme-border">
                <dt class="text-sm text-theme-textLight">{{ $t('dashboard.workspace_role') }}</dt>
                <dd class="text-sm font-medium text-theme-text capitalize text-right">{{ roleLabel }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-3 border-b border-theme-border">
                <dt class="text-sm text-theme-textLight">{{ $t('dashboard.workspace_plan') }}</dt>
                <dd class="text-sm font-medium text-theme-text capitalize text-right">{{ planLabel }}</dd>
              </div>
              <div class="flex justify-between gap-4 py-3">
                <dt class="text-sm text-theme-textLight">{{ $t('dashboard.workspace_members') }}</dt>
                <dd class="text-sm font-medium text-theme-text text-right">{{ memberCount }}</dd>
              </div>
            </dl>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  statsAPI,
  billingAPI,
  workspaceAPI,
} from '@/api'
import AppNav from '@/components/AppNav.vue'
import { user } from '@/lib/auth'
import { saveWorkspaceToStorage } from '@/utils/storage'
import { useWorkspaceContext } from '@/lib/permission'

const { t } = useI18n()

const calendlyUrl = 'https://calendly.com/'
const WHATS_NEW_KEY = 'dashboard_whats_new_v1'

const { currentWorkspace, canManage } = useWorkspaceContext()
const workspace = ref({ ...currentWorkspace })
const loading = ref(true)
const plan = ref('hobby')
const showWhatsNew = ref(localStorage.getItem(WHATS_NEW_KEY) !== '1')

function dismissWhatsNew() {
  showWhatsNew.value = false
  localStorage.setItem(WHATS_NEW_KEY, '1')
}

const stats = ref({
  apiCallsMonth: 0,
  documentsStored: 0,
})

const PLAN_LABELS = {
  hobby: 'Hobby',
  starter: 'Starter',
  growth: 'Growth',
  pro: 'Pro',
}

const planLabel = computed(() => PLAN_LABELS[plan.value] || plan.value || 'Hobby')

const roleLabel = computed(() => workspace.value?.role || '—')

const memberCount = computed(() => {
  const count = workspace.value?.membersCount
  if (count == null) return loading.value ? '…' : '—'
  return count.toLocaleString()
})

const attentionItems = computed(() => [
  {
    id: 'quota',
    level: 'warn',
    title: t('dashboard.attention_quota_title'),
    description: t('dashboard.attention_quota_desc'),
    action: t('dashboard.attention_quota_action'),
    to: '/billing',
  },
  {
    id: 'doc',
    level: 'error',
    title: t('dashboard.attention_doc_title'),
    description: t('dashboard.attention_doc_desc'),
    action: t('dashboard.attention_doc_action'),
    to: '/knowledge',
  },
  {
    id: 'key',
    level: 'error',
    title: t('dashboard.attention_key_title'),
    description: t('dashboard.attention_key_desc'),
    action: t('dashboard.attention_key_action'),
    to: '/settings',
  },
  {
    id: 'invite',
    level: 'info',
    title: t('dashboard.attention_invite_title'),
    description: t('dashboard.attention_invite_desc'),
    action: t('dashboard.attention_invite_action'),
    to: '/members',
  },
  {
    id: 'integration',
    level: 'info',
    title: t('dashboard.attention_integration_title'),
    description: t('dashboard.attention_integration_desc'),
    action: t('dashboard.attention_integration_action'),
    to: '/integrations',
  },
  {
    id: 'queue',
    level: 'error',
    title: t('dashboard.attention_queue_title'),
    description: t('dashboard.attention_queue_desc'),
    action: t('dashboard.attention_queue_action'),
    to: '/knowledge',
  },
  {
    id: 'trial',
    level: 'warn',
    title: t('dashboard.attention_trial_title'),
    description: t('dashboard.attention_trial_desc'),
    action: t('dashboard.attention_trial_action'),
    to: '/billing',
  },
  {
    id: 'stale_key',
    level: 'info',
    title: t('dashboard.attention_stale_key_title'),
    description: t('dashboard.attention_stale_key_desc'),
    action: t('dashboard.attention_stale_key_action'),
    to: '/settings',
  },
  {
    id: 'latency',
    level: 'warn',
    title: t('dashboard.attention_latency_title'),
    description: t('dashboard.attention_latency_desc'),
    action: t('dashboard.attention_latency_action'),
    to: '/audits',
  },
  {
    id: 'empty_kb',
    level: 'info',
    title: t('dashboard.attention_empty_kb_title'),
    description: t('dashboard.attention_empty_kb_desc'),
    action: t('dashboard.attention_empty_kb_action'),
    to: '/knowledge',
  },
])

function attentionDotClass(level) {
  return {
    error: 'bg-red-500',
    warn: 'bg-amber-500',
    info: 'bg-sky-500',
  }[level] || 'bg-theme-textLight'
}

const activityIcons = {
  doc: 'M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z',
  member: 'M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z',
  key: 'M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z',
}

const recentActivity = computed(() => [
  {
    id: 1,
    icon: 'doc',
    actor: 'System',
    action: t('dashboard.activity_indexed'),
    target: 'onboarding-guide.pdf',
    time: t('dashboard.activity_time_18m'),
  },
  {
    id: 2,
    icon: 'member',
    actor: 'Ahmed',
    action: t('dashboard.activity_invited'),
    target: 'maya@example.com',
    time: t('dashboard.activity_time_3h'),
  },
  {
    id: 3,
    icon: 'key',
    actor: 'Ahmed',
    action: t('dashboard.activity_created_key'),
    target: 'production-ci',
    time: t('dashboard.activity_time_1d'),
  },
])

async function loadDashboard() {
  loading.value = true

  const tasks = [
    statsAPI.get(currentWorkspace.id)
      .then((res) => {
        stats.value = {
          apiCallsMonth: res.data.apiCallsMonth ?? 0,
          documentsStored: res.data.documentsStored ?? 0,
        }
      })
      .catch(() => {}),

    workspaceAPI.get(currentWorkspace.id)
      .then((res) => {
        workspace.value = { ...currentWorkspace, ...res.data, role: res.data.role ?? currentWorkspace.role }
        saveWorkspaceToStorage(workspace.value)
      })
      .catch(() => {}),
  ]

  if (canManage.value) {
    tasks.push(
      billingAPI.status(currentWorkspace.id)
        .then((res) => { plan.value = res.data.plan || 'hobby' })
        .catch(() => {}),
    )
  }

  await Promise.allSettled(tasks)
  loading.value = false
}

onMounted(loadDashboard)
</script>
