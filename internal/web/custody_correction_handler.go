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
	Title                       string
	Error                       string
	BaseTitle                   string
	BaseDescription             string
	EffectivePersonnelLabel     string
	EffectivePersonnelIsActive  bool
	HasInactiveEffectiveRecords bool
	InactiveEffectiveWarning    string
	Receipt                     custodyEditReceiptView
	CorrectionID                string
	FormAction                  string
	CorrectedPersonnelID        string
	CorrectedNotes              string
	PersonnelOptions            []correctionPersonnelOptionView
	AssetOptions                []custodyAssetOptionView
	LineRows                    []custodyLineFormRowView
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

type custodyTransactionEditFormState struct {
	CorrectionID         string
	CorrectedPersonnelID string
	CorrectedNotes       string
	LineRows             []custodyLineFormRowView
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
			"O ID da correção está ausente. Recarregue a página e tente novamente.",
		)
		return
	}

	correctedPersonnelID := domain.PersonnelID(strings.TrimSpace(r.FormValue("corrected_personnel_id")))
	correctedNotes := strings.TrimSpace(r.FormValue("corrected_notes"))

	lineRows := custodyLineFormRowsFromRequest(r)
	lines, parseErr := parseCustodyLineCommandsFromRequest(r)
	if parseErr != nil {
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
				Error:                humanizeCustodyCorrectionError(parseErr),
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
		StatusFilter: string(ports.RecordStatusFilterActive),
	})
	if err != nil {
		return custodyTransactionEditPageData{}, err
	}

	assets, err := s.services.ListAssets.Execute(r.Context(), app.ListAssetsCommand{
		Limit:        500,
		StatusFilter: string(ports.RecordStatusFilterActive),
	})
	if err != nil {
		return custodyTransactionEditPageData{}, err
	}

	baseTitle := "Editando transação original"
	baseDescription := "Esta transação não possui edição anterior. O formulário é baseado no recibo original."

	effectivePersonnelID := receipt.PersonnelID
	effectiveNotes := receipt.Notes
	effectiveLines := correctionLineRowsFromReceiptLines(receipt.Lines)

	if receipt.HasCorrection {
		baseTitle = "Editando edição mais recente"
		baseDescription = "Esta transação já possui uma edição. O formulário é baseado na correção mais recente, não diretamente na transação original."
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

	if state.CorrectedNotes == "" {
		state.CorrectedNotes = effectiveNotes
	}

	activePersonnelIDs := activePersonnelIDSet(personnel)
	activeAssetIDs := activeAssetIDSet(assets)
	effectivePersonnelLabel, effectivePersonnelIsActive := effectivePersonnelLabelFromReceipt(receipt)

	state.CorrectedPersonnelID = correctionSelectedPersonnelID(
		effectivePersonnelID,
		activePersonnelIDs,
		state.CorrectedPersonnelID,
	)
	state.LineRows = correctionLineRowsForForm(effectiveLines, activeAssetIDs, state.LineRows)

	hasInactiveEffectiveAsset := false
	for _, row := range state.LineRows {
		if row.NeedsReplacement {
			hasInactiveEffectiveAsset = true
			break
		}
	}

	hasInactiveEffectiveRecords := !effectivePersonnelIsActive || hasInactiveEffectiveAsset
	inactiveEffectiveWarning := ""

	switch {
	case !effectivePersonnelIsActive && hasInactiveEffectiveAsset:
		inactiveEffectiveWarning = "Esta transação referencia atualmente militar e materiais inativos. Escolha substitutos ativos para salvar esta edição ou reative primeiro os registros inativos."
	case !effectivePersonnelIsActive:
		inactiveEffectiveWarning = "Esta transação referencia atualmente um militar inativo. Escolha um militar ativo para salvar esta edição ou reative primeiro o militar atual."
	case hasInactiveEffectiveAsset:
		inactiveEffectiveWarning = "Esta transação referencia atualmente materiais inativos. Escolha materiais ativos substitutos para salvar esta edição ou reative primeiro os materiais inativos."
	}

	typeLabel := custodyTransactionTypeLabel(receipt.TransactionType)

	data := custodyTransactionEditPageData{
		privateLayoutData:           newPrivateLayoutData(r),
		Title:                       custodyTransactionTypeActionLabel(receipt.TransactionType),
		Error:                       state.Error,
		BaseTitle:                   baseTitle,
		BaseDescription:             baseDescription,
		EffectivePersonnelLabel:     effectivePersonnelLabel,
		EffectivePersonnelIsActive:  effectivePersonnelIsActive,
		HasInactiveEffectiveRecords: hasInactiveEffectiveRecords,
		InactiveEffectiveWarning:    inactiveEffectiveWarning,
		CorrectionID:                state.CorrectionID,
		FormAction:                  "/custody/transactions/" + string(transactionID) + "/corrections",
		CorrectedPersonnelID:        state.CorrectedPersonnelID,
		CorrectedNotes:              state.CorrectedNotes,
		Receipt: custodyEditReceiptView{
			ID:        string(receipt.ID),
			TypeLabel: typeLabel,
			EditLabel: custodyTransactionTypeActionLabel(receipt.TransactionType),
			CreatedAt: formatDateTime(receipt.CreatedAt),
		},
		PersonnelOptions: newCorrectionPersonnelOptions(personnel, state.CorrectedPersonnelID),
		AssetOptions:     newCustodyAssetOptions(assets),
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
		displayName := militaryDisplayName(item.Rank(), item.Alias())
		sectionShortLabel := item.Section().Abbreviation()

		options = append(options, correctionPersonnelOptionView{
			ID:       string(item.ID()),
			Label:    newPersonnelOptionLabel(displayName, item.FullName(), sectionShortLabel),
			Selected: string(item.ID()) == selectedID,
		})
	}

	return options
}

