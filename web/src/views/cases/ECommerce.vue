<script setup>
import CaseLayout from '@/components/CaseLayout.vue'
</script>

<template>
  <CaseLayout
    title="AI Memory for E-Commerce | Ziee"
    description="Shoppers shouldn't re-explain size and style every visit. Personal memory for recommendations that actually land."
  >
    <template #default="{ startPath }">
      <section class="relative mx-auto max-w-6xl px-6 pb-20 pt-12 md:pt-16">
        <nav class="reveal mb-8 text-sm text-muted">
          <router-link to="/#use-cases" class="hover:text-theme-text">Use cases</router-link>
          <span class="mx-2">/</span>
          <span class="text-theme-text">E-Commerce</span>
        </nav>

        <div class="grid items-center gap-12 lg:grid-cols-2">
          <div>
            <p class="use-case-label opacity-0 animate-fade-up">E-Commerce</p>
            <h1 class="hero-gradient-text mt-4 text-4xl font-semibold leading-tight tracking-tight opacity-0 animate-fade-up animate-delay-100 md:text-5xl md:leading-[1.1]">
              Shopping That Actually Knows You
            </h1>
            <p class="mt-6 text-lg text-muted opacity-0 animate-fade-up animate-delay-200">
              Size, style, budget, purchase history, remembered across visits so recommendations feel personal, not generic.
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
                  Show me jackets from this brand.
                </div>
              </div>
              <div class="mt-4 flex items-start justify-end gap-3">
                <div class="rounded-2xl rounded-tr-sm bg-primary-800 px-4 py-3 text-sm text-white">
                  You usually wear size M in this brand. Want me to filter for that?
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
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Why Shopping Bots Feel Robotic</h2>
          </div>
          <div class="mt-12 grid gap-6 md:grid-cols-3">
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Preferences Reset Every Visit</h3>
              <p class="mt-2 text-sm text-muted">Shoppers re-explain size, style, and budget like they've never been here before.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Recommendations Miss the Mark</h3>
              <p class="mt-2 text-sm text-muted">Without memory, suggestions ignore what they've bought, returned, or said they love.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Support Has No Context</h3>
              <p class="mt-2 text-sm text-muted">Order questions and product issues start from zero across chat, email, and app.</p>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Features</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Shopping That Remembers You</h2>
        </div>
        <div class="mt-12 grid gap-6 md:grid-cols-2">
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Size & Fit Memory</h3>
            <p class="mt-2 text-sm text-muted">Recalls sizing across brands and past returns or exchanges.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Style Preferences</h3>
            <p class="mt-2 text-sm text-muted">Remembers colors, categories, and brands a shopper prefers or avoids.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Purchase History Context</h3>
            <p class="mt-2 text-sm text-muted">Surfaces relevant prior orders when answering questions or suggesting items.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Cross-Channel Continuity</h3>
            <p class="mt-2 text-sm text-muted">Carries context from chat to email to in-app support seamlessly.</p>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="grid gap-8 md:grid-cols-2">
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">Higher Conversion, Less Friction</h3>
              <p class="mt-3 text-muted">Personalized agents reduce repetitive questions and surface better recommendations.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Remembers shopper preferences automatically</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Relevant context at recommendation time</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Works across support and sales agents</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">Integrates With Your Store</h3>
              <p class="mt-3 text-muted">Add memory to concierge bots, support agents, and recommendation flows.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Per-shopper agent memory</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Fast retrieval for live chat</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Workspace-scoped data isolation</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Why Ziee</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Commerce That Converts</h2>
        </div>
        <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div class="reveal step-card">
            <h3 class="font-semibold">Knows Their Size</h3>
            <p class="mt-2 text-sm text-muted">Sizing across brands, plus returns and exchanges. No more guesswork.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Style That Sticks</h3>
            <p class="mt-2 text-sm text-muted">Colors, categories, and brands they love, and the ones they don't.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Order History on Tap</h3>
            <p class="mt-2 text-sm text-muted">Past purchases surface when answering questions or making suggestions.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Every Channel Connected</h3>
            <p class="mt-2 text-sm text-muted">Chat, email, and in-app support share the same shopper memory.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Gets Better Over Time</h3>
            <p class="mt-2 text-sm text-muted">Every interaction sharpens recommendations, without manual profile work.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Fits Your Storefront</h3>
            <p class="mt-2 text-sm text-muted">Add memory to concierge bots and support agents in an afternoon.</p>
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
            <router-link to="/cases/sales" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">Sales & CRM</p>
              <h3 class="mt-3 text-lg font-semibold">Deal Memory Assistant</h3>
              <p class="mt-2 text-sm text-muted">Never walk into a follow-up cold. Every objection and stakeholder, recalled instantly.</p>
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
