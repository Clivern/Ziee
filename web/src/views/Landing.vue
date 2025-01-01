<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { highlightCode } from '@/lib/code'
import { openCookiePreferences } from '@/lib/cookies'
import { loadUserFromStorage } from '@/utils/storage'
import LandingNav from '@/components/LandingNav.vue'

const codeTab = ref('python')
const featureTab = ref('efficiency')
const copyLabel = ref('Copy')
const rootRef = ref(null)

const highlightedCode = computed(() => highlightCode(codeSamples[codeTab.value], codeTab.value))

const startPath = computed(() => {
  if (loadUserFromStorage()) return '/select-workspace'
  return '/login'
})

const plans = [
  {
    id: 'hobby',
    name: 'Hobby',
    price: '$0',
    unit: '/mo',
    description: 'Free workspace plan to get started with core capabilities.',
    features: [
      'Core workspace access',
      'Basic usage limits',
      'Community support',
      'Upgrade anytime',
    ],
  },
  {
    id: 'starter',
    name: 'Starter',
    price: '$19',
    unit: '/mo',
    description: 'For individuals and small projects that need more room.',
    popular: true,
    features: [
      'Everything in Hobby',
      'Higher usage limits',
      'Managed infrastructure',
      'Email support',
    ],
  },
  {
    id: 'growth',
    name: 'Growth',
    price: '$69',
    unit: '/mo',
    description: 'More capacity for growing teams and heavier workloads.',
    features: [
      'Everything in Starter',
      'Expanded usage limits',
      'Priority capacity',
      'Faster support',
    ],
  },
  {
    id: 'pro',
    name: 'Pro',
    price: '€219',
    unit: '/mo',
    description: 'Advanced plan for teams that need the highest monthly capacity.',
    features: [
      'Everything in Growth',
      'Highest usage limits',
      'Premium capacity',
      'Priority support',
    ],
  },
]

const codeSamples = {
  python: `# pip install ziee
from ziee import MergeClient

client = MergeClient(
    api_key="your-access-key",
    workspace="your-workspace",
)

# Enqueue an agent branch for autonomous merge
client.enqueue(
    branch="agent/feature-checkout",
    target="main",
    policy="auto-merge-on-green",
)

status = client.status(
    branch="agent/feature-checkout",
)
print(status)`,
  curl: `# Enqueue a merge for an agent branch
curl -X POST "https://your-instance/api/v1/merge/enqueue" \\
  -H "Authorization: Bearer YOUR_ACCESS_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "branch": "agent/feature-checkout",
    "target": "main",
    "policy": "auto-merge-on-green"
  }'`,
  go: `// go get github.com/actx0/ziee/sdk/go
client, _ := ziee.New(ziee.Config{
    APIKey:    "your-access-key",
    Workspace: "your-workspace",
})

_, _ = client.Enqueue(ctx, ziee.EnqueueRequest{
    Branch: "agent/feature-checkout",
    Target: "main",
    Policy: "auto-merge-on-green",
})

status, _ := client.Status(ctx, ziee.StatusRequest{
    Branch: "agent/feature-checkout",
})`,
}

let revealObserver

async function copyCode() {
  try {
    await navigator.clipboard.writeText(codeSamples[codeTab.value])
    copyLabel.value = 'Copied!'
    setTimeout(() => {
      copyLabel.value = 'Copy'
    }, 1600)
  } catch {
    /* clipboard unavailable */
  }
}

onMounted(() => {
  const items = rootRef.value?.querySelectorAll('.reveal') ?? []
  revealObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible')
          revealObserver.unobserve(entry.target)
        }
      })
    },
    { threshold: 0.12, rootMargin: '0px 0px -40px 0px' }
  )
  items.forEach((item) => revealObserver.observe(item))
})

onUnmounted(() => {
  revealObserver?.disconnect()
})
</script>

