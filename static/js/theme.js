(() => {
  const storageKey = "cordell-theme";
  const root = document.documentElement;
  const toggleButtons = document.querySelectorAll("[data-theme-toggle]");
  const allowedThemes = ["light", "dark", "sepia"];
  const themeLabels = {
    light: "Tema: Claro",
    dark: "Tema: Escuro",
    sepia: "Tema: Sépia",
  };
  const themeNames = {
    light: "claro",
    dark: "escuro",
    sepia: "sépia",
  };

  function normalizeTheme(value) {
    if (allowedThemes.includes(value)) {
      return value;
    }

    return "light";
  }

  function currentTheme() {
    return normalizeTheme(root.dataset.theme);
  }

  function nextTheme(theme) {
    const currentIndex = allowedThemes.indexOf(normalizeTheme(theme));

    return allowedThemes[(currentIndex + 1) % allowedThemes.length];
  }

  function applyTheme(theme) {
    const normalizedTheme = normalizeTheme(theme);
    const nextNormalizedTheme = nextTheme(normalizedTheme);

    root.dataset.theme = normalizedTheme;

    for (const button of toggleButtons) {
      button.textContent = themeLabels[normalizedTheme];
      button.setAttribute("aria-pressed", normalizedTheme === "light" ? "false" : "true");
      button.setAttribute("aria-label", `Mudar para tema ${themeNames[nextNormalizedTheme]}`);
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
    const nextNormalizedTheme = nextTheme(currentTheme());

    saveTheme(nextNormalizedTheme);
    applyTheme(nextNormalizedTheme);
  }

  for (const button of toggleButtons) {
    button.addEventListener("click", toggleTheme);
  }

  applyTheme(currentTheme());
})();
