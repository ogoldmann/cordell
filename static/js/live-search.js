(() => {
  const debounce = (callback, delay) => {
    let timeoutID

    return (...args) => {
      window.clearTimeout(timeoutID)

      timeoutID = window.setTimeout(() => {
        callback(...args)
      }, delay)
    }
  }

  const buildURL = (form) => {
    const formData = new FormData(form)
    const params = new URLSearchParams()

    for (const [key, value] of formData.entries()) {
      const trimmedValue = String(value).trim()

      if (trimmedValue !== '') {
        params.set(key, trimmedValue)
      }
    }

    const action = form.getAttribute('action') || window.location.pathname
    const url = new URL(action, window.location.origin)

    url.search = params.toString()

    return url
  }

  const setStatus = (statusElement, message) => {
    if (!statusElement) {
      return
    }

    statusElement.textContent = message
  }

  const setupLiveSearchForm = (form) => {
    const targetSelector = form.dataset.liveSearchTarget
    const target = document.querySelector(targetSelector)

    if (!target) {
      return
    }

    const statusSelector = form.dataset.liveSearchStatus
    const statusElement = statusSelector ? document.querySelector(statusSelector) : null
    const delay = Number(form.dataset.liveSearchDelay || 250)

    let controller = null

    const runSearch = async () => {
      const url = buildURL(form)

      if (controller) {
        controller.abort()
      }

      controller = new AbortController()

      setStatus(statusElement, 'Pesquisando...')

      try {
        const response = await fetch(url.toString(), {
          method: 'GET',
          headers: {
            'X-Cordell-Partial': '1',
            'X-Requested-With': 'fetch',
          },
          signal: controller.signal,
        })

        if (!response.ok) {
          throw new Error(`Unexpected response: ${response.status}`)
        }

        const html = await response.text()

        target.innerHTML = html

        const visibleURL = new URL(url.toString())
        visibleURL.searchParams.delete('partial')

        window.history.replaceState({}, '', visibleURL.toString())

        setStatus(statusElement, 'Resultados atualizados.')
      } catch (error) {
        if (error.name === 'AbortError') {
          return
        }

        console.error(error)
        setStatus(statusElement, 'Não foi possível atualizar os resultados.')
      }
    }

    const debouncedRunSearch = debounce(runSearch, delay)

    form.addEventListener('input', (event) => {
      const targetElement = event.target

      if (!(targetElement instanceof HTMLInputElement || targetElement instanceof HTMLSelectElement)) {
        return
      }

      debouncedRunSearch()
    })

    form.addEventListener('submit', (event) => {
      event.preventDefault()
      runSearch()
    })
  }

  document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('[data-live-search-form]').forEach(setupLiveSearchForm)
  })
})()
