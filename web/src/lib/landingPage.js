import { computed, onMounted, onUnmounted } from 'vue'
import { loadUserFromStorage } from '@/utils/storage'

export function useStartPath() {
  return computed(() => (loadUserFromStorage() ? '/select-workspace' : '/login'))
}

export function useLandingPage(rootRef, meta) {
  let revealObserver

  function updatePageMeta() {
    if (!meta?.title) return
    document.title = meta.title
    document.querySelector('meta[name="description"]')?.setAttribute('content', meta.description)
    document.querySelector('meta[property="og:title"]')?.setAttribute('content', meta.title)
    document.querySelector('meta[property="og:description"]')?.setAttribute('content', meta.description)
    document.querySelector('meta[name="twitter:title"]')?.setAttribute('content', meta.title)
    document.querySelector('meta[name="twitter:description"]')?.setAttribute('content', meta.description)
  }

  function setupReveal() {
    revealObserver?.disconnect()
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
  }

  onMounted(() => {
    window.scrollTo({ top: 0, left: 0 })
    updatePageMeta()
    requestAnimationFrame(setupReveal)
  })

  onUnmounted(() => {
    revealObserver?.disconnect()
  })
}
