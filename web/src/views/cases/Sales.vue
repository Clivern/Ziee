<script setup>
import CaseLayout from '@/components/CaseLayout.vue'
</script>

<template>
  <CaseLayout
    title="AI Memory for Sales & CRM | Ziee"
    description="Deals run on relationships. Give sales agents memory for every stakeholder, objection, and commitment across the full cycle."
  >
    <template #default="{ startPath }">
      <section class="relative mx-auto max-w-6xl px-6 pb-20 pt-12 md:pt-16">
        <nav class="reveal mb-8 text-sm text-muted">
          <router-link to="/#use-cases" class="hover:text-theme-text">Use cases</router-link>
          <span class="mx-2">/</span>
          <span class="text-theme-text">Sales & CRM</span>
        </nav>

        <div class="grid items-center gap-12 lg:grid-cols-2">
          <div>
            <p class="use-case-label opacity-0 animate-fade-up">Sales & CRM</p>
            <h1 class="hero-gradient-text mt-4 text-4xl font-semibold leading-tight tracking-tight opacity-0 animate-fade-up animate-delay-100 md:text-5xl md:leading-[1.1]">
              Never Walk Into a Follow-Up Cold
            </h1>
            <p class="mt-6 text-lg text-muted opacity-0 animate-fade-up animate-delay-200">
              Deals run on relationships. Ziee remembers every stakeholder, objection, and commitment across weeks or months.
            </p>
            <div class="mt-8 opacity-0 animate-fade-up animate-delay-300">
              <router-link :to="startPath" class="landing-btn">Get Started</router-link>
            </div>
          </div>

          <div class="reveal">
            <div class="glass-card p-6">
              <div class="flex items-start gap-3">
                <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-primary-100 text-sm font-medium text-primary-800">U</div>
                <div class="rounded-2xl rounded-tl-sm bg-theme-hover px-4 py-3 text-sm text-theme-text">
                  They mentioned budget approval happens in Q3. Can we follow up then?
                </div>
              </div>
              <div class="mt-4 flex items-start justify-end gap-3">
                <div class="rounded-2xl rounded-tr-sm bg-primary-800 px-4 py-3 text-sm text-white">
                  Got it. I'll flag this deal for a Q3 follow-up when budget opens.
                </div>
                <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-primary-800 text-sm font-medium text-white">A</div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="reveal mx-auto max-w-2xl text-center">
            <p class="section-label">Problems</p>
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Why Deals Die in the CRM</h2>
          </div>
          <div class="mt-12 grid gap-6 md:grid-cols-3">
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Every Call Starts Over</h3>
              <p class="mt-2 text-sm text-muted">Reps re-discover the same objections and politics instead of advancing the deal.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Handoffs Kill Momentum</h3>
              <p class="mt-2 text-sm text-muted">When ownership changes, the relationship context walks out the door with the last rep.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">The CRM Is Always Stale</h3>
              <p class="mt-2 text-sm text-muted">The nuance from calls and emails never makes it into structured fields, if it gets logged at all.</p>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Features</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Deal Memory That Moves With You</h2>
        </div>
        <div class="mt-12 grid gap-6 md:grid-cols-2">
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Stakeholder Mapping</h3>
            <p class="mt-2 text-sm text-muted">Remembers who cares about what, who blocks deals, and who champions them.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Objection Tracking</h3>
            <p class="mt-2 text-sm text-muted">Surfaces past concerns and how they were addressed across the full cycle.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Timeline Awareness</h3>
            <p class="mt-2 text-sm text-muted">Recalls budget windows, renewal dates, and follow-up commitments automatically.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Personalized Outreach</h3>
            <p class="mt-2 text-sm text-muted">Drafts replies that reference prior conversations instead of generic templates.</p>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="grid gap-8 md:grid-cols-2">
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">Less Admin, More Selling</h3>
              <p class="mt-3 text-muted">Capture conversation context automatically instead of manual CRM updates after every call.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Extracts deal signals from conversations</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Surfaces the right context at call time</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Reduces repetitive note-taking</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">Works With Your CRM Stack</h3>
              <p class="mt-3 text-muted">Drop Ziee into your existing sales workflows without replacing your CRM.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>SDK and API integration</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Agent-scoped memory namespaces</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Fast retrieval for live conversations</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Why Ziee</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Close Faster With Full Deal Context</h2>
        </div>
        <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div class="reveal step-card">
            <h3 class="font-semibold">Instant Deal Recall</h3>
            <p class="mt-2 text-sm text-muted">Every objection, milestone, and preference, one search away before the call starts.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Seamless Rep Handoffs</h3>
            <p class="mt-2 text-sm text-muted">New owner inherits the full relationship history, not a half-empty CRM record.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Built for Long Cycles</h3>
            <p class="mt-2 text-sm text-muted">Memory that spans weeks and months, not just the current session.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Stalled Deal Alerts</h3>
            <p class="mt-2 text-sm text-muted">Surface recurring blockers and quiet deals before they go dark.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Timing That Lands</h3>
            <p class="mt-2 text-sm text-muted">Follow up when budget opens, not when someone finally remembers to set a reminder.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Fits Your Stack</h3>
            <p class="mt-2 text-sm text-muted">Layer onto your CRM and sales tools. No rip-and-replace.</p>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="reveal mx-auto max-w-2xl text-center">
            <p class="section-label">More use cases</p>
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">See it in other industries</h2>
          </div>
          <div class="mt-12 grid gap-6 md:grid-cols-3">
            <router-link to="/cases/customer-support" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">Customer Support</p>
              <h3 class="mt-3 text-lg font-semibold">Context-Aware Support Agent</h3>
              <p class="mt-2 text-sm text-muted">Stop making customers repeat themselves across tickets and channels.</p>
            </router-link>
            <router-link to="/cases/healthcare" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">Healthcare</p>
              <h3 class="mt-3 text-lg font-semibold">Patient Care Companion</h3>
              <p class="mt-2 text-sm text-muted">Allergies, preferences, and visit history, remembered between every touchpoint.</p>
            </router-link>
            <router-link to="/cases/devtools" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">DevTools</p>
              <h3 class="mt-3 text-lg font-semibold">Coding Agent Memory</h3>
              <p class="mt-2 text-sm text-muted">Your agent remembers the conventions, so it stops re-introducing bad patterns.</p>
            </router-link>
          </div>
        </div>
      </section>
    </template>
  </CaseLayout>
</template>
