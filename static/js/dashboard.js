(() => {
  const welcomeSeenKey = 'cordell-dashboard-welcome-seen'

  const updateNavbarHeight = () => {
    const navbar = document.querySelector('.cordell-navbar')

    if (!navbar) {
      return
    }

    document.documentElement.style.setProperty(
      '--cordell-navbar-height',
      `${navbar.offsetHeight}px`,
    )
  }

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

  const setupWelcome = (root) => {
    const welcome = root.querySelector('[data-dashboard-welcome]')

    if (!welcome) {
      return
    }

    const alreadySeen = window.sessionStorage.getItem(welcomeSeenKey) === 'true'

    if (alreadySeen) {
      welcome.dataset.welcomeState = 'hidden'
      welcome.setAttribute('aria-hidden', 'true')
      root.dataset.welcomeVisible = 'false'
      return
    }

    welcome.dataset.welcomeState = 'visible'
    welcome.removeAttribute('aria-hidden')
    root.dataset.welcomeVisible = 'true'

    const text = welcome.dataset.dashboardWelcomeText || welcome.textContent.trim()
    typeText(welcome, text)

    window.sessionStorage.setItem(welcomeSeenKey, 'true')
  }

  const setupSearch = (root) => {
    const input = root.querySelector('#dashboard-global-search-input')

    const updateSearchState = () => {
      const active = input && input.value.trim() !== ''
      root.dataset.searchActive = String(active)
    }

    if (!input) {
      return
    }

    window.requestAnimationFrame(() => {
      input.focus({ preventScroll: true })
    })

    input.addEventListener('input', updateSearchState)
    input.addEventListener('search', updateSearchState)

    updateSearchState()
  }

  const setupDashboard = (root) => {
    updateNavbarHeight()
    setupWelcome(root)
    setupSearch(root)

    window.addEventListener('resize', updateNavbarHeight)
  }

  document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('[data-dashboard-home]').forEach(setupDashboard)
  })
})()
