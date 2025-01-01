<script setup>
import CaseLayout from '@/components/CaseLayout.vue'
</script>

<template>
  <CaseLayout
    title="AI Memory for DevTools & Coding Agents | Ziee"
    description="Your coding agent forgets conventions between sessions. Ziee remembers codebase context, decisions, and team standards."
  >
    <template #default="{ startPath }">
      <section class="relative mx-auto max-w-6xl px-6 pb-20 pt-12 md:pt-16">
        <nav class="reveal mb-8 text-sm text-muted">
          <router-link to="/#use-cases" class="hover:text-theme-text">Use cases</router-link>
          <span class="mx-2">/</span>
          <span class="text-theme-text">DevTools</span>
        </nav>

        <div class="grid items-center gap-12 lg:grid-cols-2">
          <div>
            <p class="use-case-label opacity-0 animate-fade-up">DevTools</p>
            <h1 class="hero-gradient-text mt-4 text-4xl font-semibold leading-tight tracking-tight opacity-0 animate-fade-up animate-delay-100 md:text-5xl md:leading-[1.1]">
              Coding Agents That Don't Forget the Rules
            </h1>
            <p class="mt-6 text-lg text-muted opacity-0 animate-fade-up animate-delay-200">
              Conventions, architecture decisions, and codebase context, persisted across sessions so your agent stays on-pattern.
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
                  We use pnpm workspaces and prefer functional components. Stick to that.
                </div>
              </div>
              <div class="mt-4 flex items-start justify-end gap-3">
                <div class="rounded-2xl rounded-tr-sm bg-primary-800 px-4 py-3 text-sm text-white">
                  On it. Functional components in the dashboard package, pnpm workspace layout.
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
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Why Coding Agents Lose the Plot</h2>
          </div>
          <div class="mt-12 grid gap-6 md:grid-cols-3">
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Conventions Get Forgotten</h3>
              <p class="mt-2 text-sm text-muted">Agents re-introduce patterns the team already rejected or ignore project standards.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Context Window Limits</h3>
              <p class="mt-2 text-sm text-muted">Large codebases exceed what fits in a single prompt, causing inconsistent output.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Decisions Do Not Persist</h3>
              <p class="mt-2 text-sm text-muted">Architecture choices and trade-offs discussed in prior sessions vanish.</p>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Features</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Code Agents That Stay Coherent</h2>
        </div>
        <div class="mt-12 grid gap-6 md:grid-cols-2">
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Convention Memory</h3>
            <p class="mt-2 text-sm text-muted">Remembers lint rules, folder structure, naming patterns, and team preferences.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Decision Recall</h3>
            <p class="mt-2 text-sm text-muted">Surfaces prior architecture choices and why they were made.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Relevant Code Context</h3>
            <p class="mt-2 text-sm text-muted">Retrieves the most relevant memories instead of stuffing the full repo into prompts.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Multi-Session Continuity</h3>
            <p class="mt-2 text-sm text-muted">Picks up where the last session left off without re-explaining the codebase.</p>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="grid gap-8 md:grid-cols-2">
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">90% Less Redundant Context</h3>
              <p class="mt-3 text-muted">Compress conversation and codebase knowledge into retrievable memories instead of repeating it every turn.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Semantic memory compression</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Sub-10ms vector retrieval</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Drop-in SDK, no pipeline rewrites</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">Built for Agent Workflows</h3>
              <p class="mt-3 text-muted">Scoped memory per repo, agent, or workspace with full observability.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Agent and workspace namespaces</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Audit logs and metrics</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>API and SDK integration</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Why Ziee</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Agents That Stay on Pattern</h2>
        </div>
        <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div class="reveal step-card">
            <h3 class="font-semibold">Team Conventions, Locked In</h3>
            <p class="mt-2 text-sm text-muted">Lint rules, folder structure, naming, remembered so the agent stops fighting your standards.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Architecture Decisions Stick</h3>
            <p class="mt-2 text-sm text-muted">Prior trade-offs and why they were made, surfaced when the agent needs them.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Multi-Day Coherence</h3>
            <p class="mt-2 text-sm text-muted">Pick up a long task tomorrow without re-explaining the entire codebase.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">90% Less Context Bloat</h3>
            <p class="mt-2 text-sm text-muted">Retrieve what matters instead of stuffing the full repo into every prompt.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Shared Team Memory</h3>
            <p class="mt-2 text-sm text-muted">Every developer's agent aligned to the same project standards.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Fully Observable</h3>
            <p class="mt-2 text-sm text-muted">Audit every memory write and retrieval. Know exactly what the agent knows.</p>
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
            <router-link to="/cases/e-commerce" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">E-Commerce</p>
              <h3 class="mt-3 text-lg font-semibold">Personal Shopper Agent</h3>
              <p class="mt-2 text-sm text-muted">Size, style, and purchase history, so recommendations actually feel personal.</p>
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
