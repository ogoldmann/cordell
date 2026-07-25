(() => {
  const typeText = (element, text) => {
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    if (prefersReducedMotion) {
      element.textContent = text
      element.dataset.typingDone = 'true'
      return
    }

    element.textContent = ''

    let index = 0

    const tick = () => {
      element.textContent = text.slice(0, index)

      if (index >= text.length) {
        element.dataset.typingDone = 'true'
        return
      }

      index += 1
      window.setTimeout(tick, 24 + Math.random() * 28)
    }

    window.setTimeout(tick, 180)
  }

  const setupDashboard = (root) => {
    const input = root.querySelector('#dashboard-global-search-input')
    const welcome = root.querySelector('[data-dashboard-welcome]')

    if (welcome) {
      const text = welcome.dataset.dashboardWelcomeText || welcome.textContent.trim()
      typeText(welcome, text)
    }

    const updateSearchState = () => {
      const active = input && input.value.trim() !== ''
      root.dataset.searchActive = String(active)
    }

    if (input) {
      window.requestAnimationFrame(() => {
        input.focus({ preventScroll: true })
      })

      input.addEventListener('input', updateSearchState)
      input.addEventListener('search', updateSearchState)

      updateSearchState()
    }
  }

  document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('[data-dashboard-home]').forEach(setupDashboard)
  })
})()
