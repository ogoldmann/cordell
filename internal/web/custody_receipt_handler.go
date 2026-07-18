package web

import (
	"errors"
	"net/http"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/ports"

	"github.com/go-chi/chi/v5"
)

type custodyReceiptPageData struct {
	privateLayoutData
	Title   string
	Receipt custodyReceiptView
}

func (s *Server) handleShowCustodyReceipt(w http.ResponseWriter, r *http.Request) {
	transactionID := domain.CustodyTransactionID(chi.URLParam(r, "id"))

	receipt, err := s.services.GetCustodyReceipt.Execute(r.Context(), app.GetCustodyReceiptCommand{
		ID: transactionID,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyTransactionID) {
			http.NotFound(w, r)
			return
		}

		s.logger.Error("failed to show custody receipt", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	receiptView := newCustodyReceiptView(receipt)

	data := custodyReceiptPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             receiptView.TypeLabel + " receipt",
		Receipt:           receiptView,
	}

	if err := s.renderer.Render(w, http.StatusOK, "custody_receipt_show.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}
