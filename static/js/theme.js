(() => {
  const storageKey = "cordell-theme";
  const root = document.documentElement;
  const toggleButtons = document.querySelectorAll("[data-theme-toggle]");

  function normalizeTheme(value) {
    if (value === "dark" || value === "light") {
      return value;
    }

    return "light";
  }

  function currentTheme() {
    return normalizeTheme(root.dataset.theme);
  }

  function applyTheme(theme) {
    const normalizedTheme = normalizeTheme(theme);

    root.dataset.theme = normalizedTheme;

    for (const button of toggleButtons) {
      const isDark = normalizedTheme === "dark";
      const label = isDark ? "Theme: Dark" : "Theme: Light";
      const nextTheme = isDark ? "light" : "dark";

      button.textContent = label;
      button.setAttribute("aria-pressed", isDark ? "true" : "false");
      button.setAttribute("aria-label", `Switch to ${nextTheme} theme`);
    }
  }

  function saveTheme(theme) {
    try {
      localStorage.setItem(storageKey, theme);
    } catch {
      // Ignore storage errors. Theme switching should still work in memory.
    }
  }

  function toggleTheme() {
    const nextTheme = currentTheme() === "dark" ? "light" : "dark";

    saveTheme(nextTheme);
    applyTheme(nextTheme);
  }

  for (const button of toggleButtons) {
    button.addEventListener("click", toggleTheme);
  }

  applyTheme(currentTheme());
})();