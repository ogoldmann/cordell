(() => {
  const storageKey = 'cordell-theme'
  const allowedThemes = ['light', 'dark', 'sepia']

  const isAllowedTheme = (theme) => {
    return allowedThemes.includes(theme)
  }

  const getStoredTheme = () => {
    const storedTheme = window.localStorage.getItem(storageKey)

    if (isAllowedTheme(storedTheme)) {
      return storedTheme
    }

    return null
  }

  const getInitialTheme = () => {
    return getStoredTheme() || 'dark'
  }

  const applyTheme = (theme) => {
    const nextTheme = isAllowedTheme(theme) ? theme : 'dark'

    document.documentElement.dataset.theme = nextTheme
    document.documentElement.style.colorScheme = nextTheme === 'dark' ? 'dark' : 'light'

    window.localStorage.setItem(storageKey, nextTheme)

    document.dispatchEvent(
      new CustomEvent('cordell:theme-changed', {
        detail: {
          theme: nextTheme,
        },
      }),
    )
  }

  const updateThemeControls = (theme) => {
    document.querySelectorAll('[data-theme-option]').forEach((control) => {
      const isSelected = control.dataset.themeOption === theme

      control.dataset.selected = String(isSelected)
      control.setAttribute('aria-pressed', String(isSelected))
    })
  }

  const setupThemeControls = () => {
    document.querySelectorAll('[data-theme-option]').forEach((control) => {
      control.addEventListener('click', () => {
        applyTheme(control.dataset.themeOption)
        updateThemeControls(document.documentElement.dataset.theme)
      })
    })
  }

  const initialTheme = getInitialTheme()

  applyTheme(initialTheme)

  document.addEventListener('DOMContentLoaded', () => {
    updateThemeControls(document.documentElement.dataset.theme)
    setupThemeControls()
  })

  document.addEventListener('cordell:theme-changed', (event) => {
    updateThemeControls(event.detail.theme)
  })

  window.CordellTheme = {
    apply: applyTheme,
    current: () => document.documentElement.dataset.theme,
  }
})()