<template>
  <div ref="rootRef" class="landing-page relative min-h-screen overflow-x-hidden">
    <div class="pointer-events-none absolute inset-0 bg-hero-glow" />
    <div class="pointer-events-none absolute inset-x-0 top-0 h-[520px] hero-grid opacity-60" />

    <LandingNav home-page />

    <main>
      <section class="relative mx-auto max-w-6xl px-6 pb-20 pt-16 md:pt-24">
        <div class="mx-auto max-w-3xl text-center">
          <div class="opacity-0 animate-fade-up">
            <h1 class="hero-gradient-text text-4xl font-semibold leading-tight tracking-tight md:text-6xl md:leading-[1.1]">
              The Autonomous Merge Layer for Agent-Scale Delivery
            </h1>
          </div>
          <p class="mx-auto mt-6 max-w-2xl text-lg text-muted opacity-0 animate-fade-up animate-delay-200">
            Merge agent branches safely, resolve conflicts automatically, and keep shipping continuous as parallel agents multiply.
          </p>
          <div class="mt-8 flex flex-wrap items-center justify-center gap-3 opacity-0 animate-fade-up animate-delay-300">
            <router-link :to="startPath" class="landing-btn">Get Started</router-link>
          </div>
        </div>

        <div class="reveal mx-auto mt-16 max-w-3xl">
          <div class="code-window">
            <div class="flex items-center justify-between border-b border-white/10 bg-primary-800 px-4 py-3">
              <div class="flex items-center gap-2">
                <span class="h-3 w-3 rounded-full bg-red-500/80" />
                <span class="h-3 w-3 rounded-full bg-amber-500/80" />
                <span class="h-3 w-3 rounded-full bg-emerald-500/80" />
              </div>
              <div class="flex gap-1">
                <button
                  v-for="tab in ['python', 'curl', 'go']"
                  :key="tab"
                  type="button"
                  class="code-tab"
                  :class="{ active: codeTab === tab }"
                  @click="codeTab = tab"
                >
                  {{ tab === 'curl' ? 'cURL' : tab.charAt(0).toUpperCase() + tab.slice(1) }}
                </button>
              </div>
              <button type="button" class="landing-btn-outline px-3 py-1.5 text-xs" @click="copyCode">
                {{ copyLabel }}
              </button>
            </div>
            <div class="overflow-x-auto p-5">
              <pre class="code-line" v-html="highlightedCode" />
            </div>
          </div>
        </div>
      </section>

      <section class="section-alt py-12">
        <div class="mx-auto grid max-w-6xl grid-cols-2 gap-8 px-6 md:grid-cols-4">
          <div v-for="stat in [
            { value: '24/7', label: 'Autonomous merge queue' },
            { value: 'Policy', label: 'Guardrails on every merge' },
            { value: 'Managed', label: 'Cloud infrastructure' },
            { value: 'Monthly', label: 'Flexible billing' },
          ]" :key="stat.label" class="reveal text-center">
            <p class="text-3xl font-semibold md:text-4xl">{{ stat.value }}</p>
            <p class="mt-1 text-sm text-muted">{{ stat.label }}</p>
          </div>
        </div>
      </section>

      <section id="developers" class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Built for developers</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Proof, not promises</h2>
          <p class="mt-4 text-muted">
            Ziee merges agent output without manual queue babysitting. Fewer blocked PRs, faster delivery, and merge throughput that scales with agent output.
          </p>
        </div>

        <div class="reveal mt-12">
          <div class="mx-auto flex w-fit gap-2 rounded-xl border border-theme-border bg-theme-hover p-1">
            <button
              v-for="tab in ['efficiency', 'visibility', 'control']"
              :key="tab"
              type="button"
              class="feature-tab capitalize"
              :class="{ active: featureTab === tab }"
              @click="featureTab = tab"
            >
              {{ tab }}
            </button>
          </div>

          <div class="mt-8 grid gap-8 md:grid-cols-2">
            <template v-if="featureTab === 'efficiency'">
              <div class="glass-card bg-card-glow p-8">
                <h3 class="text-xl font-semibold">Memory Compression Engine</h3>
                <p class="mt-3 text-muted">
                  Automatically condenses chat history into compact memories that cut tokens and latency while keeping the right context for each agent turn.
                </p>
                <ul class="mt-6 space-y-3 text-sm text-theme-text">
                  <li v-for="item in [
                    'Vector search powered by Qdrant for sub-10ms retrieval',
                    'Semantic deduplication across sessions',
                    'Drop-in SDK, no pipeline rewrites',
                  ]" :key="item" class="flex items-start gap-2">
                    <span class="mt-1 text-primary-700">→</span>
                    {{ item }}
                  </li>
                </ul>
              </div>
              <div class="glass-card flex items-center justify-center p-8">
                <div class="w-full space-y-3 font-mono text-xs">
                  <div class="rounded-lg bg-theme-hover p-4 text-primary-500">Before: 12,400 tokens / request</div>
                  <div class="rounded-lg border border-primary-300 bg-primary-50 p-4 text-primary-800">After: 1,240 tokens / request</div>
                  <div class="h-2 overflow-hidden rounded-full bg-primary-200">
                    <div class="h-full w-[10%] rounded-full bg-primary-800" />
                  </div>
                  <p class="text-center text-primary-500">90% context compression</p>
                </div>
              </div>
            </template>

            <div v-else-if="featureTab === 'visibility'" class="glass-card bg-card-glow p-8 md:col-span-2">
              <div class="grid gap-8 md:grid-cols-2">
                <div>
                  <h3 class="text-xl font-semibold">Full Observability</h3>
                  <p class="mt-3 text-muted">
                    Every memory write, retrieval, and agent session is auditable. Prometheus metrics, usage events, and audit logs built in from day one.
                  </p>
                </div>
                <div class="space-y-2">
                  <div class="log-line">memory.search · 4.2ms · agent:support-bot</div>
                  <div class="log-line">memory.add · 8.1ms · 3 facts extracted</div>
                  <div class="log-line-muted">session.start · agent:onboarding</div>
                </div>
              </div>
            </div>

            <div v-else class="glass-card bg-card-glow p-8 md:col-span-2">
              <div class="grid gap-8 md:grid-cols-2">
                <div>
                  <h3 class="text-xl font-semibold">Workspace-Scoped Control</h3>
                  <p class="mt-3 text-muted">
                    Multi-tenant by design. Workspace access keys, role-based permissions, prompt versioning, and knowledge bases, all isolated per team.
                  </p>
                </div>
                <div class="grid grid-cols-2 gap-3 text-sm">
                  <div v-for="item in [
                    { title: 'Access Keys', text: 'Scoped API credentials' },
                    { title: 'Prompt Versions', text: 'Production labels' },
                    { title: 'Knowledge Base', text: 'Document ingestion' },
                    { title: 'Agent Sessions', text: 'Persistent threads' },
                  ]" :key="item.title" class="inner-card">
                    <p class="font-medium">{{ item.title }}</p>
                    <p class="mt-1 text-muted">{{ item.text }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="reveal mx-auto max-w-2xl text-center">
            <p class="section-label">How it works</p>
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Enqueue. Merge. Ship.</h2>
          </div>
          <div class="mt-16 grid gap-6 md:grid-cols-3">
            <div v-for="(step, index) in [
              { title: 'Enqueue', text: 'Point agent branches at Ziee with a target branch and merge policy. No custom merge scripts required.' },
              { title: 'Merge', text: 'Ziee resolves conflicts, runs checks, and lands changes automatically when policy conditions are met.' },
              { title: 'Ship', text: 'Keep delivery continuous as parallel agents multiply. Every merge is auditable and policy-bound.' },
            ]" :key="step.title" class="reveal step-card">
              <div class="step-badge">{{ index + 1 }}</div>
              <h3 class="text-lg font-semibold">{{ step.title }}</h3>
              <p class="mt-2 text-sm text-muted">{{ step.text }}</p>
            </div>
          </div>
        </div>
      </section>

      <section id="use-cases" class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Use cases</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Built for every kind of agent</h2>
        </div>
        <div class="mt-12 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <router-link to="/cases/customer-support" class="reveal use-case-card block p-6 transition hover:border-primary-300">
            <p class="use-case-label">Customer Support</p>
            <h3 class="mt-3 text-lg font-semibold">Context-Aware Support Agent</h3>
            <p class="mt-2 text-sm text-muted">Stop making customers repeat themselves across tickets and channels.</p>
          </router-link>
          <router-link to="/cases/sales" class="reveal use-case-card block p-6 transition hover:border-primary-300">
            <p class="use-case-label">Sales & CRM</p>
            <h3 class="mt-3 text-lg font-semibold">Deal Memory Assistant</h3>
            <p class="mt-2 text-sm text-muted">Never walk into a follow-up cold. Every objection and stakeholder, recalled instantly.</p>
          </router-link>
          <router-link to="/cases/healthcare" class="reveal use-case-card block p-6 transition hover:border-primary-300">
            <p class="use-case-label">Healthcare</p>
            <h3 class="mt-3 text-lg font-semibold">Patient Care Companion</h3>
            <p class="mt-2 text-sm text-muted">Allergies, preferences, and visit history, remembered between every touchpoint.</p>
          </router-link>
          <router-link to="/cases/education" class="reveal use-case-card block p-6 transition hover:border-primary-300">
            <p class="use-case-label">Education</p>
            <h3 class="mt-3 text-lg font-semibold">Adaptive Learning Tutor</h3>
            <p class="mt-2 text-sm text-muted">A tutor that remembers what clicked last time teaches better every session.</p>
          </router-link>
          <router-link to="/cases/devtools" class="reveal use-case-card block p-6 transition hover:border-primary-300">
            <p class="use-case-label">DevTools</p>
            <h3 class="mt-3 text-lg font-semibold">Coding Agent Memory</h3>
            <p class="mt-2 text-sm text-muted">Your agent remembers the conventions, so it stops re-introducing bad patterns.</p>
          </router-link>
          <router-link to="/cases/e-commerce" class="reveal use-case-card block p-6 transition hover:border-primary-300">
            <p class="use-case-label">E-Commerce</p>
            <h3 class="mt-3 text-lg font-semibold">Personal Shopper Agent</h3>
            <p class="mt-2 text-sm text-muted">Size, style, and purchase history, so recommendations actually feel personal.</p>
          </router-link>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="reveal mx-auto max-w-2xl text-center">
            <p class="section-label">Enterprise</p>
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Built for scale. Designed for control.</h2>
            <p class="mt-4 text-muted">
              Agent-scale delivery is infrastructure. Ziee gives teams governance, reliability, and observability so engineers ship instead of babysitting merge queues.
            </p>
          </div>
          <div class="mt-12 grid gap-6 md:grid-cols-3">
            <div v-for="item in [
              { title: 'Managed Cloud', text: 'Fully managed memory infrastructure. No servers to provision, patch, or scale yourself.' },
              { title: 'Auditable', text: 'Audit events, usage tracking, and Prometheus metrics on every operation. Know what happened, when, and for whom.' },
              { title: 'Multi-Tenant', text: 'Workspaces, roles, access keys, and billing per team. Isolated memory namespaces with shared infrastructure.' },
            ]" :key="item.title" class="reveal glass-card p-6">
              <h3 class="font-semibold">{{ item.title }}</h3>
              <p class="mt-2 text-sm text-muted">{{ item.text }}</p>
            </div>
          </div>
        </div>
      </section>

      <section id="pricing" class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Pricing</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Start free. Scale when you're ready.</h2>
          <p class="mt-4 text-muted">All plans run on Ziee managed cloud. Billed monthly.</p>
        </div>

        <div class="mt-12 grid gap-6 sm:grid-cols-2 xl:grid-cols-4">
          <div
            v-for="plan in plans"
            :key="plan.id"
            class="reveal glass-card relative flex flex-col p-8"
            :class="{ 'pricing-popular': plan.popular }"
          >
            <span v-if="plan.popular" class="pricing-badge">Popular</span>
            <p class="text-sm font-medium text-muted">{{ plan.name }}</p>
            <p class="mt-2 text-3xl font-semibold">
              {{ plan.price }}<span class="text-lg font-normal text-muted">{{ plan.unit }}</span>
            </p>
            <p class="mt-2 text-sm text-muted">{{ plan.description }}</p>
            <ul class="mt-6 flex-1 space-y-2 text-sm text-theme-text">
              <li v-for="feature in plan.features" :key="feature">{{ feature }}</li>
            </ul>
            <router-link
              :to="startPath"
              class="mt-8 w-full"
              :class="plan.popular ? 'landing-btn' : 'landing-btn-outline'"
            >
              {{ plan.id === 'hobby' ? 'Get started free' : 'Start trial' }}
            </router-link>
          </div>
        </div>

        <div class="reveal glass-card mx-auto mt-8 flex max-w-4xl flex-col items-start justify-between gap-6 p-8 md:flex-row md:items-center">
          <div>
            <p class="text-sm font-medium text-muted">Enterprise</p>
            <h3 class="mt-1 text-2xl font-semibold">Custom plans for larger teams</h3>
            <p class="mt-2 text-sm text-muted">
              SSO & SAML, dedicated infrastructure, SLA & compliance, and custom integrations.
            </p>
          </div>
          <a href="mailto:hello@clivern.com" class="landing-btn-outline shrink-0">Talk to us</a>
        </div>
      </section>

      <section id="get-started" class="relative overflow-hidden py-24">
        <div class="pointer-events-none absolute inset-0 bg-hero-glow" />
        <div class="relative mx-auto max-w-3xl px-6 text-center">
          <h2 class="reveal text-4xl font-semibold tracking-tight md:text-5xl">
            Ship at <span class="hero-gradient-text">Agent Scale</span>
          </h2>
          <p class="reveal mt-4 text-lg text-muted">
            The autonomous merge layer for agent-scale delivery. Merge safely. Ship continuously.
          </p>
          </p>
          <div class="reveal mt-8 flex flex-wrap items-center justify-center gap-3">
            <router-link :to="startPath" class="landing-btn">Get Started</router-link>
            <a href="#pricing" class="landing-btn-outline">See Pricing</a>
          </div>
        </div>
      </section>
    </main>

    <footer class="border-t border-theme-border bg-theme-hover">
      <div class="mx-auto max-w-6xl px-6 py-12">
        <div class="flex flex-col gap-8 md:flex-row md:justify-between">
          <div>
            <router-link to="/" class="flex items-center gap-2.5">
              <img src="/logo.png" alt="Ziee" class="h-7 w-auto" />
              <span class="font-semibold">Ziee</span>
            </router-link>
            <p class="mt-3 max-w-xs text-sm text-muted">The autonomous merge layer for agent-scale delivery.</p>
          </div>
          <div class="grid grid-cols-2 gap-8 sm:grid-cols-4">
            <div>
              <p class="text-sm font-medium">Product</p>
              <ul class="mt-3 space-y-2 text-sm text-muted">
                <li><a href="#developers" class="hover:text-theme-text">Developers</a></li>
                <li><a href="#pricing" class="hover:text-theme-text">Pricing</a></li>
                <li><a href="#use-cases" class="hover:text-theme-text">Use Cases</a></li>
              </ul>
            </div>
            <div>
              <p class="text-sm font-medium">Use Cases</p>
              <ul class="mt-3 space-y-2 text-sm text-muted">
                <li><router-link to="/cases/customer-support" class="hover:text-theme-text">Customer Support</router-link></li>
                <li><router-link to="/cases/sales" class="hover:text-theme-text">Sales & CRM</router-link></li>
                <li><router-link to="/cases/healthcare" class="hover:text-theme-text">Healthcare</router-link></li>
                <li><router-link to="/cases/education" class="hover:text-theme-text">Education</router-link></li>
                <li><router-link to="/cases/devtools" class="hover:text-theme-text">DevTools</router-link></li>
                <li><router-link to="/cases/e-commerce" class="hover:text-theme-text">E-Commerce</router-link></li>
              </ul>
            </div>
            <div>
              <p class="text-sm font-medium">Resources</p>
              <ul class="mt-3 space-y-2 text-sm text-muted">
                <li><router-link to="/docs" class="hover:text-theme-text">Docs</router-link></li>
                <li><router-link to="/status" class="hover:text-theme-text">Status</router-link></li>
                <li><a href="https://github.com/actx0/ziee/issues" class="hover:text-theme-text" target="_blank" rel="noopener">Support</a></li>
              </ul>
            </div>
            <div>
              <p class="text-sm font-medium">Company</p>
              <ul class="mt-3 space-y-2 text-sm text-muted">
                <li><a href="mailto:hello@clivern.com" class="hover:text-theme-text">Contact</a></li>
                <li>
                  <button type="button" class="hover:text-theme-text" @click="openCookiePreferences">
                    {{ $t('cookies.settings_link') }}
                  </button>
                </li>
              </ul>
            </div>
          </div>
        </div>
        <div class="mt-10 border-t border-theme-border pt-6 text-center text-sm text-primary-500">
          Copyright © 2026 Ziee. All Rights Reserved.
        </div>
      </div>
    </footer>
  </div>
</template>
