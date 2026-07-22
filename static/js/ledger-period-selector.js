(() => {
  function initLedgerPeriodSelector() {
    const scope = document.querySelector("[data-ledger-period-scope]");
    if (!scope) return;

    const yearSelect = scope.querySelector("[data-ledger-year-select]");
    const monthSelect = scope.querySelector("[data-ledger-month-select]");

    if (!yearSelect || !monthSelect) return;

    function updateMonthOptions() {
      const selectedYear = yearSelect.value;
      const isAllPeriods = selectedYear === "all";

      monthSelect.disabled = isAllPeriods;

      Array.from(monthSelect.options).forEach((option) => {
        const optionYear = option.dataset.year;

        if (option.value === "") {
          option.hidden = !isAllPeriods;
          return;
        }

        option.hidden = isAllPeriods || optionYear !== selectedYear;
      });

      if (isAllPeriods) {
        monthSelect.value = "";
        return;
      }

      const selectedOption = monthSelect.selectedOptions[0];
      if (selectedOption && !selectedOption.hidden) {
        return;
      }

      const firstVisibleOption = Array.from(monthSelect.options).find((option) => {
        return option.value !== "" && !option.hidden;
      });

      if (firstVisibleOption) {
        monthSelect.value = firstVisibleOption.value;
      }
    }

    yearSelect.addEventListener("change", updateMonthOptions);
    updateMonthOptions();
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initLedgerPeriodSelector);
  } else {
    initLedgerPeriodSelector();
  }
})();
