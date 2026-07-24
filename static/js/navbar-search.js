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

  const setupNavbarSearch = (root) => {
    const form = root.querySelector('[data-navbar-search-form]')
    const input = root.querySelector('[data-navbar-search-input]')
    const results = root.querySelector('[data-navbar-search-results]')

    if (!form || !input || !results) {
      return
    }

    let controller = null

    const closeResults = () => {
      results.innerHTML = ''
      results.hidden = true
    }

    const openResults = () => {
      results.hidden = false
    }

    const fetchSuggestions = async () => {
      const query = input.value.trim()

      if (query === '') {
        closeResults()
        return
      }

      if (controller) {
        controller.abort()
      }

      controller = new AbortController()

      const url = new URL('/search/suggestions', window.location.origin)
      url.searchParams.set('q', query)

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

        results.innerHTML = await response.text()
        openResults()
      } catch (error) {
        if (error.name === 'AbortError') {
          return
        }

        console.error(error)
        closeResults()
      }
    }

    const debouncedFetchSuggestions = debounce(fetchSuggestions, 200)

    input.addEventListener('input', () => {
      debouncedFetchSuggestions()
    })

    input.addEventListener('focus', () => {
      if (input.value.trim() !== '' && results.innerHTML.trim() !== '') {
        openResults()
      }
    })

    input.addEventListener('keydown', (event) => {
      if (event.key === 'Escape') {
        closeResults()
        input.blur()
      }
    })

    form.addEventListener('submit', (event) => {
      const query = input.value.trim()

      if (query === '') {
        event.preventDefault()
        closeResults()
      }
    })

    document.addEventListener('click', (event) => {
      if (!root.contains(event.target)) {
        closeResults()
      }
    })
  }

  document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('[data-navbar-search]').forEach(setupNavbarSearch)
  })
})()
