package web

import (
	"errors"
	"net/http"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/ports"
)

type checkoutNewPageData struct {
	privateLayoutData
	Title               string
	Error               string
	Personnel           []personnelView
	SelectedPersonnelID string
	AssetOptions        []custodyAssetOptionView
	LineRows            []custodyLineFormRowView
	Notes               string
	TransactionID       string
}

func (s *Server) handleNewCheckoutForm(w http.ResponseWriter, r *http.Request) {
	lineRows := defaultCustodyLineFormRows()
	if assetID := r.URL.Query().Get("asset_id"); assetID != "" {
		lineRows = []custodyLineFormRowView{
			{
				AssetID:  assetID,
				Quantity: "1",
			},
		}
	}

	data, err := s.buildCheckoutNewPageData(r, checkoutNewPageData{
		Title:               checkoutLabel(),
		SelectedPersonnelID: r.URL.Query().Get("personnel_id"),
		LineRows:            lineRows,
	})
	if err != nil {
		s.logger.Error("failed to build checkout form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := s.renderer.Render(w, http.StatusOK, "checkout_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreateCheckout(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderCheckoutFormWithError(
			w,
			r,
			http.StatusBadRequest,
			"Envio de formulário inválido.",
		)
		return
	}

	personnelID := r.FormValue("personnel_id")
	notes := r.FormValue("notes")
	transactionID := domain.CustodyTransactionID(r.FormValue("transaction_id"))
	if transactionID == "" {
		s.renderCheckoutFormWithError(
			w,
			r,
			http.StatusBadRequest,
			"O ID da transação está ausente. Recarregue a página e tente novamente.",
		)
		return
	}

	lines, err := parseCustodyLineCommandsFromRequest(r)
	if err != nil {
		s.renderCheckoutFormWithError(
			w,
			r,
			http.StatusBadRequest,
			humanizeCustodyLineFormError(err),
		)
		return
	}

	currentOperator, ok := currentOperatorFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	result, err := s.services.RegisterCheckout.Execute(r.Context(), app.RegisterCheckoutCommand{
		TransactionID: transactionID,
		PersonnelID:   domain.PersonnelID(personnelID),
		OperatorID:    currentOperator.ID(),
		Lines:         lines,
		Notes:         notes,
	})
	if err != nil {
		s.renderCheckoutFormWithError(
			w,
			r,
			http.StatusBadRequest,
			humanizeCheckoutError(err),
		)
		return
	}

	if result.Created {
		s.recordAuditEventOrLog(
			r,
			domain.AuditEventCustodyCheckoutCreated,
			domain.AuditEntityCustodyTransaction,
			string(result.Transaction.ID()),
			map[string]string{
				"personnel_id": string(result.Transaction.PersonnelID()),
			},
		)
	}

	http.Redirect(
		w,
		r,
		"/custody/transactions/"+string(result.Transaction.ID()),
		http.StatusSeeOther,
	)
}

func (s *Server) renderCheckoutFormWithError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message string,
) {
	data, err := s.buildCheckoutNewPageData(r, checkoutNewPageData{
		Title:               checkoutLabel(),
		Error:               message,
		SelectedPersonnelID: r.FormValue("personnel_id"),
		LineRows:            custodyLineFormRowsFromRequest(r),
		Notes:               r.FormValue("notes"),
		TransactionID:       r.FormValue("transaction_id"),
	})
	if err != nil {
		s.logger.Error("failed to rebuild checkout form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := s.renderer.Render(w, status, "checkout_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) buildCheckoutNewPageData(
	r *http.Request,
	data checkoutNewPageData,
) (checkoutNewPageData, error) {
	data.privateLayoutData = newPrivateLayoutData(r)
	if data.TransactionID == "" {
		transactionID, err := newFormTransactionID()
		if err != nil {
			return checkoutNewPageData{}, err
		}

		data.TransactionID = string(transactionID)
	}

	personnel, err := s.services.ListPersonnel.Execute(r.Context(), app.ListPersonnelCommand{
		Limit:        100,
		StatusFilter: string(ports.RecordStatusFilterActive),
	})
	if err != nil {
		return checkoutNewPageData{}, err
	}

	assets, err := s.services.ListAssets.Execute(r.Context(), app.ListAssetsCommand{
		Limit:        100,
		StatusFilter: string(ports.RecordStatusFilterActive),
	})
	if err != nil {
		return checkoutNewPageData{}, err
	}

	data.Personnel = make([]personnelView, 0, len(personnel))
	for _, item := range personnel {
		if !item.Active() {
			continue
		}

		data.Personnel = append(data.Personnel, newPersonnelView(item))
	}

	data.AssetOptions = newCustodyAssetOptions(assets)
	if len(data.LineRows) == 0 {
		data.LineRows = defaultCustodyLineFormRows()
	}

	return data, nil
}

func humanizeCheckoutError(err error) string {
	switch {
	case errors.Is(err, domain.ErrEmptyPersonnelID):
		return "Selecione um militar."
	case errors.Is(err, domain.ErrInactivePersonnel):
		return "Este militar está inativo e não pode receber nova cautela."
	case errors.Is(err, domain.ErrInactiveAsset):
		return "Um dos materiais selecionados está inativo e não pode ser cautelado."
	case errors.Is(err, domain.ErrEmptyAssetID):
		return "Selecione pelo menos um material."
	case errors.Is(err, domain.ErrInvalidQuantity):
		return "Informe uma quantidade válida."
	case errors.Is(err, ports.ErrNotFound):
		return "Militar ou material selecionado não foi encontrado."
	default:
		return "Não foi possível registrar a cautela."
	}
}
