package web

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/ports"

	"github.com/go-chi/chi/v5"
)

type custodyTransactionEditPageData struct {
	privateLayoutData
	Title                string
	Error                string
	BaseTitle            string
	BaseDescription      string
	Receipt              custodyEditReceiptView
	CorrectionID         string
	FormAction           string
	CorrectedPersonnelID string
	CorrectedNotes       string
	PersonnelOptions     []correctionPersonnelOptionView
	AssetOptions         []correctionAssetOptionView
	LineRows             []correctionLineRowView
}

type custodyEditReceiptView struct {
	ID        string
	TypeLabel string
	EditLabel string
	CreatedAt string
}

type correctionPersonnelOptionView struct {
	ID       string
	Label    string
	Selected bool
}

type correctionAssetOptionView struct {
	ID    string
	Label string
}

type correctionLineRowView struct {
	AssetID  string
	Quantity string
}

type custodyTransactionEditFormState struct {
	CorrectionID         string
	CorrectedPersonnelID string
	CorrectedNotes       string
	LineRows             []correctionLineRowView
	Error                string
}

func (s *Server) handleEditCustodyTransactionForm(w http.ResponseWriter, r *http.Request) {
	transactionID := domain.CustodyTransactionID(chi.URLParam(r, "id"))

	correctionID, err := newFormCorrectionID()
	if err != nil {
		s.logger.Error("failed to generate correction id", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data, err := s.newCustodyTransactionEditPageData(r, transactionID, custodyTransactionEditFormState{
		CorrectionID: string(correctionID),
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyTransactionID) {
			http.NotFound(w, r)
			return
		}

		s.logger.Error("failed to build custody transaction edit form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := s.renderer.Render(w, http.StatusOK, "custody_transaction_edit.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreateCustodyCorrection(w http.ResponseWriter, r *http.Request) {
	transactionID := domain.CustodyTransactionID(chi.URLParam(r, "id"))

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	correctionID := domain.CustodyCorrectionID(strings.TrimSpace(r.FormValue("correction_id")))
	if correctionID == "" {
		s.renderCustodyCorrectionFormError(
			w,
			r,
			http.StatusBadRequest,
			transactionID,
			"Form correction ID is missing. Please reload the page and try again.",
		)
		return
	}

	correctedPersonnelID := domain.PersonnelID(strings.TrimSpace(r.FormValue("corrected_personnel_id")))
	correctedNotes := strings.TrimSpace(r.FormValue("corrected_notes"))

	lineRows := correctionLineRowsFromRequest(r)
	lines, parseErr := parseCorrectionLineCommands(lineRows)
	if parseErr != "" {
		s.renderCustodyCorrectionFormErrorWithState(
			w,
			r,
			http.StatusBadRequest,
			transactionID,
			custodyTransactionEditFormState{
				CorrectionID:         string(correctionID),
				CorrectedPersonnelID: string(correctedPersonnelID),
				CorrectedNotes:       correctedNotes,
				LineRows:             lineRows,
				Error:                parseErr,
			},
		)
		return
	}

	currentOperator, ok := currentOperatorFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	result, err := s.services.RegisterCustodyCorrection.Execute(r.Context(), app.RegisterCustodyCorrectionCommand{
		CorrectionID:           correctionID,
		CorrectedTransactionID: transactionID,
		OperatorID:             currentOperator.ID(),
		CorrectedPersonnelID:   correctedPersonnelID,
		Lines:                  lines,
		CorrectedNotes:         correctedNotes,
	})
	if err != nil {
		s.renderCustodyCorrectionFormErrorWithState(
			w,
			r,
			http.StatusBadRequest,
			transactionID,
			custodyTransactionEditFormState{
				CorrectionID:         string(correctionID),
				CorrectedPersonnelID: string(correctedPersonnelID),
				CorrectedNotes:       correctedNotes,
				LineRows:             lineRows,
				Error:                humanizeCustodyCorrectionError(err),
			},
		)
		return
	}

	if result.Created {
		s.recordAuditEventOrLog(
			r,
			domain.AuditEventCustodyCorrectionCreated,
			domain.AuditEntityCustodyTransaction,
			string(transactionID),
			custodyCorrectionAuditMetadata(result),
		)
	}

	http.Redirect(w, r, "/custody/transactions/"+string(transactionID), http.StatusSeeOther)
}

func (s *Server) newCustodyTransactionEditPageData(
	r *http.Request,
	transactionID domain.CustodyTransactionID,
	state custodyTransactionEditFormState,
) (custodyTransactionEditPageData, error) {
	receipt, err := s.services.GetCustodyReceipt.Execute(r.Context(), app.GetCustodyReceiptCommand{
		ID: transactionID,
	})
	if err != nil {
		return custodyTransactionEditPageData{}, err
	}

	personnel, err := s.services.ListPersonnel.Execute(r.Context(), app.ListPersonnelCommand{
		Limit:        500,
		StatusFilter: string(ports.RecordStatusFilterAll),
	})
	if err != nil {
		return custodyTransactionEditPageData{}, err
	}

	assets, err := s.services.ListAssets.Execute(r.Context(), app.ListAssetsCommand{
		Limit:        500,
		StatusFilter: string(ports.RecordStatusFilterAll),
	})
	if err != nil {
		return custodyTransactionEditPageData{}, err
	}

	baseTitle := "Editing original transaction"
	baseDescription := "This transaction has no previous edit. The form is based on the original receipt."

	effectivePersonnelID := receipt.PersonnelID
	effectiveNotes := receipt.Notes
	effectiveLines := correctionLineRowsFromReceiptLines(receipt.Lines)

	if receipt.HasCorrection {
		baseTitle = "Editing latest edit"
		baseDescription = "This transaction already has an edit. The form is based on the latest correction, not directly on the original transaction."
		effectivePersonnelID = receipt.Correction.CorrectedPersonnelID
		effectiveNotes = receipt.Correction.CorrectedNotes
		effectiveLines = correctionLineRowsFromCorrectionLines(receipt.Correction.Lines)
	}

	if state.CorrectionID == "" {
		correctionID, err := newFormCorrectionID()
		if err != nil {
			return custodyTransactionEditPageData{}, err
		}
		state.CorrectionID = string(correctionID)
	}

	if state.CorrectedPersonnelID == "" {
		state.CorrectedPersonnelID = string(effectivePersonnelID)
	}

	if state.CorrectedNotes == "" {
		state.CorrectedNotes = effectiveNotes
	}

	if len(state.LineRows) == 0 {
		state.LineRows = ensureAtLeastOneCorrectionRow(effectiveLines)
	}

	typeLabel := custodyTransactionTypeLabel(receipt.TransactionType)

	data := custodyTransactionEditPageData{
		privateLayoutData:    newPrivateLayoutData(r),
		Title:                "Edit " + strings.ToLower(typeLabel),
		Error:                state.Error,
		BaseTitle:            baseTitle,
		BaseDescription:      baseDescription,
		CorrectionID:         state.CorrectionID,
		FormAction:           "/custody/transactions/" + string(transactionID) + "/corrections",
		CorrectedPersonnelID: state.CorrectedPersonnelID,
		CorrectedNotes:       state.CorrectedNotes,
		Receipt: custodyEditReceiptView{
			ID:        string(receipt.ID),
			TypeLabel: typeLabel,
			EditLabel: "Edit " + strings.ToLower(typeLabel),
			CreatedAt: formatDateTime(receipt.CreatedAt),
		},
		PersonnelOptions: newCorrectionPersonnelOptions(personnel, state.CorrectedPersonnelID),
		AssetOptions:     newCorrectionAssetOptions(assets),
		LineRows:         state.LineRows,
	}

	return data, nil
}

func newCorrectionPersonnelOptions(
	personnel []domain.Personnel,
	selectedID string,
) []correctionPersonnelOptionView {
	options := make([]correctionPersonnelOptionView, 0, len(personnel))

	for _, item := range personnel {
		label := militaryDisplayName(item.Rank(), item.Alias()) + " - " + item.FullName()
		if !item.Active() {
			label += " (Inactive)"
		}

		options = append(options, correctionPersonnelOptionView{
			ID:       string(item.ID()),
			Label:    label,
			Selected: string(item.ID()) == selectedID,
		})
	}

	return options
}

func newCorrectionAssetOptions(assets []domain.Asset) []correctionAssetOptionView {
	options := make([]correctionAssetOptionView, 0, len(assets))

	for _, item := range assets {
		label := item.Name()
		if !item.Active() {
			label += " (Inactive)"
		}

		options = append(options, correctionAssetOptionView{
			ID:    string(item.ID()),
			Label: label,
		})
	}

	return options
}

func correctionLineRowsFromReceiptLines(lines []app.CustodyReceiptLine) []correctionLineRowView {
	rows := make([]correctionLineRowView, 0, len(lines))

	for _, line := range lines {
		rows = append(rows, correctionLineRowView{
			AssetID:  string(line.AssetID),
			Quantity: strconv.Itoa(line.Quantity),
		})
	}

	return rows
}

func correctionLineRowsFromCorrectionLines(lines []app.CustodyCorrectionContextLine) []correctionLineRowView {
	rows := make([]correctionLineRowView, 0, len(lines))

	for _, line := range lines {
		rows = append(rows, correctionLineRowView{
			AssetID:  string(line.AssetID),
			Quantity: strconv.Itoa(line.Quantity),
		})
	}

	return rows
}

func ensureAtLeastOneCorrectionRow(rows []correctionLineRowView) []correctionLineRowView {
	if len(rows) > 0 {
		return rows
	}

	return []correctionLineRowView{{}}
}

func correctionLineRowsFromRequest(r *http.Request) []correctionLineRowView {
	assetIDs := r.PostForm["asset_id"]
	quantities := r.PostForm["quantity"]

	rowCount := len(assetIDs)
	if len(quantities) > rowCount {
		rowCount = len(quantities)
	}

	rows := make([]correctionLineRowView, 0, rowCount)

	for index := 0; index < rowCount; index++ {
		assetID := ""
		if index < len(assetIDs) {
			assetID = strings.TrimSpace(assetIDs[index])
		}

		quantity := ""
		if index < len(quantities) {
			quantity = strings.TrimSpace(quantities[index])
		}

		rows = append(rows, correctionLineRowView{
			AssetID:  assetID,
			Quantity: quantity,
		})
	}

	return rows
}

func parseCorrectionLineCommands(rows []correctionLineRowView) ([]app.CustodyLineCommand, string) {
	lines := make([]app.CustodyLineCommand, 0, len(rows))

	for _, row := range rows {
		assetID := strings.TrimSpace(row.AssetID)
		quantityText := strings.TrimSpace(row.Quantity)

		if assetID == "" && quantityText == "" {
			continue
		}

		if assetID == "" || quantityText == "" {
			return nil, "Each correction line must include both asset and quantity."
		}

		quantity, err := strconv.Atoi(quantityText)
		if err != nil || quantity <= 0 {
			return nil, "Each correction quantity must be a positive number."
		}

		lines = append(lines, app.CustodyLineCommand{
			AssetID:  domain.AssetID(assetID),
			Quantity: quantity,
		})
	}

	if len(lines) == 0 {
		return nil, "At least one correction line is required."
	}

	return lines, ""
}

func (s *Server) renderCustodyCorrectionFormError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	transactionID domain.CustodyTransactionID,
	message string,
) {
	s.renderCustodyCorrectionFormErrorWithState(w, r, status, transactionID, custodyTransactionEditFormState{
		Error: message,
	})
}

func (s *Server) renderCustodyCorrectionFormErrorWithState(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	transactionID domain.CustodyTransactionID,
	state custodyTransactionEditFormState,
) {
	data, err := s.newCustodyTransactionEditPageData(r, transactionID, state)
	if err != nil {
		s.logger.Error("failed to render custody correction form error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := s.renderer.Render(w, status, "custody_transaction_edit.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func humanizeCustodyCorrectionError(err error) string {
	switch {
	case errors.Is(err, ports.ErrNotFound):
		return "The transaction, personnel, asset, or operator could not be found."
	case errors.Is(err, domain.ErrInactiveOperator):
		return "Inactive operators cannot register custody corrections."
	case errors.Is(err, domain.ErrInactivePersonnel):
		return "This edit would assign checkout custody to an inactive personnel. Reactivate the personnel first or choose an active personnel."
	case errors.Is(err, domain.ErrInactiveAsset):
		return "This edit would assign checkout custody to an inactive asset. Reactivate the asset first or choose an active asset."
	case errors.Is(err, domain.ErrInsufficientCustodyBalance):
		return "This edit cannot be applied because it would make a custody balance negative. This can happen when later custody activity already consumed part of the balance affected by the edit."
	case errors.Is(err, domain.ErrEmptyPersonnelID):
		return "Corrected personnel is required."
	case errors.Is(err, domain.ErrEmptyAssetID):
		return "Each correction line must include an asset."
	case errors.Is(err, domain.ErrEmptyTransactionLines):
		return "At least one correction line is required."
	case errors.Is(err, domain.ErrInvalidQuantity):
		return "Each correction quantity must be a positive number."
	default:
		return "Could not save this edit. Please review the corrected personnel, assets, quantities, and current custody state."
	}
}

func custodyCorrectionAuditMetadata(result app.RegisterCustodyCorrectionResult) map[string]string {
	return map[string]string{
		"correction_id":            string(result.Correction.ID()),
		"corrected_transaction_id": string(result.Correction.CorrectedTransactionID()),
		"corrected_personnel_id":   string(result.Correction.CorrectedPersonnelID()),
	}
}
