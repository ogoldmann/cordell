(() => {
  const setupOperatorMenu = (root) => {
    const button = root.querySelector('[data-operator-menu-button]')
    const dropdown = root.querySelector('[data-operator-menu-dropdown]')

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
    document.querySelectorAll('[data-operator-menu]').forEach(setupOperatorMenu)
  })
})()
