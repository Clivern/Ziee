<script setup>
import CaseLayout from '@/components/CaseLayout.vue'
</script>

<template>
  <CaseLayout
    title="AI Memory for Healthcare | Ziee"
    description="Patients shouldn't re-explain allergies and preferences at every visit. Persistent context for care agents across channels."
  >
    <template #default="{ startPath }">
      <section class="relative mx-auto max-w-6xl px-6 pb-20 pt-12 md:pt-16">
        <nav class="reveal mb-8 text-sm text-muted">
          <router-link to="/#use-cases" class="hover:text-theme-text">Use cases</router-link>
          <span class="mx-2">/</span>
          <span class="text-theme-text">Healthcare</span>
        </nav>

        <div class="grid items-center gap-12 lg:grid-cols-2">
          <div>
            <p class="use-case-label opacity-0 animate-fade-up">Healthcare</p>
            <h1 class="hero-gradient-text mt-4 text-4xl font-semibold leading-tight tracking-tight opacity-0 animate-fade-up animate-delay-100 md:text-5xl md:leading-[1.1]">
              Care That Remembers Between Visits
            </h1>
            <p class="mt-6 text-lg text-muted opacity-0 animate-fade-up animate-delay-200">
              Patients shouldn't re-explain allergies and preferences at every touchpoint. Ziee keeps care context alive, safely and automatically.
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
                  I'm allergic to latex. Please make sure that's on my record.
                </div>
              </div>
              <div class="mt-4 flex items-start justify-end gap-3">
                <div class="rounded-2xl rounded-tr-sm bg-primary-800 px-4 py-3 text-sm text-white">
                  Noted. I'll flag the latex allergy for scheduling, triage, and every future visit.
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
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Why Care Feels Like Starting Over</h2>
          </div>
          <div class="mt-12 grid gap-6 md:grid-cols-3">
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Critical Details, Repeated</h3>
              <p class="mt-2 text-sm text-muted">Allergies, medications, and preferences get re-collected at every portal login and phone call.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Visits Don't Connect</h3>
              <p class="mt-2 text-sm text-muted">What happened last time rarely shapes this interaction, even when it should.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Generic by Default</h3>
              <p class="mt-2 text-sm text-muted">Without memory, every patient gets the same script instead of care that builds over time.</p>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Features</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Continuity of Care, Built In</h2>
        </div>
        <div class="mt-12 grid gap-6 md:grid-cols-2">
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Patient Preference Memory</h3>
            <p class="mt-2 text-sm text-muted">Recalls communication style, scheduling preferences, and care priorities.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Allergy & Safety Recall</h3>
            <p class="mt-2 text-sm text-muted">Surfaces critical safety information at the right moment in every interaction.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Visit-to-Visit Context</h3>
            <p class="mt-2 text-sm text-muted">Connects prior conversations and outcomes to current care decisions.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Sensitive Data Controls</h3>
            <p class="mt-2 text-sm text-muted">Workspace-scoped memory with time-aware retention for sensitive information.</p>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="grid gap-8 md:grid-cols-2">
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">Built for Sensitive Workloads</h3>
              <p class="mt-3 text-muted">Memory infrastructure designed for environments where accuracy and control matter.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Workspace isolation per care team</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Auditable memory operations</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Time-aware retention policies</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">Integrates With Care Workflows</h3>
              <p class="mt-3 text-muted">Add persistent context to triage bots, scheduling agents, and patient-facing assistants.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Sub-10ms retrieval for live interactions</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Multi-agent memory namespaces</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Drop-in SDK integration</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Why Ziee</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Safer, More Personal Care</h2>
        </div>
        <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div class="reveal step-card">
            <h3 class="font-semibold">Full Patient Context</h3>
            <p class="mt-2 text-sm text-muted">History, preferences, and prior interactions, available when it counts.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Safety First</h3>
            <p class="mt-2 text-sm text-muted">Allergies and contraindications surface at the right moment, every time.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Every Channel, One Record</h3>
            <p class="mt-2 text-sm text-muted">Portal, phone, and chat share the same memory. No gaps.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Care That Adapts</h3>
            <p class="mt-2 text-sm text-muted">Responses improve based on what worked for this patient before.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Governed Retention</h3>
            <p class="mt-2 text-sm text-muted">You control what's remembered, for how long, and who can access it.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Works With What You Have</h3>
            <p class="mt-2 text-sm text-muted">Add memory to triage bots and assistants without replacing your EHR stack.</p>
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
            <router-link to="/cases/education" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">Education</p>
              <h3 class="mt-3 text-lg font-semibold">Adaptive Learning Tutor</h3>
              <p class="mt-2 text-sm text-muted">A tutor that remembers what clicked last time teaches better every session.</p>
            </router-link>
            <router-link to="/cases/sales" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">Sales & CRM</p>
              <h3 class="mt-3 text-lg font-semibold">Deal Memory Assistant</h3>
              <p class="mt-2 text-sm text-muted">Never walk into a follow-up cold. Every objection and stakeholder, recalled instantly.</p>
            </router-link>
          </div>
        </div>
      </section>
    </template>
  </CaseLayout>
</template>
