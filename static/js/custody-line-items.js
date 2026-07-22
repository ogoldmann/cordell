(() => {
  function initCustodyLineItems() {
    document.querySelectorAll("[data-custody-line-items]").forEach((scope) => {
      const container = scope.querySelector("[data-custody-line-container]");
      const template = scope.querySelector("[data-custody-line-template]");
      const addButton = scope.querySelector("[data-add-custody-line]");

      if (!container || !template || !addButton) return;

      function rows() {
        return Array.from(container.querySelectorAll("[data-custody-line-row]"));
      }

      function selects() {
        return Array.from(container.querySelectorAll("[data-custody-asset-select]"));
      }

      function selectedAssetIDsExcept(currentSelect) {
        return selects()
          .filter((select) => select !== currentSelect)
          .map((select) => select.value)
          .filter((value) => value !== "");
      }

      function updateAssetOptions() {
        selects().forEach((select) => {
          const selectedElsewhere = new Set(selectedAssetIDsExcept(select));

          Array.from(select.options).forEach((option) => {
            if (option.value === "") {
              option.hidden = false;
              option.disabled = false;
              return;
            }

            const unavailable = selectedElsewhere.has(option.value);

            option.hidden = unavailable;
            option.disabled = unavailable;
          });
        });
      }

      function updateRemoveButtons() {
        const currentRows = rows();

        currentRows.forEach((row) => {
          const removeButton = row.querySelector("[data-remove-custody-line]");
          if (!removeButton) return;

          const disabled = currentRows.length <= 1;
          removeButton.disabled = disabled;
          removeButton.classList.toggle("opacity-50", disabled);
          removeButton.classList.toggle("cursor-not-allowed", disabled);
        });
      }

      function syncRows() {
        updateRemoveButtons();
        updateAssetOptions();
      }

      function focusNewRow(row) {
        const select = row.querySelector("[data-custody-asset-select]");
        if (select) select.focus();
      }

      addButton.addEventListener("click", () => {
        const fragment = template.content.cloneNode(true);
        const newRow = fragment.querySelector("[data-custody-line-row]");

        container.appendChild(fragment);
        syncRows();

        if (newRow) {
          focusNewRow(newRow);
        }
      });

      container.addEventListener("click", (event) => {
        const removeButton = event.target.closest("[data-remove-custody-line]");
        if (!removeButton) return;

        const currentRows = rows();
        if (currentRows.length <= 1) return;

        const row = removeButton.closest("[data-custody-line-row]");
        if (!row) return;

        row.remove();
        syncRows();
      });

      container.addEventListener("change", (event) => {
        if (!event.target.closest("[data-custody-asset-select]")) return;

        updateAssetOptions();
      });

      syncRows();
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initCustodyLineItems);
  } else {
    initCustodyLineItems();
  }
})();
