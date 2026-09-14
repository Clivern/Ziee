<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="66 30 368 404"
    fill="currentColor"
    :role="alt ? 'img' : 'presentation'"
    :aria-label="alt || undefined"
    :aria-hidden="alt ? undefined : true"
  >
    <circle cx="249.8" cy="82" r="46" />
    <circle cx="249.8" cy="190" r="46" />
    <path d="m334.64 261.9c-12.44-12.44-11.39-12-29-12h-111.64c-17.62 0-16.57-.43-29 12l-59.52 59.53a161.1 161.1 0 0 0 60.34 65.93c15.25-15.24 61-60.95 68.57-68.57 6.25-6.24 9-8.66 12.39-9.7a9.87 9.87 0 0 1 3-.56 9.71 9.71 0 0 1 3 .56c3.39 1 6.15 3.46 12.39 9.7l68.57 68.57a161.06 161.06 0 0 0 60.35-65.93z" />
  </svg>
</template>

<script setup>
defineProps({
  alt: {
    type: String,
    default: ''
  }
})
</script>
