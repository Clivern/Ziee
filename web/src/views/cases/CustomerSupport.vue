<script setup>
import CaseLayout from '@/components/CaseLayout.vue'
</script>

<template>
  <CaseLayout
    title="AI Memory for Customer Support | Ziee"
    description="Stop making customers repeat themselves. Give support agents persistent memory across tickets, channels, and sessions."
  >
    <template #default="{ startPath }">
      <section class="relative mx-auto max-w-6xl px-6 pb-20 pt-12 md:pt-16">
        <nav class="reveal mb-8 text-sm text-muted">
          <router-link to="/#use-cases" class="hover:text-theme-text">Use cases</router-link>
          <span class="mx-2">/</span>
          <span class="text-theme-text">Customer Support</span>
        </nav>

        <div class="grid items-center gap-12 lg:grid-cols-2">
          <div>
            <p class="use-case-label opacity-0 animate-fade-up">Customer Support</p>
            <h1 class="hero-gradient-text mt-4 text-4xl font-semibold leading-tight tracking-tight opacity-0 animate-fade-up animate-delay-100 md:text-5xl md:leading-[1.1]">
              Stop Making Customers Repeat Themselves
            </h1>
            <p class="mt-6 text-lg text-muted opacity-0 animate-fade-up animate-delay-200">
              Every ticket, channel, and session connected. Agents pick up exactly where the last conversation left off.
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
                  I have a new phone number: 555-346-0123
                </div>
              </div>
              <div class="mt-4 flex items-start justify-end gap-3">
                <div class="rounded-2xl rounded-tr-sm bg-primary-800 px-4 py-3 text-sm text-white">
                  Saved. I'll use 555-346-0123 on your account from now on.
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
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Why Support Still Breaks Down</h2>
          </div>
          <div class="mt-12 grid gap-6 md:grid-cols-3">
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Same Story, Every Ticket</h3>
              <p class="mt-2 text-sm text-muted">Customers re-explain issues agents already have in the thread, wasting time on both sides.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Context Dies at Handoff</h3>
              <p class="mt-2 text-sm text-muted">Move from chat to email to phone and everything the customer already said disappears.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Patterns Slip Through</h3>
              <p class="mt-2 text-sm text-muted">Recurring bugs and repeat escalations go unnoticed when nothing connects across sessions.</p>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Features</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Smarter Replies, Less Repetition</h2>
        </div>
        <div class="mt-12 grid gap-6 md:grid-cols-2">
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Full Issue History</h3>
            <p class="mt-2 text-sm text-muted">Agents remember past conversations, problems, and resolutions without digging through tickets.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Early Pattern Detection</h3>
            <p class="mt-2 text-sm text-muted">Spots recurring issues across users so teams can act before things escalate.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Seamless Multi-Channel Context</h3>
            <p class="mt-2 text-sm text-muted">Carries memory across chat, email, and phone interactions with preserved context.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Personalized Responses</h3>
            <p class="mt-2 text-sm text-muted">Replies reflect tone, urgency, and past behavior to feel more human.</p>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="grid gap-8 md:grid-cols-2">
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">More Context, Fewer Tokens</h3>
              <p class="mt-3 text-muted">Run smarter agents without stuffing every ticket into the prompt. Retrieve only what matters, in under 10ms.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Memory lives outside the context window</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Up to 90% less redundant context per reply</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Handles long threads without losing the plot</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">Live in Minutes, Not Months</h3>
              <p class="mt-3 text-muted">Drop Ziee into your existing helpdesk, chatbot, or agent stack. No pipeline rewrite required.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Sub-10ms retrieval for live conversations</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Workspace-scoped isolation and access keys</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Python, Go, and HTTP SDKs</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Why Ziee</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">What Changes With Ziee</h2>
        </div>
        <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div class="reveal step-card">
            <h3 class="font-semibold">Full Ticket Recall</h3>
            <p class="mt-2 text-sm text-muted">Past issues, tone, and resolutions, surfaced instantly, not dug up manually.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">One Thread, Every Channel</h3>
            <p class="mt-2 text-sm text-muted">Chat, email, and voice share the same memory. No more starting over.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Bots and Humans, Same Brain</h3>
            <p class="mt-2 text-sm text-muted">Hand off between chatbot and live agent without losing a single detail.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Spot Issues Early</h3>
            <p class="mt-2 text-sm text-muted">Detect repeat complaints across customers before they become outages.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Smart Retention</h3>
            <p class="mt-2 text-sm text-muted">Remember what matters. Forget what shouldn't stick around.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Ship in an Afternoon</h3>
            <p class="mt-2 text-sm text-muted">Add memory with a few API calls. No support flow rebuild.</p>
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
            <router-link to="/cases/e-commerce" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">E-Commerce</p>
              <h3 class="mt-3 text-lg font-semibold">Personal Shopper Agent</h3>
              <p class="mt-2 text-sm text-muted">Size, style, and purchase history, so recommendations actually feel personal.</p>
            </router-link>
          </div>
        </div>
      </section>
    </template>
  </CaseLayout>
</template>
