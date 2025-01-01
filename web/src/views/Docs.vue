<template>
  <div class="landing-page min-h-screen bg-theme-bg">
    <header class="landing-header sticky top-0 z-50">
      <nav class="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <router-link to="/" class="flex items-center gap-2.5">
          <img src="/logo.png" alt="Ziee" class="h-8 w-auto" />
          <span class="text-lg font-semibold tracking-tight">Ziee</span>
        </router-link>
        <router-link to="/" class="landing-link">Back to home</router-link>
      </nav>
    </header>

    <main class="mx-auto flex max-w-2xl flex-col items-center px-6 py-24 text-center">
      <p class="section-label">Documentation</p>
      <h1 class="mt-3 text-3xl font-semibold tracking-tight text-theme-text">Docs coming soon</h1>
      <p class="mt-4 text-muted">
        Guides, API references, and SDK documentation for building with Ziee are on the way.
      </p>
    </main>
  </div>
</template>
