(() => {
  function initSubmitConfirmModal() {
    const modal = document.querySelector("[data-submit-confirm-modal]");
    if (!modal) return;

    const title = modal.querySelector("[data-submit-confirm-title]");
    const description = modal.querySelector("[data-submit-confirm-description]");
    const warning = modal.querySelector("[data-submit-confirm-warning]");
    const confirmButton = modal.querySelector("[data-submit-confirm-confirm-button]");

    if (!title || !description || !warning || !confirmButton) return;

    let pendingForm = null;

    document.querySelectorAll("[data-submit-confirm-trigger]").forEach((form) => {
      form.addEventListener("submit", (event) => {
        if (form.dataset.submitConfirmConfirmed === "true") {
          return;
        }

        event.preventDefault();

        pendingForm = form;

        title.textContent = form.dataset.submitConfirmTitle || "Save edit?";
        description.textContent = form.dataset.submitConfirmDescription || "";

        const warningText = form.dataset.submitConfirmWarning || "";
        if (warningText.trim() === "") {
          warning.hidden = true;
          warning.textContent = "";
        } else {
          warning.hidden = false;
          warning.textContent = warningText;
        }

        confirmButton.textContent = form.dataset.submitConfirmLabel || "Save edit";

        modal.showModal();
      });
    });

    confirmButton.addEventListener("click", () => {
      if (!pendingForm) return;

      pendingForm.dataset.submitConfirmConfirmed = "true";
      modal.close();

      pendingForm.requestSubmit();
    });

    modal.addEventListener("close", () => {
      pendingForm = null;
    });

    modal.addEventListener("click", (event) => {
      if (event.target === modal) {
        modal.close();
      }
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initSubmitConfirmModal);
  } else {
    initSubmitConfirmModal();
  }
})();
