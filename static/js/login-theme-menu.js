(() => {
  const setupLoginThemeMenu = (root) => {
    const button = root.querySelector('[data-login-theme-menu-button]')
    const dropdown = root.querySelector('[data-login-theme-menu-dropdown]')

    if (!button || !dropdown) {
      return
    }

    const open = () => {
      dropdown.hidden = false
      button.setAttribute('aria-expanded', 'true')
      root.dataset.open = 'true'
    }

    const close = () => {
      dropdown.hidden = true
      button.setAttribute('aria-expanded', 'false')
      root.dataset.open = 'false'
    }

    const toggle = () => {
      if (dropdown.hidden) {
        open()
        return
      }

      close()
    }

    button.addEventListener('click', (event) => {
      event.stopPropagation()
      toggle()
    })

    dropdown.addEventListener('click', (event) => {
      event.stopPropagation()
    })

    document.addEventListener('click', () => {
      close()
    })

    document.addEventListener('keydown', (event) => {
      if (event.key === 'Escape') {
        close()
        button.focus()
      }
    })
  }

  document.addEventListener('DOMContentLoaded', () => {
    window.sessionStorage.removeItem('cordell-dashboard-welcome-seen')

    document.querySelectorAll('[data-login-theme-menu]').forEach(setupLoginThemeMenu)
  })
})()
