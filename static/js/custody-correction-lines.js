(() => {
  function initCorrectionLines() {
    const container = document.querySelector("[data-correction-lines]");
    const template = document.querySelector("[data-correction-line-template]");
    const addButton = document.querySelector("[data-add-correction-line]");

    if (!container || !template || !addButton) return;

    function rows() {
      return Array.from(container.querySelectorAll("[data-correction-line-row]"));
    }

    function updateRemoveButtons() {
      const currentRows = rows();

      currentRows.forEach((row) => {
        const removeButton = row.querySelector("[data-remove-correction-line]");
        if (!removeButton) return;

        const disabled = currentRows.length <= 1;
        removeButton.disabled = disabled;
        removeButton.classList.toggle("opacity-50", disabled);
        removeButton.classList.toggle("cursor-not-allowed", disabled);
      });
    }

    function focusNewRow(row) {
      const select = row.querySelector("select[name='asset_id']");
      if (select) select.focus();
    }

    addButton.addEventListener("click", () => {
      const fragment = template.content.cloneNode(true);
      const newRow = fragment.querySelector("[data-correction-line-row]");

      container.appendChild(fragment);

      updateRemoveButtons();

      if (newRow) {
        focusNewRow(newRow);
      }
    });

    container.addEventListener("click", (event) => {
      const removeButton = event.target.closest("[data-remove-correction-line]");
      if (!removeButton) return;

      const currentRows = rows();
      if (currentRows.length <= 1) return;

      const row = removeButton.closest("[data-correction-line-row]");
      if (!row) return;

      row.remove();
      updateRemoveButtons();
    });

    updateRemoveButtons();
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initCorrectionLines);
  } else {
    initCorrectionLines();
  }
})();
