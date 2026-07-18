(() => {
  function setButtonStyle(button, style) {
    button.className = "inline-flex w-full items-center justify-center rounded-xl px-5 py-3 font-bold transition sm:w-auto";

    if (style === "warning") {
      button.className += " border border-amber-300 bg-amber-50 text-amber-800 hover:border-amber-400 dark:border-amber-900/60 dark:bg-amber-950/40 dark:text-amber-200";
      return;
    }

    button.className += " bg-cordell-primary text-white hover:bg-cordell-primary-hover dark:bg-cordell-dark-primary dark:text-cordell-dark-background";
  }

  function initConfirmModal() {
    const modal = document.querySelector("[data-confirm-modal]");
    if (!modal) return;

    const modalForm = modal.querySelector("[data-confirm-modal-form]");
    const modalTitle = modal.querySelector("[data-confirm-modal-title]");
    const modalDescription = modal.querySelector("[data-confirm-modal-description]");
    const modalWarning = modal.querySelector("[data-confirm-modal-warning]");
    const modalButton = modal.querySelector("[data-confirm-modal-button]");

    if (!modalForm || !modalTitle || !modalDescription || !modalWarning || !modalButton) return;

    document.querySelectorAll("[data-confirm-trigger]").forEach((triggerForm) => {
      triggerForm.addEventListener("submit", (event) => {
        event.preventDefault();

        modalForm.action = triggerForm.action;
        modalTitle.textContent = triggerForm.dataset.confirmTitle || "Confirm action";
        modalDescription.textContent = triggerForm.dataset.confirmDescription || "";
        modalButton.textContent = triggerForm.dataset.confirmLabel || "Confirm";

        const warning = triggerForm.dataset.confirmWarning || "";
        if (warning.trim() === "") {
          modalWarning.hidden = true;
          modalWarning.textContent = "";
        } else {
          modalWarning.hidden = false;
          modalWarning.textContent = warning;
        }

        setButtonStyle(modalButton, triggerForm.dataset.confirmStyle || "primary");

        modal.showModal();
      });
    });

    modal.addEventListener("click", (event) => {
      if (event.target === modal) {
        modal.close();
      }
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initConfirmModal);
  } else {
    initConfirmModal();
  }
})();