func activePersonnelIDSet(personnel []domain.Personnel) map[string]struct{} {
	ids := make(map[string]struct{}, len(personnel))

	for _, item := range personnel {
		ids[string(item.ID())] = struct{}{}
	}

	return ids
}

func activeAssetIDSet(assets []domain.Asset) map[string]struct{} {
	ids := make(map[string]struct{}, len(assets))

	for _, item := range assets {
		ids[string(item.ID())] = struct{}{}
	}

	return ids
}

func correctionSelectedPersonnelID(
	effectivePersonnelID domain.PersonnelID,
	activePersonnelIDs map[string]struct{},
	stateSelectedPersonnelID string,
) string {
	if strings.TrimSpace(stateSelectedPersonnelID) != "" {
		return stateSelectedPersonnelID
	}

	if _, ok := activePersonnelIDs[string(effectivePersonnelID)]; !ok {
		return ""
	}

	return string(effectivePersonnelID)
}

func effectivePersonnelLabelFromReceipt(receipt app.CustodyReceipt) (string, bool) {
	if receipt.HasCorrection {
		return militaryDisplayName(
				receipt.Correction.CorrectedPersonnelRank,
				receipt.Correction.CorrectedPersonnelAlias,
			) + " - " + receipt.Correction.CorrectedPersonnelFullName,
			receipt.Correction.CorrectedPersonnelActive
	}

	return militaryDisplayName(
			receipt.PersonnelRank,
			receipt.PersonnelAlias,
		) + " - " + receipt.PersonnelFullName,
		receipt.PersonnelActive
}

func correctionLineRowsFromReceiptLines(lines []app.CustodyReceiptLine) []custodyLineFormRowView {
	rows := make([]custodyLineFormRowView, 0, len(lines))

	for _, line := range lines {
		rows = append(rows, custodyLineFormRowView{
			AssetID:              string(line.AssetID),
			Quantity:             strconv.Itoa(line.Quantity),
			CurrentAssetLabel:    line.AssetName,
			CurrentAssetIsActive: line.AssetActive,
			NeedsReplacement:     !line.AssetActive,
		})
	}

	return rows
}

func correctionLineRowsFromCorrectionLines(lines []app.CustodyCorrectionContextLine) []custodyLineFormRowView {
	rows := make([]custodyLineFormRowView, 0, len(lines))

	for _, line := range lines {
		rows = append(rows, custodyLineFormRowView{
			AssetID:              string(line.AssetID),
			Quantity:             strconv.Itoa(line.Quantity),
			CurrentAssetLabel:    line.AssetName,
			CurrentAssetIsActive: line.AssetActive,
			NeedsReplacement:     !line.AssetActive,
		})
	}

	return rows
}

func correctionLineRowsForForm(
	effectiveLines []custodyLineFormRowView,
	activeAssetIDs map[string]struct{},
	stateRows []custodyLineFormRowView,
) []custodyLineFormRowView {
	if len(stateRows) > 0 {
		return ensureAtLeastOneCustodyLineFormRow(stateRows)
	}

	rows := make([]custodyLineFormRowView, 0, len(effectiveLines))

	for _, line := range effectiveLines {
		_, isActive := activeAssetIDs[line.AssetID]

		row := line
		row.CurrentAssetIsActive = isActive

		if !isActive {
			row.AssetID = ""
			row.NeedsReplacement = true
		}

		if row.Quantity == "" {
			row.Quantity = "1"
		}

		rows = append(rows, row)
	}

	return ensureAtLeastOneCustodyLineFormRow(rows)
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
		return "A transação, militar, material ou operador não foi encontrado."
	case errors.Is(err, domain.ErrInactiveOperator):
		return "Operadores inativos não podem registrar correções de custódia."
	case errors.Is(err, domain.ErrInactivePersonnel):
		return "Esta edição atribuiria custódia de cautela a um militar inativo. Reative o militar primeiro ou escolha um militar ativo."
	case errors.Is(err, domain.ErrInactiveAsset):
		return "Esta edição atribuiria custódia de cautela a um material inativo. Reative o material primeiro ou escolha um material ativo."
	case errors.Is(err, domain.ErrInsufficientCustodyBalance):
		return "Esta edição não pode ser aplicada porque deixaria um saldo de custódia negativo. Isso pode acontecer quando uma atividade posterior já consumiu parte do saldo afetado pela edição."
	case errors.Is(err, domain.ErrEmptyPersonnelID):
		return "O militar corrigido é obrigatório."
	case errors.Is(err, domain.ErrEmptyAssetID):
		return "Cada linha de correção deve incluir um material."
	case errors.Is(err, errNoCustodyLineSubmitted):
		return "Pelo menos uma linha de correção é obrigatória."
	case errors.Is(err, domain.ErrEmptyTransactionLines):
		return "Pelo menos uma linha de correção é obrigatória."
	case errors.Is(err, domain.ErrInvalidQuantity):
		return "Cada quantidade da correção deve ser positiva."
	default:
		return "Não foi possível salvar esta edição. Revise o militar corrigido, materiais, quantidades e estado atual da custódia."
	}
}

func custodyCorrectionAuditMetadata(result app.RegisterCustodyCorrectionResult) map[string]string {
	return map[string]string{
		"correction_id":            string(result.Correction.ID()),
		"corrected_transaction_id": string(result.Correction.CorrectedTransactionID()),
		"corrected_personnel_id":   string(result.Correction.CorrectedPersonnelID()),
	}
}
