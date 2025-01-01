<script setup>
import CaseLayout from '@/components/CaseLayout.vue'
</script>

<template>
  <CaseLayout
    title="AI Memory for Education | Ziee"
    description="A tutor that remembers what clicked last time. Adaptive learning memory for student pace, style, and progress."
  >
    <template #default="{ startPath }">
      <section class="relative mx-auto max-w-6xl px-6 pb-20 pt-12 md:pt-16">
        <nav class="reveal mb-8 text-sm text-muted">
          <router-link to="/#use-cases" class="hover:text-theme-text">Use cases</router-link>
          <span class="mx-2">/</span>
          <span class="text-theme-text">Education</span>
        </nav>

        <div class="grid items-center gap-12 lg:grid-cols-2">
          <div>
            <p class="use-case-label opacity-0 animate-fade-up">Education</p>
            <h1 class="hero-gradient-text mt-4 text-4xl font-semibold leading-tight tracking-tight opacity-0 animate-fade-up animate-delay-100 md:text-5xl md:leading-[1.1]">
              A Tutor That Actually Remembers You
            </h1>
            <p class="mt-6 text-lg text-muted opacity-0 animate-fade-up animate-delay-200">
              Generic tutoring resets every session. Ziee remembers each student's pace, style, and struggles, so every lesson builds on the last.
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
                  I'm stuck on this calculus problem.
                </div>
              </div>
              <div class="mt-4 flex items-start justify-end gap-3">
                <div class="rounded-2xl rounded-tr-sm bg-primary-800 px-4 py-3 text-sm text-white">
                  You learn best with visual examples. Let's try a diagram for this concept.
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
            <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Why Most Tutors Feel the Same</h2>
          </div>
          <div class="mt-12 grid gap-6 md:grid-cols-3">
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Blank Slate Every Time</h3>
              <p class="mt-2 text-sm text-muted">The agent forgets what the student struggled with or nailed last session.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">One Pace for Everyone</h3>
              <p class="mt-2 text-sm text-muted">Without memory, tutoring can't adapt to how this student actually learns.</p>
            </div>
            <div class="reveal glass-card p-6">
              <h3 class="text-lg font-semibold">Progress Evaporates</h3>
              <p class="mt-2 text-sm text-muted">Real skill-building needs continuity that a single session can't provide.</p>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Features</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Tutoring That Remembers</h2>
        </div>
        <div class="mt-12 grid gap-6 md:grid-cols-2">
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Learning Style Adaptation</h3>
            <p class="mt-2 text-sm text-muted">Remembers whether a student prefers examples, diagrams, or step-by-step explanations.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Mastery Tracking</h3>
            <p class="mt-2 text-sm text-muted">Recalls concepts mastered and areas that need reinforcement.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Pace-Aware Guidance</h3>
            <p class="mt-2 text-sm text-muted">Adjusts difficulty based on prior performance and engagement signals.</p>
          </div>
          <div class="reveal use-case-card p-6">
            <h3 class="text-lg font-semibold">Long-Horizon Progress</h3>
            <p class="mt-2 text-sm text-muted">Builds on prior lessons instead of treating each session independently.</p>
          </div>
        </div>
      </section>

      <section class="section-alt py-24">
        <div class="mx-auto max-w-6xl px-6">
          <div class="grid gap-8 md:grid-cols-2">
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">Personalization at Scale</h3>
              <p class="mt-3 text-muted">Give every student a tutor that remembers them without manual profile management.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Automatic memory extraction from sessions</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Relevant context at lesson time</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Lower token use with compressed memory</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
            <div class="reveal glass-card bg-card-glow p-8">
              <h3 class="text-xl font-semibold">Fits Your Learning Platform</h3>
              <p class="mt-3 text-muted">Integrate memory into LMS bots, tutoring apps, and study assistants.</p>
              <ul class="mt-6 space-y-3 text-sm text-theme-text">
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Agent-scoped student memory</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>Fast retrieval for live tutoring</li>
                <li class="flex items-start gap-2"><span class="mt-1 text-primary-700">→</span>SDK for Python, Go, and HTTP</li>
              </ul>
              <router-link :to="startPath" class="landing-btn mt-8 inline-flex">Get Started</router-link>
            </div>
          </div>
        </div>
      </section>

      <section class="mx-auto max-w-6xl px-6 py-24">
        <div class="reveal mx-auto max-w-2xl text-center">
          <p class="section-label">Why Ziee</p>
          <h2 class="mt-3 text-3xl font-semibold tracking-tight md:text-4xl">Learning That Compounds</h2>
        </div>
        <div class="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div class="reveal step-card">
            <h3 class="font-semibold">Knows Their Style</h3>
            <p class="mt-2 text-sm text-muted">Diagrams, examples, or step-by-step. The tutor remembers what works for this student.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Sessions That Stack</h3>
            <p class="mt-2 text-sm text-muted">Each lesson picks up where the last one left off. No re-explaining basics.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Mastery Tracking</h3>
            <p class="mt-2 text-sm text-muted">Knows what's solid and what needs another pass before moving on.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Pace-Aware Teaching</h3>
            <p class="mt-2 text-sm text-muted">Adjusts difficulty based on how this student has performed over time.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Weeks of Coherence</h3>
            <p class="mt-2 text-sm text-muted">Context that spans a semester, not just a single chat window.</p>
          </div>
          <div class="reveal step-card">
            <h3 class="font-semibold">Plugs Into Your LMS</h3>
            <p class="mt-2 text-sm text-muted">Add memory to tutoring apps and study bots without rebuilding your platform.</p>
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
            <router-link to="/cases/customer-support" class="reveal use-case-card block p-6 transition hover:border-primary-300">
              <p class="use-case-label">Customer Support</p>
              <h3 class="mt-3 text-lg font-semibold">Context-Aware Support Agent</h3>
              <p class="mt-2 text-sm text-muted">Stop making customers repeat themselves across tickets and channels.</p>
            </router-link>
          </div>
        </div>
      </section>
    </template>
  </CaseLayout>
</template>
