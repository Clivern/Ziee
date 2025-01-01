<script setup>
import { ref } from 'vue'
import LandingNav from '@/components/LandingNav.vue'
import { openCookiePreferences } from '@/lib/cookies'
import { useLandingPage, useStartPath } from '@/lib/landingPage'

const props = defineProps({
  title: { type: String, required: true },
  description: { type: String, required: true },
})

const rootRef = ref(null)
const startPath = useStartPath()

useLandingPage(rootRef, props)
</script>

<template>
  <div ref="rootRef" class="landing-page relative min-h-screen overflow-x-hidden">
    <div class="pointer-events-none absolute inset-0 bg-hero-glow" />
    <div class="pointer-events-none absolute inset-x-0 top-0 h-[520px] hero-grid opacity-60" />

    <LandingNav />

    <main>
      <slot :start-path="startPath" />
    </main>

    <section class="relative overflow-hidden py-24">
      <div class="pointer-events-none absolute inset-0 bg-hero-glow" />
      <div class="relative mx-auto max-w-3xl px-6 text-center">
        <h2 class="reveal text-4xl font-semibold tracking-tight md:text-5xl">
          Ship at <span class="hero-gradient-text">Agent Scale</span>
        </h2>
        <p class="reveal mt-4 text-lg text-muted">
          Ziee merges agent work safely so every delivery cycle keeps moving forward.
        </p>
        <div class="reveal mt-8 flex flex-wrap items-center justify-center gap-3">
          <router-link :to="startPath" class="landing-btn">Get Started</router-link>
          <router-link to="/#pricing" class="landing-btn-outline">See Pricing</router-link>
        </div>
      </div>
    </section>

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
          <div class="grid grid-cols-2 gap-8 sm:grid-cols-3">
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
